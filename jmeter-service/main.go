package main

import (
	"context"
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/kelseyhightower/envconfig"
)

const eventbroker = "EVENTBROKER"
const configurationService = "CONFIGURATION_SERVICE"

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
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(os.Args[1:], env))
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptncommon.NewLogger(shkeptncontext, event.Context.GetID(), "jmeter-service")

	data := &keptnv2.TestTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if event.Type() != keptnv2.GetTriggeredEventType(keptnv2.TestTaskName) {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	if data.Test.TestStrategy == TestStrategy_RealUser {
		logger.Info("Received '" + TestStrategy_RealUser + "' test strategy, hence no tests are triggered")
		return nil
	}
	go runTests(event, shkeptncontext, *data, logger)

	return nil
}

//
// This method executes the correct tests based on the passed testStrategy in the deployment finished event
// The method will always try to execute a health check workload first, then execute the workload based on the passed testStrategy
//
func runTests(event cloudevents.Event, shkeptncontext string, data keptnv2.TestTriggeredEventData, logger *keptncommon.Logger) {

	sendTestsStartedEvent(shkeptncontext, event, logger)

	testInfo := getTestInfo(data, shkeptncontext)
	startedAt := time.Now()

	// load the workloads from JMeterConf
	var err error
	var jmeterconf *JMeterConf
	jmeterconf, err = getJMeterConf(testInfo.Project, testInfo.Stage, testInfo.Service, logger)

	// get the service endpoint we need to run the test against
	var serviceUrl *url.URL
	serviceUrl, err = getServiceURL(data)
	if err != nil {
		msg := fmt.Sprintf("Error getting service url to run test against: %s", err.Error())
		logger.Error(msg)
		if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, msg, logger); err != nil {
			logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
		}
		return
	}

	// first we run a health check workload. If that fails we stop the rest
	var healthCheckWorkload *Workload
	var res bool
	healthCheckWorkload, err = getWorkload(jmeterconf, TestStrategy_HealthCheck)
	if healthCheckWorkload != nil {
		res, err = runWorkload(serviceUrl, testInfo, healthCheckWorkload, logger)
		if err != nil {
			msg := fmt.Sprintf("could not run test workload: %s", err.Error())
			logger.Error(msg)
			if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, msg, logger); err != nil {
				logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
			}
			return
		}

		if !res {
			if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, keptnv2.ResultFailed, logger); err != nil {
				logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
			}
			return
		}
		logger.Info("Health Check test passed = " + strconv.FormatBool(res) + ". " + testInfo.ToString())
	} else {
		logger.Info("No Health Check test workload configuration found. Skipping Health Check")
	}

	// now lets execute the workload based on the passed testStrategy
	res = false
	var testStrategy = strings.ToLower(data.Test.TestStrategy)

	if testStrategy == "" {
		// no testStrategy passed at all -> we just send a successful test finished event!
		logger.Info("No testStrategy specified therefore skipping further test execution and sending back success")
		res = true
	} else {
		// lets get the workload configuration for the test strategy
		var teststrategyWorkload *Workload
		teststrategyWorkload, err = getWorkload(jmeterconf, testStrategy)
		if teststrategyWorkload != nil {
			res, err = runWorkload(serviceUrl, testInfo, teststrategyWorkload, logger)
			if err != nil {
				msg := fmt.Sprintf("could not run test workload: %s", err.Error())
				logger.Error(msg)
				if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, msg, logger); err != nil {
					logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
				}
				return
			} else {
				logger.Info(fmt.Sprintf("Tests for %s with status = %s.%s", testStrategy, strconv.FormatBool(res), testInfo.ToString()))
			}
		} else {
			// no workload for that test strategy!!
			res = false
			logger.Error(fmt.Sprintf("No workload definition found for testStrategy %s", testStrategy))
		}
	}

	// now lets send the test finished event
	if !res {
		if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, keptnv2.ResultFailed, logger); err != nil {
			logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
		}
		return
	}

	if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, keptnv2.ResultPass, logger); err != nil {
		logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
	}
}

//
// Extracts relevant information from the data object
//
func getTestInfo(data keptnv2.TestTriggeredEventData, shkeptncontext string) *TestInfo {
	return &TestInfo{
		Project:      data.Project,
		Service:      data.Service,
		Stage:        data.Stage,
		TestStrategy: data.Test.TestStrategy,
		Context:      shkeptncontext,
	}
}

//
// returns the service URL that is either passed via the DeploymentURI* parameters or constructs one based on keptn naming structure
//
func getServiceURL(data keptnv2.TestTriggeredEventData) (*url.URL, error) {

	if len(data.Deployment.DeploymentURIsLocal) > 0 && data.Deployment.DeploymentURIsLocal[0] != "" {
		return url.Parse(data.Deployment.DeploymentURIsLocal[0])

	} else if len(data.Deployment.DeploymentURIsPublic) > 0 && data.Deployment.DeploymentURIsPublic[0] != "" {
		return url.Parse(data.Deployment.DeploymentURIsPublic[0])
	}

	return nil, errors.New("no deployment URI included in event")
}

