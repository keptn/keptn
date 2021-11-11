package controller

import (
	"fmt"
	"github.com/ghodss/yaml"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"
)

// HandlerBase provides basic functionality for all handlers
type HandlerBase struct {
	keptnHandler     *keptnv2.Keptn
	helmExecutor     helm.HelmExecutor
	configServiceURL string
}

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var GracefulShutdownKey = gracefulShutdownKeyType{}

// NewHandlerBase creates a new HandlerBase
func NewHandlerBase(keptnHandler *keptnv2.Keptn, helmExecutor helm.HelmExecutor, configServiceURL string) *HandlerBase {

	return &HandlerBase{
		keptnHandler:     keptnHandler,
		helmExecutor:     helmExecutor,
		configServiceURL: configServiceURL,
	}
}

func (h *HandlerBase) getKeptnHandler() *keptnv2.Keptn {
	return h.keptnHandler
}

func (h *HandlerBase) getHelmExecutor() helm.HelmExecutor {
	return h.helmExecutor
}

func (h *HandlerBase) getConfigServiceURL() string {
	return h.configServiceURL
}

func (h *HandlerBase) getGeneratedChart(e keptnv2.EventData) (*chart.Chart, string, error) {
	helmChartName := helm.GetChartName(e.Service, true)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, h.configServiceURL)
}

func (h *HandlerBase) getUserChart(e keptnv2.EventData) (*chart.Chart, string, error) {
	helmChartName := helm.GetChartName(e.Service, false)
	// Read chart
	return keptnutils.GetChart(e.Project, e.Service, e.Stage, helmChartName, h.configServiceURL)
}

func (h *HandlerBase) getUserManagedEndpoints(event keptnv2.EventData) (*keptnv2.Endpoints, error) {
	endpointsResource, err := h.getKeptnHandler().ResourceHandler.GetServiceResource(event.Project, event.Stage, event.Service, "helm/endpoints.yaml")
	if err != nil {
		// do not fail if the resource is not available
		if err == keptnapi.ResourceNotFoundError {
			return nil, nil
		}
		return nil, fmt.Errorf("Could not fetch endpoints resource: %s", err.Error())
	}
	if endpointsResource == nil {
		return nil, nil
	}
	endpoints := &keptnv2.Endpoints{}
	err = yaml.Unmarshal([]byte(endpointsResource.ResourceContent), endpoints)
	if err != nil {
		return nil, fmt.Errorf("could not parse endpoints.yaml: %s", err.Error())
	}
	return endpoints, nil
}

func (h *HandlerBase) existsGeneratedChart(e keptnv2.EventData) (bool, error) {
	genChartName := helm.GetChartName(e.Service, true)
	return helm.DoesChartExist(e, genChartName, h.getConfigServiceURL())
}

// HandleError logs the error and sends a finished-event
func (h *HandlerBase) handleError(triggerID string, err error, taskName string, finishedEventData interface{}) {
	h.keptnHandler.Logger.Error(err.Error())
	if err := h.sendEvent(triggerID, keptnv2.GetFinishedEventType(taskName), finishedEventData); err != nil {
		h.keptnHandler.Logger.Error(err.Error())
	}
}

func (h *HandlerBase) sendEvent(triggerID, ceType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(ceType)
	event.SetSource("helm-service")
	event.SetDataContentType(cloudevents.ApplicationJSON)

	event.SetExtension("triggeredid", triggerID)
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
	} else if strategy == keptnevents.UserManaged {
		return "user_managed"
	}
	return ""
}

func (h *HandlerBase) upgradeChart(ch *chart.Chart, event keptnv2.EventData, strategy keptnevents.DeploymentStrategy) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	releasename := helm.GetReleaseName(event.Project, event.Stage, event.Service, generated)
	namespace := event.Project + "-" + event.Stage

	return h.helmExecutor.UpgradeChart(ch, releasename, namespace,
		getKeptnValues(event.Project, event.Stage, event.Service, getDeploymentName(strategy, generated)))
}

func (h *HandlerBase) upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy, replicas int) error {
	generated := strings.HasSuffix(ch.Name(), "-generated")
	releasename := helm.GetReleaseName(event.Project, event.Stage, event.Service, generated)
	namespace := event.Project + "-" + event.Stage

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
