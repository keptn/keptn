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
	mesh             mesh.Mesh
	keptnHandler     *keptnv2.Keptn
	configServiceURL string
}

// NewOnboarder creates a new Onboarder
func NewOnboarder(keptnHandler *keptnv2.Keptn, mesh mesh.Mesh, configServiceURL string) *Onboarder {
	return &Onboarder{
		mesh:             mesh,
		keptnHandler:     keptnHandler,
		configServiceURL: configServiceURL,
	}
}

// HandleEvent onboards a new service
func (o *Onboarder) HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn)) error {

	event := &keptnv2.ServiceCreateTriggeredEventData{}
	if err := ce.DataAs(event); err != nil {
		o.keptnHandler.Logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	// Check whether Helm chart is provided
	if len(event.Helm.Chart) == 0 {
		// Event does not contain a Helm chart
		return nil
	}

	// Only close logger/websocket, if there is a chart which needs to be onboarded
	defer closeLogger(o.keptnHandler)

	// Check if project exists
	projHandler := configutils.NewProjectHandler(o.configServiceURL)
	if _, err := projHandler.GetProject(models.Project{ProjectName: event.Project}); err != nil {
		o.keptnHandler.Logger.Error(fmt.Sprintf("Could not retrieve project %s: %s", event.Project, *err.Message))
		return errors.New(*err.Message)
	}

	// Check service name
	if err := o.checkAndSetServiceName(event); err != nil {
		o.keptnHandler.Logger.Error(fmt.Sprintf("Invalid service name: %s", err.Error()))
		return err
	}

	// Check stages
	stageHandler := configutils.NewStageHandler(o.configServiceURL)
	stages, err := stageHandler.GetAllStages(event.Project)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when getting all stages: " + err.Error())
		return err
	}
	if len(stages) == 0 {
		o.keptnHandler.Logger.Error("Cannot onboard service because no stage is available")
		return errors.New("Cannot onboard service because no stage is available")
	}

	// Initialize Namespace
	namespaceMng := NewNamespaceManager(o.keptnHandler.Logger)
	if err := namespaceMng.InitNamespaces(event.Project, stages); err != nil {
		o.keptnHandler.Logger.Error(err.Error())
		return err
	}

	// Onboard service in all namespaces
	for _, stage := range stages {
		if err := o.onboardService(stage.StageName, event); err != nil {
			o.keptnHandler.Logger.Error(err.Error())
			return err
		}
	}

	o.keptnHandler.Logger.Info(fmt.Sprintf("Finished creating service %s in project %s", event.Service, event.Project))
	return nil
}

func (o *Onboarder) checkAndSetServiceName(event *keptnv2.ServiceCreateTriggeredEventData) error {

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

func (o *Onboarder) onboardService(stageName string, event *keptnv2.ServiceCreateTriggeredEventData) error {

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

func (o *Onboarder) OnboardGeneratedChart(helmManifest string, event keptnv2.EventData, strategy keptnevents.DeploymentStrategy) (*chart.Chart, error) {

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
