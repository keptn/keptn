package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"

	"github.com/keptn/keptn/shipyard-controller/common"
	dbmock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

func Test_WhenTimeOfEventIsOlder_EventIsSentImmediately(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			if filter.CurrentTriggeredID != "" {
				return []models.SequenceExecution{
					{
						ID: "my-id",
						Status: models.SequenceExecutionStatus{
							State: apimodels.SequenceTriggeredState,
						},
					},
				}, nil
			}
			return nil, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}
	eventSender := &fake.EventSender{}
	mockClock := clock.NewMock()

	mockClock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          10 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}
	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event), TimeStamp: timeBefore}

	err := dispatcher.Add(dispatcherEvent, false)
	require.Nil(t, err)
	require.Equal(t, 1, len(eventSender.SentEvents))
}

func Test_WhenTimeOfEventIsOlder_EventIsSentImmediatelyButSequenceIsPaused(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return true
		},
		QueueEventFunc: func(item models.QueueItem) error {
			return nil
		},
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			if filter.CurrentTriggeredID != "" {
				return []models.SequenceExecution{
					{
						ID: "my-id",
						Status: models.SequenceExecutionStatus{
							State: apimodels.SequencePaused,
						},
					},
				}, nil
			}
			return nil, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return true
		},
	}

	eventSender := &fake.EventSender{}
	mockClock := clock.NewMock()

	mockClock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          10 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}
	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event), TimeStamp: timeBefore}

	err := dispatcher.Add(dispatcherEvent, false)

	require.Nil(t, err)
	require.Empty(t, eventSender.SentEvents)

	require.Len(t, eventQueueRepo.QueueEventCalls(), 1)
}

func Test_EventIsSentImmediatelyButOtherSequenceIsRunning(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
		QueueEventFunc: func(item models.QueueItem) error {
			return nil
		},
		GetQueuedEventsFunc: func(timestamp time.Time) ([]models.QueueItem, error) {
			return nil, nil
		},
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			return []models.SequenceExecution{
				{
					ID:       "my-task-sequence-execution-id",
					Sequence: keptnv2.Sequence{},
					Status: models.SequenceExecutionStatus{
						State:         apimodels.SequenceStartedState,
						PreviousTasks: nil,
						CurrentTask:   models.TaskExecutionState{},
					},
					Scope: models.EventScope{
						EventData: keptnv2.EventData{
							Project: "my-project",
							Stage:   "my-stage",
							Service: "my-service",
						},
						KeptnContext: "my-other-context-id",
					},
				},
			}, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	eventSender := &fake.EventSender{}
	mockClock := clock.NewMock()

	mockClock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          10 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}
	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event), TimeStamp: timeBefore}

	err := dispatcher.Add(dispatcherEvent, false)

	require.Nil(t, err)
	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Len(t, eventQueueRepo.QueueEventCalls(), 1)
}

func Test_EventIsSentImmediatelyButOtherSequenceIsRunningDifferentService(t *testing.T) {

	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)
	eventSender := &fake.EventSender{
		SentEvents: []event.Event{},
	}
	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
		QueueEventFunc: func(item models.QueueItem) error {
			eventSender.SentEvents = append(eventSender.SentEvents, event.Event{})
			return nil
		},
		GetQueuedEventsFunc: func(timestamp time.Time) ([]models.QueueItem, error) {
			return nil, nil
		},
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			return []models.SequenceExecution{}, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	mockClock := clock.NewMock()

	mockClock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          15 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}
	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-otherservice",
	}
	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	event.Shkeptncontext = "my-context-id"
	event.Triggeredid = "my-otherid"
	dispatcherEvent := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event), TimeStamp: timeAfter}

	err := dispatcher.Add(dispatcherEvent, false)

	require.Nil(t, err)
	require.Equal(t, 1, len(eventSender.SentEvents))
	require.Len(t, eventQueueRepo.QueueEventCalls(), 1)
}

