package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

const maxRepoReadRetries = 5

const couldNotGetActiveSequencesErrMsg = "unable to get active task executions for project %s in stage %s for Keptn context %s: %w"
const noActiveSequencesErrMsg = "no active task executions for project %s in stage %s for Keptn context %s found"

var shipyardControllerInstance *shipyardController

//go:generate moq -pkg fake -skip-ensure -out ./fake/shipyardcontroller.go . IShipyardController
type IShipyardController interface {
	GetAllTriggeredEvents(filter common.EventFilter) ([]apimodels.KeptnContextExtendedCE, error)
	GetTriggeredEventsOfProject(project string, filter common.EventFilter) ([]apimodels.KeptnContextExtendedCE, error)
	HandleIncomingEvent(event apimodels.KeptnContextExtendedCE, waitForCompletion bool) error
	ControlSequence(controlSequence apimodels.SequenceControl) error
	StartTaskSequence(event apimodels.KeptnContextExtendedCE) error
	StartDispatchers(ctx context.Context, mode common.SDMode)
	StopDispatchers()
}

type shipyardController struct {
	eventRepo                  db.EventRepo
	sequenceExecutionRepo      db.SequenceExecutionRepo
	projectMvRepo              db.ProjectMVRepo
	eventDispatcher            IEventDispatcher
	sequenceDispatcher         ISequenceDispatcher
	sequenceTimeoutChan        chan apimodels.SequenceTimeout
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
	sequenceTimeoutChannel chan apimodels.SequenceTimeout,
	shipyardRetriever IShipyardRetriever,
) *shipyardController {
	if shipyardControllerInstance == nil {
		cbConnectionInstance := db.GetMongoDBConnectionInstance()
		shipyardControllerInstance = &shipyardController{
			eventRepo:             db.NewMongoDBEventsRepo(cbConnectionInstance),
			sequenceExecutionRepo: db.NewMongoDBSequenceExecutionRepo(cbConnectionInstance),
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
}

func (sc *shipyardController) ControlSequence(controlSequence apimodels.SequenceControl) error {
	switch controlSequence.State {
	case apimodels.AbortSequence:
		log.Info("Processing ABORT sequence control")
		return sc.cancelSequence(controlSequence)
	case apimodels.PauseSequence:
		log.Info("Processing PAUSE sequence control")
		sc.onSequencePaused(models.EventScope{
			EventData: keptnv2.EventData{
				Project: controlSequence.Project,
				Stage:   controlSequence.Stage,
			},
			KeptnContext: controlSequence.KeptnContext,
		})
		return sc.pauseSequence(controlSequence)
	case apimodels.ResumeSequence:
		log.Info("Processing RESUME sequence control")
		sc.onSequenceResumed(models.EventScope{
			EventData: keptnv2.EventData{
				Project: controlSequence.Project,
				Stage:   controlSequence.Stage,
			},
			KeptnContext: controlSequence.KeptnContext,
		})
		return sc.resumeSequence(controlSequence)
	}
	return nil
}

func (sc shipyardController) StartDispatchers(ctx context.Context, mode common.SDMode) {
	sc.eventDispatcher.Run(ctx)
	sc.sequenceDispatcher.Run(ctx, mode, sc.StartTaskSequence)
}

func (sc shipyardController) StopDispatchers() {
	sc.eventDispatcher.Stop()
	sc.sequenceDispatcher.Stop()
}

func (sc *shipyardController) HandleIncomingEvent(event apimodels.KeptnContextExtendedCE, waitForCompletion bool) error {
	statusType, err := ExtractEventKind(event)
	if err != nil {
		return err
	}
	// only create channel if waitForCompletion is set to true
	var done chan error
	if waitForCompletion {
		done = make(chan error)
	}

	log.Infof("Received event of type %s from %s", *event.Type, *event.Source)
	log.Debugf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, ObjToJSON(event))
	cb := getCompletionCallback(waitForCompletion, done)

	switch statusType {
	case string(common.TriggeredEvent):
		go func() {
			err := sc.handleSequenceTriggered(event)
			if err != nil {
				log.WithError(err).Error("Unable to handle sequence '.triggered' event")
			}
			cb(err)
		}()
	case string(common.StartedEvent), string(common.FinishedEvent):
		go func() {
			err := sc.handleTaskEvent(event)
			if err != nil {
				if errors.Is(err, ErrSequenceNotFound) || errors.Is(err, models.ErrInvalidEventScope) {
					log.Infof("Unable to handle task event: %v", err)
				} else {
					log.Errorf("Unable to handle task event: %v", err)
				}
			}
			cb(err)
		}()
	default:
		return nil
	}
	return completeEventHandler(waitForCompletion, done)
}

func completeEventHandler(waitForCompletion bool, done chan error) error {
	if waitForCompletion {
		return <-done
	}
	return nil
}

func getCompletionCallback(waitForCompletion bool, done chan error) func(err error) {
	return func(err error) {
		if waitForCompletion {
			done <- err
		}
	}
}

func (sc *shipyardController) handleSequenceTriggered(event apimodels.KeptnContextExtendedCE) error {
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

	sc.appendLatestCommitIDToEvent(*eventScope, &eventScope.WrappedEvent)
	if err := sc.eventRepo.InsertEvent(eventScope.Project, eventScope.WrappedEvent, common.TriggeredEvent); err != nil {
		log.Infof("could not store event that triggered task sequence: %s", err.Error())
	}

	inputProperties := map[string]interface{}{}

	if err := keptnv2.Decode(event.Data, &inputProperties); err != nil {
		return err
	}

	sequenceExecution := models.SequenceExecution{
		ID:       uuid.New().String(),
		Sequence: *sequence,
		Status: models.SequenceExecutionStatus{
			State:         apimodels.SequenceTriggeredState,
			PreviousTasks: []models.TaskExecutionResult{},
		},
		InputProperties: inputProperties,
		Scope:           *eventScope,
		TriggeredAt:     time.Now().UTC(),
	}
	sequenceExecution.Scope.TriggeredID = event.ID
	sequenceExecution.Scope.GitCommitID = eventScope.WrappedEvent.GitCommitID

	if sc.sequenceExecutionRepo.IsContextPaused(*eventScope) {
		sequenceExecution.Pause()
	}

	// insert the sequence execution, but only if there is no sequence with the same triggeredID already there
	if err := sc.sequenceExecutionRepo.Upsert(sequenceExecution, &models.SequenceExecutionUpsertOptions{CheckUniqueTriggeredID: true}); err != nil {
		return fmt.Errorf("could not store task sequence execution: %w", err)
	}

	sc.onSequenceTriggered(eventScope.WrappedEvent)
	err = sc.sequenceDispatcher.Add(models.QueueItem{
		Scope:     *eventScope,
		EventID:   eventScope.WrappedEvent.ID,
		Timestamp: eventScope.WrappedEvent.Time,
	})
	if errors.Is(err, ErrSequenceBlockedWaiting) {
		sc.onSequenceWaiting(eventScope.WrappedEvent)
		return nil
	}

	return err
}

func (sc *shipyardController) appendLatestCommitIDToEvent(eventScope models.EventScope, event *apimodels.KeptnContextExtendedCE) {
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

func (sc *shipyardController) handleTaskEvent(event apimodels.KeptnContextExtendedCE) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("unable to handle 'task.finished' event: %w", err)
	}

	if !keptnv2.IsTaskEventType(eventScope.EventType) {
		return nil
	}

	sequenceExecution, err := sc.getOpenSequenceExecution(*eventScope)
	if err != nil {
		return fmt.Errorf("unable to handle %s event: %w", eventScope.EventType, err)
	}

	if sequenceExecution == nil {
		log.Infof("The received %s event with keptn context %s is not accociated with a task that was previously triggered",
			eventScope.EventType, eventScope.KeptnContext)
		return ErrSequenceNotFound
	}

	if keptnv2.IsStartedEventType(*event.Type) {
		sc.onSequenceTaskStarted(eventScope.WrappedEvent)
	}

	return sc.onTaskProgress(event, *sequenceExecution, eventScope)
}

