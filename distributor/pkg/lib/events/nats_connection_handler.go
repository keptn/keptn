package events

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"sort"
	"sync"
)

type NatsConnectionHandler struct {
	natsConnection *nats.Conn
	subscriptions  []*nats.Subscription
	topics         []string
	natsURL        string
	messageHandler func(m *nats.Msg)
	mux            sync.Mutex
}

// NewNatsConnectionHandler creates a new NATS connection handler to a NATS
// broker at the gieven URL
func NewNatsConnectionHandler(natsURL string) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		natsURL: natsURL,
	}
	return nch
}

// Connect will try to establish a connection to the NATS broker.
// Note that this will automatically indefinitely handle reconnection attempts
func (nch *NatsConnectionHandler) Connect() error {
	var err error
	nch.natsConnection, err = nats.Connect(nch.natsURL, nats.MaxReconnects(-1))
	if err != nil {
		return fmt.Errorf("failed to create NATS connection: %w", err)
	}
	return nil
}

// RemoveAllSubscriptions removes all current subscriptions from the NATS handler
func (nch *NatsConnectionHandler) RemoveAllSubscriptions() error {
	if nch.natsConnection == nil {
		return fmt.Errorf("unable to remove all subscriptions, because no connection to nats was established")
	}
	for _, sub := range nch.subscriptions {
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.Subject)
	}
	nch.subscriptions = nch.subscriptions[:0]
	return nil
}

// SubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
func (nch *NatsConnectionHandler) SubscribeToTopics(topics []string) error {
	return nch.QueueSubscribeToTopics(topics, "")
}

// QueueSubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
// The queueGroup parameter defines a NATS queue group to join when subscribing to the topic(s).
// Only one member of a queue group will be able to receive a published message via NATS.
// Note, that passing queueGroup = "" is equivalent to not using any queue group at all.
func (nch *NatsConnectionHandler) QueueSubscribeToTopics(topics []string, queueGroup string) error {
	nch.mux.Lock()
	defer nch.mux.Unlock()
	if nch.natsConnection == nil {
		return errors.New("unable to remove all subscriptions, because no connection to nats was established")
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		nch.RemoveAllSubscriptions()
		nch.topics = topics

		for _, topic := range nch.topics {
			logger.Infof("Subscribing to topic <%s> with queue group <%s>", topic, queueGroup)
			sub, err := nch.natsConnection.QueueSubscribe(topic, queueGroup, nch.messageHandler)
			if err != nil {
				return errors.New("failed to subscribe to topic: " + err.Error())
			}
			nch.subscriptions = append(nch.subscriptions, sub)
		}
	}
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
