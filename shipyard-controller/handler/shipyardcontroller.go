package handler

import (
	"context"
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	modelsv2 "github.com/keptn/keptn/shipyard-controller/models/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

const maxRepoReadRetries = 5

var shipyardControllerInstance *shipyardController

//go:generate moq -pkg fake -skip-ensure -out ./fake/shipyardcontroller.go . IShipyardController
type IShipyardController interface {
	GetAllTriggeredEvents(filter common.EventFilter) ([]models.Event, error)
	GetTriggeredEventsOfProject(project string, filter common.EventFilter) ([]models.Event, error)
	HandleIncomingEvent(event models.Event, waitForCompletion bool) error
	ControlSequence(controlSequence models.SequenceControl) error
	StartTaskSequence(event models.Event) error
}

type shipyardController struct {
	eventRepo                  db.EventRepo
	taskSequenceRepo           db.TaskSequenceRepo
	taskSequenceV2Repo         db.TaskSequenceV2Repo
	projectMvRepo              db.ProjectMVRepo
	eventDispatcher            IEventDispatcher
	sequenceDispatcher         ISequenceDispatcher
	sequenceTimeoutChan        chan models.SequenceTimeout
	sequenceTriggeredHooks     []sequencehooks.ISequenceTriggeredHook
	sequenceStartedHooks       []sequencehooks.ISequenceStartedHook
	sequenceTaskTriggeredHooks []sequencehooks.ISequenceTaskTriggeredHook
	sequenceTaskStartedHooks   []sequencehooks.ISequenceTaskStartedHook
	sequenceTaskFinishedHooks  []sequencehooks.ISequenceTaskFinishedHook
	subSequenceFinishedHooks   []sequencehooks.ISubSequenceFinishedHook
	sequenceFinishedHooks      []sequencehooks.ISequenceFinishedHook
	sequenceAbortedHooks       []sequencehooks.ISequenceAbortedHook
	sequenceTimoutHooks        []sequencehooks.ISequenceTimeoutHook
	sequencePausedHooks        []sequencehooks.ISequencePausedHook
	sequenceResumedHooks       []sequencehooks.ISequenceResumedHook
	shipyardRetriever          IShipyardRetriever
}

func GetShipyardControllerInstance(
	ctx context.Context,
	eventDispatcher IEventDispatcher,
	sequenceDispatcher ISequenceDispatcher,
	sequenceTimeoutChannel chan models.SequenceTimeout,
	shipyardRetriever IShipyardRetriever,
) *shipyardController {
	if shipyardControllerInstance == nil {
		cbConnectionInstance := db.GetMongoDBConnectionInstance()
		shipyardControllerInstance = &shipyardController{
			eventRepo:        db.NewMongoDBEventsRepo(cbConnectionInstance),
			taskSequenceRepo: db.NewTaskSequenceMongoDBRepo(cbConnectionInstance),
			projectMvRepo: db.NewProjectMVRepo(
				db.NewMongoDBKeyEncodingProjectsRepo(cbConnectionInstance),
				db.NewMongoDBEventsRepo(cbConnectionInstance)),
			taskSequenceV2Repo:  db.NewMongoDBTaskSequenceV2Repo(cbConnectionInstance),
			eventDispatcher:     eventDispatcher,
			sequenceDispatcher:  sequenceDispatcher,
			sequenceTimeoutChan: sequenceTimeoutChannel,
			shipyardRetriever:   shipyardRetriever,
		}
		shipyardControllerInstance.run(ctx)
	}
	return shipyardControllerInstance
}

func (sc *shipyardController) run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case timeoutSequence := <-sc.sequenceTimeoutChan:
				err := sc.timeoutSequence(timeoutSequence)
				if err != nil {
					log.WithError(err).Error("Unable to cancel sequence")
					return
				}
				break
			}
		}
	}()
}

func (sc shipyardController) StartDispatchers(ctx context.Context) {

	sc.eventDispatcher.Run(ctx)
	sc.sequenceDispatcher.Run(ctx, sc.StartTaskSequence)
}

//func (sc shipyardController) StopDispatchers() {
//	sc.eventDispatcher.
//	sc.sequenceDispatcher.Run(context.Background(), sc.StartTaskSequence)
//}

