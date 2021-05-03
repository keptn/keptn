package handler

import (
	"context"
	"github.com/benbjohnson/clock"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventdispatcher.go . IEventDispatcher
// IEventDispatcher is responsible for dispatching events to be sent to the event broker
type IEventDispatcher interface {
	Add(event models.DispatcherEvent) error
	Run(ctx context.Context)
}

// EventDispatcher is an implementation of IEventDispatcher
// It regularly fetches (queued) events from the database and eventually
// forwards them to the event broker
type EventDispatcher struct {
	eventRepo      db.EventRepo
	eventQueueRepo db.EventQueueRepo
	eventSender    keptncommon.EventSender
	theClock       clock.Clock
	syncInterval   time.Duration
}

// NewEventDispatcher creates a new EventDispatcher
func NewEventDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	eventSender keptncommon.EventSender,
	syncInterval time.Duration,

) IEventDispatcher {
	return &EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		theClock:       clock.New(),
		syncInterval:   syncInterval,
	}
}

// Add adds a DispatcherEvent to the event queue
func (e *EventDispatcher) Add(event models.DispatcherEvent) error {

	if e.theClock.Now().UTC().After(event.TimeStamp) {
		// send event immediately
		return e.eventSender.SendEvent(event.Event)
	}

	ed, err := models.ConvertToEvent(event.Event)
	if err != nil {
		return err
	}
	eventScope, err := models.NewEventScope(*ed)
	if err != nil {
		return err
	}

	return e.eventQueueRepo.QueueEvent(models.QueueItem{
		Scope:     *eventScope,
		EventID:   event.Event.ID(),
		Timestamp: event.TimeStamp,
	})
}

// Run starts the event dispatcher loop which will periodically fetch (queued) events
// from the database and eventually forward/send them to the event broker
// The fetch interval is configured when creating a EventDispatcher using the "syncInterval" field
func (e *EventDispatcher) Run(ctx context.Context) {
	ticker := e.theClock.Ticker(e.syncInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("cancelling event dispatcher loop")
				return
			case <-ticker.C:
				log.Debugf("%.2f seconds have passed. Dispatching events", e.syncInterval.Seconds())
				e.dispatchEvents()
			}
		}
	}()
}

func (e *EventDispatcher) dispatchEvents() {

	events, err := e.eventQueueRepo.GetQueuedEvents(e.theClock.Now().UTC())
	if err != nil {
		log.Debugf("could not fetch event queue: %s", err.Error())
	}

	for _, queueItem := range events {
		log.Infof("Dispatching event with ID %s", queueItem.EventID)
		events, err := e.eventRepo.GetEvents(queueItem.Scope.Project, common.EventFilter{ID: &queueItem.EventID}, common.TriggeredEvent)
		if err != nil {
			log.Errorf("could not fetch event with ID %s: %s", queueItem.EventID, err.Error())
			continue
		}

		if len(events) == 0 {
			log.Debugf("could not find event with ID %s in project %s", queueItem.EventID, queueItem.Scope.Project)
			continue
		}
		triggeredEvent := events[0]

		ce := &cloudevents.Event{}
		if err := keptnv2.Decode(triggeredEvent, ce); err != nil {
			log.Errorf("could not convert triggeredEvent to CloudEvent: %s", err.Error())
			continue
		}

		if err := e.eventSender.SendEvent(*ce); err != nil {
			log.Errorf("could not send CloudEvent: %s", err.Error())
			continue
		}

		if err := e.eventQueueRepo.DeleteQueuedEvent(queueItem.EventID); err != nil {
			log.Errorf("could not delete event from event queue: %s", err.Error())
			continue
		}
	}

}
