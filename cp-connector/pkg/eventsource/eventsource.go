package eventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"reflect"
	"sort"
	"sync"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
	natseventsource "github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/nats-io/nats.go"
)

// EventSource is anything that can be used
// to get events from the Keptn Control Plane
type EventSource interface {
	// Start triggers the execution of the EventSource
	Start(context.Context, types.RegistrationData, chan types.EventUpdate, *sync.WaitGroup) error
	// OnSubscriptionUpdate can be called to tell the EventSource that
	// the current subscriptions have been changed
	OnSubscriptionUpdate([]string)
	// Sender returns a component that gives the possiblity to send events back
	// to the Keptn Control plane
	Sender() types.EventSender
	//Stop is stopping the EventSource
	Stop() error
}

// NATSEventSource is an implementation of EventSource
// that is using the NATS event broker internally
type NATSEventSource struct {
	currentSubjects []string
	connector       natseventsource.NATS
	eventProcessFn  natseventsource.ProcessEventFn
	queueGroup      string
	logger          logger.Logger
}

// New creates a new NATSEventSource
func New(natsConnector natseventsource.NATS, opts ...func(source *NATSEventSource)) *NATSEventSource {
	e := &NATSEventSource{
		currentSubjects: []string{},
		connector:       natsConnector,
		eventProcessFn:  func(event *nats.Msg) error { return nil },
		logger:          logger.NewDefaultLogger(),
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

// WithLogger sets the logger to use
func WithLogger(logger logger.Logger) func(*NATSEventSource) {
	return func(ns *NATSEventSource) {
		ns.logger = logger
	}
}

func (n *NATSEventSource) Start(ctx context.Context, registrationData types.RegistrationData, eventChannel chan types.EventUpdate, wg *sync.WaitGroup) error {
	n.queueGroup = registrationData.Name
	n.eventProcessFn = func(event *nats.Msg) error {
		keptnEvent := models.KeptnContextExtendedCE{}
		if err := json.Unmarshal(event.Data, &keptnEvent); err != nil {
			return fmt.Errorf("could not unmarshal message: %w", err)
		}
		eventChannel <- types.EventUpdate{
			KeptnEvent: keptnEvent,
			MetaData:   types.EventUpdateMetaData{event.Sub.Subject},
		}
		return nil
	}

	if err := n.connector.QueueSubscribeMultiple(n.currentSubjects, n.queueGroup, n.eventProcessFn); err != nil {
		return fmt.Errorf("could not start NATS event source: %w", err)
	}
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := n.connector.UnsubscribeAll(); err != nil {
			n.logger.Errorf("Unable to unsubscribe from NATS: %v", err)
			return
		}
		n.logger.Debug("Unsubscribed from NATS")
	}()
	return nil
}

func (n *NATSEventSource) OnSubscriptionUpdate(subjects []string) {
	s := dedup(subjects)
	n.logger.Debugf("Updating subscriptions")
	if !isEqual(n.currentSubjects, s) {
		n.logger.Debugf("Cleaning up %d old subscriptions", len(n.currentSubjects))
		err := n.connector.UnsubscribeAll()
		n.logger.Debug("Unsubscribed from previous subscriptions")
		if err != nil {
			n.logger.Errorf("Could not handle subscription update: %v", err)
			return
		}
		n.logger.Debugf("Subscribing to %d topics", len(s))
		if err := n.connector.QueueSubscribeMultiple(s, n.queueGroup, n.eventProcessFn); err != nil {
			n.logger.Errorf("Could not handle subscription update: %v", err)
			return
		}
		n.currentSubjects = s
		n.logger.Debugf("Subscription to %d topics successful", len(s))
	}
}

func (n *NATSEventSource) Sender() types.EventSender {
	return n.connector.Publish
}

func (n *NATSEventSource) Stop() error {
	return n.connector.Disconnect()
}

func isEqual(a1 []string, a2 []string) bool {
	sort.Strings(a2)
	sort.Strings(a1)
	return reflect.DeepEqual(a1, a2)
}

func dedup(elements []string) []string {
	result := make([]string, 0, len(elements))
	temp := map[string]struct{}{}
	for _, el := range elements {
		if _, ok := temp[el]; !ok {
			temp[el] = struct{}{}
			result = append(result, el)
		}
	}
	return result
}