func (sc *shipyardController) ControlSequence(controlSequence models.SequenceControl) error {
	switch controlSequence.State {
	case models.AbortSequence:
		log.Info("Processing ABORT sequence control")
		return sc.cancelSequence(controlSequence)
	case models.PauseSequence:
		log.Info("Processing PAUSE sequence control")
		sc.onSequencePaused(models.EventScope{
			EventData: keptnv2.EventData{
				Project: controlSequence.Project,
				Stage:   controlSequence.Stage,
			},
			KeptnContext: controlSequence.KeptnContext,
		})
	case models.ResumeSequence:
		log.Info("Processing RESUME sequence control")
		sc.onSequenceResumed(models.EventScope{
			EventData: keptnv2.EventData{
				Project: controlSequence.Project,
				Stage:   controlSequence.Stage,
			},
			KeptnContext: controlSequence.KeptnContext,
		})
	}
	return nil
}

func (sc *shipyardController) HandleIncomingEvent(event models.Event, waitForCompletion bool) error {
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(event.Data, eventData)
	if err != nil {
		log.Errorf("Could not parse event data: %v", err)
		return err
	}

	statusType, err := keptnv2.ParseEventKind(*event.Type)
	if err != nil {
		return err
	}
	done := make(chan error)

	log.Infof("Received event of type %s from %s", *event.Type, *event.Source)
	log.Debugf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, ObjToJSON(event))

	switch statusType {
	case string(common.TriggeredEvent):
		go func() {
			err := sc.handleSequenceTriggered(event)
			if err != nil {
				log.WithError(err).Error("Unable to handle sequence '.triggered' event")
			}
			done <- err
		}()
	case string(common.StartedEvent):
		go func() {
			err := sc.handleTaskEvent(event)
			if err != nil {
				log.WithError(err).Error("Unable to handle task '.started' event")
			}
			done <- err
		}()
	case string(common.FinishedEvent):
		go func() {
			err := sc.handleTaskEvent(event)
			if err != nil {
				log.WithError(err).Error("Unable to handle task '.finished' event")
			}
			done <- err
		}()
	default:
		return nil
	}
	if waitForCompletion {
		return <-done
	}
	return nil
}

func (sc *shipyardController) handleSequenceTriggered(event models.Event) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("unable to create event scope: %w", err)
	}

	// only process 'sequence.triggered' events
	if !keptnv2.IsSequenceEventType(eventScope.EventType) {
		return nil
	}

	log.Infof("Checking if sequence '.triggered' event should start a sequence in project %s", eventScope.Project)
	_, taskSequenceName, _, err := keptnv2.ParseSequenceEventType(eventScope.EventType)
	if err != nil {
		return fmt.Errorf("unable to parse seuqnce event of type %s: %w", eventScope.EventType, err)
	}

	// fetching cached shipyard file from project git repo
	shipyard, err := sc.shipyardRetriever.GetShipyard(eventScope.Project)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Shipyard file: %v", err)
		log.Errorf(msg)
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}

	// check if the sequence is available in the given stage
	sequence, err := GetTaskSequenceInStage(eventScope.Stage, taskSequenceName, shipyard)
	if err != nil {
		msg := fmt.Sprintf("Unable to start sequence %s: %v", taskSequenceName, err)
		log.Error(msg)
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}

	if err := sc.eventRepo.InsertEvent(eventScope.Project, eventScope.WrappedEvent, common.TriggeredEvent); err != nil {
		log.Infof("could not store event that triggered task sequence: %s", err.Error())
	}

	sequenceExecution := modelsv2.TaskSequence{
		ID:       uuid.New().String(),
		Sequence: *sequence,
		Status: modelsv2.TaskSequenceStatus{
			State:         models.SequenceTriggeredState,
			PreviousTasks: []modelsv2.TaskExecutionResult{},
		},
		Scope: modelsv2.EventScope{
			KeptnContext: eventScope.KeptnContext,
			Project:      eventScope.Project,
			Stage:        eventScope.Stage,
			Service:      eventScope.Service,
		},
	}

	if err := sc.taskSequenceV2Repo.Upsert(sequenceExecution); err != nil {
		return fmt.Errorf("could not store task sequence execution: %w", err)
	}
	sc.onSequenceTriggered(eventScope.WrappedEvent)
	return sc.sequenceDispatcher.Add(models.QueueItem{
		Scope:     *eventScope,
		EventID:   eventScope.WrappedEvent.ID,
		Timestamp: common.ParseTimestamp(eventScope.WrappedEvent.Time, nil),
	})
}

