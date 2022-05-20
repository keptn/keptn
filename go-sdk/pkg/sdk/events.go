package sdk

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"time"
)

func createStartedEvent(source string, parentEvent models.KeptnContextExtendedCE) (*models.KeptnContextExtendedCE, error) {
	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	startedEventType, err := keptnv2.ReplaceEventTypeKind(*parentEvent.Type, "started")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.started' event for parent event %s: %w", parentEvent.ID, err)
	}
	eventData := keptnv2.EventData{}
	parentEvent.DataAs(&eventData)
	return createEvent(source, startedEventType, parentEvent, eventData), nil
}

func createFinishedEvent(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}) (*models.KeptnContextExtendedCE, error) {
	if parentEvent.Type == nil {
		return nil, fmt.Errorf("unable to get keptn event type from event %s", parentEvent.ID)
	}

	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	finishedEventType, err := keptnv2.ReplaceEventTypeKind(*parentEvent.Type, "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, *parentEvent.Type)
	}
	var genericEventData map[string]interface{}
	err = keptnv2.Decode(eventData, &genericEventData)
	if err != nil || genericEventData == nil {
		return nil, fmt.Errorf("unable to decode generic event data")
	}

	if genericEventData["status"] == nil || genericEventData["status"] == "" {
		genericEventData["status"] = "succeeded"
	}

	if genericEventData["result"] == nil || genericEventData["result"] == "" {
		genericEventData["result"] = "pass"
	}
	return createEvent(source, finishedEventType, parentEvent, genericEventData), nil
}

func createFinishedEventWithError(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	if errVal == nil {
		errVal = &Error{}
	}
	commonEventData := keptnv2.EventData{}
	if eventData == nil {
		parentEvent.DataAs(&commonEventData)
	}
	commonEventData.Result = errVal.ResultType
	commonEventData.Status = errVal.StatusType
	commonEventData.Message = errVal.Message

	finishedEventType, err := keptnv2.ReplaceEventTypeKind(*parentEvent.Type, "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event for parent event %s: %w", parentEvent.ID, err)
	}
	return createEvent(source, finishedEventType, parentEvent, commonEventData), nil
}

func createErrorEvent(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	if errVal == nil {
		errVal = &Error{}
	}

	if keptnv2.IsTaskEventType(*parentEvent.Type) && keptnv2.IsTriggeredEventType(*parentEvent.Type) {
		errorFinishedEvent, err := createFinishedEventWithError(source, parentEvent, eventData, errVal)
		if err != nil {
			return nil, err
		}
		return errorFinishedEvent, nil
	}
	errorLogEvent, err := createErrorLogEvent(source, parentEvent, eventData, errVal)
	if err != nil {
		return nil, err
	}
	return errorLogEvent, nil
}

func createErrorLogEvent(source string, parentEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	if parentEvent.Type == nil {
		return nil, fmt.Errorf("unable to get keptn event type from parent event %s", parentEvent.ID)
	}

	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	if errVal == nil {
		errVal = &Error{}
	}

	if keptnv2.IsTaskEventType(*parentEvent.Type) && keptnv2.IsTriggeredEventType(*parentEvent.Type) {
		errorFinishedEvent, err := createFinishedEventWithError(source, parentEvent, eventData, errVal)
		if err != nil {
			return nil, err
		}
		return errorFinishedEvent, nil
	}
	errorEventData := keptnv2.ErrorLogEvent{}
	if eventData == nil {
		parentEvent.DataAs(&errorEventData)
	}
	if keptnv2.IsTaskEventType(*parentEvent.Type) {
		taskName, _, err := keptnv2.ParseTaskEventType(*parentEvent.Type)
		if err == nil && taskName != "" {
			errorEventData.Task = taskName
		}
	}
	errorEventData.Message = errVal.Message
	if parentEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from parent event %s", parentEvent.ID)
	}
	return createEvent(source, keptnv2.ErrorLogEventName, parentEvent, errorEventData), nil
}

func createEvent(source string, eventType string, parentEvent models.KeptnContextExtendedCE, eventData interface{}) *models.KeptnContextExtendedCE {
	return &models.KeptnContextExtendedCE{
		ID:                 uuid.NewString(),
		Triggeredid:        parentEvent.ID,
		Shkeptncontext:     parentEvent.Shkeptncontext,
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               eventData,
		Source:             strutils.Stringp(source),
		Shkeptnspecversion: shkeptnspecversion,
		Specversion:        cloudeventsversion,
		Time:               time.Now().UTC(),
		Type:               strutils.Stringp(eventType),
	}
}
