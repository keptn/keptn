package main

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
	"net/url"
)

// EventHandler handles events of type 'test.triggered' and kicks off
// the TestRunner to execute the JMeter tests
type EventHandler struct {
	testRunner *TestRunner
}

func (e *EventHandler) handleEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	if err := event.Context.ExtensionAs("shkeptncontext", &shkeptncontext); err != nil {
		return err
	}

	if event.Type() != keptnv2.GetTriggeredEventType(keptnv2.TestTaskName) {
		logger.Warnf("Received unexpected keptn event: %s", event.Type())
		return nil
	}
	data := &keptnv2.TestTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		logger.Errorf("Unable to decode 'test.triggered' event data: %v", err)
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

	//go func() {
	if err := e.testRunner.RunTests(ctx, *testInfo); err != nil {
		logger.Errorf("Unable to run JMeter tests: %v", err)
	}
	//}()
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
