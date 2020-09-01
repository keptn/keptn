package controller

import (
	"errors"
	"fmt"

	configutils "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/keptn/keptn/helm-service/controller/helm"

	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

// DeleteHandler handles sh.keptn.internal.event.service.delete events
type DeleteHandler struct {
	keptnHandler     *keptn.Keptn
	helmExecutor     helm.HelmExecutor
	configServiceURL string
}

// NewDeleteHandler creates a new DeleteHandler
func NewDeleteHandler(keptnHandler *keptn.Keptn, configServiceURL string) *DeleteHandler {
	helmExecutor := helm.NewHelmV3Executor(keptnHandler.Logger)
	return &DeleteHandler{keptnHandler: keptnHandler, helmExecutor: helmExecutor, configServiceURL: configServiceURL}
}

// HandleEvent takes the sh.keptn.internal.event.service.delete event and deletes the service in all stages
func (d *DeleteHandler) HandleEvent(ce cloudevents.Event, loggingDone chan bool) error {

	defer func() { loggingDone <- true }()
	serviceDeleteEvent := keptn.ServiceDeleteEventData{}

	err := ce.DataAs(&serviceDeleteEvent)
	if err != nil {
		errMsg := "service.delete event not well-formed: " + err.Error()
		d.keptnHandler.Logger.Error(errMsg)
		return errors.New(errMsg)
	}

	stageHandler := configutils.NewStageHandler(d.configServiceURL)
	stages, err := stageHandler.GetAllStages(serviceDeleteEvent.Project)
	if err != nil {
		d.keptnHandler.Logger.Error("Error when getting all stages: " + err.Error())
		return err
	}

	allReleasesSuccessfullyUnistalled := true
	for _, stage := range stages {
		d.keptnHandler.Logger.Info(fmt.Sprintf("Uninstalling Helm releases for service %s in "+
			"stage %s and project %s", serviceDeleteEvent.Service, stage.StageName, serviceDeleteEvent.Project))

		namespace := serviceDeleteEvent.Project + "-" + stage.StageName
		releaseName := namespace + "-" + serviceDeleteEvent.Service
		if err := d.helmExecutor.UninstallRelease(releaseName, namespace); err != nil {
			d.keptnHandler.Logger.Error(err.Error())
			allReleasesSuccessfullyUnistalled = false
		}
		if err := d.helmExecutor.UninstallRelease(releaseName+"-generated", namespace); err != nil {
			d.keptnHandler.Logger.Error(err.Error())
			allReleasesSuccessfullyUnistalled = false
		}
	}

	if allReleasesSuccessfullyUnistalled {
		d.keptnHandler.Logger.Info(fmt.Sprintf("All Helm releases for service %s in project %s successfully uninstalled",
			serviceDeleteEvent.Service, serviceDeleteEvent.Project))
	}

	return nil
}
