package handlers

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/api/models"
	"golang.org/x/exp/slices"
)

var ErrUnknownEventType = errors.New("invalid event: unknown event type")
var ErrDisallowedEventKind = errors.New("invalid event: event action is not allowed")
var ErrCommonEventDataInvalid = errors.New("invalid event: could not parse common event data")
var ErrCommonEventDataMissing = errors.New("invalid event: mandatory field(s) 'project', 'stage' or 'service' missing")
var ErrStageMismatch = errors.New("invalid event: stage name in event data and in event type does not match")
var ErrResultFieldMissing = errors.New("invalid event: result field is not set")
var ErrStatusFieldMissing = errors.New("invalid event: status field is not set")

type ValidateFn func(models.KeptnContextExtendedCE) error

var allowedTaskActions = []string{"started", "finished"}
var allowedSequenceActions = []string{"triggered"}
var sequenceEventValidators = map[string]ValidateFn{"triggered": validateSequenceTriggeredEvent}
var taskEventValidators = map[string]ValidateFn{
	"started":  validateTaskStartedEvent,
	"finished": validateTaskFinishedEvent,
}

func Validate(e models.KeptnContextExtendedCE) error {
	if v0_2_0.IsSequenceEventType(*e.Type) {
		return validate(e, allowedSequenceActions, sequenceEventValidators)
	}
	if v0_2_0.IsTaskEventType(*e.Type) {
		return validate(e, allowedTaskActions, taskEventValidators)
	}
	if *e.Type == "sh.keptn.log.error" {
		return validateErrorLogEvent(e)
	}
	return ErrUnknownEventType
}

func validateErrorLogEvent(event models.KeptnContextExtendedCE) error {
	// TODO: implement what makes sense
	return nil
}

func validate(event models.KeptnContextExtendedCE, allowedKinds []string, validators map[string]ValidateFn) error {
	kind, err := v0_2_0.ParseEventKind(*event.Type)
	if err != nil {
		return err
	}
	if !slices.Contains(allowedKinds, kind) {
		return fmt.Errorf("invalid event kind/action: %s, %w", kind, ErrDisallowedEventKind)
	}
	if _, ok := validators[kind]; !ok {
		return nil
	}
	return validators[kind](event)

}

func validateTaskStartedEvent(e models.KeptnContextExtendedCE) error {
	var eventData v0_2_0.EventData
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return ErrCommonEventDataInvalid
	}
	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return ErrCommonEventDataMissing
	}
	return nil
}

func validateTaskFinishedEvent(e models.KeptnContextExtendedCE) error {
	var eventData v0_2_0.EventData
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return ErrCommonEventDataInvalid
	}
	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return ErrCommonEventDataMissing
	}
	if eventData.Result == "" {
		return ErrResultFieldMissing
	}
	if eventData.Status == "" {
		return ErrStatusFieldMissing
	}
	return nil
}

func validateSequenceTriggeredEvent(e models.KeptnContextExtendedCE) error {
	var eventData v0_2_0.EventData
	if err := v0_2_0.Decode(e.Data, &eventData); err != nil {
		return ErrCommonEventDataInvalid
	}
	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return ErrCommonEventDataMissing
	}
	stage, _, _, err := v0_2_0.ParseSequenceEventType(*e.Type)
	if err != nil {
		return ErrUnknownEventType
	}
	if eventData.Stage != stage {
		return ErrStageMismatch
	}
	return nil
}
