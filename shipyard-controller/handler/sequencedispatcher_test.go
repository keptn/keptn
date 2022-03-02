package handler_test

import (
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSequenceDispatcher(t *testing.T) {
	theClock := clock.NewMock()

	startSequenceCalls := []models.Event{}
	triggeredEvents := []models.Event{
		{
			Data: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			ID:             "my-event-id",
			Shkeptncontext: "my-context-id",
			Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
		},
	}
	currentTaskSequences := []models.TaskExecution{}
	mockQueue := []models.QueueItem{}

	mockEventRepo := &db_mock.EventRepoMock{
		GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
			return triggeredEvents, nil
		},
	}

	mockEventQueueRepo := &db_mock.EventQueueRepoMock{
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	mockSequenceQueueRepo := &db_mock.SequenceQueueRepoMock{
		QueueSequenceFunc: func(item models.QueueItem) error {
			mockQueue = append(mockQueue, item)
			return nil
		},
		GetQueuedSequencesFunc: func() ([]models.QueueItem, error) {
			return mockQueue, nil
		},
		DeleteQueuedSequencesFunc: func(itemFilter models.QueueItem) error {
			for index := range mockQueue {
				if mockQueue[index].EventID == itemFilter.EventID {
					mockQueue = append(mockQueue[:index], mockQueue[index+1:]...)
				}
			}
			return nil
		},
	}

	mockTaskSequenceRepo := &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return currentTaskSequences, nil
		},
	}

	sequenceDispatcher := handler.NewSequenceDispatcher(mockEventRepo, mockEventQueueRepo, mockSequenceQueueRepo, mockTaskSequenceRepo, 10*time.Second, theClock, common.SDModeRW)

	sequenceDispatcher.Run(context.Background(), common.SDModeRW, func(event models.Event) error {
		startSequenceCalls = append(startSequenceCalls, event)
		return nil
	})

	// check if repos are queried
	theClock.Add(11 * time.Second)
	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 1)
	// since no elements have been added to the queue yet, the other repos should not have been queried at this point
	require.Empty(t, mockTaskSequenceRepo.GetTaskExecutionsCalls())
	require.Empty(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls())

	// now, let's add a sequence to the queue - should be started immediately since no other sequences are running currently
	queueItem := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-id",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID: "my-event-id",
	}
	err := sequenceDispatcher.Add(queueItem)
	require.Nil(t, err)
	require.Len(t, mockTaskSequenceRepo.GetTaskExecutionsCalls(), 1)
	require.Equal(t, mockTaskSequenceRepo.GetTaskExecutionsCalls()[0].Project, queueItem.Scope.Project)
	require.Equal(t, mockTaskSequenceRepo.GetTaskExecutionsCalls()[0].Filter, models.TaskExecution{Stage: queueItem.Scope.Stage, Service: queueItem.Scope.Service})

	require.Len(t, mockEventRepo.GetEventsCalls(), 1)
	require.Equal(t, mockEventRepo.GetEventsCalls()[0].Project, queueItem.Scope.Project)
	require.Equal(t, *mockEventRepo.GetEventsCalls()[0].Filter.ID, queueItem.EventID)
	require.Equal(t, mockEventRepo.GetEventsCalls()[0].Status[0], common.TriggeredEvent)

	// if the sequence has been dispatched immediately, we do not need to insert it into the queue
	require.Empty(t, mockSequenceQueueRepo.QueueSequenceCalls())

	// has the queueItem been removed properly?
	require.Len(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls(), 1)
	require.Equal(t, queueItem, mockSequenceQueueRepo.DeleteQueuedSequencesCalls()[0].ItemFilter)

	require.Eventually(t, func() bool {
		return len(startSequenceCalls) == 1
	}, 5*time.Second, 1*time.Second)

	require.Equal(t, triggeredEvents[0], startSequenceCalls[0])

	// now we have a sequence running
	currentTaskSequences = append(currentTaskSequences, models.TaskExecution{
		TaskSequenceName: "delivery",
		TriggeredEventID: "my-event-id",
		Stage:            "my-stage",
		Service:          "my-service",
		KeptnContext:     "my-context-id",
		Task: models.Task{
			Task: keptnv2.Task{
				Name: "my-task",
			},
			TaskIndex: 1,
		},
	})

	// dispatch another task sequence - this one should not start since there is currently another one in progress
	queueItem2 := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-id-2",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
			WrappedEvent: models.Event{
				Type: common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
			},
		},
		EventID: "my-event-id-2",
	}
	err = sequenceDispatcher.Add(queueItem2)
	require.ErrorIs(t, err, handler.ErrSequenceBlockedWaiting)
	require.Len(t, mockTaskSequenceRepo.GetTaskExecutionsCalls(), 2)
	require.Equal(t, mockTaskSequenceRepo.GetTaskExecutionsCalls()[1].Project, queueItem.Scope.Project)
	require.Equal(t, mockTaskSequenceRepo.GetTaskExecutionsCalls()[1].Filter, models.TaskExecution{Stage: queueItem.Scope.Stage, Service: queueItem.Scope.Service})

	// GetEvents and DeleteQueuedSequences should not have been called again at this point
	require.Len(t, mockEventRepo.GetEventsCalls(), 1)
	require.Len(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls(), 1)

	// item should have been added to queue
	require.Len(t, mockSequenceQueueRepo.QueueSequenceCalls(), 1)

	// lets' now stop the dispatcher function and verify that add only queues
	// and does not call  dispatchSequence
	sequenceDispatcher.Stop()

	queueItem3 := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-otherproject",
				Stage:   "my-otherstage",
				Service: "my-otherservice",
			},
			KeptnContext: "my-context-id-3",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
			WrappedEvent: models.Event{
				Type: common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
			},
		},
		EventID: "my-event-id-3",
	}
	err = sequenceDispatcher.Add(queueItem3)
	// no new call to check sequences because Read disabled
	require.Len(t, mockTaskSequenceRepo.GetTaskExecutionsCalls(), 2)
	//new item should have been added to queue
	require.Len(t, mockSequenceQueueRepo.QueueSequenceCalls(), 2)

}

