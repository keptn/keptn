package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"net/url"
	"strings"
	"time"
)

const maxRepoReadRetries = 30

var errNoMatchingEvent = errors.New("no matching event found")

var shipyardControllerInstance *shipyardController

type IShipyardController interface {
	GetAllTriggeredEvents(filter db.EventFilter) ([]models.Event, error)
	GetTriggeredEventsOfProject(project string, filter db.EventFilter) ([]models.Event, error)
	HandleIncomingEvent(event models.Event) error
}

type shipyardController struct {
	projectRepo      db.ProjectRepo
	eventRepo        db.EventRepo
	taskSequenceRepo db.TaskSequenceRepo
	logger           *keptncommon.Logger
}

func GetShipyardControllerInstance() *shipyardController {
	if shipyardControllerInstance == nil {
		logger := keptncommon.NewLogger("", "", "shipyard-controller")
		shipyardControllerInstance = &shipyardController{
			projectRepo: &db.ProjectMongoDBRepo{
				Logger: logger,
			},
			eventRepo: &db.MongoDBEventsRepo{
				Logger: logger,
			},
			taskSequenceRepo: &db.TaskSequenceMongoDBRepo{
				Logger: logger,
			},
			logger: logger,
		}
	}
	return shipyardControllerInstance
}

func (sc *shipyardController) GetAllTriggeredEvents(filter db.EventFilter) ([]models.Event, error) {
	projects, err := sc.projectRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []models.Event{}
	for _, project := range projects {
		sc.logger.Info(fmt.Sprintf("Retrieving all .triggered events of project %s with filter: %s", project, printObject(filter)))
		events, err := sc.eventRepo.GetEvents(project, filter, db.TriggeredEvent)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (sc *shipyardController) GetTriggeredEventsOfProject(project string, filter db.EventFilter) ([]models.Event, error) {
	sc.logger.Info(fmt.Sprintf("Retrieving all .triggered events with filter: %s", printObject(filter)))
	return sc.eventRepo.GetEvents(project, filter, db.TriggeredEvent)
}

func (sc *shipyardController) HandleIncomingEvent(event models.Event) error {
	// check if the status type is either 'triggered', 'started', or 'finished'
	split := strings.Split(*event.Type, ".")

	statusType := split[len(split)-1]

	switch statusType {
	case string(db.TriggeredEvent):
		return sc.handleTriggeredEvent(event)
	case string(db.StartedEvent):
		return sc.handleStartedEvent(event)
	case string(db.FinishedEvent):
		return sc.handleFinishedEvent(event)
	default:
		return nil
	}
}

func getEventScope(event models.Event) (*keptnv2.EventData, error) {
	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}
	data := &keptnv2.EventData{}
	err = json.Unmarshal(marshal, data)
	if err != nil {
		return nil, err
	}
	if data.Project == "" {
		return nil, errors.New("event does not contain a project")
	}
	if data.Stage == "" {
		return nil, errors.New("event does not contain a stage")
	}
	if data.Service == "" {
		return nil, errors.New("event does not contain a service")
	}
	return data, nil
}

func (sc *shipyardController) handleFinishedEvent(event models.Event) error {

	if *event.Source == "shipyard-controller" {
		sc.logger.Info("Received event from myself. Ignoring.")
		return nil
	}

	eventScope, err := getEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, printObject(event)))

	trimmedEventType := strings.TrimSuffix(*event.Type, string(db.FinishedEvent))
	// get corresponding 'started' event for the incoming 'finished' event
	filter := db.EventFilter{
		Type:        trimmedEventType + string(db.StartedEvent),
		TriggeredID: &event.Triggeredid,
	}
	startedEvents, err := sc.getEvents(eventScope.Project, filter, db.StartedEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		sc.logger.Error(msg)
		return errors.New(msg)
	} else if startedEvents == nil || len(startedEvents) == 0 {
		msg := "no matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		sc.logger.Error(msg)
		return errNoMatchingEvent
	}

	// persist the .finished event
	err = sc.eventRepo.InsertEvent(eventScope.Project, event, db.FinishedEvent)
	if err != nil {
		sc.logger.Error("Could not store .finished event: " + err.Error())
	}

	for _, startedEvent := range startedEvents {
		if *event.Source == *startedEvent.Source {
			err = sc.eventRepo.DeleteEvent(eventScope.Project, startedEvent.ID, db.StartedEvent)
			if err != nil {
				msg := "could not delete '.started' event with ID " + startedEvent.ID + ": " + err.Error()
				sc.logger.Error(msg)
				return errors.New(msg)
			}
		}
	}
	// check if this was the last '.started' event
	if len(startedEvents) == 1 {
		triggeredEventFilter := db.EventFilter{
			Type: trimmedEventType + string(db.TriggeredEvent),
			ID:   &event.Triggeredid,
		}
		triggeredEvents, err := sc.getEvents(eventScope.Project, triggeredEventFilter, db.TriggeredEvent, maxRepoReadRetries)
		if err != nil {
			msg := "could not retrieve '.triggered' event with ID " + event.Triggeredid + ": " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}
		if triggeredEvents == nil || len(triggeredEvents) == 0 {
			msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
			sc.logger.Error(msg)
			return errNoMatchingEvent
		}
		// if the previously deleted '.started' event was the last, the '.triggered' event can be removed
		sc.logger.Info("triggered event will be deleted")
		err = sc.eventRepo.DeleteEvent(eventScope.Project, triggeredEvents[0].ID, db.TriggeredEvent)
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
		shipyard, err := common.GetShipyard(eventScope)
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

		sc.logger.Info("retrieving all .finished events for task " + trimmedEventType + " triggered by " + event.Triggeredid + " to aggregate data")
		allFinishedEventsForTask, err := sc.eventRepo.GetEvents(eventScope.Project, db.EventFilter{
			Type:    "",
			Stage:   &eventScope.Stage,
			Service: &eventScope.Service,
			// TriggeredID: &event.Triggeredid,
			KeptnContext: &event.Shkeptncontext,
		}, db.FinishedEvent)
		if err != nil {
			msg := "Could not retrieve " + *event.Type + " events: " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}

		sc.logger.Info(fmt.Sprintf("Found %d events. Aggregating their properties for next task ", len(allFinishedEventsForTask)))

		finishedEventsData := []interface{}{}

		continueTaskSequence := true
		for index, finishedEvent := range allFinishedEventsForTask {
			finishedEventScope, err := getEventScope(finishedEvent)
			if err != nil {
				sc.logger.Error("Could not determine scope of .finished event with ID " + finishedEvent.ID + ": " + err.Error())
				continueTaskSequence = false
			} else if finishedEventScope.Status == keptnv2.StatusErrored {
				sc.logger.Info("Finished event with ID " + finishedEvent.ID + " reported an error. Will abort task sequence " + sequence.Name + " with KeptnContext " + event.Shkeptncontext)
				continueTaskSequence = false
			}
			marshal, _ := json.Marshal(allFinishedEventsForTask[index].Data)
			var tmp interface{}
			_ = json.Unmarshal(marshal, &tmp)
			finishedEventsData = append(finishedEventsData, tmp)
		}

		if !continueTaskSequence {
			sc.logger.Info("Aborting task sequence " + sequence.Name + " with KeptnContext " + event.Shkeptncontext)
			return sc.completeTaskSequence(event.Shkeptncontext, &keptnv2.EventData{
				Project: eventScope.Project,
				Stage:   eventScope.Stage,
				Service: eventScope.Service,
				Labels:  eventScope.Labels,
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultFailed,
				Message: "",
			}, sequence.Name)
		}

		split := strings.Split(trimmedEventType, ".")

		if len(split) < 2 {
			msg := "Could not determine task name "
			sc.logger.Error(msg)
			return errors.New(msg)
		}
		return sc.proceedTaskSequence(eventScope, sequence, event, shipyard, finishedEventsData, split[len(split)-2])
	}
	return nil
}

