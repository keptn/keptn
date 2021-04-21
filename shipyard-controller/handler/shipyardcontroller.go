package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/common/eventutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"net/url"
	"time"
)

const maxRepoReadRetries = 30

var errNoMatchingEvent = errors.New("no matching event found")
var shipyardControllerInstance *shipyardController

type IShipyardController interface {
	GetAllTriggeredEvents(filter common.EventFilter) ([]models.Event, error)
	GetTriggeredEventsOfProject(project string, filter common.EventFilter) ([]models.Event, error)
	HandleIncomingEvent(event models.Event) error
}

type shipyardController struct {
	projectRepo        db.ProjectRepo
	eventRepo          db.EventRepo
	taskSequenceRepo   db.TaskSequenceRepo
	eventsDbOperations db.EventsDbOperations
	logger             *keptncommon.Logger
	eventDispatcher    IEventDispatcher
}

// TODO inject eventDispatcher
func GetShipyardControllerInstance(eventDispatcher IEventDispatcher) *shipyardController {
	if shipyardControllerInstance == nil {
		logger := keptncommon.NewLogger("", "", "shipyard-controller")
		shipyardControllerInstance = &shipyardController{
			projectRepo:      &db.MongoDBProjectsRepo{Logger: logger},
			eventRepo:        &db.MongoDBEventsRepo{Logger: logger},
			taskSequenceRepo: &db.TaskSequenceMongoDBRepo{Logger: logger},
			eventsDbOperations: &db.ProjectsMaterializedView{
				ProjectRepo:     &db.MongoDBProjectsRepo{Logger: logger},
				EventsRetriever: &db.MongoDBEventsRepo{Logger: logger},
				Logger:          logger,
			},
			logger:          logger,
			eventDispatcher: eventDispatcher,
		}
	}
	return shipyardControllerInstance
}

func (sc *shipyardController) HandleIncomingEvent(event models.Event) error {
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(event.Data, eventData)
	if err != nil {
		sc.logger.Error("Could not parse event data: " + err.Error())
		return err
	}

	if eventData.Project != "" && eventData.Stage != "" && eventData.Service != "" {
		go func() {
			common.LockProject(eventData.Project)
			defer common.UnlockProject(eventData.Project)
			if err := sc.eventsDbOperations.UpdateEventOfService(
				event.Data,
				*event.Type,
				event.Shkeptncontext,
				event.ID,
				event.Triggeredid,
			); err != nil {
				sc.logger.Error(fmt.Sprintf("could not update event for project %s: %s", eventData.Project, err.Error()))
			}
		}()
	}

	statusType, err := eventutils.ParseEventKind(*event.Type)
	if err != nil {
		return err
	}
	switch statusType {
	case string(common.TriggeredEvent):
		return sc.handleTriggeredEvent(event)
	case string(common.StartedEvent):
		return sc.handleStartedEvent(event)
	case string(common.FinishedEvent):
		return sc.handleFinishedEvent(event)
	default:
		return nil
	}
}

func (sc *shipyardController) handleStartedEvent(event models.Event) error {
	sc.logger.Info("Received .started event: " + *event.Type)
	// eventScope contains all properties (project, stage, service) that are needed to determine the current state within a task sequence
	// if those are not present the next action can not be determined
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, printObject(event)))

	triggeredEventType, err := eventutils.ReplaceEventTypeKind(*event.Type, string(common.TriggeredEvent))
	if err != nil {
		return err
	}

	// get corresponding 'triggered' event for the incoming 'started' event
	filter := common.EventFilter{
		Type: triggeredEventType,
		ID:   &event.Triggeredid,
	}

	events, err := sc.getEvents(eventScope.Project, filter, common.TriggeredEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		sc.logger.Error(msg)
		return errors.New(msg)
	} else if len(events) == 0 {
		msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		sc.logger.Error(msg)
		return errNoMatchingEvent
	}

	return sc.eventRepo.InsertEvent(eventScope.Project, event, common.StartedEvent)
}

