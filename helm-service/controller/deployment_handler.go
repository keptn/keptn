package controller

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/configuration_changer"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	"helm.sh/helm/v3/pkg/chart"
)

// DeploymentHandler is a handler for doing the deployment and
// optionally first change the configuration
type DeploymentHandler struct {
	Handler
	mesh                  mesh.Mesh
	generatedChartHandler *helm.GeneratedChartGenerator
}

// NewDeploymentHandler creates a new DeploymentHandler
func NewDeploymentHandler(keptnHandler *keptnv2.Keptn, mesh mesh.Mesh, configServiceURL string) DeploymentHandler {
	generatedChartHandler := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	return DeploymentHandler{
		Handler:           NewHandlerBase(keptnHandler, configServiceURL),
		mesh:                  mesh,
		generatedChartHandler: generatedChartHandler,
	}
}

// HandleEvent handles deployment.triggered events by first changing the new configuration and
// afterwards applying the configuration in the cluster
func (h DeploymentHandler) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) {

	defer closeLogger(h.GetKeptnHandler())

	e := keptnv2.DeploymentTriggeredEventData{}
	if err := ce.DataAs(&e); err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Send deployment started event
	if err := h.SendEvent(ce.ID(), keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName), h.getStartedEventData(e.EventData)); err != nil {
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	var userChart *chart.Chart
	var err error
	gitVersion := ""
	if len(e.ConfigurationChange.Values) > 0 {
		valuesUpdater := configuration_changer.NewValuesUpdater(e.ConfigurationChange.Values)
		userChart, gitVersion, err = configuration_changer.NewConfigurationChanger(h.GetConfigServiceURL()).UpdateChart(e.EventData,
			false, valuesUpdater)
		if err != nil {
			err = fmt.Errorf("failed to update values: %v", err)
			h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
			return
		}
	} else {
		// Read chart
		// TODO set gitVersion
		userChart, err = h.GetUserChart(e.EventData)
		if err != nil {
			err = fmt.Errorf("failed to load chart: %v", err)
			h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
			return
		}
	}

	deploymentStrategy, err := keptnevents.GetDeploymentStrategy(e.Deployment.DeploymentStrategy)
	if err != nil {
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Upgrade user chart
	if err := h.upgradeChart(userChart, e.EventData, deploymentStrategy); err != nil {
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	genChart, err := h.catchupGeneratedChartOnboarding(deploymentStrategy, e.EventData)
	if err != nil {
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Upgrade generated chart
	if err := h.upgradeChart(genChart, e.EventData, deploymentStrategy); err != nil {
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Send finished event
	data := h.getFinishedEventDataForSuccess(e.EventData, gitVersion,
		getDeploymentName(deploymentStrategy, false), deploymentStrategy)
	if err := h.SendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), data); err != nil {
		h.HandleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}
}

// catchupGeneratedChartOnboarding checks if generated chart already exists and if not, it onboards the chart
func (h DeploymentHandler) catchupGeneratedChartOnboarding(deploymentStrategy keptnevents.DeploymentStrategy,
	event keptnv2.EventData) (*chart.Chart, error) {

	exists, err := h.ExistsGeneratedChart(event)
	if err != nil {
		return nil, err
	}

	if exists {
		return h.GetGeneratedChart(event)
	} else {
		// Chart does not exist yet, onboard it now
		userChartManifest, err := h.GetHelmExecutor().GetManifest(helm.GetReleaseName(event.Project, event.Stage, event.Service, false),
			event.Project+"-"+event.Stage)
		if err != nil {
			return nil, err
		}
		onboarder := NewOnboarder(h.GetKeptnHandler(), h.mesh, h.GetConfigServiceURL())
		return onboarder.OnboardGeneratedChart(userChartManifest, event, deploymentStrategy)
	}
}

func (h DeploymentHandler) getStartedEventData(inEventData keptnv2.EventData) keptnv2.DeploymentStartedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.DeploymentStartedEventData{EventData: inEventData}
}

func (h DeploymentHandler) getFinishedEventDataForSuccess(inEventData keptnv2.EventData, gitCommit string,
	deploymentName string, deploymentStrategy keptnevents.DeploymentStrategy) keptnv2.DeploymentFinishedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = keptnv2.ResultPass
	inEventData.Message = "Successfully deployed"
	return keptnv2.DeploymentFinishedEventData{
		EventData: inEventData,
		Deployment: keptnv2.DeploymentData{
			DeploymentStrategy:   deploymentStrategy.String(),
			DeploymentURIsPublic: mesh.GetPublicDeploymentURI(inEventData),
			DeploymentURIsLocal:  mesh.GetLocalDeploymentURI(inEventData),
			DeploymentNames:      []string{deploymentName},
			GitCommit:            gitCommit,
		},
	}
}

func (h DeploymentHandler) getFinishedEventDataForError(eventData keptnv2.EventData, err error) keptnv2.DeploymentFinishedEventData {

	eventData.Status = keptnv2.StatusErrored
	eventData.Result = keptnv2.ResultFailed
	eventData.Message = err.Error()
	return keptnv2.DeploymentFinishedEventData{
		EventData: eventData,
	}
}