func printObject(obj interface{}) string {
	indent, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(indent)
}

func merge(in1, in2 interface{}) interface{} {
	switch in1 := in1.(type) {
	case []interface{}:
		in2, ok := in2.([]interface{})
		if !ok {
			return in1
		}
		return append(in1, in2...)
	case map[string]interface{}:
		in2, ok := in2.(map[string]interface{})
		if !ok {
			return in1
		}
		for k, v2 := range in2 {
			if v1, ok := in1[k]; ok {
				in1[k] = merge(v1, v2)
			} else {
				in1[k] = v2
			}
		}
	case nil:
		in2, ok := in2.(map[string]interface{})
		if ok {
			return in2
		}
	}
	return in1
}

func (sc *shipyardController) getEvents(project string, filter db.EventFilter, status db.EventStatus, nrRetries int) ([]models.Event, error) {
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

func (sc *shipyardController) handleStartedEvent(event models.Event) error {

	sc.logger.Info("Received .started event: " + *event.Type)
	eventScope, err := getEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, printObject(event)))

	trimmedEventType := strings.TrimSuffix(*event.Type, string(db.StartedEvent))
	// get corresponding 'triggered' event for the incoming 'started' event
	filter := db.EventFilter{
		Type: trimmedEventType + string(db.TriggeredEvent),
		ID:   &event.Triggeredid,
	}

	events, err := sc.getEvents(eventScope.Project, filter, db.TriggeredEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		sc.logger.Error(msg)
		return errors.New(msg)
	} else if events == nil || len(events) == 0 {
		msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		sc.logger.Error(msg)
		return errNoMatchingEvent
	}

	return sc.eventRepo.InsertEvent(eventScope.Project, event, db.StartedEvent)
}

