package controlplane

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	nats2 "github.com/keptn/keptn/lib-cp-connector/pkg/nats"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type EventSourceMock struct {
	StartFn                func(context.Context, RegistrationData, chan models.KeptnContextExtendedCE) error
	OnSubscriptionUpdateFn func([]string)
	SenderFn               func() EventSender
	StopFn                 func() error
}

func (e *EventSourceMock) Start(ctx context.Context, data RegistrationData, ces chan models.KeptnContextExtendedCE) error {
	if e.StartFn != nil {
		return e.StartFn(ctx, data, ces)
	}
	panic("implement me")
}

func (e *EventSourceMock) OnSubscriptionUpdate(strings []string) {
	if e.OnSubscriptionUpdateFn != nil {
		e.OnSubscriptionUpdateFn(strings)
		return
	}
	panic("implement me")
}

func (e *EventSourceMock) Sender() EventSender {
	if e.SenderFn != nil {
		return e.SenderFn()
	}
	panic("implement me")
}

func (e *EventSourceMock) Stop() error {
	if e.StopFn != nil {
		return e.StopFn()
	}
	panic("implement me")
}

type NATSConnectorMock struct {
	SubscribeFn                 func(string, nats2.ProcessEventFn) error
	QueueSubscribeFn            func(string, string, nats2.ProcessEventFn) error
	SubscribeMultipleFn         func([]string, nats2.ProcessEventFn) error
	QueueSubscribeMultipleFn    func(string, []string, nats2.ProcessEventFn) error
	QueueSubscribeMultipleCalls int
	PublishFn                   func(ce models.KeptnContextExtendedCE) error
	PublishCalls                int
	DisconnectFn                func() error
	DisconnectCalls             int
	UnsubscribeAllFn            func() error
	UnsubscribeAllCalls         int
	QueueGroup                  string
	ProcessEventFn              nats2.ProcessEventFn
}