// TODO this will become obsolete
func (sc *shipyardController) handleTaskStarted(event models.Event) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("unable to handle 'task.started' event: %w", err)
	}

	if !keptnv2.IsTaskEventType(eventScope.EventType) {
		return nil
	}

	wasTriggered, err := sc.wasTaskTriggered(*eventScope)
	if err != nil {
		return fmt.Errorf("unable to handle %s event: %w", eventScope.EventType, err)
	}

	if !wasTriggered {
		log.Infof("The received %s event keptn context %s is not accociated with a task that was previously triggered",
			eventScope.EventType, eventScope.KeptnContext)
		return nil
	}

	sc.onSequenceTaskStarted(eventScope.WrappedEvent)
	return sc.eventRepo.InsertEvent(eventScope.Project, eventScope.WrappedEvent, common.StartedEvent)
}

func (sc *shipyardController) handleTaskEvent(event models.Event) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("unable to handle 'task.finished' event: %w", err)
	}

	if !keptnv2.IsTaskEventType(eventScope.EventType) {
		return nil
	}

	sequenceExecution, err := sc.getOpenTaskExecution(*eventScope)
	if err != nil {
		return fmt.Errorf("unable to handle %s event: %w", eventScope.EventType, err)
	}

	if sequenceExecution == nil {
		log.Infof("The received %s event with keptn context %s is not accociated with a task that was previously triggered",
			eventScope.EventType, eventScope.KeptnContext)
		return nil
	}

	if keptnv2.IsFinishedEventType(*event.Type) {
		// TODO this should become obsolete by also storing the relevant event data in the sequenceExecution
		// for now keep it because it's needed for aggregating data for next task.triggered event
		err = sc.eventRepo.InsertEvent(eventScope.Project, eventScope.WrappedEvent, common.FinishedEvent)
		if err != nil {
			log.Errorf("unable to store %s event: %v ", eventScope.EventType, err.Error())
		}
	} else if keptnv2.IsStartedEventType(*event.Type) {
		sc.onSequenceTaskStarted(eventScope.WrappedEvent)
	}

	return sc.onTaskProgress(event, *sequenceExecution, eventScope)
}

func (sc *shipyardController) onTaskProgress(event models.Event, sequenceExecution modelsv2.TaskSequence, eventScope *models.EventScope) error {
	taskEvent := modelsv2.TaskEvent{
		EventType: *event.Type,
		Source:    *event.Source,
		Result:    string(eventScope.Result),
		Status:    string(eventScope.Status),
		Time:      event.Time,
	}
	if keptnv2.IsFinishedEventType(taskEvent.EventType) {
		// TODO this will need refactoring
		eventData := map[string]interface{}{}
		err := keptnv2.Decode(event.Data, &eventData)
		if err != nil {
			return err
		}
		if taskProperties, ok := eventData[sequenceExecution.Status.CurrentTask.Name]; ok {
			taskEvent.Properties = taskProperties.(map[string]interface{})
		}
	}
	updatedSequenceExecution, err := sc.taskSequenceV2Repo.AppendTaskEvent(sequenceExecution, taskEvent)
	if err != nil {
		return err
	}

	// now check if the number of .started events matches the number of finished events - if yes, that means were done
	// note: this should also work with multiple replicas because the `AppendTaskEvent` updates the list of events and returns the resulting state
	// atomically, so ONLY the thread that appended the last event to reach the completion state of the task will get the state required for further proceeding with the task sequence
	if !updatedSequenceExecution.Status.CurrentTask.IsFinished() {
		return nil
	}

	triggeredEventType, err := keptnv2.ReplaceEventTypeKind(eventScope.EventType, string(common.TriggeredEvent))
	if err != nil {
		return err
	}

	triggeredEvents, err := sc.eventRepo.GetEventsWithRetry(eventScope.Project, common.EventFilter{Type: triggeredEventType, ID: &eventScope.TriggeredID}, common.TriggeredEvent, maxRepoReadRetries)
	if err != nil {
		return fmt.Errorf("unable to retrieve associated task '.triggered' event with ID %s: %w", eventScope.TriggeredID, err)
	}
	if len(triggeredEvents) == 0 {
		return fmt.Errorf("no matching task '.triggered' event found for event %s with triggered ID %s", eventScope.WrappedEvent.ID, eventScope.TriggeredID)
	}

	// the '.triggered' event can now be removed
	err = sc.eventRepo.DeleteEvent(eventScope.Project, triggeredEvents[0].ID, common.TriggeredEvent)
	if err != nil {
		return fmt.Errorf("unable to delete associated task '.triggered' event with ID %s: %w", eventScope.TriggeredID, err)
	}

	finishedEventsData, err := sc.getFinishedEventData(*eventScope)
	if err != nil {
		return fmt.Errorf("unable to gather task '.finished' event data: %w", err)
	}

	sc.onSequenceTaskFinished(eventScope.WrappedEvent)
	return sc.proceedTaskSequence(*eventScope, *updatedSequenceExecution, finishedEventsData)
}

