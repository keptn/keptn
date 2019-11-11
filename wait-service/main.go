package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

const timeout = 60
const eventbroker = "EVENTBROKER"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type Client struct {
	httpClient *http.Client
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
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
		log.Fatalf("failed to create transport: %v", err)
	}

	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func newClient() *Client {
	client := Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	return &client
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "wait-service")

	data := &keptnevents.DeploymentFinishedEventData{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if event.Type() != keptnevents.DeploymentFinishedEventType {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	go waitDuration(event, shkeptncontext, *data, logger)

	return nil
}

// waitDuration just waits for a the time defined in environment variable WAIT_DURATION
func waitDuration(event cloudevents.Event, shkeptncontext string, data keptnevents.DeploymentFinishedEventData, logger *keptnutils.Logger) {

	startedAt := time.Now()

	switch strings.ToLower(data.TestStrategy) {
	case "real-user":
		duration, err := retrieveDuration("WAIT_DURATION")
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err.Error()))
			duration = 0
		}
		logger.Debug(fmt.Sprintf("Start to wait %d seconds.", duration))
		time.Sleep(time.Duration(duration) * time.Second)
		logger.Debug(fmt.Sprintf("Waiting %d seconds is over.", duration))

		if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, logger); err != nil {
			logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". ")
		}
	case "":
		logger.Info("No test strategy specified, hence no tests are triggered. ")
	default:
		logger.Error(fmt.Sprintf("Unknown test strategy '%s'. ", data.TestStrategy))
	}
}

// retrieveDuration reads the definition of the duration from environment variable WAIT_TIME.
// Then converts the value, which can have unit hour [h], minute [m], or second[s], into seconds.
func retrieveDuration(environmentVariable string) (int, error) {
	durationStr := os.Getenv(environmentVariable)
	if durationStr == "" {
		return 0, fmt.Errorf("Failed to retrieve value from  environment variable: %s", environmentVariable)
	}

	if strings.Contains(durationStr, "s") {
		durationStr = strings.TrimSuffix(durationStr, "s")
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return 0, fmt.Errorf("Failed to convert value %s into integer", durationStr)
		}
		return duration, nil
	} else if strings.Contains(durationStr, "m") {
		durationStr = strings.TrimSuffix(durationStr, "m")
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return 0, fmt.Errorf("Failed to convert value %s into integer", durationStr)
		}
		return duration * 60, nil

	} else if strings.Contains(durationStr, "h") {
		durationStr = strings.TrimSuffix(durationStr, "h")
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return 0, fmt.Errorf("Failed to convert value %s into integer", durationStr)
		}
		return duration * 60 * 60, nil
	}

	return 0, fmt.Errorf("Value of environment variable: %s not correct. Please set value based on the pattern: [duration][unit] e.g.: 1h, 10m, 50s", environmentVariable)
}

// getServiceEndpoint retrieves an endpoint stored in an environment variable and sets http as default scheme
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

// sendTestsFinishedEvent sends a Cloud Event of type sh.keptn.events.tests-finished to the event broker
func sendTestsFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event, startedAt time.Time, logger *keptnutils.Logger) error {

	source, _ := url.Parse("wait-service")
	contentType := "application/json"

	var testFinishedData interface{}
	if err := incomingEvent.DataAs(&testFinishedData); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	testFinishedData.(map[string]interface{})["startedat"] = startedAt

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.TestsFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: testFinishedData,
	}

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
		return errors.New("Failed to create transport: " + err.Error())
	}

	client, err := client.New(transport)
	if err != nil {
		return errors.New("Failed to create HTTP client: " + err.Error())
	}

	logger.Debug("Send sh.keptn.events.tests-finished event to event broker")
	if _, err := client.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}
