package controlplane

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/lib-cp-connector/pkg/nats"
	"log"
)

type EventSenderKeyType struct{}

var EventSenderKey = EventSenderKeyType{}

type EventSender func(ce models.KeptnContextExtendedCE) error

type EventSource interface {
	Start(context.Context, chan models.KeptnContextExtendedCE) error
	OnSubscriptionUpdate([]string)
	Sender() EventSender
}

type NATSEventSource struct {
	//TODO: should be list of string ( topics/subjects )
	currentSubscriptions []string
	connector            *nats.NatsConnector
	eventProcessFn       nats.ProcessEventFn
}

func NewNATSEventSource(natsConnector *nats.NatsConnector) *NATSEventSource {
	return &NATSEventSource{
		currentSubscriptions: []string{},
		connector:            natsConnector,
		eventProcessFn:       func(event models.KeptnContextExtendedCE) error { return nil },
	}
}

func (n *NATSEventSource) Start(ctx context.Context, eventChannel chan models.KeptnContextExtendedCE) error {
	n.eventProcessFn = func(event models.KeptnContextExtendedCE) error {
		eventChannel <- event
		return nil
	}
	if err := n.connector.SubscribeMultiple(n.currentSubscriptions, n.eventProcessFn); err != nil {
		return fmt.Errorf("could not start NATS event source: %w", err)
	}
	return nil
}

func (n *NATSEventSource) OnSubscriptionUpdate(subscriptions []string) {
	n.currentSubscriptions = subscriptions
	err := n.connector.UnsubscribeAll()
	if err != nil {
		log.Printf("error during handling of subscription update: %v\n", err)
	}
	if err := n.connector.SubscribeMultiple(n.currentSubscriptions, n.eventProcessFn); err != nil {
		log.Printf("error during handling of subscription update: %v\n", err)
	}
}

func (n *NATSEventSource) Sender() EventSender {
	return n.connector.Publish
}

type HTTPEventSource struct{}

func (H HTTPEventSource) Start(ctx context.Context, ces chan models.KeptnContextExtendedCE) error {
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
