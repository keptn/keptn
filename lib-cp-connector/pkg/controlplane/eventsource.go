package controlplane

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/lib-cp-connector/pkg/nats"
	"log"
	"reflect"
	"sort"
)

type EventSenderKeyType struct{}

var EventSenderKey = EventSenderKeyType{}

type EventSender func(ce models.KeptnContextExtendedCE) error

type EventSource interface {
	Start(context.Context, RegistrationData, chan models.KeptnContextExtendedCE) error
	OnSubscriptionUpdate([]string)
	Sender() EventSender
}

type NATSEventSource struct {
	currentSubjects []string
	connector       *nats.NatsConnector
	eventProcessFn  nats.ProcessEventFn
	queueGroup      string
}

func NewNATSEventSource(natsConnector *nats.NatsConnector) *NATSEventSource {
	return &NATSEventSource{
		currentSubjects: []string{},
		connector:       natsConnector,
		eventProcessFn:  func(event models.KeptnContextExtendedCE) error { return nil },
	}
}

func (n *NATSEventSource) Start(ctx context.Context, registrationData RegistrationData, eventChannel chan models.KeptnContextExtendedCE) error {
	n.queueGroup = registrationData.Name
	n.eventProcessFn = func(event models.KeptnContextExtendedCE) error {
		eventChannel <- event
		return nil
	}
	if err := n.connector.QueueSubscribeMultiple(n.queueGroup, n.currentSubjects, n.eventProcessFn); err != nil {
		return fmt.Errorf("could not start NATS event source: %w", err)
	}
	return nil
}

func (n *NATSEventSource) OnSubscriptionUpdate(subjects []string) {
	s := dedup(subjects)
	if !isEqual(n.currentSubjects, s) {
		n.currentSubjects = s
		err := n.connector.UnsubscribeAll()
		if err != nil {
			log.Printf("error during handling of subscription update: %v\n", err)
		}
		if err := n.connector.QueueSubscribeMultiple(n.queueGroup, n.currentSubjects, n.eventProcessFn); err != nil {
			log.Printf("error during handling of subscription update: %v\n", err)
		}
	}
}

func (n *NATSEventSource) Sender() EventSender {
	return n.connector.Publish
}

type HTTPEventSource struct{}

func (H HTTPEventSource) Start(ctx context.Context, registrationData RegistrationData, ces chan models.KeptnContextExtendedCE) error {
	//TODO implement me
	panic("implement me")
}

func (H HTTPEventSource) OnSubscriptionUpdate(strings []string) {
	//TODO implement me
	panic("implement me")
}

func (H HTTPEventSource) Sender() EventSender {
	//TODO implement me
	panic("implement me")
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
