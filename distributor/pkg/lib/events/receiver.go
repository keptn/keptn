package events

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"os"
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
		logger.Infof("Received a message for topic [%s]\n", m.Subject)
		e, err := DecodeCloudEvent(m.Data)
		if e != nil && err == nil {
			var subscriptionForTopic *models.EventSubscription
			for _, subscription := range n.currentSubscriptions {
				if subscription.Event == m.Sub.Subject { // need to check against the name of the subscription because this can be a wildcard as well
					subscriptionForTopic = &subscription
					break
				}
			}
			err = n.sendEvent(*e, subscriptionForTopic)
			if err != nil {
				logger.Errorf("Could not send CloudEvent: %v", err)
			}
		}
	}()
}

// TODO: remove duplication of this method (poller.go)
func (n *NATSEventReceiver) sendEvent(event cloudevents.Event, subscription *models.EventSubscription) error {
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
