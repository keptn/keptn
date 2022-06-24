package httpeventsource

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/eventsource/httpeventsource/fake"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestEventSourceCanBeStoppedViaContext(t *testing.T) {
	eventGetSender := &fake.EventAPIMock{}
	eventGetSender.GetFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{{
			Type: strutils.Stringp("sh.keptn.event.task.triggered"),
		}}, nil
	}
	eventChan := make(chan types.EventUpdate)
	ctx, cancel := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	err := New(clock.New(), eventGetSender).Start(ctx, types.RegistrationData{}, eventChan, make(chan error), wg)
	require.NoError(t, err)
	cancel()
	<-eventChan
	wg.Wait()
}

func TestEventSourceCanBeStopped(t *testing.T) {
	eventGetSender := &fake.EventAPIMock{}
	eventGetSender.GetFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{{
			Type: strutils.Stringp("sh.keptn.event.task.triggered"),
		}}, nil
	}
	eventChan := make(chan types.EventUpdate)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	es := New(clock.New(), eventGetSender)
	es.Start(context.TODO(), types.RegistrationData{}, eventChan, make(chan error), wg)
	es.Stop()
	<-eventChan
	wg.Wait()
}

func TestAPICallFailsAfterMaxAttempts(t *testing.T) {
	eventGetSender := &fake.EventAPIMock{}
	eventGetSender.GetFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return nil, fmt.Errorf("error")
	}

	eventChan := make(chan types.EventUpdate)
	errChan := make(chan error)
	eventsource := New(clock.New(), eventGetSender)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	eventsource.maxAttempts = 2

	err := eventsource.Start(context.TODO(), types.RegistrationData{}, eventChan, errChan, wg)
	eventsource.OnSubscriptionUpdate([]models.EventSubscription{{Event: "sh.keptn.event.task.triggered"}})
	require.NoError(t, err)
	<-errChan
	wg.Wait()
}

func TestAPIReceiveEvents(t *testing.T) {
	eventGetSender := &fake.EventAPIMock{}
	eventGetSender.GetFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{
			{
				Type: strutils.Stringp("sh.keptn.event.task.triggered"),
			},
			{
				Type: strutils.Stringp("sh.keptn.event.task2.triggered"),
			}}, nil
	}
	clock := clock.NewMock()
	eventsource := New(clock, eventGetSender)
	eventChan := make(chan types.EventUpdate)

	err := eventsource.Start(context.TODO(), types.RegistrationData{}, eventChan, make(chan error), &sync.WaitGroup{})
	eventsource.OnSubscriptionUpdate([]models.EventSubscription{{ID: "id1", Event: "sh.keptn.event.task.triggered"}, {ID: "id2", Event: "sh.keptn.event.task2.triggered"}})
	require.NoError(t, err)
	clock.Add(time.Second)
	<-eventChan
	clock.Add(time.Second)
	<-eventChan
}

func TestAPIReceiveEventsWithMoreAdvancedFilters(t *testing.T) {
	eventGetSender := &fake.EventAPIMock{}
	eventGetSender.GetFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{
			{
				Data: v0_2_0.EventData{
					Project: "project1",
					Stage:   "stage1",
					Service: "service1",
				},
				Type: strutils.Stringp("sh.keptn.event.task.triggered"),
			},
			{
				Type: strutils.Stringp("sh.keptn.event.task2.triggered"),
			}}, nil
	}
	clock := clock.NewMock()
	eventsource := New(clock, eventGetSender)
	eventChan := make(chan types.EventUpdate)

	err := eventsource.Start(context.TODO(), types.RegistrationData{}, eventChan, make(chan error), &sync.WaitGroup{})
	eventsource.OnSubscriptionUpdate([]models.EventSubscription{{ID: "id1", Event: "sh.keptn.event.task.triggered", Filter: models.EventSubscriptionFilter{
		Projects: []string{"project1"},
		Stages:   []string{"stage1"},
		Services: []string{"service1"},
	}}})
	require.NoError(t, err)
	clock.Add(time.Second)
	<-eventChan
}

func TestAPIPassEventOnlyOnce(t *testing.T) {
	eventGetSender := &fake.EventAPIMock{}
	eventGetSender.GetFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{
			{
				Type: strutils.Stringp("sh.keptn.event.task.triggered"),
			},
		}, nil
	}
	clock := clock.NewMock()
	eventsource := New(clock, eventGetSender)
	eventChan := make(chan types.EventUpdate)

	err := eventsource.Start(context.TODO(), types.RegistrationData{}, eventChan, make(chan error), &sync.WaitGroup{})
	eventsource.OnSubscriptionUpdate([]models.EventSubscription{{ID: "id1", Event: "sh.keptn.event.task.triggered"}})
	require.NoError(t, err)

	eventsReceived := 0
	go func() {
		for {
			<-eventChan
			eventsReceived++
		}
	}()
	clock.Add(time.Second)
	clock.Add(time.Second)
	time.Sleep(time.Second)
	require.Equal(t, 1, eventsReceived)
}

func TestEventSourceGetSender(t *testing.T) {
	senderCalled := false
	sender := func(keptnContextExtendedCE models.KeptnContextExtendedCE) error {
		senderCalled = true
		return nil
	}
	eventGetSender := &fake.EventAPIMock{
		SendFunc: sender,
	}
	err := New(clock.New(), eventGetSender).Sender()(models.KeptnContextExtendedCE{})
	require.NoError(t, err)
	require.True(t, senderCalled)

}
