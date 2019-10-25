package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/events"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {

	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {

	switch event.Type() {
	case keptnevents.TestsFinishedEventType:
		return startEvaluation(event) // backwards compatibility to Keptn versions <= 0.5.x
	case keptnevents.StartEvaluationEventType:
		return startEvaluation(event) // new event type in Keptn versions >= 0.6
	case keptnevents.InternalGetSLIDoneEventType:
		return evaluateSLIValues(event)
	default:
		return errors.New("received unknown event type")
	}
}

func evaluateSLIValues(event cloudevents.Event) error {
	// get results of previous evaluations from data store (mongodb-datastore.keptn-datastore.svc.cluster.local)

	// compare the results based on the evaluation strategy

	// send the evaluation-done-event
	return nil
}

func startEvaluation(event cloudevents.Event) error {
	// get the SLI provider that has been configured for the project (e.g. 'dynatrace' or 'prometheus')

	// send a new event to trigger the SLI retrieval
	return nil
}

func sendEvaluationDoneEvent(shkeptncontext string, project string,
	service string, nextStage string, image string) error {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	valuesCanary := make(map[string]interface{})
	valuesCanary["image"] = image
	canary := keptnevents.Canary{Action: keptnevents.Set, Value: 100}
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:      project,
		Service:      service,
		Stage:        nextStage,
		ValuesCanary: valuesCanary,
		Canary:       &canary,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.EvaluationDoneEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return sendEvent(event)
}

func sendInternalGetSLIEvent(shkeptncontext string, project string,
	service string, stage string, sliProvider string) error {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	getSLIEvent := keptnevents.InternalGetSLIEventData{
		SLIProvider: sliProvider,
		Project:     project,
		Service:     service,
		Stage:       stage,
		Indicators:  []string{},
	}
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.InternalGetSLIEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: getSLIEvent,
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
