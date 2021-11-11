package main

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/common/retry"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
	"net"
	"net/url"
	"time"
)

type EventHandler struct {
	testRunner *TestRunner
}

func (e *EventHandler) handleEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	data := &keptnv2.TestTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		logger.WithError(err).Error("Unable to decode 'test.triggered' event data")
		return nil
	}

	if event.Type() != keptnv2.GetTriggeredEventType(keptnv2.TestTaskName) {
		logger.Errorf("Received unexpected keptn event: %s", event.Type())
		return nil
	}
	if data.Test.TestStrategy == TestStrategy_RealUser {
		logger.Infof("Received '%s' test strategy, hence no tests are triggered", TestStrategy_RealUser)
		return nil
	}

	testInfo, err := createTestInfo(*data, shkeptncontext, event.ID())
	if err != nil {
		logger.Errorf("Unable to create test info: %v", err)
		return nil
	}

	go e.testRunner.RunTests(*testInfo)
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
