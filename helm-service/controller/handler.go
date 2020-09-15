package controller

import (
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

type Handler interface {
	GetKeptnHandler() *keptnv2.Keptn
	GetHelmExecutor() helm.HelmExecutor
	GetConfigServiceURL() string

	GetGeneratedChart(e keptnv2.EventData) (*chart.Chart, error)
	GetUserChart(e keptnv2.EventData) (*chart.Chart, error)
	ExistsGeneratedChart(e keptnv2.EventData) (bool, error)
	HandleError(triggerId string, err error, taskName string, finishedEventData interface{})
	SendEvent(triggerId, ceType string, data interface{}) error
	upgradeChart(ch *chart.Chart, event keptnv2.EventData,
		strategy keptnevents.DeploymentStrategy) error
	upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
		strategy keptnevents.DeploymentStrategy, replicas int) error
}
