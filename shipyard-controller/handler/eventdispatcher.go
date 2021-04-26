package handler

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventdispatcher.go . IEventDispatcher
type IEventDispatcher interface {
	Add(event models.DispatcherEvent) error
	Run(ctx context.Context)
}

type EventDispatcher struct {
	eventRepo      db.EventRepo
	eventQueueRepo db.EventQueueRepo
	eventSender    keptncommon.EventSender
	logger         keptncommon.LoggerInterface
	theClock       clock.Clock
	syncInterval   time.Duration
}

func NewEventDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	eventSender keptncommon.EventSender,
	syncInterval time.Duration,
	logger keptncommon.LoggerInterface,
) IEventDispatcher {
	return &EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		eventSender:    eventSender,
		logger:         logger,
		theClock:       clock.New(),
		syncInterval:   syncInterval,
	}
}

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

func (e *EventDispatcher) Run(ctx context.Context) {
	ticker := e.theClock.Ticker(e.syncInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				e.logger.Info("cancelling event dispatcher loop")
				return
			case <-ticker.C:
				e.logger.Debug(fmt.Sprintf("%.2f seconds have passed. Dispatching events", e.syncInterval.Seconds()))
				e.dispatchEvents()
			}
		}
	}()
}

func (e *EventDispatcher) dispatchEvents() {

	events, err := e.eventQueueRepo.GetQueuedEvents(e.theClock.Now().UTC())
	if err != nil {
		e.logger.Debug(fmt.Sprintf("could not fetch event queue: %s", err.Error()))
	}

	for _, queueItem := range events {
		e.logger.Info("Dispatching event with ID " + queueItem.EventID)
		events, err := e.eventRepo.GetEvents(queueItem.Scope.Project, common.EventFilter{ID: &queueItem.EventID}, common.TriggeredEvent)
		if err != nil {
			e.logger.Error(fmt.Sprintf("could not fetch event with ID %s: %s", queueItem.EventID, err.Error()))
			continue
		}

		if len(events) == 0 {
			e.logger.Info(fmt.Sprintf("could not find event with ID %s in project %s", queueItem.EventID, queueItem.Scope.Project))
			continue
		}
		triggeredEvent := events[0]

		ce := &cloudevents.Event{}
		if err := keptnv2.Decode(triggeredEvent, ce); err != nil {
			e.logger.Error(fmt.Sprintf("could not convert triggeredEvent to CloudEvent: %s", err.Error()))
			continue
		}

		if err := e.eventSender.SendEvent(*ce); err != nil {
			e.logger.Error(fmt.Sprintf("could not send CloudEvent: %s", err.Error()))
			continue
		}

		if err := e.eventQueueRepo.DeleteQueuedEvent(queueItem.EventID); err != nil {
			e.logger.Error(fmt.Sprintf("could not delete event from event queue: %s", err.Error()))
			continue
		}

	}

}
