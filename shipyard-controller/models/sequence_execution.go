package models

import (
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
)

// SequenceExecution contains all required information needed by the shipyard controller on how to preceed within a task sequence.
// An instance of SequenceExecution represents the execution of a sequence for a certain keptnContext within a stage.
// This means that, e.g. for a multi-stage sequence, multiple instances of this struct are maintained (one for each sequence in a given stage).
// Also, for multiple iterations of a sequence, each iteration will get a new instance.
type SequenceExecution struct {
	ID string `json:"_id" bson:"_id"`
	// SchemaVersion indicates the scheme that is used for the internal representation of the sequence execution
	SchemaVersion string `json:"schemaVersion" bson:"schemaVersion"`
	// Sequence contains the complete sequence definition
	Sequence keptnv2.Sequence        `json:"sequence" bson:"sequence"`
	Status   SequenceExecutionStatus `json:"status" bson:"status"`
	Scope    EventScope              `json:"scope" bson:"scope"`
	// InputProperties contains properties of the event which triggered the task sequence
	InputProperties map[string]interface{} `json:"inputProperties" bson:"inputProperties"`
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
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

func (r TaskExecutionResult) IsFailed() bool {
	return r.Result == keptnv2.ResultFailed
}

func (r TaskExecutionResult) IsErrored() bool {
	return r.Status == keptnv2.StatusErrored
}

type TaskExecutionState struct {
	Name        string      `json:"name" bson:"name"`
	TriggeredID string      `json:"triggeredID" bson:"triggeredID"`
	Events      []TaskEvent `json:"events" bson:"events"`
}

// GetNextTaskOfSequence returns the next task of a sequence, based on its current execution state. If no task is remaining, or if a previous task
// could not be completed successfully, it will return nil.
func (e *SequenceExecution) GetNextTaskOfSequence() *keptnv2.Task {
	if e.GetLastTaskExecutionResult().IsFailed() || e.GetLastTaskExecutionResult().IsErrored() {
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

func (e *SequenceExecution) GetLastTaskExecutionResult() TaskExecutionResult {
	if len(e.Status.PreviousTasks) == 0 {
		return TaskExecutionResult{}
	}
	return e.Status.PreviousTasks[len(e.Status.PreviousTasks)-1]
}

// CompleteCurrentTask completes the current task and appends the aggregated result of the current task to the list of already completed tasks.
func (e *SequenceExecution) CompleteCurrentTask() (keptnv2.ResultType, keptnv2.StatusType) {
	var result keptnv2.ResultType
	var status keptnv2.StatusType
	if e.Status.CurrentTask.IsFailed() {
		result = keptnv2.ResultFailed
	} else if e.Status.CurrentTask.IsWarning() {
		result = keptnv2.ResultWarning
	} else {
		result = keptnv2.ResultPass
	}
	if e.Status.CurrentTask.IsErrored() {
		status = keptnv2.StatusErrored
	} else {
		status = keptnv2.StatusSucceeded
	}

	var mergedProperties interface{}

	for _, taskEvent := range e.Status.CurrentTask.Events {
		if keptnv2.IsFinishedEventType(taskEvent.EventType) && taskEvent.Properties != nil {
			mergedProperties = common.Merge(mergedProperties, taskEvent.Properties)
		}
	}

	executionResult := TaskExecutionResult{
		Name:        e.Status.CurrentTask.Name,
		TriggeredID: e.Status.CurrentTask.TriggeredID,
		Result:      result,
		Status:      status,
	}
	if mergedPropertiesMap, ok := mergedProperties.(map[string]interface{}); ok {
		executionResult.Properties = mergedPropertiesMap
	}
	e.Status.PreviousTasks = append(
		e.Status.PreviousTasks,
		executionResult,
	)
	e.Status.CurrentTask = TaskExecutionState{}
	return result, status
}

// GetNextTriggeredEventData generates a map representing the event payload for the next task.triggered event. For this, it will merge the following properties:
// - The payload provided by the event that triggered the sequence
// - The properties of the task, defined in the sequence definition
// - The results of the already completed tasks of the sequence
func (e *SequenceExecution) GetNextTriggeredEventData() map[string]interface{} {
	eventPayload := map[string]interface{}{}

	if e.InputProperties != nil {
		inputProperties := common.CopyMap(e.InputProperties)
		eventPayload = common.Merge(eventPayload, inputProperties).(map[string]interface{})
	}

	eventPayload["project"] = e.Scope.Project
	eventPayload["stage"] = e.Scope.Stage
	eventPayload["service"] = e.Scope.Service

	if len(e.Status.PreviousTasks) > 0 {
		for _, previousTask := range e.Status.PreviousTasks {
			eventPayload = common.Merge(eventPayload, previousTask.Properties).(map[string]interface{})
		}
		lastTaskIndex := len(e.Status.PreviousTasks) - 1
		eventPayload["result"] = e.Status.PreviousTasks[lastTaskIndex].Result
		eventPayload["status"] = e.Status.PreviousTasks[lastTaskIndex].Status
	}

	nextTask := e.GetNextTaskOfSequence()
	if nextTask != nil && nextTask.Properties != nil {
		eventPayload[nextTask.Name] = common.Merge(eventPayload[nextTask.Name], nextTask.Properties)
	}

	// remove any messages set by previous task executors
	if eventPayload["message"] != nil {
		eventPayload["message"] = ""
	}

	return eventPayload
}

func (e *SequenceExecution) IsPaused() bool {
	return e.Status.State == models.SequencePaused
}

// CanBePaused determines whether a sequence can be paused, based on its current state. E.g. a finished sequence cannot be paused
func (e *SequenceExecution) CanBePaused() bool {
	return e.Status.State == models.SequenceStartedState || e.Status.State == models.SequenceWaitingState || e.Status.State == models.SequenceTriggeredState || e.Status.State == models.SequenceWaitingForApprovalState
}

// Pause tries to pause the sequence execution, based on its current state. If it was successful, returns true. If it could not be paused, false is returned
func (e *SequenceExecution) Pause() bool {
	if !e.CanBePaused() {
		return false
	}
	e.Status.StateBeforePause = e.Status.State
	e.Status.State = models.SequencePaused
	return true
}

// Resume tries to resume the sequence execution, based on its current state. If it was successful, returns true. If it could not be paused, false is returned
func (e *SequenceExecution) Resume() bool {
	if !e.IsPaused() {
		return false
	}
	e.Status.State = e.Status.StateBeforePause
	return true
}

// SetNextCurrentTask updates the Current task of the sequence and sets the current state appropriately, considering the special logic that should be applied for approval tasks
func (e *SequenceExecution) SetNextCurrentTask(taskName, triggeredEventID string) {
	e.Status.CurrentTask = TaskExecutionState{
		Name:        taskName,
		TriggeredID: triggeredEventID,
		Events:      []TaskEvent{},
	}

	// special handling for approval events
	nextState := models.SequenceStartedState
	if taskName == keptnv2.ApprovalTaskName {
		nextState = models.SequenceWaitingForApprovalState
	}

	if e.IsPaused() {
		e.Status.StateBeforePause = nextState
	} else {
		e.Status.State = nextState
	}
}

// IsFinished indicates if a task is finished, i.e. the number of task.started and task.finished events line up
func (e *TaskExecutionState) IsFinished() bool {
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

func (e *TaskExecutionState) IsFailed() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == keptnv2.ResultFailed {
				return true
			}
		}
	}
	return false
}

func (e *TaskExecutionState) IsWarning() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == keptnv2.ResultWarning {
				return true
			}
		}
	}
	return false
}

func (e *TaskExecutionState) IsPassed() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == keptnv2.ResultFailed || event.Result == keptnv2.ResultWarning {
				return false
			}
		}
	}
	return true
}

func (e *TaskExecutionState) IsErrored() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Status == keptnv2.StatusErrored {
				return true
			}
		}
	}
	return false
}

type TaskEvent struct {
	EventType  string                 `json:"eventType" bson:"eventType"`
	Source     string                 `json:"source" bson:"source"`
	Result     keptnv2.ResultType     `json:"result" bson:"result"`
	Status     keptnv2.StatusType     `json:"status" bson:"status"`
	Time       string                 `json:"time" bson:"time"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

type SequenceExecutionFilter struct {
	Scope              EventScope
	Status             []string
	Name               string
	CurrentTriggeredID string
}

type SequenceExecutionUpsertOptions struct {
	CheckUniqueTriggeredID bool
	Replace                bool
}