func (sc *shipyardController) handleTriggeredEvent(event models.Event) error {

	if *event.Source == "shipyard-controller" {
		sc.logger.Info("Received event from myself. Ignoring.")
		return nil
	}

	eventScope, err := getEventScope(event)
	if err != nil {
		sc.logger.Error("Could not determine eventScope of event: " + err.Error())
		return err
	}
	sc.logger.Info(fmt.Sprintf("Context of event %s, sent by %s: %s", *event.Type, *event.Source, printObject(event)))

	sc.logger.Info("received event of type " + *event.Type + " from " + *event.Source)
	split := strings.Split(*event.Type, ".")

	sc.logger.Info("Checking if .triggered event should start a sequence in project " + eventScope.Project)
	// get stage and taskSequenceName - cannot tell if this is actually a task sequence triggered event though
	var stageName, taskSequenceName string
	if len(split) >= 3 {
		taskSequenceName = split[len(split)-2]
		stageName = split[len(split)-3]
	}

	shipyard, err := common.GetShipyard(eventScope)
	if err != nil {
		msg := "could not retrieve shipyard: " + err.Error()
		sc.logger.Error(msg)
		return sc.sendTaskSequenceFinishedEvent(event.Shkeptncontext, &keptnv2.EventData{
			Project: eventScope.Project,
			Stage:   eventScope.Stage,
			Service: eventScope.Service,
			Labels:  eventScope.Labels,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: msg,
		}, taskSequenceName)
	}

	// validate the shipyard version - only shipyard files following the '0.2.0' spec are supported by the shipyard controller
	err = common.ValidateShipyardVersion(shipyard)
	if err != nil {
		// if the validation has not been successful: send a <task-sequence>.finished event with status=errored
		sc.logger.Error("invalid shipyard version: " + err.Error())
		return sc.sendTaskSequenceFinishedEvent(event.Shkeptncontext, &keptnv2.EventData{
			Project: eventScope.Project,
			Stage:   eventScope.Stage,
			Service: eventScope.Service,
			Labels:  eventScope.Labels,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: "Found shipyard.yaml with invalid version. Please upgrade the shipyard.yaml of the project using the Keptn CLI: 'keptn upgrade project " + eventScope.Project + " --shipyard'. '",
		}, taskSequenceName)
	}

	taskSequence, err := sc.getTaskSequenceInStage(stageName, taskSequenceName, shipyard)
	if err != nil && err == errNoTaskSequence {
		sc.logger.Info("no task sequence with name " + taskSequenceName + " found in stage " + stageName)
		return err
	} else if err != nil && err == errNoStage {
		sc.logger.Info("no stage with name " + stageName + " found in project " + eventScope.Project)
		return err
	}

	if err := sc.eventRepo.InsertEvent(eventScope.Project, event, db.TriggeredEvent); err != nil {
		sc.logger.Info("could not store event that triggered task sequence: " + err.Error())
	}

	eventScope.Stage = stageName

	eventMap := map[string]interface{}{}

	marshal, err := json.Marshal(event.Data)
	if err != nil {
		sc.logger.Info("could not marshal incoming event: " + err.Error())
		return err
	}
	if err := json.Unmarshal(marshal, &eventMap); err != nil {
		sc.logger.Info("could not convert incoming event to map[string]interface{}: " + err.Error())
		return err
	}

	return sc.proceedTaskSequence(eventScope, taskSequence, event, shipyard, []interface{}{eventMap}, "")
}

