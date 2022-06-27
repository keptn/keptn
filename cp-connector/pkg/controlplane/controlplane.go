package controlplane

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/eventmatcher"
	"github.com/keptn/keptn/cp-connector/pkg/eventsource"
	"github.com/keptn/keptn/cp-connector/pkg/logforwarder"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
	"github.com/keptn/keptn/cp-connector/pkg/subscriptionsource"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type EventSender = types.EventSender
type EventSenderKeyType = types.EventSenderKeyType
type RegistrationData = types.RegistrationData

const tmpDataDistributorKey = "distributor"

var ErrEventHandleFatal = errors.New("fatal event handling error")

// Integration represents a Keptn Service that wants to receive events from the Keptn Control plane
type Integration interface {
	// OnEvent is called when a new event was received
	OnEvent(context.Context, models.KeptnContextExtendedCE) error

	// RegistrationData is called to get the initial registration data
	RegistrationData() types.RegistrationData
}

// ControlPlane can be used to connect to the Keptn Control Plane
type ControlPlane struct {
	subscriptionSource   subscriptionsource.SubscriptionSource
	eventSource          eventsource.EventSource
	currentSubscriptions []models.EventSubscription
	logger               logger.Logger
	registered           bool
	integrationID        string
	logForwarder         logforwarder.LogForwarder
}

// WithLogger sets the logger to use
func WithLogger(logger logger.Logger) func(plane *ControlPlane) {
	return func(ns *ControlPlane) {
		ns.logger = logger
	}
}

