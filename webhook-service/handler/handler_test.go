package handler_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/handler"
	fake2 "github.com/keptn/keptn/webhook-service/handler/fake"
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/keptn/keptn/webhook-service/lib/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const webHookContent = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionID: "my-subscription-id"
      sendFinished: true
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithStartedEvent = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.started"
      subscriptionID: "my-subscription-id"
      sendFinished: true
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithFinishedEvent = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.finished"
      subscriptionID: "my-subscription-id"
      sendFinished: true
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
      subscriptionID: "my-subscription-id"
      sendFinished: true
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithMultipleRequestsAndDisabledFinished = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionID: "my-subscription-id"
      sendFinished: false
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithMultipleRequestsAndDisabledStarted = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionID: "my-subscription-id"
      sendStarted: false
      sendFinished: false
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentWithMissingTemplateData = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionID: "my-subscription-id"
      sendFinished: true
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.unavailable}} {{.env.mysecret}}"`

const webHookContentWithNoMatchingSubscriptionID = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionID: "my-unmatched-subscription-id"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`

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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 1)
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

func Test_HandleIncomingStartedEvent(t *testing.T) {
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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithStartedEvent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.started.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 1)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	require.Empty(t, fakeKeptn.GetEventSender().SentEvents)
}

func Test_HandleIncomingStartedEventWithResultingError(t *testing.T) {
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
		return "", errors.New("oops")
	}

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithStartedEvent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.started.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 1)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	require.Empty(t, fakeKeptn.GetEventSender().SentEvents)
}

func Test_HandleIncomingFinishedEvent(t *testing.T) {
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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithFinishedEvent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.finished.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 1)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	require.Empty(t, fakeKeptn.GetEventSender().SentEvents)
}

func Test_HandleIncomingFinishedEventWithResultingError(t *testing.T) {
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
		return "", errors.New("oops")
	}

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithFinishedEvent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.finished.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 1)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	require.Empty(t, fakeKeptn.GetEventSender().SentEvents)
}

func Test_HandleIncomingTriggeredEvent_SendMultipleRequests(t *testing.T) {
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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequests})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 2)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)

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

func Test_HandleIncomingTriggeredEvent_SendMultipleRequestsDisableFinished(t *testing.T) {
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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequestsAndDisabledFinished})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 4)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[2].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[3].CurlCmd)

	//verify sent events
	require.Equal(t, 4, len(fakeKeptn.GetEventSender().SentEvents))
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[1].Type())
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[2].Type())
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[3].Type())
}

func Test_HandleIncomingTriggeredEvent_SendMultipleRequestsDisableStarted(t *testing.T) {
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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequestsAndDisabledStarted})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 4)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[2].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[3].CurlCmd)

	//verify sent events
	require.Empty(t, len(fakeKeptn.GetEventSender().SentEvents))
}

func Test_HandleIncomingTriggeredEvent_SendMultipleRequestsDisableFinishedOneRequestFails(t *testing.T) {
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
		// make the second request fail
		if len(curlExecutorMock.CurlCalls()) == 2 {
			return "", errors.New("oops")
		}
		return "success", nil
	}

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequestsAndDisabledFinished})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Len(t, curlExecutorMock.CurlCalls(), 2)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://localhost:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)

	//verify sent events
	require.Equal(t, 7, len(fakeKeptn.GetEventSender().SentEvents))
	// each request should have a .started event
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[1].Type())
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[2].Type())
	assert.Equal(t, "sh.keptn.event.webhook.started", fakeKeptn.GetEventSender().SentEvents[3].Type())
	// apart from the first request which has been successful, every request should have a failed .finished event
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[4].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[5].Type())
	assert.Equal(t, "sh.keptn.event.webhook.finished", fakeKeptn.GetEventSender().SentEvents[6].Type())
}

func Test_HandleIncomingTriggeredEvent_NoMatchingWebhookFound(t *testing.T) {
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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithNoMatchingSubscriptionID})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

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

func TestTaskHandler_Execute_WebhookCannotBeRetrieved(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.FailingResourceHandler{})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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

func TestTaskHandler_Execute_NoSubscriptionIDInEvent(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.FailingResourceHandler{})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-no-subscription-id.json"))

	//verify sent events
	require.Equal(t, 0, len(fakeKeptn.GetEventSender().SentEvents))

	require.Empty(t, curlExecutorMock.CurlCalls())
}

func TestTaskHandler_Execute_InvalidEvent(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.FailingResourceHandler{})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/invalid-event.json"))

	//verify sent events
	require.Empty(t, fakeKeptn.GetEventSender().SentEvents)

	require.Empty(t, curlExecutorMock.CurlCalls())
}

func TestTaskHandler_CannotReadSecret(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "", errors.New("unable to read secret :(")
	}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMissingTemplateData})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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
	curlValidatorMock := &fake.ICurlValidatorMock{}
	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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

func TestTaskHandler_CurlExecutorFailsHideSecret(t *testing.T) {
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
		return "", errors.New("unable to execute curl call containing secret my-secret-value")
	}
	curlValidatorMock := &fake.ICurlValidatorMock{}
	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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
	assert.NotContains(t, eventData.Message, "my-secret-value")
}

func TestTaskHandler_Execute_WebhookConfigInService(t *testing.T) {
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

	resourceHandlerMock := &fake2.IResourceHandlerMock{}
	resourceHandlerMock.GetResourceFunc = func(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
		return &models.Resource{
			Metadata:        &models.Version{Version: "CommitID"},
			ResourceContent: webHookContent,
			ResourceURI:     nil,
		}, nil
	}

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")

	fakeKeptn.SetResourceHandler(resourceHandlerMock)
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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
	assert.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	assert.Equal(t, keptnv2.ResultPass, eventData.Result)

	require.Len(t, resourceHandlerMock.GetResourceCalls(), 1)

	// need to use reflection package to inspect passed parameters since the properties are unexported
	scopeVals := reflect.ValueOf(resourceHandlerMock.GetResourceCalls()[0].Scope)
	require.Equal(t, "myservice", scopeVals.FieldByName("service").String())
	require.Equal(t, "mystage", scopeVals.FieldByName("stage").String())
	require.Equal(t, "myproject", scopeVals.FieldByName("project").String())
}

func TestTaskHandler_Execute_WebhookConfigInStage(t *testing.T) {
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

	resourceHandlerMock := &fake2.IResourceHandlerMock{}
	nrResourceRequests := 0
	resourceHandlerMock.GetResourceFunc = func(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
		// the first request, i.e. when checking  at service level, will return no result
		if nrResourceRequests == 0 {
			nrResourceRequests++
			return nil, nil
		}
		// the second request (stage-level) will return something
		return &models.Resource{
			Metadata:        &models.Version{Version: "CommitID"},
			ResourceContent: webHookContent,
			ResourceURI:     nil,
		}, nil
	}

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")

	fakeKeptn.SetResourceHandler(resourceHandlerMock)
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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
	assert.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	assert.Equal(t, keptnv2.ResultPass, eventData.Result)

	require.Len(t, resourceHandlerMock.GetResourceCalls(), 2)

	// need to use reflection package to inspect passed parameters since the properties are unexported
	scopeVals1 := reflect.ValueOf(resourceHandlerMock.GetResourceCalls()[0].Scope)
	require.Equal(t, "myservice", scopeVals1.FieldByName("service").String())
	require.Equal(t, "mystage", scopeVals1.FieldByName("stage").String())
	require.Equal(t, "myproject", scopeVals1.FieldByName("project").String())

	// for the second request, we should check for the resource at stage level
	scopeVals2 := reflect.ValueOf(resourceHandlerMock.GetResourceCalls()[1].Scope)
	require.Equal(t, "", scopeVals2.FieldByName("service").String())
	require.Equal(t, "mystage", scopeVals2.FieldByName("stage").String())
	require.Equal(t, "myproject", scopeVals2.FieldByName("project").String())
}

func TestTaskHandler_Execute_WebhookConfigInProject(t *testing.T) {
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

	resourceHandlerMock := &fake2.IResourceHandlerMock{}
	nrResourceRequests := 0
	resourceHandlerMock.GetResourceFunc = func(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
		// the first request, i.e. when checking  at service level, will return no result
		if nrResourceRequests <= 1 {
			nrResourceRequests++
			return nil, nil
		}
		// the third request (stage-level) will return something
		return &models.Resource{
			Metadata:        &models.Version{Version: "CommitID"},
			ResourceContent: webHookContent,
			ResourceURI:     nil,
		}, nil
	}

	curlValidatorMock := &fake.ICurlValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")

	fakeKeptn.SetResourceHandler(resourceHandlerMock)
	fakeKeptn.AddTaskHandler("*", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

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
	assert.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	assert.Equal(t, keptnv2.ResultPass, eventData.Result)

	require.Len(t, resourceHandlerMock.GetResourceCalls(), 3)

	// need to use reflection package to inspect passed parameters since the properties are unexported
	scopeVals1 := reflect.ValueOf(resourceHandlerMock.GetResourceCalls()[0].Scope)
	require.Equal(t, "myservice", scopeVals1.FieldByName("service").String())
	require.Equal(t, "mystage", scopeVals1.FieldByName("stage").String())
	require.Equal(t, "myproject", scopeVals1.FieldByName("project").String())

	// for the second request, we should check for the resource at stage level
	scopeVals2 := reflect.ValueOf(resourceHandlerMock.GetResourceCalls()[1].Scope)
	require.Equal(t, "", scopeVals2.FieldByName("service").String())
	require.Equal(t, "mystage", scopeVals2.FieldByName("stage").String())
	require.Equal(t, "myproject", scopeVals2.FieldByName("project").String())

	// for the third request, we should check for the resource at project level
	scopeVals3 := reflect.ValueOf(resourceHandlerMock.GetResourceCalls()[2].Scope)
	require.Equal(t, "", scopeVals3.FieldByName("service").String())
	require.Equal(t, "", scopeVals3.FieldByName("stage").String())
	require.Equal(t, "myproject", scopeVals3.FieldByName("project").String())
}

func Test_createRequest(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	curlValidatorMock := &fake.ICurlValidatorMock{}
	curlValidatorMock.ValidateFunc = func(request lib.Request) error {
		return nil
	}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, curlValidatorMock, secretReaderMock)

	tests := []struct {
		name    string
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "valid alpha input",
			data:    "curl command",
			want:    "curl command",
			wantErr: false,
		},
		{
			name: "valid beta input #1",
			data: lib.Request{
				Headers: []lib.Header{
					{
						Key:   "key",
						Value: "value",
					},
				},
				Method:  "POST",
				Options: "--some-options",
				Payload: "some payload",
				URL:     "http://local:8080",
			},
			want:    "curl --request POST --header 'key: value' --data 'some payload' --some-options http://local:8080",
			wantErr: false,
		},
		{
			name: "valid beta input #2",
			data: lib.Request{
				Headers: []lib.Header{
					{
						Key:   "key",
						Value: "value",
					},
				},
				Method: "POST",
				URL:    "http://local:8080",
			},
			want:    "curl --request POST --header 'key: value' http://local:8080",
			wantErr: false,
		},
		{
			name: "valid beta input #3",
			data: lib.Request{
				Method: "POST",
				URL:    "http://local:8080",
			},
			want:    "curl --request POST http://local:8080",
			wantErr: false,
		},
		{
			name:    "invalid input",
			data:    1,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := taskHandler.CreateRequest(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
