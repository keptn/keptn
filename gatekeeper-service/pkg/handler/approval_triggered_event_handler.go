package handler

import (
	"fmt"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

type ApprovalTriggeredEventHandler struct {
	logger *keptnevents.Logger
}

func NewApprovalTriggeredEventHandler(l *keptnevents.Logger) *ApprovalTriggeredEventHandler {
	return &ApprovalTriggeredEventHandler{logger: l}
}

func (a *ApprovalTriggeredEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.ApprovalTriggeredEventType
}

func (a *ApprovalTriggeredEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn, shipyard *keptnevents.Shipyard) {

	data := &keptnevents.ApprovalTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		a.logger.Error(fmt.Sprintf("failed to parse ApprovalTriggeredEventData: %v", err))
		return
	}

	outgoingEvents := a.handleApprovalTriggeredEvent(*data, event.Context.GetID(), keptnHandler.KeptnContext, *shipyard)
	sendEvents(keptnHandler, outgoingEvents, a.logger)
}

func (a *ApprovalTriggeredEventHandler) handleApprovalTriggeredEvent(inputEvent keptnevents.ApprovalTriggeredEventData, triggeredId, shkeptncontext string,
	shipyard keptnevents.Shipyard) []cloudevents.Event {

	outgoingEvents := make([]cloudevents.Event, 0)
	if inputEvent.Result == PassResult && a.getApprovalStrategyForPass(inputEvent.Stage, shipyard) == keptnevents.Automatic ||
		inputEvent.Result == WarningResult && a.getApprovalStrategyForWarning(inputEvent.Stage, shipyard) == keptnevents.Automatic {
		// Pass
		a.logger.Info(fmt.Sprintf("Automatically approve image %s for service %s of project %s and current stage %s",
			inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		outgoingEvents = append(outgoingEvents, *a.getApprovalFinishedEvent(inputEvent, PassResult, triggeredId, shkeptncontext))
	} else if inputEvent.Result == FailResult {
		// Handle case if an ApprovalTriggered event was sent even the evaluation result is failed
		a.logger.Info(fmt.Sprintf("Disapprove image %s for service %s of project %s and current stage %s because"+
			"the evaluation result is fail", inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		outgoingEvents = append(outgoingEvents, *a.getApprovalFinishedEvent(inputEvent, FailResult, triggeredId, shkeptncontext))
	}

	return outgoingEvents
}

func (a *ApprovalTriggeredEventHandler) getApprovalStrategyForPass(stageName string, shipyard keptnevents.Shipyard) keptnevents.ApprovalStrategy {
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName && stage.ApprovalStrategy != nil {
			return stage.ApprovalStrategy.Pass
		}
	}
	// Implements the default behavior if the Shipyard does not specify an ApprovalStrategy
	return keptnevents.Automatic
}

func (a *ApprovalTriggeredEventHandler) getApprovalStrategyForWarning(stageName string, shipyard keptnevents.Shipyard) keptnevents.ApprovalStrategy {
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName && stage.ApprovalStrategy != nil {
			return stage.ApprovalStrategy.Warning
		}
	}
	// Implements the default behavior if the Shipyard does not specify an ApprovalStrategy
	return keptnevents.Automatic
}

func (a *ApprovalTriggeredEventHandler) getApprovalFinishedEvent(inputEvent keptnevents.ApprovalTriggeredEventData, result, triggeredId, shkeptncontext string) *cloudevents.Event {

	approvalFinishedEvent := keptnevents.ApprovalFinishedEventData{
		Project:            inputEvent.Project,
		Service:            inputEvent.Service,
		Stage:              inputEvent.Stage,
		TestStrategy:       inputEvent.TestStrategy,
		DeploymentStrategy: inputEvent.DeploymentStrategy,
		Tag:                inputEvent.Tag,
		Image:              inputEvent.Image,
		Labels:             inputEvent.Labels,
		Approval: keptnevents.ApprovalData{
			TriggeredID: triggeredId,
			Result:      result,
			Status:      SucceededResult,
		},
	}
	return getCloudEvent(approvalFinishedEvent, keptnevents.ApprovalFinishedEventType, shkeptncontext)
}