func TestSequenceDispatcher_Remove(t *testing.T) {
	mockSequenceQueueRepo := &db_mock.SequenceQueueRepoMock{
		DeleteQueuedSequencesFunc: func(itemFilter models.QueueItem) error {
			return nil
		},
	}

	sequenceDispatcher := handler.NewSequenceDispatcher(nil, nil, mockSequenceQueueRepo, nil, 10*time.Second, nil, common.SDModeRW)

	myScope := models.EventScope{
		EventData:    keptnv2.EventData{Project: "my-project"},
		KeptnContext: "my-context",
	}
	err := sequenceDispatcher.Remove(myScope)

	require.Nil(t, err)

	require.Len(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls(), 1)
	require.Equal(t, models.QueueItem{Scope: myScope}, mockSequenceQueueRepo.DeleteQueuedSequencesCalls()[0].ItemFilter)

}

func TestSequenceDispatcher_AddError(t *testing.T) {
	sequencePaused := false
	theClock := clock.NewMock()
	currentTaskSequences := []models.TaskExecution{
		{
			TaskSequenceName: "delivery",
			TriggeredEventID: "my-event-id",
			Stage:            "my-stage",
			Service:          "my-service",
			KeptnContext:     "my-context-id",
			Task: models.Task{
				Task: keptnv2.Task{
					Name: "my-task",
				},
				TaskIndex: 1,
			},
		},
	}
	startSequenceCalls := []models.Event{}
	mockQueue := []models.QueueItem{}
	triggeredEvents := []models.Event{
		{
			Data: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			ID:             "my-event-id",
			Shkeptncontext: "my-context-id",
			Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
		},
	}

	mockEventRepo := &db_mock.EventRepoMock{
		GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
			return triggeredEvents, nil
		},
	}
	mockEventQueueRepo := getEventQueueRepo(&sequencePaused)

	mockSequenceQueueRepo := &db_mock.SequenceQueueRepoMock{
		QueueSequenceFunc: func(item models.QueueItem) error {
			return errors.New("could not append item!")
		},
		GetQueuedSequencesFunc: func() ([]models.QueueItem, error) {
			return mockQueue, nil
		},
		DeleteQueuedSequencesFunc: func(itemFilter models.QueueItem) error {
			for index := range mockQueue {
				if mockQueue[index].EventID == itemFilter.EventID {
					mockQueue = append(mockQueue[:index], mockQueue[index+1:]...)
				}
			}
			return nil
		},
	}

	mockTaskSequenceRepo := getTaskSequenceRepo(currentTaskSequences)

	sequenceDispatcher := handler.NewSequenceDispatcher(mockEventRepo, mockEventQueueRepo, mockSequenceQueueRepo, mockTaskSequenceRepo, 10*time.Second, theClock, common.SDModeRW)

	sequenceDispatcher.Run(context.Background(), common.SDModeRW, func(event models.Event) error {
		startSequenceCalls = append(startSequenceCalls, event)
		return nil
	})

	// test failure in branch blocked
	queueItem := getQueueItem("myid1")
	err := sequenceDispatcher.Add(queueItem)
	require.Error(t, err, "could not append item!")

	theClock.Add(11 * time.Second)

	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 1)
	require.Len(t, mockTaskSequenceRepo.GetTaskExecutionsCalls(), 1)
	require.Empty(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls())

	// test failure in branch pause
	queueItem2 := getQueueItem("myid2")
	sequencePaused = true
	err2 := sequenceDispatcher.Add(queueItem2)
	require.Error(t, err2, "could not append item!")
	theClock.Add(11 * time.Second)
	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 2)
	require.Len(t, mockTaskSequenceRepo.GetTaskExecutionsCalls(), 1)
	require.Empty(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls())

}

func getQueueItem(id string) models.QueueItem {
	return models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: id,
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID: id,
	}
}

func getEventQueueRepo(reply *bool) *db_mock.EventQueueRepoMock {
	return &db_mock.EventQueueRepoMock{
		IsSequenceOfEventPausedFunc: func(eventScope models.EventScope) bool {
			return *reply
		},
	}
}

func getTaskSequenceRepo(currentTaskSequences []models.TaskExecution) *db_mock.TaskSequenceRepoMock {
	return &db_mock.TaskSequenceRepoMock{
		GetTaskExecutionsFunc: func(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
			return currentTaskSequences, nil
		},
	}
}