func (sc *shipyardController) handleTriggeredEvent(event models.Event) error {
	if *event.Source == "shipyard-controller" {
		sc.logger.Info("Received event from myself. Ignoring.")
		return nil
	}

	// eventScope contains all properties (project, stage, service) that are needed to determine the current state within a task sequence
	// if those are not present the next action can not be determined
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, printObject(event)))

	sc.logger.Info("received event of type " + *event.Type + " from " + *event.Source)
	sc.logger.Info("Checking if .triggered event should start a sequence in project " + eventScope.Project)
	// get stage and taskSequenceName - cannot tell if this is actually a task sequence triggered event though

	stageName, taskSequenceName, _, err := eventutils.ParseSequenceEventType(*event.Type)
	if err != nil {
		return err
	}

	shipyard, err := common.GetShipyard(eventScope.Project)
	if err != nil {
		msg := "could not retrieve shipyard: " + err.Error()
		sc.logger.Error(msg)
		return sc.sendTaskSequenceFinishedEvent(&models.EventScope{
			EventData: keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   eventScope.Stage,
				Service: eventScope.Service,
				Labels:  eventScope.Labels,
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultFailed,
				Message: msg,
			},
			KeptnContext: event.Shkeptncontext,
		}, taskSequenceName, event.ID)
	}

	// update the shipyard content of the project
	shipyardContent, err := yaml.Marshal(shipyard)
	if err != nil {
		// log the error but continue
		sc.logger.Error(fmt.Sprintf("could not encode shipyard file of project %s: %s", eventScope.Project, err.Error()))
	}
	if err := sc.eventsDbOperations.UpdateShipyard(eventScope.Project, string(shipyardContent)); err != nil {
		// log the error but continue
		sc.logger.Error(fmt.Sprintf("could not update shipyard content of project %s: %s", eventScope.Project, err.Error()))
	}

	// validate the shipyard version - only shipyard files following the '0.2.0' spec are supported by the shipyard controller
	err = common.ValidateShipyardVersion(shipyard)
	if err != nil {
		// if the validation has not been successful: send a <task-sequence>.finished event with status=errored
		sc.logger.Error("invalid shipyard version: " + err.Error())
		return sc.sendTaskSequenceFinishedEvent(&models.EventScope{
			EventData: keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   eventScope.Stage,
				Service: eventScope.Service,
				Labels:  eventScope.Labels,
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultFailed,
				Message: "Found shipyard.yaml with invalid version. Please upgrade the shipyard.yaml of the project using the Keptn CLI: 'keptn upgrade project " + eventScope.Project + " --shipyard'. '",
			},
			KeptnContext: event.Shkeptncontext,
		}, taskSequenceName, event.ID)
	}

	taskSequence, err := sc.getTaskSequenceInStage(stageName, taskSequenceName, shipyard)
	if err != nil {
		sc.logger.Error(err.Error())
		return err
	}

	if err := sc.eventRepo.InsertEvent(eventScope.Project, event, common.TriggeredEvent); err != nil {
		sc.logger.Info("could not store event that triggered task sequence: " + err.Error())
	}

	eventScope.Stage = stageName
	return sc.proceedTaskSequence(eventScope, taskSequence, shipyard, []interface{}{}, "")
}

