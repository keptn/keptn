package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/sirupsen/logrus"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/retry"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const (
	configurationService = "CONFIGURATION_SERVICE"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port     int    `envconfig:"RCV_PORT" default:"8080"`
	Path     string `envconfig:"RCV_PATH" default:"/"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logger.Fatalf("Failed to process env var: %s", err)
	}

	logger.SetLevel(logger.InfoLevel)

	if os.Getenv(env.LogLevel) != "" {
		logLevel, err := logger.ParseLevel(os.Getenv(env.LogLevel))
		if err != nil {
			logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			logger.SetLevel(logLevel)
		}
	}
	os.Exit(_main(os.Args[1:], env))
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	data := &keptnv2.TestTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		logger.WithError(err).Error("Got Data Error")
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
	ctx.Value("Wg").(*sync.WaitGroup).Add(1)
	go runTests(ctx, event, shkeptncontext, *data)

	return nil
}

//
// This method executes the correct tests based on the passed testStrategy in the deployment finished event
// The method will always try to execute a health check workload first, then execute the workload based on the passed testStrategy
//
func runTests(ctx context.Context, event cloudevents.Event, shkeptncontext string, data keptnv2.TestTriggeredEventData) {
	defer ctx.Value("Wg").(*sync.WaitGroup).Done()
	sendTestsStartedEvent(shkeptncontext, event)

	testInfo := getTestInfo(data, shkeptncontext)
	startedAt := time.Now()

	go func() {
		<-ctx.Done()
		logger.Error("Error sending test finished event" + testInfo.ToString())
		if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, "received a SIGTERM/SIGINT, jmeter terminated before the end of the test"); err != nil {
			logger.Error(fmt.Sprintf("Error sending test finished event: %s", err.Error()) + ". " + testInfo.ToString())
		}
		ctx.Value("Wg").(*sync.WaitGroup).Done()
		return
	}()

	// load the workloads from JMeterConf
	var err error
	var jmeterconf *JMeterConf
	jmeterconf, err = getJMeterConf(testInfo.Project, testInfo.Stage, testInfo.Service)

	// get the service endpoint we need to run the test against
	var serviceUrl *url.URL
	serviceUrl, err = getServiceURL(data)
	if err != nil {
		msg := fmt.Sprintf("Error getting service url to run test against: %s", err.Error())
		logger.Error(msg)
		if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, msg); err != nil {
			logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
		}
		return
	}

	// first we run a health check workload. If that fails we stop the rest
	var healthCheckWorkload *Workload
	var res bool
	healthCheckWorkload, err = getWorkload(jmeterconf, TestStrategy_HealthCheck)
	if healthCheckWorkload != nil {
		// do a basic health check, verifying whether the endpoint is available or not
		err := checkEndpointAvailable(5*time.Second, serviceUrl)
		if err != nil {
			msg := fmt.Sprintf("Jmeter-service cannot reach URL %s: %s", serviceUrl, err.Error())
			logger.Error(msg)
			if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, msg, keptnv2.ResultFailed); err != nil {
				logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
			}
			return
		}

		res, err = runWorkload(serviceUrl, testInfo, healthCheckWorkload)
		if err != nil {
			msg := fmt.Sprintf("could not run test workload: %s", err.Error())
			logger.Error(msg)
			if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, msg); err != nil {
				logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
			}
			return
		}

		if !res {
			msg := fmt.Sprintf("Tests for %s with status = %s.%s", TestStrategy_HealthCheck, strconv.FormatBool(res), testInfo.ToString())
			if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, msg, keptnv2.ResultFailed); err != nil {
				logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
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
			res, err = runWorkload(serviceUrl, testInfo, teststrategyWorkload)
			if err != nil {
				msg := fmt.Sprintf("could not run test workload: %s", err.Error())
				logger.Error(msg)
				if err := sendErroredTestsFinishedEvent(shkeptncontext, event, startedAt, msg); err != nil {
					logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
				}
				return
			} else {
				logger.Infof("Tests for %s with status = %s.%s", testStrategy, strconv.FormatBool(res), testInfo.ToString())
			}
		} else {
			// no workload for that test strategy!!
			res = false
			logger.Errorf("No workload definition found for testStrategy %s", testStrategy)
		}
	}

	// now lets send the test finished event
	msg := fmt.Sprintf("Tests for %s with status = %s.%s", testStrategy, strconv.FormatBool(res), testInfo.ToString())
	if !res {
		if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, msg, keptnv2.ResultFailed); err != nil {
			logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
		}
		return
	}

	if err := sendTestsFinishedEvent(shkeptncontext, event, startedAt, msg, keptnv2.ResultPass); err != nil {
		logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
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
		newurl, err := url.Parse(data.Deployment.DeploymentURIsLocal[0])
		if newurl.Path == "" {
			newurl.Path += "/"
		}
		return newurl, err

	} else if len(data.Deployment.DeploymentURIsPublic) > 0 && data.Deployment.DeploymentURIsPublic[0] != "" {
		newurl, err := url.Parse(data.Deployment.DeploymentURIsPublic[0])
		if newurl.Path == "" {
			newurl.Path += "/"
		}
		return newurl, err
	}

	return nil, errors.New("no deployment URI included in event")
}

//
// executes the actual JMeter tests based on the workload configuration
//
func runWorkload(serviceURL *url.URL, testInfo *TestInfo, workload *Workload) (bool, error) {
	// for testStrategy functional we enforce a 0% error policy!
	breakOnFunctionalIssues := workload.TestStrategy == TestStrategy_Functional

	logger.Infof(
		"Running workload testStrategy=%s, vuser=%d, loopcount=%d, thinktime=%d, funcvalidation=%t, acceptederrors=%f, avgrtvalidation=%d, script=%s",
		workload.TestStrategy, workload.VUser, workload.LoopCount, workload.ThinkTime, breakOnFunctionalIssues, workload.AcceptedErrorRate, workload.AvgRtValidation, workload.Script)
	if runlocal {
		logger.Info("LOCALLY: not executing actual tests!")
		return true, nil
	}

	// the resultdirectory is unique as it contains context but also gives some human readable context such as teststrategy and service
	// this will also be used for TSN parameter
	resultDirectory := fmt.Sprintf("%s_%s_%s_%s_%s", testInfo.Project, testInfo.Service, testInfo.Stage, workload.TestStrategy, testInfo.Context)

	// lets first remove all potentially left over result files from previous runs -> we keep them between runs for troubleshooting though
	err := os.RemoveAll(resultDirectory)
	if err != nil {
		return false, err
	}

	err = os.RemoveAll(resultDirectory + "_result.tlf")
	if err != nil {
		return false, err
	}

	err = os.RemoveAll("output.txt")
	if err != nil {
		return false, err
	}

	return executeJMeter(testInfo, workload, resultDirectory, serviceURL, resultDirectory, breakOnFunctionalIssues)
}

func checkEndpointAvailable(timeout time.Duration, serviceURL *url.URL) error {
	if serviceURL == nil {
		return fmt.Errorf("url to check for reachability is nil")
	}

	// serviceURL.Host does not contain the port in case of serviceURL=http://1.2.3.4/ (without port)
	// hence we need to manually construct hostWithPort here
	hostWithPort := fmt.Sprintf("%s:%s", serviceURL.Hostname(), derivePort(serviceURL))

	var err error = nil

	_ = retry.Retry(func() error {
		if _, err = net.DialTimeout("tcp", hostWithPort, timeout); err != nil {
			return err
		}

		return nil
	}, retry.DelayBetweenRetries(time.Second*5), retry.NumberOfRetries(3))

	return err
}

func sendTestsStartedEvent(shkeptncontext string, incomingEvent cloudevents.Event) error {
	source, _ := url.Parse("jmeter-service")

	testStartedEventData := keptnv2.TestStartedEventData{}
	testTriggeredEventData := keptnv2.TestTriggeredEventData{}
	// fill in data from incoming event (e.g., project, service, stage)
	if err := incomingEvent.DataAs(&testTriggeredEventData); err != nil {
		logger.WithError(err).Error("Got Data Error")
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

func sendTestsFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event, startedAt time.Time, msg string, result keptnv2.ResultType) error {
	source, _ := url.Parse("jmeter-service")

	testFinishedData := keptnv2.TestFinishedEventData{}
	testTriggeredEventData := keptnv2.TestTriggeredEventData{}
	// fill in data from incoming event (e.g., project, service, stage)
	if err := incomingEvent.DataAs(&testTriggeredEventData); err != nil {
		logger.WithError(err).Error("Got Data Error")
		return err
	}
	testFinishedData.EventData = testTriggeredEventData.EventData
	// fill in timestamps
	testFinishedData.Test.Start = startedAt.Format(time.RFC3339)
	testFinishedData.Test.End = time.Now().Format(time.RFC3339)
	// set test result
	testFinishedData.Result = result
	testFinishedData.Status = keptnv2.StatusSucceeded
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

func sendErroredTestsFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event, startedAt time.Time, msg string) error {
	source, _ := url.Parse("jmeter-service")

	testFinishedData := keptnv2.TestFinishedEventData{}
	testTriggeredEventData := keptnv2.TestTriggeredEventData{}
	// fill in data from incoming event (e.g., project, service, stage)
	if err := incomingEvent.DataAs(&testTriggeredEventData); err != nil {
		logger.WithError(err).Error("Got Data Error")
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

	ctx := getGracefulContext()

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(keptnapi.HealthEndpointHandler))
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
	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{})

	if err != nil {
		return errors.New("Failed to initialize Keptn handler: " + err.Error())
	}

	return keptnHandler.SendCloudEvent(event)
}

func getGracefulContext() context.Context {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), "Wg", wg)))

	go func() {
		<-ch
		log.Fatal("Container termination triggered, starting graceful shutdown")
		cancel()
	}()

	return ctx
}
