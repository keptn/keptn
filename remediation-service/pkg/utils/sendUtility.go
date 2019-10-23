package utils

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const eventbroker = "EVENTBROKER"

// CreateAndSendConfigurationChangedEvent creates ConfigurationChangeEvent and sends it
func CreateAndSendConfigurationChangedEvent(problem *keptnevents.ProblemEventData,
	shkeptncontext string, changedTemplates []*chart.Template) error {

	source, _ := url.Parse("https://github.com/keptn/keptn/remediation-service")
	contentType := "application/json"

	changedFiles := make(map[string]string)
	for _, template := range changedTemplates {
		changedFiles[template.Name] = string(template.Data)
	}

	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:                   problem.Project,
		Service:                   problem.Service,
		Stage:                     problem.Stage,
		FileChangesGeneratedChart: changedFiles,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return sendEvent(event)
}

// SendTestsFinishedEvent sends a Cloud Event of type sh.keptn.events.tests-finished to the event broker
func SendTestsFinishedEvent(shkeptncontext string, project string, stage string, service string) error {

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	testFinishedData := keptnevents.TestsFinishedEventData{}
	testFinishedData.Project = project
	testFinishedData.Stage = stage
	testFinishedData.Service = service
	testFinishedData.TestStrategy = "real-user"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        "sh.keptn.events.tests-finished",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: testFinishedData,
	}

	return sendEvent(event)
}

func sendEvent(event cloudevents.Event) error {
	endPoint, err := getServiceEndpoint(eventbroker)
	if err != nil {
		return errors.New("Failed to retrieve endpoint of eventbroker. %s" + err.Error())
	}

	if endPoint.Host == "" {
		return errors.New("Host of eventbroker not set")
	}

	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(endPoint.String()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(transport)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

// getServiceEndpoint gets an endpoint stored in an environment variable and sets http as default scheme
func getServiceEndpoint(service string) (url.URL, error) {
	url, err := url.Parse(os.Getenv(service))
	if err != nil {
		return *url, fmt.Errorf("Failed to retrieve value from ENVIRONMENT_VARIABLE: %s", service)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	return *url, nil
}
