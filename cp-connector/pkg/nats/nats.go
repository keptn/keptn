package nats

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

var _ NATS = (*NatsConnector)(nil)

const (
	EnvVarNatsURL        = "NATS_URL"
	EnvVarNatsURLDefault = "nats://keptn-nats"
	CloudEventsVersionV1 = "1.0"
)

type NATS interface {
	Subscribe(subject string, fn ProcessEventFn) error
	QueueSubscribe(queueGroup string, subject string, fn ProcessEventFn) error
	SubscribeMultiple(subjects []string, fn ProcessEventFn) error
	QueueSubscribeMultiple(subjects []string, queueGroup string, fn ProcessEventFn) error
	Publish(event models.KeptnContextExtendedCE) error
	Disconnect() error
	UnsubscribeAll() error
}

var (
	ErrSubAlreadySubscribed   = errors.New("already subscribed")
	ErrSubNilMessageProcessor = errors.New("message processor is nil")
	ErrSubEmptySubject        = errors.New("empty subject")
	ErrPubEventTypeMissing    = errors.New("event is missing the event type")
)

// ProcessEventFn is used to process a received keptn event
type ProcessEventFn func(msg *nats.Msg) error

// NatsConnector can be used to subscribe to certain events
// on the NATS event system
type NatsConnector struct {
	connection    *nats.Conn
	connectURL    string
	subscriptions map[string]*nats.Subscription
	logger        logger.Logger
}

// WithLogger sets the logger to use
func WithLogger(logger logger.Logger) func(*NatsConnector) {
	return func(n *NatsConnector) {
		n.logger = logger
	}
}

// New returns an initialised NatsConnector with a nil connection
func New(connectURL string, opts ...func(connector *NatsConnector)) *NatsConnector {
	nc := &NatsConnector{
		connection:    &nats.Conn{},
		connectURL:    connectURL,
		subscriptions: make(map[string]*nats.Subscription),
		logger:        logger.NewDefaultLogger(),
	}
	for _, o := range opts {
		o(nc)
	}
	return nc
}

// NewFromEnv returns a NatsConnector to NATS.
// The URL is read from the environment variable "NATS_URL"
// If the URL is not set via the environment variable "NATS_URL",
// it falls back to the default URL "nats://keptn-nats"
func NewFromEnv() *NatsConnector {
	natsURL := os.Getenv(EnvVarNatsURL)
	if natsURL == "" {
		natsURL = EnvVarNatsURLDefault
	}
	return New(natsURL)
}

// getOrCreateConnection connects a NatsConnector or returns the existing connection to NATS
// Note that this will automatically and indefinitely try to reconnect
// as soon as it looses connection
func (nc *NatsConnector) getOrCreateConnection() (*nats.Conn, error) {

	if !nc.connection.IsConnected() {
		var err error
		nc.connection, err = nats.Connect(nc.connectURL, nats.MaxReconnects(-1))

		if err != nil {
			return nil, fmt.Errorf("could not connect to NATS: %w", err)
		}
	}

	return nc.connection, nil
}

// UnsubscribeAll deletes all current subscriptions
func (nc *NatsConnector) UnsubscribeAll() error {
	for _, s := range nc.subscriptions {
		if err := s.Unsubscribe(); err != nil {
			return fmt.Errorf("unable to unsubscribe from subject %s: %w", s.Subject, err)
		}
	}
	nc.subscriptions = make(map[string]*nats.Subscription)
	return nil
}

// Subscribe adds a subscription to a specific subject to the NatsConnector.
// It takes the subject as string (usually the event type) and a function fn
// being called when an event is received
func (nc *NatsConnector) Subscribe(subject string, fn ProcessEventFn) error {
	return nc.QueueSubscribe(subject, "", fn)
}

// QueueSubscribe adds a queue subscription to the NatsConnector
func (nc *NatsConnector) QueueSubscribe(subject string, queueGroup string, fn ProcessEventFn) error {
	if subject == "" {
		return ErrSubEmptySubject
	}
	if fn == nil {
		return ErrSubNilMessageProcessor
	}
	return nc.queueSubscribe(subject, queueGroup, fn)
}

// SubscribeMultiple adds multiple subscriptions to the NatsConnector
func (nc *NatsConnector) SubscribeMultiple(subjects []string, fn ProcessEventFn) error {
	return nc.QueueSubscribeMultiple(subjects, "", fn)
}

// QueueSubscribeMultiple adds multiple queue subscriptions to the NatsConnector
func (nc *NatsConnector) QueueSubscribeMultiple(subjects []string, queueGroup string, fn ProcessEventFn) error {
	if fn == nil {
		return ErrSubNilMessageProcessor
	}

	if _, err := nc.getOrCreateConnection(); err != nil {
		return err
	}
	for _, sub := range subjects {
		nc.logger.Debug("Subscribing to topic %s", sub)
		if err := nc.queueSubscribe(sub, queueGroup, fn); err != nil {
			return fmt.Errorf("could not subscribe to subject %s: %w", sub, err)
		}
		nc.logger.Debug("Successfully subscribed to topic %s", sub)
	}
	return nil
}

// Publish sends a keptn event to the message broker
func (nc *NatsConnector) Publish(event models.KeptnContextExtendedCE) error {
	if event.Type == nil || *event.Type == "" {
		return ErrPubEventTypeMissing
	}
	// ensure that the mandatory fields time, id and specversion are set in the CloudEvent
	event.Time = time.Now().UTC()
	event.Specversion = CloudEventsVersionV1
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	serializedEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not publish event: %w", err)
	}
	conn, err := nc.getOrCreateConnection()
	if err != nil {
		return fmt.Errorf("could not connect to NATS to publish event: %w", err)
	}
	return conn.Publish(*event.Type, serializedEvent)
}

// Disconnect disconnects/closes the connection to NATS
func (nc *NatsConnector) Disconnect() error {
	connection, err := nc.getOrCreateConnection()
	if err != nil {
		return fmt.Errorf("could not disconnect from NATS: %w", err)
	}
	connection.Close()
	return nil
}

func (nc *NatsConnector) queueSubscribe(subject string, queueGroup string, fn ProcessEventFn) error {
	conn, err := nc.getOrCreateConnection()
	if err != nil {
		return fmt.Errorf("could not queue: %w", err)
	}
	sub, err := conn.QueueSubscribe(subject, queueGroup, func(m *nats.Msg) {
		err := fn(m)
		if err != nil {
			nc.logger.Errorf("Could not process message %s: %v\n", string(m.Data), err)
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
