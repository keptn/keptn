package controlplane

import (
	"context"
	"errors"
	"fmt"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
)

const tmpDataDistributorKey = "distributor"

var ErrEventHandleFatal = errors.New("fatal event handling error")

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
	registered           bool
	integrationID        string
	logForwarder         LogForwarder
}

// New creates a new ControlPlane
// It is using a SubscriptionSource source to get information about current uniform subscriptions
// as well as an EventSource to actually receive events from Keptn
// and a LogForwarder to forward error logs
func New(subscriptionSource SubscriptionSource, eventSource EventSource, logForwarder LogForwarder) *ControlPlane {
	return &ControlPlane{
		subscriptionSource:   subscriptionSource,
		eventSource:          eventSource,
		currentSubscriptions: []models.EventSubscription{},
		logger:               logger.NewDefaultLogger(),
		logForwarder:         logForwarder,
		registered:           false,
	}
}

// Register is initially used to register the Keptn integration to the Control Plane
func (cp *ControlPlane) Register(ctx context.Context, integration Integration) error {
	eventUpdates := make(chan EventUpdate)
	subscriptionUpdates := make(chan []models.EventSubscription)

	var err error
	registrationData := integration.RegistrationData()
	cp.integrationID, err = cp.subscriptionSource.Register(models.Integration(registrationData))
	if err != nil {
		return fmt.Errorf("could not register integration: %w", err)
	}
	registrationData.ID = cp.integrationID

	if err := cp.eventSource.Start(ctx, registrationData, eventUpdates); err != nil {
		return err
	}
	if err := cp.subscriptionSource.Start(ctx, registrationData, subscriptionUpdates); err != nil {
		return err
	}
	cp.registered = true
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
			cp.registered = false
			return nil
		}
	}
}

// IsRegistered can be called to detect whether the controlPlane is registered and ready to receive events
func (cp *ControlPlane) IsRegistered() bool {
	return cp.registered
}

func (cp *ControlPlane) handle(ctx context.Context, eventUpdate EventUpdate, integration Integration) error {
	for _, subscription := range cp.currentSubscriptions {
		if subscription.Event == eventUpdate.MetaData.Subject {
			matcher := NewEventMatcherFromSubscription(subscription)
			if matcher.Matches(eventUpdate.KeptnEvent) {
				if err := cp.forwardMatchedEvent(ctx, eventUpdate, integration, subscription); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (cp *ControlPlane) getSender(sender EventSender) EventSender {
	if cp.logForwarder != nil {
		return func(ce models.KeptnContextExtendedCE) error {
			err := cp.logForwarder.Forward(ce, cp.integrationID)
			if err != nil {
				cp.logger.Warnf("could not forward event")
			}
			return sender(ce)
		}
	} else {
		return sender
	}
}

func (cp *ControlPlane) forwardMatchedEvent(ctx context.Context, eventUpdate EventUpdate, integration Integration, subscription models.EventSubscription) error {
	err := eventUpdate.KeptnEvent.AddTemporaryData(
		tmpDataDistributorKey,
		AdditionalSubscriptionData{
			SubscriptionID: subscription.ID,
		},
		models.AddTemporaryDataOptions{
			OverwriteIfExisting: true,
		},
	)
	if err != nil {
		cp.logger.Warnf("Could not append subscription data to event: %v", err)
	}
	if err := integration.OnEvent(context.WithValue(ctx, EventSenderKey, cp.getSender(cp.eventSource.Sender())), eventUpdate.KeptnEvent); err != nil {
		if errors.Is(err, ErrEventHandleFatal) {
			cp.logger.Errorf("Fatal error during handling of event: %v", err)
			return err
		}
		cp.logger.Warnf("Error during handling of event: %v", err)
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
