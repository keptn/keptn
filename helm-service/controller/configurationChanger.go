package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

// ChangeConfiguration
func ChangeConfiguration(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string) error {

	if len(ce.ValuesPrimary) > 0 {
		changeValues(ce, true, ce.ValuesPrimary, logger, configServiceURL)
	}
	if len(ce.ValuesCanary) > 0 {
		changeValues(ce, false, ce.ValuesCanary, logger, configServiceURL)
	}

	// Change canary
	if ce.Canary != nil {
		changeCanary(ce, mesh, logger, configServiceURL)
	}

	return nil
}

func changeValues(ce *keptnevents.ConfigurationChangeEventData, generated bool, newValues map[string]*chart.Value,
	logger *keptnutils.Logger, configServiceURL string) error {

	// Read chart
	chart, err := helm.GetChart(ce.Project, ce.Service, ce.Stage, helm.GetChartName(ce.Service, generated), configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when reading chart %s from project %s: %s", helm.GetChartName(ce.Service, generated), ce.Project, err.Error()))
		return err
	}
	// Change values
	for k, v := range newValues {
		chart.Values.Values[k] = v
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when packaging modified chart %s from project %s: %s", helm.GetChartName(ce.Service, generated), ce.Project, err.Error()))
		return err
	}
	err = helm.StoreChart(ce.Project, ce.Service, ce.Stage, helm.GetChartName(ce.Service, generated), chartData, configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when storing modified chart %s from project %s: %s", helm.GetChartName(ce.Service, generated), ce.Project, err.Error()))
		return err
	}
	return ApplyConfiguration(ce, generated, logger, configServiceURL)
}

func changeCanary(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string) error {

	switch ce.Canary.Action {
	case keptnevents.Discard:
		if err := setCanaryWeight(ce, mesh, logger, configServiceURL, 0); err != nil {
			return err
		}
		if err := ApplyConfiguration(ce, true, logger, configServiceURL); err != nil {
			return err
		}
		if err := DeleteRelease(ce, false, logger); err != nil {
			return err
		}

	case keptnevents.Promote:
		if err := setCanaryWeight(ce, mesh, logger, configServiceURL, 100); err != nil {
			return err
		}
		if err := ApplyConfiguration(ce, true, logger, configServiceURL); err != nil {
			return err
		}
		if err := copyCanaryToPrimary(ce, mesh, logger, configServiceURL); err != nil {
			return err
		}
		if err := ApplyConfiguration(ce, true, logger, configServiceURL); err != nil {
			return err
		}
		if err := setCanaryWeight(ce, mesh, logger, configServiceURL, 0); err != nil {
			return err
		}
		if err := ApplyConfiguration(ce, true, logger, configServiceURL); err != nil {
			return err
		}

	case keptnevents.Set:
		err := setCanaryWeight(ce, mesh, logger, configServiceURL, ce.Canary.Value)
		if err != nil {
			return err
		}
		err = ApplyConfiguration(ce, true, logger, configServiceURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyCanaryToPrimary(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string) error {

	// TODO: Implement this function
	// TODO: In generatedChartHander adapt service
	return nil
}

func setCanaryWeight(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string,
	canaryWeight int32) error {

	// Read chart
	chart, err := helm.GetChart(ce.Project, ce.Service, ce.Stage, helm.GetChartName(ce.Service, true), configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when reading chart %s from project %s: %s", helm.GetChartName(ce.Service, true), ce.Project, err.Error()))
		return err
	}

	// Set weights in all virtualservices
	for _, template := range chart.Templates {
		if strings.HasPrefix(template.Name, "templates/") &&
			strings.HasSuffix(template.Name, mesh.GetVirtualServiceSuffix()) {

			vs, err := mesh.UpdateWeights(template.Data, canaryWeight)
			if err != nil {
				logger.Error(fmt.Sprintf("Error when setting new weights in VirtualService %s from chart %s and project %s: %s",
					template.Name, helm.GetChartName(ce.Service, true), ce.Project, err.Error()))
				return err
			}
			template.Data = vs
		}
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when packaging modified chart %s from project %s: %s", helm.GetChartName(ce.Service, true), ce.Project, err.Error()))
		return err
	}
	err = helm.StoreChart(ce.Project, ce.Service, ce.Stage, helm.GetChartName(ce.Service, true), chartData, configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when storing modified chart %s from project %s: %s", helm.GetChartName(ce.Service, true), ce.Project, err.Error()))
		return err
	}
	return nil
}

// DeleteRelease deletes a helm reease
func DeleteRelease(ce *keptnevents.ConfigurationChangeEventData, generated bool, logger *keptnutils.Logger) error {
	releaseName := helm.GetReleaseName(ce.Project, ce.Service, ce.Stage, generated)
	if _, err := keptnutils.ExecuteCommand("helm", []string{"delete", "--purge", releaseName}); err != nil {
		logger.Error(fmt.Sprintf("Error when deleting release %s: %s", releaseName, err.Error()))
		return err
	}
	return nil
}

// ApplyConfiguration applies the chart of the provided service.
// Furthermore, this function waits until all deployments in the namespace are ready.
func ApplyConfiguration(ce *keptnevents.ConfigurationChangeEventData, generated bool, logger *keptnutils.Logger, configServiceURL string) error {

	ch, err := helm.GetChart(ce.Project, ce.Service, ce.Stage, helm.GetChartName(ce.Service, generated), configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when reading chart %s: %s", helm.GetChartName(ce.Service, generated), err.Error()))
		return err
	}

	helmChartDir, err := ioutil.TempDir("", "")
	if err != nil {
		logger.Error(fmt.Sprintf("Error when creating temporary directory: %s", err.Error()))
		return err
	}
	defer os.RemoveAll(helmChartDir)

	chartPath, err := chartutil.Save(ch, helmChartDir)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when saving chart into temporary directory %s: %s", helmChartDir, err.Error()))
		return err
	}

	releaseName := helm.GetReleaseName(ce.Project, ce.Service, ce.Stage, generated)
	namespace := helm.GetNamespace(ce.Project, ce.Stage, generated)
	if _, err := keptnutils.ExecuteCommand("helm", []string{"upgrade", "--install", releaseName,
		chartPath, "--namespace", namespace, "--wait"}); err != nil {
		logger.Error(fmt.Sprintf("Error when upgrading chart %s in namespace %s: %s",
			releaseName, namespace, err.Error()))
		return err
	}

	useInClusterConfig := false
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	}

	if err := keptnutils.WaitForDeploymentsInNamespace(useInClusterConfig, namespace); err != nil {
		logger.Error(fmt.Sprintf("Error when waiting for deployments in namespace %s: %s", namespace, err.Error()))
		return err
	}
	logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	return nil
}
