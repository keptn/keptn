package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	dbmock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

func TestSequenceDispatcher(t *testing.T) {
	theClock := clock.NewMock()

	startSequenceCalls := []apimodels.KeptnContextExtendedCE{}
	triggeredEvents := []apimodels.KeptnContextExtendedCE{
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

	mockQueue := []models.QueueItem{}

	mockEventRepo := &dbmock.EventRepoMock{
		GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
			return triggeredEvents, nil
		},
	}

	currentSequenceExecutions := []models.SequenceExecution{}

	mockSequenceQueueRepo := &dbmock.SequenceQueueRepoMock{
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

	mockSequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			return currentSequenceExecutions, nil
		},
		GetByTriggeredIDFunc: func(project string, triggeredID string) (*models.SequenceExecution, error) {
			return &models.SequenceExecution{
				ID: "my-id",
				Status: models.SequenceExecutionStatus{
					State: apimodels.SequenceTriggeredState,
				},
			}, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	sequenceDispatcher := handler.NewSequenceDispatcher(mockEventRepo, mockSequenceQueueRepo, mockSequenceExecutionRepo, 10*time.Second, theClock, common.SDModeRW)

	sequenceDispatcher.Run(context.Background(), common.SDModeRW, func(event apimodels.KeptnContextExtendedCE) error {
		startSequenceCalls = append(startSequenceCalls, event)
		return nil
	})

	// check if repos are queried
	theClock.Add(11 * time.Second)
	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 1)
	// since no elements have been added to the queue yet, the other repos should not have been queried at this point
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
	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 1)
	require.Equal(t, mockSequenceExecutionRepo.GetCalls()[0].Filter.Scope.Project, queueItem.Scope.Project)
	require.Equal(t, mockSequenceExecutionRepo.GetCalls()[0].Filter.Scope.Stage, queueItem.Scope.Stage)

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

	// dispatch another task sequence - this one should start since the service is different
	queueItemPar := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-other-service",
			},
			KeptnContext: "my-other-context-id",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
			WrappedEvent: apimodels.KeptnContextExtendedCE{
				Type: common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
			},
		},
		EventID: "my-event-other-id",
	}
	err = sequenceDispatcher.Add(queueItemPar)

	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 2)

	// GetEvents and DeleteQueuedSequences should have been called again at this point
	require.Len(t, mockEventRepo.GetEventsCalls(), 2)
	require.Len(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls(), 2)

	// item should not have been added to queue
	require.Len(t, mockSequenceQueueRepo.QueueSequenceCalls(), 0)

	// now we have a sequence running
	currentSequenceExecutions = append(currentSequenceExecutions, models.SequenceExecution{
		ID: "my-id",
		Sequence: keptnv2.Sequence{
			Name: "delivery",
		},
		Status: models.SequenceExecutionStatus{
			State: apimodels.SequenceStartedState,
		},
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-id",
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
			WrappedEvent: apimodels.KeptnContextExtendedCE{
				Type: common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
			},
		},
		EventID: "my-event-id-2",
	}
	err = sequenceDispatcher.Add(queueItem2)

	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 3)

	// GetEvents and DeleteQueuedSequences should not have been called again at this point
	require.Len(t, mockEventRepo.GetEventsCalls(), 2)
	require.Len(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls(), 2)

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
			WrappedEvent: apimodels.KeptnContextExtendedCE{
				Type: common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
			},
		},
		EventID: "my-event-id-3",
	}
	err = sequenceDispatcher.Add(queueItem3)
	// no new call to check sequences because Read disabled
	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 3)
	//new item should have been added to queue
	require.Len(t, mockSequenceQueueRepo.QueueSequenceCalls(), 2)

}

