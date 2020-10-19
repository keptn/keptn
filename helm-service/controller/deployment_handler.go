package controller

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
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
	onboarder             Onboarder
}

// NewDeploymentHandler creates a new DeploymentHandler
func NewDeploymentHandler(keptnHandler *keptnv2.Keptn, mesh mesh.Mesh, onboarder Onboarder, configServiceURL string) *DeploymentHandler {
	generatedChartHandler := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	return &DeploymentHandler{
		Handler:               NewHandlerBase(keptnHandler, configServiceURL),
		mesh:                  mesh,
		onboarder:             onboarder,
		generatedChartHandler: generatedChartHandler,
	}
}

// HandleEvent handles deployment.triggered events by first changing the new configuration and
// afterwards applying the configuration in the cluster
func (h *DeploymentHandler) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) {

	defer closeLogger(h.getKeptnHandler())

	e := keptnv2.DeploymentTriggeredEventData{}
	if err := ce.DataAs(&e); err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Send deployment started event
	h.getKeptnHandler().Logger.Info(fmt.Sprintf("Starting deployment for service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
	if err := h.sendEvent(ce.ID(), keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName), h.getStartedEventData(e.EventData)); err != nil {
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	if e.Result == keptnv2.ResultFailed {
		h.getKeptnHandler().Logger.Info(fmt.Sprintf("No deployment done for service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
		data := h.getFinishedEventDataForNoDeployment(e.EventData)
		if err := h.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), data); err != nil {
			h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
			return
		}
		return
	}

	var userChart *chart.Chart
	var err error
	gitVersion := ""
	if len(e.ConfigurationChange.Values) > 0 {
		h.getKeptnHandler().Logger.Info(fmt.Sprintf("Updating values for service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
		valuesUpdater := configurationchanger.NewValuesManipulator(e.ConfigurationChange.Values)
		userChart, gitVersion, err = configurationchanger.NewConfigurationChanger(h.getConfigServiceURL()).UpdateChart(e.EventData,
			false, valuesUpdater)
		if err != nil {
			err = fmt.Errorf("failed to update values: %v", err)
			h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
			return
		}
	} else {
		// Read chart
		// TODO set gitVersion
		userChart, err = h.getUserChart(e.EventData)
		if err != nil {
			err = fmt.Errorf("failed to load chart: %v", err)
			h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
			return
		}
	}

	deploymentStrategy, err := keptnevents.GetDeploymentStrategy(e.Deployment.DeploymentStrategy)
	if err != nil {
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Upgrade user chart
	if err := h.upgradeChart(userChart, e.EventData, deploymentStrategy); err != nil {
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	if err := h.upgradeGeneratedChart(deploymentStrategy, e); err != nil {
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Send finished event
	data := h.getFinishedEventDataForSuccess(e.EventData, gitVersion,
		getDeploymentName(deploymentStrategy, false), deploymentStrategy)
	if err := h.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), data); err != nil {
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}
	h.getKeptnHandler().Logger.Info(fmt.Sprintf("Deployment finished for service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
}

func (h *DeploymentHandler) upgradeGeneratedChart(deploymentStrategy keptnevents.DeploymentStrategy, e keptnv2.DeploymentTriggeredEventData) error {

	genChart, err := h.catchupGeneratedChartOnboarding(deploymentStrategy, e.EventData)
	if err != nil {
		return err
	}

	if deploymentStrategy == keptnevents.Duplicate {
		// Route the traffic to the user-chart
		weightUpdater := configurationchanger.NewCanaryWeightManipulator(h.mesh, 100)
		genChart, _, err = configurationchanger.NewConfigurationChanger(h.getConfigServiceURL()).UpdateLoadedChart(genChart, e.EventData,
			true, weightUpdater)
		if err != nil {
			return fmt.Errorf("failed to update canary weight: %v", err)
		}
	}

	// Upgrade generated chart
	return h.upgradeChart(genChart, e.EventData, deploymentStrategy)
}

// catchupGeneratedChartOnboarding checks if generated chart already exists and if not, it onboards the chart
func (h *DeploymentHandler) catchupGeneratedChartOnboarding(deploymentStrategy keptnevents.DeploymentStrategy,
	event keptnv2.EventData) (*chart.Chart, error) {

	exists, err := h.existsGeneratedChart(event)
	if err != nil {
		return nil, err
	}

	if exists {
		return h.getGeneratedChart(event)
	}

	// Chart does not exist yet, onboard it now
	userChartManifest, err := h.getHelmExecutor().GetManifest(helm.GetReleaseName(event.Project, event.Stage, event.Service, false),
		event.Project+"-"+event.Stage)
	if err != nil {
		return nil, err
	}
	//onboarder := NewOnboarder(h.getKeptnHandler(), h.mesh, h.getConfigServiceURL())
	return h.onboarder.OnboardGeneratedChart(userChartManifest, event, deploymentStrategy)
}

func (h *DeploymentHandler) getStartedEventData(inEventData keptnv2.EventData) keptnv2.DeploymentStartedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.DeploymentStartedEventData{EventData: inEventData}
}

func (h *DeploymentHandler) getFinishedEventDataForSuccess(inEventData keptnv2.EventData, gitCommit string,
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

func (h *DeploymentHandler) getFinishedEventDataForError(eventData keptnv2.EventData, err error) keptnv2.DeploymentFinishedEventData {

	eventData.Status = keptnv2.StatusErrored
	eventData.Result = keptnv2.ResultFailed
	eventData.Message = err.Error()
	return keptnv2.DeploymentFinishedEventData{
		EventData: eventData,
	}
}

func (h *DeploymentHandler) getFinishedEventDataForNoDeployment(eventData keptnv2.EventData) keptnv2.DeploymentFinishedEventData {

	eventData.Status = keptnv2.StatusSucceeded
	eventData.Result = keptnv2.ResultFailed
	eventData.Message = "No deployment has been executed"
	return keptnv2.DeploymentFinishedEventData{
		EventData: eventData,
	}
}
