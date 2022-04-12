package controlplane

import (
	"context"
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/lib-cp-connector/pkg/logger"
)

var ErrEventHandleFatal = errors.New("fatal event handling error")
var ErrEventHandleIgnore = errors.New("event handling error")

type RegistrationData models.Integration

// Integration represents a Keptn Service that wants to receive events from the Keptn Control plane
type Integration interface {
	// OnEvent is called when a new event was received
	OnEvent(context.Context, models.KeptnContextExtendedCE) error

	// RegistrationData is called to get the initial registration data
	RegistrationData() RegistrationData
}

// ControlPlane can be used to connect to the Keptn Control Plane
type ControlPlane struct {
	subscriptionSource   SubscriptionSource
	eventSource          EventSource
	currentSubscriptions []models.EventSubscription
	logger               logger.Logger
}

// New creates a new ControlPlane
// It is using a SubscriptionSource source to get information about current uniform subscriptions
// as well as an EventSource to actually receive events from Keptn
func New(subscriptionSource SubscriptionSource, eventSource EventSource) *ControlPlane {
	return &ControlPlane{
		subscriptionSource:   subscriptionSource,
		eventSource:          eventSource,
		currentSubscriptions: []models.EventSubscription{},
		logger:               logger.NewDefaultLogger(),
	}
}

// Register is initially used to register the Keptn integration to the Control Plane
func (cp *ControlPlane) Register(ctx context.Context, integration Integration) error {
	eventUpdates := make(chan EventUpdate)
	subscriptionUpdates := make(chan []models.EventSubscription)
	if err := cp.eventSource.Start(ctx, integration.RegistrationData(), eventUpdates); err != nil {
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

func (cp *ControlPlane) handle(ctx context.Context, eventUpdate EventUpdate, integration Integration) error {
	for _, subscription := range cp.currentSubscriptions {
		if subscription.Event == eventUpdate.MetaData.Subject {
			matcher := NewEventMatcherFromSubscription(subscription)
			if matcher.Matches(eventUpdate.KeptnEvent) {
				if err := integration.OnEvent(context.WithValue(ctx, EventSenderKey, cp.eventSource.Sender()), eventUpdate.KeptnEvent); err != nil {
					if errors.Is(err, ErrEventHandleFatal) {
						cp.logger.Errorf("Fatal error during handling of event: %v", err)
						return err
					}
					if errors.Is(err, ErrEventHandleIgnore) {
						cp.logger.Warnf("Error during handling of event: %v", err)
					}
				}
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