func (sc *shipyardController) proceedTaskSequence(eventScope *keptnv2.EventData, taskSequence *keptnv2.Sequence, event models.Event, shipyard *keptnv2.Shipyard, previousFinishedEvents []interface{}, previousTask string) error {
	task, err := sc.getNextTaskOfSequence(taskSequence, previousTask)
	if err != nil && err == errNoFurtherTaskForSequence {
		// get the input for te .triggered event that triggered the previous sequence and append it to the list of previous events to gather all required data for the next stage
		events, err := sc.eventRepo.GetEvents(eventScope.Project, db.EventFilter{
			Type:         keptnv2.GetTriggeredEventType(eventScope.Stage + "." + taskSequence.Name),
			Stage:        &eventScope.Stage,
			KeptnContext: &event.Shkeptncontext,
		}, db.TriggeredEvent)

		if err != nil {
			sc.logger.Error("Could not load event that triggered task sequence " + eventScope.Stage + "." + taskSequence.Name + " with KeptnContext " + event.Shkeptncontext)
			return err
		}

		var inputEvent *models.Event
		if len(events) > 0 {
			inputEvent = &events[0]
			marshal, err := json.Marshal(inputEvent.Data)
			if err != nil {
				sc.logger.Error("Could not marshal input event: " + err.Error())
				return err
			}
			var tmp interface{}
			_ = json.Unmarshal(marshal, &tmp)
			previousFinishedEvents = append(previousFinishedEvents, tmp)
		}

		// task sequence completed -> send .finished event and check if a new task sequence should be triggered by the completion
		err = sc.completeTaskSequence(event.Shkeptncontext, eventScope, taskSequence.Name)
		if err != nil {
			sc.logger.Error("Could not complete task sequence " + eventScope.Stage + "." + taskSequence.Name + " with KeptnContext " + event.Shkeptncontext)
			return err
		}
		if eventScope.Result == keptnv2.ResultPass {
			return sc.triggerNextTaskSequences(event, eventScope, taskSequence, shipyard, previousFinishedEvents, inputEvent)
		}
		return nil
	} else if err != nil {
		sc.logger.Error("Could not get next task of sequence: " + err.Error())
		return err
	}
	return sc.sendTaskTriggeredEvent(event.Shkeptncontext, eventScope, taskSequence.Name, *task, previousFinishedEvents)
}

