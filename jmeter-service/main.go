package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/sirupsen/logrus"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/retry"
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
	os.Exit(_main(env))
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
	testInfo, err := createTestInfo(*data, shkeptncontext, event.ID())
	if err != nil {
		return err
	}
	jmeterConfig, err := getJMeterConf(*testInfo)
	if err != nil {
		return err
	}
	eventSender, err := keptnv2.NewHTTPEventSender("")
	if err != nil {
		return err
	}
	testRunner := NewTestRunner(*jmeterConfig, eventSender)
	go testRunner.RunTests(*testInfo)
	return nil
}

func createTestInfo(data keptnv2.TestTriggeredEventData, shkeptncontext string, triggeredID string) (*TestInfo, error) {
	serviceURL, err := getServiceURL(data)
	if err != nil {
		return nil, err
	}
	return &TestInfo{
		Project:           data.Project,
		Service:           data.Service,
		Stage:             data.Stage,
		TestStrategy:      data.Test.TestStrategy,
		Context:           shkeptncontext,
		TriggeredID:       triggeredID,
		TestTriggeredData: data,
		ServiceURL:        serviceURL,
	}, nil
}

// getServiceURL returns the service URL that is either passed via the DeploymentURI* parameters or constructs one based on keptn naming structure
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

func _main(env envConfig) int {
	if runlocal {
		log.Println("Running LOCALLY: env=runlocal")
	}

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

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
