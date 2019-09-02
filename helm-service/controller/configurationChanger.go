package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func init() {
	_, err := keptnutils.ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		log.Fatal(err)
	}
}

type ConfigurationChanger struct {
	generatedChartHandler *helm.GeneratedChartHandler
	canaryLevelGen        helm.CanaryLevelGenerator
	configServiceURL      string
	logger                *keptnutils.Logger
	keptnDomain           string
}

func NewConfigurationChanger(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger *keptnutils.Logger, keptnDomain string, configServiceURL string) *ConfigurationChanger {
	generatedChartHandler := helm.NewGeneratedChartHandler(mesh, canaryLevelGen, keptnDomain, configServiceURL)
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

	if len(e.ValuesPrimary) > 0 {
		if err := c.changeValues(e, true, e.ValuesPrimary); err != nil {
			c.logger.Error(err.Error())
			return err
		}
	}
	if len(e.ValuesCanary) > 0 {
		if err := c.changeValues(e, false, e.ValuesCanary); err != nil {
			c.logger.Error(err.Error())
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

	return nil
}

func (c *ConfigurationChanger) changeValues(e *keptnevents.ConfigurationChangeEventData, generated bool, newValues map[string]*chart.Value) error {

	// Read chart
	chart, err := helm.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), c.configServiceURL)
	if err != nil {
		return fmt.Errorf("Error when reading chart %s from project %s: %s",
			helm.GetChartName(e.Service, generated), e.Project, err.Error())
	}
	// Change values
	for k, v := range newValues {
		chart.Values.Values[k] = v
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		return fmt.Errorf("Error when packaging modified chart %s from project %s: %s",
			helm.GetChartName(e.Service, generated), e.Project, err.Error())
	}
	err = helm.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), chartData, c.configServiceURL)
	if err != nil {
		return fmt.Errorf("Error when storing modified chart %s from project %s: %s",
			helm.GetChartName(e.Service, generated), e.Project, err.Error())
	}
	return c.applyConfiguration(e, generated)
}

func (c *ConfigurationChanger) changeCanary(e *keptnevents.ConfigurationChangeEventData) error {

	switch e.Canary.Action {
	case keptnevents.Discard:
		if err := c.generatedChartHandler.SetCanaryWeight(e, 0); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}
		if err := c.deleteRelease(e, false); err != nil {
			return err
		}

	case keptnevents.Promote:
		if err := c.generatedChartHandler.SetCanaryWeight(e, 100); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}

		if err := c.generatedChartHandler.CopyCanaryToPrimary(e); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}
		if err := c.generatedChartHandler.SetCanaryWeight(e, 0); err != nil {
			return err
		}
		if err := c.applyConfiguration(e, true); err != nil {
			return err
		}

	case keptnevents.Set:
		if err := c.generatedChartHandler.SetCanaryWeight(e, e.Canary.Value); err != nil {
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

	if err := c.canaryLevelGen.DeleteRelease(e.Project, e.Stage, e.Service, generated, c.configServiceURL); err != nil {
		return fmt.Errorf("Error when deleting release %s: %s",
			helm.GetReleaseName(e.Project, e.Service, e.Stage, generated), err.Error())
	}
	return nil
}

// applyConfiguration applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) applyConfiguration(e *keptnevents.ConfigurationChangeEventData, generated bool) error {

	ch, err := helm.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), c.configServiceURL)
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

	releaseName := helm.GetReleaseName(e.Project, e.Service, e.Stage, generated)
	namespace := c.canaryLevelGen.GetNamespace(e.Project, e.Stage, generated)
	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}

	useInClusterConfig := false
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	}

	if err := keptnutils.WaitForDeploymentsInNamespace(useInClusterConfig, namespace); err != nil {
		return fmt.Errorf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error())
	}
	c.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	return nil
}
