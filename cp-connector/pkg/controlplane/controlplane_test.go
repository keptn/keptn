package controlplane

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane/fake"
	"github.com/stretchr/testify/require"
)

type ExampleIntegration struct {
	OnEventFn          func(ctx context.Context, ce models.KeptnContextExtendedCE) error
	RegistrationDataFn func() RegistrationData
}

func (e ExampleIntegration) OnEvent(ctx context.Context, ce models.KeptnContextExtendedCE) error {
	if e.OnEventFn != nil {
		return e.OnEventFn(ctx, ce)
	}
	panic("implement me")
}

func (e ExampleIntegration) RegistrationData() RegistrationData {
	if e.RegistrationDataFn != nil {
		return e.RegistrationDataFn()
	}
	panic("implement me")
}

func TestControlPlaneInitialRegistrationFails(t *testing.T) {
	ssm := &SubscriptionSourceMock{}
	esm := &EventSourceMock{}
	um := &fake.UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			return "", fmt.Errorf("error occured")
		},
	}
	integration := ExampleIntegration{RegistrationDataFn: func() RegistrationData { return RegistrationData{} }}
	err := New(ssm, esm, um).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneEventSourceFailsToStart(t *testing.T) {
	ssm := &SubscriptionSourceMock{}
	esm := &EventSourceMock{StartFn: func(ctx context.Context, data RegistrationData, ces chan EventUpdate) error {
		return fmt.Errorf("error occured")
	}}
	um := &fake.UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	integration := ExampleIntegration{RegistrationDataFn: func() RegistrationData { return RegistrationData{} }}
	err := New(ssm, esm, um).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneSubscriptionSourceFailsToStart(t *testing.T) {
	ssm := &SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
			return fmt.Errorf("error occured")
		},
	}
	esm := &EventSourceMock{StartFn: func(ctx context.Context, data RegistrationData, ces chan EventUpdate) error {
		return nil
	}}
	um := &fake.UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	integration := ExampleIntegration{RegistrationDataFn: func() RegistrationData { return RegistrationData{} }}
	err := New(ssm, esm, um).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneInboundEventIsForwardedToIntegration(t *testing.T) {
	var eventChan chan EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent models.KeptnContextExtendedCE
	eventUpdate := EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
			subsChan = c
			return nil
		},
	}
	esm := &EventSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, ces chan EventUpdate) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(strings []string) {},
		SenderFn:               func() EventSender { return callBackSender },
	}

	um := &fake.UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}

	controlPlane := New(ssm, esm, um)

	integration := ExampleIntegration{
		RegistrationDataFn: func() RegistrationData { return RegistrationData{} },
		OnEventFn: func(ctx context.Context, ce models.KeptnContextExtendedCE) error {
			integrationReceivedEvent = ce
			return nil
		},
	}
	go controlPlane.Register(context.TODO(), integration)
	require.Eventually(t, func() bool { return subsChan != nil }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return eventChan != nil }, time.Second, time.Millisecond*100)

	subsChan <- []models.EventSubscription{{ID: "some-id", Event: "sh.keptn.event.echo.triggered", Filter: models.EventSubscriptionFilter{}}}
	eventChan <- eventUpdate

	require.Eventually(t, func() bool {
		eventUpdate.KeptnEvent.Data = integrationReceivedEvent.Data
		return reflect.DeepEqual(eventUpdate.KeptnEvent, integrationReceivedEvent)
	}, time.Second, time.Millisecond*100)

	eventData := map[string]interface{}{}
	err := integrationReceivedEvent.DataAs(&eventData)
	require.Nil(t, err)

	require.Equal(t, map[string]interface{}{
		"temporaryData": map[string]interface{}{
			"distributor": map[string]interface{}{
				"subscriptionID": "some-id",
			},
		},
	}, eventData)
}

func TestControlPlaneIntegrationOnEventThrowsIgnoreableError(t *testing.T) {
	var eventChan chan EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent bool

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
			subsChan = c
			return nil
		},
	}
	esm := &EventSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, ces chan EventUpdate) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(strings []string) {},
		SenderFn:               func() EventSender { return callBackSender },
	}

	um := &fake.UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}

	controlPlane := New(ssm, esm, um)

	integration := ExampleIntegration{
		RegistrationDataFn: func() RegistrationData { return RegistrationData{} },
		OnEventFn: func(ctx context.Context, ce models.KeptnContextExtendedCE) error {
			integrationReceivedEvent = true
			return fmt.Errorf("could not handle event: %w", fmt.Errorf("error occured"))
		},
	}
	var controlPlaneErr error
	go func() { controlPlaneErr = controlPlane.Register(context.TODO(), integration) }()
	require.Eventually(t, func() bool { return subsChan != nil }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return eventChan != nil }, time.Second, time.Millisecond*100)

	subsChan <- []models.EventSubscription{{ID: "some-id", Event: "sh.keptn.event.echo.triggered", Filter: models.EventSubscriptionFilter{}}}
	eventChan <- EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	require.Eventually(t, func() bool { return integrationReceivedEvent }, time.Second, time.Millisecond*100)
	require.Never(t, func() bool { return controlPlaneErr != nil }, time.Second, time.Millisecond*100)
}

func TestControlPlaneIntegrationOnEventThrowsFatalError(t *testing.T) {
	var eventChan chan EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent bool

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
			subsChan = c
			return nil
		},
	}
	esm := &EventSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, ces chan EventUpdate) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(strings []string) {},
		SenderFn:               func() EventSender { return callBackSender },
	}

	um := &fake.UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}

	controlPlane := New(ssm, esm, um)

	integration := ExampleIntegration{
		RegistrationDataFn: func() RegistrationData { return RegistrationData{} },
		OnEventFn: func(ctx context.Context, ce models.KeptnContextExtendedCE) error {
			integrationReceivedEvent = true
			return fmt.Errorf("could not handle event: %w", ErrEventHandleFatal)
		},
	}
	var controlPlaneErr error
	go func() { controlPlaneErr = controlPlane.Register(context.TODO(), integration) }()
	require.Eventually(t, func() bool { return subsChan != nil }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return eventChan != nil }, time.Second, time.Millisecond*100)

	subsChan <- []models.EventSubscription{{ID: "some-id", Event: "sh.keptn.event.echo.triggered", Filter: models.EventSubscriptionFilter{}}}
	eventChan <- EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	require.Eventually(t, func() bool { return integrationReceivedEvent }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return controlPlaneErr != nil }, time.Second, time.Millisecond*100)
}

func TestControlPlane_IsRegistered(t *testing.T) {
	var eventChan chan EventUpdate
	var subsChan chan []models.EventSubscription

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
			subsChan = c
			return nil
		},
	}
	esm := &EventSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, ces chan EventUpdate) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(strings []string) {},
		SenderFn:               func() EventSender { return callBackSender },
	}

	controlPlane := New(ssm, esm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() RegistrationData { return RegistrationData{} },
		OnEventFn: func(ctx context.Context, ce models.KeptnContextExtendedCE) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.TODO())

	require.False(t, controlPlane.IsRegistered())

	go func() { _ = controlPlane.Register(ctx, integration) }()
	require.Eventually(t, func() bool { return subsChan != nil }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return eventChan != nil }, time.Second, time.Millisecond*100)
	require.True(t, controlPlane.IsRegistered())

	cancel()

	require.Eventually(t, func() bool {
		return !controlPlane.IsRegistered()
	}, time.Second, 100*time.Millisecond)
}
