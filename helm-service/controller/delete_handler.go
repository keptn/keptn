package controller

import (
	"fmt"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	configutils "github.com/keptn/go-utils/pkg/api/utils"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// DeleteHandler handles sh.keptn.internal.event.service.delete events
type DeleteHandler struct {
	Handler
}

// NewDeleteHandler creates a new DeleteHandler
func NewDeleteHandler(keptnHandler *keptnv2.Keptn, configServiceURL string) *DeleteHandler {
	return &DeleteHandler{
		Handler: NewHandlerBase(keptnHandler, configServiceURL),
	}
}

// HandleEvent takes the sh.keptn.internal.event.service.delete event and deletes the service in all stages
func (h *DeleteHandler) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) {

	defer closeLogger(h.getKeptnHandler())

	serviceDeleteEvent := keptnv2.ServiceDeleteFinishedEventData{}

	err := ce.DataAs(&serviceDeleteEvent)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		h.handleError(ce.ID(), err, keptnv2.ServiceDeleteTaskName, h.getFinishedEventDataForError(serviceDeleteEvent.EventData, err))
		return
	}

	h.getKeptnHandler().Logger.Info(fmt.Sprintf("Starting uninstalling releases for service %s of project %s",
		serviceDeleteEvent.Service, serviceDeleteEvent.Project))

	stageHandler := configutils.NewStageHandler(h.getConfigServiceURL())
	stages, err := stageHandler.GetAllStages(serviceDeleteEvent.Project)
	if err != nil {
		err = fmt.Errorf("error when getting all stages: %v", err)
		h.handleError(ce.ID(), err, keptnv2.ServiceDeleteTaskName, h.getFinishedEventDataForError(serviceDeleteEvent.EventData, err))
		return
	}

	allReleasesSuccessfullyUnistalled := true
	for _, stage := range stages {
		h.getKeptnHandler().Logger.Info(fmt.Sprintf("Uninstalling Helm releases for service %s in "+
			"stage %s and project %s", serviceDeleteEvent.Service, stage.StageName, serviceDeleteEvent.Project))

		namespace := serviceDeleteEvent.Project + "-" + stage.StageName
		releaseName := namespace + "-" + serviceDeleteEvent.Service
		if err := h.getHelmExecutor().UninstallRelease(releaseName, namespace); err != nil {
			h.getKeptnHandler().Logger.Error(err.Error())
			allReleasesSuccessfullyUnistalled = false
		}
		if err := h.getHelmExecutor().UninstallRelease(releaseName+"-generated", namespace); err != nil {
			h.getKeptnHandler().Logger.Error(err.Error())
			allReleasesSuccessfullyUnistalled = false
		}
	}

	if allReleasesSuccessfullyUnistalled {
		h.getKeptnHandler().Logger.Info(fmt.Sprintf("All Helm releases for service %s in project %s successfully uninstalled",
			serviceDeleteEvent.Service, serviceDeleteEvent.Project))
	}

	// Send finished event
	msg := fmt.Sprintf("Finished uninstalling service %s in project %s", serviceDeleteEvent.Service, serviceDeleteEvent.Project)
	data := h.getFinishedEventData(serviceDeleteEvent.EventData, keptnv2.StatusSucceeded, keptnv2.ResultPass, msg)
	if err := h.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), data); err != nil {
		h.handleError(ce.ID(), err, keptnv2.ServiceDeleteTaskName, h.getFinishedEventDataForError(serviceDeleteEvent.EventData, err))
	}
}

func (h *DeleteHandler) getFinishedEventData(inEventData keptnv2.EventData, status keptnv2.StatusType, result keptnv2.ResultType,
	message string) keptnv2.ServiceDeleteFinishedEventData {

	inEventData.Status = status
	inEventData.Result = result
	inEventData.Message = message

	return keptnv2.ServiceDeleteFinishedEventData{
		EventData: inEventData,
	}
}

func (h *DeleteHandler) getFinishedEventDataForError(inEventData keptnv2.EventData, err error) keptnv2.ServiceDeleteFinishedEventData {
	return h.getFinishedEventData(inEventData, keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
}
