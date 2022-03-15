package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
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
	projectMvRepo              db.ProjectMVRepo
	eventDispatcher            IEventDispatcher
	sequenceDispatcher         ISequenceDispatcher
	sequenceTimeoutChan        chan models.SequenceTimeout
	sequenceTriggeredHooks     []sequencehooks.ISequenceTriggeredHook
	sequenceStartedHooks       []sequencehooks.ISequenceStartedHook
	sequenceWaitingHooks       []sequencehooks.ISequenceWaitingHook
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
	sc.eventDispatcher.Run(context.Background())
	sc.sequenceDispatcher.Run(context.Background(), sc.StartTaskSequence)
}

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
			err := sc.handleTaskStarted(event)
			if err != nil {
				log.WithError(err).Error("Unable to handle task '.started' event")
			}
			done <- err
		}()
	case string(common.FinishedEvent):
		go func() {
			err := sc.handleTaskFinished(event)
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
	_, err = GetTaskSequenceInStage(eventScope.Stage, taskSequenceName, shipyard)
	if err != nil {
		msg := fmt.Sprintf("Unable to start sequence %s: %v", taskSequenceName, err)
		log.Error(msg)
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}

	sc.appendLatestCommitIDToEvent(*eventScope, &eventScope.WrappedEvent)
	if err := sc.eventRepo.InsertEvent(eventScope.Project, eventScope.WrappedEvent, common.TriggeredEvent); err != nil {
		log.Infof("could not store event that triggered task sequence: %s", err.Error())
	}

	sc.onSequenceTriggered(eventScope.WrappedEvent)
	err = sc.sequenceDispatcher.Add(models.QueueItem{
		Scope:     *eventScope,
		EventID:   eventScope.WrappedEvent.ID,
		Timestamp: common.ParseTimestamp(eventScope.WrappedEvent.Time, nil),
	})
	if err == ErrSequenceBlockedWaiting {
		sc.onSequenceWaiting(eventScope.WrappedEvent)
		return nil
	}

	return err
}

func (sc *shipyardController) appendLatestCommitIDToEvent(eventScope models.EventScope, event *models.Event) {
	// get the latest git commit ID for the stage if it is not specified in the event
	if eventScope.WrappedEvent.GitCommitID == "" {
		latestGitCommitID, err := sc.shipyardRetriever.GetLatestCommitID(eventScope.Project, eventScope.Stage)
		if err != nil {
			// log the error, but having no commit ID should not prevent the sequence from being executed
			log.Errorf("Could not determine latest commit ID for stage %s in project %s", eventScope.Project, eventScope.Stage)
		}
		event.GitCommitID = latestGitCommitID
	}
}

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

