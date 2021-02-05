package handler

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type ApprovalTriggeredEventHandler struct {
	keptn *keptnv2.Keptn
}

// NewApprovalTriggeredEventHandler returns a new approal.triggered event handler
func NewApprovalTriggeredEventHandler(keptn *keptnv2.Keptn) *ApprovalTriggeredEventHandler {
	return &ApprovalTriggeredEventHandler{keptn: keptn}
}

// IsTypeHandled godoc
func (a *ApprovalTriggeredEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName)
}

// Handle godoc
func (a *ApprovalTriggeredEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnv2.Keptn) {
	data := &keptnv2.ApprovalTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		a.keptn.Logger.Error(fmt.Sprintf("failed to parse ApprovalTriggeredEventData: %v", err))
		return
	}

	outgoingEvents := a.handleApprovalTriggeredEvent(*data, event.Context.GetID(), keptnHandler.KeptnContext)
	sendEvents(keptnHandler, outgoingEvents, a.keptn.Logger)
}

func (a *ApprovalTriggeredEventHandler) handleApprovalTriggeredEvent(inputEvent keptnv2.ApprovalTriggeredEventData,
	triggeredID, shkeptncontext string) []cloudevents.Event {

	outgoingEvents := make([]cloudevents.Event, 0)
	if inputEvent.Result == keptnv2.ResultPass && inputEvent.Approval.Pass == keptnv2.ApprovalAutomatic ||
		inputEvent.Result == keptnv2.ResultWarning && inputEvent.Approval.Warning == keptnv2.ApprovalAutomatic {

		startedEvent := a.getApprovalStartedEvent(inputEvent, triggeredID, shkeptncontext)
		outgoingEvents = append(outgoingEvents, *startedEvent)
		a.keptn.Logger.Info(fmt.Sprintf("Automatically approve release of service %s of project %s and current stage %s",
			inputEvent.Service, inputEvent.Project, inputEvent.Stage))

		finishedEvent := a.getApprovalFinishedEvent(inputEvent, keptnv2.ResultPass, triggeredID, shkeptncontext)
		outgoingEvents = append(outgoingEvents, *finishedEvent)
	} else if inputEvent.Result == keptnv2.ResultFailed {
		// Handle case if an ApprovalTriggered event was sent even the evaluation result is failed

		startedEvent := a.getApprovalStartedEvent(inputEvent, triggeredID, shkeptncontext)
		outgoingEvents = append(outgoingEvents, *startedEvent)

		a.keptn.Logger.Info(fmt.Sprintf("Disapprove release of service %s of project %s and current stage %s because"+
			"the previous step failed", inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		finishedEvent := a.getApprovalFinishedEvent(inputEvent, keptnv2.ResultFailed, triggeredID, shkeptncontext)
		outgoingEvents = append(outgoingEvents, *finishedEvent)
	}

	return outgoingEvents
}

func (a *ApprovalTriggeredEventHandler) getApprovalStartedEvent(inputEvent keptnv2.ApprovalTriggeredEventData, triggeredID, shkeptncontext string) *cloudevents.Event {
	approvalFinishedEvent := keptnv2.ApprovalStartedEventData{
		EventData: keptnv2.EventData{
			Project: inputEvent.Project,
			Stage:   inputEvent.Stage,
			Service: inputEvent.Service,
			Labels:  inputEvent.Labels,
			Status:  keptnv2.StatusSucceeded,
			Message: "",
		},
	}

	return getCloudEvent(approvalFinishedEvent, keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, triggeredID)
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
