package controller

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
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
	HandleEvent(ce cloudevents.Event)
	getGeneratedChart(e keptnv2.EventData, commitID string) (*chart.Chart, string, error)
	getUserChart(e keptnv2.EventData, commitID string) (*chart.Chart, string, error)
	existsGeneratedChart(e keptnv2.EventData, commitID string) (bool, error)
	handleError(triggerID string, err error, taskName string, finishedEventData interface{})
	sendEvent(triggerID, ceType string, data interface{}) error
	upgradeChart(ch *chart.Chart, event keptnv2.EventData,
		strategy keptnevents.DeploymentStrategy) error
	upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
		strategy keptnevents.DeploymentStrategy, replicas int) error
	getUserManagedEndpoints(event keptnv2.EventData) (*keptnv2.Endpoints, error)
}
