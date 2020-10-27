package controller

import (
	"errors"
	"fmt"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/ghodss/yaml"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"
)

// ConfigurationChanger is a container of variables required for changing the configuration of a service
type ConfigurationChanger struct {
	mesh                  mesh.Mesh
	generatedChartHandler *helm.GeneratedChartHandler
	keptnHandler          *keptnevents.Keptn
	helmExecutor          helm.HelmExecutor
	configServiceURL      string
}

// NewConfigurationChanger creates a new ConfigurationChanger
func NewConfigurationChanger(mesh mesh.Mesh, keptnHandler *keptnevents.Keptn, configServiceURL string) *ConfigurationChanger {
	generatedChartHandler := helm.NewGeneratedChartHandler(mesh, keptnHandler.Logger)
	helmExecutor := helm.NewHelmV3Executor(keptnHandler.Logger)
	return &ConfigurationChanger{
		mesh:                  mesh,
		generatedChartHandler: generatedChartHandler,
		keptnHandler:          keptnHandler,
		helmExecutor:          helmExecutor,
		configServiceURL:      configServiceURL,
	}
}

// ChangeAndApplyConfiguration changes the configuration and applies it in the cluster
func (c *ConfigurationChanger) ChangeAndApplyConfiguration(ce cloudevents.Event, loggingDone chan bool) error {

	defer func() { loggingDone <- true }()

	e := &keptnevents.ConfigurationChangeEventData{}
	if err := ce.DataAs(e); err != nil {
		c.keptnHandler.Logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	keptnHandler, err := keptnevents.NewKeptn(&ce, keptnevents.KeptnOpts{})
	if err != nil {
		c.keptnHandler.Logger.Error("Could not initialize keptn handler: " + err.Error())
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" && e.Stage == "" {
		stage, err := getFirstStage(keptnHandler)
		keptnHandler.KeptnBase.Stage = stage
		if err != nil {
			c.keptnHandler.Logger.Error(fmt.Sprintf("Error when reading shipyard: %s", err.Error()))
			return err
		}
		e.Stage = stage
	}

	genChart, err := c.getGeneratedChart(e)
	if err != nil {
		c.keptnHandler.Logger.Error(err.Error())
		return err
	}
	deploymentStrategy, err := getDeploymentStrategyOfService(genChart)
	if err != nil {
		c.keptnHandler.Logger.Error(err.Error())
		return err
	}

	if len(e.ValuesCanary) > 0 {
		err := c.applyValuesCanary(e, genChart, deploymentStrategy)
		if err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
	}

	if len(e.FileChangesUserChart) > 0 {
		ch, err := c.updateChart(e, false, changeUserChart)
		if err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
		if err := c.upgradeChart(ch, *e, deploymentStrategy); err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
	}

	if len(e.FileChangesGeneratedChart) > 0 {
		ch, err := c.updateChart(e, true, changeGeneratedChart)
		if err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
		if err := c.upgradeChart(ch, *e, deploymentStrategy); err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
	}

	if len(e.FileChangesUmbrellaChart) > 0 {
		if err := c.updateUmbrellaChart(e); err != nil {
			c.keptnHandler.Logger.Error(err.Error())
		}
	}

	// Change canary
	if e.Canary != nil {
		c.keptnHandler.Logger.Debug(fmt.Sprintf("Canary action %s for service %s in stage %s of project %s was received",
			e.Canary.Action, e.Service, e.Stage, e.Project))

		if deploymentStrategy == keptnevents.Duplicate {
			c.keptnHandler.Logger.Debug(fmt.Sprintf("Apply canary action %s for service %s in stage %s of project %s", e.Canary.Action, e.Service, e.Stage, e.Project))

			if err := c.changeCanary(e, deploymentStrategy); err != nil {
				c.keptnHandler.Logger.Error(err.Error())
				return err
			}
		} else {
			c.keptnHandler.Logger.Debug(fmt.Sprintf("Discard canary action %s for service %s in stage %s of project %s because service is not duplicated",
				e.Canary.Action, e.Service, e.Stage, e.Project))
			if os.Getenv("PRE_WORKFLOW_ENGINE") != "true" {
				c.keptnHandler.Logger.Error(
					fmt.Sprintf("Cannot process received canary instructions as no duplicate deployment for service %s in stage %s of project %s is available",
						e.Service, e.Stage, e.Project))
			}
		}
	}

	// Send deployment finished event
	// Note that this condition also stops the keptn-flow if an artifact is discarded
	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		!(e.Canary != nil && (e.Canary.Action == keptnevents.Discard || e.Canary.Action == keptnevents.Promote)) {

		testStrategy, err := getTestStrategy(keptnHandler, e.Stage)
		if err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}

		image := ""
		tag := ""
		labels := e.Labels

		for k, v := range e.ValuesCanary {
			if k == "image" {
				splittedImage := strings.Split(v.(string), ":")
				if len(splittedImage) > 0 {
					image = splittedImage[0]
					tag = splittedImage[1]
				}
			}
		}
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(keptnHandler, testStrategy, deploymentStrategy, image, tag, labels, mesh.GetIngressHostnameSuffix(), mesh.GetIngressProtocol(), mesh.GetIngressPort()); err != nil {
			c.keptnHandler.Logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
	}

	return nil
}

func (c *ConfigurationChanger) updateUmbrellaChart(e *keptnevents.ConfigurationChangeEventData) error {

	umbrellaChartHandler := helm.NewUmbrellaChartHandler(c.configServiceURL)

	if err := umbrellaChartHandler.UpdateUmbrellaChart(e); err != nil {
		return err
	}

	ch, err := umbrellaChartHandler.GetUmbrellaChart(e.Project, e.Stage)
	if err != nil {
		return fmt.Errorf("error when getting umbrella chart: %s", err)
	}
	if err := c.helmExecutor.UpgradeChart(ch, helm.GetUmbrellaReleaseName(e.Project, e.Stage),
		helm.GetUmbrellaNamespace(e.Project, e.Stage), nil); err != nil {
		return fmt.Errorf("error when applying umbrella chart in stage %s: %s", e.Stage, err.Error())
	}
	return nil
}

func (c *ConfigurationChanger) applyValuesCanary(e *keptnevents.ConfigurationChangeEventData,
	genChart *chart.Chart, deploymentStrategy keptnevents.DeploymentStrategy) error {
	ch, err := c.updateChart(e, false, changeValue)
	if err != nil {
		return err
	}
	if err := c.upgradeChart(ch, *e, deploymentStrategy); err != nil {
		return err
	}
	onboarder := NewOnboarder(c.mesh, c.keptnHandler, c.configServiceURL)
	if onboarder.IsGeneratedChartEmpty(genChart) {
		userChartManifest, err := c.helmExecutor.GetManifest(helm.GetReleaseName(e.Project, e.Stage, e.Service, false),
			e.Project+"-"+e.Stage)
		if err != nil {
			return err
		}
		genChart, err = onboarder.OnboardGeneratedService(userChartManifest, e.Project, e.Stage, e.Service, deploymentStrategy)
		if err != nil {
			return err
		}
		if deploymentStrategy == keptnevents.Direct {
			if err := c.upgradeChart(genChart, *e, deploymentStrategy); err != nil {
				return err
			}
		}
	}
	return nil
}

func changeUserChart(e *keptnevents.ConfigurationChangeEventData, chart *chart.Chart) error {
	return applyFileChanges(e.FileChangesUserChart, chart)
}

func changeGeneratedChart(e *keptnevents.ConfigurationChangeEventData, chart *chart.Chart) error {
	return applyFileChanges(e.FileChangesGeneratedChart, chart)
}

func applyFileChanges(newFileContent map[string]string, ch *chart.Chart) error {

	for _, template := range ch.Templates {
		if val, ok := newFileContent[template.Name]; ok {
			template.Data = []byte(val)
			delete(newFileContent, template.Name)
		}
	}

	for uri, content := range newFileContent {
		if strings.HasPrefix(uri, "templates/") {
			// Add a new file to templates/
			template := &chart.File{Name: uri, Data: []byte(content)}
			ch.Templates = append(ch.Templates, template)
		} else if uri == "values.yaml" {
			values, err := loadValues(content)
			if err != nil {
				return err
			}
			ch.Values = values
		} else {
			return errors.New(fmt.Sprintf("Unsupported update of file %s", uri))
		}
	}
	return nil
}

func loadValues(valuesString string) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(valuesString), &values); err != nil {
		return nil, fmt.Errorf("Cannot load values: %v", err)
	}
	return values, nil
}

