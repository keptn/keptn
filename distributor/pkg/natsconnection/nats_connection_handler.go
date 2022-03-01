package natsconnection

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"reflect"
	"sort"
	"sync"
)

type NatsConnectionHandler struct {
	MessageHandler func(m *nats.Msg)
	natsConnection *nats.Conn
	subscriptions  []*nats.Subscription
	topics         []string
	natsURL        string
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
	disconnectLogger := func(con *nats.Conn, err error) {
		if err != nil {
			logger.Errorf("Disconnected from NATS due to an error: %v", err)
		} else {
			logger.Info("Disconnected from NATS")
		}
	}
	reconnectLogger := func(*nats.Conn) {
		logger.Info("Reconnected to NATS")
	}
	var err error
	nch.natsConnection, err = nats.Connect(nch.natsURL, nats.ReconnectHandler(reconnectLogger), nats.DisconnectErrHandler(disconnectLogger), nats.RetryOnFailedConnect(true), nats.MaxReconnects(-1))
	if err != nil {
		return fmt.Errorf("failed to create NATS connection: %w", err)
	}
	return nil
}

// RemoveAllSubscriptions removes all current subscriptions from the NATS handler
func (nch *NatsConnectionHandler) RemoveAllSubscriptions() error {
	if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
		return fmt.Errorf("could not remove all subscriptions, because not connected to NATS")
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
	if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
		return errors.New("could not remove all subscriptions, because not connected to NATS")
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		err := nch.RemoveAllSubscriptions()
		if err != nil {
			return fmt.Errorf("could not remove all subscriptions: %w", err)
		}
		nch.topics = topics

		for _, topic := range nch.topics {
			logger.Infof("Subscribing to topic '%s' with queue group '%s'", topic, queueGroup)
			sub, err := nch.natsConnection.QueueSubscribe(topic, queueGroup, nch.MessageHandler)
			if err != nil {
				return errors.New("failed to subscribe to topic: " + err.Error())
			}
			nch.subscriptions = append(nch.subscriptions, sub)
		}
	}
	return nil
}

func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a2)
	sort.Strings(a1)
	return reflect.DeepEqual(a1, a2)
}
