package handler

import (
	"context"
	"github.com/benbjohnson/clock"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
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
	eventRepo      db.EventRepo
	eventQueueRepo db.EventQueueRepo
	theClock       clock.Clock
	syncInterval   time.Duration
}

// NewSequenceDispatcher creates a new SequenceDispatcher
func NewSequenceDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	eventSender keptncommon.EventSender,
	syncInterval time.Duration,

) ISequenceDispatcher {
	return &SequenceDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		theClock:       clock.New(),
		syncInterval:   syncInterval,
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

}
