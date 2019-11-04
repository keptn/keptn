package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"

	"github.com/keptn/keptn/helm-service/controller/mesh"
)

// Onboarder is a container of variables required for onboarding a new service
type Onboarder struct {
	mesh           mesh.Mesh
	logger         keptnutils.LoggerInterface
	canaryLevelGen helm.CanaryLevelGenerator
	keptnDomain    string
}

// NewOnboarder creates a new Onboarder
func NewOnboarder(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger keptnutils.LoggerInterface, keptnDomain string) *Onboarder {
	return &Onboarder{mesh: mesh, canaryLevelGen: canaryLevelGen, logger: logger, keptnDomain: keptnDomain}
}

// DoOnboard onboards a new service
func (o *Onboarder) DoOnboard(ce cloudevents.Event, loggingDone chan bool) error {

	defer func() { loggingDone <- true }()

	event := &keptnevents.ServiceCreateEventData{}
	if err := ce.DataAs(event); err != nil {
		o.logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if _, ok := event.DeploymentStrategies["*"]; ok {
		deplStrategies, err := FixDeploymentStrategies(event.Project, event.DeploymentStrategies["*"])
		if err != nil {
			o.logger.Error(fmt.Sprintf("Error when getting deployment strategies: %s" + err.Error()))
			return err
		}
		event.DeploymentStrategies = deplStrategies
	} else if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" && (event.DeploymentStrategies == nil || len(event.DeploymentStrategies) == 0) {
		deplStrategies, err := GetDeploymentStrategies(event.Project)
		if err != nil {
			o.logger.Error(fmt.Sprintf("Error when getting deployment strategies: %s" + err.Error()))
			return err
		}
		event.DeploymentStrategies = deplStrategies
	}

	o.logger.Info(fmt.Sprintf("Start creating service %s in project %s", event.Service, event.Project))

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		o.logger.Error(fmt.Sprintf("Error when getting config service url: %s", err.Error()))
		return err
	}

	stageHandler := keptnutils.NewStageHandler(url.String())
	stages, err := stageHandler.GetAllStages(event.Project)
	if err != nil {
		o.logger.Error("Error when getting all stages: " + err.Error())
		return err
	}

	firstService, err := o.isFirstServiceOfProject(event, stages)
	if err != nil {
		o.logger.Error("Error when checking whether any service was created before: " + err.Error())
		return err
	}
	if firstService {
		o.logger.Info("Create Helm umbrella charts")
		umbrellaChartHandler := helm.NewUmbrellaChartHandler(o.mesh)
		if err := o.initAndApplyUmbrellaChart(event, umbrellaChartHandler, stages); err != nil {
			o.logger.Error(fmt.Sprintf("Error when initalizing and applying umbrella charts for project %s: %s", event.Project, err.Error()))
			return err
		}
	}

	for _, stage := range stages {
		if err := o.onboardService(stage.StageName, event, url.String()); err != nil {
			return err
		}
	}

	o.logger.Info(fmt.Sprintf("Finished creating service %s in project %s", event.Service, event.Project))
	return nil
}

func (o *Onboarder) onboardService(stageName string, event *keptnevents.ServiceCreateEventData,
	configServiceURL string) error {

	serviceHandler := keptnutils.NewServiceHandler(configServiceURL)
	helmChartData, err := base64.StdEncoding.DecodeString(event.HelmChart)
	if err != nil {
		o.logger.Error("Error when decoding the Helm chart")
	}

	o.logger.Debug("Creating new keptn service " + event.Service + " in stage " + stageName)
	respErr, err := serviceHandler.CreateService(event.Project, stageName, event.Service)
	if respErr != nil {
		return errors.New(*respErr.Message)
	}
	if err != nil {
		return err
	}

	o.logger.Debug("Storing the Helm chart provided by the user in stage " + stageName)
	if err := keptnutils.StoreChart(event.Project, event.Service, stageName, helm.GetChartName(event.Service, false),
		helmChartData, configServiceURL); err != nil {
		o.logger.Error("Error when storing the Helm chart: " + err.Error())
		return err
	}

	if err := o.updateUmbrellaChart(event.Project, stageName, helm.GetChartName(event.Service, false)); err != nil {
		return err
	}

	chartGenerator := helm.NewGeneratedChartHandler(o.mesh, o.canaryLevelGen, o.keptnDomain)
	o.logger.Debug(fmt.Sprintf("For stage %s with deployment strategy %s, an empty chart is generated", stageName, event.DeploymentStrategies[stageName].String()))
	generatedChart := chartGenerator.GenerateEmptyChart(event.Project, stageName, event.Service, event.DeploymentStrategies[stageName])

	helmChartName := helm.GetChartName(event.Service, true)
	o.logger.Debug(fmt.Sprintf("Storing the keptn generated Helm chart %s for stage %s", helmChartName, stageName))

	generatedChartData, err := keptnutils.PackageChart(generatedChart)
	if err != nil {
		o.logger.Error("Error when packing the managed chart: " + err.Error())
		return err
	}

	if err := keptnutils.StoreChart(event.Project, event.Service, stageName, helmChartName,
		generatedChartData, configServiceURL); err != nil {
		o.logger.Error("Error when storing the Helm chart: " + err.Error())
		return err
	}
	return o.updateUmbrellaChart(event.Project, stageName, helmChartName)
}

