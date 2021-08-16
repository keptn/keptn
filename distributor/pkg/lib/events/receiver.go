package events

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"
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
	if n.env.PubSubTopic == "" {
		logger.Warn("No pubsub topic defined. No need to create NATS client connection.")
		ctx.Wg.Done()
		return
	}

	n.natsConnectionHandler.MessageHandler = n.handleMessage
	err := n.natsConnectionHandler.SubscribeToTopics(strings.Split(n.env.PubSubTopic, ","))

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
		//case <-uptimeTicker.C:
		//	_ = nch.SubscribeToTopics()
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
	logger.Infof("Got subscription update... ")
	logger.Info(len(subscriptions))
	var topics []string
	for _, s := range subscriptions {
		topics = append(topics, s.Event)
	}
	err := n.natsConnectionHandler.SubscribeToTopics(topics)
	if err != nil {
		logger.Errorf("Unable to subscribe to topics %v", topics)
	}

}

func (n *NATSEventReceiver) handleMessage(m *nats.Msg) {
	go func() {
		logger.Infof("Received a message for topic [%s]\n", m.Subject)
		e, err := DecodeCloudEvent(m.Data)

		if e != nil && err == nil {
			err = n.sendEvent(*e)
			if err != nil {
				logger.Errorf("Could not send CloudEvent: %v", err)
			}
		}
	}()
}

// TODO: remove duplication of this method (poller.go)
func (n *NATSEventReceiver) sendEvent(event cloudevents.Event) error {
	if !n.eventMatcher.Matches(event) {
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
