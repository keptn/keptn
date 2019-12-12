package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/ghodss/yaml"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

// ConfigurationChanger is a container of variables required for changing the configuration of a service
type ConfigurationChanger struct {
	mesh                  mesh.Mesh
	generatedChartHandler *helm.GeneratedChartHandler
	canaryLevelGen        helm.CanaryLevelGenerator
	logger                keptnutils.LoggerInterface
	keptnDomain           string
}

// NewConfigurationChanger creates a new ConfigurationChanger
func NewConfigurationChanger(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger keptnutils.LoggerInterface, keptnDomain string) *ConfigurationChanger {
	generatedChartHandler := helm.NewGeneratedChartHandler(mesh, canaryLevelGen, keptnDomain)
	return &ConfigurationChanger{mesh: mesh, generatedChartHandler: generatedChartHandler, canaryLevelGen: canaryLevelGen,
		logger: logger, keptnDomain: keptnDomain}
}

// ChangeAndApplyConfiguration changes the configuration and applies it in the cluster
func (c *ConfigurationChanger) ChangeAndApplyConfiguration(ce cloudevents.Event, loggingDone chan bool) error {

	defer func() { loggingDone <- true }()

	e := &keptnevents.ConfigurationChangeEventData{}
	if err := ce.DataAs(e); err != nil {
		c.logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" && e.Stage == "" {
		stage, err := getFirstStage(e.Project)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Error when reading shipyard: %s", err.Error()))
			return err
		}
		e.Stage = stage
	}

	genChart, err := getGeneratedChart(e)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	deploymentStrategy, err := getDeploymentStrategyOfService(genChart)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	if len(e.ValuesCanary) > 0 {
		err := c.applyValuesCanary(e, genChart, deploymentStrategy)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
	}

	if len(e.FileChangesUserChart) > 0 {
		ch, err := c.updateChart(e, false, changeUserChart)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, deploymentStrategy, false); err != nil {
			c.logger.Error(err.Error())
			return err
		}
	}

	if len(e.FileChangesGeneratedChart) > 0 {
		ch, err := c.updateChart(e, true, changeGeneratedChart)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, deploymentStrategy, true); err != nil {
			c.logger.Error(err.Error())
			return err
		}
	}

	if len(e.FileChangesUmbrellaChart) > 0 {
		if err := c.updateUmbrellaChart(e); err != nil {
			c.logger.Error(err.Error())
		}
	}

	// Change canary
	if e.Canary != nil {
		c.logger.Debug(fmt.Sprintf("Canary action %s for service %s in stage %s of project %s was received",
			e.Canary.Action, e.Service, e.Stage, e.Project))

		if deploymentStrategy == keptnevents.Duplicate {
			c.logger.Debug(fmt.Sprintf("Apply canary action %s for service %s in stage %s of project %s", e.Canary.Action, e.Service, e.Stage, e.Project))

			if err := c.changeCanary(e, deploymentStrategy); err != nil {
				c.logger.Error(err.Error())
				return err
			}
		} else {
			c.logger.Debug(fmt.Sprintf("Discard canary action %s for service %s in stage %s of project %s because service is not duplicated",
				e.Canary.Action, e.Service, e.Stage, e.Project))
			if os.Getenv("PRE_WORKFLOW_ENGINE") != "true" {
				c.logger.Error(
					fmt.Sprintf("Cannot process received canary instructions as no duplicate deployment for service %s in stage %s of project %s is available",
						e.Service, e.Stage, e.Project))
			}
		}
	}

	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		strings.HasSuffix(ce.Source(), "remediation-service") {
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, "real-user", deploymentStrategy, "", ""); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
		return nil
	}

	// Send deployment finished event
	// Note that this condition also stops the keptn-flow if an artifact is discarded
	if os.Getenv("PRE_WORKFLOW_ENGINE") == "true" &&
		!(e.Canary != nil && (e.Canary.Action == keptnevents.Discard || e.Canary.Action == keptnevents.Promote)) {

		testStrategy, err := getTestStrategy(e.Project, e.Stage)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}

		image := ""
		tag := ""

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
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, testStrategy, deploymentStrategy, image, tag); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
	}

	return nil
}