func (ncm *NATSConnectorMock) Subscribe(subject string, fn nats2.ProcessEventFn) error {
	if ncm.SubscribeFn != nil {
		return ncm.SubscribeFn(subject, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) QueueSubscribe(queueGroup string, subject string, fn nats2.ProcessEventFn) error {
	if ncm.QueueSubscribeFn != nil {
		return ncm.QueueSubscribeFn(queueGroup, subject, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) SubscribeMultiple(subjects []string, fn nats2.ProcessEventFn) error {
	if ncm.SubscribeMultipleFn != nil {
		return ncm.SubscribeMultipleFn(subjects, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) QueueSubscribeMultiple(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error {
	ncm.ProcessEventFn = fn
	ncm.QueueSubscribeMultipleCalls++

	if ncm.QueueSubscribeMultipleFn != nil {
		return ncm.QueueSubscribeMultipleFn(queueGroup, subjects, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) Publish(event models.KeptnContextExtendedCE) error {
	ncm.PublishCalls++
	if ncm.PublishFn != nil {
		return ncm.PublishFn(event)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) Disconnect() error {
	ncm.DisconnectCalls++
	if ncm.DisconnectFn != nil {
		return ncm.DisconnectFn()
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) UnsubscribeAll() error {
	ncm.UnsubscribeAllCalls++
	if ncm.UnsubscribeAllFn != nil {
		return ncm.UnsubscribeAllFn()
	}
	panic("implement me")
}

func TestEventSourceForwardsEventToChannel(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventChannel := make(chan models.KeptnContextExtendedCE)
	eventSource := NewNATSEventSource(natsConnectorMock)
	eventSource.Start(context.TODO(), RegistrationData{}, eventChannel)
	eventSource.OnSubscriptionUpdate([]string{"a"})
	event := models.KeptnContextExtendedCE{ID: "id"}
	go natsConnectorMock.ProcessEventFn(event)
	eventFromChan := <-eventChannel
	require.Equal(t, eventFromChan, event)
}

func TestEventSourceCancelDisconnectsFromBroker(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error { return nil },
		DisconnectFn:             func() error { return nil },
	}
	ctx, cancel := context.WithCancel(context.TODO())
	NewNATSEventSource(natsConnectorMock).Start(ctx, RegistrationData{}, make(chan models.KeptnContextExtendedCE))
	cancel()
	require.Eventually(t, func() bool { return natsConnectorMock.DisconnectCalls == 1 }, 2*time.Second, 100*time.Millisecond)
}

func TestEventSourceCancelDisconnectFromBrokerFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error { return nil },
		DisconnectFn:             func() error { return fmt.Errorf("error occured") },
	}
	ctx, cancel := context.WithCancel(context.TODO())
	NewNATSEventSource(natsConnectorMock).Start(ctx, RegistrationData{}, make(chan models.KeptnContextExtendedCE))
	cancel()
	require.Eventually(t, func() bool { return natsConnectorMock.DisconnectCalls == 1 }, 2*time.Second, 100*time.Millisecond)
}

func TestEventSourceQueueSubscribeFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{QueueSubscribeMultipleFn: func(s string, strings []string, fn nats2.ProcessEventFn) error { return fmt.Errorf("error occured") }}
	eventSource := NewNATSEventSource(natsConnectorMock)
	err := eventSource.Start(context.TODO(), RegistrationData{}, make(chan models.KeptnContextExtendedCE))
	require.Error(t, err)
}

func TestEventSourceOnSubscriptionUpdate(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventSource := NewNATSEventSource(natsConnectorMock)
	err := eventSource.Start(context.TODO(), RegistrationData{}, make(chan models.KeptnContextExtendedCE))
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]string{"a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.QueueSubscribeMultipleCalls)
}

func TestEventSourceOnSubscriptiOnUpdateUnsubscribeAllFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return fmt.Errorf("error occured") },
	}
	eventSource := NewNATSEventSource(natsConnectorMock)
	err := eventSource.Start(context.TODO(), RegistrationData{}, make(chan models.KeptnContextExtendedCE))
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]string{"a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
}

func TestEventSourceOnSubscriptionUpdateQueueSubscribeMultipleFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventSource := NewNATSEventSource(natsConnectorMock)
	err := eventSource.Start(context.TODO(), RegistrationData{}, make(chan models.KeptnContextExtendedCE))
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	natsConnectorMock.QueueSubscribeMultipleFn = func(queueGroup string, subjects []string, fn nats2.ProcessEventFn) error {
		return fmt.Errorf("error occured")
	}
	eventSource.OnSubscriptionUpdate([]string{"a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.QueueSubscribeMultipleCalls)
}

func TestEventSourceGetSender(t *testing.T) {
	event := models.KeptnContextExtendedCE{ID: "id"}
	natsConnectorMock := &NATSConnectorMock{
		PublishFn: func(ce models.KeptnContextExtendedCE) error {
			require.Equal(t, event, ce)
			return nil
		},
	}
	sendFn := NewNATSEventSource(natsConnectorMock).Sender()
	require.NotNil(t, sendFn)
	err := sendFn(event)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.PublishCalls)
}

func TestEventSourceSenderFails(t *testing.T) {
	event := models.KeptnContextExtendedCE{ID: "id"}
	natsConnectorMock := &NATSConnectorMock{
		PublishFn: func(ce models.KeptnContextExtendedCE) error {
			require.Equal(t, event, ce)
			return fmt.Errorf("error occured")
		},
	}
	sendFn := NewNATSEventSource(natsConnectorMock).Sender()
	require.NotNil(t, sendFn)
	err := sendFn(event)
	require.Error(t, err)
	require.Equal(t, 1, natsConnectorMock.PublishCalls)
}

func TestEventSourceStopDisconnectsFromEventBroker(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		DisconnectFn: func() error {
			return nil
		},
	}
	err := NewNATSEventSource(natsConnectorMock).Stop()
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.DisconnectCalls)
}

func TestEventSourceStopFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		DisconnectFn: func() error {
			return fmt.Errorf("error occured")
		},
	}
	err := NewNATSEventSource(natsConnectorMock).Stop()
	require.Error(t, err)
	require.Equal(t, 1, natsConnectorMock.DisconnectCalls)
}
