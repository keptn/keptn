package handler

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/types"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	logger "github.com/sirupsen/logrus"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type ApprovalTriggeredEventHandler struct {
	keptn *keptnv2.Keptn
}

// NewApprovalTriggeredEventHandler returns a new approval.triggered event handler
func NewApprovalTriggeredEventHandler(ctx context.Context, event cloudevents.Event) (*ApprovalTriggeredEventHandler, error) {
	eventSender := ctx.Value(types.EventSenderKey).(controlplane.EventSender)
	keptnHandlerV2, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{EventSender: &CPEventSender{Sender: eventSender}})
	if err != nil {
		logger.WithError(err).Error("Failed to initialize Keptn handler")
		return nil, err
	}
	return &ApprovalTriggeredEventHandler{keptn: keptnHandlerV2}, nil
}

// IsTypeHandled godoc
func (a *ApprovalTriggeredEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName)
}

// getResult returns the result from either the preceeding task, or if result is empty, evaluation.result
func getResult(data keptnv2.ApprovalTriggeredEventData, event cloudevents.Event) keptnv2.ResultType {
	if data.Result != "" {
		return data.Result
	}

	// handle the case of no result being present (see https://github.com/keptn/keptn/issues/4391)
	// check if evaluation.finished event data are present
	evaluationFinishedData := &keptnv2.EvaluationFinishedEventData{}
	if err := event.DataAs(evaluationFinishedData); err == nil {
		return keptnv2.ResultType(evaluationFinishedData.Evaluation.Result)
	}

	// no suitable result -> we will stay empty
	return ""
}

// Handle godoc
func (a *ApprovalTriggeredEventHandler) Handle(event cloudevents.Event) error {
	data := &keptnv2.ApprovalTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		logger.WithError(err).Error("failed to parse ApprovalTriggeredEventData")
		return err
	}

	// handle the case of no result being present (see https://github.com/keptn/keptn/issues/4391)
	data.Result = getResult(*data, event)

	outgoingEvents := a.handleApprovalTriggeredEvent(*data, event.Context.GetID(), a.keptn.KeptnContext)
	sendEvents(a.keptn, outgoingEvents)

	return nil
}

func (a *ApprovalTriggeredEventHandler) handleApprovalTriggeredEvent(inputEvent keptnv2.ApprovalTriggeredEventData,
	triggeredID, shkeptncontext string) []cloudevents.Event {
	outgoingEvents := make([]cloudevents.Event, 0)

	startedEvent := a.getApprovalStartedEvent(inputEvent, triggeredID, shkeptncontext)
	outgoingEvents = append(outgoingEvents, *startedEvent)

	if inputEvent.Result == keptnv2.ResultPass && inputEvent.Approval.Pass == keptnv2.ApprovalAutomatic ||
		inputEvent.Result == keptnv2.ResultWarning && inputEvent.Approval.Warning == keptnv2.ApprovalAutomatic {

		logger.Info(fmt.Sprintf("Automatically approve release of service %s of project %s and current stage %s",
			inputEvent.Service, inputEvent.Project, inputEvent.Stage))

		finishedEvent := a.getApprovalFinishedEvent(inputEvent, keptnv2.ResultPass, triggeredID, shkeptncontext)
		outgoingEvents = append(outgoingEvents, *finishedEvent)
	} else if inputEvent.Result == keptnv2.ResultFailed {
		// Handle case if an ApprovalTriggered event was sent even the evaluation result is failed
		logger.Info(fmt.Sprintf("Disapprove release of service %s of project %s and current stage %s because"+
			"the previous step failed", inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		finishedEvent := a.getApprovalFinishedEvent(inputEvent, keptnv2.ResultFailed, triggeredID, shkeptncontext)
		outgoingEvents = append(outgoingEvents, *finishedEvent)
	}

	return outgoingEvents
}

func (a *ApprovalTriggeredEventHandler) getApprovalStartedEvent(inputEvent keptnv2.ApprovalTriggeredEventData, triggeredID, shkeptncontext string) *cloudevents.Event {
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

	return getCloudEvent(approvalStartedEvent, keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, triggeredID)
}

func (a *ApprovalTriggeredEventHandler) getApprovalFinishedEvent(inputEvent keptnv2.ApprovalTriggeredEventData,
	result keptnv2.ResultType, triggeredID, shkeptncontext string) *cloudevents.Event {
	approvalFinishedEvent := keptnv2.ApprovalFinishedEventData{
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

	return getCloudEvent(approvalFinishedEvent, keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, triggeredID)
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
