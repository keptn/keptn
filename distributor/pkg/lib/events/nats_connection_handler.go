package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"sort"
	"sync"
)

const streamName = "sh"

type NatsConnectionHandler struct {
	natsConnection *nats.Conn
	subscriptions  []*PullSubscription // TODO should be an interface
	topics         []string
	natsURL        string
	messageHandler func(m *nats.Msg)
	mux            sync.Mutex
	ctx            context.Context
	jetStream      nats.JetStreamContext
}

func NewNatsConnectionHandler(natsURL string, ctx context.Context) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		natsURL: natsURL,
		ctx:     ctx,
	}
	return nch
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	for _, sub := range nch.subscriptions {
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.subscription.Subject)
	}
	nch.subscriptions = nch.subscriptions[:0]
}

// SubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
func (nch *NatsConnectionHandler) SubscribeToTopics(topics []string) error {
	return nch.QueueSubscribeToTopics(topics, "default")
}

// QueueSubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
// The queueGroup parameter defines a NATS queue group to join when subscribing to the topic(s).
// Only one member of a queue group will be able to receive a published message via NATS.
// Note, that passing queueGroup = "" is equivalent to not using any queue group at all.
func (nch *NatsConnectionHandler) QueueSubscribeToTopics(topics []string, queueGroup string) error {
	nch.mux.Lock()
	defer nch.mux.Unlock()
	if nch.natsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
		var err error
		nch.RemoveAllSubscriptions()

		nch.natsConnection.Close()
		logger.Infof("Connecting to NATS server at %s ...", nch.natsURL)
		nch.natsConnection, err = nats.Connect(nch.natsURL)

		if err != nil {
			return errors.New("failed to create NATS connection: " + err.Error())
		}

		err = nch.setupJetStreamContext()
		if err != nil {
			return err
		}
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		nch.RemoveAllSubscriptions()
		nch.topics = topics

		for _, topic := range nch.topics {
			subscription := NewPullSubscription(queueGroup, topic, nch.jetStream, nch.messageHandler)
			if err := subscription.Activate(); err != nil {
				return fmt.Errorf("could not start subscription: %s", err.Error())
			}
			nch.subscriptions = append(nch.subscriptions, subscription)
		}
	}
	return nil
}

func (nch *NatsConnectionHandler) setupJetStreamContext() error {
	js, err := nch.natsConnection.JetStream()
	if err != nil {
		return fmt.Errorf("failed to create nats jetStream context: %s", err.Error())
	}

	stream, err := js.StreamInfo(streamName)
	//if err != nil {
	//	return fmt.Errorf("failed to retrieve stream info: %s", err.Error())
	//}
	if stream == nil {
		logger.Infof("creating stream %q", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"sh.keptn.>"},
		})
		if err != nil {
			return fmt.Errorf("failed to add stream: %s", err.Error())
		}
	} else {
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"sh.keptn.>"},
		})
		if err != nil {
			return fmt.Errorf("failed to update stream: %s", err.Error())
		}
	}
	nch.jetStream = js
	return nil
}

func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	if len(a1) == len(a2) {
		for i, v := range a1 {
			if v != a2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
