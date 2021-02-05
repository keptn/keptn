package lib

import (
	"errors"
	"fmt"
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

func NewNatsConnectionHandler(natsUrl string, topics []string) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		Topics:  topics,
		NatsURL: natsUrl,
	}

	return nch
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	nch.mux.Lock()
	defer nch.mux.Unlock()
	for _, sub := range nch.Subscriptions {
		// Unsubscribe
		_ = sub.Unsubscribe()
		fmt.Println("Unsubscribed from NATS topic: " + sub.Subject)
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
		fmt.Println("Connecting to NATS server at " + nch.NatsURL + "...")
		nch.NatsConnection, err = nats.Connect(nch.NatsURL)

		if err != nil {
			return errors.New("failed to create NATS connection: " + err.Error())
		}

		fmt.Println("Connected to NATS server")

		for _, topic := range nch.Topics {
			fmt.Println("Subscribing to topic " + topic + "...")
			sub, err := nch.NatsConnection.Subscribe(topic, nch.MessageHandler)
			if err != nil {
				return errors.New("failed to subscribe to topic: " + err.Error())
			}
			fmt.Println("Subscribed to topic " + topic)
			nch.Subscriptions = append(nch.Subscriptions, sub)
		}
	}
	return nil
}
