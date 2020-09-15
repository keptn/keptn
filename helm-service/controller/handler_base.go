package controller

import (
	"fmt"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"
)

type HandlerBase struct {
	keptnHandler     *keptnv2.Keptn
	helmExecutor     helm.HelmExecutor
	configServiceURL string
}

func NewHandlerBase(keptnHandler *keptnv2.Keptn, configServiceURL string) *HandlerBase {
	helmExecutor := helm.NewHelmV3Executor(keptnHandler.Logger)
	return &HandlerBase{
		keptnHandler:     keptnHandler,
		helmExecutor:     helmExecutor,
		configServiceURL: configServiceURL,
	}
}

func (h HandlerBase) GetKeptnHandler() *keptnv2.Keptn {
	return h.keptnHandler
}

func (h HandlerBase) GetHelmExecutor() helm.HelmExecutor {
	return h.helmExecutor
}

func (h HandlerBase) GetConfigServiceURL() string {
	return h.configServiceURL
}

func (h HandlerBase) GetGeneratedChart(e keptnv2.EventData) (*chart.Chart, error) {
	helmChartName := helm.GetChartName(e.Service, true)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, h.configServiceURL)
}

func (h HandlerBase) GetUserChart(e keptnv2.EventData) (*chart.Chart, error) {
	helmChartName := helm.GetChartName(e.Service, false)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, h.configServiceURL)
}

func (h HandlerBase) ExistsGeneratedChart(e keptnv2.EventData) (bool, error) {
	genChartName := helm.GetChartName(e.Service, true)
	return helm.DoesChartExist(e, genChartName, h.GetConfigServiceURL())
}

// HandleError logs the error and sends a finished-event
func (h HandlerBase) HandleError(triggerId string, err error, taskName string, finishedEventData interface{}) {
	h.keptnHandler.Logger.Error(err.Error())
	if err := h.SendEvent(triggerId, keptnv2.GetFinishedEventType(taskName), finishedEventData); err != nil {
		h.keptnHandler.Logger.Error(err.Error())
	}
}

func (h HandlerBase) SendEvent(triggerId, ceType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(ceType)
	event.SetSource("helm-service")
	event.SetDataContentType(cloudevents.ApplicationJSON)

	event.SetExtension("triggeredid", triggerId)
	event.SetExtension("shkeptncontext", h.keptnHandler.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, data)
	return h.keptnHandler.SendCloudEvent(event)
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

func (h HandlerBase) upgradeChart(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	releasename := helm.GetReleaseName(event.Project, event.Stage, event.Service, generated)
	namespace := event.Project + "-" + event.Stage

	h.keptnHandler.Logger.Info(fmt.Sprintf("Start upgrading release %s in namespace %s", releasename, namespace))

	return h.helmExecutor.UpgradeChart(ch, releasename, namespace,
		getKeptnValues(event.Project, event.Stage, event.Service,
			getDeploymentName(strategy, generated)))
}

func (h HandlerBase) upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy, replicas int) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	releasename := helm.GetReleaseName(event.Project, event.Stage, event.Service, generated)
	namespace := event.Project + "-" + event.Stage

	h.keptnHandler.Logger.Info(fmt.Sprintf("Start upgrading release %s in namespace %s with replicas %d",
		releasename, namespace, replicas))

	return h.helmExecutor.UpgradeChart(ch, releasename, namespace,
		addReplicas(getKeptnValues(event.Project, event.Stage, event.Service,
			getDeploymentName(strategy, generated)), replicas))
}

func getKeptnValues(project, stage, service, deploymentName string) map[string]interface{} {
	return map[string]interface{}{
		"keptn": map[string]interface{}{
			"project":    project,
			"stage":      stage,
			"service":    service,
			"deployment": deploymentName,
		},
	}
}

func addReplicas(vals map[string]interface{}, replicas int) map[string]interface{} {
	vals["replicaCount"] = replicas
	return vals
}
