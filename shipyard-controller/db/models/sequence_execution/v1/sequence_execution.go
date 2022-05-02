package v1

import (
	"encoding/json"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
)

const SchemaVersionV1 = "1"

type SchemaVersion struct {
	SchemaVersion string `json:"schemaVersion" bson:"schemaVersion"`
}
type JsonStringEncodedSequenceExecution struct {
	ID string `json:"_id" bson:"_id"`
	// SchemaVersion indicates the version of the schema - needed to decide if items in collection need to be migrated
	SchemaVersion `bson:",inline"`
	// Sequence contains the complete sequence definition
	Sequence Sequence                `json:"sequence" bson:"sequence"`
	Status   SequenceExecutionStatus `json:"status" bson:"status"`
	Scope    models.EventScope       `json:"scope" bson:"scope"`
	// EncodedInputProperties contains properties of the event which triggered the task sequence
	EncodedInputProperties string `json:"encodedInputProperties" bson:"encodedInputProperties"`
}

type Sequence struct {
	Name  string `json:"name" bson:"name"`
	Tasks []Task `json:"tasks" bson:"tasks"`
}

func (s Sequence) DecodeTasks() []keptnv2.Task {
	tasks := []keptnv2.Task{}

	for _, task := range s.Tasks {
		newTask := keptnv2.Task{
			Name:           task.Name,
			TriggeredAfter: task.TriggeredAfter,
		}
		if task.EncodedProperties != "" {
			properties := map[string]interface{}{}
			if err := json.Unmarshal([]byte(task.EncodedProperties), &properties); err == nil {
				newTask.Properties = properties
			}
		}
		tasks = append(tasks, newTask)
	}
	return tasks
}

type Task struct {
	Name              string `json:"name" bson:"name"`
	TriggeredAfter    string `json:"triggeredAfter,omitempty" bson:"triggeredAfter,omitempty"`
	EncodedProperties string `json:"encodedProperties" bson:"encodedProperties"`
}

type SequenceExecutionStatus struct {
	State string `json:"state" bson:"state"` // triggered, waiting, suspended (approval in progress), paused, finished, cancelled, timedOut
	// StateBeforePause is needed to keep track of the state before a sequence has been paused. Example: when a sequence has been paused while being queued, and then resumed, it should not be set to started immediately, but to the state it had before
	StateBeforePause string `json:"stateBeforePause" bson:"stateBeforePause"`
	// PreviousTasks contains the results of all completed tasks of the sequence
	PreviousTasks []TaskExecutionResult `json:"previousTasks" bson:"previousTasks"`
	// CurrentTask represents the state of the currently active task
	CurrentTask TaskExecutionState `json:"currentTask" bson:"currentTask"`
}

func (s SequenceExecutionStatus) DecodePreviousTasks() []models.TaskExecutionResult {
	result := []models.TaskExecutionResult{}

	for _, previousTask := range s.PreviousTasks {
		newPreviousTask := models.TaskExecutionResult{
			Name:        previousTask.Name,
			TriggeredID: previousTask.TriggeredID,
			Result:      previousTask.Result,
			Status:      previousTask.Status,
		}

		if previousTask.EncodedProperties != "" {
			properties := map[string]interface{}{}
			if err := json.Unmarshal([]byte(previousTask.EncodedProperties), &properties); err == nil {
				newPreviousTask.Properties = properties
			}
		}

		result = append(result, newPreviousTask)
	}
	return result
}

type TaskExecutionResult struct {
	Name        string             `json:"name" bson:"name"`
	TriggeredID string             `json:"triggeredID" bson:"triggeredID"`
	Result      keptnv2.ResultType `json:"result" bson:"result"`
	Status      keptnv2.StatusType `json:"status" bson:"status"`
	// EncodedProperties contains the aggregated results of the task's executors
	EncodedProperties string `json:"encodedProperties" bson:"encodedProperties"`
}

type TaskExecutionState struct {
	Name        string      `json:"name" bson:"name"`
	TriggeredID string      `json:"triggeredID" bson:"triggeredID"`
	Events      []TaskEvent `json:"events" bson:"events"`
}

func (s TaskExecutionState) DecodeEvents() []models.TaskEvent {
	result := []models.TaskEvent{}

	for _, event := range s.Events {
		newEvent := models.TaskEvent{
			EventType: event.EventType,
			Source:    event.Source,
			Result:    event.Result,
			Status:    event.Status,
			Time:      event.Time,
		}
		if event.EncodedProperties != "" {
			properties := map[string]interface{}{}
			if err := json.Unmarshal([]byte(event.EncodedProperties), &properties); err == nil {
				newEvent.Properties = properties
			}
		}
		result = append(result, newEvent)
	}
	return result
}

type TaskEvent struct {
	EventType         string             `json:"eventType" bson:"eventType"`
	Source            string             `json:"source" bson:"source"`
	Result            keptnv2.ResultType `json:"result" bson:"result"`
	Status            keptnv2.StatusType `json:"status" bson:"status"`
	Time              string             `json:"time" bson:"time"`
	EncodedProperties string             `json:"encodedProperties" bson:"encodedProperties"`
}

func (e JsonStringEncodedSequenceExecution) ToSequenceExecution() models.SequenceExecution {
	result := models.SequenceExecution{
		ID: e.ID,
		Sequence: keptnv2.Sequence{
			Name:  e.Sequence.Name,
			Tasks: e.Sequence.DecodeTasks(),
		},
		Status: models.SequenceExecutionStatus{
			State:            e.Status.State,
			StateBeforePause: e.Status.StateBeforePause,
			PreviousTasks:    e.Status.DecodePreviousTasks(),
			CurrentTask: models.TaskExecutionState{
				Name:        e.Status.CurrentTask.Name,
				TriggeredID: e.Status.CurrentTask.TriggeredID,
				Events:      e.Status.CurrentTask.DecodeEvents(),
			},
		},
		Scope: e.Scope,
	}
	inputProperties := map[string]interface{}{}
	err := json.Unmarshal([]byte(e.EncodedInputProperties), &inputProperties)
	if err == nil {
		result.InputProperties = inputProperties
	}
	return result
}
