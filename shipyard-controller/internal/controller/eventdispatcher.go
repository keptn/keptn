package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"strings"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/benbjohnson/clock"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

// IEventDispatcher is responsible for dispatching events to be sent to the event broker
//
//go:generate moq -pkg fake -skip-ensure -out ./fake/eventdispatcher.go . IEventDispatcher
type IEventDispatcher interface {
	Add(event models.DispatcherEvent, skipQueue bool) error
	Run(ctx context.Context)
	Stop()
}

// EventDispatcher is an implementation of IEventDispatcher
// It regularly fetches (queued) events from the database and eventually
// forwards them to the event broker
type EventDispatcher struct {
	eventRepo             db.EventRepo
	eventQueueRepo        db.EventQueueRepo
	sequenceExecutionRepo db.SequenceExecutionRepo
	eventSender           keptncommon.EventSender
	theClock              clock.Clock
	syncInterval          time.Duration
	ticker                *clock.Ticker
}

// NewEventDispatcher creates a new EventDispatcher
func NewEventDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	sequenceExecutionRepo db.SequenceExecutionRepo,
	eventSender keptncommon.EventSender,
	syncInterval time.Duration,
) *EventDispatcher {
	return &EventDispatcher{
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		sequenceExecutionRepo: sequenceExecutionRepo,
		eventSender:           eventSender,
		theClock:              clock.New(),
		syncInterval:          syncInterval,
	}
}

// Add adds a DispatcherEvent to the event queue
func (e *EventDispatcher) Add(event models.DispatcherEvent, skipQueue bool) error {
	ed, err := models.ConvertToEvent(event.Event)
	if err != nil {
		return err
	}
	eventScope, err := models.NewEventScope(*ed)
	if err != nil {
		return err
	}

	if skipQueue {
		if err := e.eventSender.Send(context.TODO(), event.Event); err != nil {
			return err
		}
		log.
			WithFields(log.Fields{
				"source":       eventScope.EventSource,
				"keptncontext": eventScope.KeptnContext,
				"project":      eventScope.Project,
				"service":      eventScope.Service,
				"stage":        eventScope.Stage,
			}).
			Infof("[DISPATCHED] Event '%s'", eventScope.EventType)
		return nil
	}
	if e.theClock.Now().UTC().Equal(event.TimeStamp) || e.theClock.Now().UTC().After(event.TimeStamp) {
		// try to send event immediately
		if err := e.tryToSendEvent(*eventScope, event); err != nil {
			// if the event cannot be sent because it is blocked by other sequences,
			// we'll add it to the queue and try to send it again later
			if !strings.Contains(err.Error(), common.OtherActiveSequencesRunning) && err != common.ErrSequencePaused {
				// in all other cases, return the error
				return err
			}
		} else {
			log.
				WithFields(log.Fields{
					"source":       eventScope.EventSource,
					"keptncontext": eventScope.KeptnContext,
					"project":      eventScope.Project,
					"service":      eventScope.Service,
					"stage":        eventScope.Stage,
				}).
				Infof("[DISPATCHED] Event '%s'", eventScope.EventType)
			return nil
		}
	}

	return e.eventQueueRepo.QueueEvent(models.QueueItem{
		Scope:     *eventScope,
		EventID:   event.Event.ID(),
		Timestamp: event.TimeStamp,
	})
}

func (e *EventDispatcher) OnSequenceAborted(eventScope models.EventScope) {
	e.cleanupQueueOfSequence(eventScope)
}

func (e *EventDispatcher) OnSequenceFinished(event apimodels.KeptnContextExtendedCE) {
	e.cleanupQueueOfSequence(models.EventScope{KeptnContext: event.Shkeptncontext})
}

func (e *EventDispatcher) OnSequenceTimeout(event apimodels.KeptnContextExtendedCE) {
	e.cleanupQueueOfSequence(models.EventScope{KeptnContext: event.Shkeptncontext})
}

// Run starts the event dispatcher loop which will periodically fetch (queued) events
// from the database and eventually forward/send them to the event broker
// The fetch interval is configured when creating a EventDispatcher using the "syncInterval" field
func (e *EventDispatcher) Run(ctx context.Context) {
	e.ticker = e.theClock.Ticker(e.syncInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("cancelling event dispatcher loop")
				return
			case <-e.ticker.C:
				log.Debugf("%.2f seconds have passed. Dispatching events", e.syncInterval.Seconds())
				e.dispatchEvents()
			}
		}
	}()
}