func (c *ConfigurationChanger) getGeneratedChart(e *keptnevents.ConfigurationChangeEventData) (*chart.Chart, error) {
	helmChartName := helm.GetChartName(e.Service, true)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, c.configServiceURL)
}

func getDeploymentName(strategy keptnevents.DeploymentStrategy, generatedChart bool) string {

	if strategy == keptnevents.Duplicate && generatedChart {
		return "primary"
	} else if strategy == keptnevents.Duplicate && !generatedChart {
		return "canary"
	} else if strategy == keptnevents.Direct {
		return "direct"
	}
	return ""
}

func getDeploymentStrategyOfService(ch *chart.Chart) (keptnevents.DeploymentStrategy, error) {
	for _, keyword := range ch.Metadata.Keywords {
		if keyword == "deployment_strategy=duplicate" {
			return keptnevents.Duplicate, nil
		}
		if keyword == "deployment_strategy=direct" {
			return keptnevents.Direct, nil
		}
	}
	return keptnevents.Duplicate, errors.New("Cannot find deployment_strategy in keywords")
}

func (c *ConfigurationChanger) updateChart(e *keptnevents.ConfigurationChangeEventData, generated bool,
	editChart func(*keptnevents.ConfigurationChangeEventData, *chart.Chart) error) (*chart.Chart, error) {

	helmChartName := helm.GetChartName(e.Service, generated)
	c.keptnHandler.Logger.Info(fmt.Sprintf("Start updating chart %s of stage %s", helmChartName, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, c.configServiceURL)
	if err != nil {
		return nil, err
	}

	// Edit chart
	err = editChart(e, chart)
	if err != nil {
		return nil, err
	}

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return nil, err
	}
	if _, err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helmChartName, chartData, c.configServiceURL); err != nil {
		return nil, err
	}
	c.keptnHandler.Logger.Info(fmt.Sprintf("Finished updating chart %s of stage %s", helmChartName, e.Stage))
	return chart, nil
}