func (sc *shipyardController) wasTaskTriggered(eventScope models.EventScope) (bool, error) {
	sequenceExecutions, err := sc.getOpenTaskExecution(eventScope)
	if err != nil {
		return false, err
	}
	if sequenceExecutions == nil {
		return false, nil
	}
	return true, nil
}

func (sc *shipyardController) cancelSequence(cancel models.SequenceControl) error {
	sequenceExecutions, err := sc.taskSequenceV2Repo.Get(modelsv2.GetTaskSequenceFilter{Scope: modelsv2.EventScope{
		KeptnContext: cancel.KeptnContext,
		Project:      cancel.Project,
		Stage:        cancel.Stage,
	}})

	if err != nil {
		return fmt.Errorf("unable to get active task executions for project %s in stage %s for keptn context %s", cancel.Project, cancel.Stage, cancel.KeptnContext)
	}

	if len(sequenceExecutions) == 0 {
		log.Infof("no active task execution for context %s found.", cancel.KeptnContext)
		return nil
	}
	// TODO check state and remove from dispatcher if sequence is currently in queue, i.e. waiting or triggered

	// delete all open .triggered events for the task sequence
	for _, sequenceExecution := range sequenceExecutions {
		err := sc.eventRepo.DeleteEvent(cancel.Project, sequenceExecution.Status.CurrentTask.TriggeredID, common.TriggeredEvent)
		if err != nil {
			// log the error, but continue
			log.WithError(err).Error("could not delete event")
		}

		sequenceTriggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(models.EventScope{
			EventData: keptnv2.EventData{
				Project: cancel.Project,
				Stage:   sequenceExecution.Scope.Stage,
			},
			KeptnContext: cancel.KeptnContext,
		}, sequenceExecution.Sequence.Name)
		if err != nil {
			return err
		}

		if sequenceTriggeredEvent != nil {
			return sc.forceTaskSequenceCompletion(*sequenceTriggeredEvent, sequenceExecution)
		}
	}
	return nil
}

func (sc *shipyardController) forceTaskSequenceCompletion(sequenceTriggeredEvent models.Event, sequenceExecution modelsv2.TaskSequence) error {
	sc.onSequenceAborted(sequenceTriggeredEvent)
	scope, err := models.NewEventScope(sequenceTriggeredEvent)
	if err != nil {
		return err
	}

	scope.Result = keptnv2.ResultPass
	scope.Status = keptnv2.StatusAborted

	return sc.completeTaskSequence(*scope, sequenceExecution, sequenceTriggeredEvent.ID, models.SequenceFinished)
}

func (sc *shipyardController) timeoutSequence(timeout models.SequenceTimeout) error {
	log.Infof("sequence %s has been timed out", timeout.KeptnContext)
	eventScope, err := models.NewEventScope(timeout.LastEvent)
	if err != nil {
		return err
	}

	eventScope.Status = keptnv2.StatusErrored
	eventScope.Result = keptnv2.ResultFailed
	eventScope.Message = fmt.Sprintf("sequence timed out while waiting for task %s to receive a correlating .started or .finished event", *timeout.LastEvent.Type)

	sequenceExecutions, err := sc.taskSequenceV2Repo.Get(modelsv2.GetTaskSequenceFilter{
		CurrentTriggeredID: timeout.LastEvent.ID,
		Scope: modelsv2.EventScope{
			KeptnContext: eventScope.KeptnContext,
			Project:      eventScope.Project,
			Stage:        eventScope.Stage,
		}})

	if err != nil {
		return fmt.Errorf("Could not retrieve task executions associated to eventID %s: %s", timeout.LastEvent.ID, err.Error())
	}

	if len(sequenceExecutions) == 0 {
		log.Infof("No task executions associated with eventID %s found", timeout.LastEvent.ID)
		return nil
	}

	sequenceExecution := sequenceExecutions[0]
	sc.onSequenceTimeout(timeout.LastEvent)
	taskSequenceTriggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(*eventScope, sequenceExecution.Sequence.Name)
	if err != nil {
		return err
	}
	if taskSequenceTriggeredEvent != nil {
		if err := sc.completeTaskSequence(*eventScope, sequenceExecution, taskSequenceTriggeredEvent.ID, models.TimedOut); err != nil {
			return err
		}
	}
	return nil
}