func (sc *shipyardController) onTaskProgress(event apimodels.KeptnContextExtendedCE, sequenceExecution models.SequenceExecution, eventScope *models.EventScope) error {
	taskEvent := models.TaskEvent{
		EventType: *event.Type,
		Source:    *event.Source,
		Result:    eventScope.Result,
		Status:    eventScope.Status,
		Time:      timeutils.GetKeptnTimeStamp(event.Time),
	}
	if keptnv2.IsFinishedEventType(taskEvent.EventType) {
		eventData := map[string]interface{}{}
		err := keptnv2.Decode(event.Data, &eventData)
		if err != nil {
			return err
		}
		taskEvent.Properties = eventData
	}
	updatedSequenceExecution, err := sc.sequenceExecutionRepo.AppendTaskEvent(sequenceExecution, taskEvent)
	if err != nil {
		return err
	}

	// now check if the number of .started events matches the number of finished events - if yes, that means were done
	// note: this should also work with multiple replicas because the `AppendTaskEvent` updates the list of events and returns the resulting state
	// atomically, so ONLY the thread that appended the last event to reach the completion state of the task will get the state required for further proceeding with the task sequence
	if !updatedSequenceExecution.Status.CurrentTask.IsFinished() {
		return nil
	}

	result, status := updatedSequenceExecution.CompleteCurrentTask()

	eventScope.Result = result
	eventScope.Status = status

	triggeredEventType, err := keptnv2.ReplaceEventTypeKind(eventScope.EventType, string(common.TriggeredEvent))
	if err != nil {
		return err
	}

	// Todo: as a follow up refactoring task, we can check if we can also get rid of the -triggeredEvents collection by replacing it with the items in the sequenceExecution collection - the items of that collection should contain all information about open tasks and their payloads
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

	sc.onSequenceTaskFinished(eventScope.WrappedEvent)
	return sc.proceedTaskSequence(*eventScope, *updatedSequenceExecution)
}