// IsGeneratedChartEmpty checks whether the generated chart is empty
func (c *Onboarder) IsGeneratedChartEmpty(chart *chart.Chart) (bool) {

	return len(chart.Templates) == 0
}

func (o *Onboarder) OnboardGeneratedService(helmUpgradeMsg string, project string, stageName string,
	service string, strategy keptnevents.DeploymentStrategy) (*chart.Chart, error) {

	chartGenerator := helm.NewGeneratedChartHandler(o.mesh, o.canaryLevelGen, o.keptnDomain)

	helmChartName := helm.GetChartName(service, true)
	o.logger.Debug(fmt.Sprintf("Generating the keptn-managed Helm chart %s for stage %s", helmChartName, stageName))

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	var generatedChart *chart.Chart
	if strategy == keptnevents.Duplicate {
		o.logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, " +
			"a chart for a duplicate deployment strategy is generated", service, stageName, strategy.String()))
		generatedChart, err = chartGenerator.GenerateDuplicateManagedChart(helmUpgradeMsg, project, stageName, service)
		if err != nil {
			o.logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	} else {
		o.logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, a mesh chart is generated",
			service, stageName, strategy.String()))
		generatedChart, err = chartGenerator.GenerateMeshChart(helmUpgradeMsg, project, stageName, service)
		if err != nil {
			o.logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	}

	o.logger.Debug(fmt.Sprintf("Storing the keptn generated Helm chart %s for stage %s", helmChartName, stageName))
	generatedChartData, err := keptnutils.PackageChart(generatedChart)
	if err != nil {
		o.logger.Error("Error when packing the managed chart: " + err.Error())
		return nil, err
	}

	if err := keptnutils.StoreChart(project, service, stageName, helmChartName,
		generatedChartData, url.String()); err != nil {
		o.logger.Error("Error when storing the Helm chart: " + err.Error())
		return nil, err
	}
	return generatedChart, nil
}

func (o *Onboarder) updateUmbrellaChart(project, stage, helmChartName string) error {

	umbrellaChartHandler := helm.NewUmbrellaChartHandler(o.mesh)
	o.logger.Debug(fmt.Sprintf("Updating the Umbrella chart with the new Helm chart %s in stage %s", helmChartName, stage))
	// if err := helm.AddChartInUmbrellaRequirements(event.Project, helmChartName, stage, url.String()); err != nil {
	// 	o.logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
	// 	return err
	// }
	if err := umbrellaChartHandler.AddChartInUmbrellaValues(project, helmChartName, stage); err != nil {
		o.logger.Error("Error when adding the chart in the Umbrella values file: " + err.Error())
		return err
	}
	return nil
}

func (o *Onboarder) initAndApplyUmbrellaChart(event *keptnevents.ServiceCreateEventData,
	umbrellaChartHandler *helm.UmbrellaChartHandler, stages []*models.Stage) error {

	// Initalize the umbrella chart
	if err := umbrellaChartHandler.InitUmbrellaChart(event, stages); err != nil {
		return fmt.Errorf("Error when initializing the umbrella chart: %s", err.Error())
	}

	for _, stage := range stages {
		// Apply the umbrella chart
		umbrellaChart, err := ioutil.TempDir("", "")
		if err != nil {
			return fmt.Errorf("Error when creating a temporary directory: %s", err.Error())
		}
		if err := umbrellaChartHandler.GetUmbrellaChart(umbrellaChart, event.Project, stage.StageName); err != nil {
			return fmt.Errorf("Error when getting umbrella chart: %s", err)
		}

		configChanger := NewConfigurationChanger(o.mesh, o.canaryLevelGen, o.logger, o.keptnDomain)
		if err := configChanger.ApplyDirectory(umbrellaChart, helm.GetUmbrellaReleaseName(event.Project, stage.StageName),
			helm.GetUmbrellaNamespace(event.Project, stage.StageName)); err != nil {
			return fmt.Errorf("Error when applying umbrella chart in stage %s: %s", stage.StageName, err.Error())
		}
		if err := os.RemoveAll(umbrellaChart); err != nil {
			return err
		}
	}
	return nil
}

func (o *Onboarder) isFirstServiceOfProject(event *keptnevents.ServiceCreateEventData, stages []*models.Stage) (bool, error) {

	if len(stages) == 0 {
		return false, errors.New("Cannot onboard service because no stage is available")
	}
	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return false, err
	}
	svcHandler := keptnutils.NewServiceHandler(url.String())
	// Use any stage for checking whether there is already a service created
	services, err := svcHandler.GetAllServices(event.Project, stages[0].StageName)
	if err != nil {
		return false, err
	}
	return len(services) == 0, nil
}
