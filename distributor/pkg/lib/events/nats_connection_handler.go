package events

import (
	"errors"
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

func NewNatsConnectionHandler(natsURL string) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		natsURL: natsURL,
	}
	return nch
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	for _, sub := range nch.subscriptions {
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.Subject)
	}
	nch.subscriptions = nch.subscriptions[:0]
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