func (sc *shipyardController) wasTaskTriggered(eventScope models.EventScope) (bool, error) {
	taskContext, err := sc.getOpenSequenceExecution(eventScope)
	if err != nil {
		return false, err
	}
	if taskContext == nil {
		return false, nil
	}
	return true, nil
}

func (sc *shipyardController) cancelSequence(cancel apimodels.SequenceControl) error {
	sc.onSequenceAborted(models.EventScope{
		KeptnContext: cancel.KeptnContext,
		EventData:    keptnv2.EventData{Project: cancel.Project, Stage: cancel.Stage},
	})
	sequenceExecutions, err := sc.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{
		KeptnContext: cancel.KeptnContext,
		EventData: keptnv2.EventData{
			Project: cancel.Project,
			Stage:   cancel.Stage,
		},
	}})

	if err != nil {
		return fmt.Errorf(couldNotGetActiveSequencesErrMsg, cancel.Project, cancel.Stage, cancel.KeptnContext, err)
	}

	if len(sequenceExecutions) == 0 {
		log.Infof(noActiveSequencesErrMsg, cancel.Project, cancel.Stage, cancel.KeptnContext)
		return nil
	}

	err = sc.sequenceDispatcher.Remove(
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

	// delete all open .triggered events for the task sequence
	for _, sequenceExecution := range sequenceExecutions {
		err := sc.eventRepo.DeleteEvent(cancel.Project, sequenceExecution.Status.CurrentTask.TriggeredID, common.TriggeredEvent)
		if err != nil {
			// log the error, but continue
			log.WithError(err).Error("could not delete event")
		}

		if err := sc.forceTaskSequenceCompletion(sequenceExecution); err != nil {
			log.Errorf("Could not complete sequence execution %s: %v", sequenceExecution.Scope.KeptnContext, err)
		}
	}
	return nil
}