func (sc *shipyardController) handleTaskFinished(event models.Event) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("unable to handle 'task.finished' event: %w", err)
	}

	if !keptnv2.IsTaskEventType(eventScope.EventType) {
		return nil
	}

	wasTriggered, err := sc.wasTaskTriggered(*eventScope)
	if err != nil {
		return fmt.Errorf("unable to handle %s event: %w", eventScope.EventType, err)
	}

	if !wasTriggered {
		log.Infof("The received %s event with keptn context %s is not accociated with a task that was previously triggered",
			eventScope.EventType, eventScope.KeptnContext)
		return nil
	}

	common.LockServiceInStageOfProject(eventScope.Project, eventScope.Stage, eventScope.Service+":taskFinisher")
	defer common.UnlockServiceInStageOfProject(eventScope.Project, eventScope.Stage, eventScope.Service+":taskFinisher")

	startedEvents, err := sc.eventRepo.GetStartedEventsForTriggeredID(*eventScope)
	if err != nil {
		return fmt.Errorf("error while retrieving matching '.started' event for event %s with triggered ID %s: %w", eventScope.WrappedEvent.ID, eventScope.TriggeredID, err)
	}
	if len(startedEvents) == 0 {
		return fmt.Errorf("no matching '.started' event found for event %s with triggered ID %s", eventScope.WrappedEvent.ID, eventScope.TriggeredID)
	}

	err = sc.eventRepo.InsertEvent(eventScope.Project, eventScope.WrappedEvent, common.FinishedEvent)
	if err != nil {
		log.Errorf("unable to store %s event: %v ", eventScope.EventType, err.Error())
	}

	err = sc.eventRepo.DeleteEvent(eventScope.Project, startedEvents[0].ID, common.StartedEvent)
	if err != nil {
		return fmt.Errorf("unable to delete associated task '.started' event with ID %s: %w", startedEvents[0].ID, err)
	}

	// if this was the last '.started' event
	if len(startedEvents) == 1 {
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

		// if the previously deleted '.started' event was the last, the '.triggered' event can be removed
		err = sc.eventRepo.DeleteEvent(eventScope.Project, triggeredEvents[0].ID, common.TriggeredEvent)
		if err != nil {
			return fmt.Errorf("unable to delete associated task '.triggered' event with ID %s: %w", eventScope.TriggeredID, err)
		}

		taskExecution, err := sc.getOpenTaskExecution(*eventScope)
		if err != nil {
			return fmt.Errorf("unable to get open task execution: %w", err)
		}

		shipyard, err := sc.shipyardRetriever.GetCachedShipyard(eventScope.Project)
		if err != nil {
			return fmt.Errorf("unable to fetch shipyard: %w", err)
		}

		sequence, err := GetTaskSequenceInStage(taskExecution.Stage, taskExecution.TaskSequenceName, shipyard)
		if err != nil {
			return fmt.Errorf("unable to find task %s.%s found in shipyard: %w", taskExecution.Stage, taskExecution.TaskSequenceName, err)
		}

		finishedEventsData, err := sc.getFinishedEventData(*eventScope)
		if err != nil {
			return fmt.Errorf("unable to gather task '.finished' event data: %w", err)
		}

		sc.onSequenceTaskFinished(eventScope.WrappedEvent)
		return sc.proceedTaskSequence(*eventScope, sequence, finishedEventsData, taskExecution)
	}
	return nil
}

func (sc *shipyardController) wasTaskTriggered(eventScope models.EventScope) (bool, error) {
	taskContext, err := sc.getOpenTaskExecution(eventScope)
	if err != nil {
		return false, err
	}
	if taskContext == nil {
		return false, nil
	}
	return true, nil
}

func (sc *shipyardController) cancelSequence(cancel models.SequenceControl) error {
	sc.onSequenceAborted(models.EventScope{
		KeptnContext: cancel.KeptnContext,
		EventData:    keptnv2.EventData{Project: cancel.Project, Stage: cancel.Stage},
	})
	taskExecutions, err := sc.taskSequenceRepo.GetTaskExecutions(cancel.Project,
		models.TaskExecution{
			KeptnContext: cancel.KeptnContext,
			Stage:        cancel.Stage,
		})
	if err != nil {
		return fmt.Errorf("unable to get active task executions for project %s in stage %s for keptn context %s", cancel.Project, cancel.Stage, cancel.KeptnContext)
	}
	if len(taskExecutions) == 0 {
		log.Infof("no active task execution for context %s found. Trying to remove it from the queue", cancel.KeptnContext)
		return sc.cancelQueuedSequence(cancel)
	}

	err = sc.taskSequenceRepo.DeleteTaskExecutions(cancel.KeptnContext, cancel.Project, cancel.Stage)
	if err != nil {
		// log the error, but continue with the rest
		log.Errorf("Could not delete open task executions for sequence %s in proect %s: %v", cancel.KeptnContext, cancel.Project, err)
	}

	// delete all open .triggered events for the task sequence
	for _, sequenceEvent := range taskExecutions {
		err := sc.eventRepo.DeleteEvent(cancel.Project, sequenceEvent.TriggeredEventID, common.TriggeredEvent)
		if err != nil {
			// log the error, but continue
			log.WithError(err).Error("could not delete event")
		}
	}

	lastTaskOfSequence := getLastTaskOfSequence(taskExecutions)
	sequenceTriggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(models.EventScope{
		EventData: keptnv2.EventData{
			Project: cancel.Project,
			Stage:   lastTaskOfSequence.Stage,
		},
		KeptnContext: cancel.KeptnContext,
	}, lastTaskOfSequence.TaskSequenceName)
	if err != nil {
		return err
	}

	if sequenceTriggeredEvent != nil {
		return sc.forceTaskSequenceCompletion(*sequenceTriggeredEvent, lastTaskOfSequence.TaskSequenceName)
	}
	return nil
}

