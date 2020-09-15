package controller

import (
	"errors"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	configutils "github.com/keptn/go-utils/pkg/api/utils"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

// DeleteHandler handles sh.keptn.internal.event.service.delete events
type DeleteHandler struct {
	Handler
}

// NewDeleteHandler creates a new DeleteHandler
func NewDeleteHandler(keptnHandler *keptnv2.Keptn, configServiceURL string) DeleteHandler {
	return DeleteHandler{
		Handler: NewHandlerBase(keptnHandler, configServiceURL),
	}
}

// HandleEvent takes the sh.keptn.internal.event.service.delete event and deletes the service in all stages
func (h DeleteHandler) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) error {

	defer closeLogger(h.GetKeptnHandler())

	serviceDeleteEvent := keptn.ServiceDeleteEventData{}

	err := ce.DataAs(&serviceDeleteEvent)
	if err != nil {
		errMsg := "service.delete event not well-formed: " + err.Error()
		h.GetKeptnHandler().Logger.Error(errMsg)
		return errors.New(errMsg)
	}

	stageHandler := configutils.NewStageHandler(h.GetConfigServiceURL())
	stages, err := stageHandler.GetAllStages(serviceDeleteEvent.Project)
	if err != nil {
		h.GetKeptnHandler().Logger.Error("Error when getting all stages: " + err.Error())
		return err
	}

	allReleasesSuccessfullyUnistalled := true
	for _, stage := range stages {
		h.GetKeptnHandler().Logger.Info(fmt.Sprintf("Uninstalling Helm releases for service %s in "+
			"stage %s and project %s", serviceDeleteEvent.Service, stage.StageName, serviceDeleteEvent.Project))

		namespace := serviceDeleteEvent.Project + "-" + stage.StageName
		releaseName := namespace + "-" + serviceDeleteEvent.Service
		if err := h.GetHelmExecutor().UninstallRelease(releaseName, namespace); err != nil {
			h.GetKeptnHandler().Logger.Error(err.Error())
			allReleasesSuccessfullyUnistalled = false
		}
		if err := h.GetHelmExecutor().UninstallRelease(releaseName+"-generated", namespace); err != nil {
			h.GetKeptnHandler().Logger.Error(err.Error())
			allReleasesSuccessfullyUnistalled = false
		}
	}

	if allReleasesSuccessfullyUnistalled {
		h.GetKeptnHandler().Logger.Info(fmt.Sprintf("All Helm releases for service %s in project %s successfully uninstalled",
			serviceDeleteEvent.Service, serviceDeleteEvent.Project))
	}

	return nil
}
