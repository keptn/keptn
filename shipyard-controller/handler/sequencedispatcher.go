package handler

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/sequencedispatcher.go . ISequenceDispatcher
// IEventDispatcher is responsible for dispatching events to be sent to the event broker
type ISequenceDispatcher interface {
	Add(event models.DispatcherEvent) error
	Run(ctx context.Context)
}

type SequenceDispatcher struct {
	eventRepo         db.EventRepo
	sequenceQueue     db.SequenceQueueRepo
	sequenceStateRepo db.SequenceStateRepo
	theClock          clock.Clock
	syncInterval      time.Duration
}

// NewSequenceDispatcher creates a new SequenceDispatcher
func NewSequenceDispatcher(
	eventRepo db.EventRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	sequenceStateRepo db.SequenceStateRepo,
	syncInterval time.Duration,

) ISequenceDispatcher {
	return &SequenceDispatcher{
		eventRepo:         eventRepo,
		sequenceQueue:     sequenceQueueRepo,
		sequenceStateRepo: sequenceStateRepo,
		theClock:          clock.New(),
		syncInterval:      syncInterval,
	}
}

func (sd *SequenceDispatcher) Add(event models.DispatcherEvent) error {
	panic("implement me")
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
			log.WithError(err).Error("could not dispatch sequence with keptnContext %s", queuedSequence.EventID)
		}
	}
}

func (sd *SequenceDispatcher) dispatchSequence(queuedSequence models.QueueItem) error {
	// fetch all sequences that are currently running in the project where the sequence should run
	_, err := sd.sequenceStateRepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project: queuedSequence.Scope.Project,
			State:   models.SequenceTriggeredState,
		},
	})
	if err != nil {
		return err
	}

	// check if there is a free slot in the stage
	return nil
}
