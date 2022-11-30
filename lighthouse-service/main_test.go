package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk/connector/controlplane"
	eventsource "github.com/keptn/go-utils/pkg/sdk/connector/eventsource/nats"
	gofake "github.com/keptn/go-utils/pkg/sdk/connector/fake"
	"github.com/keptn/go-utils/pkg/sdk/connector/logforwarder"
	gonats "github.com/keptn/go-utils/pkg/sdk/connector/nats"
	"github.com/keptn/go-utils/pkg/sdk/connector/subscriptionsource"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
	lighthousefake "github.com/keptn/keptn/lighthouse-service/event_handler/fake"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakek8s "k8s.io/client-go/kubernetes/fake"
)

const mongoDBGetByTypeContent = `{
  "contenttype": "application/json",
  "data": {
    "evaluation": {
      "comparedEvents": [
        "b7ab5914-b24a-4f8d-8af3-eb98240de3cf"
      ],
      "indicatorResults": [
        {
          "displayName": "",
          "keySli": false,
          "passTargets": [
            {
              "criteria": "<=+75%",
              "targetValue": 350,
              "violated": false
            },
            {
              "criteria": "<800",
              "targetValue": 800,
              "violated": false
            }
          ],
          "score": 1,
          "status": "pass",
          "value": {
            "comparedValue": 200,
            "metric": "response_time_p95",
            "success": true,
            "value": 200
          },
          "warningTargets": [
            {
              "criteria": "<=1000",
              "targetValue": 1000,
              "violated": false
            },
            {
              "criteria": "<=+100%",
              "targetValue": 400,
              "violated": false
            }
          ]
        },
        {
          "displayName": "",
          "keySli": false,
          "passTargets": [
            {
              "criteria": "<=+100%",
              "targetValue": 400,
              "violated": false
            },
            {
              "criteria": ">=-80%",
              "targetValue": 40,
              "violated": false
            }
          ],
          "score": 1,
          "status": "pass",
          "value": {
            "comparedValue": 200,
            "metric": "throughput",
            "success": true,
            "value": 200
          },
          "warningTargets": null
        },
        {
          "displayName": "",
          "keySli": false,
          "passTargets": null,
          "score": 0,
          "status": "info",
          "value": {
            "comparedValue": 0,
            "metric": "error_rate",
            "success": true,
            "value": 0
          },
          "warningTargets": null
        }
      ],
      "result": "pass",
      "score": 100,
      "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjAuMS4xIgpjb21wYXJpc29uOgogIGFnZ3JlZ2F0ZV9mdW5jdGlvbjogImF2ZyIKICBjb21wYXJlX3dpdGg6ICJzaW5nbGVfcmVzdWx0IgogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIgogIG51bWJlcl9vZl9jb21wYXJpc29uX3Jlc3VsdHM6IDEKZmlsdGVyOgpvYmplY3RpdmVzOgogIC0gc2xpOiAicmVzcG9uc2VfdGltZV9wOTUiCiAgICBrZXlfc2xpOiBmYWxzZQogICAgcGFzczogICAgICAgICAgICAgIyBwYXNzIGlmIChyZWxhdGl2ZSBjaGFuZ2UgPD0gNzUlIEFORCBhYnNvbHV0ZSB2YWx1ZSBpcyA8IDc1bXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSs3NSUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw4MDAiICAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yCiAgICB3YXJuaW5nOiAgICAgICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDIwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9MTAwMCIKICAgICAgICAgIC0gIjw9KzEwMCUiCiAgICB3ZWlnaHQ6IDEKICAtIHNsaTogInRocm91Z2hwdXQiCiAgICBwYXNzOgogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD0rMTAwJSIKICAgICAgICAgIC0gIj49LTgwJSIKICAtIHNsaTogImVycm9yX3JhdGUiCnRvdGFsX3Njb3JlOgogIHBhc3M6ICIxMDAlIgogIHdhcm5pbmc6ICI2NSUi",
      "timeEnd": "2022-01-26T10:10:53.931Z",
      "timeStart": "2022-01-26T10:05:53.931Z"
    },
    "project": "keptn-quality-gates",
    "result": "pass",
    "service": "my-service",
    "stage": "hardening",
    "status": "succeeded",
    "temporaryData": {
      "distributor": {
        "subscriptionID": ""
      }
    }
  },
  "id": "c09c80ae-a70d-4f1f-b2ff-f09f4ca38bdd",
  "shkeptncontext": "b95e0772-053a-4722-9feb-fc7dfa435872",
  "shkeptnspecversion": "0.2.4",
  "source": "lighthouse-service",
  "specversion": "1.0",
  "time": "2022-08-31T16:14:10.531773605Z",
  "triggeredid": "92ced52e-9273-472c-b48b-b4c92ef37211",
  "gitcommitid": "ddeda67224e97244e28bd74c322f9919ef650c53",
  "type": "sh.keptn.event.evaluation.finished"
}`

