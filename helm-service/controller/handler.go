package controller

import (
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

// Handler provides methods for handling received Keptn events
type Handler interface {
	getKeptnHandler() *keptnv2.Keptn
	getHelmExecutor() helm.HelmExecutor
	getConfigServiceURL() string

	getGeneratedChart(e keptnv2.EventData) (*chart.Chart, error)
	getUserChart(e keptnv2.EventData) (*chart.Chart, error)
	existsGeneratedChart(e keptnv2.EventData) (bool, error)
	handleError(triggerId string, err error, taskName string, finishedEventData interface{})
	sendEvent(triggerId, ceType string, data interface{}) error
	upgradeChart(ch *chart.Chart, event keptnv2.EventData,
		strategy keptnevents.DeploymentStrategy) error
	upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
		strategy keptnevents.DeploymentStrategy, replicas int) error
}