// RunWithGracefulShutdown starts the controlplane component which takes care of registering
// the integration and handling events and subscriptions. Further, it supports graceful shutdown handling
// when receiving a SIGHUB, SIGINT, SIGQUIT, SIGARBT or SIGTERM signal.
//
// This call is blocking.
//
//If you want to start the controlplane component with an own context you need to call the Regiser(ctx,integration)
// method on your own
func RunWithGracefulShutdown(controlPlane *ControlPlane, integration Integration, shutdownTimeout time.Duration) error {
	ctxShutdown, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctxShutdown, _ = signal.NotifyContext(ctxShutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() {
		<-ctxShutdown.Done()
		time.Sleep(shutdownTimeout) // shutdown timeout
		log.Printf("failed to gracefully shutdown")
		os.Exit(1)
	}()

	return controlPlane.Register(ctxShutdown, integration)
}

// New creates a new ControlPlane
// It is using a SubscriptionSource source to get information about current uniform subscriptions
// as well as an EventSource to actually receive events from Keptn
// and a LogForwarder to forward error logs
func New(subscriptionSource subscriptionsource.SubscriptionSource, eventSource eventsource.EventSource, logForwarder logforwarder.LogForwarder, opts ...func(plane *ControlPlane)) *ControlPlane {
	cp := &ControlPlane{
		subscriptionSource:   subscriptionSource,
		eventSource:          eventSource,
		currentSubscriptions: []models.EventSubscription{},
		logger:               logger.NewDefaultLogger(),
		logForwarder:         logForwarder,
		registered:           false,
	}
	for _, o := range opts {
		o(cp)
	}
	return cp
}

// Register is initially used to register the Keptn integration to the Control Plane
func (cp *ControlPlane) Register(ctx context.Context, integration Integration) error {
	eventUpdates := make(chan types.EventUpdate)
	subscriptionUpdates := make(chan []models.EventSubscription)
	errC := make(chan error)

	var err error
	registrationData := integration.RegistrationData()
	cp.logger.Debugf("Registering integration %s", integration.RegistrationData().Name)
	cp.integrationID, err = cp.subscriptionSource.Register(models.Integration(registrationData))
	if err != nil {
		return fmt.Errorf("could not register integration: %w", err)
	}
	cp.logger.Debugf("Registered with integration ID %s", cp.integrationID)
	registrationData.ID = cp.integrationID

	// WaitGroup used for synchronized shutdown of eventsource and subscription source
	// during cancellation of the context
	wg := &sync.WaitGroup{}
	wg.Add(2)

	cp.logger.Debugf("Starting event source for integration ID %s", cp.integrationID)
	if err := cp.eventSource.Start(ctx, registrationData, eventUpdates, errC, wg); err != nil {
		return err
	}
	cp.logger.Debugf("Event source started with data: %+v", registrationData)
	cp.logger.Debugf("Starting subscription source for integration ID %s", cp.integrationID)
	if err := cp.subscriptionSource.Start(ctx, registrationData, subscriptionUpdates, errC, wg); err != nil {
		return err
	}
	cp.logger.Debug("Subscription source started")
	cp.registered = true
	for {
		select {
		// event updates
		case event := <-eventUpdates:
			cp.logger.Debug("Got new event update")
			err := cp.handle(ctx, event, integration)
			if errors.Is(err, ErrEventHandleFatal) {
				return err
			}

		// subscription updates
		case subscriptions := <-subscriptionUpdates:
			cp.logger.Debugf("ControlPlane: Got a subscription update with %d subscriptions", len(subscriptions))
			cp.currentSubscriptions = subscriptions
			cp.eventSource.OnSubscriptionUpdate(subscriptions)

		// control plane cancelled via context
		case <-ctx.Done():
			cp.logger.Debug("Controlplane cancelled via context. Unregistering...")
			wg.Wait()
			cp.registered = false
			return nil

		// control plane cancelled via error in either one of the sub components
		case e := <-errC:
			cp.logger.Debugf("Stopping control plane due to error: %v", e)
			cp.cleanup()
			cp.logger.Debug("Waiting for components to shutdown")
			wg.Wait()
			cp.registered = false
			return nil
		}
	}
}

// IsRegistered can be called to detect whether the controlPlane is registered and ready to receive events
func (cp *ControlPlane) IsRegistered() bool {
	return cp.registered
}

func (cp *ControlPlane) handle(ctx context.Context, eventUpdate types.EventUpdate, integration Integration) error {
	cp.logger.Debugf("Received an event of type: %s", *eventUpdate.KeptnEvent.Type)
	for _, subscription := range cp.currentSubscriptions {
		if subscription.Event == eventUpdate.MetaData.Subject {
			cp.logger.Debugf("Check if event matches subscription %s", subscription.ID)
			matcher := eventmatcher.New(subscription)
			if matcher.Matches(eventUpdate.KeptnEvent) {
				cp.logger.Info("Forwarding matched event update: ", eventUpdate.KeptnEvent.ID)
				if err := cp.forwardMatchedEvent(ctx, eventUpdate, integration, subscription); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (cp *ControlPlane) getSender(sender types.EventSender) types.EventSender {
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

func (cp *ControlPlane) forwardMatchedEvent(ctx context.Context, eventUpdate types.EventUpdate, integration Integration, subscription models.EventSubscription) error {
	err := eventUpdate.KeptnEvent.AddTemporaryData(
		tmpDataDistributorKey,
		types.AdditionalSubscriptionData{
			SubscriptionID: subscription.ID,
		},
		models.AddTemporaryDataOptions{
			OverwriteIfExisting: true,
		},
	)
	if err != nil {
		cp.logger.Warnf("Could not append subscription data to event: %v", err)
	}
	if err := integration.OnEvent(context.WithValue(ctx, types.EventSenderKey, cp.getSender(cp.eventSource.Sender())), eventUpdate.KeptnEvent); err != nil {
		if errors.Is(err, ErrEventHandleFatal) {
			cp.logger.Errorf("Fatal error during handling of event: %v", err)
			return err
		}
		cp.logger.Warnf("Error during handling of event: %v", err)
	}
	return nil
}

func (cp *ControlPlane) cleanup() {
	cp.logger.Info("Stopping subscription source...")
	if err := cp.subscriptionSource.Stop(); err != nil {
		log.Fatalf("Unable to stop subscription source: %v", err)
	}
	cp.logger.Info("Stopping event source...")
	if err := cp.eventSource.Stop(); err != nil {
		log.Fatalf("Unable to stop event source: %v", err)
	}
}