const qualityGatesShortSLOFileContent = `---
spec_version: "0.1.1"
comparison:
  aggregate_function: "avg"
  compare_with: "single_result"
  include_result_with_score: "pass"
  number_of_comparison_results: 1
filter:
objectives:
  - sli: "response_time_p95"
    key_sli: false
    pass:             # pass if (relative change <= 75% AND absolute value is < 75ms)
      - criteria:
          - "<=+75%"  # relative values require a prefixed sign (plus or minus)
          - "<800"     # absolute values only require a logical operator
    warning:          # if the response time is below 200ms, the result should be a warning
      - criteria:
          - "<=1000"
          - "<=+100%"
    weight: 1
  - sli: "throughput"
    pass:
      - criteria:
          - "<=+100%"
          - ">=-80%"
  - sli: "error_rate"
total_score:
  pass: "100%"
  warning: "65%"`

const noObjectivesSLOFileContent = `---
spec_version: "1.0"
comparison:
  aggregate_function: "avg"
  compare_with: "single_result"
  include_result_with_score: "pass"
  number_of_comparison_results: 1
filter:
objectives:
total_score:
  pass: "90%"
  warning: "75%"`

const natsTestPort = 8370

var keptnContext = "context"
var projectName = "quality-gates-invalid-finish"
var serviceName = "my-service"
var stageName = "dev"
var configurationService *httptest.Server
var mongodbService *httptest.Server

func TestMain(m *testing.M) {
	test := testing.T{}

	natsServer := setupNatsServer(natsTestPort, test.TempDir())
	defer startFakeConfigurationService()()
	defer startFakeMongoDBService()()
	defer natsServer.Shutdown()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	natsConnector := gonats.New(fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort))
	gonats.WithLogger(log)(natsConnector)
	eventSource := eventsource.New(natsConnector, eventsource.WithLogger(log))

	subscriptionSource := subscriptionsource.NewFixedSubscriptionSource(
		subscriptionsource.WithFixedSubscriptions(
			apimodels.EventSubscription{Event: "sh.keptn.event.evaluation.triggered"},
			apimodels.EventSubscription{Event: "sh.keptn.event.get-sli.finished"},
			apimodels.EventSubscription{Event: "sh.keptn.event.monitoring.configure"},
		),
	)

	logHandler := &gofake.LogAPIMock{}
	logForwarder := logforwarder.New(logHandler)

	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder, controlplane.WithLogger(log))

	test.Setenv("RESOURCE_SERVICE", configurationService.URL)
	mongo := strings.TrimPrefix(mongodbService.URL, "http://")
	test.Setenv("MONGODB_DATASTORE", mongo)

	fakeK8sClient := fakek8s.NewSimpleClientset(
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "lighthouse-config",
			},
			Data: map[string]string{
				"sli-provider": "my-sli-provider"},
		},
	)

	mockedEventStore := &lighthousefake.EventStoreMock{
		GetEventsFunc: func(filter *keptnapi.EventFilter) ([]*apimodels.KeptnContextExtendedCE, *apimodels.Error) {
			return []*apimodels.KeptnContextExtendedCE{
				{
					Contenttype: "application/json",
					Data: keptnv2.EvaluationTriggeredEventData{
						EventData: keptnv2.EventData{
							Project: projectName,
							Stage:   stageName,
							Service: serviceName,
						},
						Test: keptnv2.Test{},
						Evaluation: keptnv2.Evaluation{
							End:       "2022-01-26T10:10:53.931Z",
							Start:     "2022-01-26T10:05:53.931Z",
							Timeframe: "",
						},
						Deployment: keptnv2.Deployment{},
					},
					Extensions:         nil,
					ID:                 uuid.NewString(),
					Shkeptncontext:     keptnContext,
					Shkeptnspecversion: "0.2.3",
					Source:             strutils.Stringp("fakeshipyard"),
					Specversion:        "1.0",
					Time:               time.Now(),
					Triggeredid:        "",
					GitCommitID:        "",
					Type:               strutils.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
				},
			}, nil
		},
	}

	lighthouseService := LighthouseService{
		KubeAPI: fakeK8sClient,
		EventStore: func(k *keptnv2.Keptn) event_handler.EventStore {
			return mockedEventStore
		},
		env: envConfig{
			ConfigurationServiceURL: configurationService.URL,
			LogLevel:                logrus.DebugLevel.String(),
		},
	}

	go _main(controlPlane, log, lighthouseService)
	// need to wait until the subscriptions are ready
	time.Sleep(2 * time.Second)
	m.Run()
}

