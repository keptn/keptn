package handler

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencedispatcher.go . ISequenceDispatcher
// ISequenceDispatcher is responsible for dispatching events to be sent to the event broker
type ISequenceDispatcher interface {
	Add(queueItem models.QueueItem) error
	Run(ctx context.Context, startSequenceFunc func(event models.Event) error)
	Remove(eventScope models.EventScope) error
}

type SequenceDispatcher struct {
	eventRepo          db.EventRepo
	eventQueueRepo     db.EventQueueRepo
	sequenceQueue      db.SequenceQueueRepo
	sequenceRepo       db.TaskSequenceRepo
	theClock           clock.Clock
	syncInterval       time.Duration
	startSequenceFunc  func(event models.Event) error
	shipyardController shipyardController
	mutex              sync.Mutex
}

// NewSequenceDispatcher creates a new SequenceDispatcher
func NewSequenceDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	sequenceRepo db.TaskSequenceRepo,
	syncInterval time.Duration,
	theClock clock.Clock,

) ISequenceDispatcher {
	return &SequenceDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceQueue:  sequenceQueueRepo,
		sequenceRepo:   sequenceRepo,
		theClock:       theClock,
		syncInterval:   syncInterval,
		mutex:          sync.Mutex{},
	}
}

func (sd *SequenceDispatcher) Add(queueItem models.QueueItem) error {
	// try to dispatch the sequence immediately
	if err := sd.dispatchSequence(queueItem); err != nil {
		if err == ErrSequenceBlocked {
			// if the sequence is currently blocked, insert it into the queue
			if err2 := sd.sequenceQueue.QueueSequence(queueItem); err2 != nil {
				return err2
			}
		} else {
			return err
		}
	}
	return nil
}

func (sd *SequenceDispatcher) Remove(eventScope models.EventScope) error {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	return sd.sequenceQueue.DeleteQueuedSequences(models.QueueItem{
		Scope: eventScope,
	})
}

func (sd *SequenceDispatcher) Run(ctx context.Context, startSequenceFunc func(event models.Event) error) {
	ticker := sd.theClock.Ticker(sd.syncInterval)
	sd.startSequenceFunc = startSequenceFunc
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("cancelling sequence dispatcher loop")
				return
			case <-ticker.C:
				log.Debugf("%.2f seconds have passed. Dispatching sequences", sd.syncInterval.Seconds())
				sd.dispatchSequences()
			}
		}
	}()
}

func (sd *SequenceDispatcher) dispatchSequences() {
	queuedSequences, err := sd.sequenceQueue.GetQueuedSequences()
	if err != nil {
		if err == db.ErrNoEventFound {
			// if no sequences are in the queue, we can return here
			return
		}
		log.WithError(err).Error("could not load queued sequences")
		return
	}

	for _, queuedSequence := range queuedSequences {
		if err := sd.dispatchSequence(queuedSequence); err != nil {
			log.WithError(err).Errorf("could not dispatch sequence with keptnContext %s", queuedSequence.EventID)
		}
	}
}

func (sd *SequenceDispatcher) dispatchSequence(queuedSequence models.QueueItem) error {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()
	// first, check if the sequence is currently paused
	if sd.eventQueueRepo.IsSequenceOfEventPaused(queuedSequence.Scope) {
		log.Infof("Sequence %s is currently paused. Will not start it yet.", queuedSequence.Scope.KeptnContext)
		return ErrSequenceBlocked
	}
	// fetch all sequences that are currently running in the stage of the project where the sequence should run
	taskExecutions, err := sd.sequenceRepo.GetTaskExecutions(queuedSequence.Scope.Project, models.TaskExecution{
		Stage:   queuedSequence.Scope.Stage,
		Service: queuedSequence.Scope.Service,
	})
	if err != nil {
		return err
	}

	// if there is a sequence running in the stage, we cannot trigger this sequence yet
	if sd.areActiveSequencesBlockingQueuedSequences(taskExecutions) {
		log.Infof("sequence %s cannot be started yet because sequences are still running in stage %s", queuedSequence.Scope.KeptnContext, queuedSequence.Scope.Stage)
		return ErrSequenceBlocked
	}

	events, err := sd.eventRepo.GetEvents(queuedSequence.Scope.Project, common.EventFilter{
		ID: &queuedSequence.EventID,
	}, common.TriggeredEvent)

	if err != nil {
		return err
	}

	if len(events) == 0 {
		return fmt.Errorf("sequence.triggered event with ID %s cannot be found anymore", queuedSequence.EventID)
	}

	sequenceTriggeredEvent := events[0]

	if err := sd.startSequenceFunc(sequenceTriggeredEvent); err != nil {
		return fmt.Errorf("could not start task sequence %s: %s", queuedSequence.EventID, err.Error())
	}

	return sd.sequenceQueue.DeleteQueuedSequences(queuedSequence)
}

func (sd *SequenceDispatcher) areActiveSequencesBlockingQueuedSequences(sequenceTasks []models.TaskExecution) bool {
	if len(sequenceTasks) == 0 {
		// if there is no sequence currently running, we do not need to block
		return false
	}

	tasksGroupedByContext := groupSequenceMappingsByContext(sequenceTasks)

	for _, tasksOfContext := range tasksGroupedByContext {
		lastTaskOfSequence := getLastTaskOfSequence(tasksOfContext)
		// first, check if the other sequence is currently paused
		if sd.eventQueueRepo.IsSequenceOfEventPaused(
			models.EventScope{
				KeptnContext: lastTaskOfSequence.KeptnContext,
				EventData:    keptnv2.EventData{Stage: lastTaskOfSequence.Stage},
			}) {
			// if it is indeed paused, no need to consider it
			continue
		}
		if lastTaskOfSequence.Task.Name != keptnv2.ApprovalTaskName {
			// if there is a sequence running that is not waiting for an approval, we need to block
			return true
		}
	}
	// do not block if all active sequences are currently handling an approval task
	return false
}

func groupSequenceMappingsByContext(sequenceTasks []models.TaskExecution) map[string][]models.TaskExecution {
	result := map[string][]models.TaskExecution{}
	for index := range sequenceTasks {
		result[sequenceTasks[index].KeptnContext] = append(result[sequenceTasks[index].KeptnContext], sequenceTasks[index])
	}
	return result
}

func getLastTaskOfSequence(sequenceTasks []models.TaskExecution) models.TaskExecution {
	lastTask := models.TaskExecution{
		Task: models.Task{TaskIndex: -1},
	}
	for index := range sequenceTasks {
		if sequenceTasks[index].Task.TaskIndex > lastTask.Task.TaskIndex {
			lastTask = sequenceTasks[index]
		}
	}

	return lastTask
}