func (sc *shipyardController) triggerSequenceFailed(eventScope models.EventScope, msg string, taskSequenceName string) error {
	event := eventScope.WrappedEvent
	sc.onSequenceTriggered(event) //TODO: remove?
	finishedEvent := event
	finishedEventData := keptnv2.EventData{
		Project: eventScope.Project,
		Stage:   eventScope.Stage,
		Service: eventScope.Service,
		Labels:  eventScope.Labels,
		Status:  keptnv2.StatusErrored,
		Result:  keptnv2.ResultFailed,
		Message: msg,
	}
	finishedEvent.Data = finishedEventData

	sc.onSequenceFinished(finishedEvent)
	return sc.sendTaskSequenceFinishedEvent(models.EventScope{
		EventData:    finishedEventData,
		KeptnContext: event.Shkeptncontext,
	}, taskSequenceName, event.ID)
}

func (sc *shipyardController) StartTaskSequence(event models.Event) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return err
	}

	_, taskSequenceName, _, err := keptnv2.ParseSequenceEventType(*event.Type)
	if err != nil {
		return err
	}

	sc.onSequenceStarted(event)

	sequenceExecutions, err := sc.taskSequenceV2Repo.Get(
		modelsv2.GetTaskSequenceFilter{
			Scope: modelsv2.EventScope{
				KeptnContext: eventScope.KeptnContext,
				Project:      eventScope.Project,
				Stage:        eventScope.Stage,
			},
			Name:   taskSequenceName,
			Status: []string{models.SequenceTriggeredState},
		},
	)
	if err != nil {
		msg := fmt.Sprintf("could not get sequence execution state %s: %s", taskSequenceName, err.Error())
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}
	if len(sequenceExecutions) == 0 {
		msg := fmt.Sprintf("no sequence execution state found for sequence %s", taskSequenceName)
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}
	sequenceExecution := sequenceExecutions[0]
	sequenceExecution.Status.State = models.SequenceStartedState
	if err := sc.taskSequenceV2Repo.Upsert(sequenceExecution); err != nil {
		msg := fmt.Sprintf("could not update sequence execution state %s: %s", taskSequenceName, err.Error())
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}
	sc.onSequenceStarted(event)
	return sc.proceedTaskSequence(*eventScope, sequenceExecution, []interface{}{})
}

func (sc *shipyardController) getOpenTaskExecution(eventScope models.EventScope) (*modelsv2.TaskSequence, error) {
	sequenceExecutions, err := sc.taskSequenceV2Repo.Get(modelsv2.GetTaskSequenceFilter{
		Scope: modelsv2.EventScope{
			KeptnContext: eventScope.KeptnContext,
			Project:      eventScope.Project,
			Stage:        eventScope.Stage,
		},
		CurrentTriggeredID: eventScope.TriggeredID,
	})
	if err != nil {
		return nil, err
	}
	if len(sequenceExecutions) == 0 {
		return nil, nil
	}
	return &sequenceExecutions[0], nil
}

func (sc *shipyardController) getFinishedEventData(eventScope models.EventScope) ([]interface{}, error) {
	// TODO: make use of event data in sequenceExecution!!
	allFinishedEventsForTask, err := sc.eventRepo.GetEvents(eventScope.Project, common.EventFilter{
		Type:         "",
		Stage:        &eventScope.Stage,
		Service:      &eventScope.Service,
		KeptnContext: &eventScope.KeptnContext,
	}, common.FinishedEvent)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve %s events: %s", eventScope.EventType, err.Error())
	}

	log.Infof("Found %d events. Aggregating their properties for next task ", len(allFinishedEventsForTask))

	finishedEventsData := []interface{}{}

	for index := range allFinishedEventsForTask {
		marshal, _ := json.Marshal(allFinishedEventsForTask[index].Data)
		var tmp interface{}
		_ = json.Unmarshal(marshal, &tmp)
		finishedEventsData = append(finishedEventsData, tmp)
	}
	return finishedEventsData, nil
}