func startFakeConfigurationService() func() {
	configurationService = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))

	return configurationService.Close
}

func startFakeMongoDBService() func() {
	mongodbService = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mongoDBGetByTypeContent))
	}))

	return mongodbService.Close
}

func Test_ErroredFinishedPayloadSend(t *testing.T) {
	natsClient := qualityGatesGenericTestStart(t)

	t.Log("sending invalid get-sli.finished event")
	payload := apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultPass,
				Message: "some-silly-msg",
			},
			GetSLI: keptnv2.GetSLIFinished{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}

	marshal, err := json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName), marshal)
	require.Nil(t, err)

	t.Log("expecting evaluation.finished event")
	var evaluationFinishedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if event != nil {
			evaluationFinishedEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got evaluation.finished event")
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationFinishedEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, stageName, evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusErrored, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.EventData.Result)
	require.Equal(t, "evaluation performed by lighthouse received an unexpected error: some-silly-msg", evaluationFinishedPayload.EventData.Message)

	go func() {
		natsClient.Close()
	}()
}

func Test_AbortedFinishedPayloadSend(t *testing.T) {
	natsClient := qualityGatesGenericTestStart(t)

	t.Log("sending invalid get-sli.finished event")
	payload := apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusAborted,
				Result:  keptnv2.ResultPass,
				Message: "some-silly-msg",
			},
			GetSLI: keptnv2.GetSLIFinished{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}

	marshal, err := json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName), marshal)
	require.Nil(t, err)

	t.Log("expecting evaluation.finished event")
	var evaluationFinishedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if event != nil {
			evaluationFinishedEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got evaluation.finished event")
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationFinishedEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, stageName, evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusErrored, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.EventData.Result)
	require.Equal(t, "evaluation performed by lighthouse was aborted", evaluationFinishedPayload.EventData.Message)

	go func() {
		natsClient.Close()
	}()
}

func Test_SLIFailResultSend(t *testing.T) {
	natsClient := qualityGatesGenericTestStart(t)

	t.Log("sending invalid get-sli.finished event")
	payload := apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultFailed,
				Message: "some-silly-msg",
			},
			GetSLI: keptnv2.GetSLIFinished{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}

	marshal, err := json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName), marshal)
	require.Nil(t, err)

	t.Log("expecting evaluation.finished event")
	var evaluationFinishedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if event != nil {
			evaluationFinishedEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got evaluation.finished event")
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationFinishedEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, stageName, evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.EventData.Result)
	require.Equal(t, "lighthouse failed because SLI failed with message some-silly-msg", evaluationFinishedPayload.EventData.Message)

	go func() {
		natsClient.Close()
	}()
}

