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
	Add(event models.DispatcherEvent, skipQueue bool) error
	Run(ctx context.Context)
}

// EventDispatcher is an implementation of IEventDispatcher
// It regularly fetches (queued) events from the database and eventually
// forwards them to the event broker
type EventDispatcher struct {
	eventRepo      db.EventRepo
	eventQueueRepo db.EventQueueRepo
	sequenceRepo   db.TaskSequenceRepo
	eventSender    keptncommon.EventSender
	theClock       clock.Clock
	syncInterval   time.Duration
}

// NewEventDispatcher creates a new EventDispatcher
func NewEventDispatcher(
	eventRepo db.EventRepo,
	eventQueueRepo db.EventQueueRepo,
	sequenceRepo db.TaskSequenceRepo,
	eventSender keptncommon.EventSender,
	syncInterval time.Duration,
) *EventDispatcher {
	return &EventDispatcher{
		eventRepo:      eventRepo,
		eventQueueRepo: eventQueueRepo,
		sequenceRepo:   sequenceRepo,
		eventSender:    eventSender,
		theClock:       clock.New(),
		syncInterval:   syncInterval,
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
		return e.eventSender.SendEvent(event.Event)
	}
	if e.theClock.Now().UTC().Equal(event.TimeStamp) || e.theClock.Now().UTC().After(event.TimeStamp) {
		// try to send event immediately
		if err := e.tryToSendEvent(*eventScope, event); err != nil {
			// if the event cannot be sent because it is blocked by other sequences,
			// we'll add it to the queue and try to send it again later
			if err != ErrOtherActiveSequencesRunning && err != ErrSequencePaused {
				// in all other cases, return the error
				return err
			}
		} else {
			return nil
		}
	}

	return e.eventQueueRepo.QueueEvent(models.QueueItem{
		Scope:     *eventScope,
		EventID:   event.Event.ID(),
		Timestamp: event.TimeStamp,
	})
}

func (e *EventDispatcher) OnSequenceFinished(event models.Event) {
	e.cleanupQueueOfSequence(event)
}

func (e *EventDispatcher) OnSequenceTimeout(event models.Event) {
	e.cleanupQueueOfSequence(event)
}

func (e *EventDispatcher) OnSequencePaused(pause models.EventScope) {
	err := e.eventQueueRepo.CreateOrUpdateEventQueueState(models.EventQueueSequenceState{
		State: models.SequencePaused,
		Scope: pause,
	})
	if err != nil {
		log.WithError(err).Error("could not set sequence state to 'paused'")
	}
}

func (e *EventDispatcher) OnSequenceResumed(resume models.EventScope) {
	err := e.eventQueueRepo.CreateOrUpdateEventQueueState(models.EventQueueSequenceState{
		State: models.SequenceStartedState,
		Scope: resume,
	})
	if err != nil {
		log.WithError(err).Error("could not set sequence state to 'started'")
	}
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

		// set the time of the cloud event to the current time since it is being sent now
		ce.SetTime(e.theClock.Now().UTC())

		eventScope, err := models.NewEventScope(triggeredEvent)
		if err != nil {
			log.WithError(err).Error("could not determine scope of event")
			continue
		}

		if err := e.tryToSendEvent(*eventScope, models.DispatcherEvent{Event: *ce, TimeStamp: time.Now().UTC()}); err != nil {
			log.Errorf("could not send CloudEvent: %s", err.Error())
			continue
		}

		if err := e.eventQueueRepo.DeleteQueuedEvent(queueItem.EventID); err != nil {
			log.Errorf("could not delete event from event queue: %s", err.Error())
			continue
		}
	}
}

func (e *EventDispatcher) tryToSendEvent(eventScope models.EventScope, event models.DispatcherEvent) error {
	if e.eventQueueRepo.IsSequenceOfEventPaused(eventScope) {
		log.Infof("sequence %s is currently paused. will not send event %s", eventScope.KeptnContext, event.Event.ID())
		return ErrSequencePaused
	}
	taskExecutions, err := e.sequenceRepo.GetTaskExecutions(eventScope.Project, models.TaskExecution{
		Stage:   eventScope.Stage,
		Service: eventScope.Service,
	})
	if err != nil {
		return err
	}
	if e.isCurrentEventBlockedByOtherTasks(eventScope, taskExecutions, event) {
		return ErrOtherActiveSequencesRunning
	}

	return e.eventSender.SendEvent(event.Event)
}

func (e *EventDispatcher) isCurrentEventBlockedByOtherTasks(eventScope models.EventScope, runningSequencesInStage []models.TaskExecution, queuedEvent models.DispatcherEvent) bool {
	runningSequencesInStage = removeSequencesOfSameContext(eventScope.KeptnContext, runningSequencesInStage)
	// if there is another sequence running in the stage, we cannot send the event
	tasksGroupedByContext := groupSequenceMappingsByContext(runningSequencesInStage)

	for otherKeptnContext, tasksOfContext := range tasksGroupedByContext {
		// if the other sequence is currently paused, it should not block the current event
		if e.eventQueueRepo.IsSequenceOfEventPaused(models.EventScope{
			KeptnContext: otherKeptnContext,
			EventData: keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   eventScope.Stage,
			},
		}) {
			continue
		}
		lastTaskOfSequence := getLastTaskOfSequence(tasksOfContext)
		if lastTaskOfSequence.Task.Name != keptnv2.ApprovalTaskName {
			if e.isCurrentEventOverrulingOtherEvent(lastTaskOfSequence, queuedEvent) {
				continue
			}
			log.Infof("event %s cannot be sent because there are other sequences running in stage %s for service %s - blocked by event %s", eventScope.KeptnContext, eventScope.Stage, eventScope.Service, lastTaskOfSequence.TriggeredEventID)
			return true
		}
	}
	return false
}

func (e *EventDispatcher) isCurrentEventOverrulingOtherEvent(lastTaskOfSequence models.TaskExecution, queuedEvent models.DispatcherEvent) bool {
	otherQueuedEvents, err := e.eventQueueRepo.GetQueuedEvents(e.theClock.Now().UTC())
	if err != nil {
		log.Debugf("could not fetch event queue: %s", err.Error())
		return false
	} else if len(otherQueuedEvents) == 0 {
		return false
	}
	for _, otherEvent := range otherQueuedEvents {
		if otherEvent.EventID == lastTaskOfSequence.TriggeredEventID && otherEvent.Timestamp.Before(queuedEvent.TimeStamp) {
			return true
		}
	}
	return false
}

func (e *EventDispatcher) cleanupQueueOfSequence(event models.Event) {
	err := e.eventQueueRepo.DeleteEventQueueStates(models.EventQueueSequenceState{Scope: models.EventScope{
		KeptnContext: event.Shkeptncontext,
	}})
	if err != nil {
		log.WithError(err).Errorf("could not clear event queue states for context %s", event.Shkeptncontext)
	}
	err = e.eventQueueRepo.DeleteQueuedEvents(models.EventScope{KeptnContext: event.Shkeptncontext})
	if err != nil {
		log.WithError(err).Errorf("could not clear event queue for context %s", event.Shkeptncontext)
	}
}

func removeSequencesOfSameContext(keptnContext string, sequenceTasks []models.TaskExecution) []models.TaskExecution {
	result := []models.TaskExecution{}
	for index := range sequenceTasks {
		if sequenceTasks[index].KeptnContext != keptnContext {
			result = append(result, sequenceTasks[index])
		}
	}
	return result
}
