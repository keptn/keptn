package handlers

import (
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/api/models"
	"golang.org/x/exp/slices"
)

// EventValidationError is a custom error used to represent
// errors during validation of keptn events sent to the API
type EventValidationError struct{ Msg string }

func (e EventValidationError) Error() string { return e.Msg }

func (e EventValidationError) Unwrap() error { return e }

type validateFn func(models.KeptnContextExtendedCE) error

// allow only "sh.keptn.event.<task>.<started|finished>" events
var allowedTaskActions = []string{"started", "finished", "invalidated"}

// allow only "sh.keptn.event.<sequence>.triggered" events
var allowedSequenceActions = []string{"triggered"}

// specify validation functions for sequence triggered events
var sequenceEventValidators = map[string]validateFn{"triggered": validateSequenceTriggeredEvent}

// specify validation functions for task.started and task.finished events
var taskEventValidators = map[string]validateFn{
	"started":  validateTaskStartedEvent,
	"finished": validateTaskFinishedEvent,
}

// Validate takes a KeptnContextExtendedCE value and validates its content.
// If the event is valid, the returned error is nil. If it is not valid
// a EventValidationError is returned
func Validate(e models.KeptnContextExtendedCE) error {
	// validate "special" events
	if *e.Type == "sh.keptn.log.error" || *e.Type == "sh.keptn.events.problem" {
		//no specific validation
		return nil
	}
	if *e.Type == "sh.keptn.event.monitoring.configure" {
		return validateMonitoringConfigureEvent(e)
	}

	// validate events related to sequences
	if v0_2_0.IsSequenceEventType(*e.Type) {
		return validate(e, allowedSequenceActions, sequenceEventValidators)
	}

	// validate events related to tasks
	if v0_2_0.IsTaskEventType(*e.Type) {
		return validate(e, allowedTaskActions, taskEventValidators)
	}

	// received unknown type of event
	return &EventValidationError{Msg: "unknown event type"}
}

func validate(event models.KeptnContextExtendedCE, allowedActions []string, validators map[string]validateFn) error {
	action, err := v0_2_0.ParseEventKind(*event.Type)
	if err != nil {
		return err
	}
	if !slices.Contains(allowedActions, action) {
		return &EventValidationError{Msg: "unknown event type action: " + action}
	}
	if _, ok := validators[action]; !ok {
		return nil
	}
	return validators[action](event)

}

// validateTaskStartedEvent contains logic to validate "sh.keptn.event.task.started" events
func validateTaskStartedEvent(e models.KeptnContextExtendedCE) error {
	var eventData v0_2_0.EventData
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return &EventValidationError{Msg: "could not parse common event data"}
	}
	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return &EventValidationError{Msg: "mandatory field(s) 'project', 'stage' or 'service' missing"}
	}

	if e.Triggeredid == "" {
		return &EventValidationError{Msg: "mandatory field 'triggeredid' missing"}
	}
	return nil
}

// validateTaskFinishedEvent contains logic that validates "sh.keptn.event.task.finished" events
func validateTaskFinishedEvent(e models.KeptnContextExtendedCE) error {
	var eventData v0_2_0.EventData
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return &EventValidationError{Msg: "could not parse common event data"}
	}
	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return &EventValidationError{Msg: "mandatory field(s) 'project', 'stage' or 'service' missing"}
	}
	if eventData.Result == "" {
		return &EventValidationError{Msg: "result field is not set"}
	}
	if eventData.Status == "" {
		return &EventValidationError{Msg: "status field is not set"}
	}

	if e.Triggeredid == "" {
		return &EventValidationError{Msg: "mandatory field 'triggeredid' missing"}
	}
	return nil
}

// validateSequenceTriggeredEvent contains logic that validates "sh.keptn.dev.sequence.triggered" events
func validateSequenceTriggeredEvent(e models.KeptnContextExtendedCE) error {
	var eventData v0_2_0.EventData
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return &EventValidationError{Msg: "could not parse common event data"}
	}
	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return &EventValidationError{Msg: "mandatory field(s) 'project', 'stage' or 'service' missing"}
	}
	stage, _, _, err := v0_2_0.ParseSequenceEventType(*e.Type)
	if err != nil {
		return &EventValidationError{Msg: "unknown event type"}
	}
	if eventData.Stage != stage {
		return &EventValidationError{Msg: "stage name in event data and in event type does not match"}
	}
	return nil
}

func validateMonitoringConfigureEvent(e models.KeptnContextExtendedCE) error {
	type data struct {
		v0_2_0.EventData
		Type string `json:"type,omitempty"`
	}

	var eventData data
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return &EventValidationError{Msg: "could not parse common event data"}
	}

	if eventData.Project == "" || eventData.Type == "" {
		return &EventValidationError{Msg: "mandatory fields(s) 'project', 'service' or 'type' missing"}
	}

	return nil
}
