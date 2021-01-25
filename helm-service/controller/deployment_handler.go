package controller

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	"helm.sh/helm/v3/pkg/chart"
	corev1 "k8s.io/api/core/v1"
	"math"
)

// DeploymentHandler is a handler for doing the deployment and
// optionally first change the configuration
type DeploymentHandler struct {
	Handler
	mesh                  mesh.Mesh
	generatedChartHandler helm.ChartGenerator
	onboarder             Onboarder
}

// NewDeploymentHandler creates a new DeploymentHandler
func NewDeploymentHandler(keptnHandler *keptnv2.Keptn, mesh mesh.Mesh, onboarder Onboarder, chartGenerator helm.ChartGenerator, configServiceURL string) *DeploymentHandler {
	return &DeploymentHandler{
		Handler:               NewHandlerBase(keptnHandler, configServiceURL),
		mesh:                  mesh,
		onboarder:             onboarder,
		generatedChartHandler: chartGenerator,
	}
}

// HandleEvent handles deployment.triggered events by first changing the new configuration and
// afterwards applying the configuration in the cluster
func (h *DeploymentHandler) HandleEvent(ce cloudevents.Event) {

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
		userChart, gitVersion, err = h.getUserChart(e.EventData)
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
	data, err := h.getFinishedEventDataForSuccess(e.EventData, gitVersion,
		getDeploymentName(deploymentStrategy, false), deploymentStrategy)
	if err != nil {
		h.handleError(ce.ID(), err, keptnv2.DeploymentTaskName, h.getFinishedEventDataForError(e.EventData, err))
		return
	}
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
		generatedChart, _, err := h.getGeneratedChart(event)
		return generatedChart, err
	}

	// Chart does not exist yet, onboard it now
	userChartManifest, err := h.getHelmExecutor().GetManifest(helm.GetReleaseName(event.Project, event.Stage, event.Service, false),
		event.Project+"-"+event.Stage)
	if err != nil {
		return nil, err
	}

	return h.onboarder.OnboardGeneratedChart(userChartManifest, event, deploymentStrategy)
}

func (h *DeploymentHandler) getDeploymentURIs(e keptnv2.EventData) ([]string, []string, error) {
	userChartManifest, err := h.getHelmExecutor().GetManifest(helm.GetReleaseName(e.Project, e.Stage, e.Service, false),
		e.Project+"-"+e.Stage)

	if err != nil {
		return nil, nil, err
	}
	services := helm.GetServices(userChartManifest)
	if len(services) > 0 {
		if len(services[0].Spec.Ports) > 0 {
			lowestPort, foundPort := getPortOfService(services[0])
			if foundPort {
				localDeploymentURI := mesh.GetLocalDeploymentURI(e, fmt.Sprintf("%d", lowestPort))
				publicDeploymentURI := mesh.GetPublicDeploymentURI(e)
				return localDeploymentURI, publicDeploymentURI, nil
			}
			return nil, nil, errors.New("deployed service does not contain a valid port definition")

		}
	}
	return nil, nil, nil
}

func getPortOfService(service *corev1.Service) (int32, bool) {
	lowestPort := int32(math.MaxInt32)
	foundPort := false
	for _, port := range service.Spec.Ports {
		if port.Protocol == corev1.ProtocolTCP && port.Port < lowestPort {
			lowestPort = port.Port
			foundPort = true
		}
	}
	// if no port explicitly marked as TCP port could be found, take the lowest port number of all specified ports
	if !foundPort {
		for _, port := range service.Spec.Ports {
			if port.Port < lowestPort {
				lowestPort = port.Port
				foundPort = true
			}
		}
	}
	return lowestPort, foundPort
}

func (h *DeploymentHandler) getStartedEventData(inEventData keptnv2.EventData) keptnv2.DeploymentStartedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.DeploymentStartedEventData{EventData: inEventData}
}

func (h *DeploymentHandler) getFinishedEventDataForSuccess(inEventData keptnv2.EventData, gitCommit string,
	deploymentName string, deploymentStrategy keptnevents.DeploymentStrategy) (*keptnv2.DeploymentFinishedEventData, error) {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = keptnv2.ResultPass
	inEventData.Message = "Successfully deployed"

	localURIs, publicURIs, err := h.getDeploymentURIs(inEventData)
	if err != nil {
		return nil, fmt.Errorf("could not determine deployment URIs: %s", err.Error())
	}
	return &keptnv2.DeploymentFinishedEventData{
		EventData: inEventData,
		Deployment: keptnv2.DeploymentData{
			DeploymentStrategy:   deploymentStrategy.String(),
			DeploymentURIsPublic: publicURIs,
			DeploymentURIsLocal:  localURIs,
			DeploymentNames:      []string{deploymentName},
			GitCommit:            gitCommit,
		},
	}, nil
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