func changeValue(e *keptnevents.ConfigurationChangeEventData, chart *chart.Chart) error {

	// Change values
	for k, v := range e.ValuesCanary {
		chart.Values[k] = v
	}
	return nil
}

func (c *ConfigurationChanger) setCanaryWeight(e *keptnevents.ConfigurationChangeEventData, canaryWeight int32) (*chart.Chart, error) {

	c.keptnHandler.Logger.Info(fmt.Sprintf("Start updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), c.configServiceURL)
	if err != nil {
		return nil, err
	}

	c.generatedChartHandler.UpdateCanaryWeight(chart, canaryWeight)

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return nil, err
	}
	if _, err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), chartData, c.configServiceURL); err != nil {
		return nil, err
	}
	c.keptnHandler.Logger.Info(fmt.Sprintf("Finished updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	return chart, nil
}

func (c *ConfigurationChanger) upgradeChart(ch *chart.Chart, configChange keptnevents.ConfigurationChangeEventData,
	strategy keptnevents.DeploymentStrategy) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	return c.helmExecutor.UpgradeChart(ch,
		helm.GetReleaseName(configChange.Project, configChange.Stage, configChange.Service, generated),
		configChange.Project+"-"+configChange.Stage,
		getKeptnValues(configChange.Project, configChange.Stage, configChange.Service,
			getDeploymentName(strategy, generated)))
}

