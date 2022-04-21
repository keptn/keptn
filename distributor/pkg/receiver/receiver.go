package receiver

import (
	"context"
	"errors"
	"fmt"

	"github.com/keptn/keptn/distributor/pkg/model"
	nats2 "github.com/keptn/keptn/distributor/pkg/natsconnection"
	"github.com/keptn/keptn/distributor/pkg/poller"
	"github.com/keptn/keptn/distributor/pkg/utils"
	"sync"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
)

// EventReceiver is responsible for receive and process events from Keptn
type EventReceiver interface {
	Start(ctx *utils.ExecutionContext)
}

// NATSEventReceiver receives events directly from the NATS broker and sends the cloud event to the
// the keptn service
type NATSEventReceiver struct {
	env                   config.EnvConfig
	eventSender           poller.EventSender
	eventMatcher          *utils.EventMatcher
	natsConnectionHandler *nats2.NatsConnectionHandler
	ceCache               *utils.Cache
	mutex                 *sync.Mutex
	currentSubscriptions  []models.EventSubscription
	pullSubscriptions     bool
}

func New(env config.EnvConfig, eventSender poller.EventSender, pullSubscriptions bool) *NATSEventReceiver {
	eventMatcher := utils.NewEventMatcherFromEnv(env)
	nch := nats2.NewNatsConnectionHandler(env.PubSubURL)

	return &NATSEventReceiver{
		env:                   env,
		eventSender:           eventSender,
		eventMatcher:          eventMatcher,
		ceCache:               utils.NewCache(),
		mutex:                 &sync.Mutex{},
		natsConnectionHandler: nch,
		pullSubscriptions:     pullSubscriptions,
	}
}

func (n *NATSEventReceiver) Start(ctx *utils.ExecutionContext) error {
	if n.env.PubSubRecipient == "" {
		return errors.New("could not start NatsEventReceiver: no pubsub recipient defined")
	}

	if err := n.natsConnectionHandler.Connect(); err != nil {
		return fmt.Errorf("could not Start NatsEventReceiver: %w", err)
	}
	n.natsConnectionHandler.MessageHandler = n.handleMessage
	err := n.natsConnectionHandler.QueueSubscribeToTopics(n.env.PubSubTopics(), n.env.PubSubGroup)
	if err != nil {
		return fmt.Errorf("could not subscribe to events: %w", err)
	}

	defer func() {
		ctx.Wg.Done()
		err := n.natsConnectionHandler.RemoveAllSubscriptions()
		if err != nil {
			logger.WithError(err).Error("Could not remove subscriptions.")
		}
		logger.Info("Terminating NATS event receiver")
	}()

	<-ctx.Done()
	return nil
}

func (n *NATSEventReceiver) UpdateSubscriptions(subscriptions []models.EventSubscription) {
	n.currentSubscriptions = subscriptions
	topics := []string{}
	for _, s := range subscriptions {
		topics = append(topics, s.Event)
	}
	err := n.natsConnectionHandler.QueueSubscribeToTopics(topics, n.env.PubSubGroup)
	if err != nil {
		logger.Errorf("Could not subscribe to topics %v: %v", topics, err)
	}
}

func (n *NATSEventReceiver) handleMessage(m *nats.Msg) {
	go func() {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		logger.Infof("Received a message for topic [%s]\n", m.Subject)

		// decode to cloudevent
		cloudEvent, err := utils.DecodeNATSMessage(m.Data)
		if err != nil {
			return
		}

		// decode to keptn event
		keptnEvent, err := v0_2_0.ToKeptnEvent(*cloudEvent)
		if err != nil {
			return
		}

		// determine subscription for the received message
		subscriptions := n.getSubscriptionsFromReceivedMessage(m, *cloudEvent)
		if len(subscriptions) > 0 {
			if err := n.sendEventForSubscriptions(subscriptions, keptnEvent); err != nil {
				logger.Errorf("Could not send cloud event: %v", err)
			}
		} else if !n.pullSubscriptions {
			// forward keptn event
			if err := n.sendEvent(keptnEvent, nil); err != nil {
				logger.Errorf("Could not send cloud event: %v", err)
			}
		}
	}()
}

func (n *NATSEventReceiver) sendEventForSubscriptions(subscriptions []models.EventSubscription, keptnEvent models.KeptnContextExtendedCE) error {
	for i, subscription := range subscriptions {
		// check if the event with the given ID has already been sent for the subscription
		if n.ceCache.Contains(subscription.ID, keptnEvent.ID) {
			// Skip this event as it has already been sent
			logger.Debugf("CloudEvent with ID %s has already been sent", keptnEvent.ID)
			continue
		}
		// add to CloudEvents cache
		n.ceCache.Add(subscription.ID, keptnEvent.ID)

		defer func() {
			// after some time, remove the cache entry
			go func() {
				<-time.After(10 * time.Second)
				n.ceCache.Remove(subscription.ID, keptnEvent.ID)
			}()
		}()

		// add subscription ID as additional information to the keptn event
		if err := keptnEvent.AddTemporaryData("distributor", model.AdditionalSubscriptionData{SubscriptionID: subscription.ID}, models.AddTemporaryDataOptions{OverwriteIfExisting: true}); err != nil {
			logger.Errorf("Could not add temporary information about subscriptions to event: %v", err)
		}
		// forward keptn event
		if err := n.sendEvent(keptnEvent, &subscriptions[i]); err != nil {
			logger.Errorf("Could not send event for subscription %s: %v", subscription.ID, err)
		}
	}
	return nil
}

func (n *NATSEventReceiver) getSubscriptionsFromReceivedMessage(m *nats.Msg, event cloudevents.Event) []models.EventSubscription {
	subscriptionsForTopic := []models.EventSubscription{}
	for _, subscription := range n.currentSubscriptions {
		if subscription.Event == m.Sub.Subject { // need to check against the name of the subscription because this can be a wildcard as well
			matcher := utils.NewEventMatcherFromSubscription(subscription)
			if matcher.Matches(event) {
				subscriptionsForTopic = append(subscriptionsForTopic, subscription)
			}
		}
	}
	return subscriptionsForTopic
}

func (n *NATSEventReceiver) sendEvent(e models.KeptnContextExtendedCE, subscription *models.EventSubscription) error {
	event := v0_2_0.ToCloudEvent(e)
	if subscription != nil {
		matcher := utils.NewEventMatcherFromSubscription(*subscription)
		if !matcher.Matches(event) {
			return nil
		}
	} else if !n.eventMatcher.Matches(event) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	logger.Infof("Sending CloudEvent with ID %s to %s", event.ID(), n.env.PubSubRecipient)
	if err := n.eventSender.Send(ctx, event); err != nil {
		return err
	}
	return nil
}
