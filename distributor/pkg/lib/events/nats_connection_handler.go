package events

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsConnectionHandler struct {
	NatsConnection *nats.Conn
	Subscriptions  []*nats.Subscription
	Topics         []string
	NatsURL        string
	MessageHandler func(m *nats.Msg)

	uptimeTicker *time.Ticker
	mux          sync.Mutex
}

func NewNatsConnectionHandler(natsURL string, topics []string) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		Topics:  topics,
		NatsURL: natsURL,
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
// the NatsConnectionHandler will subscribe to the topics configured at instatiation time
func (nch *NatsConnectionHandler) SubscribeToTopics(topics ...string) error {
	if nch.NatsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if len(nch.Topics) == 0 {
		return errors.New("no PubSub Topics defined")
	}

	if nch.NatsConnection == nil || !nch.NatsConnection.IsConnected() {
		var err error
		nch.RemoveAllSubscriptions()
		nch.mux.Lock()
		defer nch.mux.Unlock()

		if len(topics) > 0 {
			nch.Topics = topics
		}

		logger.Infof("Connecting to NATS server at %s ...", nch.NatsURL)
		nch.NatsConnection, err = nats.Connect(nch.NatsURL)

		if err != nil {
			return errors.New("failed to create NATS connection: " + err.Error())
		}

		logger.Info("Connected to NATS server")
	}

	for _, topic := range nch.Topics {
		logger.Infof("Subscribing to topic %s ...", topic)
		sub, err := nch.NatsConnection.Subscribe(topic, nch.MessageHandler)
		if err != nil {
			return errors.New("failed to subscribe to topic: " + err.Error())
		}
		logger.Infof("Subscribed to topic %s", topic)
		nch.Subscriptions = append(nch.Subscriptions, sub)
	}
	//}
	return nil
}
