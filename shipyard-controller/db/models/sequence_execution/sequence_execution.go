package sequence_execution

import (
	"encoding/json"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type SequenceExecution struct {
	ID string `json:"_id" bson:"_id"`
	// SchemaVersion indicates the version of the schema - needed to decide if items in collection need to be migrated
	SchemaVersion string `json:"schemaVersion" bson:"schemaVersion"`
	// Sequence contains the complete sequence definition
	Sequence Sequence                `json:"sequence" bson:"sequence"`
	Status   SequenceExecutionStatus `json:"status" bson:"status"`
	Scope    models.EventScope       `json:"scope" bson:"scope"`
	// InputProperties contains properties of the event which triggered the task sequence
	InputProperties string `json:"inputProperties" bson:"inputProperties"`
}

func (e SequenceExecution) ToSequenceExecution() (*models.SequenceExecution, error) {
	// TODO
	return nil, nil
}

type Sequence struct {
	Name  string `json:"name" bson:"name"`
	Tasks []Task `json:"tasks" bson:"tasks"`
}

type Task struct {
	Name           string `json:"name" bson:"name"`
	TriggeredAfter string `json:"triggeredAfter,omitempty" bson:"triggeredAfter,omitempty"`
	Properties     string `json:"properties" bson:"properties"`
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

type TaskExecutionResult struct {
	Name        string             `json:"name" bson:"name"`
	TriggeredID string             `json:"triggeredID" bson:"triggeredID"`
	Result      keptnv2.ResultType `json:"result" bson:"result"`
	Status      keptnv2.StatusType `json:"status" bson:"status"`
	// Properties contains the aggregated results of the task's executors
	Properties string `json:"properties" bson:"properties"`
}

type TaskExecutionState struct {
	Name        string      `json:"name" bson:"name"`
	TriggeredID string      `json:"triggeredID" bson:"triggeredID"`
	Events      []TaskEvent `json:"events" bson:"events"`
}

type TaskEvent struct {
	EventType  string             `json:"eventType" bson:"eventType"`
	Source     string             `json:"source" bson:"source"`
	Result     keptnv2.ResultType `json:"result" bson:"result"`
	Status     keptnv2.StatusType `json:"status" bson:"status"`
	Time       string             `json:"time" bson:"time"`
	Properties string             `json:"properties" bson:"properties"`
}

func FromSequenceExecution(se models.SequenceExecution) SequenceExecution {
	newSE := SequenceExecution{
		ID: se.ID,
		Sequence: Sequence{
			Name:  se.Sequence.Name,
			Tasks: transformTasks(se.Sequence.Tasks),
		},
		Status: transformStatus(se.Status),
		Scope:  se.Scope,
	}
	if se.InputProperties != nil {
		inputPropertiesJsonString, err := json.Marshal(se.InputProperties)
		if err == nil {
			newSE.InputProperties = string(inputPropertiesJsonString)
		}
	}
	return newSE
}

func transformTasks(tasks []keptnv2.Task) []Task {
	result := []Task{}

	for _, task := range tasks {
		newTask := Task{
			Name:           task.Name,
			TriggeredAfter: task.TriggeredAfter,
		}
		if task.Properties != nil {
			taskPropertiesString, err := json.Marshal(task.Properties)
			if err == nil {
				newTask.Properties = string(taskPropertiesString)
			}
		}
		result = append(result, newTask)
	}
	return result
}

func transformStatus(status models.SequenceExecutionStatus) SequenceExecutionStatus {
	newStatus := SequenceExecutionStatus{
		State:            status.State,
		StateBeforePause: status.StateBeforePause,
		PreviousTasks:    transformPreviousTasks(status.PreviousTasks),
		CurrentTask:      transformCurrentTask(status.CurrentTask),
	}

	return newStatus
}

func transformCurrentTask(task models.TaskExecutionState) TaskExecutionState {
	newTaskExecutionState := TaskExecutionState{
		Name:        task.Name,
		TriggeredID: task.TriggeredID,
		Events:      transformTaskEvents(task.Events),
	}
	return newTaskExecutionState
}

func transformTaskEvents(events []models.TaskEvent) []TaskEvent {
	newTaskEvents := []TaskEvent{}

	for _, e := range events {
		newTaskEvent := TaskEvent{
			EventType: e.EventType,
			Source:    e.Source,
			Result:    e.Result,
			Status:    e.Status,
			Time:      e.Time,
		}

		if e.Properties != nil {
			properties, err := json.Marshal(e.Properties)
			if err == nil {
				newTaskEvent.Properties = string(properties)
			}
		}
		newTaskEvents = append(newTaskEvents, newTaskEvent)
	}
	return newTaskEvents
}

func transformPreviousTasks(tasks []models.TaskExecutionResult) []TaskExecutionResult {
	newPreviousTasks := []TaskExecutionResult{}

	for _, t := range tasks {
		newPreviousTask := TaskExecutionResult{
			Name:        t.Name,
			TriggeredID: t.TriggeredID,
			Result:      t.Result,
			Status:      t.Status,
		}

		if t.Properties != nil {
			properties, err := json.Marshal(t.Properties)
			if err == nil {
				newPreviousTask.Properties = string(properties)
			}
		}
		newPreviousTasks = append(newPreviousTasks, newPreviousTask)
	}
	return newPreviousTasks
}
