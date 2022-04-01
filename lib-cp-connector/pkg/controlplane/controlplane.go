package controlplane

import (
	"context"
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	"log"
)

var ErrEventHandleFatal = errors.New("fatal event handling error")
var ErrEventHandleIgnore = errors.New("event handling error")

type ControlPlaneOptions struct {
	KeptnAPIEndpoint string
	KeptnAPIToken    string
	NATSEndpoint     string
}

type ControlPlane struct {
	subscriptionSource *SubscriptionSource
	eventSource        EventSource
}

func New(subscriptionSource *SubscriptionSource, eventSource EventSource) *ControlPlane {
	return &ControlPlane{
		subscriptionSource: subscriptionSource,
		eventSource:        eventSource,
	}
}

func (cp *ControlPlane) Register(ctx context.Context, integration Integration) error {
	eventUpdates := make(chan models.KeptnContextExtendedCE)
	subscriptionUpdates := make(chan []models.EventSubscription)
	if err := cp.eventSource.Start(ctx, eventUpdates); err != nil {
		return err
	}
	if err := cp.subscriptionSource.Start(ctx, integration.RegistrationData(), subscriptionUpdates); err != nil {
		return err
	}
	for {
		select {
		case event := <-eventUpdates:
			err := integration.OnEvent(context.WithValue(ctx, EventSenderKey, cp.eventSource.Sender()), event)
			if errors.Is(err, ErrEventHandleFatal) {
				return err
			}
			if errors.Is(err, ErrEventHandleIgnore) {
				log.Print("error during handling of event")
			}
		case subscriptions := <-subscriptionUpdates:
			cp.eventSource.OnSubscriptionUpdate(subscriptions)
		case <-ctx.Done():
			return nil
		}
	}
}
