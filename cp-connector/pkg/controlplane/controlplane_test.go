package controlplane

import (
	"context"
	"fmt"
	fake2 "github.com/keptn/keptn/cp-connector/pkg/fake"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/stretchr/testify/require"
)

type ExampleIntegration struct {
	OnEventFn          func(ctx context.Context, ce models.KeptnContextExtendedCE) error
	RegistrationDataFn func() types.RegistrationData
}

func (e ExampleIntegration) OnEvent(ctx context.Context, ce models.KeptnContextExtendedCE) error {
	if e.OnEventFn != nil {
		return e.OnEventFn(ctx, ce)
	}
	panic("implement me")
}

func (e ExampleIntegration) RegistrationData() types.RegistrationData {
	if e.RegistrationDataFn != nil {
		return e.RegistrationDataFn()
	}
	panic("implement me")
}

type LogForwarderMock struct {
	ForwardFn func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error
}

func (l LogForwarderMock) Forward(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
	if l.ForwardFn != nil {
		return l.ForwardFn(keptnEvent, integrationID)
	}
	panic("implement me")
}

func TestControlPlaneInitialRegistrationFails(t *testing.T) {
	ssm := &fake2.SubscriptionSourceMock{
		RegisterFn: func(integration models.Integration) (string, error) {
			return "", fmt.Errorf("some err")
		},
	}
	esm := &fake2.EventSourceMock{}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}
	integration := ExampleIntegration{RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} }}
	err := New(ssm, esm, fm).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneEventSourceFailsToStart(t *testing.T) {
	ssm := &fake2.SubscriptionSourceMock{
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			return fmt.Errorf("error occured")
		}}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}
	integration := ExampleIntegration{RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} }}
	err := New(ssm, esm, fm).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneSubscriptionSourceFailsToStart(t *testing.T) {
	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			return fmt.Errorf("error occured")
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
		return nil
	}}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}
	integration := ExampleIntegration{RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} }}
	err := New(ssm, esm, fm).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneInboundEventIsForwardedToIntegration(t *testing.T) {
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent models.KeptnContextExtendedCE
	eventUpdate := types.EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: types.EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			subsChan = c
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(strings []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
	}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}

	controlPlane := New(ssm, esm, fm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
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

func TestControlPlaneInboundEventIsForwardedToIntegrationWithoutLogForwarder(t *testing.T) {
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent models.KeptnContextExtendedCE
	eventUpdate := types.EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: types.EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			subsChan = c
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(strings []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
	}

	controlPlane := New(ssm, esm, nil)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
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

func TestControlPlaneIntegrationIDIsForwarded(t *testing.T) {
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent models.KeptnContextExtendedCE
	eventUpdate := types.EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: types.EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			if data.ID != "some-other-id" {
				return fmt.Errorf("error occured")
			}
			subsChan = c
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-other-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			if data.ID != "some-other-id" {
				return fmt.Errorf("error occured")
			}
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(subscriptions []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
	}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}

	controlPlane := New(ssm, esm, fm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
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
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent bool

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			subsChan = c
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(subscriptions []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
	}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}

	controlPlane := New(ssm, esm, fm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
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
	eventChan <- types.EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: types.EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	require.Eventually(t, func() bool { return integrationReceivedEvent }, time.Second, time.Millisecond*100)
	require.Never(t, func() bool { return controlPlaneErr != nil }, time.Second, time.Millisecond*100)
}

func TestControlPlaneIntegrationOnEventThrowsFatalError(t *testing.T) {
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription
	var integrationReceivedEvent bool

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			subsChan = c
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			eventChan = ces
			return nil
		},
		OnSubscriptionUpdateFn: func(subscriptions []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
	}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}

	controlPlane := New(ssm, esm, fm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
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
	eventChan <- types.EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: types.EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}

	require.Eventually(t, func() bool { return integrationReceivedEvent }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return controlPlaneErr != nil }, time.Second, time.Millisecond*100)
}

func TestControlPlane_IsRegistered(t *testing.T) {
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription

	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			subsChan = c
			go func() {
				<-ctx.Done()
				wg.Done()
			}()
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-id", nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, ces chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			eventChan = ces
			go func() {
				<-ctx.Done()
				wg.Done()
			}()
			return nil
		},
		OnSubscriptionUpdateFn: func(subscriptions []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
	}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}

	controlPlane := New(ssm, esm, fm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
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

func TestControlPlane_StoppedByReceivingErrEvent(t *testing.T) {
	var eventChan chan types.EventUpdate
	var subsChan chan []models.EventSubscription
	var errorC chan error
	var eventSourceStopCalled bool
	var subscriptionSourceStopCalled bool
	//var integrationReceivedEvent models.KeptnContextExtendedCE
	//eventUpdate := types.EventUpdate{KeptnEvent: models.KeptnContextExtendedCE{ID: "some-id", Type: strutils.Stringp("sh.keptn.event.echo.triggered")}, MetaData: types.EventUpdateMetaData{Subject: "sh.keptn.event.echo.triggered"}}
	callBackSender := func(ce models.KeptnContextExtendedCE) error { return nil }

	ssm := &fake2.SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, subC chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
			subsChan = subC
			errorC = errC
			return nil
		},
		RegisterFn: func(integration models.Integration) (string, error) {
			return "some-other-id", nil
		},
		StopFn: func() error {
			subscriptionSourceStopCalled = true
			return nil
		},
	}
	esm := &fake2.EventSourceMock{
		StartFn: func(ctx context.Context, data types.RegistrationData, evC chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
			eventChan = evC
			errorC = errC
			go func() {
				errorC <- fmt.Errorf("LKJ") //SOMETHING WENT WRONG
			}()
			return nil
		},
		OnSubscriptionUpdateFn: func(subscriptions []models.EventSubscription) {},
		SenderFn:               func() types.EventSender { return callBackSender },
		StopFn: func() error {
			eventSourceStopCalled = true
			return nil
		},
	}
	fm := &LogForwarderMock{
		ForwardFn: func(keptnEvent models.KeptnContextExtendedCE, integrationID string) error {
			return nil
		},
	}

	controlPlane := New(ssm, esm, fm)

	integration := ExampleIntegration{
		RegistrationDataFn: func() types.RegistrationData { return types.RegistrationData{} },
		OnEventFn:          func(ctx context.Context, ce models.KeptnContextExtendedCE) error { return nil },
	}

	go controlPlane.Register(context.TODO(), integration)
	require.Eventually(t, func() bool { return subsChan != nil }, time.Second, time.Millisecond*100)
	require.Eventually(t, func() bool { return eventChan != nil }, time.Second, time.Millisecond*100)

	go func() {
		fmt.Println("printing to channel")
		fmt.Println(errorC)
		errorC <- fmt.Errorf("some-error")
	}()

	fmt.Println("waiting for err on channel")
	<-errorC

	require.Eventually(t, func() bool {
		return subscriptionSourceStopCalled && eventSourceStopCalled
	}, time.Second, 100*time.Millisecond)
}
