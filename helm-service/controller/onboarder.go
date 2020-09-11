package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/go-utils/pkg/api/models"

	"helm.sh/helm/v3/pkg/chart"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
)

// Onboarder is a container of variables required for onboarding a new service
type Onboarder struct {
	HandlerBase
	mesh mesh.Mesh
}

// NewOnboarder creates a new Onboarder
func NewOnboarder(keptnHandler *keptnv2.Keptn, mesh mesh.Mesh, configServiceURL string) Onboarder {
	return Onboarder{
		HandlerBase: NewHandlerBase(keptnHandler, configServiceURL),
		mesh:        mesh,
	}
}

// HandleEvent onboards a new service
func (o Onboarder) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) {

	e := &keptnv2.ServiceCreateTriggeredEventData{}
	if err := ce.DataAs(e); err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Check whether Helm chart is provided
	if len(e.Helm.Chart) == 0 {
		// Event does not contain a Helm chart
		return
	}

	// Send service.create started eveny
	if err := o.SendEvent(ce.ID(), keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName),
		o.getStartedEventData(e.EventData)); err != nil {
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Only close logger/websocket, if there is a chart which needs to be onboarded
	defer closeLogger(o.keptnHandler)

	// Check if project exists
	projHandler := configutils.NewProjectHandler(o.configServiceURL)
	if _, err := projHandler.GetProject(models.Project{ProjectName: e.Project}); err != nil {
		err := fmt.Errorf("failed not retrieve project %s: %s", e.Project, *err.Message)
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Check service name
	if err := o.checkAndSetServiceName(e); err != nil {
		err := fmt.Errorf("invalid service name: %s", err.Error())
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Check stages
	stages, err := o.getStages(e)
	if err != nil {
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Initialize Namespace
	namespaceMng := NewNamespaceManager(o.keptnHandler.Logger)
	if err := namespaceMng.InitNamespaces(e.Project, stages); err != nil {
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Onboard service in all namespaces
	for _, stage := range stages {
		if err := o.onboardService(stage, e); err != nil {
			o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
			return
		}
	}

	// Send finished event
	msg := fmt.Sprintf("Finished creating service %s in project %s", e.Service, e.Project)
	// TODO Set git version
	data := o.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded, keptnv2.ResultPass, msg, "")
	if err := o.SendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), data); err != nil {
		o.HandleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}
}

// getStages returns a list of stages where the service should be onboarded
// If the stage of the incoming event is empty, all available stages are returned
func (o Onboarder) getStages(e *keptnv2.ServiceCreateTriggeredEventData) ([]string, error) {
	stageHandler := configutils.NewStageHandler(o.configServiceURL)
	allStages, err := stageHandler.GetAllStages(e.Project)
	if err != nil {
		return nil, fmt.Errorf("failed to retriev stages: %v", err.Error())
	}
	var stages []string = nil
	for _, availableStage := range allStages {
		if availableStage.StageName == e.Stage || e.Stage == "" {
			stages = append(stages, availableStage.StageName)
		}
	}

	if len(stages) == 0 {
		return nil, errors.New("Cannot onboard service because no stage is available")
	}
	return stages, nil
}

func (o Onboarder) checkAndSetServiceName(event *keptnv2.ServiceCreateTriggeredEventData) error {

	errorMsg := "Service name contains upper case letter(s) or special character(s).\n " +
		"Keptn relies on the following conventions: " +
		"start with a lower case letter, then lower case letters, numbers, and hyphens are allowed."

	helmChartData, err := base64.StdEncoding.DecodeString(event.Helm.Chart)
	if err != nil {
		return fmt.Errorf("Error when decoding the Helm Chart: %v", err)
	}
	ch, err := keptnutils.LoadChart(helmChartData)
	if err != nil {
		return fmt.Errorf("Error when loading Helm Chart: %v", err)
	}
	services, err := keptnutils.GetRenderedServices(ch)
	if err != nil {
		return fmt.Errorf("Error when rendering services: %v", err)
	}
	if len(services) != 1 {
		return fmt.Errorf("Helm Chart has to contain exactly one Kubernetes service, but it contains %d services", len(services))
	}
	k8sServiceName := services[0].Name
	if !keptncommon.ValidateKeptnEntityName(k8sServiceName) {
		return errors.New(errorMsg)
	}
	if event.Service == "" {
		// Set service name in event
		event.Service = k8sServiceName
	}
	if k8sServiceName != event.Service {
		return fmt.Errorf("Provided Keptn service name \"%s\" "+
			"does not match Kubernetes service name \"%s\"", event.Service, k8sServiceName)
	}
	return nil
}

func (o Onboarder) onboardService(stageName string, event *keptnv2.ServiceCreateTriggeredEventData) error {

	serviceHandler := configutils.NewServiceHandler(o.configServiceURL)
	const retries = 2
	var err error
	for i := 0; i < retries; i++ {
		_, err = serviceHandler.GetService(event.Project, stageName, event.Service)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return err
	}

	helmChartData, err := base64.StdEncoding.DecodeString(event.Helm.Chart)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when decoding the Helm Chart")
		return err
	}

	o.keptnHandler.Logger.Debug("Storing the Helm Chart provided by the user in stage " + stageName)
	if _, err := keptnutils.StoreChart(event.Project, event.Service, stageName, helm.GetChartName(event.Service, false),
		helmChartData, o.configServiceURL); err != nil {
		o.keptnHandler.Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return err
	}
	return nil
}

