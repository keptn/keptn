package handler

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencedispatcher.go . ISequenceDispatcher
// ISequenceDispatcher is responsible for dispatching events to be sent to the event broker
type ISequenceDispatcher interface {
	Add(queueItem models.QueueItem) error
	Run(ctx context.Context)
}

type SequenceDispatcher struct {
	eventRepo     db.EventRepo
	sequenceQueue db.SequenceQueueRepo
	sequenceRepo  db.TaskSequenceRepo
	theClock      clock.Clock
	syncInterval  time.Duration
	eventChannel  chan models.Event
}

// NewSequenceDispatcher creates a new SequenceDispatcher
func NewSequenceDispatcher(
	eventRepo db.EventRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	sequenceStateRepo db.TaskSequenceRepo,
	syncInterval time.Duration,
	eventChannel chan models.Event,

) ISequenceDispatcher {
	return &SequenceDispatcher{
		eventRepo:     eventRepo,
		sequenceQueue: sequenceQueueRepo,
		sequenceRepo:  sequenceStateRepo,
		theClock:      clock.New(),
		syncInterval:  syncInterval,
		eventChannel:  eventChannel,
	}
}

func (sd *SequenceDispatcher) Add(queueItem models.QueueItem) error {
	if err := sd.sequenceQueue.QueueSequence(queueItem); err != nil {
		return err
	}
	return sd.dispatchSequence(queueItem)
}

func (sd *SequenceDispatcher) Run(ctx context.Context) {
	ticker := sd.theClock.Ticker(sd.syncInterval)
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
	// fetch all sequences that are currently running in the stage of the project where the sequence should run
	runningSequencesInStage, err := sd.sequenceRepo.GetTaskSequences(queuedSequence.Scope.Project, models.TaskSequenceEvent{
		Stage: queuedSequence.Scope.Stage,
	})
	if err != nil {
		return err
	}
	/// if there is a sequence running in the stage, we cannot trigger this sequence yet
	if runningSequencesInStage != nil && len(runningSequencesInStage) > 0 {
		log.Infof("sequence %s cannot be started yet because sequences are still running in stage %s", queuedSequence.Scope.KeptnContext, queuedSequence.Scope.Stage)
		return nil
	}

	events, err := sd.eventRepo.GetEvents(queuedSequence.Scope.Project, common.EventFilter{
		ID: &queuedSequence.EventID,
	}, common.TriggeredEvent)

	if err != nil {
		return err
	}

	if events == nil || len(events) == 0 {
		return fmt.Errorf("sequence.triggered event with ID %s cannot be found anymore", queuedSequence.EventID)
	}

	sequenceTriggeredEvent := events[0]

	sd.eventChannel <- sequenceTriggeredEvent

	if err := sd.sequenceQueue.DeleteQueuedSequences(queuedSequence); err != nil {
		return err
	}

	return nil
}