func Test_EventIsSentImmediatelyAndOtherSequenceIsRunningButIsPaused(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
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
						State: apimodels.SequenceStartedState, // overall sequence is running
					},
					{
						Scope: models.EventScope{
							KeptnContext: "my-other-context-id",
							EventData:    keptnv2.EventData{Stage: "my-stage"}, // but in this stage, it has been paused
						},
						State: apimodels.SequencePaused,
					},
				}, nil
			}
			return nil, nil
		},
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			if filter.Scope.KeptnContext == "my-context-id" {
				return []models.SequenceExecution{
					{
						ID:       "",
						Sequence: keptnv2.Sequence{},
						Status: models.SequenceExecutionStatus{
							State: apimodels.SequenceStartedState,
						},
						Scope:           models.EventScope{},
						InputProperties: nil,
					},
				}, nil
			}
			return nil, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			if eventScope.KeptnContext == "my-other-context-id" {
				return true
			}
			return false
		},
	}

	eventSender := &fake.EventSender{}
	mockClock := clock.NewMock()

	mockClock.Set(timeAfter)

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          10 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}
	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	event.Shkeptncontext = "my-context-id"
	dispatcherEvent := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event), TimeStamp: timeBefore}

	err := dispatcher.Add(dispatcherEvent, false)

	require.Nil(t, err)
	require.Len(t, eventSender.SentEvents, 1)
	require.Len(t, eventQueueRepo.QueueEventCalls(), 0)
}

func Test_WhenTimeOfEventIsYounger_EventIsQueued(t *testing.T) {

	timeBefore := time.Date(2021, 4, 21, 15, 00, 00, 0, time.UTC)
	timeAfter := time.Date(2021, 4, 21, 15, 00, 00, 1, time.UTC)

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{}
	eventSender := &fake.EventSender{}
	mockClock := clock.NewMock()

	eventQueueRepo.QueueEventFunc = func(item models.QueueItem) error {
		return nil
	}
	mockClock.Set(timeBefore)

	dispatcher := EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		theClock:       mockClock,
		syncInterval:   10 * time.Second,
	}

	data := keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}

	event, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).Build()
	dispatcherEvent := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event), TimeStamp: timeAfter}

	err := dispatcher.Add(dispatcherEvent, false)

	require.Nil(t, err)
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

	event1, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task"), "source", data).WithKeptnContext("my-context").Build()
	dispatcherEvent1 := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event1), TimeStamp: timeAfter1}
	event2, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task2"), "source", data).WithKeptnContext("my-context").Build()
	dispatcherEvent2 := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event2), TimeStamp: timeAfter2}
	event3, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task3"), "source", data).WithKeptnContext("my-context").Build()
	dispatcherEvent3 := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event3), TimeStamp: timeAfter3}

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
	}
	eventSender := &fake.EventSender{}

	mockClock := clock.NewMock()
	mockClock.Set(timeNow)

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

	eventRepo.GetEventsFunc = func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
		return []apimodels.KeptnContextExtendedCE{{Shkeptncontext: "my-context", ID: *filter.ID, Specversion: "1.0", Source: stringp("source"), Type: stringp("my-type"), Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		}}}, nil
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			if filter.Scope.KeptnContext == "my-context" {
				return []models.SequenceExecution{
					{
						Status: models.SequenceExecutionStatus{
							State: apimodels.SequenceStartedState,
						},
					},
				}, nil
			}
			return nil, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          10 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}

	_ = dispatcher.Add(dispatcherEvent1, false)
	_ = dispatcher.Add(dispatcherEvent2, false)
	_ = dispatcher.Add(dispatcherEvent3, false)

	require.Equal(t, 0, len(eventSender.SentEvents))
	require.Equal(t, 3, len(eventQueueRepo.QueueEventCalls()))
	dispatcher.Run(context.Background())
	mockClock.Add(9 * time.Second)
	require.Equal(t, 0, len(eventSender.SentEvents))
	mockClock.Add(2 * time.Second)
	require.Eventually(t, func() bool {
		return len(eventSender.SentEvents) == 3
	}, 5*time.Second, 1*time.Second)
	// check if the time stamp of the event has been set to the time at which it has been sent
	require.WithinDuration(t, mockClock.Now().UTC(), eventSender.SentEvents[0].Time(), time.Second)
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
	dispatcherEvent1 := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event1), TimeStamp: timeAfter1}
	event2, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task2"), "source", data).Build()
	dispatcherEvent2 := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event2), TimeStamp: timeAfter2}
	event3, _ := keptnv2.KeptnEvent(keptnv2.GetStartedEventType("task3"), "source", data).Build()
	dispatcherEvent3 := models.DispatcherEvent{Event: keptnv2.ToCloudEvent(event3), TimeStamp: timeAfter3}

	eventRepo := &dbmock.EventRepoMock{}
	eventQueueRepo := &dbmock.EventQueueRepoMock{
		GetEventQueueSequenceStatesFunc: func(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
			return nil, nil
		},
	}
	eventSender := &fake.EventSender{}

	mockClock := clock.NewMock()
	mockClock.Set(timeNow)

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

	eventRepo.GetEventsFunc = func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
		// first event is not found
		if *filter.ID == event1.ID {
			return []apimodels.KeptnContextExtendedCE{}, nil
		}
		// fetching for second event fails
		if *filter.ID == event2.ID {
			return nil, fmt.Errorf("error")
		}
		return []apimodels.KeptnContextExtendedCE{{ID: *filter.ID, Specversion: "1.0", Source: stringp("source"), Type: stringp("my-type"), Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		}}}, nil
	}

	sequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			if filter.CurrentTriggeredID == "" {
				return nil, nil
			}
			return []models.SequenceExecution{
				{
					Status: models.SequenceExecutionStatus{
						State: apimodels.SequenceStartedState,
					},
				},
			}, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	dispatcher := EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		eventSender:           eventSender,
		theClock:              mockClock,
		syncInterval:          10 * time.Second,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}

	_ = dispatcher.Add(dispatcherEvent1, false)
	_ = dispatcher.Add(dispatcherEvent2, false)
	_ = dispatcher.Add(dispatcherEvent3, false)
	dispatcher.Run(context.Background())

	mockClock.Add(10 * time.Second)
	require.Eventually(t, func() bool {
		return len(eventSender.SentEvents) == 1
	}, 5*time.Second, 1*time.Second)
	require.Equal(t, 1, len(eventQueueRepo.DeleteQueuedEventCalls()))
	require.Equal(t, event3.ID, eventQueueRepo.DeleteQueuedEventCalls()[0].EventID)
}