func (o Onboarder) OnboardGeneratedChart(helmManifest string, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy) (*chart.Chart, error) {

	chartGenerator := helm.NewGeneratedChartGenerator(o.mesh, o.keptnHandler.Logger)

	helmChartName := helm.GetChartName(event.Service, true)
	o.keptnHandler.Logger.Debug(fmt.Sprintf("Generating the Keptn-managed Helm Chart %s for stage %s", helmChartName, event.Stage))

	var generatedChart *chart.Chart
	var err error
	if strategy == keptnevents.Duplicate {
		o.keptnHandler.Logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, "+
			"a chart for a duplicate deployment strategy is generated", event.Service, event.Stage, strategy.String()))
		generatedChart, err = chartGenerator.GenerateDuplicateChart(helmManifest, event.Project, event.Stage, event.Service)
		if err != nil {
			o.keptnHandler.Logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
		// inject Istio to the namespace for blue-green deployments
		namespaceMng := NewNamespaceManager(o.keptnHandler.Logger)
		if err := namespaceMng.InjectIstio(event.Project, event.Stage); err != nil {
			return nil, err
		}
	} else {
		o.keptnHandler.Logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, a mesh chart is generated",
			event.Service, event.Stage, strategy.String()))
		generatedChart, err = chartGenerator.GenerateMeshChart(helmManifest, event.Project, event.Stage, event.Service)
		if err != nil {
			o.keptnHandler.Logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	}

	o.keptnHandler.Logger.Debug(fmt.Sprintf("Storing the Keptn-generated Helm Chart %s for stage %s", helmChartName, event.Stage))
	generatedChartData, err := keptnutils.PackageChart(generatedChart)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when packing the managed chart: " + err.Error())
		return nil, err
	}

	if _, err := keptnutils.StoreChart(event.Project, event.Service, event.Stage, helmChartName,
		generatedChartData, o.configServiceURL); err != nil {
		o.keptnHandler.Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return nil, err
	}
	return generatedChart, nil
}

func (o Onboarder) getStartedEventData(inEventData keptnv2.EventData) keptnv2.ServiceCreateStartedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.ServiceCreateStartedEventData{EventData: inEventData}
}

func (o Onboarder) getFinishedEventData(inEventData keptnv2.EventData, status keptnv2.StatusType, result keptnv2.ResultType,
	message string, gitCommit string) keptnv2.ServiceCreateFinishedEventData {

	inEventData.Status = status
	inEventData.Result = result
	inEventData.Message = message

	return keptnv2.ServiceCreateFinishedEventData{
		EventData: inEventData,
		Helm:      keptnv2.HelmData{GitCommit: gitCommit},
	}
}

func (o Onboarder) getFinishedEventDataForError(inEventData keptnv2.EventData, err error) keptnv2.ServiceCreateFinishedEventData {
	return o.getFinishedEventData(inEventData, keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error(), "")
}