func (sc *shipyardController) GetAllTriggeredEvents(filter common.EventFilter) ([]models.Event, error) {
	projects, err := sc.projectMvRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []models.Event{}
	for _, project := range projects {
		events, err := sc.eventRepo.GetEvents(project.ProjectName, filter, common.TriggeredEvent)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (sc *shipyardController) GetTriggeredEventsOfProject(projectName string, filter common.EventFilter) ([]models.Event, error) {
	project, err := sc.projectMvRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	} else if project == nil {
		return nil, ErrProjectNotFound
	}
	events, err := sc.eventRepo.GetEvents(projectName, filter, common.TriggeredEvent)
	if err != nil && err != db.ErrNoEventFound {
		return nil, err
	} else if err != nil && err == db.ErrNoEventFound {
		return []models.Event{}, nil
	}
	return events, nil
}

func (sc *shipyardController) proceedTaskSequence(eventScope models.EventScope, sequenceExecution modelsv2.TaskSequence, eventHistory []interface{}) error {

	// get the input for the .triggered event that triggered the previous sequence and append it to the list of previous events to gather all required data for the next stage
	inputEvent, eventHistory, err := sc.appendTriggerEventProperties(sequenceExecution, eventHistory)
	if err != nil {
		return err
	}

	task := sequenceExecution.GetNextTaskOfSequence()
	if task == nil {
		if sequenceExecution.Status.CurrentTask.IsFailed() {
			eventScope.Result = keptnv2.ResultFailed
		}
		if sequenceExecution.Status.CurrentTask.IsErrored() {
			eventScope.Status = keptnv2.StatusErrored
		}
		// task sequence completed -> send .finished event and check if a new task sequence should be triggered by the completion
		err = sc.completeTaskSequence(eventScope, sequenceExecution, inputEvent.ID, models.SequenceFinished)
		if err != nil {
			log.Errorf("Could not complete task sequence %s.%s with KeptnContext %s: %s", eventScope.Stage, sequenceExecution.Sequence.Name, eventScope.KeptnContext, err.Error())
			return err
		}

		lastTaskName := ""
		if lastTask := sequenceExecution.GetLastTaskExecutionResult(); lastTask != nil {
			lastTaskName = lastTask.Name
		}
		return sc.triggerNextTaskSequences(eventScope, sequenceExecution.Sequence, eventHistory, inputEvent, lastTaskName)
	}
	//return sc.sendTaskTriggeredEvent(eventScope, sequenceExecution.Sequence.Name, *task, eventHistory)
	return sc.triggerTask(eventScope, sequenceExecution, *task, eventHistory)
}

// this function retrieves the .triggered event for the task sequence and appends its properties to the existing .finished events
// this ensures that all parameters set in the .triggered event are received by all execution plane services, instead of just the first one
func (sc *shipyardController) appendTriggerEventProperties(sequenceExecution modelsv2.TaskSequence, eventHistory []interface{}) (*models.Event, []interface{}, error) {
	// TODO only make use of one EventScope object to avoid these conversions
	triggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(
		models.EventScope{
			KeptnContext: sequenceExecution.Scope.KeptnContext,
			EventData: keptnv2.EventData{
				Stage:   sequenceExecution.Scope.Stage,
				Project: sequenceExecution.Scope.Project,
			},
		}, sequenceExecution.Sequence.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load event that triggered task sequence %s.%s with KeptnContext %s: %w", sequenceExecution.Scope.Stage, sequenceExecution.Sequence.Name, sequenceExecution.Scope.KeptnContext, err)
	}
	if triggeredEvent != nil {
		marshal, err := json.Marshal(triggeredEvent.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to marshal triggered event: %w", err)
		}
		var tmp interface{}
		_ = json.Unmarshal(marshal, &tmp)
		eventHistory = append(eventHistory, tmp)
	}

	return triggeredEvent, eventHistory, nil
}

func (sc *shipyardController) triggerNextTaskSequences(eventScope models.EventScope, completedSequence keptnv2.Sequence, eventHistory []interface{}, inputEvent *models.Event, previousTask string) error {
	shipyard, err := sc.shipyardRetriever.GetCachedShipyard(eventScope.Project)
	if err != nil {
		return err
	}
	nextSequences := GetTaskSequencesByTrigger(eventScope, completedSequence.Name, shipyard, previousTask)

	if len(nextSequences) == 0 {
		sc.onSequenceFinished(*inputEvent)
	}

	for _, sequence := range nextSequences {
		newScope := &models.EventScope{
			EventData: keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   sequence.StageName,
				Service: eventScope.Service,
			},
			KeptnContext: eventScope.KeptnContext,
		}

		err := sc.sendTaskSequenceTriggeredEvent(newScope, sequence.Sequence.Name, inputEvent, eventHistory)
		if err != nil {
			log.Errorf("could not send event %s.%s.triggered: %s",
				newScope.Stage, sequence.Sequence.Name, err.Error())
			continue
		}
	}
	return nil
}

