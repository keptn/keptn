package controlplane

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"time"
)

type EventSource interface {
	Start(ctx context.Context) error
	RegisterIntegration(integration *Integration) error
	OnSubscriptionUpdate(subscriptions []models.EventSubscription)
	Subscriptions() []models.EventSubscription
}

type HTTPEventSource struct {
	integration          *Integration
	currentSubscriptions []models.EventSubscription
}

func (e *HTTPEventSource) Start(ctx context.Context) error {
	for _, sub := range e.currentSubscriptions {
		events := e.pollEvents(ctx, sub)
		for _, ev := range events {
			e.integration.OnEvent(ev)
		}
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			for _, sub := range e.currentSubscriptions {
				events := e.pollEvents(ctx, sub)
				for _, ev := range events {
					e.integration.OnEvent(ev)
				}
			}
		}
	}
}

func (e *HTTPEventSource) OnSubscriptionUpdate(subscriptions []models.EventSubscription) {
	e.currentSubscriptions = subscriptions
}

func (e *HTTPEventSource) RegisterIntegration(integration *Integration) error {
	e.integration = integration
	return nil
}

func (e *HTTPEventSource) Subscriptions() []models.EventSubscription {
	return e.currentSubscriptions
}

func (e *HTTPEventSource) pollEvents(ctx context.Context, subscription models.EventSubscription) []models.KeptnContextExtendedCE {
	return []models.KeptnContextExtendedCE{
		{ID: "eventID"},
	}
}

type NATSEventSource struct {
}

func (n NATSEventSource) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (n NATSEventSource) RegisterIntegration(integration *Integration) error {
	//TODO implement me
	panic("implement me")
}

func (n NATSEventSource) OnSubscriptionUpdate(subscriptions []models.EventSubscription) {
	//TODO implement me
	panic("implement me")
}

func (n NATSEventSource) Subscriptions() []models.EventSubscription {
	//TODO implement me
	panic("implement me")
}
