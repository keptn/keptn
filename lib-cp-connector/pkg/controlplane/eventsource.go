package controlplane

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/lib-cp-connector/pkg/logger"
	natseventsource "github.com/keptn/keptn/lib-cp-connector/pkg/nats"
	"github.com/nats-io/nats.go"
	"reflect"
	"sort"
)

type EventSenderKeyType struct{}

var EventSenderKey = EventSenderKeyType{}

type EventSender func(ce models.KeptnContextExtendedCE) error

// EventUpdate wraps a new Keptn event received from the Event source
type EventUpdate struct {
	KeptnEvent models.KeptnContextExtendedCE
	MetaData   EventUpdateMetaData
}

// EventUpdateMetaData is additional metadata for bound to the
// event received from the event source
type EventUpdateMetaData struct {
	Subject string
}

// EventSource is anything that can be used
// to get events from the Keptn Control Plane
type EventSource interface {
	// Start triggers the execution of the EventSource
	Start(context.Context, RegistrationData, chan EventUpdate) error
	// OnSubscriptionUpdate can be called to tell the EventSource that
	// the current subscriptions have been changed
	OnSubscriptionUpdate([]string)
	// Sender returns a component that gives the possiblity to send events back
	// to the Keptn Control plane
	Sender() EventSender
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

// NewNATSEventSource creates a new NATSEventSource
func NewNATSEventSource(natsConnector natseventsource.NATS) *NATSEventSource {
	return &NATSEventSource{
		currentSubjects: []string{},
		connector:       natsConnector,
		eventProcessFn:  func(event *nats.Msg) error { return nil },
		logger:          logger.NewDefaultLogger(),
	}
}

func (n *NATSEventSource) Start(ctx context.Context, registrationData RegistrationData, eventChannel chan EventUpdate) error {
	n.queueGroup = registrationData.Name
	n.eventProcessFn = func(event *nats.Msg) error {
		keptnEvent := models.KeptnContextExtendedCE{}
		if err := json.Unmarshal(event.Data, &keptnEvent); err != nil {
			return fmt.Errorf("could not unmarshal message: %w", err)
		}
		eventChannel <- EventUpdate{
			KeptnEvent: keptnEvent,
			MetaData:   EventUpdateMetaData{event.Sub.Subject},
		}
		return nil
	}
	if err := n.connector.QueueSubscribeMultiple(n.currentSubjects, n.queueGroup, n.eventProcessFn); err != nil {
		return fmt.Errorf("could not start NATS event source: %w", err)
	}
	go func() {
		<-ctx.Done()
		if err := n.connector.Disconnect(); err != nil {
			n.logger.Errorf("Unable to disconnect from NATS: %v", err)
			return
		}
	}()
	return nil
}

func (n *NATSEventSource) OnSubscriptionUpdate(subjects []string) {
	s := dedup(subjects)
	if !isEqual(n.currentSubjects, s) {
		err := n.connector.UnsubscribeAll()
		if err != nil {
			n.logger.Errorf("Could not handle subscription update: %v", err)
			return
		}
		if err := n.connector.QueueSubscribeMultiple(subjects, n.queueGroup, n.eventProcessFn); err != nil {
			n.logger.Errorf("Could not handle subscription update: %v", err)
			return
		}
		n.currentSubjects = s
	}
}

func (n *NATSEventSource) Sender() EventSender {
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
