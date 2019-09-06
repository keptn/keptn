package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/ghodss/yaml"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"k8s.io/helm/pkg/chartutil"
)

func init() {
	_, err := keptnutils.ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		log.Fatal(err)
	}
}

// ConfigurationChanger is a container of variables required for changing the configuration of a service
type ConfigurationChanger struct {
	generatedChartHandler *helm.GeneratedChartHandler
	canaryLevelGen        helm.CanaryLevelGenerator
	configServiceURL      string
	logger                keptnutils.LoggerInterface
	keptnDomain           string
}

// NewConfigurationChanger creates a new ConfigurationChanger
func NewConfigurationChanger(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger keptnutils.LoggerInterface, keptnDomain string, configServiceURL string) *ConfigurationChanger {
	generatedChartHandler := helm.NewGeneratedChartHandler(mesh, canaryLevelGen, keptnDomain)
	return &ConfigurationChanger{generatedChartHandler: generatedChartHandler, canaryLevelGen: canaryLevelGen,
		configServiceURL: configServiceURL, logger: logger, keptnDomain: keptnDomain}
}

// ChangeAndApplyConfiguration changes the configuration and applies it in the cluster
func (c *ConfigurationChanger) ChangeAndApplyConfiguration(ce cloudevents.Event) error {

	e := &keptnevents.ConfigurationChangeEventData{}
	if err := ce.DataAs(e); err != nil {
		c.logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if os.Getenv("pre_workflow_engine") == "true" && e.Stage == "" {
		stage, err := getFirstStage(e.Project, c.configServiceURL)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Error when reading shipyard: %s" + err.Error()))
			return err
		}
		e.Stage = stage
	}

	if len(e.ValuesPrimary) > 0 {
		if err := c.changeValues(e, true, e.ValuesPrimary); err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}
	}
	if len(e.ValuesCanary) > 0 {
		if err := c.changeValues(e, false, e.ValuesCanary); err != nil {
			c.logger.Error(err.Error())
			return err
		}
		if err := c.applyConfiguration(e, false); err != nil {
			return err
		}
	}

	// Change canary
	if e.Canary != nil {
		if err := c.changeCanary(e); err != nil {
			c.logger.Error(err.Error())
			return err
		}
	}

	// Send deployment finished event
	// Note that this condition also stops the keptn-flow if an artifact is discarded
	if os.Getenv("pre_workflow_engine") == "true" &&
		!(e.Canary != nil && e.Canary.Action == keptnevents.Discard) {
		var shkeptncontext string
		ce.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
		if err := sendDeploymentFinishedEvent(shkeptncontext, e.Project, e.Stage, e.Service, c.configServiceURL); err != nil {
			c.logger.Error(fmt.Sprintf("Cannot send deployment finished event: %s", err.Error()))
			return err
		}
	}

	return nil
}

func (c *ConfigurationChanger) changeValues(e *keptnevents.ConfigurationChangeEventData, generated bool, newValues map[string]interface{}) error {

	helmChartName := helm.GetChartName(e.Service, generated)
	c.logger.Info(fmt.Sprintf("Start updating values of chart %s", helmChartName))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, c.configServiceURL)
	if err != nil {
		return err
	}

	values := make(map[string]interface{})
	yaml.Unmarshal([]byte(chart.Values.Raw), &values)

	// Change values
	for k, v := range newValues {
		values[k] = v
	}

	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return err
	}
	chart.Values.Raw = string(valuesData)

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helmChartName, chartData, c.configServiceURL); err != nil {
		return err
	}
	c.logger.Info(fmt.Sprintf("Finished updating values of chart %s", helmChartName))
	return nil
}

func (c *ConfigurationChanger) setCanaryWeight(e *keptnevents.ConfigurationChangeEventData, canaryWeight int32) error {

	c.logger.Info(fmt.Sprintf("Start updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	// Read chart
	chart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), c.configServiceURL)
	if err != nil {
		return err
	}

	c.generatedChartHandler.UpdateCanaryWeight(chart, canaryWeight)

	// Store chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return err
	}
	if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), chartData, c.configServiceURL); err != nil {
		return err
	}
	c.logger.Info(fmt.Sprintf("Finished updating canary weight to %d for service %s of project %s in stage %s",
		canaryWeight, e.Service, e.Project, e.Stage))
	return nil
}

func (c *ConfigurationChanger) changeCanary(e *keptnevents.ConfigurationChangeEventData) error {

	switch e.Canary.Action {
	case keptnevents.Discard:
		if err := c.setCanaryWeight(e, 0); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}
		if err := c.deleteRelease(e, false); err != nil {
			return err
		}

	case keptnevents.Promote:
		if err := c.setCanaryWeight(e, 100); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}

		userChart, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, false), c.configServiceURL)
		if err != nil {
			return err
		}
		genChartData, err := c.generatedChartHandler.GenerateManagedChart(userChart, e.Project, e.Stage)
		if err != nil {
			return err
		}
		genChart, err := keptnutils.LoadChart(genChartData)
		if err != nil {
			return err
		}
		if err := c.generatedChartHandler.UpdateCanaryWeight(genChart, int32(100)); err != nil {
			return err
		}
		genChartData, err = keptnutils.PackageChart(genChart)
		if err != nil {
			return err
		}
		if err := keptnutils.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), genChartData, c.configServiceURL); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}

		if err := c.setCanaryWeight(e, 0); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}
		if err := c.deleteRelease(e, false); err != nil {
			return err
		}

	case keptnevents.Set:
		if err := c.setCanaryWeight(e, e.Canary.Value); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}
	}

	return nil
}

// deleteRelease deletes a helm release
func (c *ConfigurationChanger) deleteRelease(e *keptnevents.ConfigurationChangeEventData, generated bool) error {
	c.logger.Info(fmt.Sprintf("Start deleting deployment/release of service %s of project %s in stage %s", e.Service, e.Project, e.Stage))
	if err := c.canaryLevelGen.DeleteRelease(e.Project, e.Stage, e.Service, generated, c.configServiceURL); err != nil {
		return fmt.Errorf("Error when deleting release %s: %s",
			helm.GetReleaseName(e.Project, e.Stage, e.Service, generated), err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished deleting deployment/release of service %s of project %s in stage %s", e.Service, e.Project, e.Stage))
	return nil
}

// applyConfiguration applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) applyConfiguration(e *keptnevents.ConfigurationChangeEventData, generated bool) error {

	releaseName := helm.GetReleaseName(e.Project, e.Stage, e.Service, generated)
	namespace := c.canaryLevelGen.GetNamespace(e.Project, e.Stage, generated)
	c.logger.Info(fmt.Sprintf("Start upgrading chart %s in namespace %s", releaseName, namespace))

	ch, err := keptnutils.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), c.configServiceURL)
	if err != nil {
		return fmt.Errorf("Error when reading chart %s: %s", helm.GetChartName(e.Service, generated), err.Error())
	}

	helmChartDir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("Error when creating temporary directory: %s", err.Error())
	}
	defer os.RemoveAll(helmChartDir)

	chartPath, err := chartutil.Save(ch, helmChartDir)
	if err != nil {
		return fmt.Errorf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error())
	}

	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}

	if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
		return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	return nil
}

// ApplyDirectory applies the provided directory
func ApplyDirectory(chartPath, releaseName, namespace string) error {

	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}

	if err := keptnutils.WaitForDeploymentsInNamespace(getInClusterConfig(), namespace); err != nil {
		return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	return nil
}

func getInClusterConfig() bool {
	if os.Getenv("env") == "production" {
		return true
	}
	return false
}