func (sc *shipyardController) handleFinishedEvent(event models.Event) error {
	if *event.Source == "shipyard-controller" {
		sc.logger.Info("Received event from myself. Ignoring.")
		return nil
	}

	// eventScope contains all properties (project, stage, service) that are needed to determine the current state within a task sequence
	// if those are not present the next action can not be determined
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, printObject(event)))

	startedEventType, err := eventutils.ReplaceEventTypeKind(*event.Type, string(common.StartedEvent))
	if err != nil {
		return err
	}
	// get corresponding 'started' event for the incoming 'finished' event
	filter := common.EventFilter{
		Type:        startedEventType,
		TriggeredID: &event.Triggeredid,
	}
	startedEvents, err := sc.getEvents(eventScope.Project, filter, common.StartedEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		sc.logger.Error(msg)
		return errors.New(msg)
	} else if len(startedEvents) == 0 {
		msg := "no matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		sc.logger.Error(msg)
		return errNoMatchingEvent
	}

	// persist the .finished event
	err = sc.eventRepo.InsertEvent(eventScope.Project, event, common.FinishedEvent)
	if err != nil {
		sc.logger.Error("Could not store .finished event: " + err.Error())
	}

	for _, startedEvent := range startedEvents {
		if *event.Source == *startedEvent.Source {
			err = sc.eventRepo.DeleteEvent(eventScope.Project, startedEvent.ID, common.StartedEvent)
			if err != nil {
				msg := "could not delete '.started' event with ID " + startedEvent.ID + ": " + err.Error()
				sc.logger.Error(msg)
				return errors.New(msg)
			}
		}
	}
	// check if this was the last '.started' event
	if len(startedEvents) == 1 {
		triggeredEventType, err := eventutils.ReplaceEventTypeKind(*event.Type, string(common.TriggeredEvent))
		if err != nil {
			return err
		}

		triggeredEventFilter := common.EventFilter{
			Type: triggeredEventType,
			ID:   &event.Triggeredid,
		}
		triggeredEvents, err := sc.getEvents(eventScope.Project, triggeredEventFilter, common.TriggeredEvent, maxRepoReadRetries)
		if err != nil {
			msg := "could not retrieve '.triggered' event with ID " + event.Triggeredid + ": " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}
		if len(triggeredEvents) == 0 {
			msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
			sc.logger.Error(msg)
			return errNoMatchingEvent
		}
		// if the previously deleted '.started' event was the last, the '.triggered' event can be removed
		sc.logger.Info("triggered event will be deleted")
		err = sc.eventRepo.DeleteEvent(eventScope.Project, triggeredEvents[0].ID, common.TriggeredEvent)
		if err != nil {
			msg := "Could not delete .triggered event with ID " + event.Triggeredid + ": " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}

		// get the taskSequence related to the triggeredID and proceed with the next task
		sc.logger.Info("Retrieving task sequence related to triggeredID " + event.Triggeredid)
		eventToSequence, err := sc.taskSequenceRepo.GetTaskSequence(eventScope.Project, event.Triggeredid)
		if err != nil {
			msg := "Could not retrieve task sequence associated to eventID " + event.Triggeredid + ": " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}

		if eventToSequence == nil {
			sc.logger.Info("No task event associated with eventID " + event.Triggeredid + " found.")
			return nil
		}
		sc.logger.Info("Task sequence related to eventID " + event.Triggeredid + ": " + eventToSequence.Stage + "." + eventToSequence.TaskSequenceName)
		sc.logger.Info("Trying to fetch shipyard and get next task")
		shipyard, err := common.GetShipyard(eventScope.Project)
		if err != nil {
			msg := "Could not retrieve shipyard of project " + eventScope.Project + ": " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}
		sequence, err := sc.getTaskSequenceInStage(eventToSequence.Stage, eventToSequence.TaskSequenceName, shipyard)
		if err != nil {
			msg := "No task eventToSequence " + eventToSequence.Stage + "." + eventToSequence.TaskSequenceName + " found in shipyard: " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}

		task, _, err := eventutils.ParseTaskEventType(*event.Type)
		if err != nil {
			return err
		}
		sc.logger.Info("retrieving all .finished events for task " + task + " triggered by " + event.Triggeredid + " to aggregate data")
		allFinishedEventsForTask, err := sc.eventRepo.GetEvents(eventScope.Project, common.EventFilter{
			Type:    "",
			Stage:   &eventScope.Stage,
			Service: &eventScope.Service,
			// TriggeredID: &event.Triggeredid,
			KeptnContext: &event.Shkeptncontext,
		}, common.FinishedEvent)
		if err != nil {
			msg := "Could not retrieve " + *event.Type + " events: " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}

		sc.logger.Info(fmt.Sprintf("Found %d events. Aggregating their properties for next task ", len(allFinishedEventsForTask)))

		finishedEventsData := []interface{}{}

		for index := range allFinishedEventsForTask {
			marshal, _ := json.Marshal(allFinishedEventsForTask[index].Data)
			var tmp interface{}
			_ = json.Unmarshal(marshal, &tmp)
			finishedEventsData = append(finishedEventsData, tmp)
		}

		return sc.proceedTaskSequence(eventScope, sequence, shipyard, finishedEventsData, task)
	}
	return nil
}

