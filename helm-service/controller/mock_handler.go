package controller

import (
	"errors"
	"fmt"
	keptnevents "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

// MockedHandler mocks typical tasks of a handler
type MockedHandler struct {
	keptnHandler                        *keptnv2.Keptn
	helmExecutor                        helm.HelmExecutor
	configServiceURL                    string
	options                             MockedHandlerOptions
	sentCloudEvents                     []cloudevents.Event
	handledErrorEvents                  []interface{}
	upgradeChartInvocations             []upgradeChartData
	upgradeChartWithReplicasInvocations []upgradeChartWithReplicaData
}

// MockedHandlerOption function is used to configure the mock
type MockedHandlerOption func(*MockedHandlerOptions)

// MockedHandlerOptions contains configuration items for the mock
type MockedHandlerOptions struct {
	SendEventBehavior func(eventType string) bool
}

func sendingEventSucceeds(eventType string) bool {
	return true
}

// NewMockedHandler creates a MockedHandler
func NewMockedHandler(keptnHandler *keptnv2.Keptn, configServiceURL string, options ...MockedHandlerOption) *MockedHandler {
	helmExecutor := helm.NewHelmMockExecutor()
	opt := MockedHandlerOptions{
		SendEventBehavior: sendingEventSucceeds,
	}

	for _, o := range options {
		o(&opt)
	}

	return &MockedHandler{
		keptnHandler:            keptnHandler,
		helmExecutor:            helmExecutor,
		configServiceURL:        configServiceURL,
		options:                 opt,
		upgradeChartInvocations: []upgradeChartData{},
	}
}

func (h *MockedHandler) getKeptnHandler() *keptnv2.Keptn {
	return h.keptnHandler
}

func (h *MockedHandler) getHelmExecutor() helm.HelmExecutor {
	return h.helmExecutor
}

func (h *MockedHandler) getConfigServiceURL() string {
	return h.configServiceURL
}

func (h *MockedHandler) getGeneratedChart(e keptnv2.EventData) (*chart.Chart, string, error) {
	ch := helm.GetTestGeneratedChart()
	return &ch, "GENERATED_CHART_GIT_ID", nil
}

func (h *MockedHandler) getUserChart(e keptnv2.EventData) (*chart.Chart, string, error) {
	ch := helm.GetTestUserChart()
	return &ch, "USER_CHART_GIT_ID", nil
}

func (h *MockedHandler) existsGeneratedChart(e keptnv2.EventData) (bool, error) {
	return true, nil
}

// HandleError logs the error and sends a finished-event
func (h *MockedHandler) handleError(triggerID string, err error, taskName string, finishedEventData interface{}) {
	fmt.Println("HandleError: " + err.Error())
	h.handledErrorEvents = append(h.handledErrorEvents, finishedEventData)
}

func (h *MockedHandler) sendEvent(triggerID, ceType string, data interface{}) error {
	if h.options.SendEventBehavior != nil && !h.options.SendEventBehavior(ceType) {
		return errors.New("Failed at sending event of type " + ceType)
	}

	event := cloudevents.NewEvent()
	event.SetType(ceType)
	event.SetSource("helm-service")
	event.SetDataContentType(cloudevents.ApplicationJSON)

	event.SetExtension("triggeredid", triggerID)
	event.SetExtension("shkeptncontext", h.keptnHandler.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, data)

	fmt.Println("Send Event: ")
	fmt.Println(event.String())

	h.sentCloudEvents = append(h.sentCloudEvents, event)
	return nil
}

type upgradeChartData struct {
	ch       *chart.Chart
	event    keptnv2.EventData
	strategy keptnevents.DeploymentStrategy
}

type upgradeChartWithReplicaData struct {
	upgradeChartData
	replicas int
}

func (h *MockedHandler) upgradeChart(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy) error {
	ucd := upgradeChartData{
		ch:       ch,
		event:    event,
		strategy: strategy,
	}
	h.upgradeChartInvocations = append(h.upgradeChartInvocations, ucd)

	return nil
}

func (h *MockedHandler) upgradeChartWithReplicas(ch *chart.Chart, event keptnv2.EventData,
	strategy keptnevents.DeploymentStrategy, replicas int) error {
	ucdr := upgradeChartWithReplicaData{
		upgradeChartData: upgradeChartData{
			ch:       ch,
			event:    event,
			strategy: strategy,
		},
		replicas: replicas,
	}
	h.upgradeChartWithReplicasInvocations = append(h.upgradeChartWithReplicasInvocations, ucdr)

	return nil
}
