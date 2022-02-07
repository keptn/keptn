package v2

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type TaskSequence struct {
	ID       string             `json:"_id" bson:"_id"`
	Sequence keptnv2.Sequence   `json:"sequence" bson:"sequence"`
	Status   TaskSequenceStatus `json:"status" bson:"status"`
	Scope    EventScope         `json:"scope" bson:"scope"`
}

type EventScope struct {
	KeptnContext string `json:"keptnContext" bson:"keptnContext"`
	GitCommitID  string `json:"gitCommitID" bson:"gitCommitID"`
	Project      string `json:"project" bson:"project"`
	Stage        string `json:"stage" bson:"stage"`
	Service      string `json:"service" bson:"service"`
}

type TaskSequenceStatus struct {
	State         string                `json:"state" bson:"state"` // triggered, waiting, suspended (approval in progress), paused, finished, cancelled, timedOut
	PreviousTasks []TaskExecutionResult `json:"previousTasks" bson:"previousTasks"`
	CurrentTask   TaskExecution         `json:"currentTask" bson:"currentTask"`
}

type TaskExecutionResult struct {
	Name        string                 `json:"name" bson:"name"`
	TriggeredID string                 `json:"triggeredID" bson:"triggeredID"`
	Result      string                 `json:"result" bson:"result"`
	Status      string                 `json:"status" bson:"status"`
	Properties  map[string]interface{} `json:"properties" bson:"properties"`
}

type TaskExecution struct {
	Name        string      `json:"name" bson:"name"`
	TriggeredID string      `json:"triggeredID" bson:"triggeredID"`
	Events      []TaskEvent `json:"events" bson:"events"`
}

func (e TaskSequence) GetNextTaskOfSequence() *keptnv2.Task {
	if e.Status.CurrentTask.IsFailed() {
		return nil
	}
	nextTaskIndex := 0
	if e.Status.PreviousTasks != nil && len(e.Status.PreviousTasks) > 0 {
		nextTaskIndex = len(e.Status.PreviousTasks)
	}

	if len(e.Sequence.Tasks) > nextTaskIndex {
		return &e.Sequence.Tasks[nextTaskIndex]
	}
	return nil
}

func (e TaskSequence) GetLastTaskExecutionResult() *TaskExecutionResult {
	if len(e.Status.PreviousTasks) == 0 {
		return nil
	}
	return &e.Status.PreviousTasks[len(e.Status.PreviousTasks)-1]
}

// IsFinished indicates if a task is finished, i.e. the number of task.started and task.finished events line up
func (e TaskExecution) IsFinished() bool {
	if len(e.Events) == 0 {
		return false
	}
	nrStartedEvents := 0
	nrFinishedEvents := 0
	for _, event := range e.Events {
		if keptnv2.IsStartedEventType(event.EventType) {
			nrStartedEvents++
		} else if keptnv2.IsFinishedEventType(event.EventType) {
			nrFinishedEvents++
		}
	}

	if nrFinishedEvents == nrStartedEvents && nrFinishedEvents > 0 {
		return true
	}
	return false
}

func (e TaskExecution) IsFailed() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == string(keptnv2.ResultFailed) {
				return true
			}
		}
	}
	return false
}

func (e TaskExecution) IsErrored() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Status == string(keptnv2.StatusErrored) {
				return true
			}
		}
	}
	return false
}

type TaskEvent struct {
	EventType  string                 `json:"eventType" bson:"eventType"`
	Source     string                 `json:"source" bson:"source"`
	Result     string                 `json:"result" bson:"result"`
	Status     string                 `json:"status" bson:"status"`
	Time       string                 `json:"time" bson:"time"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

type GetTaskSequenceFilter struct {
	Scope              EventScope
	Status             []string
	Name               string
	CurrentTriggeredID string
}