func (sc *shipyardController) triggerNextTaskSequences(event models.Event, eventScope *keptnv2.EventData, completedSequence *keptnv2.Sequence, shipyard *keptnv2.Shipyard, previousFinishedEvents []interface{}, inputEvent *models.Event) error {

	nextSequences := sc.getTaskSequencesByTrigger(eventScope, completedSequence.Name, shipyard)

	for _, sequence := range nextSequences {
		newScope := &keptnv2.EventData{
			Project: eventScope.Project,
			Stage:   sequence.StageName,
			Service: eventScope.Service,
		}
		err := sc.sendTaskSequenceTriggeredEvent(event.Shkeptncontext, newScope, sequence.Sequence.Name, inputEvent)
		if err != nil {
			sc.logger.Error("could not send event " + newScope.Stage + "." + sequence.Sequence.Name + ".triggered: " + err.Error())
			continue
		}

		err = sc.proceedTaskSequence(newScope, &sequence.Sequence, event, shipyard, previousFinishedEvents, "")
		if err != nil {
			sc.logger.Error("could not proceed task sequence " + newScope.Stage + "." + sequence.Sequence.Name + ".triggered: " + err.Error())
			continue
		}
	}

	return nil
}

func (sc *shipyardController) completeTaskSequence(keptnContext string, eventScope *keptnv2.EventData, taskSequenceName string) error {
	err := sc.taskSequenceRepo.DeleteTaskSequenceMapping(keptnContext, eventScope.Project, eventScope.Stage, taskSequenceName)
	if err != nil {
		return err
	}

	sc.logger.Info("Deleting all task.finished events of task sequence " + taskSequenceName + " with context " + keptnContext)
	// delete all finished events of this sequence
	finishedEvents, err := sc.eventRepo.GetEvents(eventScope.Project, db.EventFilter{
		Stage:        &eventScope.Stage,
		KeptnContext: &keptnContext,
	}, db.FinishedEvent)

	if err != nil {
		sc.logger.Error("could not retrieve task.finished events: " + err.Error())
		return err
	}

	for _, event := range finishedEvents {
		err = sc.eventRepo.DeleteEvent(eventScope.Project, event.ID, db.FinishedEvent)
		if err != nil {
			sc.logger.Error("could not delete " + *event.Type + " event with ID " + event.ID + ": " + err.Error())
			return err
		}
	}

	return sc.sendTaskSequenceFinishedEvent(keptnContext, eventScope, taskSequenceName)
}

var errNoFurtherTaskForSequence = errors.New("no further task for sequence")
var errNoTaskSequence = errors.New("no task sequence found")
var errNoStage = errors.New("no stage found")

func (sc *shipyardController) getTaskSequencesByTrigger(eventScope *keptnv2.EventData, completedTaskSequence string, shipyard *keptnv2.Shipyard) []NextTaskSequence {
	result := []NextTaskSequence{}
	for _, stage := range shipyard.Spec.Stages {
		for tsIndex, taskSequence := range stage.Sequences {
			for _, trigger := range taskSequence.Triggers {
				if trigger == eventScope.Stage+"."+completedTaskSequence+".finished" {
					result = append(result, NextTaskSequence{
						Sequence:  stage.Sequences[tsIndex],
						StageName: stage.Name,
					})
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
					Name:     "evaluation",
					Triggers: nil,
					Tasks: []keptnv2.Task{
						{
							Name: keptnv2.EvaluationTaskName,
						},
					},
				}, nil
			}
			return nil, errNoTaskSequence
		}
	}
	return nil, errNoStage
}

func (sc *shipyardController) getNextTaskOfSequence(taskSequence *keptnv2.Sequence, previousTask string) (*keptnv2.Task, error) {
	if len(taskSequence.Tasks) == 0 {
		sc.logger.Info("Task sequence " + taskSequence.Name + " does not contain any tasks.")
		return nil, errNoFurtherTaskForSequence
	}
	if previousTask == "" {
		sc.logger.Info("Returning first task of task sequence " + taskSequence.Name)
		return &taskSequence.Tasks[0], nil
	}
	sc.logger.Info("Getting task that should be executed after task " + previousTask)
	for index := range taskSequence.Tasks {
		if taskSequence.Tasks[index].Name == previousTask {
			if len(taskSequence.Tasks) > index+1 {
				sc.logger.Info("found next task: " + taskSequence.Tasks[index+1].Name)
				return &taskSequence.Tasks[index+1], nil
			}
			break
		}
	}
	sc.logger.Info("No further tasks detected")
	return nil, errNoFurtherTaskForSequence
}