func TestEventDispatcher_OnSequenceFinished(t *testing.T) {
	eventQueueRepo := &dbmock.EventQueueRepoMock{
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

	dispatcher.OnSequenceFinished(apimodels.KeptnContextExtendedCE{Shkeptncontext: "my-context"})

	require.Len(t, eventQueueRepo.DeleteEventQueueStatesCalls(), 1)
	require.Len(t, eventQueueRepo.DeleteQueuedEventsCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.DeleteEventQueueStatesCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, "my-context", eventQueueRepo.DeleteQueuedEventsCalls()[0].Scope.KeptnContext)
}

func TestEventDispatcher_OnSequenceTimeout(t *testing.T) {
	eventQueueRepo := &dbmock.EventQueueRepoMock{
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

	dispatcher.OnSequenceTimeout(apimodels.KeptnContextExtendedCE{Shkeptncontext: "my-context"})

	require.Len(t, eventQueueRepo.DeleteEventQueueStatesCalls(), 1)
	require.Len(t, eventQueueRepo.DeleteQueuedEventsCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.DeleteEventQueueStatesCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, "my-context", eventQueueRepo.DeleteQueuedEventsCalls()[0].Scope.KeptnContext)
}

func TestEventDispatcher_OnSequenceFinished_DeletingStateFailsButDeletingQueueShouldBeCalled(t *testing.T) {
	eventQueueRepo := &dbmock.EventQueueRepoMock{
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

	dispatcher.OnSequenceFinished(apimodels.KeptnContextExtendedCE{Shkeptncontext: "my-context"})

	require.Len(t, eventQueueRepo.DeleteEventQueueStatesCalls(), 1)
	require.Len(t, eventQueueRepo.DeleteQueuedEventsCalls(), 1)

	require.Equal(t, "my-context", eventQueueRepo.DeleteEventQueueStatesCalls()[0].State.Scope.KeptnContext)
	require.Equal(t, "my-context", eventQueueRepo.DeleteQueuedEventsCalls()[0].Scope.KeptnContext)
}