func TestSequenceDispatcher_Remove(t *testing.T) {
	mockSequenceQueueRepo := &dbmock.SequenceQueueRepoMock{
		DeleteQueuedSequencesFunc: func(itemFilter models.QueueItem) error {
			return nil
		},
	}

	sequenceDispatcher := handler.NewSequenceDispatcher(nil, mockSequenceQueueRepo, nil, 10*time.Second, nil, common.SDModeRW)

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
	theClock := clock.NewMock()
	sequencePaused := false
	currentSequenceExecutions := []models.SequenceExecution{{
		ID: "my-id",
		Sequence: keptnv2.Sequence{
			Name: "delivery",
		},
		Status: models.SequenceExecutionStatus{
			State: apimodels.SequenceStartedState,
		},
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-id",
		},
	}}
	startSequenceCalls := []apimodels.KeptnContextExtendedCE{}
	mockQueue := []models.QueueItem{}
	triggeredEvents := []apimodels.KeptnContextExtendedCE{
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

	mockEventRepo := &dbmock.EventRepoMock{
		GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
			return triggeredEvents, nil
		},
	}

	mockSequenceQueueRepo := &dbmock.SequenceQueueRepoMock{
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

	mockSequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			return currentSequenceExecutions, nil
		},
		GetByTriggeredIDFunc: func(project string, triggeredID string) (*models.SequenceExecution, error) {
			return &models.SequenceExecution{
				ID: "my-id",
				Status: models.SequenceExecutionStatus{
					State: apimodels.SequenceTriggeredState,
				},
			}, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return sequencePaused
		},
	}

	sequenceDispatcher := handler.NewSequenceDispatcher(mockEventRepo, mockSequenceQueueRepo, mockSequenceExecutionRepo, 10*time.Second, theClock, common.SDModeRW)

	sequenceDispatcher.Run(context.Background(), common.SDModeRW, func(event apimodels.KeptnContextExtendedCE) error {
		startSequenceCalls = append(startSequenceCalls, event)
		return nil
	})

	// test failure in branch blocked
	queueItem := getQueueItem("myid1")
	err := sequenceDispatcher.Add(queueItem)
	require.Error(t, err, "could not append item!")

	theClock.Add(11 * time.Second)

	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 2)
	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 1)
	require.Empty(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls())

	// test failure in branch pause
	queueItem2 := getQueueItem("myid2")
	sequencePaused = true
	err2 := sequenceDispatcher.Add(queueItem2)
	require.Error(t, err2, "could not append item!")
	theClock.Add(11 * time.Second)
	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 4)
	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 1)
	require.Empty(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls())

}

func TestSequenceDispatcher_QueueIsNotEmpty(t *testing.T) {
	theClock := clock.NewMock()

	startSequenceCalls := []apimodels.KeptnContextExtendedCE{}
	triggeredEvents := []apimodels.KeptnContextExtendedCE{
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

	mockQueue := []models.QueueItem{}

	mockEventRepo := &dbmock.EventRepoMock{
		GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
			return triggeredEvents, nil
		},
	}

	currentSequenceExecutions := []models.SequenceExecution{}

	mockSequenceQueueRepo := &dbmock.SequenceQueueRepoMock{
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

	mockSequenceExecutionRepo := &dbmock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			return currentSequenceExecutions, nil
		},
		GetByTriggeredIDFunc: func(project string, triggeredID string) (*models.SequenceExecution, error) {
			return &models.SequenceExecution{
				ID: "my-id",
				Status: models.SequenceExecutionStatus{
					State: apimodels.SequenceTriggeredState,
				},
			}, nil
		},
		IsContextPausedFunc: func(eventScope models.EventScope) bool {
			return false
		},
	}

	sequenceDispatcher := handler.NewSequenceDispatcher(mockEventRepo, mockSequenceQueueRepo, mockSequenceExecutionRepo, 10*time.Second, theClock, common.SDModeRW)

	sequenceDispatcher.Run(context.Background(), common.SDModeRW, func(event apimodels.KeptnContextExtendedCE) error {
		startSequenceCalls = append(startSequenceCalls, event)
		return nil
	})

	// check if repos are queried
	theClock.Add(11 * time.Second)
	// queue repo should have been queried
	require.Len(t, mockSequenceQueueRepo.GetQueuedSequencesCalls(), 1)
	// since no elements have been added to the queue yet, the other repos should not have been queried at this point// since one element has been added to the queue, the other sequence previously in the queue was dispatched
	require.Len(t, mockSequenceQueueRepo.DeleteQueuedSequencesCalls(), 0)

	mockQueue = append(mockQueue, models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-id2",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID: "my-event-id2",
	})

	// now, let's add a sequence to the queue - should not be started immediately since there is another queued sequence
	queueItem := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context-id2",
			EventType:    keptnv2.GetTriggeredEventType("dev.delivery"),
		},
		EventID: "my-event-id2",
	}
	err := sequenceDispatcher.Add(queueItem)
	require.Equal(t, err.Error(), "sequence is currently blocked by waiting for another sequence to end")
	require.Len(t, mockSequenceExecutionRepo.GetCalls(), 0)

	require.Len(t, mockEventRepo.GetEventsCalls(), 0)

	// the sequence should be inserted into the queue
	require.Len(t, mockSequenceQueueRepo.QueueSequenceCalls(), 1)
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
