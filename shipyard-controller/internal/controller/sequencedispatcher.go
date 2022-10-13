package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/benbjohnson/clock"
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
	shipyardController    ShipyardController
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
	_, seqName, _, err := keptnv2.ParseSequenceEventType(queueItem.Scope.EventType)
	if err != nil {
		return fmt.Errorf("could not add event to sequence dispatcher: %w", err)
	}
	if sd.mode == common.SDModeRW {
		//if there is only one shipyard we can both read and write,
		//so we try to dispatch the sequence immediately
		if err := sd.dispatchSequence(queueItem); err != nil {
			if errors.Is(err, common.ErrSequenceBlocked) {
				//if the sequence is currently blocked, insert it into the queue
				log.
					WithFields(log.Fields{
						"source":       queueItem.Scope.EventSource,
						"keptncontext": queueItem.Scope.KeptnContext,
						"project":      queueItem.Scope.Project,
						"service":      queueItem.Scope.Service,
						"stage":        queueItem.Scope.Stage,
					}).
					Infof("[QUEUED    ] Sequence '%s' in stage '%s': %v", seqName, queueItem.Scope.Stage, err)
				return sd.add(queueItem)
			} else if errors.Is(err, common.ErrSequenceBlockedWaiting) {
				//if the sequence is currently blocked and should wait, insert it into the queue
				log.
					WithFields(log.Fields{
						"source":       queueItem.Scope.EventSource,
						"keptncontext": queueItem.Scope.KeptnContext,
						"project":      queueItem.Scope.Project,
						"service":      queueItem.Scope.Service,
						"stage":        queueItem.Scope.Stage,
					}).
					Infof("[QUEUED    ] Sequence '%s' in stage '%s': %v", seqName, queueItem.Scope.Stage, err)
				if err2 := sd.add(queueItem); err2 != nil {
					return err2
				}
				return common.ErrSequenceBlockedWaiting
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
			if errors.Is(err, common.ErrSequenceBlocked) || errors.Is(err, common.ErrSequenceBlockedWaiting) {
				log.Debugf("Could not dispatch sequence with keptnContext %s. Sequence is currently blocked by other sequence", queuedSequence.Scope.KeptnContext)
			} else {
				log.WithError(err).Errorf("Could not dispatch sequence with keptnContext %s", queuedSequence.Scope.KeptnContext)
			}
		} else {
			_, seqName, _, err := keptnv2.ParseSequenceEventType(queuedSequence.Scope.EventType)
			if err != nil {
				log.Errorf("could not parse sequence event type: %v", err)
			}
			log.
				WithFields(log.Fields{
					"source":       queuedSequence.Scope.EventSource,
					"keptncontext": queuedSequence.Scope.KeptnContext,
					"project":      queuedSequence.Scope.Project,
					"service":      queuedSequence.Scope.Service,
					"stage":        queuedSequence.Scope.Stage,
				}).
				Infof("[DISPATCHED] Sequence '%s' in stage `%s`", seqName, queuedSequence.Scope.Stage)
		}
	}
}

func (sd *SequenceDispatcher) isSequenceBlocked(queueItem models.QueueItem) (bool, string, error) {
	// searching for running sequences
	startedSequenceExecutions, err := sd.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: queueItem.Scope.Project,
				Stage:   queueItem.Scope.Stage,
				Service: queueItem.Scope.Service,
			},
		},
		Status: []string{apimodels.SequenceStartedState},
	})
	if err != nil {
		log.Errorf("Could not load started sequences for project %s, service %s, stage %s: %v", queueItem.Scope.Project, queueItem.Scope.Service, queueItem.Scope.Stage, err)
		return true, "", err
	}

	if len(startedSequenceExecutions) > 0 && queueItem.Scope.KeptnContext != startedSequenceExecutions[0].Scope.KeptnContext {
		log.Debugf("Sequence with KeptnContext %s blocked due to started sequence with KeptnContext %s in stage %s", queueItem.Scope.KeptnContext, startedSequenceExecutions[0].Scope.KeptnContext, queueItem.Scope.Stage)
		return true, startedSequenceExecutions[0].Scope.KeptnContext, nil
	}

	//searching for triggered sequences which were triggered before the actual sequence
	triggeredSequenceExecutions, err := sd.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: queueItem.Scope.Project,
				Stage:   queueItem.Scope.Stage,
				Service: queueItem.Scope.Service,
			},
		},
		Status:      []string{apimodels.SequenceTriggeredState},
		TriggeredAt: queueItem.Timestamp,
	})
	if err != nil {
		log.Errorf("Could not load triggered sequences for project %s, service %s, stage %s: %v", queueItem.Scope.Project, queueItem.Scope.Service, queueItem.Scope.Stage, err)
		return true, "", err
	}

	if len(triggeredSequenceExecutions) == 1 {
		if triggeredSequenceExecutions[0].Scope.KeptnContext != queueItem.Scope.KeptnContext {
			log.Debugf("Sequence with KeptnContext %s is blocked due to triggered sequence with KeptnContext %s in stage %s", queueItem.Scope.KeptnContext, triggeredSequenceExecutions[0].Scope.KeptnContext, queueItem.Scope.Stage)
			return true, startedSequenceExecutions[0].Scope.KeptnContext, nil
		}
	}

	if len(triggeredSequenceExecutions) > 1 {
		log.Debugf("Sequence with KeptnContext %s is blocked due to triggered sequences in stage %s with KeptnContext %s", queueItem.Scope.KeptnContext, queueItem.Scope.Stage, triggeredSequenceExecutions[0].Scope.KeptnContext)
		return true, triggeredSequenceExecutions[0].Scope.KeptnContext, nil
	}

	return false, "", nil
}

func (sd *SequenceDispatcher) dispatchSequence(queueItem models.QueueItem) error {
	// first, check if the sequence is currently paused
	sequenceExecution, err := sd.sequenceExecutionRepo.GetByTriggeredID(queueItem.Scope.Project, queueItem.EventID)
	if err != nil {
		return err
	}

	if sequenceExecution == nil {
		return common.ErrSequenceNotFound
	}

	if sequenceExecution.IsPaused() || sd.sequenceExecutionRepo.IsContextPaused(queueItem.Scope) {
		//log.Infof("Sequence %s is paused and will not be triggered yet.", queueItem.Scope.KeptnContext)
		return fmt.Errorf("sequence is paused: %w", common.ErrSequenceBlocked)
	}

	sequenceBlocked, blockingSequenceContext, err := sd.isSequenceBlocked(queueItem)
	if err != nil {
		return err
	}

	if sequenceBlocked {
		return fmt.Errorf("blocked by context: %s: %w", blockingSequenceContext, common.ErrSequenceBlockedWaiting)
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
