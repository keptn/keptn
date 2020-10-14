package controller

import (
	"fmt"

	keptnevents "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

// MockHandler mocks typical tasks of a handler
type MockHandler struct {
	keptnHandler     *keptnv2.Keptn
	helmExecutor     helm.HelmExecutor
	configServiceURL string
}

// NewMockHandler creates a MockHandler
func NewMockHandler(keptnHandler *keptnv2.Keptn, configServiceURL string) *MockHandler {
	helmExecutor := helm.NewHelmMockExecutor()
	return &MockHandler{
		keptnHandler:     keptnHandler,
		helmExecutor:     helmExecutor,
		configServiceURL: configServiceURL,
	}
}

func (h *MockHandler) getKeptnHandler() *keptnv2.Keptn {
	return h.keptnHandler
}

func (h *MockHandler) getHelmExecutor() helm.HelmExecutor {
	return h.helmExecutor
}

func (h *MockHandler) getConfigServiceURL() string {
	return h.configServiceURL
}

func (h *MockHandler) getGeneratedChart(e keptnv2.EventData) (*chart.Chart, error) {
	ch := helm.GetTestGeneratedChart()
	return &ch, nil
}

func (h *MockHandler) getUserChart(e keptnv2.EventData) (*chart.Chart, error) {
	ch := helm.GetTestUserChart()
	return &ch, nil
}

func (h *MockHandler) existsGeneratedChart(e keptnv2.EventData) (bool, error) {
	return true, nil
}

// HandleError logs the error and sends a finished-event
func (h *MockHandler) handleError(triggerID string, err error, taskName string, finishedEventData interface{}) {
	fmt.Println("HandleError: " + err.Error())
}

var sentCloudEvents []cloudevents.Event

func (h *MockHandler) sendEvent(triggerID, ceType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(ceType)
	event.SetSource("helm-service")
	event.SetDataContentType(cloudevents.ApplicationJSON)

	event.SetExtension("triggeredid", triggerID)
	event.SetExtension("shkeptncontext", h.keptnHandler.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, data)

	fmt.Println("Send Event: ")
	fmt.Println(event.String())

	sentCloudEvents = append(sentCloudEvents, event)
	return nil
}

func (h *MockHandler) upgradeChart(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy) error {

	return nil
}

func (h *MockHandler) upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy, replicas int) error {

	return nil
}
