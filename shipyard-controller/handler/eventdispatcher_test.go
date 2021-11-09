package handler

import (
	"context"
	"errors"
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
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}
	sequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return []models.TaskExecution{
				{
					TaskSequenceName: "delivery",
					TriggeredEventID: "my-triggered-id",
					Task: models.Task{
						Task: keptnv2.Task{
							Name: "deployment",
						},
						TaskIndex: 0,
					},
					Stage:        "my-stage",
					Service:      "my-service",
					KeptnContext: "my-context-id",
				},
			}, nil
		},
	}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()

	clock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceRepo:   sequenceRepo,
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
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeBefore}

	dispatcher.Add(dispatcherEvent, false)
	require.Equal(t, 1, len(eventSender.SentEvents))
}

func Test_WhenTimeOfEventIsOlder_EventIsSentImmediatelyButSequenceIsPaused(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return true
		},
		QueueEventFunc: func(item models.QueueItem) error {
			return nil
		},
	}
	sequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return []models.TaskExecution{
				{
					TaskSequenceName: "delivery",
					TriggeredEventID: "my-triggered-id",
					Task: models.Task{
						Task: keptnv2.Task{
							Name: "deployment",
						},
						TaskIndex: 0,
					},
					Stage:        "my-stage",
					Service:      "my-service",
					KeptnContext: "my-context-id",
				},
			}, nil
		},
	}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()

	clock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceRepo:   sequenceRepo,
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
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeBefore}

	dispatcher.Add(dispatcherEvent, false)
	require.Empty(t, eventSender.SentEvents)

	require.Len(t, eventQueueRepo.QueueEventCalls(), 1)
}

func Test_EventIsSentImmediatelyButOtherSequenceIsRunning(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		QueueEventFunc: func(item models.QueueItem) error {
			return nil
		},
		GetQueuedEventsFunc: func(timestamp time.Time) ([]models.QueueItem, error) {
			return nil, nil
		},
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}
	sequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return []models.TaskExecution{
				{
					TaskSequenceName: "delivery",
					TriggeredEventID: "my-triggered-id",
					Task: models.Task{
						Task: keptnv2.Task{
							Name: "deployment",
						},
						TaskIndex: 0,
					},
					Stage:        "my-stage",
					Service:      "my-service",
					KeptnContext: "my-other-context-id",
				},
			}, nil
		},
	}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()

	clock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceRepo:   sequenceRepo,
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
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeBefore}

	dispatcher.Add(dispatcherEvent, false)
	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Len(t, eventQueueRepo.QueueEventCalls(), 1)
}

func Test_EventIsSentImmediatelyAndOtherSequenceIsRunningButIsPaused(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &db_mock.EventRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		QueueEventFunc: func(item models.QueueItem) error {
			return nil
		},
		GetQueuedEventsFunc: func(timestamp time.Time) ([]models.QueueItem, error) {
			return nil, nil
		},
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			if filter.Scope.KeptnContext == "my-other-context-id" {
				return []models.EventQueueSequenceState{
					{
						Scope: models.EventScope{
							KeptnContext: "my-other-context-id",
						},
						State: models.SequenceStartedState, // overall sequence is running
					},
					{
						Scope: models.EventScope{
							KeptnContext: "my-other-context-id",
							EventData:    keptnv2.EventData{Stage: "my-stage"}, // but in this stage, it has been paused
						},
						State: models.SequencePaused,
					},
				}, nil
			}
			return nil, nil
		},
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			if eventScope.KeptnContext == "my-other-context-id" {
				return true
			}
			return false
		},
	}
	sequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return []models.TaskExecution{
				{
					TaskSequenceName: "delivery",
					TriggeredEventID: "my-triggered-id",
					Task: models.Task{
						Task: keptnv2.Task{
							Name: "deployment",
						},
						TaskIndex: 0,
					},
					Stage:        "my-stage",
					Service:      "my-service",
					KeptnContext: "my-other-context-id",
				},
			}, nil
		},
	}
	eventSender := &fake.EventSender{}
	clock := clock.NewMock()

	clock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceRepo:   sequenceRepo,
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
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{keptnv2.ToCloudEvent(event), timeBefore}

	dispatcher.Add(dispatcherEvent, false)
	require.Len(t, eventSender.SentEvents, 1)
	require.Len(t, eventQueueRepo.QueueEventCalls(), 0)
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

	dispatcher.Add(dispatcherEvent, false)
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
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}
	eventSender := &fake.EventSender{}
	sequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return []models.TaskExecution{}, nil
		},
	}
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
		return []models.Event{{ID: *filter.ID, Specversion: "1.0", Source: stringp("source"), Type: stringp("my-type"), Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		}}}, nil
	}

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		sequenceRepo:   sequenceRepo,
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	dispatcher.Add(dispatcherEvent1, false)
	dispatcher.Add(dispatcherEvent2, false)
	dispatcher.Add(dispatcherEvent3, false)

	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Equal(t, 3, len(eventQueueRepo.QueueEventCalls()))
	dispatcher.Run(context.Background())
	clock.Add(9 * time.Second)
	require.Equal(t, 0, len(eventSender.SentEvents))
	clock.Add(2 * time.Second)
	require.Eventually(t, func() bool {
		return len(eventSender.SentEvents) == 3
	}, 5*time.Second, 1*time.Second)
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
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}
	eventSender := &fake.EventSender{}
	sequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return []models.TaskExecution{}, nil
		},
	}
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
		return []models.Event{{ID: *filter.ID, Specversion: "1.0", Source: stringp("source"), Type: stringp("my-type"), Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		}}}, nil
	}

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceRepo:   sequenceRepo,
		eventSender:    eventSender,
		theClock:       clock,
		syncInterval:   10 * time.Second,
	}

	dispatcher.Add(dispatcherEvent1, false)
	dispatcher.Add(dispatcherEvent2, false)
	dispatcher.Add(dispatcherEvent3, false)
	dispatcher.Run(context.Background())

	clock.Add(10 * time.Second)
	require.Eventually(t, func() bool {
		return len(eventSender.SentEvents) == 1
	}, 5*time.Second, 1*time.Second)
	require.Equal(t, 1, len(eventQueueRepo.DeleteQueuedEventCalls()))
	require.Equal(t, event3.ID, eventQueueRepo.DeleteQueuedEventCalls()[0].EventID)
}