func qualityGatesGenericTestStart(t *testing.T) *testNatsClient {
	setupFakeConfigurationService()

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
	natsClient, err := newTestNatsClient(natsURL, t)
	require.Nil(t, err)

	t.Log("sending evaluation.triggered event")
	payload := apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.EvaluationTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
			},
			Test: keptnv2.Test{},
			Evaluation: keptnv2.Evaluation{
				End:       "2022-01-26T10:10:53.931Z",
				Start:     "2022-01-26T10:05:53.931Z",
				Timeframe: "",
			},
			Deployment: keptnv2.Deployment{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
	}

	marshal, err := json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), marshal)
	require.Nil(t, err)

	t.Log("expecting evaluation.started event")
	var evaluationStartedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName))
		if event != nil {
			evaluationStartedEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got evaluation.started event")
	evaluationStartedPayload := &keptnv2.EvaluationStartedEventData{}
	err = keptnv2.Decode(evaluationStartedEvent.Data, evaluationStartedPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationStartedEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationStartedPayload.EventData.Project)
	require.Equal(t, stageName, evaluationStartedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationStartedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationStartedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), evaluationStartedPayload.EventData.Result)
	require.Empty(t, evaluationStartedPayload.EventData.Message)

	t.Log("expecting get-sli.triggered event")
	var getSLITriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
		if event != nil {
			getSLITriggeredEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got get-sli.triggered event")
	getSLIPayload := &keptnv2.GetSLITriggeredEventData{}
	err = keptnv2.Decode(getSLITriggeredEvent.Data, getSLIPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, getSLITriggeredEvent.Shkeptncontext)
	require.Equal(t, "my-sli-provider", getSLIPayload.GetSLI.SLIProvider)
	require.NotEmpty(t, getSLIPayload.GetSLI.Start)
	require.NotEmpty(t, getSLIPayload.GetSLI.End)
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "response_time_p95")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "throughput")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "error_rate")
	require.Equal(t, projectName, getSLIPayload.EventData.Project)
	require.Equal(t, stageName, getSLIPayload.EventData.Stage)
	require.Equal(t, serviceName, getSLIPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), getSLIPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), getSLIPayload.EventData.Result)
	require.Empty(t, getSLIPayload.EventData.Message)

	t.Log("sending get-sli.started event")
	payload = apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}

	marshal, err = json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName), marshal)
	require.Nil(t, err)

	return natsClient
}

type testNatsClient struct {
	*nats.Conn
	t              *testing.T
	receivedEvents []apimodels.KeptnContextExtendedCE
	sync.RWMutex
}

func setupNatsServer(port int, storeDir string) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	opts.StoreDir = storeDir
	svr := natsserver.RunServer(&opts)

	connect, _ := nats.Connect(svr.ClientURL())

	js, _ := connect.JetStream()

	js.DeleteStream("keptn")

	return svr
}

func newTestNatsClient(natsURL string, t *testing.T) (*testNatsClient, error) {
	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	tnc := &testNatsClient{
		t:    t,
		Conn: natsConn,
	}

	_, err = tnc.Subscribe("sh.keptn.>", func(msg *nats.Msg) {
		tnc.onEvent(msg)
	})
	if err != nil {
		return nil, err
	}

	return tnc, nil
}

func (n *testNatsClient) onEvent(msg *nats.Msg) {
	n.Lock()
	defer n.Unlock()

	n.t.Logf("Received event of type: %s", msg.Subject)
	ev := &apimodels.KeptnContextExtendedCE{}

	if err := json.Unmarshal(msg.Data, ev); err == nil {
		n.receivedEvents = append(n.receivedEvents, *ev)
	}
}

func (n *testNatsClient) getLatestEventOfType(keptnContext, projectName, stage, eventType string) *apimodels.KeptnContextExtendedCE {
	var result *apimodels.KeptnContextExtendedCE
	n.Lock()
	defer n.Unlock()
	for index := range n.receivedEvents {
		if n.receivedEvents[index].Shkeptncontext == keptnContext && *n.receivedEvents[index].Type == eventType {
			ed := &keptnv2.EventData{}
			err := keptnv2.Decode(n.receivedEvents[index].Data, ed)
			require.Nil(n.t, err)
			if ed.Project == projectName && ed.Stage == stage {
				result = &n.receivedEvents[index]
			}
		}
	}
	return result
}

func setupFakeConfigurationService() {
	_setupFakeConfigurationService(false)
}

