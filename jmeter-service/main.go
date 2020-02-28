package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "jmeter-service")

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

	if data.TestStrategy == "real-user" {
		logger.Info("Received 'real-user' test strategy, hence no tests are triggered")
		return nil
	}
	go runTests(event, shkeptncontext, *data, logger)

	return nil
}

func runTests(event cloudevents.Event, shkeptncontext string, data keptnevents.DeploymentFinishedEventData, logger *keptnutils.Logger) {

	testInfo := getTestInfo(data)
	id := uuid.New().String()
	startedAt := time.Now()

	var res bool
	var err error
	res, err = runHealthCheck(data, id, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if !res {
		if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, "fail", logger); err != nil {
			logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
		}
		return
	}
	logger.Info("Health Check test passed = " + strconv.FormatBool(res) + ". " + testInfo.ToString())

	var sendEvent = false

	switch strings.ToLower(data.TestStrategy) {
	case "functional":
		res, err = runFunctionalCheck(data, id, logger)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Info("Functional test passed = " + strconv.FormatBool(res) + ". " + testInfo.ToString())
		sendEvent = true

	case "performance":
		res, err = runPerformanceCheck(data, id, logger)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Info("Performance test passed = " + strconv.FormatBool(res) + ". " + testInfo.ToString())
		sendEvent = true

	case "":
		logger.Info("No test strategy specified, hence no tests are triggered. " + testInfo.ToString())
		sendEvent = true

	default:
		logger.Error("Unknown test strategy '" + data.TestStrategy + "'" + ". " + testInfo.ToString())
	}

	if sendEvent {
		if !res {
			if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, "fail", logger); err != nil {
				logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
			}
			return
		}
		if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, "pass", logger); err != nil {
			logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
		}
	}
}

func getTestInfo(data keptnevents.DeploymentFinishedEventData) *TestInfo {
	return &TestInfo{
		Project:      data.Project,
		Service:      data.Service,
		Stage:        data.Stage,
		TestStrategy: data.TestStrategy,
	}
}

func getServiceURL(data keptnevents.DeploymentFinishedEventData) (*url.URL, error) {

	if data.DeploymentURILocal != "" {
		return url.Parse(data.DeploymentURILocal)

	} else if data.DeploymentURIPublic != "" {
		return url.Parse(data.DeploymentURIPublic)
	}

	// Use educated guess of the service url based on stage, service name, deployment type
	serviceURL := data.Service + "." + data.Project + "-" + data.Stage
	if data.DeploymentStrategy == "blue_green_service" {
		serviceURL = data.Service + "-canary" + "." + data.Project + "-" + data.Stage
	}
	serviceURL = "http://" + serviceURL + "/health"
	return url.Parse(serviceURL)
}

func runHealthCheck(data keptnevents.DeploymentFinishedEventData, id string, logger *keptnutils.Logger) (bool, error) {
	os.RemoveAll("HealthCheck_" + data.Service)
	os.RemoveAll("HealthCheck_" + data.Service + "_result.tlf")
	os.RemoveAll("output.txt")

	testInfo := getTestInfo(data)
	url, err := getServiceURL(data)
	if err != nil {
		return false, err
	}
	return executeJMeter(testInfo, "jmeter/basiccheck.jmx", "HealthCheck_"+data.Service, url,
		1, 1, 250, "HealthCheck_"+id,
		true, 0, logger)
}

func runFunctionalCheck(data keptnevents.DeploymentFinishedEventData, id string, logger *keptnutils.Logger) (bool, error) {

	os.RemoveAll("FuncCheck_" + data.Service)
	os.RemoveAll("FuncCheck_" + data.Service + "_result.tlf")
	os.RemoveAll("output.txt")

	testInfo := getTestInfo(data)
	url, err := getServiceURL(data)
	if err != nil {
		return false, err
	}
	return executeJMeter(testInfo, "jmeter/load.jmx", "FuncCheck_"+data.Service, url,
		1, 1, 250, "FuncCheck_"+id, true, 0, logger)
}

func runPerformanceCheck(data keptnevents.DeploymentFinishedEventData, id string, logger *keptnutils.Logger) (bool, error) {

	os.RemoveAll("PerfCheck_" + data.Service)
	os.RemoveAll("PerfCheck_" + data.Service + "_result.tlf")
	os.RemoveAll("output.txt")

	testInfo := getTestInfo(data)
	url, err := getServiceURL(data)
	if err != nil {
		return false, err
	}
	return executeJMeter(testInfo, "jmeter/load.jmx", "PerfCheck_"+data.Service, url,
		10, 500, 250, "PerfCheck_"+id,
		false, 0, logger)
}

func getGatewayFromConfigmap() (string, error) {

	api, err := keptnutils.GetKubeAPI(true)
	if err != nil {
		return "", err
	}

	cm, err := api.ConfigMaps("keptn").Get("keptn-domain", metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(cm.Data["app_domain"]), nil
}

func sendTestsFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event, startedAt time.Time, result string, logger *keptnutils.Logger) error {
	source, _ := url.Parse("jmeter-service")
	contentType := "application/json"

	testFinishedData := keptnevents.TestsFinishedEventData{}
	// fill in data from incoming event (e.g., project, service, stage, teststrategy, deploymentstrategy)
	if err := incomingEvent.DataAs(&testFinishedData); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	// fill in timestamps
	testFinishedData.Start = startedAt.Format(time.RFC3339)
	testFinishedData.End = time.Now().Format(time.RFC3339)
	// set test result
	testFinishedData.Result = result

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

	return sendEvent(event)
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