func (sc *shipyardController) GetAllTriggeredEvents(filter common.EventFilter) ([]models.Event, error) {
	projects, err := sc.projectRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []models.Event{}
	for _, project := range projects {
		sc.logger.Info(fmt.Sprintf("Retrieving all .triggered events of project %s with filter: %s", project.ProjectName, printObject(filter)))
		events, err := sc.eventRepo.GetEvents(project.ProjectName, filter, common.TriggeredEvent)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (sc *shipyardController) GetTriggeredEventsOfProject(project string, filter common.EventFilter) ([]models.Event, error) {
	sc.logger.Info(fmt.Sprintf("Retrieving all .triggered events with filter: %s", printObject(filter)))
	return sc.eventRepo.GetEvents(project, filter, common.TriggeredEvent)
}

func (sc *shipyardController) getEvents(project string, filter common.EventFilter, status common.EventStatus, nrRetries int) ([]models.Event, error) {
	sc.logger.Info(string("Trying to get " + status + " events"))
	for i := 0; i <= nrRetries; i++ {
		startedEvents, err := sc.eventRepo.GetEvents(project, filter, status)
		if err != nil && err == db.ErrNoEventFound {
			sc.logger.Info(string("No matching " + status + " events found. Retrying in 2s."))
			<-time.After(2 * time.Second)
		} else {
			return startedEvents, err
		}
	}
	return nil, nil
}

func (sc *shipyardController) proceedTaskSequence(eventScope *models.EventScope, taskSequence *keptnv2.Sequence, shipyard *keptnv2.Shipyard, eventHistory []interface{}, previousTask string) error {
	// get the input for the .triggered event that triggered the previous sequence and append it to the list of previous events to gather all required data for the next stage
	inputEvent, eventHistory, err := sc.appendTriggerEventProperties(eventScope, taskSequence, eventHistory)
	if err != nil {
		return err
	}

	task := sc.getNextTaskOfSequence(taskSequence, previousTask, eventScope, eventHistory)
	if task == nil {
		// task sequence completed -> send .finished event and check if a new task sequence should be triggered by the completion
		err = sc.completeTaskSequence(eventScope, taskSequence.Name, inputEvent.ID)
		if err != nil {
			sc.logger.Error("Could not complete task sequence " + eventScope.Stage + "." + taskSequence.Name + " with KeptnContext " + eventScope.KeptnContext)
			return err
		}
		return sc.triggerNextTaskSequences(eventScope, taskSequence, shipyard, eventHistory, inputEvent)
	}
	return sc.sendTaskTriggeredEvent(eventScope, taskSequence.Name, *task, eventHistory)
}

// this function retrieves the .triggered event for the task sequence and appends its properties to the existing .finished events
// this ensures that all parameters set in the .triggered event are received by all execution plane services, instead of just the first one
func (sc *shipyardController) appendTriggerEventProperties(eventScope *models.EventScope, taskSequence *keptnv2.Sequence, eventHistory []interface{}) (*models.Event, []interface{}, error) {
	events, err := sc.eventRepo.GetEvents(eventScope.Project, common.EventFilter{
		Type:         keptnv2.GetTriggeredEventType(eventScope.Stage + "." + taskSequence.Name),
		Stage:        &eventScope.Stage,
		KeptnContext: &eventScope.KeptnContext,
	}, common.TriggeredEvent)

	if err != nil {
		sc.logger.Error("Could not load event that triggered task sequence " + eventScope.Stage + "." + taskSequence.Name + " with KeptnContext " + eventScope.KeptnContext)
		return nil, nil, err
	}

	var inputEvent *models.Event
	if len(events) > 0 {
		inputEvent = &events[0]
		marshal, err := json.Marshal(inputEvent.Data)
		if err != nil {
			sc.logger.Error("Could not marshal input event: " + err.Error())
			return nil, nil, err
		}
		var tmp interface{}
		_ = json.Unmarshal(marshal, &tmp)
		eventHistory = append(eventHistory, tmp)
	}
	return inputEvent, eventHistory, nil
}

func (sc *shipyardController) triggerNextTaskSequences(eventScope *models.EventScope, completedSequence *keptnv2.Sequence, shipyard *keptnv2.Shipyard, eventHistory []interface{}, inputEvent *models.Event) error {
	nextSequences := getTaskSequencesByTrigger(eventScope, completedSequence.Name, shipyard)

	for _, sequence := range nextSequences {
		newScope := &models.EventScope{
			EventData: keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   sequence.StageName,
				Service: eventScope.Service,
			},
			KeptnContext: eventScope.KeptnContext,
		}

		err := sc.sendTaskSequenceTriggeredEvent(newScope, sequence.Sequence.Name, inputEvent)
		if err != nil {
			sc.logger.Error("could not send event " + newScope.Stage + "." + sequence.Sequence.Name + ".triggered: " + err.Error())
			continue
		}

		err = sc.proceedTaskSequence(newScope, &sequence.Sequence, shipyard, eventHistory, "")
		if err != nil {
			sc.logger.Error("could not proceed task sequence " + newScope.Stage + "." + sequence.Sequence.Name + ".triggered: " + err.Error())
			continue
		}
	}

	return nil
}

func (sc *shipyardController) completeTaskSequence(eventScope *models.EventScope, taskSequenceName, triggeredID string) error {
	err := sc.taskSequenceRepo.DeleteTaskSequenceMapping(eventScope.KeptnContext, eventScope.Project, eventScope.Stage, taskSequenceName)
	if err != nil {
		return err
	}

	sc.logger.Info("Deleting all task.finished events of task sequence " + taskSequenceName + " with context " + eventScope.KeptnContext)
	// delete all finished events of this sequence
	finishedEvents, err := sc.eventRepo.GetEvents(eventScope.Project, common.EventFilter{
		Stage:        &eventScope.Stage,
		KeptnContext: &eventScope.KeptnContext,
	}, common.FinishedEvent)

	if err != nil {
		sc.logger.Error("could not retrieve task.finished events: " + err.Error())
		return err
	}

	for _, event := range finishedEvents {
		err = sc.eventRepo.DeleteEvent(eventScope.Project, event.ID, common.FinishedEvent)
		if err != nil {
			sc.logger.Error("could not delete " + *event.Type + " event with ID " + event.ID + ": " + err.Error())
			return err
		}
	}

	return sc.sendTaskSequenceFinishedEvent(eventScope, taskSequenceName, triggeredID)
}

func getTaskSequencesByTrigger(eventScope *models.EventScope, completedTaskSequence string, shipyard *keptnv2.Shipyard) []NextTaskSequence {
	var result []NextTaskSequence
	for _, stage := range shipyard.Spec.Stages {
		for tsIndex, taskSequence := range stage.Sequences {
			for _, trigger := range taskSequence.TriggeredOn {
				if trigger.Event == eventScope.Stage+"."+completedTaskSequence+".finished" {
					appendSequence := false
					// default behavior if no selector is available: 'pass', as well as 'warning' results trigger this sequence
					if trigger.Selector.Match == nil {
						if eventScope.Result == keptnv2.ResultPass || eventScope.Result == keptnv2.ResultWarning {
							appendSequence = true
						}
					} else {
						// if a selector is there, compare the 'result' property
						if string(eventScope.Result) == trigger.Selector.Match["result"] {
							appendSequence = true
						}
					}
					if appendSequence {
						result = append(result, NextTaskSequence{
							Sequence:  stage.Sequences[tsIndex],
							StageName: stage.Name,
						})
					}
				}
			}
		}
	}
	return result
}

func (sc *shipyardController) getTaskSequenceInStage(stageName, taskSequenceName string, shipyard *keptnv2.Shipyard) (*keptnv2.Sequence, error) {
	for _, stage := range shipyard.Spec.Stages {
		if stage.Name == stageName {
			for _, taskSequence := range stage.Sequences {
				if taskSequence.Name == taskSequenceName {
					sc.logger.Info("Found matching task sequence " + taskSequence.Name + " in stage " + stage.Name)
					return &taskSequence, nil
				}
			}
			// provide built-int task sequence for evaluation
			if taskSequenceName == keptnv2.EvaluationTaskName {
				return &keptnv2.Sequence{
					Name:        "evaluation",
					TriggeredOn: nil,
					Tasks: []keptnv2.Task{
						{
							Name: keptnv2.EvaluationTaskName,
						},
					},
				}, nil
			}
			return nil, fmt.Errorf("no task sequence with name %s found in stage %s", taskSequenceName, stageName)
		}
	}
	return nil, fmt.Errorf("no stage with name %s", stageName)
}

func (sc *shipyardController) getNextTaskOfSequence(taskSequence *keptnv2.Sequence, previousTask string, eventScope *models.EventScope, eventHistory []interface{}) *keptnv2.Task {
	if previousTask != "" {
		for _, e := range eventHistory {
			eventData := keptnv2.EventData{}
			_ = keptnv2.Decode(e, &eventData)

			if eventData.Status == keptnv2.StatusErrored || eventData.Result == keptnv2.ResultFailed {
				eventScope.Status = eventData.Status
				eventScope.Result = eventData.Result
				return nil
			}
		}
	}

	if len(taskSequence.Tasks) == 0 {
		sc.logger.Info("Task sequence " + taskSequence.Name + " does not contain any tasks.")
		return nil
	}
	if previousTask == "" {
		sc.logger.Info("Returning first task of task sequence " + taskSequence.Name)
		return &taskSequence.Tasks[0]
	}
	sc.logger.Info("Getting task that should be executed after task " + previousTask)
	for index := range taskSequence.Tasks {
		if taskSequence.Tasks[index].Name == previousTask {
			if len(taskSequence.Tasks) > index+1 {
				sc.logger.Info("found next task: " + taskSequence.Tasks[index+1].Name)
				return &taskSequence.Tasks[index+1]
			}
			break
		}
	}
	sc.logger.Info("No further tasks detected")
	return nil
}

func (sc *shipyardController) sendTaskSequenceTriggeredEvent(eventScope *models.EventScope, taskSequenceName string, inputEvent *models.Event) error {
	eventPayload := map[string]interface{}{}

	eventPayload["project"] = eventScope.Project
	eventPayload["stage"] = eventScope.Stage
	eventPayload["service"] = eventScope.Service

	var mergedPayload interface{}
	if inputEvent != nil {
		marshal, err := json.Marshal(inputEvent.Data)
		if err != nil {
			return fmt.Errorf("could not marshal input event: %s ", err.Error())
		}
		tmp := map[string]interface{}{}
		if err := json.Unmarshal(marshal, &tmp); err != nil {
			return fmt.Errorf("could not convert input event: %s ", err.Error())
		}
		mergedPayload = common.Merge(eventPayload, tmp)
	}

	source, _ := url.Parse("shipyard-controller")
	eventType := eventScope.Stage + "." + taskSequenceName

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetType(keptnv2.GetTriggeredEventType(eventType))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", eventScope.KeptnContext)
	if mergedPayload != nil {
		event.SetData(cloudevents.ApplicationJSON, mergedPayload)
	} else {
		event.SetData(cloudevents.ApplicationJSON, eventPayload)
	}

	toEvent, err := models.ConvertToEvent(event)
	if err != nil {
		return fmt.Errorf("could not store event that triggered task sequence: " + err.Error())
	}
	if err := sc.eventRepo.InsertEvent(eventScope.Project, *toEvent, common.TriggeredEvent); err != nil {
		return fmt.Errorf("could not store event that triggered task sequence: " + err.Error())
	}

	return common.SendEvent(event)
}

func (sc *shipyardController) sendTaskSequenceFinishedEvent(eventScope *models.EventScope, taskSequenceName, triggeredID string) error {
	source, _ := url.Parse("shipyard-controller")
	eventType := eventScope.Stage + "." + taskSequenceName

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetFinishedEventType(eventType))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", eventScope.KeptnContext)
	event.SetExtension("triggeredid", triggeredID)
	event.SetData(cloudevents.ApplicationJSON, eventScope.EventData)

	return common.SendEvent(event)
}

