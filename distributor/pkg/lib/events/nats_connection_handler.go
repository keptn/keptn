package events

import (
	"errors"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"sort"
	"sync"
)

type NatsConnectionHandler struct {
	NatsConnection *nats.Conn
	Subscriptions  []*nats.Subscription
	topics         []string
	natsURL        string
	MessageHandler func(m *nats.Msg)
	mux            sync.Mutex
}

func NewNatsConnectionHandler(natsURL string) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		natsURL: natsURL,
	}

	return nch
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	nch.mux.Lock()
	defer nch.mux.Unlock()
	for _, sub := range nch.Subscriptions {
		// Unsubscribe
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.Subject)
	}
	nch.NatsConnection.Close()
	nch.Subscriptions = nch.Subscriptions[:0]
}

// SubscribeToTopics expresses interest in the given subject on the NATS message broker.
// Note, that when you pass in subjects via the topics parameter, the NatsConnectionHandler will
// try to subscribe to these topics. If you don't pass any subjects via the topics parameter
// the NatsConnectionHandler will subscribe to the topics configured at instantiation time
func (nch *NatsConnectionHandler) SubscribeToTopics(topics []string) error {
	if nch.natsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if nch.NatsConnection == nil || !nch.NatsConnection.IsConnected() {
		var err error
		nch.RemoveAllSubscriptions()
		nch.mux.Lock()
		defer nch.mux.Unlock()

		logger.Infof("Connecting to NATS server at %s ...", nch.natsURL)
		nch.NatsConnection, err = nats.Connect(nch.natsURL)

		if err != nil {
			return errors.New("failed to create NATS connection: " + err.Error())
		}

		logger.Info("Connected to NATS server")
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		nch.topics = topics

		for _, topic := range nch.topics {
			logger.Infof("Subscribing to topic %s ...", topic)
			sub, err := nch.NatsConnection.Subscribe(topic, nch.MessageHandler)
			if err != nil {
				return errors.New("failed to subscribe to topic: " + err.Error())
			}
			logger.Infof("Subscribed to topic %s", topic)
			nch.Subscriptions = append(nch.Subscriptions, sub)
		}
	}
	//}
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