//
// executes the actual JMEter tests based on the workload configuration
//
func runWorkload(serviceURL *url.URL, testInfo *TestInfo, workload *Workload, logger *keptncommon.Logger) (bool, error) {

	// for testStrategy functional we enforce a 0% error policy!
	breakOnFunctionalIssues := workload.TestStrategy == TestStrategy_Functional

	logger.Info(
		fmt.Sprintf("Running workload testStrategy=%s, vuser=%d, loopcount=%d, thinktime=%d, funcvalidation=%t, acceptederrors=%f, avgrtvalidation=%d, script=%s",
			workload.TestStrategy, workload.VUser, workload.LoopCount, workload.ThinkTime, breakOnFunctionalIssues, workload.AcceptedErrorRate, workload.AvgRtValidation, workload.Script))
	if runlocal {
		logger.Info("LOCALLY: not executing actual tests!")
		return true, nil
	}

	// the resultdirectory is unique as it contains context but also gives some human readable context such as teststrategy and service
	// this will also be used for TSN parameter
	resultDirectory := fmt.Sprintf("%s_%s_%s_%s_%s", testInfo.Project, testInfo.Service, testInfo.Stage, workload.TestStrategy, testInfo.Context)

	// lets first remove all potentially left over result files from previous runs -> we keept them between runs for troubleshooting though
	os.RemoveAll(resultDirectory)
	os.RemoveAll(resultDirectory + "_result.tlf")
	os.RemoveAll("output.txt")

	return executeJMeter(testInfo, workload, resultDirectory, serviceURL, resultDirectory, breakOnFunctionalIssues, logger)
}

func sendTestsStartedEvent(shkeptncontext string, incomingEvent cloudevents.Event, logger *keptncommon.Logger) error {
	source, _ := url.Parse("jmeter-service")

	testStartedEventData := keptnv2.TestStartedEventData{}
	testTriggeredEventData := keptnv2.TestTriggeredEventData{}
	// fill in data from incoming event (e.g., project, service, stage)
	if err := incomingEvent.DataAs(&testTriggeredEventData); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	testStartedEventData.EventData = testTriggeredEventData.EventData
	testStartedEventData.Status = keptnv2.StatusSucceeded

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetStartedEventType(keptnv2.TestTaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetExtension("triggeredid", incomingEvent.ID())
	event.SetData(cloudevents.ApplicationJSON, testStartedEventData)

	return sendEvent(event)
}

func sendTestsFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event, startedAt time.Time, result keptnv2.ResultType, logger *keptncommon.Logger) error {
	source, _ := url.Parse("jmeter-service")

	testFinishedData := keptnv2.TestFinishedEventData{}
	testTriggeredEventData := keptnv2.TestTriggeredEventData{}
	// fill in data from incoming event (e.g., project, service, stage)
	if err := incomingEvent.DataAs(&testTriggeredEventData); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	testFinishedData.EventData = testTriggeredEventData.EventData
	// fill in timestamps
	testFinishedData.Test.Start = startedAt.Format(time.RFC3339)
	testFinishedData.Test.End = time.Now().Format(time.RFC3339)
	// set test result
	testFinishedData.Result = result
	testFinishedData.Status = keptnv2.StatusSucceeded

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetFinishedEventType(keptnv2.TestTaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetExtension("triggeredid", incomingEvent.ID())
	event.SetData(cloudevents.ApplicationJSON, testFinishedData)

	return sendEvent(event)
}

func sendErroredTestsFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event, startedAt time.Time, msg string, logger *keptncommon.Logger) error {
	source, _ := url.Parse("jmeter-service")

	testFinishedData := keptnv2.TestFinishedEventData{}
	testTriggeredEventData := keptnv2.TestTriggeredEventData{}
	// fill in data from incoming event (e.g., project, service, stage)
	if err := incomingEvent.DataAs(&testTriggeredEventData); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	testFinishedData.EventData = testTriggeredEventData.EventData
	// fill in timestamps
	testFinishedData.Test.Start = startedAt.Format(time.RFC3339)
	testFinishedData.Test.End = time.Now().Format(time.RFC3339)
	// set test result
	testFinishedData.Result = keptnv2.ResultFailed
	testFinishedData.Status = keptnv2.StatusErrored
	testFinishedData.Message = msg

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetFinishedEventType(keptnv2.TestTaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetExtension("triggeredid", incomingEvent.ID())
	event.SetData(cloudevents.ApplicationJSON, testFinishedData)

	return sendEvent(event)
}

func _main(args []string, env envConfig) int {

	if runlocal {
		log.Println("Running LOCALLY: env=runlocal")
	}

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(ctx, gotEvent))

	return 0
}

func sendEvent(event cloudevents.Event) error {
	if runlocal {
		log.Println("LOCALLY: Sending Event")
		return nil
	}

	endPoint, err := getServiceEndpoint(eventbroker)
	if err != nil {
		return errors.New("Failed to retrieve endpoint of eventbroker. %s" + err.Error())
	}

	if endPoint.Host == "" {
		return errors.New("Host of eventbroker not set")
	}
	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		EventBrokerURL: endPoint.String(),
	})

	if err != nil {
		return errors.New("Failed to initialize Keptn handler: " + err.Error())
	}

	return keptnHandler.SendCloudEvent(event)
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
