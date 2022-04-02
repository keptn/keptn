package controlplane

import (
	"context"
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	"log"
)

var ErrEventHandleFatal = errors.New("fatal event handling error")
var ErrEventHandleIgnore = errors.New("event handling error")

type AdditionalSubscriptionData struct {
	SubscriptionID string `json:"subscriptionID"`
}

type ControlPlaneOptions struct {
	KeptnAPIEndpoint string
	KeptnAPIToken    string
	NATSEndpoint     string
}

type ControlPlane struct {
	subscriptionSource   *SubscriptionSource
	eventSource          EventSource
	currentSubscriptions []models.EventSubscription
}

func New(subscriptionSource *SubscriptionSource, eventSource EventSource) *ControlPlane {
	return &ControlPlane{
		subscriptionSource:   subscriptionSource,
		eventSource:          eventSource,
		currentSubscriptions: []models.EventSubscription{},
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
			err := cp.handle(ctx, event, integration)
			if errors.Is(err, ErrEventHandleFatal) {
				return err
			}
		case subscriptions := <-subscriptionUpdates:
			cp.currentSubscriptions = subscriptions
			cp.eventSource.OnSubscriptionUpdate(subjects(subscriptions))
		case <-ctx.Done():
			return nil
		}
	}
}

func (cp *ControlPlane) handle(ctx context.Context, event models.KeptnContextExtendedCE, integration Integration) error {
	subscriptionsForTopic := []models.EventSubscription{}
	for _, subscription := range cp.currentSubscriptions {
		if subscription.Event == *event.Type { // need to check against the name of the subscription because this can be a wildcard as well
			matcher := NewEventMatcherFromSubscription(subscription)
			if matcher.Matches(event) {
				subscriptionsForTopic = append(subscriptionsForTopic, subscription)
			}
		}
	}

	for _, t := range subscriptionsForTopic {
		if err := event.AddTemporaryData("distributor", AdditionalSubscriptionData{SubscriptionID: t.ID}, models.AddTemporaryDataOptions{OverwriteIfExisting: true}); err != nil {
			log.Printf("Could not add temporary information about subscriptions to event: %v\n", err)
		}
		if err := integration.OnEvent(context.WithValue(ctx, EventSenderKey, cp.eventSource.Sender()), event); err != nil {
			if errors.Is(err, ErrEventHandleFatal) {
				return err
			}
			if errors.Is(err, ErrEventHandleIgnore) {
				log.Print("error during handling of event")
			}
		}
	}
	return nil
}

func subjects(subscriptions []models.EventSubscription) []string {
	var ret []string
	for _, s := range subscriptions {
		ret = append(ret, s.Event)
	}
	return ret
}
