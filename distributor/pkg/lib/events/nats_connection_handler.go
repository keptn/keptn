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

func (nch *NatsConnectionHandler) SubscribeToTopics() error {
	if nch.NatsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if len(nch.Topics) == 0 {
		return errors.New("no PubSub Topics defined")
	}

	var err error

	if nch.NatsConnection == nil || !nch.NatsConnection.IsConnected() {
		nch.RemoveAllSubscriptions()
		nch.mux.Lock()
		defer nch.mux.Unlock()
		logger.Infof("Connecting to NATS server at %s ...", nch.NatsURL)
		nch.NatsConnection, err = nats.Connect(nch.NatsURL)

		if err != nil {
			return errors.New("failed to create NATS connection: " + err.Error())
		}

		logger.Info("Connected to NATS server")

		for _, topic := range nch.Topics {
			logger.Infof("Subscribing to topic %s ...", topic)
			sub, err := nch.NatsConnection.Subscribe(topic, nch.MessageHandler)
			if err != nil {
				return errors.New("failed to subscribe to topic: " + err.Error())
			}
			logger.Infof("Subscribed to topic %s", topic)
			nch.Subscriptions = append(nch.Subscriptions, sub)
		}
	}
	return nil
}
