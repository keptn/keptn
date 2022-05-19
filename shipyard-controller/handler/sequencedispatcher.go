package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/benbjohnson/clock"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencedispatcher.go . ISequenceDispatcher
// ISequenceDispatcher is responsible for dispatching events to be sent to the event broker
type ISequenceDispatcher interface {
	Add(queueItem models.QueueItem) error
	Run(ctx context.Context, mode common.SDMode, startSequenceFunc func(event apimodels.KeptnContextExtendedCE) error)
	Remove(eventScope models.EventScope) error
	Stop()
}

type SequenceDispatcher struct {
	eventRepo             db.EventRepo
	sequenceQueue         db.SequenceQueueRepo
	sequenceExecutionRepo db.SequenceExecutionRepo
	theClock              clock.Clock
	syncInterval          time.Duration
	startSequenceFunc     func(event apimodels.KeptnContextExtendedCE) error
	shipyardController    shipyardController
	ticker                *clock.Ticker
	mode                  common.SDMode
}

// NewSequenceDispatcher creates a new SequenceDispatcher
func NewSequenceDispatcher(
	eventRepo db.EventRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	sequenceExecutionRepo db.SequenceExecutionRepo,
	syncInterval time.Duration,
	theClock clock.Clock,
	mode common.SDMode,
) ISequenceDispatcher {
	return &SequenceDispatcher{
		eventRepo:             eventRepo,
		sequenceQueue:         sequenceQueueRepo,
		sequenceExecutionRepo: sequenceExecutionRepo,
		theClock:              theClock,
		syncInterval:          syncInterval,
		mode:                  mode,
	}
}

func (sd *SequenceDispatcher) Add(queueItem models.QueueItem) error {
	if sd.mode == common.SDModeRW {
		//if there is only one shipyard we can both read and write,
		//so we try to dispatch the sequence immediately if the sequence queue is empty
		sequenceBlocked, err := sd.areSequencesInState(queueItem, []string{apimodels.SequenceTriggeredState, apimodels.SequenceStartedState})
		if err != nil || sequenceBlocked {
			return sd.addItemToQueue(queueItem)
		}
		if err := sd.dispatchSequence(queueItem); err != nil {
			if errors.Is(err, ErrSequenceBlocked) {
				//if the sequence is currently blocked, insert it into the queue
				return sd.add(queueItem)
			} else if errors.Is(err, ErrSequenceBlockedWaiting) {
				//if the sequence is currently blocked and should wait, insert it into the queue
				return sd.addItemToQueue(queueItem)
			} else {
				return err
			}
		}
		return nil
	} else {
		//if there are multiple shipyard we should only write
		return sd.add(queueItem)
	}
}

func (sd *SequenceDispatcher) addItemToQueue(queueItem models.QueueItem) error {
	if err2 := sd.add(queueItem); err2 != nil {
		return err2
	}
	return ErrSequenceBlockedWaiting
}

func (sd *SequenceDispatcher) add(queueItem models.QueueItem) error {
	return sd.sequenceQueue.QueueSequence(queueItem)
}

func (sd *SequenceDispatcher) Remove(eventScope models.EventScope) error {
	return sd.sequenceQueue.DeleteQueuedSequences(models.QueueItem{
		Scope: eventScope,
	})
}

func (sd *SequenceDispatcher) SetStartSequenceCallback(startSequenceFunc func(event apimodels.KeptnContextExtendedCE) error) {
	sd.startSequenceFunc = startSequenceFunc
}

func (sd *SequenceDispatcher) Run(ctx context.Context, mode common.SDMode, startSequenceFunc func(event apimodels.KeptnContextExtendedCE) error) {
	// at each run the dispatcher needs to know if it is a leader or not
	sd.mode = mode
	sd.ticker = sd.theClock.Ticker(sd.syncInterval)
	sd.startSequenceFunc = startSequenceFunc
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Cancelling sequence dispatcher loop")
				return
			case <-sd.ticker.C:
				log.Debugf("%.2f seconds have passed. Dispatching sequences", sd.syncInterval.Seconds())
				sd.dispatchSequences()
			}
		}
	}()
}

func (sd *SequenceDispatcher) Stop() {
	// as soon as a new leader is elected dispatcher should only write
	sd.mode = common.SDModeW
	if sd.ticker == nil {
		return
	}
	sd.ticker.Stop()
}

func (sd *SequenceDispatcher) dispatchSequences() {
	queuedSequences, err := sd.sequenceQueue.GetQueuedSequences()
	if err != nil {
		if errors.Is(err, db.ErrNoEventFound) {
			// if no sequences are in the queue, we can return here
			return
		}
		log.WithError(err).Error("Could not load queued sequences")
		return
	}

	for _, queuedSequence := range queuedSequences {
		if err := sd.dispatchSequence(queuedSequence); err != nil {
			if errors.Is(err, ErrSequenceBlocked) || errors.Is(err, ErrSequenceBlockedWaiting) {
				log.Infof("Could not dispatch sequence with keptnContext %s. Sequence is currently blocked by other sequence", queuedSequence.Scope.KeptnContext)
			} else {
				log.WithError(err).Errorf("Could not dispatch sequence with keptnContext %s", queuedSequence.Scope.KeptnContext)
			}
		}
	}
}

func (sd *SequenceDispatcher) areSequencesInState(queueItem models.QueueItem, status []string) (bool, error) {
	sequenceExecutions, err := sd.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: queueItem.Scope.Project,
				Stage:   queueItem.Scope.Stage,
				Service: queueItem.Scope.Service,
			},
		},
		Status: status,
	})

	if err != nil {
		return true, err
	}

	if len(sequenceExecutions) == 1 {
		if sequenceExecutions[0].Scope.KeptnContext != queueItem.Scope.KeptnContext {
			return true, nil
		}
	}

	if len(sequenceExecutions) > 1 {
		return true, nil
	}

	return false, nil
}

func (sd *SequenceDispatcher) dispatchSequence(queueItem models.QueueItem) error {
	// first, check if the sequence is currently paused
	sequenceExecution, err := sd.sequenceExecutionRepo.GetByTriggeredID(queueItem.Scope.Project, queueItem.EventID)
	if err != nil {
		return err
	}

	if sequenceExecution == nil {
		return ErrSequenceNotFound
	}

	if sequenceExecution.IsPaused() || sd.sequenceExecutionRepo.IsContextPaused(queueItem.Scope) {
		log.Infof("Sequence %s is currently paused. Will not start it yet.", queueItem.Scope.KeptnContext)
		return ErrSequenceBlocked
	}

	sequenceBlocked, err := sd.areSequencesInState(queueItem, []string{apimodels.SequenceStartedState})
	if err != nil {
		return err
	}

	if sequenceBlocked {
		log.Infof("Sequence %s cannot be started yet because sequences are still running in stage %s", queueItem.Scope.KeptnContext, queueItem.Scope.Stage)
		return ErrSequenceBlockedWaiting
	}

	events, err := sd.eventRepo.GetEvents(queueItem.Scope.Project, common.EventFilter{
		ID: &queueItem.EventID,
	}, common.TriggeredEvent)

	if err != nil {
		return err
	}

	if len(events) == 0 {
		return fmt.Errorf("sequence.triggered event with ID %s cannot be found anymore", queueItem.EventID)
	}

	sequenceTriggeredEvent := events[0]

	if err := sd.startSequenceFunc(sequenceTriggeredEvent); err != nil {
		return fmt.Errorf("could not start task sequence %s: %s", queueItem.EventID, err.Error())
	}

	return sd.sequenceQueue.DeleteQueuedSequences(queueItem)
}