func (c *ConfigurationChanger) upgradeChartWithReplicas(ch *chart.Chart, configChange keptnevents.ConfigurationChangeEventData,
	strategy keptnevents.DeploymentStrategy, replicas int) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	return c.helmExecutor.UpgradeChart(ch,
		helm.GetReleaseName(configChange.Project, configChange.Stage, configChange.Service, generated),
		configChange.Project+"-"+configChange.Stage,
		addReplicas(getKeptnValues(configChange.Project, configChange.Stage, configChange.Service,
			getDeploymentName(strategy, generated)), replicas))
}

func (c *ConfigurationChanger) changeCanary(e *keptnevents.ConfigurationChangeEventData,
	deploymentStrategy keptnevents.DeploymentStrategy) error {

	switch e.Canary.Action {
	case keptnevents.Discard:
		ch, err := c.setCanaryWeight(e, 0)
		if err != nil {
			return err
		}
		if err := c.upgradeChart(ch, *e, deploymentStrategy); err != nil {
			return err
		}
		userChart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, false), c.configServiceURL)
		if err != nil {
			return err
		}
		if err := c.upgradeChartWithReplicas(userChart, *e, deploymentStrategy, 0); err != nil {
			return err
		}

	case keptnevents.Promote:
		ch, err := c.setCanaryWeight(e, 100)
		if err != nil {
			return err
		}
		if err := c.upgradeChart(ch, *e, deploymentStrategy); err != nil {
			return err
		}

		chartGenerator := helm.NewGeneratedChartHandler(c.mesh, c.keptnHandler.Logger)
		userChartManifest, err := c.helmExecutor.GetManifest(helm.GetReleaseName(e.Project, e.Stage, e.Service, false),
			e.Project+"-"+e.Stage)
		if err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
		genChart, err := chartGenerator.GenerateDuplicateManagedChart(userChartManifest, e.Project, e.Stage, e.Service)
		if err != nil {
			c.keptnHandler.Logger.Error(err.Error())
			return err
		}
		if err := c.generatedChartHandler.UpdateCanaryWeight(genChart, int32(100)); err != nil {
			return err
		}
		genChartData, err := keptnutils.PackageChart(genChart)
		if err != nil {
			return err
		}
		if _, err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), genChartData, c.configServiceURL); err != nil {
			return err
		}
		if err := c.upgradeChart(genChart, *e, deploymentStrategy); err != nil {
			return err
		}

		genChart, err = c.setCanaryWeight(e, 0)
		if err != nil {
			return err
		}
		if err := c.upgradeChart(genChart, *e, deploymentStrategy); err != nil {
			return err
		}
		userChart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, false), c.configServiceURL)
		if err != nil {
			return err
		}
		if err := c.upgradeChartWithReplicas(userChart, *e, deploymentStrategy, 0); err != nil {
			return err
		}

	case keptnevents.Set:
		ch, err := c.setCanaryWeight(e, e.Canary.Value)
		if err != nil {
			return err
		}
		if err := c.upgradeChart(ch, *e, deploymentStrategy); err != nil {
			return err
		}
	}

	return nil
}

func getKeptnValues(project, stage, service, deploymentName string) map[string]interface{} {

	return map[string]interface{}{
		"keptn": map[string]interface{}{
			"project":    project,
			"stage":      stage,
			"service":    service,
			"deployment": deploymentName,
		},
	}
}

func addReplicas(vals map[string]interface{}, replicas int) map[string]interface{} {
	vals["replicaCount"] = replicas
	return vals
}
