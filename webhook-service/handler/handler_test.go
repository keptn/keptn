package handler_test

import (
	"encoding/json"
	"errors"
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
	"io/ioutil"
	"log"
	"reflect"
	"testing"
	"time"
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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

const webHookContentBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
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
		- url: http://local:8080
		  method: GET
		  headers:
            - key: x-token
              value: "{{.env.secretKey}}"
			  - key: project
              value: "{{.data.project}}"`

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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

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
        - "curl http://local:8080 {{.unavailable}} {{.env.mysecret}}"`

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
        - "curl http://local:8080 {{.data.project}} {{.env.mysecret}}"`

func newWebhookTriggeredEvent(filename string) models.KeptnContextExtendedCE {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	event := models.KeptnContextExtendedCE{}
	err = json.Unmarshal(content, &event)
	_ = err
	return event
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn("test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()

	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))
	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 1 }, 30*time.Second, time.Millisecond*10)
	require.Eventually(t, func() bool {
		return curlExecutorMock.CurlCalls()[0].CurlCmd == "curl http://local:8080 myproject my-secret-value"
	}, 30*time.Second, time.Millisecond*10)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, keptnv2.GetStartedEventType("webhook"))
	fakeKeptn.AssertSentEventType(t, 1, keptnv2.GetFinishedEventType("webhook"))
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn("test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithStartedEvent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.started", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)
	err := fakeKeptn.Start()
	if err != nil {
		t.Fatal(err)
	}
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.started.json"))
	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 1 }, 30*time.Second, time.Millisecond*10)
	require.Eventually(t, func() bool {
		return curlExecutorMock.CurlCalls()[0].CurlCmd == "curl http://local:8080 myproject my-secret-value"
	}, 30*time.Second, time.Millisecond*10)
	fakeKeptn.AssertNumberOfEventSent(t, 0)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithStartedEvent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.started", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.started.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 1 }, 30*time.Second, time.Millisecond*10)
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 0)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithFinishedEvent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.finished", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.finished.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 1 }, 30*time.Second, time.Millisecond*10)
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 0)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithFinishedEvent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.finished", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)
	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.finished.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 1 }, 30*time.Second, time.Millisecond*10)
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 0)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequests})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 2 }, 30*time.Second, time.Millisecond*10)
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)

	////verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequestsAndDisabledFinished})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 4 }, 30*time.Second, time.Millisecond*10)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[2].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[3].CurlCmd)

	////verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 4)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 2, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 3, "sh.keptn.event.webhook.started")

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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequestsAndDisabledStarted})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 4 }, 30*time.Second, time.Millisecond*10)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[2].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[3].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 0)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMultipleRequestsAndDisabledFinished})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 2 }, 30*time.Second, time.Millisecond*10)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)
	assert.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[1].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 7)
	//require.Equal(t, 7, len(fakeKeptn.GetEventSender().SentEvents))
	// each request should have a .started event
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 2, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 3, "sh.keptn.event.webhook.started")
	// apart from the first request which has been successful, every request should have a failed .finished event
	fakeKeptn.AssertSentEventType(t, 4, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventType(t, 5, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventType(t, 6, "sh.keptn.event.webhook.finished")

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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithNoMatchingSubscriptionID})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool { return len(curlExecutorMock.CurlCalls()) == 0 }, 30*time.Second, time.Millisecond*10)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)
}

func TestTaskHandler_Execute_WebhookCannotBeRetrieved(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.FailingResourceHandler{})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)

}

func TestTaskHandler_Execute_NoSubscriptionIDInEvent(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.FailingResourceHandler{})
	fakeKeptn.AddTaskHandler("sh.keptn.event.webhook.triggered", taskHandler)
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-no-subscription-id.json"))

	fakeKeptn.AssertNumberOfEventSent(t, 0)

}

func TestTaskHandler_Execute_InvalidEvent(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.FailingResourceHandler{})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/invalid-event.json"))

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 0)
	require.Empty(t, curlExecutorMock.CurlCalls())
}

func TestTaskHandler_CannotReadSecret(t *testing.T) {
	templateEngineMock := &fake.ITemplateEngineMock{}
	secretReaderMock := &fake.ISecretReaderMock{}
	secretReaderMock.ReadSecretFunc = func(name string, key string) (string, error) {
		return "", errors.New("unable to read secret :(")
	}
	curlExecutorMock := &fake.ICurlExecutorMock{}
	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn("test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentWithMissingTemplateData})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool {
		return len(secretReaderMock.ReadSecretCalls()) > 0
	}, 30*time.Second, 10*time.Millisecond)
	require.NotEmpty(t, templateEngineMock.ParseTemplateCalls())
	require.Empty(t, curlExecutorMock.CurlCalls())

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)
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
	requestValidatorMock := &fake.RequestValidatorMock{}
	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool {
		return len(secretReaderMock.ReadSecretCalls()) > 0
	}, 30*time.Second, 10*time.Millisecond)
	require.NotEmpty(t, templateEngineMock.ParseTemplateCalls())
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)

}

func TestTaskHandler_RequestValidatorFails(t *testing.T) {
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
	requestValidatorMock := &fake.RequestValidatorMock{}
	requestValidatorMock.ValidateFunc = func(request lib.Request) error {
		return errors.New("validation failed")
	}
	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContentBeta})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool {
		return len(secretReaderMock.ReadSecretCalls()) == 0
	}, 30*time.Second, 10*time.Millisecond)
	require.Empty(t, templateEngineMock.ParseTemplateCalls())

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)
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
	requestValidatorMock := fake.RequestValidatorMock{}
	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")
	fakeKeptn.SetResourceHandler(sdk.StringResourceHandler{ResourceContent: webHookContent})
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	require.Eventually(t, func() bool {
		return len(secretReaderMock.ReadSecretCalls()) > 0
	}, 30*time.Second, 10*time.Millisecond)
	require.NotEmpty(t, templateEngineMock.ParseTemplateCalls())
	require.Equal(t, "curl http://local:8080 myproject my-secret-value", curlExecutorMock.CurlCalls()[0].CurlCmd)

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")

	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultFailed)
	eventData := &keptnv2.EventData{}
	keptnv2.EventDataAs(fakeKeptn.TestEventSource.SentEvents[1], eventData)
	require.NotContains(t, eventData.Message, "my-secret-value")
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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")

	fakeKeptn.SetResourceHandler(resourceHandlerMock)
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)

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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")

	fakeKeptn.SetResourceHandler(resourceHandlerMock)
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)

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

	requestValidatorMock := &fake.RequestValidatorMock{}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

	fakeKeptn := sdk.NewFakeKeptn(
		"test-webhook-svc")

	fakeKeptn.SetResourceHandler(resourceHandlerMock)
	fakeKeptn.AddTaskHandlerWithSubscriptionID("sh.keptn.event.webhook.triggered", taskHandler, "my-subscription-id")
	fakeKeptn.SetAutomaticResponse(false)

	fakeKeptn.Start()
	fakeKeptn.NewEvent(newWebhookTriggeredEvent("test/events/test-webhook.triggered-0.json"))

	//verify sent events
	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.webhook.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.webhook.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, keptnv2.StatusSucceeded)
	fakeKeptn.AssertSentEventResult(t, 1, keptnv2.ResultPass)

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
	requestValidatorMock := &fake.RequestValidatorMock{}
	requestValidatorMock.ValidateFunc = func(request lib.Request) error {
		return nil
	}

	taskHandler := handler.NewTaskHandler(templateEngineMock, curlExecutorMock, requestValidatorMock, secretReaderMock)

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
			name:    "invalid alpha input #1",
			data:    "curl http:localhost",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid alpha input #2",
			data:    "curl kubernetes.svc",
			want:    "",
			wantErr: true,
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
