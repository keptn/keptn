package handler

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/keptn/keptn/shipyard-controller/common"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_WhenTimeOfEventIsOlder_EventIsSentImmediately(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()

	clock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}
	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeBefore}

	dispatcher.Add(dispatcherEvent)
	require.Equal(t, 1, len(eventSender.SentEvents))
}

func Test_WhenTimeOfEventIsYounger_EventIsQueued(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()

	eventQueueRepo.QueueEventFunc = func(item models.QueueItem) error {
		return nil
	}
	clock.Set(timeBefore)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}

	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeAfter}

	dispatcher.Add(dispatcherEvent)
	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Equal(t, 1, len(eventQueueRepo.QueueEventCalls()))
}

func Test_WhenSyncTimeElapses_EventsAreDispatched(t *testing.T) {

	timeNow := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter1 := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)
	timeAfter2 := time.Date(2021, 4, 21, 15, 00, 00, 2, time.UTC)
	timeAfter3 := time.Date(2021, 4, 21, 15, 00, 00, 3, time.UTC)

	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}

	event1, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	dispatcherEvent1 := models.DispatcherEvent{keptnv2.ToCloudEvent(event1), timeAfter1}
	event2, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task2"), "source", data).Build()
	dispatcherEvent2 := models.DispatcherEvent{keptnv2.ToCloudEvent(event2), timeAfter2}
	event3, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task3"), "source", data).Build()
	dispatcherEvent3 := models.DispatcherEvent{keptnv2.ToCloudEvent(event3), timeAfter3}

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()
	clock.Set(timeNow)

	eventQueueRepo.QueueEventFunc = func(item models.QueueItem) error {
		return nil
	}

	eventQueueRepo.GetQueuedEventsFunc = func(timestamp time.Time) ([]models.QueueItem, error) {
		var items []models.QueueItem
		for _, i := range eventQueueRepo.QueueEventCalls() {
			if timestamp.After(i.Item.Timestamp) {
				items = append(items, i.Item)
			}
		}
		return items, nil
	}

	eventQueueRepo.DeleteQueuedEventFunc = func(eventID string) error {
		return nil
	}

	eventRepo.GetEventsFunc = func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
		return []models.Event{{ID: *filter.ID, Specversion: "1.0"}}, nil
	}

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	dispatcher.Add(dispatcherEvent1)
	dispatcher.Add(dispatcherEvent2)
	dispatcher.Add(dispatcherEvent3)

	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Equal(t, 3, len(eventQueueRepo.QueueEventCalls()))
	dispatcher.Run(context.Background())
	clock.Add(9 * time.Second)
	require.Equal(t, 0, len(eventSender.SentEvents))
	clock.Add(1 * time.Second)
	require.Equal(t, 3, len(eventSender.SentEvents))
	// check if the time stamp of the event has been set to the time at which it has been sent
	require.WithinDuration(t, clock.Now().UTC(), eventSender.SentEvents[0].Time(), time.Second)
	require.Equal(t, 3, len(eventQueueRepo.DeleteQueuedEventCalls()))
}

func Test_WhenAnEventCouldNotBeFetched_NextEventIsProcessed(t *testing.T) {

	timeNow := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter1 := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)
	timeAfter2 := time.Date(2021, 4, 21, 15, 00, 00, 2, time.UTC)
	timeAfter3 := time.Date(2021, 4, 21, 15, 00, 00, 3, time.UTC)

	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}

	event1, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	dispatcherEvent1 := models.DispatcherEvent{keptnv2.ToCloudEvent(event1), timeAfter1}
	event2, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task2"), "source", data).Build()
	dispatcherEvent2 := models.DispatcherEvent{keptnv2.ToCloudEvent(event2), timeAfter2}
	event3, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task3"), "source", data).Build()
	dispatcherEvent3 := models.DispatcherEvent{keptnv2.ToCloudEvent(event3), timeAfter3}

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()
	clock.Set(timeNow)

	eventQueueRepo.QueueEventFunc = func(item models.QueueItem) error {
		return nil
	}

	eventQueueRepo.GetQueuedEventsFunc = func(timestamp time.Time) ([]models.QueueItem, error) {
		var items []models.QueueItem
		for _, i := range eventQueueRepo.QueueEventCalls() {
			if timestamp.After(i.Item.Timestamp) {
				items = append(items, i.Item)
			}
		}
		return items, nil
	}

	eventQueueRepo.DeleteQueuedEventFunc = func(eventID string) error {
		return nil
	}

	eventRepo.GetEventsFunc = func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
		// first event is not found
		if *filter.ID == event1.ID {
			return []models.Event{}, nil
		}
		// fetching for second event fails
		if *filter.ID == event2.ID {
			return nil, fmt.Errorf("error")
		}
		return []models.Event{{ID: *filter.ID, Specversion: "1.0"}}, nil
	}

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	dispatcher.Add(dispatcherEvent1)
	dispatcher.Add(dispatcherEvent2)
	dispatcher.Add(dispatcherEvent3)
	dispatcher.Run(context.Background())

	clock.Add(10 * time.Second)
	require.Equal(t, 1, len(eventSender.SentEvents))
	require.Equal(t, 1, len(eventQueueRepo.DeleteQueuedEventCalls()))
	require.Equal(t, event3.ID, eventQueueRepo.DeleteQueuedEventCalls()[0].EventID)
}
