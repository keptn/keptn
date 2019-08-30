package controller

import (
	"fmt"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// ChangeConfiguration
func ChangeConfiguration(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string) error {

	if len(ce.ValuesPrimary) > 0 {
		changeValues(ce, ce.Service+"-generated", ce.ValuesPrimary, logger, configServiceURL)
	}
	if len(ce.ValuesCanary) > 0 {
		changeValues(ce, ce.Service, ce.ValuesCanary, logger, configServiceURL)
	}

	// Change canary
	if ce.Canary != nil {
		changeCanary(ce, mesh, logger, configServiceURL)
	}

	return nil
}

func changeValues(ce *keptnevents.ConfigurationChangeEventData, chartName string, newValues map[string]*chart.Value,
	logger *keptnutils.Logger, configServiceURL string) error {

	// Read chart
	chart, err := helm.GetChart(ce.Project, ce.Service, ce.Stage, chartName, configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when reading chart %s from project %s: %s", chartName, ce.Project, err.Error()))
		return err
	}
	// Change values
	for k, v := range newValues {
		chart.Values.Values[k] = v
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when packaging modified chart %s from project %s: %s", chartName, ce.Project, err.Error()))
		return err
	}
	err = helm.StoreChart(ce.Project, ce.Service, ce.Stage, chartName, chartData, configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when storing modified chart %s from project %s: %s", chartName, ce.Project, err.Error()))
		return err
	}
	return nil
}

func changeCanary(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string) error {

	switch ce.Canary.Action {
	case keptnevents.Discard:
		return setCanaryWeight(ce, mesh, logger, configServiceURL, 0)
	case keptnevents.Promote:
		return copyCanaryToPrimary(ce, mesh, logger, configServiceURL)
	case keptnevents.Set:
		return setCanaryWeight(ce, mesh, logger, configServiceURL, ce.Canary.Value)
	}

	return nil
}

func copyCanaryToPrimary(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string) error {

	// Set 100% to primary
	return nil
}

func setCanaryWeight(ce *keptnevents.ConfigurationChangeEventData, mesh mesh.Mesh, logger *keptnutils.Logger, configServiceURL string,
	canaryWeight int32) error {

	chartName := ce.Service + "-generated"
	// Read chart
	chart, err := helm.GetChart(ce.Project, ce.Service, ce.Stage, chartName, configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when reading chart %s from project %s: %s", chartName, ce.Project, err.Error()))
		return err
	}

	// Set weights in all virtualservices
	for _, template := range chart.Templates {
		if strings.HasPrefix(template.Name, "templates/") &&
			strings.HasSuffix(template.Name, mesh.GetVirtualServiceSuffix()) {

			vs, err := mesh.UpdateWeights(template.Data, canaryWeight)
			if err != nil {
				logger.Error(fmt.Sprintf("Error when setting new weights in VirtualService %s from chart %s and project %s: %s",
					template.Name, chartName, ce.Project, err.Error()))
				return err
			}
			template.Data = vs
		}
	}

	// Store chart
	chartData, err := helm.PackageChart(chart)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when packaging modified chart %s from project %s: %s", chartName, ce.Project, err.Error()))
		return err
	}
	err = helm.StoreChart(ce.Project, ce.Service, ce.Stage, chartName, chartData, configServiceURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when storing modified chart %s from project %s: %s", chartName, ce.Project, err.Error()))
		return err
	}
	return nil
}
