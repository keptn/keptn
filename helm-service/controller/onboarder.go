package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"os"
	"time"

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
func NewOnboarder(mesh mesh.Mesh, keptnHandler *keptnv2.Keptn, configServiceURL string) *Onboarder {
	return &Onboarder{
		mesh:             mesh,
		keptnHandler:     keptnHandler,
		configServiceURL: configServiceURL,
	}
}

// DoOnboard onboards a new service
func (o *Onboarder) DoOnboard(ce cloudevents.Event, loggingDone chan bool) error {

	event := &keptnevents.ServiceCreateEventData{}
	if err := ce.DataAs(event); err != nil {
		o.keptnHandler.Logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	// Check whether Helm chart is provided
	if len(event.HelmChart) == 0 {
		// Event does not contain a Helm chart
		return nil
	}

	// Only close logger/websocket, if there is a chart which needs to be onboarded
	defer func() { loggingDone <- true }()

	keptnHandler, err := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	if err != nil {
		o.keptnHandler.Logger.Error("Could not initialize Keptn handler: " + err.Error())
		return err
	}

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

	// Get deployment strategy
	if _, ok := event.DeploymentStrategies["*"]; ok {
		// Uses the provided deployment strategy for ALL stages
		deplStrategies, err := fixDeploymentStrategies(keptnHandler, event.DeploymentStrategies["*"])
		if err != nil {
			o.keptnHandler.Logger.Error(err.Error())
			return err
		}
		event.DeploymentStrategies = deplStrategies
	} else if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" && len(event.DeploymentStrategies) == 0 {
		deplStrategies, err := getDeploymentStrategies(keptnHandler)
		if err != nil {
			o.keptnHandler.Logger.Error(err.Error())
			return err
		}
		event.DeploymentStrategies = deplStrategies
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

	// Initialize Umbrella chart
	umbrellaChartHandler := helm.NewUmbrellaChartHandler(o.configServiceURL)
	isUmbrellaChartAvailable, err := umbrellaChartHandler.IsUmbrellaChartAvailableInAllStages(event.Project, stages)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when getting Helm Chart for stages. " + err.Error())
		return err
	}
	if !isUmbrellaChartAvailable {
		o.keptnHandler.Logger.Info(fmt.Sprintf("Create umbrella Helm Chart for project %s", event.Project))
		// Initialize the umbrella chart
		if err := umbrellaChartHandler.InitUmbrellaChart(event, stages); err != nil {
			return fmt.Errorf("Error when initializing the umbrella chart for project %s: %s", event.Project, err.Error())
		}
	}

	// Onboard service in all namespaces
	for _, stage := range stages {
		if err := o.onboardService(stage.StageName, event); err != nil {
			o.keptnHandler.Logger.Error(err.Error())
			return err
		}
		if event.DeploymentStrategies[stage.StageName] == keptnevents.Duplicate {
			// inject Istio to the namespace for blue-green deployments
			if err := namespaceMng.InjectIstio(event.Project, stage.StageName); err != nil {
				o.keptnHandler.Logger.Error(err.Error())
				return err
			}
		}
	}

	o.keptnHandler.Logger.Info(fmt.Sprintf("Finished creating service %s in project %s", event.Service, event.Project))
	return nil
}

