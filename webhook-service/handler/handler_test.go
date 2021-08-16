package handler_test

import (
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	fakekeptn "github.com/keptn/keptn/go-sdk/pkg/sdk/fake"
	"github.com/keptn/keptn/webhook-service/handler"
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/keptn/keptn/webhook-service/lib/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

const webHookContent = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.deployment.triggered"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.project}} {{.env.mysecret}}"`

const webHookContentWithMissingTemplateData = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.deployment.triggered"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.unavailable}} {{.env.mysecret}}"`

func TestTaskHandler_Execute(t *testing.T) {

	templateEngineMock := &fake.ITemplateEngineMock{ParseTemplateFunc: func(data interface{}, templateStr string) (string, error) {
		tplE := &lib.TemplateEngine{}
		return tplE.ParseTemplate(data, templateStr)
	}}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler,
			map[string]interface{}{},
		),
	)

	mockResourceHandler := &fakekeptn.ResourceHandlerMock{
		GetProjectResourceFunc: func(project string, resourceURI string) (*models.Resource, error) {
			return &models.Resource{ResourceURI: &resourceURI, ResourceContent: webHookContent}, nil
		},
		GetServiceResourceFunc: func(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
			return &models.Resource{ResourceURI: &resourceURI, ResourceContent: webHookContent}, nil
		},
		GetStageResourceFunc: func(project string, stage string, resourceURI string) (*models.Resource, error) {
			return &models.Resource{ResourceURI: &resourceURI, ResourceContent: webHookContent}, nil
		},
	}

	fakeKeptn.TestResourceHandler = mockResourceHandler

	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "my-secret-value", nil
	}

	curlExecutorMock.CurlFunc = func(curlCmd string) (string, error) {
		return "success", nil
	}

	result, sdkErr := taskHandler.Execute(fakeKeptn, &keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))

	require.Nil(t, sdkErr)
	require.NotNil(t, result)

	require.Equal(t, "curl http://localhost:8080 my-project my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
}

func TestTaskHandler_Execute_WebhookCannotBeRetrieved(t *testing.T) {

	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler,
			map[string]interface{}{},
		),
	)

	mockResourceHandler := &fakekeptn.ResourceHandlerMock{
		GetProjectResourceFunc: func(project string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
		GetServiceResourceFunc: func(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
		GetStageResourceFunc: func(project string, stage string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
	}

	fakeKeptn.TestResourceHandler = mockResourceHandler

	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "my-secret-value", nil
	}

	curlExecutorMock.CurlFunc = func(curlCmd string) (string, error) {
		return "success", nil
	}

	result, sdkErr := taskHandler.Execute(fakeKeptn, &keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))

	require.NotNil(t, sdkErr)
	require.Nil(t, result)

	require.Empty(t, curlExecutorMock.CurlCalls())
	require.Empty(t, secretReaderMock.ReadSecretCalls())
	require.Empty(t, templateEngineMock.ParseTemplateCalls())
}

func TestTaskHandler_CannotReadSecret(t *testing.T) {

	templateEngineMock := &fake.ITemplateEngineMock{ParseTemplateFunc: func(data interface{}, templateStr string) (string, error) {
		tplE := &lib.TemplateEngine{}
		return tplE.ParseTemplate(data, templateStr)
	}}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler,
			map[string]interface{}{},
		),
	)

	mockResourceHandler := &fakekeptn.ResourceHandlerMock{
		GetProjectResourceFunc: func(project string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
		GetServiceResourceFunc: func(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
		GetStageResourceFunc: func(project string, stage string, resourceURI string) (*models.Resource, error) {
			return &models.Resource{ResourceURI: &resourceURI, ResourceContent: webHookContent}, nil
		},
	}

	fakeKeptn.TestResourceHandler = mockResourceHandler

	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "", errors.New("oops")
	}

	curlExecutorMock.CurlFunc = func(curlCmd string) (string, error) {
		return "success", nil
	}

	result, sdkErr := taskHandler.Execute(fakeKeptn, &keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))

	require.NotNil(t, sdkErr)
	require.Nil(t, result)

	require.NotEmpty(t, secretReaderMock.ReadSecretCalls())
	require.Empty(t, curlExecutorMock.CurlCalls())
	require.Empty(t, templateEngineMock.ParseTemplateCalls())
}

func TestTaskHandler_IncompleteDataForTemplate(t *testing.T) {

	templateEngineMock := &fake.ITemplateEngineMock{ParseTemplateFunc: func(data interface{}, templateStr string) (string, error) {
		tplE := &lib.TemplateEngine{}
		return tplE.ParseTemplate(data, templateStr)
	}}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler,
			map[string]interface{}{},
		),
	)

	mockResourceHandler := &fakekeptn.ResourceHandlerMock{
		GetProjectResourceFunc: func(project string, resourceURI string) (*models.Resource, error) {
			return &models.Resource{ResourceURI: &resourceURI, ResourceContent: webHookContentWithMissingTemplateData}, nil
		},
		GetServiceResourceFunc: func(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
		GetStageResourceFunc: func(project string, stage string, resourceURI string) (*models.Resource, error) {
			return nil, nil
		},
	}

	fakeKeptn.TestResourceHandler = mockResourceHandler

	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "my-secret-value", nil
	}

	curlExecutorMock.CurlFunc = func(curlCmd string) (string, error) {
		return "success", nil
	}

	result, sdkErr := taskHandler.Execute(fakeKeptn, &keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))

	require.NotNil(t, sdkErr)
	require.Nil(t, result)

	require.NotEmpty(t, secretReaderMock.ReadSecretCalls())
	require.NotEmpty(t, templateEngineMock.ParseTemplateCalls())
	require.Empty(t, curlExecutorMock.CurlCalls())
}
