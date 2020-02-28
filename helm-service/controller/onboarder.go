package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"

	cloudevents "github.com/cloudevents/sdk-go"

	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
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

	if err := o.checkAndSetServiceName(event); err != nil {
		o.logger.Error(fmt.Sprintf("Invalid service name: %s", err.Error()))
		return err
	}

	if _, ok := event.DeploymentStrategies["*"]; ok {
		// Uses the provided deployment strategy for ALL stages
		deplStrategies, err := fixDeploymentStrategies(event.Project, event.DeploymentStrategies["*"])
		if err != nil {
			o.logger.Error(fmt.Sprintf("Error when getting deployment strategies: %s" + err.Error()))
			return err
		}
		event.DeploymentStrategies = deplStrategies
	} else if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" && len(event.DeploymentStrategies) == 0 {
		deplStrategies, err := getDeploymentStrategies(event.Project)
		if err != nil {
			o.logger.Error(fmt.Sprintf("Error when getting deployment strategies: %s" + err.Error()))
			return err
		}
		event.DeploymentStrategies = deplStrategies
	}

	o.logger.Debug(fmt.Sprintf("Start creating service %s in project %s", event.Service, event.Project))

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		o.logger.Error(fmt.Sprintf("Error when getting config service url: %s", err.Error()))
		return err
	}

	stageHandler := configutils.NewStageHandler(url.String())
	stages, err := stageHandler.GetAllStages(event.Project)
	if err != nil {
		o.logger.Error("Error when getting all stages: " + err.Error())
		return err
	}

	if len(stages) == 0 {
		o.logger.Error("Cannot onboard service because no stage is available")
		return errors.New("Cannot onboard service because no stage is available")
	}

	if event.HelmChart != "" {
		umbrellaChartHandler := helm.NewUmbrellaChartHandler(url.String())
		isUmbrellaChartAvailable, err := umbrellaChartHandler.IsUmbrellaChartAvailableInAllStages(event.Project, stages)
		if err != nil {
			o.logger.Error("Error when getting Helm chart for stages. " + err.Error())
			return err
		}
		if !isUmbrellaChartAvailable {
			o.logger.Info("Create Helm umbrella charts")
			// Initalize the umbrella chart
			if err := umbrellaChartHandler.InitUmbrellaChart(event, stages); err != nil {
				return fmt.Errorf("Error when initializing the umbrella chart for project %s: %s", event.Project, err.Error())
			}
		}
	}

	kubeClient, err := keptnutils.GetKubeAPI(true)
	if err != nil {
		return err
	}
	for _, stage := range stages {
		if err := o.onboardService(stage.StageName, event, url.String()); err != nil {
			o.logger.Error(err.Error())
			return err
		}
		if event.DeploymentStrategies[stage.StageName] == keptnevents.Duplicate && event.HelmChart != "" {
			// inject Istio to the namespace for blue-green deployments
			namespace, err := kubeClient.Namespaces().Get(helm.GetUmbrellaNamespace(event.Project, stage.StageName), v1.GetOptions{})
			if err != nil {
				o.logger.Error(err.Error())
				return err
			}

			if namespace != nil {
				o.logger.Debug(fmt.Sprintf("Inject Istio to the %s namespace for blue-green deployments", helm.GetUmbrellaNamespace(event.Project, stage.StageName)))

				if namespace.ObjectMeta.Labels == nil {
					namespace.ObjectMeta.Labels = make(map[string]string)
				}

				namespace.ObjectMeta.Labels["istio-injection"] = "enabled"
				_, err = kubeClient.Namespaces().Update(namespace)
				if err != nil {
					o.logger.Error(err.Error())
					return err
				}
			}
		}
	}

	o.logger.Info(fmt.Sprintf("Finished creating service %s in project %s", event.Service, event.Project))
	return nil
}

func (o *Onboarder) checkAndSetServiceName(event *keptnevents.ServiceCreateEventData) error {

	errorMsg := "Service name contains upper case letter(s) or special character(s).\n " +
		"Keptn relies on the following conventions: " +
		"start with a lower case letter, then lower case letters, numbers, and hyphens are allowed."

	if event.HelmChart == "" {
		// Case when only a service is created but not onboarded (i.e. no Helm chart is available)
		if !keptnutils.ValidateKeptnEntityName(event.Service) {
			return errors.New(errorMsg)
		}
		return nil
	}

	helmChartData, err := base64.StdEncoding.DecodeString(event.HelmChart)
	if err != nil {
		return fmt.Errorf("Error when decoding the Helm chart: %v", err)
	}
	ch, err := keptnutils.LoadChart(helmChartData)
	if err != nil {
		return fmt.Errorf("Error when loading Helm chart: %v", err)
	}
	services, err := keptnutils.GetRenderedServices(ch)
	if err != nil {
		return fmt.Errorf("Error when rendering services: %v", err)
	}
	if len(services) != 1 {
		return fmt.Errorf("Helm chart has to contain exactly one Kubernetes service but has %d", len(services))
	}
	k8sServiceName := services[0].Name
	if !keptnutils.ValidateKeptnEntityName(k8sServiceName) {
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

func (o *Onboarder) onboardService(stageName string, event *keptnevents.ServiceCreateEventData,
	configServiceURL string) error {

	serviceHandler := configutils.NewServiceHandler(configServiceURL)

	o.logger.Debug("Creating new keptn service " + event.Service + " in stage " + stageName)
	respErr, err := serviceHandler.CreateService(event.Project, stageName, event.Service)
	if respErr != nil {
		return errors.New(*respErr.Message)
	}
	if err != nil {
		return err
	}

	if event.HelmChart != "" {
		helmChartData, err := base64.StdEncoding.DecodeString(event.HelmChart)
		if err != nil {
			o.logger.Error("Error when decoding the Helm chart")
			return err
		}

		o.logger.Debug("Storing the Helm chart provided by the user in stage " + stageName)
		if err := keptnutils.StoreChart(event.Project, event.Service, stageName, helm.GetChartName(event.Service, false),
			helmChartData, configServiceURL); err != nil {
			o.logger.Error("Error when storing the Helm chart: " + err.Error())
			return err
		}

		if err := o.updateUmbrellaChart(event.Project, stageName, helm.GetChartName(event.Service, false), configServiceURL); err != nil {
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
		return o.updateUmbrellaChart(event.Project, stageName, helmChartName, configServiceURL)
	}

	return nil
}

// IsGeneratedChartEmpty checks whether the generated chart is empty
func (c *Onboarder) IsGeneratedChartEmpty(chart *chart.Chart) bool {

	return len(chart.Templates) == 0
}

func (o *Onboarder) OnboardGeneratedService(helmManifest string, project string, stageName string,
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
		o.logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, "+
			"a chart for a duplicate deployment strategy is generated", service, stageName, strategy.String()))
		generatedChart, err = chartGenerator.GenerateDuplicateManagedChart(helmManifest, project, stageName, service)
		if err != nil {
			o.logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	} else {
		o.logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, a mesh chart is generated",
			service, stageName, strategy.String()))
		generatedChart, err = chartGenerator.GenerateMeshChart(helmManifest, project, stageName, service)
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

func (o *Onboarder) updateUmbrellaChart(project, stage, helmChartName, configServiceURL string) error {

	umbrellaChartHandler := helm.NewUmbrellaChartHandler(configServiceURL)
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