func (sc *shipyardController) sendTaskSequenceTriggeredEvent(keptnContext string, eventScope *keptnv2.EventData, taskSequenceName string, inputEvent *models.Event) error {

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
		mergedPayload = merge(eventPayload, tmp)
	}

	source, _ := url.Parse("shipyard-controller")
	eventType := eventScope.Stage + "." + taskSequenceName

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetType(keptnv2.GetTriggeredEventType(eventType))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", keptnContext)
	if mergedPayload != nil {
		event.SetData(cloudevents.ApplicationJSON, mergedPayload)
	} else {
		event.SetData(cloudevents.ApplicationJSON, eventPayload)
	}

	toEvent, err := models.ConvertToEvent(event)
	if err != nil {
		return fmt.Errorf("could not store event that triggered task sequence: " + err.Error())
	}
	if err := sc.eventRepo.InsertEvent(eventScope.Project, *toEvent, db.TriggeredEvent); err != nil {
		return fmt.Errorf("could not store event that triggered task sequence: " + err.Error())
	}

	return common.SendEvent(event)
}

func (sc *shipyardController) sendTaskSequenceFinishedEvent(keptnContext string, eventScope *keptnv2.EventData, taskSequenceName string) error {
	source, _ := url.Parse("shipyard-controller")
	eventType := eventScope.Stage + "." + taskSequenceName

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetFinishedEventType(eventType))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", keptnContext)
	event.SetData(cloudevents.ApplicationJSON, eventScope)

	return common.SendEvent(event)
}

func (sc *shipyardController) sendTaskTriggeredEvent(keptnContext string, eventScope *keptnv2.EventData, taskSequenceName string, task keptnv2.Task, previousFinishedEvents []interface{}) error {

	eventPayload := map[string]interface{}{}

	eventPayload["project"] = eventScope.Project
	eventPayload["stage"] = eventScope.Stage
	eventPayload["service"] = eventScope.Service

	eventPayload[task.Name] = task.Properties

	var mergedPayload interface{}
	mergedPayload = nil
	if previousFinishedEvents != nil {
		for index := range previousFinishedEvents {
			if mergedPayload == nil {
				mergedPayload = merge(eventPayload, previousFinishedEvents[index])
			} else {
				mergedPayload = merge(mergedPayload, previousFinishedEvents[index])
			}
		}
	}

	// make sure the result from the previous event is used
	eventPayload["result"] = eventScope.Result
	eventPayload["status"] = eventScope.Status

	source, _ := url.Parse("shipyard-controller")

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType(keptnv2.GetTriggeredEventType(task.Name))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", keptnContext)
	event.SetData(cloudevents.ApplicationJSON, eventPayload)

	marshal, err := json.Marshal(event)

	storeEvent := &models.Event{}
	err = json.Unmarshal(marshal, &storeEvent)
	if err != nil {
		sc.logger.Error("could not transform CloudEvent for storage in mongodb: " + err.Error())
		return err
	}

	err = sc.eventRepo.InsertEvent(eventScope.Project, *storeEvent, db.TriggeredEvent)
	if err != nil {
		sc.logger.Error("Could not store event: " + err.Error())
		return err
	}

	err = sc.taskSequenceRepo.CreateTaskSequenceMapping(eventScope.Project, models.TaskSequenceEvent{
		TaskSequenceName: taskSequenceName,
		TriggeredEventID: event.ID(),
		Stage:            eventScope.Stage,
		KeptnContext:     keptnContext,
	})
	if err != nil {
		sc.logger.Error("Could not store mapping between eventID and task: " + err.Error())
		return err
	}

	return common.SendEvent(event)
}
