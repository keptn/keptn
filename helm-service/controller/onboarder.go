package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/keptn/keptn/helm-service/pkg/namespacemanager"
	"github.com/keptn/keptn/helm-service/pkg/types"
	"time"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"helm.sh/helm/v3/pkg/chart"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/keptn/keptn/helm-service/pkg/helm"
)

// Onboarder is responsible for onboarding a service
type Onboarder interface {

	// OnboardGeneratedChart generates the generated chart using the Helm manifests of the user chart
	// as well as the specified deployment strategy
	OnboardGeneratedChart(helmManifest string, event keptnv2.EventData, strategy keptnevents.DeploymentStrategy) (*chart.Chart, error)

	// OnboardService
	OnboardService(stageName string, event *keptnv2.ServiceCreateFinishedEventData) error
}

// onboarder is an implemntation of Onboarder
type onboarder struct {
	Handler
	namespaceManager namespacemanager.INamespaceManager
	serviceHandler   types.IServiceHandler
	chartStorer      types.IChartStorer
	chartGenerator   helm.ChartGenerator
	chartPackager    types.IChartPackager
}

// NewOnboarder creates a new onboarder instance
func NewOnboarder(
	keptnHandler *keptnv2.Keptn,
	namespaceManager namespacemanager.INamespaceManager,
	serviceHandler types.IServiceHandler,
	chartStorer types.IChartStorer,
	chartGenerator helm.ChartGenerator,
	chartPackager types.IChartPackager,
	configServiceURL string) Onboarder {

	return &onboarder{
		Handler:          NewHandlerBase(keptnHandler, configServiceURL),
		namespaceManager: namespaceManager,
		serviceHandler:   serviceHandler,
		chartStorer:      chartStorer,
		chartGenerator:   chartGenerator,
		chartPackager:    chartPackager,
	}
}

func (o *onboarder) OnboardService(stageName string, event *keptnv2.ServiceCreateFinishedEventData) error {

	const retries = 2
	var err error
	for i := 0; i < retries; i++ {
		_, err = o.serviceHandler.GetService(event.Project, stageName, event.Service)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return err
	}

	helmChartData, err := base64.StdEncoding.DecodeString(event.Helm.Chart)
	if err != nil {
		o.getKeptnHandler().Logger.Error("Error when decoding the Helm Chart")
		return err
	}

	o.getKeptnHandler().Logger.Debug("Storing the Helm Chart provided by the user in stage " + stageName)

	storeOpts := keptnutils.StoreChartOptions{
		Project:   event.Project,
		Service:   event.Service,
		Stage:     stageName,
		ChartName: helm.GetChartName(event.Service, false),
		HelmChart: helmChartData,
	}

	if _, err := o.chartStorer.Store(storeOpts); err != nil {
		o.getKeptnHandler().Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return err
	}
	return nil
}

// OnboardGeneratedChart generates the generated chart using the Helm manifests of the user chart
// as well as the specified deployment strategy
func (o *onboarder) OnboardGeneratedChart(helmManifest string, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy) (*chart.Chart, error) {

	helmChartName := helm.GetChartName(event.Service, true)
	o.getKeptnHandler().Logger.Debug(fmt.Sprintf("Generating the Keptn-managed Helm Chart %s for stage %s", helmChartName, event.Stage))

	var generatedChart *chart.Chart
	var err error
	if strategy == keptnevents.Duplicate {
		o.getKeptnHandler().Logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, "+
			"a chart for a duplicate deployment strategy is generated", event.Service, event.Stage, strategy.String()))
		generatedChart, err = o.chartGenerator.GenerateDuplicateChart(helmManifest, event.Project, event.Stage, event.Service)
		if err != nil {
			o.getKeptnHandler().Logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
		// inject Istio to the namespace for blue-green deployments
		if err := o.namespaceManager.InjectIstio(event.Project, event.Stage); err != nil {
			return nil, err
		}
	} else {
		o.getKeptnHandler().Logger.Debug(fmt.Sprintf("For service %s in stage %s with deployment strategy %s, a mesh chart is generated",
			event.Service, event.Stage, strategy.String()))
		generatedChart, err = o.chartGenerator.GenerateMeshChart(helmManifest, event.Project, event.Stage, event.Service)
		if err != nil {
			o.getKeptnHandler().Logger.Error("Error when generating the managed chart: " + err.Error())
			return nil, err
		}
	}

	o.getKeptnHandler().Logger.Debug(fmt.Sprintf("Storing the Keptn-generated Helm Chart %s for stage %s", helmChartName, event.Stage))
	generatedChartData, err := o.chartPackager.Package(generatedChart)
	if err != nil {
		o.getKeptnHandler().Logger.Error("Error when packing the managed chart: " + err.Error())
		return nil, err
	}

	storeOpts := keptnutils.StoreChartOptions{
		Project:   event.Project,
		Service:   event.Service,
		Stage:     event.Stage,
		ChartName: helmChartName,
		HelmChart: generatedChartData,
	}

	if _, err := o.chartStorer.Store(storeOpts); err != nil {
		o.getKeptnHandler().Logger.Error("Error when storing the Helm Chart: " + err.Error())
		return nil, err
	}
	return generatedChart, nil
}