func (e *EventDispatcher) Stop() {
	if e.ticker == nil {
		return
	}
	e.ticker.Stop()
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

		// set the time of the cloud event to the current time since it is being sent now
		ce.SetTime(e.theClock.Now().UTC())

		eventScope, err := models.NewEventScope(triggeredEvent)
		if err != nil {
			log.WithError(err).Error("could not determine scope of event")
			continue
		}

		if err := e.tryToSendEvent(*eventScope, models.DispatcherEvent{Event: *ce, TimeStamp: time.Now().UTC()}); err != nil {
			log.Errorf("could not dispatch event with type '%s' and context '%s': %s", eventScope.EventType, eventScope.KeptnContext, err.Error())
			continue
		} else {
			log.
				WithFields(log.Fields{
					"source":       eventScope.EventSource,
					"keptncontext": eventScope.KeptnContext,
					"project":      eventScope.Project,
					"service":      eventScope.Service,
					"stage":        eventScope.Stage,
				}).
				Infof("[DISPATCHED] Event '%s'", eventScope.EventType)
		}

		if err := e.eventQueueRepo.DeleteQueuedEvent(queueItem.EventID); err != nil {
			log.Errorf("could not delete event with id %s from event queue: %s", queueItem.EventID, err.Error())
			continue
		}
	}
}

func (e *EventDispatcher) tryToSendEvent(eventScope models.EventScope, event models.DispatcherEvent) error {
	if e.sequenceExecutionRepo.IsContextPaused(eventScope) {
		log.Debugf("sequence %s is currently paused. will not send event %s", eventScope.KeptnContext, event.Event.ID())
		return common.ErrSequencePaused
	}
	sequenceExecutions, err := e.sequenceExecutionRepo.Get(
		models.SequenceExecutionFilter{
			Scope:              eventScope,
			CurrentTriggeredID: event.Event.ID(),
		},
	)
	if err != nil {
		return err
	}
	if len(sequenceExecutions) == 0 {
		return common.ErrSequenceNotFound
	}

	filter := models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   eventScope.Stage,
				Service: eventScope.Service,
			},
		},
		Status: []string{apimodels.SequenceStartedState},
	}

	startedSequenceExecutions, err := e.sequenceExecutionRepo.Get(filter)
	if err != nil {
		return err
	}

	err2 := checkStarted(startedSequenceExecutions, event, e)
	if err2 != nil {
		return err2
	}

	return e.eventSender.Send(context.TODO(), event.Event)
}

func checkStarted(startedSequenceExecutions []models.SequenceExecution, event models.DispatcherEvent, e *EventDispatcher) error {
	if startedSequenceExecutions != nil && len(startedSequenceExecutions) > 0 {
		// if there is another sequence with the state 'started'
		for _, otherSequence := range startedSequenceExecutions {
			if otherSequence.Status.CurrentTask.TriggeredID != event.Event.ID() {
				if !e.isCurrentEventOverrulingOtherEvent(otherSequence, event) {
					return errors.New(fmt.Sprint(common.OtherActiveSequencesRunning, otherSequence.Scope.KeptnContext))
				}
			}
		}
	}
	return nil
}

func (e *EventDispatcher) isCurrentEventOverrulingOtherEvent(otherSequence models.SequenceExecution, queuedEvent models.DispatcherEvent) bool {
	otherQueuedEvents, err := e.eventQueueRepo.GetQueuedEvents(e.theClock.Now().UTC())
	if err != nil {
		log.Debugf("could not fetch event queue: %s", err.Error())
		return false
	} else if len(otherQueuedEvents) == 0 {
		return false
	}
	for _, otherEvent := range otherQueuedEvents {
		if otherEvent.EventID == otherSequence.Status.CurrentTask.TriggeredID && otherEvent.Timestamp.Before(queuedEvent.TimeStamp) {
			return true
		}
	}
	return false
}

func (e *EventDispatcher) cleanupQueueOfSequence(eventScope models.EventScope) {
	err := e.eventQueueRepo.DeleteEventQueueStates(models.EventQueueSequenceState{Scope: models.EventScope{
		KeptnContext: eventScope.KeptnContext,
	}})
	if err != nil {
		log.WithError(err).Errorf("could not clear event queue states for context %s", eventScope.KeptnContext)
	}
	err = e.eventQueueRepo.DeleteQueuedEvents(models.EventScope{KeptnContext: eventScope.KeptnContext})
	if err != nil {
		log.WithError(err).Errorf("could not clear event queue for context %s", eventScope.KeptnContext)
	}
}