func (c *ConfigurationChanger) updateUmbrellaChart(e *keptnevents.ConfigurationChangeEventData) error {

	umbrellaChartHandler := helm.NewUmbrellaChartHandler(c.mesh)

	if err := umbrellaChartHandler.UpdateUmbrellaChart(e); err != nil {
		return err
	}

	umbrellaChart, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("error when creating a temporary directory: %s", err.Error())
	}
	if err := umbrellaChartHandler.GetUmbrellaChart(umbrellaChart, e.Project, e.Stage); err != nil {
		return fmt.Errorf("error when getting umbrella chart: %s", err)
	}
	if err := c.ApplyDirectory(umbrellaChart, helm.GetUmbrellaReleaseName(e.Project, e.Stage),
		helm.GetUmbrellaNamespace(e.Project, e.Stage)); err != nil {
		return fmt.Errorf("error when applying umbrella chart in stage %s: %s", e.Stage, err.Error())
	}
	return os.RemoveAll(umbrellaChart)
}

func (c *ConfigurationChanger) applyValuesCanary(e *keptnevents.ConfigurationChangeEventData,
	genChart *chart.Chart, deploymentStrategy keptnevents.DeploymentStrategy) error {
	ch, err := c.updateChart(e, false, changeValue)
	if err != nil {
		return err
	}
	err = c.ApplyChart(ch, e.Project, e.Stage, e.Service, deploymentStrategy, false)
	if err != nil {
		return err
	}
	onboarder := NewOnboarder(c.mesh, c.canaryLevelGen, c.logger, c.keptnDomain)
	if onboarder.IsGeneratedChartEmpty(genChart) {
		manifest, err := c.getManifest(e.Project, e.Stage, e.Service, false)
		if err != nil {
			return err
		}

		genChart, err = onboarder.OnboardGeneratedService(manifest, e.Project, e.Stage, e.Service, deploymentStrategy)
		if err != nil {
			return err
		}
		if deploymentStrategy == keptnevents.Direct {
			err := c.ApplyChart(genChart, e.Project, e.Stage, e.Service, deploymentStrategy, true)
			if err != nil {
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

func getGeneratedChart(e *keptnevents.ConfigurationChangeEventData) (*chart.Chart, error) {
	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	helmChartName := helm.GetChartName(e.Service, true)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, url.String())
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

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	helmChartName := helm.GetChartName(e.Service, generated)
	c.logger.Info(fmt.Sprintf("Start updating chart %s of stage %s", helmChartName, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, url.String())
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
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helmChartName, chartData, url.String()); err != nil {
		return nil, err
	}
	c.logger.Info(fmt.Sprintf("Finished updating chart %s of stage %s", helmChartName, e.Stage))
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

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	c.logger.Info(fmt.Sprintf("Start updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), url.String())
	if err != nil {
		return nil, err
	}

	c.generatedChartHandler.UpdateCanaryWeight(chart, canaryWeight)

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return nil, err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), chartData, url.String()); err != nil {
		return nil, err
	}
	c.logger.Info(fmt.Sprintf("Finished updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	return chart, nil
}

func (c *ConfigurationChanger) changeCanary(e *keptnevents.ConfigurationChangeEventData,
	deploymentStrategy keptnevents.DeploymentStrategy) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	switch e.Canary.Action {
	case keptnevents.Discard:
		ch, err := c.setCanaryWeight(e, 0)
		if err != nil {
			return err
		}
		if err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, deploymentStrategy, true); err != nil {
			return err
		}
		userChart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, false), url.String())
		if err != nil {
			return err
		}
		if err := c.ApplyChartWithReplicas(userChart, e.Project, e.Stage, e.Service,
			deploymentStrategy, false, 0); err != nil {
			return err
		}

	case keptnevents.Promote:
		ch, err := c.setCanaryWeight(e, 100)
		if err != nil {
			return err
		}
		if err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, deploymentStrategy, true); err != nil {
			return err
		}

		chartGenerator := helm.NewGeneratedChartHandler(c.mesh, c.canaryLevelGen, c.keptnDomain)
		manifest, err := c.getManifest(e.Project, e.Stage, e.Service, false)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		genChart, err := chartGenerator.GenerateDuplicateManagedChart(manifest, e.Project, e.Stage, e.Service)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.generatedChartHandler.UpdateCanaryWeight(genChart, int32(100)); err != nil {
			return err
		}
		genChartData, err := keptnutils.PackageChart(genChart)
		if err != nil {
			return err
		}
		if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), genChartData, url.String()); err != nil {
			return err
		}
		if err := c.ApplyChart(genChart, e.Project, e.Stage, e.Service, deploymentStrategy, true); err != nil {
			return err
		}

		genChart, err = c.setCanaryWeight(e, 0)
		if err != nil {
			return err
		}
		if err := c.ApplyChart(genChart, e.Project, e.Stage, e.Service, deploymentStrategy, true); err != nil {
			return err
		}
		userChart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, false), url.String())
		if err != nil {
			return err
		}
		if err := c.ApplyChartWithReplicas(userChart, e.Project, e.Stage, e.Service,
			deploymentStrategy, false, 0); err != nil {
			return err
		}

	case keptnevents.Set:
		ch, err := c.setCanaryWeight(e, e.Canary.Value)
		if err != nil {
			return err
		}
		if err := c.ApplyChart(ch, e.Project, e.Stage, e.Service, deploymentStrategy, true); err != nil {
			return err
		}
	}

	return nil
}

