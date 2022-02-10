package handler

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	logger "github.com/sirupsen/logrus"
)

type ApprovalTriggeredEventHandler struct {
}

func NewApprovalTriggeredEventHandler() *ApprovalTriggeredEventHandler {
	return &ApprovalTriggeredEventHandler{}
}

func (a *ApprovalTriggeredEventHandler) Execute(k sdk.IKeptn, event sdk.KeptnEvent) (interface{}, *sdk.Error) {

	data := &keptnv2.ApprovalTriggeredEventData{}

	if err := keptnv2.Decode(event.Data, data); err != nil {
		outgoingEvent := a.getApprovalFinishedEvent(*data, keptnv2.ResultFailed)
		outgoingEvent.Status = keptnv2.StatusErrored
		k.SendFinishedEvent(event, outgoingEvent)
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "Could not decode input event data"}
	}

	// handle the case of no result being present (see https://github.com/keptn/keptn/issues/4391)
	data.Result = getResult(*data, event)

	startedEvent := a.getApprovalStartedEvent(*data, event)
	k.SendStartedEvent(*startedEvent)

	outgoingEvent := a.handleApprovalTriggeredEvent(*data, event)

	k.SendFinishedEvent(*startedEvent, outgoingEvent)
	return outgoingEvent, nil
}

// IsTypeHandled godoc
func (a *ApprovalTriggeredEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName)
}

// getResult returns the result from either the preceeding task, or if result is empty, evaluation.result
func getResult(data keptnv2.ApprovalTriggeredEventData, event sdk.KeptnEvent) keptnv2.ResultType {
	if data.Result != "" {
		return data.Result
	}

	// handle the case of no result being present (see https://github.com/keptn/keptn/issues/4391)
	// check if evaluation.finished event data are present
	evaluationFinishedData := &keptnv2.EvaluationFinishedEventData{}

	if err := keptnv2.Decode(event.Data, evaluationFinishedData); err == nil {
		return keptnv2.ResultType(evaluationFinishedData.Evaluation.Result)
	}

	// no suitable result -> we will stay empty
	return ""
}

func (a *ApprovalTriggeredEventHandler) handleApprovalTriggeredEvent(inputEvent keptnv2.ApprovalTriggeredEventData,
	event sdk.KeptnEvent) keptnv2.ApprovalFinishedEventData {
	finishedEvent := keptnv2.ApprovalFinishedEventData{}
	if inputEvent.Result == keptnv2.ResultPass && inputEvent.Approval.Pass == keptnv2.ApprovalAutomatic ||
		inputEvent.Result == keptnv2.ResultWarning && inputEvent.Approval.Warning == keptnv2.ApprovalAutomatic {

		logger.Info(fmt.Sprintf("Automatically approve release of service %s of project %s and current stage %s",
			inputEvent.Service, inputEvent.Project, inputEvent.Stage))

		finishedEvent = a.getApprovalFinishedEvent(inputEvent, keptnv2.ResultPass)
		return finishedEvent
	} else if inputEvent.Result == keptnv2.ResultFailed {
		// Handle case if an ApprovalTriggered event was sent even the evaluation result is failed
		logger.Info(fmt.Sprintf("Disapprove release of service %s of project %s and current stage %s because"+
			"the previous step failed", inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		finishedEvent = a.getApprovalFinishedEvent(inputEvent, keptnv2.ResultFailed)

	}

	return finishedEvent
}

func (a *ApprovalTriggeredEventHandler) getApprovalStartedEvent(inputEvent keptnv2.ApprovalTriggeredEventData, event sdk.KeptnEvent) *sdk.KeptnEvent {
	approvalStartedEvent := keptnv2.ApprovalStartedEventData{
		EventData: keptnv2.EventData{
			Project: inputEvent.Project,
			Stage:   inputEvent.Stage,
			Service: inputEvent.Service,
			Labels:  inputEvent.Labels,
			Status:  keptnv2.StatusSucceeded,
			Message: fmt.Sprintf("Approval strategy for result '%s': %s", string(inputEvent.Result), getApprovalStrategyForEvent(inputEvent)),
		},
	}

	return getCloudEvent(approvalStartedEvent, keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), event)
}

func (a *ApprovalTriggeredEventHandler) getApprovalFinishedEvent(inputEvent keptnv2.ApprovalTriggeredEventData,
	result keptnv2.ResultType) keptnv2.ApprovalFinishedEventData {
	return keptnv2.ApprovalFinishedEventData{
		EventData: keptnv2.EventData{
			Project: inputEvent.Project,
			Stage:   inputEvent.Stage,
			Service: inputEvent.Service,
			Labels:  inputEvent.Labels,
			Status:  keptnv2.StatusSucceeded,
			Result:  result,
			Message: "",
		},
	}

}

func getApprovalStrategyForEvent(event keptnv2.ApprovalTriggeredEventData) interface{} {
	if event.Result == keptnv2.ResultPass {
		if event.Approval.Pass != "" {
			return event.Approval.Pass
		}
		// fall back to manual if no approval strategy has been set
		return keptnv2.ApprovalManual
	}
	if event.Result == keptnv2.ResultWarning {
		if event.Approval.Warning != "" {
			return event.Approval.Warning
		}
		// fall back to manual if no approval strategy has been set
		return keptnv2.ApprovalManual
	}
	if event.Result == keptnv2.ResultFailed {
		// if we had a result=fail previously, we automatically decline
		return keptnv2.ApprovalAutomatic
	}
	// fall back to manual in al other cases
	return keptnv2.ApprovalManual
}

func getCloudEvent(data interface{}, ceType string, event sdk.KeptnEvent) *sdk.KeptnEvent {

	source := "approval-service"
	cevent := event

	cevent.Data = data
	cevent.Source = &source
	cevent.Type = &ceType

	return &cevent
}
