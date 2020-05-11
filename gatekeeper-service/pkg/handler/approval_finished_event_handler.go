package handler

import (
	"fmt"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

type ApprovalFinishedEventHandler struct {
	logger *keptnevents.Logger
}

func NewApprovalFinishedEventHandler(l *keptnevents.Logger) *ApprovalFinishedEventHandler {
	return &ApprovalFinishedEventHandler{logger: l}
}

func (a *ApprovalFinishedEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.ApprovalFinishedEventType
}

func (a *ApprovalFinishedEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn, shipyard *keptnevents.Shipyard) {
	data := &keptnevents.ApprovalFinishedEventData{}
	if err := event.DataAs(data); err != nil {
		a.logger.Error(fmt.Sprintf("failed to parse ApprovalTriggeredEventData: %v", err))
	}

	outgoingEvents := a.handleApprovalFinishedEvent(*data, keptnHandler.KeptnContext, *shipyard)
	sendEvents(keptnHandler, outgoingEvents, a.logger)
}

func (a *ApprovalFinishedEventHandler) handleApprovalFinishedEvent(inputEvent keptnevents.ApprovalFinishedEventData, shkeptncontext string,
	shipyard keptnevents.Shipyard) []cloudevents.Event {

	outgoingEvents := make([]cloudevents.Event, 0)
	if inputEvent.Approval.Status != SucceededResult {
		a.logger.Info(fmt.Sprintf("Approval finished with failed status for "+
			"image %s for service %s of project %s and current stage %s received",
			inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
	} else {
		if inputEvent.Approval.Result == PassResult {
			a.logger.Info(fmt.Sprintf("Approval for image %s for service %s of project %s and current stage %s received",
				inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))

			// TODO: Check image using inputEvent.Approval.TriggeredID
			image := inputEvent.Image
			if inputEvent.Tag != "" {
				image += ":" + inputEvent.Tag
			}
			if event := getPromotionEvent(inputEvent.Project, inputEvent.Stage, inputEvent.Service, image,
				shkeptncontext, inputEvent.Labels, shipyard, a.logger); event != nil {
				outgoingEvents = append(outgoingEvents, *event)
			}
		} else {
			a.logger.Info(fmt.Sprintf("Rejection for image %s for service %s of project %s and current stage %s received",
				inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		}
	}

	return outgoingEvents
}