func (o *Onboarder) checkAndSetServiceName(event *keptnevents.ServiceCreateEventData) error {

	errorMsg := "Service name contains upper case letter(s) or special character(s).\n " +
		"Keptn relies on the following conventions: " +
		"start with a lower case letter, then lower case letters, numbers, and hyphens are allowed."

	helmChartData, err := base64.StdEncoding.DecodeString(event.HelmChart)
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

func (o *Onboarder) onboardService(stageName string, event *keptnevents.ServiceCreateEventData) error {

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

	helmChartData, err := base64.StdEncoding.DecodeString(event.HelmChart)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when decoding the Helm Chart")
		return err
	}

	o.keptnHandler.Logger.Debug("Storing the Helm Chart provided by the user in stage " + stageName)
	if err := keptnutils.StoreChart(event.Project, event.Service, stageName, helm.GetChartName(event.Service, false),
		helmChartData, o.configServiceURL); err != nil {
		o.keptnHandler.Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return err
	}

	if err := o.updateUmbrellaChart(event.Project, stageName, helm.GetChartName(event.Service, false)); err != nil {
		return err
	}

	chartGenerator := helm.NewGeneratedChartHandler(o.mesh, o.keptnHandler.Logger)
	o.keptnHandler.Logger.Debug(fmt.Sprintf("For stage %s with deployment strategy %s, an empty chart is generated", stageName, event.DeploymentStrategies[stageName].String()))
	generatedChart := chartGenerator.GenerateEmptyChart(event.Service, event.DeploymentStrategies[stageName])

	helmChartName := helm.GetChartName(event.Service, true)
	o.keptnHandler.Logger.Debug(fmt.Sprintf("Storing the Keptn-generated Helm Chart %s for stage %s", helmChartName, stageName))

	generatedChartData, err := keptnutils.PackageChart(generatedChart)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when packing the managed chart: " + err.Error())
		return err
	}

	if err := keptnutils.StoreChart(event.Project, event.Service, stageName, helmChartName,
		generatedChartData, o.configServiceURL); err != nil {
		o.keptnHandler.Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return err
	}
	return o.updateUmbrellaChart(event.Project, stageName, helmChartName)
}

// IsGeneratedChartEmpty checks whether the generated chart is empty
func (c *Onboarder) IsGeneratedChartEmpty(chart *chart.Chart) bool {

	return len(chart.Templates) == 0
}

func (o *Onboarder) OnboardGeneratedService(helmManifest string, project string, stageName string,
	service string, strategy keptnevents.DeploymentStrategy) (*chart.Chart, error) {

	chartGenerator := helm.NewGeneratedChartHandler(o.mesh, o.keptnHandler.Logger)

	helmChartName := helm.GetChartName(service, true)
	o.keptnHandler.Logger.Debug(fmt.Sprintf("Generating the Keptn-managed Helm Chart %s for stage %s", helmChartName, stageName))

	var generatedChart *chart.Chart
	var err error
	if strategy == keptnevents.Duplicate {
		o.keptnHandler.Logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, "+
			"a chart for a duplicate deployment strategy is generated", service, stageName, strategy.String()))
		generatedChart, err = chartGenerator.GenerateDuplicateManagedChart(helmManifest, project, stageName, service)
		if err != nil {
			o.keptnHandler.Logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	} else {
		o.keptnHandler.Logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, a mesh chart is generated",
			service, stageName, strategy.String()))
		generatedChart, err = chartGenerator.GenerateMeshChart(helmManifest, project, stageName, service)
		if err != nil {
			o.keptnHandler.Logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	}

	o.keptnHandler.Logger.Debug(fmt.Sprintf("Storing the Keptn-generated Helm Chart %s for stage %s", helmChartName, stageName))
	generatedChartData, err := keptnutils.PackageChart(generatedChart)
	if err != nil {
		o.keptnHandler.Logger.Error("Error when packing the managed chart: " + err.Error())
		return nil, err
	}

	if err := keptnutils.StoreChart(project, service, stageName, helmChartName,
		generatedChartData, o.configServiceURL); err != nil {
		o.keptnHandler.Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return nil, err
	}
	return generatedChart, nil
}

func (o *Onboarder) updateUmbrellaChart(project, stage, helmChartName string) error {

	umbrellaChartHandler := helm.NewUmbrellaChartHandler(o.configServiceURL)
	o.keptnHandler.Logger.Debug(fmt.Sprintf("Updating the umbrella Helm Chart with the new Helm Chart %s in stage %s", helmChartName, stage))
	// if err := helm.AddChartInUmbrellaRequirements(event.Project, helmChartName, stage, url.String()); err != nil {
	// 	o.keptnHandler.Logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
	// 	return err
	// }
	if err := umbrellaChartHandler.AddChartInUmbrellaValues(project, helmChartName, stage); err != nil {
		o.keptnHandler.Logger.Error("Error when adding the Helm Chart in the umbrella values file: " + err.Error())
		return err
	}
	return nil
}