func (sc *shipyardController) completeTaskSequence(eventScope models.EventScope, sequenceExecution modelsv2.TaskSequence, triggeredID string, reason string) error {
	sequenceExecution.Status.State = reason
	err := sc.taskSequenceV2Repo.Upsert(sequenceExecution)

	if err != nil {
		return err
	}

	log.Infof("Deleting all task.finished events of task sequence %s with context %s", sequenceExecution.Sequence.Name, sequenceExecution.Scope.KeptnContext)
	if err := sc.eventRepo.DeleteAllFinishedEvents(eventScope); err != nil {
		return err
	}
	return sc.sendTaskSequenceFinishedEvent(eventScope, sequenceExecution.Sequence.Name, triggeredID)
}

func (sc *shipyardController) sendTaskSequenceTriggeredEvent(eventScope *models.EventScope, taskSequenceName string, inputEvent *models.Event, eventHistory []interface{}) error {
	eventPayload := map[string]interface{}{}

	eventPayload["project"] = eventScope.Project
	eventPayload["stage"] = eventScope.Stage
	eventPayload["service"] = eventScope.Service
	eventPayload["result"] = ""
	eventPayload["status"] = ""

	mergedPayload, err := GetMergedPayloadForSequenceTriggeredEvent(inputEvent, eventPayload, eventHistory)
	if err != nil {
		return err
	}

	eventType := eventScope.Stage + "." + taskSequenceName

	var event cloudevents.Event
	if mergedPayload != nil {
		event = common.CreateEventWithPayload(eventScope.KeptnContext, "", keptnv2.GetTriggeredEventType(eventType), mergedPayload)
	} else {
		event = common.CreateEventWithPayload(eventScope.KeptnContext, "", keptnv2.GetTriggeredEventType(eventType), inputEvent)
	}

	toEvent, err := models.ConvertToEvent(event)
	if err != nil {
		return fmt.Errorf("could not store event that triggered task sequence: " + err.Error())
	}
	if err := sc.eventRepo.InsertEvent(eventScope.Project, *toEvent, common.TriggeredEvent); err != nil {
		return fmt.Errorf("could not store event that triggered task sequence: " + err.Error())
	}

	return sc.eventDispatcher.Add(models.DispatcherEvent{TimeStamp: time.Now().UTC(), Event: event}, true)
}

func (sc *shipyardController) sendTaskSequenceFinishedEvent(eventScope models.EventScope, taskSequenceName, triggeredID string) error {
	eventType := eventScope.Stage + "." + taskSequenceName

	event := common.CreateEventWithPayload(eventScope.KeptnContext, triggeredID, keptnv2.GetFinishedEventType(eventType), eventScope.EventData)

	if toEvent, err := models.ConvertToEvent(event); err == nil {
		sc.onSubSequenceFinished(*toEvent)
	}

	return sc.eventDispatcher.Add(models.DispatcherEvent{TimeStamp: time.Now().UTC(), Event: event}, true)
}

func (sc *shipyardController) triggerTask(eventScope models.EventScope, sequenceExecution modelsv2.TaskSequence, task keptnv2.Task, eventHistory []interface{}) error {
	eventPayload := map[string]interface{}{}

	eventPayload["project"] = sequenceExecution.Scope.Project
	eventPayload["stage"] = sequenceExecution.Scope.Stage
	eventPayload["service"] = sequenceExecution.Scope.Service

	eventPayload[task.Name] = task.Properties

	var mergedPayload interface{}
	mergedPayload = nil
	if eventHistory != nil {
		for index := range eventHistory {
			if mergedPayload == nil {
				mergedPayload = common.Merge(eventPayload, eventHistory[index])
			} else {
				mergedPayload = common.Merge(mergedPayload, eventHistory[index])
			}
		}
	}

	// make sure the result from the previous event is used
	eventPayload["result"] = eventScope.Result
	eventPayload["status"] = eventScope.Status
	// make sure the 'message' property from the previous event is set to ""
	eventPayload["message"] = ""

	event := common.CreateEventWithPayload(eventScope.KeptnContext, "", keptnv2.GetTriggeredEventType(task.Name), eventPayload)

	storeEvent := &models.Event{}
	if err := keptnv2.Decode(event, storeEvent); err != nil {
		log.Errorf("could not transform CloudEvent for storage in mongodb: %s", err.Error())
		return err
	}

	sendTaskTimestamp := time.Now().UTC()
	if task.TriggeredAfter != "" {
		if duration, err := time.ParseDuration(task.TriggeredAfter); err == nil {
			sendTaskTimestamp = sendTaskTimestamp.Add(duration)
		} else {
			log.Errorf("could not parse triggeredAfter property: %s", err.Error())
			// TODO how do we handle this? send event immediately or not at all?
		}
		log.Infof("queueing %s event with ID %s to be sent at %s", event.Type(), event.ID(), sendTaskTimestamp.String())
	}
	storeEvent.Time = timeutils.GetKeptnTimeStamp(sendTaskTimestamp)

	if err := sc.eventRepo.InsertEvent(eventScope.Project, *storeEvent, common.TriggeredEvent); err != nil {
		log.Errorf("Could not store event: %s", err.Error())
		return err
	}

	sc.onSequenceTaskTriggered(*storeEvent)

	sequenceExecution.Status.PreviousTasks = append(
		sequenceExecution.Status.PreviousTasks,
		modelsv2.TaskExecutionResult{
			Name:        sequenceExecution.Status.CurrentTask.Name,
			TriggeredID: sequenceExecution.Status.CurrentTask.TriggeredID,
			Result:      string(eventScope.Result),
			Status:      string(eventScope.Status),
			Properties:  nil,
		},
	)

	sequenceExecution.Status.CurrentTask = modelsv2.TaskExecution{
		Name:        task.Name,
		TriggeredID: storeEvent.ID,
		Events:      []modelsv2.TaskEvent{},
	}

	// special handling for approval events
	if task.Name == "approval" {
		sequenceExecution.Status.State = models.SequenceSuspended
	} else {
		sequenceExecution.Status.State = models.SequenceStartedState
	}

	if err := sc.taskSequenceV2Repo.Upsert(sequenceExecution); err != nil {
		return err
	}
	if err := sc.eventDispatcher.Add(models.DispatcherEvent{TimeStamp: sendTaskTimestamp, Event: event}, false); err != nil {
		return err
	}
	return nil
}