func _setupFakeConfigurationService(noObjectives bool) {
	configurationService.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if strings.Contains(r.RequestURI, "/metadata.yaml") {
			res := apimodels.Resource{
				Metadata: &apimodels.Version{
					Version: "my-commit-id",
				},
			}

			marshal, _ := json.Marshal(res)
			w.Write(marshal)

			return
		} else if strings.Contains(r.RequestURI, "/slo.yaml") {
			w.WriteHeader(200)
			var encodedSLO string
			if noObjectives {
				encodedSLO = base64.StdEncoding.EncodeToString([]byte(noObjectivesSLOFileContent))
			} else {
				encodedSLO = base64.StdEncoding.EncodeToString([]byte(qualityGatesShortSLOFileContent))
			}
			res := apimodels.Resource{
				ResourceContent: encodedSLO,
			}

			marshal, _ := json.Marshal(res)
			w.Write(marshal)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
}

func Test_NoSLOObjectives(t *testing.T) {
	_setupFakeConfigurationService(true)

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
	natsClient, err := newTestNatsClient(natsURL, t)
	require.Nil(t, err)

	t.Log("sending evaluation.triggered event")
	payload := apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.EvaluationTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
			},
			Test: keptnv2.Test{},
			Evaluation: keptnv2.Evaluation{
				End:       "2022-01-26T10:10:53.931Z",
				Start:     "2022-01-26T10:05:53.931Z",
				Timeframe: "",
			},
			Deployment: keptnv2.Deployment{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
	}

	marshal, err := json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), marshal)
	require.Nil(t, err)

	t.Log("expecting evaluation.started event")
	var evaluationStartedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName))
		if event != nil {
			evaluationStartedEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got evaluation.started event")
	evaluationStartedPayload := &keptnv2.EvaluationStartedEventData{}
	err = keptnv2.Decode(evaluationStartedEvent.Data, evaluationStartedPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationStartedEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationStartedPayload.EventData.Project)
	require.Equal(t, stageName, evaluationStartedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationStartedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationStartedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), evaluationStartedPayload.EventData.Result)
	require.Empty(t, evaluationStartedPayload.EventData.Message)

	t.Log("expecting get-sli.triggered event")
	var getSLITriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
		if event != nil {
			getSLITriggeredEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got get-sli.triggered event")
	getSLIPayload := &keptnv2.GetSLITriggeredEventData{}
	err = keptnv2.Decode(getSLITriggeredEvent.Data, getSLIPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, getSLITriggeredEvent.Shkeptncontext)
	require.Equal(t, "my-sli-provider", getSLIPayload.GetSLI.SLIProvider)
	require.NotEmpty(t, getSLIPayload.GetSLI.Start)
	require.NotEmpty(t, getSLIPayload.GetSLI.End)
	require.Equal(t, projectName, getSLIPayload.EventData.Project)
	require.Equal(t, stageName, getSLIPayload.EventData.Stage)
	require.Equal(t, serviceName, getSLIPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), getSLIPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), getSLIPayload.EventData.Result)
	require.Empty(t, getSLIPayload.EventData.Message)

	t.Log("sending get-sli.started event")
	payload = apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}

	marshal, err = json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName), marshal)
	require.Nil(t, err)

	t.Log("sending invalid get-sli.finished event")
	payload = apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultFailed,
				Message: "no SLIs were requested",
			},
			GetSLI: keptnv2.GetSLIFinished{
				IndicatorValues: []*keptnv2.SLIResult{
					{
						Message: "no SLIs were requested",
						Metric:  "no metric",
						Success: false,
						Value:   0,
					},
				},
			},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.3",
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}

	marshal, err = json.Marshal(payload)
	require.Nil(t, err)

	err = natsClient.Publish(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName), marshal)
	require.Nil(t, err)

	t.Log("expecting evaluation.finished event")
	var evaluationFinishedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if event != nil {
			evaluationFinishedEvent = event
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)

	t.Log("got evaluation.finished event")
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationFinishedEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, stageName, evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.EventData.Result)
	require.Equal(t, "lighthouse failed because no SLO objective was provided", evaluationFinishedPayload.EventData.Message)

	go func() {
		natsClient.Close()
	}()
}