func (sc *shipyardController) sendTaskTriggeredEvent(eventScope *models.EventScope, taskSequenceName string, task keptnv2.Task, eventHistory []interface{}) error {
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

	source, _ := url.Parse("shipyard-controller")

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType(keptnv2.GetTriggeredEventType(task.Name))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", eventScope.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, eventPayload)

	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}
	storeEvent := &models.Event{}
	err = json.Unmarshal(marshal, &storeEvent)
	if err != nil {
		sc.logger.Error("could not transform CloudEvent for storage in mongodb: " + err.Error())
		return err
	}

	err = sc.eventRepo.InsertEvent(eventScope.Project, *storeEvent, common.TriggeredEvent)
	if err != nil {
		sc.logger.Error("Could not store event: " + err.Error())
		return err
	}

	err = sc.taskSequenceRepo.CreateTaskSequenceMapping(eventScope.Project, models.TaskSequenceEvent{
		TaskSequenceName: taskSequenceName,
		TriggeredEventID: event.ID(),
		Stage:            eventScope.Stage,
		KeptnContext:     eventScope.KeptnContext,
	})
	if err != nil {
		sc.logger.Error("Could not store mapping between eventID and task: " + err.Error())
		return err
	}

	return common.SendEvent(event)
}

func printObject(obj interface{}) string {
	indent, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(indent)
}
