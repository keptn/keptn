package handler

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

type DispatcherEvent struct {
	event     cloudevents.Event
	timestamp time.Time
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventdispatcher.go . IEventDispatcher
type IEventDispatcher interface {
	Add(event DispatcherEvent) error
	Run()
}

type EventDispatcher struct {
	eventRepo      db.EventRepo
	eventQueueRepo db.EventQueueRepo
	eventSender    keptncommon.EventSender
	syncTimer      *time.Ticker
	logger         keptncommon.LoggerInterface
}

func NewEventDispatcher() EventDispatcher {
	return EventDispatcher{}
}

func (e *EventDispatcher) Add(event DispatcherEvent) error {
	if time.Now().After(event.timestamp) {
		// send event immediately
		return e.eventSender.SendEvent(event.event)
	}
	return e.eventQueueRepo.QueueEvent(models.QueueItem{
		EventID:   event.event.ID(),
		Timestamp: event.timestamp,
	})
}

func (e *EventDispatcher) Run() {
	// TODO make sync interval configurable
	syncInterval := 10
	e.syncTimer = time.NewTicker(time.Duration(syncInterval) * time.Second)
	go func() {
		for {
			<-e.syncTimer.C
			e.logger.Info(fmt.Sprintf("%d seconds have passed. Synchronizing services", syncInterval))
			e.dispatchEvents()
		}
	}()
	e.dispatchEvents()
}

func (e *EventDispatcher) dispatchEvents() {

	events, err := e.eventQueueRepo.GetQueuedEvents(time.Now())
	if err != nil {
		e.logger.Error(fmt.Sprintf("could not fetch event queue: %s", err.Error()))
	}

	for _, queueItem := range events {
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
	}

}
