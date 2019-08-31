package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	mesh             mesh.Mesh
	logger           *keptnutils.Logger
	configServiceURL string
	canaryLevelGen   helm.CanaryLevelGenerator
}

func NewConfigurationChanger(mesh mesh.Mesh, canaryLevelGen helm.CanaryLevelGenerator,
	logger *keptnutils.Logger, configServiceURL string) *ConfigurationChanger {
	return &ConfigurationChanger{mesh: mesh, canaryLevelGen: canaryLevelGen, logger: logger, configServiceURL: configServiceURL}
}

// ChangeAndApplyConfiguration changes the configuration and applies it in the cluster
func (c *ConfigurationChanger) ChangeAndApplyConfiguration(ce cloudevents.Event) error {

	e := &keptnevents.ConfigurationChangeEventData{}
	if err := ce.DataAs(e); err != nil {
		c.logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if len(e.ValuesPrimary) > 0 {
		c.changeValues(e, true, e.ValuesPrimary)
	}
	if len(e.ValuesCanary) > 0 {
		c.changeValues(e, false, e.ValuesCanary)
	}

	// Change canary
	if e.Canary != nil {
		c.changeCanary(e)
	}

	return nil
}

func (c *ConfigurationChanger) changeValues(e *keptnevents.ConfigurationChangeEventData, generated bool, newValues map[string]*chart.Value) error {

	// Read chart
	chart, err := helm.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), c.configServiceURL)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when reading chart %s from project %s: %s", helm.GetChartName(e.Service, generated), e.Project, err.Error()))
		return err
	}
	// Change values
	for k, v := range newValues {
		chart.Values.Values[k] = v
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when packaging modified chart %s from project %s: %s", helm.GetChartName(e.Service, generated), e.Project, err.Error()))
		return err
	}
	err = helm.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), chartData, c.configServiceURL)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when storing modified chart %s from project %s: %s", helm.GetChartName(e.Service, generated), e.Project, err.Error()))
		return err
	}
	return c.applyConfiguration(e, generated)
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
		if err := c.copyCanaryToPrimary(e); err != nil {
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

	case keptnevents.Set:
		err := c.setCanaryWeight(e, e.Canary.Value)
		if err != nil {
			return err
		}
		err = c.applyConfiguration(e, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ConfigurationChanger) copyCanaryToPrimary(e *keptnevents.ConfigurationChangeEventData) error {

	// TODO: Implement this function
	return nil
}

func (c *ConfigurationChanger) setCanaryWeight(e *keptnevents.ConfigurationChangeEventData, canaryWeight int32) error {

	// Read chart
	chart, err := helm.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), c.configServiceURL)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when reading chart %s from project %s: %s", helm.GetChartName(e.Service, true), e.Project, err.Error()))
		return err
	}

	// Set weights in all virtualservices
	for _, template := range chart.Templates {
		if strings.HasPrefix(template.Name, "templates/") &&
			strings.HasSuffix(template.Name, c.mesh.GetVirtualServiceSuffix()) {

			vs, err := c.mesh.UpdateWeights(template.Data, canaryWeight)
			if err != nil {
				c.logger.Error(fmt.Sprintf("Error when setting new weights in VirtualService %s from chart %s and project %s: %s",
					template.Name, helm.GetChartName(e.Service, true), e.Project, err.Error()))
				return err
			}
			template.Data = vs
		}
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when packaging modified chart %s from project %s: %s", helm.GetChartName(e.Service, true), e.Project, err.Error()))
		return err
	}
	err = helm.StoreChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, true), chartData, c.configServiceURL)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when storing modified chart %s from project %s: %s", helm.GetChartName(e.Service, true), e.Project, err.Error()))
		return err
	}
	return nil
}

// deleteRelease deletes a helm reease
func (c *ConfigurationChanger) deleteRelease(e *keptnevents.ConfigurationChangeEventData, generated bool) error {
	releaseName := helm.GetReleaseName(e.Project, e.Service, e.Stage, generated)
	// TODO: Make differentiation between canary level
	if _, err := keptnutils.ExecuteCommand("helm", []string{"delete", "--purge", releaseName}); err != nil {
		c.logger.Error(fmt.Sprintf("Error when deleting release %s: %s", releaseName, err.Error()))
		return err
	}
	return nil
}

// applyConfiguration applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func (c *ConfigurationChanger) applyConfiguration(e *keptnevents.ConfigurationChangeEventData, generated bool) error {

	ch, err := helm.GetChart(e.Project, e.Service, e.Stage, helm.GetChartName(e.Service, generated), c.configServiceURL)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when reading chart %s: %s", helm.GetChartName(e.Service, generated), err.Error()))
		return err
	}

	helmChartDir, err := ioutil.TempDir("", "")
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when creating temporary directory: %s", err.Error()))
		return err
	}
	defer os.RemoveAll(helmChartDir)

	chartPath, err := chartutil.Save(ch, helmChartDir)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error()))
		return err
	}

	releaseName := helm.GetReleaseName(e.Project, e.Service, e.Stage, generated)
	namespace := c.canaryLevelGen.GetNamespace(e.Project, e.Stage, generated)
	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		c.logger.Error(fmt.Sprintf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error()))
		return err
	}

	useInClusterConfig := false
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	}

	if err := keptnutils.WaitForDeploymentsInNamespace(useInClusterConfig, namespace); err != nil {
		c.logger.Error(fmt.Sprintf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error()))
		return err
	}
	c.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	return nil
}
