package handler

import (
	"context"
	"errors"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"sync"
	"time"

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
	Run(ctx context.Context, startSequenceFunc func(event models.Event) error)
	Remove(eventScope models.EventScope) error
	SetStartSequenceCallback(startSequenceFunc func(event models.Event) error) // TODO: not pretty, but for simplicity let's do it this way for now
	Stop()
}

type SequenceDispatcher struct {
	eventRepo             db.EventRepo
	eventQueueRepo        db.EventQueueRepo
	sequenceQueue         db.SequenceQueueRepo
	sequenceExecutionRepo db.SequenceExecutionRepo
	theClock              clock.Clock
	syncInterval          time.Duration
	startSequenceFunc     func(event models.Event) error
	shipyardController    shipyardController
	mutex                 sync.Mutex
	ticker                *clock.Ticker
}

// NewSequenceDispatcher creates a new SequenceDispatcher
func NewSequenceDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	syncInterval time.Duration,
	theClock clock.Clock,
	sequenceExecutionRepo db.SequenceExecutionRepo,

) ISequenceDispatcher {
	return &SequenceDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		sequenceQueue:         sequenceQueueRepo,
		theClock:              theClock,
		syncInterval:          syncInterval,
		mutex:                 sync.Mutex{},
		sequenceExecutionRepo: sequenceExecutionRepo,
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
		} else if err == ErrSequenceBlockedWaiting {
			// if the sequence is currently blocked and should wait, insert it into the queue
			if err2 := sd.sequenceQueue.QueueSequence(queueItem); err2 != nil {
				return err2
			}
			return ErrSequenceBlockedWaiting
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

func (sd *SequenceDispatcher) SetStartSequenceCallback(startSequenceFunc func(event models.Event) error) {
	sd.startSequenceFunc = startSequenceFunc
}

func (sd *SequenceDispatcher) Run(ctx context.Context, startSequenceFunc func(event models.Event) error) {
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
	if sd.ticker == nil {
		return
	}
	sd.ticker.Stop()
}

func (sd *SequenceDispatcher) dispatchSequences() {
	queuedSequences, err := sd.sequenceQueue.GetQueuedSequences()
	if err != nil {
		if err == db.ErrNoEventFound {
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

func (sd *SequenceDispatcher) dispatchSequence(queuedSequence models.QueueItem) error {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	queuedSequence, err := sd.sequenceExecutionRepo.Get(
		models.SequenceExecutionFilter{
			Scope: models.EventScope{
				EventData: keptnv2.EventData{
					Project: queuedSequence.Scope.Project,
					Stage:   queuedSequence.Scope.Stage,
				},
				TriggeredID: queuedSequence.EventID,
			},
		})

	if err != nil {
		return err
	}
	//if queuedSequence.
	// first, check if the sequence is currently paused
	if sd.eventQueueRepo.IsSequenceOfEventPaused(queuedSequence.Scope) {
		log.Infof("Sequence %s is currently paused. Will not start it yet.", queuedSequence.Scope.KeptnContext)
		return ErrSequenceBlocked
	}

	startedSequenceExecutions, err := sd.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: queuedSequence.Scope.Project,
				Stage:   queuedSequence.Scope.Stage,
			},
		},
		Status: []string{models.SequenceStartedState},
	})

	if err != nil {
		return err
	}

	if len(startedSequenceExecutions) > 0 {
		log.Infof("Sequence %s cannot be started yet because sequences are still running in stage %s", queuedSequence.Scope.KeptnContext, queuedSequence.Scope.Stage)
		return ErrSequenceBlockedWaiting
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
