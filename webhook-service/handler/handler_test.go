package handler_test

import (
	"encoding/json"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	fakekeptn "github.com/keptn/keptn/go-sdk/pkg/sdk/fake"
	"github.com/keptn/keptn/webhook-service/handler"
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/keptn/keptn/webhook-service/lib/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
)

const webHookContent = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithMultipleRequests = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
          "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithMissingTemplateData = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.unavailable}} {{.env.mysecret}}"`

func newWebhookTriggeredEvent(filename string) cloudevents.Event {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := models.KeptnContextExtendedCE{}
	err = json.Unmarshal(content, &event)
	_ = err
	return keptnv2.ToCloudEvent(event)
}

func Test_HandleIncomingTriggeredEvent(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{ParseTemplateFunc: func(data interface{}, templateStr string) (string, error) {
		tplE := &lib.TemplateEngine{}
		return tplE.ParseTemplate(data, templateStr)
	}}

	secretReaderMock := &fake.ISecretReaderMock{}
	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "my-secret-value", nil
	}

	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlExecutorMock.CurlFunc = func(curlCmd string) (string, error) {
		return "success", nil
	}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler,
		),
		fakekeptn.WithResourceHandler(fakekeptn.StringResourceHandler{ResourceContent: webHookContent}),
	)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())

	finishedEvent, err := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	eventData := &keptnv2.EventData{}
	keptnv2.EventDataAs(finishedEvent, eventData)
	require.Nil(t, err)
	assert.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	assert.Equal(t, keptnv2.ResultPass, eventData.Result)
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
		),
		fakekeptn.WithResourceHandler(fakekeptn.FailingResourceHandler{}),
	)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())

	finishedEvent, err := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	eventData := &keptnv2.EventData{}
	keptnv2.EventDataAs(finishedEvent, eventData)
	require.Nil(t, err)
	assert.Equal(t, keptnv2.StatusErrored, eventData.Status)
	assert.Equal(t, keptnv2.ResultFailed, eventData.Result)
}

func TestTaskHandler_CannotReadSecret(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "", errors.New("unable to read secret :(")
	}
	curlExecutorMock := &fake.ICurlExecutorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler),
		fakekeptn.WithResourceHandler(fakekeptn.StringResourceHandler{ResourceContent: webHookContent}))

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())

	finishedEvent, err := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	eventData := &keptnv2.EventData{}
	keptnv2.EventDataAs(finishedEvent, eventData)
	require.Nil(t, err)
	assert.Equal(t, keptnv2.StatusErrored, eventData.Status)
	assert.Equal(t, keptnv2.ResultFailed, eventData.Result)
}

func TestTaskHandler_IncompleteDataForTemplate(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{ParseTemplateFunc: func(data interface{}, templateStr string) (string, error) {
		tplE := &lib.TemplateEngine{}
		return tplE.ParseTemplate(data, templateStr)
	}}
	secretReaderMock := &fake.ISecretReaderMock{}
	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "my-secret-value", nil
	}
	curlExecutorMock := &fake.ICurlExecutorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler),
		fakekeptn.WithResourceHandler(fakekeptn.StringResourceHandler{ResourceContent: webHookContentWithMissingTemplateData}))

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.NotEmpty(t, secretReaderMock.ReadSecretCalls())
	require.NotEmpty(t, templateEngineMock.ParseTemplateCalls())
	require.Empty(t, curlExecutorMock.CurlCalls())

	//verify sent events
	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())

	finishedEvent, err := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	eventData := &keptnv2.EventData{}
	keptnv2.EventDataAs(finishedEvent, eventData)
	require.Nil(t, err)
	assert.Equal(t, keptnv2.StatusErrored, eventData.Status)
	assert.Equal(t, keptnv2.ResultFailed, eventData.Result)
}

func TestTaskHandler_CurlExecutorFails(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{ParseTemplateFunc: func(data interface{}, templateStr string) (string, error) {
		tplE := &lib.TemplateEngine{}
		return tplE.ParseTemplate(data, templateStr)
	}}
	secretReaderMock := &fake.ISecretReaderMock{}
	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "my-secret-value", nil
	}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlExecutorMock.CurlFunc = func(curlCmd string) (string, error) {
		return "", errors.New("unable to execute curl call")
	}
	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, secretReaderMock)

	fakeKeptn := fakekeptn.NewFakeKeptn(
		"test-webhook-svc",
		sdk.WithHandler(
			"*",
			taskHandler),
		fakekeptn.WithResourceHandler(fakekeptn.StringResourceHandler{ResourceContent: webHookContent}))

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.NotEmpty(t, secretReaderMock.ReadSecretCalls())
	require.NotEmpty(t, templateEngineMock.ParseTemplateCalls())
	require.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	require.Equal(t, 2, len(fakeKeptn.GetEventSender().SentEvents))
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[1].Type())

	finishedEvent, err := keptnv2.ToKeptnEvent(fakeKeptn.GetEventSender().SentEvents[1])
	eventData := &keptnv2.EventData{}
	keptnv2.EventDataAs(finishedEvent, eventData)
	require.Nil(t, err)
	assert.Equal(t, keptnv2.StatusErrored, eventData.Status)
	assert.Equal(t, keptnv2.ResultFailed, eventData.Result)
}
