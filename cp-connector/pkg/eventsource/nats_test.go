package eventsource

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"sync"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	nats2 "github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

type NATSConnectorMock struct {
	SubscribeFn                 func(string, nats2.ProcessEventFn) error
	QueueSubscribeFn            func(string, string, nats2.ProcessEventFn) error
	SubscribeMultipleFn         func([]string, nats2.ProcessEventFn) error
	QueueSubscribeMultipleFn    func([]string, string, nats2.ProcessEventFn) error
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

func (ncm *NATSConnectorMock) QueueSubscribe(subject string, queueGroup string, fn nats2.ProcessEventFn) error {
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

func (ncm *NATSConnectorMock) QueueSubscribeMultiple(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error {
	ncm.ProcessEventFn = fn
	ncm.QueueSubscribeMultipleCalls++

	if ncm.QueueSubscribeMultipleFn != nil {
		return ncm.QueueSubscribeMultipleFn(subjects, queueGroup, fn)
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
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventChannel := make(chan types.EventUpdate)
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	eventSource.Start(context.TODO(), types.RegistrationData{}, eventChannel, make(chan error), wg)
	eventSource.OnSubscriptionUpdate([]string{"a"})
	event := models.KeptnContextExtendedCE{ID: "id"}
	jsonEvent, _ := event.ToJSON()
	e := &nats.Msg{Data: jsonEvent, Sub: &nats.Subscription{Subject: "subscription"}} //models.KeptnContextExtendedCE{ID: "id"}
	go natsConnectorMock.ProcessEventFn(e)
	eventFromChan := <-eventChannel
	require.Equal(t, eventFromChan.KeptnEvent, event)
}

func TestEventSourceCancelDisconnectsFromBroker(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	ctx, cancel := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	cancel()
	require.Eventually(t, func() bool { return natsConnectorMock.UnsubscribeAllCalls == 1 }, 2*time.Second, 100*time.Millisecond)
}

func TestEventSourceCallsWaitGroupDuringCancellation(t *testing.T) {
	t.Run("WaitGroup called", func(t *testing.T) {
		natsConnectorMock := &NATSConnectorMock{
			QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
			UnsubscribeAllFn:         func() error { return nil },
		}
		ctx, cancel := context.WithCancel(context.TODO())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
		cancel()
		wg.Wait()
	})
	t.Run("WaitGroup called - error in shutdown logic", func(t *testing.T) {
		natsConnectorMock := &NATSConnectorMock{
			QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
			UnsubscribeAllFn:         func() error { return fmt.Errorf("ohoh") },
		}
		ctx, cancel := context.WithCancel(context.TODO())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
		cancel()
		wg.Wait()
	})
}

func TestEventSourceCancelDisconnectFromBrokerFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return fmt.Errorf("error occured") },
	}
	ctx, cancel := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	cancel()
	require.Eventually(t, func() bool { return natsConnectorMock.UnsubscribeAllCalls == 1 }, 2*time.Second, 100*time.Millisecond)
}

func TestEventSourceQueueSubscribeFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{QueueSubscribeMultipleFn: func(strings []string, s string, fn nats2.ProcessEventFn) error { return fmt.Errorf("error occured") }}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.Error(t, err)
}

func TestEventSourceOnSubscriptionUpdate(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]string{"a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.QueueSubscribeMultipleCalls)
}

func TestEventSourceOnSubscriptionupdateWithDuplicatedSubjects(t *testing.T) {
	var receivedSubjects []string
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error {
			receivedSubjects = subjects
			return nil
		},
		UnsubscribeAllFn: func() error { return nil },
	}
	eventSource := New(natsConnectorMock)
	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), &sync.WaitGroup{})
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]string{"a", "a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.QueueSubscribeMultipleCalls)
	require.Equal(t, 1, len(receivedSubjects))
}

func TestEventSourceOnSubscriptiOnUpdateUnsubscribeAllFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return fmt.Errorf("error occured") },
	}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]string{"a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
}

func TestEventSourceOnSubscriptionUpdateQueueSubscribeMultipleFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.QueueSubscribeMultipleCalls)
	natsConnectorMock.QueueSubscribeMultipleFn = func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error {
		return fmt.Errorf("error occured")
	}
	eventSource.OnSubscriptionUpdate([]string{"a"})
	require.Equal(t, 1, natsConnectorMock.UnsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.QueueSubscribeMultipleCalls)
}

func TestEventSourceGetSender(t *testing.T) {
	event := models.KeptnContextExtendedCE{ID: "id", Type: strutils.Stringp("something")}
	natsConnectorMock := &NATSConnectorMock{
		PublishFn: func(ce models.KeptnContextExtendedCE) error {
			require.Equal(t, event, ce)
			return nil
		},
	}
	sendFn := New(natsConnectorMock).Sender()
	require.NotNil(t, sendFn)
	err := sendFn(event)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.PublishCalls)
}

func TestEventSourceSenderFails(t *testing.T) {
	event := models.KeptnContextExtendedCE{ID: "id", Type: strutils.Stringp("something")}
	natsConnectorMock := &NATSConnectorMock{
		PublishFn: func(ce models.KeptnContextExtendedCE) error {
			require.Equal(t, event, ce)
			return fmt.Errorf("error occured")
		},
	}
	sendFn := New(natsConnectorMock).Sender()
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
	err := New(natsConnectorMock).Stop()
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.DisconnectCalls)
}

func TestEventSourceStopFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		DisconnectFn: func() error {
			return fmt.Errorf("error occured")
		},
	}
	err := New(natsConnectorMock).Stop()
	require.Error(t, err)
	require.Equal(t, 1, natsConnectorMock.DisconnectCalls)
}
