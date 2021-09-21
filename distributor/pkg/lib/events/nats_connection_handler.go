package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"
)

const streamName = "sh"

type NatsConnectionHandler struct {
	NatsConnection *nats.Conn
	Subscriptions  []*nats.Subscription
	topics         []string
	natsURL        string
	MessageHandler func(m *nats.Msg)
	mux            sync.Mutex
	JetStream      nats.JetStreamContext
}

func NewNatsConnectionHandler(natsURL string) *NatsConnectionHandler {
	nch := &NatsConnectionHandler{
		natsURL: natsURL,
	}

	return nch
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	for _, sub := range nch.Subscriptions {
		// Unsubscribe
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.Subject)
	}
	nch.Subscriptions = nch.Subscriptions[:0]
}

// SubscribeToTopics expresses interest in the given subject on the NATS message broker.
// Note, that when you pass in subjects via the topics parameter, the NatsConnectionHandler will
// try to subscribe to these topics. If you don't pass any subjects via the topics parameter
// the NatsConnectionHandler will subscribe to the topics configured at instantiation time
func (nch *NatsConnectionHandler) SubscribeToTopics(topics []string) error {
	nch.mux.Lock()
	defer nch.mux.Unlock()
	if nch.natsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if nch.NatsConnection == nil || !nch.NatsConnection.IsConnected() {
		if err := nch.establishNatsConnection(); err != nil {
			return err
		}
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		nch.RemoveAllSubscriptions()
		nch.topics = topics

		for _, topic := range nch.topics {
			logger.Infof("Subscribing to topic %s ...", topic)

			pullSubscribe, err := nch.JetStream.PullSubscribe(topic, "consumer-id")
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					msgs, _ := pullSubscribe.Fetch(10, nats.Context(ctx))
					for _, msg := range msgs {
						nch.MessageHandler(msg)
					}
				}
			}()

			//sub, err := nch.JetStream.Subscribe(topic, nch.MessageHandler, nats.StartSequence(1), nats.ManualAck())
			//if err != nil {
			//	return errors.New("failed to subscribe to topic: " + err.Error())
			//}
			//logger.Infof("Subscribed to topic %s", topic)
			nch.Subscriptions = append(nch.Subscriptions, pullSubscribe)
		}
	}
	return nil
}

func (nch *NatsConnectionHandler) establishNatsConnection() error {
	var err error
	nch.RemoveAllSubscriptions()

	nch.NatsConnection.Close()
	logger.Infof("Connecting to NATS server at %s ...", nch.natsURL)
	nch.NatsConnection, err = nats.Connect(nch.natsURL)

	if err != nil {
		return errors.New("failed to create NATS connection: " + err.Error())
	}

	js, err := nch.NatsConnection.JetStream()
	if err != nil {
		return errors.New("failed to create JetStream client: " + err.Error())
	}
	nch.JetStream = js

	stream, err := js.StreamInfo(streamName)
	if err != nil {
		logger.Errorf("could not get stream info: %s", err.Error())
	}

	if stream == nil {
		logger.Infof("creating stream %q", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamName + ".>"},
		})
		if err != nil {
			return fmt.Errorf("could not create stream: %s", err.Error())
		}
	}

	js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable: "consumer-id",
	})

	logger.Info("Connected to NATS server")
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