func (sc *shipyardController) sendTaskTriggeredEvent(eventScope models.EventScope, taskSequenceName string, task models.Task, eventHistory []interface{}) error {
	common.LockServiceInStageOfProject(eventScope.Project, eventScope.Stage, eventScope.Service)
	defer common.UnlockServiceInStageOfProject(eventScope.Project, eventScope.Stage, eventScope.Service)
	eventPayload := map[string]interface{}{}

	eventPayload["project"] = eventScope.Project
	eventPayload["stage"] = eventScope.Stage
	eventPayload["service"] = eventScope.Service

	eventPayload[task.Name] = task.Properties

	var mergedPayload interface{}
	mergedPayload = nil
	if eventHistory != nil {
		for index := range eventHistory {
			if mergedPayload == nil {
				mergedPayload = common.Merge(eventPayload, eventHistory[index])
			} else {
				mergedPayload = common.Merge(mergedPayload, eventHistory[index])
			}
		}
	}

	// make sure the result from the previous event is used
	eventPayload["result"] = eventScope.Result
	eventPayload["status"] = eventScope.Status
	// make sure the 'message' property from the previous event is set to ""
	eventPayload["message"] = ""

	event := common.CreateEventWithPayload(eventScope.KeptnContext, "", keptnv2.GetTriggeredEventType(task.Name), eventPayload)

	storeEvent := &models.Event{}
	if err := keptnv2.Decode(event, storeEvent); err != nil {
		log.Errorf("could not transform CloudEvent for storage in mongodb: %s", err.Error())
		return err
	}

	sendTaskTimestamp := time.Now().UTC()
	if task.TriggeredAfter != "" {
		if duration, err := time.ParseDuration(task.TriggeredAfter); err == nil {
			sendTaskTimestamp = sendTaskTimestamp.Add(duration)
		} else {
			log.Errorf("could not parse triggeredAfter property: %s", err.Error())
			// TODO how do we handle this? send event immediately or not at all?
		}
		log.Infof("queueing %s event with ID %s to be sent at %s", event.Type(), event.ID(), sendTaskTimestamp.String())
	}
	storeEvent.Time = timeutils.GetKeptnTimeStamp(sendTaskTimestamp)

	if err := sc.eventRepo.InsertEvent(eventScope.Project, *storeEvent, common.TriggeredEvent); err != nil {
		log.Errorf("Could not store event: %s", err.Error())
		return err
	}

	sc.onSequenceTaskTriggered(*storeEvent)
	if err := sc.eventDispatcher.Add(models.DispatcherEvent{TimeStamp: sendTaskTimestamp, Event: event}, false); err != nil {
		return err
	}

	return sc.taskSequenceRepo.CreateTaskExecution(eventScope.Project, models.TaskExecution{
		TaskSequenceName: taskSequenceName,
		TriggeredEventID: event.ID(),
		Stage:            eventScope.Stage,
		Service:          eventScope.Service,
		KeptnContext:     eventScope.KeptnContext,
		Task:             task,
	})
}