// getManifest
func (c *ConfigurationChanger) getManifest(project, stage, service string, generated bool) (string, error) {

	releaseName := helm.GetReleaseName(project, stage, service, generated)
	namespace := c.canaryLevelGen.GetNamespace(project, stage, generated)
	msg, err := keptnutils.ExecuteCommand("helm", []string{"get", "manifest", releaseName,
		"--namespace", namespace})
	if err != nil {
		return "", fmt.Errorf("Error when quering the manifest of chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	return msg, nil
}

// ApplyChart applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) ApplyChart(ch *chart.Chart, project, stage, service string,
	deploymentStrategy keptnevents.DeploymentStrategy, generated bool) error {

	return c.ApplyChartWithReplicas(ch, project, stage, service, deploymentStrategy, generated, -1)
}

// ApplyChartWithReplicas applies the chart of the provided service and additionally sets the replicas
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) ApplyChartWithReplicas(ch *chart.Chart, project, stage, service string,
	deploymentStrategy keptnevents.DeploymentStrategy, generated bool, replicaCount int) error {

	releaseName := helm.GetReleaseName(project, stage, service, generated)
	namespace := c.canaryLevelGen.GetNamespace(project, stage, generated)
	c.logger.Info(fmt.Sprintf("Start upgrading chart %s in namespace %s", releaseName, namespace))

	if len(ch.Templates) > 0 {
		helmChartDir, err := ioutil.TempDir("", "")
		if err != nil {
			return fmt.Errorf("Error when creating temporary directory: %s", err.Error())
		}
		defer os.RemoveAll(helmChartDir)

		chartPath, err := chartutil.Save(ch, helmChartDir)
		if err != nil {
			return fmt.Errorf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error())
		}

		deploymentName := getDeploymentName(deploymentStrategy, generated)
		var msg string
		if replicaCount >= 0 {
			msg, err = keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
				chartPath, "--namespace", namespace, "--wait", "--force",
				"--set", "keptn.project=" + project, "--set", "keptn.stage=" + stage,
				"--set", "keptn.service=" + service, "--set", "keptn.deployment=" + deploymentName,
				"--set", "replicaCount=" + strconv.Itoa(replicaCount)})
		} else {
			msg, err = keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
				chartPath, "--namespace", namespace, "--wait", "--force",
				"--set", "keptn.project=" + project, "--set", "keptn.stage=" + stage,
				"--set", "keptn.service=" + service, "--set", "keptn.deployment=" + deploymentName})

		}
		c.logger.Debug(msg)
		if err != nil {
			return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
				releaseName, namespace, err.Error())
		}

		if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
			return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
		}
		c.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
		return nil
	}
	c.logger.Debug("Upgrade not done as this is an empty chart")
	return nil
}

// ApplyDirectory applies the provided directory
func (c *ConfigurationChanger) ApplyDirectory(chartPath, releaseName, namespace string) error {

	msg, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--reset-values", "--wait", "--force"})
	if err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	c.logger.Debug(msg)

	if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
		return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	return nil
}

func getInClusterConfig() bool {
	if os.Getenv("ENVIRONMENT") == "production" {
		return true
	}
	return false
}

func int32Ptr(i int32) *int32 { return &i }