func (sc *shipyardController) forceTaskSequenceCompletion(sequenceTriggeredEvent models.Event, taskSequenceName string) error {
	scope, err := models.NewEventScope(sequenceTriggeredEvent)
	if err != nil {
		return err
	}

	scope.Result = keptnv2.ResultPass
	scope.Status = keptnv2.StatusAborted

	return sc.completeTaskSequence(*scope, taskSequenceName, sequenceTriggeredEvent.ID)
}

func (sc *shipyardController) cancelQueuedSequence(cancel models.SequenceControl) error {
	// first, remove the sequence from the queue
	err := sc.sequenceDispatcher.Remove(
		models.EventScope{
			EventData: keptnv2.EventData{
				Project: cancel.Project,
				Stage:   cancel.Stage,
			},
			KeptnContext: cancel.KeptnContext,
		},
	)

	if err != nil {
		log.WithError(err).Errorf("could not remove sequence %s from sequence queue", cancel.KeptnContext)
	}

	// theoretically, we should not have any open task executions at this point, but call the delete function anyway,
	// just to be 100% sure no entries of this collection are blocking any new sequences
	err = sc.taskSequenceRepo.DeleteTaskExecutions(cancel.KeptnContext, cancel.Project, cancel.Stage)
	if err != nil {
		// log the error, but continue with the rest
		log.Errorf("Could not delete open task executions for sequence %s in proect %s: %v", cancel.KeptnContext, cancel.Project, err)
	}

	events, err := sc.eventRepo.GetEvents(
		cancel.Project,
		common.EventFilter{KeptnContext: &cancel.KeptnContext, Stage: &cancel.Stage},
		common.TriggeredEvent,
	)
	if err != nil {
		// if the sequence.triggered event is not available anymore, we cannot send a referencing .finished event
		if err == db.ErrNoEventFound {
			log.Infof("No sequence.triggered event for sequence %s available anymore.", cancel.KeptnContext)
			return nil
		}
		return err
	} else if len(events) == 0 {
		return ErrSequenceNotFound
	}
	// the first event of the context should be a task sequence event that contains the sequence name
	sequenceTriggeredEvent := events[0]
	if !keptnv2.IsSequenceEventType(*sequenceTriggeredEvent.Type) {
		// if the sequence.triggered event is not available anymore, we cannot send a referencing .finished event
		log.Infof("No sequence.triggered event for sequence %s available anymore.", cancel.KeptnContext)
		return nil
	}
	_, sequenceName, _, err := keptnv2.ParseSequenceEventType(*sequenceTriggeredEvent.Type)
	if err != nil {
		return err
	}

	return sc.forceTaskSequenceCompletion(sequenceTriggeredEvent, sequenceName)
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

	taskExecutions, err := sc.taskSequenceRepo.GetTaskExecutions(eventScope.Project, models.TaskExecution{TriggeredEventID: timeout.LastEvent.ID})
	if err != nil {
		return fmt.Errorf("Could not retrieve task executions associated to eventID %s: %s", timeout.LastEvent.ID, err.Error())
	}

	if len(taskExecutions) == 0 {
		log.Infof("No task executions associated with eventID %s found", timeout.LastEvent.ID)
		return nil
	}
	taskContext := taskExecutions[0]
	sc.onSequenceTimeout(timeout.LastEvent)
	taskSequenceTriggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(*eventScope, taskContext.TaskSequenceName)
	if err != nil {
		return err
	}
	if taskSequenceTriggeredEvent != nil {
		if err := sc.completeTaskSequence(*eventScope, taskContext.TaskSequenceName, taskSequenceTriggeredEvent.ID); err != nil {
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

	shipyard, err := sc.shipyardRetriever.GetCachedShipyard(eventScope.Project)
	if err != nil {
		return err
	}

	_, taskSequenceName, _, err := keptnv2.ParseSequenceEventType(*event.Type)
	if err != nil {
		return err
	}

	taskSequence, err := GetTaskSequenceInStage(eventScope.Stage, taskSequenceName, shipyard)
	if err != nil {
		msg := fmt.Sprintf("could not get definition of task sequence %s: %s", taskSequenceName, err.Error())
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}
	sc.onSequenceStarted(event)

	return sc.proceedTaskSequence(*eventScope, taskSequence, []interface{}{}, nil)
}

func (sc *shipyardController) getOpenTaskExecution(eventScope models.EventScope) (*models.TaskExecution, error) {
	for i := 0; i <= maxRepoReadRetries; i++ {
		taskExecutions, err := sc.taskSequenceRepo.GetTaskExecutions(eventScope.Project, models.TaskExecution{TriggeredEventID: eventScope.TriggeredID})
		if err != nil {
			msg := "Could not retrieve task executions associated to eventID " + eventScope.TriggeredID + ": " + err.Error()
			log.Error(msg)
			return nil, errors.New(msg)
		}

		if len(taskExecutions) == 0 {
			log.Infof("No task executions associated with eventID %s found", eventScope.TriggeredID)
			<-time.After(2 * time.Second)
		} else {
			taskContext := taskExecutions[0]
			return &taskContext, nil
		}
	}
	return nil, nil
}

func (sc *shipyardController) getFinishedEventData(eventScope models.EventScope) ([]interface{}, error) {
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

func (sc *shipyardController) proceedTaskSequence(eventScope models.EventScope, taskSequence *keptnv2.Sequence, eventHistory []interface{}, previousTask *models.TaskExecution) error {

	// get the input for the .triggered event that triggered the previous sequence and append it to the list of previous events to gather all required data for the next stage
	inputEvent, eventHistory, err := sc.appendTriggerEventProperties(eventScope, taskSequence, eventHistory)
	if err != nil {
		return err
	}

	task := GetNextTaskOfSequence(taskSequence, previousTask, &eventScope, eventHistory)
	if task == nil {
		// task sequence completed -> send .finished event and check if a new task sequence should be triggered by the completion
		err = sc.completeTaskSequence(eventScope, taskSequence.Name, inputEvent.ID)
		if err != nil {
			log.Errorf("Could not complete task sequence %s.%s with KeptnContext %s: %s", eventScope.Stage, taskSequence.Name, eventScope.KeptnContext, err.Error())
			return err
		}

		return sc.triggerNextTaskSequences(eventScope, taskSequence, eventHistory, inputEvent, previousTask.Task.Name)
	}
	// propagate the git commit id to the next task
	eventScope.GitCommitID = inputEvent.GitCommitID
	return sc.sendTaskTriggeredEvent(eventScope, taskSequence.Name, *task, eventHistory)
}

// this function retrieves the .triggered event for the task sequence and appends its properties to the existing .finished events
// this ensures that all parameters set in the .triggered event are received by all execution plane services, instead of just the first one
func (sc *shipyardController) appendTriggerEventProperties(eventScope models.EventScope, taskSequence *keptnv2.Sequence, eventHistory []interface{}) (*models.Event, []interface{}, error) {
	triggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(eventScope, taskSequence.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load event that triggered task sequence %s.%s with KeptnContext %s: %w", eventScope.Stage, taskSequence.Name, eventScope.KeptnContext, err)
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

func (sc *shipyardController) triggerNextTaskSequences(eventScope models.EventScope, completedSequence *keptnv2.Sequence, eventHistory []interface{}, inputEvent *models.Event, previousTask string) error {
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

func (sc *shipyardController) completeTaskSequence(eventScope models.EventScope, taskSequenceName, triggeredID string) error {
	err := sc.taskSequenceRepo.DeleteTaskExecution(eventScope.KeptnContext, eventScope.Project, eventScope.Stage, taskSequenceName)
	if err != nil {
		return err
	}
	log.Infof("Deleting all task.finished events of task sequence %s with context %s", taskSequenceName, eventScope.KeptnContext)
	if err := sc.eventRepo.DeleteAllFinishedEvents(eventScope); err != nil {
		return err
	}
	return sc.sendTaskSequenceFinishedEvent(eventScope, taskSequenceName, triggeredID)
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
	sc.appendLatestCommitIDToEvent(*eventScope, toEvent)
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

	if eventScope.GitCommitID != "" {
		event.SetExtension("gitcommitid", eventScope.GitCommitID)
	}
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