func (sc *shipyardController) pauseSequence(pause apimodels.SequenceControl) error {
	scope := models.EventScope{
		KeptnContext: pause.KeptnContext,
		EventData: keptnv2.EventData{
			Project: pause.Project,
			Stage:   pause.Stage,
		},
	}
	if err := sc.sequenceExecutionRepo.PauseContext(scope); err != nil {
		log.Errorf("Could not pause global scope of sequence context: %v", err)
	}
	sequenceExecutions, err := sc.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: scope})

	if err != nil {
		return fmt.Errorf(couldNotGetActiveSequencesErrMsg, pause.Project, pause.Stage, pause.KeptnContext, err)
	}

	if len(sequenceExecutions) == 0 {
		log.Infof(noActiveSequencesErrMsg, pause.Project, pause.Stage, pause.KeptnContext)
		return nil
	}

	for _, sequenceExecution := range sequenceExecutions {
		if !sequenceExecution.Pause() {
			continue
		}
		_, err = sc.sequenceExecutionRepo.UpdateStatus(sequenceExecution)
		if err != nil {
			log.Errorf("Could not update sequence execution state: %v", err)
		}
	}
	return nil
}

func (sc *shipyardController) resumeSequence(resume apimodels.SequenceControl) error {
	scope := models.EventScope{
		KeptnContext: resume.KeptnContext,
		EventData: keptnv2.EventData{
			Project: resume.Project,
			Stage:   resume.Stage,
		},
	}
	if err := sc.sequenceExecutionRepo.ResumeContext(scope); err != nil {
		log.Errorf("Could not pause global scope of sequence context: %v", err)
	}
	sequenceExecutions, err := sc.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: scope})

	if err != nil {
		return fmt.Errorf(noActiveSequencesErrMsg, resume.Project, resume.Stage, resume.KeptnContext)
	}

	if len(sequenceExecutions) == 0 {
		log.Infof("no active sequence executions for context %s found.", resume.KeptnContext)
		return nil
	}

	for _, sequenceExecution := range sequenceExecutions {
		if !sequenceExecution.Resume() {
			continue
		}
		_, err = sc.sequenceExecutionRepo.UpdateStatus(sequenceExecution)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sc *shipyardController) forceTaskSequenceCompletion(sequenceExecution models.SequenceExecution) error {
	scope := sequenceExecution.Scope

	scope.Result = keptnv2.ResultPass
	scope.Status = keptnv2.StatusAborted

	return sc.completeTaskSequence(scope, sequenceExecution, apimodels.SequenceFinished)
}

