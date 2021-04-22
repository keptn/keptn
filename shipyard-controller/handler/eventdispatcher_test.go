package handler

import (
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/lib/keptn"
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
		logger:         &keptn.Logger{},
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), nil).Build()
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
		logger:         &keptn.Logger{},
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), nil).Build()
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeAfter}

	dispatcher.Add(dispatcherEvent)
	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Equal(t, 1, len(eventQueueRepo.QueueEventCalls()))
}

func Test_WhenSyncTimeElapses_EventsAreDispatched(t *testing.T) {

	timeNow := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter1 := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)
	timeAfter2 := time.Date(2021, 4, 21, 15, 00, 00, 2, time.UTC)
	timeAfter3 := time.Date(2021, 4, 21, 15, 00, 00, 2, time.UTC)

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

	eventRepo.GetEventsFunc = func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
		return []models.Event{}, nil
	}

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		logger:         &keptn.Logger{},
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	event1, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), nil).Build()
	dispatcherEvent1 := models.DispatcherEvent{keptnv2.ToCloudEvent(event1), timeAfter1}
	event2, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task2"), nil).Build()
	dispatcherEvent2 := models.DispatcherEvent{keptnv2.ToCloudEvent(event2), timeAfter2}
	event3, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task3"), nil).Build()
	dispatcherEvent3 := models.DispatcherEvent{keptnv2.ToCloudEvent(event3), timeAfter3}

	dispatcher.Add(dispatcherEvent1)
	dispatcher.Add(dispatcherEvent2)
	dispatcher.Add(dispatcherEvent3)

	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Equal(t, 3, len(eventQueueRepo.QueueEventCalls()))
	dispatcher.Run()
	clock.Add(10 * time.Second)

}

//type EventQueueRepoDecorator struct {
//	db_mock.EventQueueRepoMock
//	QueuedEvents []models.QueueItem
//}
//
//func (e *EventQueueRepoDecorator) QueueEvent(item models.QueueItem) error {
//	e.QueuedEvents = append(e.QueuedEvents, item)
//	return e.QueueEvent(item)
//}
//
//func (e *EventQueueRepoDecorator) GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error) {
//	return e.GetQueuedEvents(timestamp)
//}
//
//func (e *EventQueueRepoDecorator) DeleteQueuedEvent(eventID string) {
//	e.DeleteQueuedEvent(eventID)
//}
//
//func (e *EventQueueRepoDecorator) DeleteQueuedEvents(scope models.EventScope) {
//	e.DeleteQueuedEvents(scope)
//}
