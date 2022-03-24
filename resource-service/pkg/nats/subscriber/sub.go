package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"os"
)

const (
	envVarNatsURL        = "NATS_URL"
	envVarNatsURLDefault = "nats://keptn-nats"
)

var (
	ErrSubAlreadySubscribed   = errors.New("already subscribed")
	ErrSubNilMessageProcessor = errors.New("message processor is nil")
	ErrSubEmptySubject        = errors.New("empty subject")
)

// ProcessEventFn is used to process a received keptn event
type ProcessEventFn func(event models.Event) error

// NatsSubscriber can be used to subscribe to certain events
// on the NATS event system
type NatsSubscriber struct {
	conn          *nats.Conn
	subscriptions map[string]*nats.Subscription
}

// Connect connects a NatsSubscriber to NATS.
// Note that this will automatically and indefinitely try to reconnect
// as soon as it looses connection
func Connect(connectURL string) (*NatsSubscriber, error) {
	conn, err := nats.Connect(connectURL, nats.MaxReconnects(-1))
	if err != nil {
		return nil, fmt.Errorf("could not connect to NATS: %w", err)
	}
	return &NatsSubscriber{
		conn:          conn,
		subscriptions: make(map[string]*nats.Subscription),
	}, nil
}

// ConnectFromEnv connects a NatsSubscriber to NATS.
// The URL is read from the environment variable "NATS_URL"
// If the URL is not set via the environment variable "NATS_URL",
// it falls back to the default URL "nats://keptn-nats"
func ConnectFromEnv() (*NatsSubscriber, error) {
	natsURL := os.Getenv(envVarNatsURL)
	if natsURL == "" {
		natsURL = envVarNatsURLDefault
	}
	return Connect(natsURL)
}

// Subscribe adds a subscription to a specific subject to the NatsSubscriber.
// It takes the subject as string (usually the event type) and a function fn
// being called when an event is received
func (nc *NatsSubscriber) Subscribe(subject string, fn ProcessEventFn) error {
	if subject == "" {
		return ErrSubEmptySubject
	}
	if fn == nil {
		return ErrSubNilMessageProcessor
	}
	sub, err := nc.conn.Subscribe(subject, func(m *nats.Msg) {
		event := &models.Event{}
		if err := json.Unmarshal(m.Data, event); err != nil {
			logger.Errorf("could not unmarshal message %s: %v", string(m.Data), err)
			return
		}
		err := fn(*event)
		if err != nil {
			logger.Errorf("Could not process message %s: %v", string(m.Data), err)
		}
	})
	if err != nil {
		return fmt.Errorf("could not subscribe to subject %s: %w", subject, err)
	}
	if nc.subscriptions == nil {
		nc.subscriptions = make(map[string]*nats.Subscription)
	}

	if _, ok := nc.subscriptions[subject]; ok {
		return ErrSubAlreadySubscribed
	}

	nc.subscriptions[subject] = sub
	return nil
}