func (sc *shipyardController) timeoutSequence(timeout apimodels.SequenceTimeout) error {
	log.Infof("sequence %s has been timed out", timeout.KeptnContext)
	eventScope, err := models.NewEventScope(timeout.LastEvent)
	if err != nil {
		return err
	}

	eventScope.Status = keptnv2.StatusErrored
	eventScope.Result = keptnv2.ResultFailed
	eventScope.Message = fmt.Sprintf("sequence timed out while waiting for task %s to receive a correlating .started or .finished event", *timeout.LastEvent.Type)

	sequenceExecutions, err := sc.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		CurrentTriggeredID: timeout.LastEvent.ID,
		Scope:              *eventScope,
	})

	if err != nil {
		return fmt.Errorf("could not sequence executions associated to eventID %s: %w", timeout.LastEvent.ID, err)
	}

	if len(sequenceExecutions) == 0 {
		log.Infof("No task executions associated with eventID %s found", timeout.LastEvent.ID)
		return nil
	}

	sequenceExecution := sequenceExecutions[0]
	sc.onSequenceTimeout(timeout.LastEvent)

	if err := sc.completeTaskSequence(sequenceExecution.Scope, sequenceExecution, apimodels.TimedOut); err != nil {
		return err
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

func (sc *shipyardController) StartTaskSequence(event apimodels.KeptnContextExtendedCE) error {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return err
	}

	_, taskSequenceName, _, err := keptnv2.ParseSequenceEventType(*event.Type)
	if err != nil {
		return err
	}

	sc.onSequenceStarted(event)

	sequenceExecutions, err := sc.sequenceExecutionRepo.Get(
		models.SequenceExecutionFilter{
			Scope:  *eventScope,
			Name:   taskSequenceName,
			Status: []string{apimodels.SequenceTriggeredState},
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
	sequenceExecution.Status.State = apimodels.SequenceStartedState
	updatedSequenceExecution, err := sc.sequenceExecutionRepo.UpdateStatus(sequenceExecution)
	if err != nil {
		msg := fmt.Sprintf("could not update sequence execution state %s: %s", taskSequenceName, err.Error())
		return sc.triggerSequenceFailed(*eventScope, msg, taskSequenceName)
	}
	sc.onSequenceStarted(event)
	return sc.proceedTaskSequence(*eventScope, *updatedSequenceExecution)
}

func (sc *shipyardController) getOpenSequenceExecution(eventScope models.EventScope) (*models.SequenceExecution, error) {
	sequenceExecutions, err := sc.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData:    keptnv2.EventData{Project: eventScope.Project},
			KeptnContext: eventScope.KeptnContext,
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

func (sc *shipyardController) GetAllTriggeredEvents(filter common.EventFilter) ([]apimodels.KeptnContextExtendedCE, error) {
	projects, err := sc.projectMvRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []apimodels.KeptnContextExtendedCE{}
	for _, project := range projects {
		events, err := sc.eventRepo.GetEvents(project.ProjectName, filter, common.TriggeredEvent)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (sc *shipyardController) GetTriggeredEventsOfProject(projectName string, filter common.EventFilter) ([]apimodels.KeptnContextExtendedCE, error) {
	project, err := sc.projectMvRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	} else if project == nil {
		return nil, ErrProjectNotFound
	}
	events, err := sc.eventRepo.GetEvents(projectName, filter, common.TriggeredEvent)
	if err != nil && err != db.ErrNoEventFound {
		return nil, err
	} else if err != nil && errors.Is(err, db.ErrNoEventFound) {
		return []apimodels.KeptnContextExtendedCE{}, nil
	}
	return events, nil
}

func (sc *shipyardController) proceedTaskSequence(eventScope models.EventScope, sequenceExecution models.SequenceExecution) error {
	// get the input for the .triggered event that triggered the previous sequence and append it to the list of previous events to gather all required data for the next stage
	inputEvent, err := sc.getSequenceTriggeredEvent(sequenceExecution)
	if err != nil {
		return err
	}

	task := sequenceExecution.GetNextTaskOfSequence()
	if task == nil {
		// task sequence completed -> send .finished event and check if a new task sequence should be triggered by the completion
		err = sc.completeTaskSequence(eventScope, sequenceExecution, apimodels.SequenceFinished)
		if err != nil {
			log.Errorf("Could not complete task sequence %s.%s with KeptnContext %s: %s", eventScope.Stage, sequenceExecution.Sequence.Name, eventScope.KeptnContext, err.Error())
			return err
		}
		return sc.triggerNextTaskSequences(eventScope, inputEvent, sequenceExecution)
	}

	return sc.triggerTask(eventScope, sequenceExecution, *task)
}

// this function retrieves the .triggered event for the task sequence and appends its properties to the existing .finished events
// this ensures that all parameters set in the .triggered event are received by all execution plane services, instead of just the first one
func (sc *shipyardController) getSequenceTriggeredEvent(sequenceExecution models.SequenceExecution) (*apimodels.KeptnContextExtendedCE, error) {
	triggeredEvent, err := sc.eventRepo.GetTaskSequenceTriggeredEvent(
		sequenceExecution.Scope,
		sequenceExecution.Sequence.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load event that triggered task sequence %s.%s with KeptnContext %s: %w", sequenceExecution.Scope.Stage, sequenceExecution.Sequence.Name, sequenceExecution.Scope.KeptnContext, err)
	}

	return triggeredEvent, nil
}

func (sc *shipyardController) triggerNextTaskSequences(eventScope models.EventScope, inputEvent *apimodels.KeptnContextExtendedCE, completedSequence models.SequenceExecution) error {
	shipyard, err := sc.shipyardRetriever.GetCachedShipyard(eventScope.Project)
	if err != nil {
		return err
	}
	nextSequences := GetTaskSequencesByTrigger(eventScope, completedSequence.Sequence.Name, shipyard, completedSequence.GetLastTaskExecutionResult().Name)

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

		err := sc.sendTaskSequenceTriggeredEvent(newScope, sequence.Sequence.Name, completedSequence)
		if err != nil {
			log.Errorf("could not send event %s.%s.triggered: %s",
				newScope.Stage, sequence.Sequence.Name, err.Error())
			continue
		}
	}
	return nil
}

func (sc *shipyardController) completeTaskSequence(eventScope models.EventScope, sequenceExecution models.SequenceExecution, reason string) error {
	sequenceExecution.Status.State = reason
	_, err := sc.sequenceExecutionRepo.UpdateStatus(sequenceExecution)

	if err != nil {
		return err
	}

	log.Infof("Deleting all task.finished events of task sequence %s with context %s", sequenceExecution.Sequence.Name, sequenceExecution.Scope.KeptnContext)
	if err := sc.eventRepo.DeleteAllFinishedEvents(eventScope); err != nil {
		return err
	}
	return sc.sendTaskSequenceFinishedEvent(eventScope, sequenceExecution.Sequence.Name, sequenceExecution.Scope.TriggeredID)
}

func (sc *shipyardController) triggerTask(eventScope models.EventScope, sequenceExecution models.SequenceExecution, task keptnv2.Task) error {
	eventPayload := sequenceExecution.GetNextTriggeredEventData()

	event := common.CreateEventWithPayload(eventScope.KeptnContext, "", keptnv2.GetTriggeredEventType(task.Name), eventPayload)
	event.SetExtension("gitcommitid", sequenceExecution.Scope.GitCommitID)

	storeEvent := &apimodels.KeptnContextExtendedCE{}
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
	storeEvent.Time = sendTaskTimestamp

	if err := sc.eventRepo.InsertEvent(eventScope.Project, *storeEvent, common.TriggeredEvent); err != nil {
		log.Errorf("Could not store event: %s", err.Error())
		return err
	}

	sc.onSequenceTaskTriggered(*storeEvent)

	sequenceExecution.Status.CurrentTask = models.TaskExecutionState{
		Name:        task.Name,
		TriggeredID: storeEvent.ID,
		Events:      []models.TaskEvent{},
	}

	// special handling for approval events
	if task.Name == "approval" {
		sequenceExecution.Status.State = apimodels.SequenceWaitingForApprovalState
	}

	if err := sc.sequenceExecutionRepo.Upsert(sequenceExecution, nil); err != nil {
		return err
	}
	if err := sc.eventDispatcher.Add(models.DispatcherEvent{TimeStamp: sendTaskTimestamp, Event: event}, false); err != nil {
		return err
	}
	return nil
}

func (sc *shipyardController) sendTaskSequenceTriggeredEvent(eventScope *models.EventScope, taskSequenceName string, completedSequence models.SequenceExecution) error {

	mergedPayload := completedSequence.GetNextTriggeredEventData()

	mergedPayload["stage"] = eventScope.Stage

	eventType := eventScope.Stage + "." + taskSequenceName

	event := common.CreateEventWithPayload(eventScope.KeptnContext, "", keptnv2.GetTriggeredEventType(eventType), mergedPayload)

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