func TestEventDispatcher_OnSequencePaused(t *testing.T) {
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		CreateOrUpdateEventQueueStateFunc: func(state models.EventQueueSequenceState) error {
			return nil
		},
	}

	dispatcher := EventDispatcher{
		eventQueueRepo: eventQueueRepo,
	}

	dispatcher.OnSequencePaused(models.EventScope{KeptnContext: "my-context"})

	require.Len(t, eventQueueRepo.CreateOrUpdateEventQueueStateCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.CreateOrUpdateEventQueueStateCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, models.SequencePaused, eventQueueRepo.CreateOrUpdateEventQueueStateCalls()[0].State.State)
}

func TestEventDispatcher_OnSequenceResumed(t *testing.T) {
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		CreateOrUpdateEventQueueStateFunc: func(state models.EventQueueSequenceState) error {
			return nil
		},
	}

	dispatcher := EventDispatcher{
		eventQueueRepo: eventQueueRepo,
	}

	dispatcher.OnSequenceResumed(models.EventScope{KeptnContext: "my-context"})

	require.Len(t, eventQueueRepo.CreateOrUpdateEventQueueStateCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.CreateOrUpdateEventQueueStateCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, models.SequenceStartedState, eventQueueRepo.CreateOrUpdateEventQueueStateCalls()[0].State.State)
}

func TestEventDispatcher_OnSequenceFinished(t *testing.T) {
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		DeleteEventQueueStatesFunc: func(state models.EventQueueSequenceState) error {
			return nil
		},
		DeleteQueuedEventsFunc: func(scope models.EventScope) error {
			return nil
		},
	}

	dispatcher := EventDispatcher{
		eventQueueRepo: eventQueueRepo,
	}

	dispatcher.OnSequenceFinished(models.Event{Shkeptncontext: "my-context"})

	require.Len(t, eventQueueRepo.DeleteEventQueueStatesCalls(), 1)
	require.Len(t, eventQueueRepo.DeleteQueuedEventsCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.DeleteEventQueueStatesCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, "my-context", eventQueueRepo.DeleteQueuedEventsCalls()[0].Scope.KeptnContext)
}

func TestEventDispatcher_OnSequenceTimeout(t *testing.T) {
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		DeleteEventQueueStatesFunc: func(state models.EventQueueSequenceState) error {
			return nil
		},
		DeleteQueuedEventsFunc: func(scope models.EventScope) error {
			return nil
		},
	}

	dispatcher := EventDispatcher{
		eventQueueRepo: eventQueueRepo,
	}

	dispatcher.OnSequenceTimeout(models.Event{Shkeptncontext: "my-context"})

	require.Len(t, eventQueueRepo.DeleteEventQueueStatesCalls(), 1)
	require.Len(t, eventQueueRepo.DeleteQueuedEventsCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.DeleteEventQueueStatesCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, "my-context", eventQueueRepo.DeleteQueuedEventsCalls()[0].Scope.KeptnContext)
}

func TestEventDispatcher_OnSequenceFinished_DeletingStateFailsButDeletingQueueShouldBeCalled(t *testing.T) {
	eventQueueRepo := &db_mock.EventQueueRepoMock{
		DeleteEventQueueStatesFunc: func(state models.EventQueueSequenceState) error {
			return errors.New("oops")
		},
		DeleteQueuedEventsFunc: func(scope models.EventScope) error {
			return nil
		},
	}

	dispatcher := EventDispatcher{
		eventQueueRepo: eventQueueRepo,
	}

	dispatcher.OnSequenceFinished(models.Event{Shkeptncontext: "my-context"})

	require.Len(t, eventQueueRepo.DeleteEventQueueStatesCalls(), 1)
	require.Len(t, eventQueueRepo.DeleteQueuedEventsCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.DeleteEventQueueStatesCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, "my-context", eventQueueRepo.DeleteQueuedEventsCalls()[0].Scope.KeptnContext)
}
