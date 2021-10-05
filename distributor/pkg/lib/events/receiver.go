package events

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

// EventReceiver is responsible for receive and process events from Keptn
type EventReceiver interface {
	Start(ctx *ExecutionContext)
}

// NATSEventReceiver receives events directly from the NATS broker and sends the cloud event to the
// the keptn service
type NATSEventReceiver struct {
	env                   config.EnvConfig
	eventSender           EventSender
	closeChan             chan bool
	eventMatcher          *EventMatcher
	natsConnectionHandler *NatsConnectionHandler
	ceCache               *CloudEventsCache
	mutex                 *sync.Mutex
	currentSubscriptions  []models.EventSubscription
}

func NewNATSEventReceiver(env config.EnvConfig, eventSender EventSender) *NATSEventReceiver {
	eventMatcher := NewEventMatcherFromEnv(env)
	nch := NewNatsConnectionHandler(env.PubSubURL)

	return &NATSEventReceiver{
		env:                   env,
		eventSender:           eventSender,
		closeChan:             make(chan bool),
		eventMatcher:          eventMatcher,
		ceCache:               NewCloudEventsCache(),
		mutex:                 &sync.Mutex{},
		natsConnectionHandler: nch,
	}
}

func (n *NATSEventReceiver) Start(ctx *ExecutionContext) {
	if n.env.PubSubRecipient == "" {
		logger.Warn("No pubsub recipient defined")
		return
	}
	n.natsConnectionHandler.messageHandler = n.handleMessage
	err := n.natsConnectionHandler.QueueSubscribeToTopics(n.env.GetPubSubTopics(), n.env.PubSubGroup)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		n.natsConnectionHandler.RemoveAllSubscriptions()
		logger.Info("Disconnected from NATS")
	}()

	for {
		select {
		case <-n.closeChan:
			return
		case <-ctx.Done():
			logger.Info("Terminating NATS event receiver")
			ctx.Wg.Done()
			return
		}
	}
}

func (n *NATSEventReceiver) UpdateSubscriptions(subscriptions []models.EventSubscription) {
	n.currentSubscriptions = subscriptions
	var topics []string
	for _, s := range subscriptions {
		topics = append(topics, s.Event)
	}
	err := n.natsConnectionHandler.QueueSubscribeToTopics(topics, n.env.PubSubGroup)
	if err != nil {
		logger.Errorf("Unable to subscribe to topics %v", topics)
	}
}

func (n *NATSEventReceiver) handleMessage(m *nats.Msg) {
	go func() {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		logger.Infof("Received a message for topic [%s]\n", m.Subject)

		// decode to cloudevent
		cloudEvent, err := DecodeNATSMessage(m.Data)
		if err != nil {
			return
		}

		// decode to keptn event
		keptnEvent, err := v0_2_0.ToKeptnEvent(*cloudEvent)
		if err != nil {
			return
		}

		// determine subscription for the received message
		subscription := n.getSubscriptionFromReceivedMessage(m, *cloudEvent)

		if subscription != nil {
			// check if the event with the given ID has already been sent for the subscription
			if n.ceCache.Contains(subscription.Event, keptnEvent.ID+"-"+subscription.ID) {
				// Skip this event as it has already been sent
				logger.Infof("CloudEvent with ID %s has already been sent", keptnEvent.ID)
				return
			}
			logger.Infof("Sending CloudEvent with ID %s to %s", keptnEvent.ID, n.env.PubSubRecipient)
			// add to CloudEvents cache
			n.ceCache.Add(subscription.Event, keptnEvent.ID+"-"+subscription.ID)

			defer func() {
				// after some time, remove the cache entry
				go func() {
					<-time.After(10 * time.Second)
					n.ceCache.Remove(subscription.Event, keptnEvent.ID+"-"+subscription.ID)
				}()
			}()
		}

		// add subscription ID as additional information to the keptn event
		if err := keptnEvent.AddTemporaryData("distributor", AdditionalSubscriptionData{SubscriptionID: subscription.ID}, models.AddTemporaryDataOptions{OverwriteIfExisting: true}); err != nil {
			logger.WithError(err).Error("Unable to add additional information about subscriptions to event")
		}

		// forward keptn event
		err = n.sendEvent(keptnEvent, subscription)
		if err != nil {
			logger.Errorf("Could not send CloudEvent: %v", err)
		}
	}()
}

func (n *NATSEventReceiver) getSubscriptionFromReceivedMessage(m *nats.Msg, event cloudevents.Event) *models.EventSubscription {
	var subscriptionForTopic models.EventSubscription
	for _, subscription := range n.currentSubscriptions {
		if subscription.Event == m.Sub.Subject { // need to check against the name of the subscription because this can be a wildcard as well
			matcher := NewEventMatcherFromSubscription(subscription)
			if matcher.Matches(event) {
				subscriptionForTopic = subscription
				break
			}
		}
	}
	return &subscriptionForTopic
}

func (n *NATSEventReceiver) sendEvent(e models.KeptnContextExtendedCE, subscription *models.EventSubscription) error {
	event := v0_2_0.ToCloudEvent(e)
	if subscription != nil {
		matcher := NewEventMatcherFromSubscription(*subscription)
		if !matcher.Matches(event) {
			return nil
		}
	} else if !n.eventMatcher.Matches(event) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx = cloudevents.ContextWithTarget(ctx, n.env.GetPubSubRecipientURL())
	ctx = cloudevents.WithEncodingStructured(ctx)
	defer cancel()

	if err := n.eventSender.Send(ctx, event); err != nil {
		logger.WithError(err).Error("Unable to send event")
		return err
	}

	logger.Infof("sent event %s", event.ID())
	return nil
}
