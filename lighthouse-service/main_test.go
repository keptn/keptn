package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
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
	"github.com/keptn/go-utils/pkg/sdk/connector/types"
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

var configurationService *httptest.Server

var EventChan chan types.EventUpdate

const natsTestPort = 8370

var keptnContext = "context"
var projectName = "quality-gates-invalid-finish"
var serviceName = "my-service"
var stageName = "dev"

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

func TestMain(m *testing.M) {
	test := testing.T{}

	natsServer := setupNatsServer(natsTestPort, test.TempDir())
	defer startFakeConfigurationService()()
	defer natsServer.Shutdown()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	natsConnector := gonats.New(fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort))
	gonats.WithLogger(log)(natsConnector)
	eventSource := eventsource.New(natsConnector, eventsource.WithLogger(log))

	subscriptionSource := subscriptionsource.NewFixedSubscriptionSource(
		subscriptionsource.WithFixedSubscriptions(
			models.EventSubscription{Event: "sh.keptn.event.evaluation.triggered"},
			models.EventSubscription{Event: "sh.keptn.event.get-sli.finished"},
			models.EventSubscription{Event: "sh.keptn.event.monitoring.configure"},
		),
	)

	logHandler := &gofake.LogAPIMock{}
	logForwarder := logforwarder.New(logHandler)

	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder, controlplane.WithLogger(log))

	os.Setenv("RESOURCE_SERVICE", configurationService.URL)
	defer os.Unsetenv("RESOURCE_SERVICE")

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
		GetEventsFunc: func(filter *keptnapi.EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
			return []*models.KeptnContextExtendedCE{
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

	go _main(controlPlane, log, envConfig{ConfigurationServiceURL: configurationService.URL, LogLevel: logrus.DebugLevel.String(), KubeAPI: fakeK8sClient, EventStore: mockedEventStore})
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

func Test_ErroredFinishedPayloadSend(t *testing.T) {
	natsClient := qualityGatesGenericTestStart(t)

	t.Log("sending invalid get-sli.finished event")
	payload := models.KeptnContextExtendedCE{
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
	var evaluationFinishedEvent *models.KeptnContextExtendedCE
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
	require.Equal(t, keptnv2.ResultType("fail"), evaluationFinishedPayload.EventData.Result)
	require.Equal(t, "evaluation performed by lighthouse received an unexpected error: some-silly-msg", evaluationFinishedPayload.EventData.Message)

	go func() {
		natsClient.Close()
	}()
}

func Test_AbortedFinishedPayloadSend(t *testing.T) {
	natsClient := qualityGatesGenericTestStart(t)

	t.Log("sending invalid get-sli.finished event")
	payload := models.KeptnContextExtendedCE{
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
	var evaluationFinishedEvent *models.KeptnContextExtendedCE
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
	require.Equal(t, keptnv2.ResultType("fail"), evaluationFinishedPayload.EventData.Result)
	require.Equal(t, "evaluation performed by lighthouse was aborted", evaluationFinishedPayload.EventData.Message)

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
	var evaluationStartedEvent *models.KeptnContextExtendedCE
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
	var getSLITriggeredEvent *models.KeptnContextExtendedCE
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
	payload = models.KeptnContextExtendedCE{
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
		//n.t.Logf(string(msg.Data))
		n.receivedEvents = append(n.receivedEvents, *ev)
	}
}

func (n *testNatsClient) triggerSequence(projectName, serviceName, stageName, sequenceName string) *apimodels.EventContext {
	source := "golang-test"
	eventType := keptnv2.GetTriggeredEventType(stageName + "." + sequenceName)
	n.t.Log("triggering task sequence")

	keptnContext := uuid.NewString()

	eventPayload := apimodels.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.EventData{
			Project: projectName,
			Stage:   stageName,
			Service: serviceName,
			Result:  keptnv2.ResultPass,
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: "0.2.0",
		Source:             &source,
		Specversion:        "1.0",
		Type:               &eventType,
	}

	marshal, err := json.Marshal(eventPayload)
	require.Nil(n.t, err)

	err = n.Publish(eventType, marshal)

	return &apimodels.EventContext{
		KeptnContext: &keptnContext,
	}
}

func (n *testNatsClient) SendEvent(event cloudevents.Event) error {
	m, _ := json.Marshal(event)
	return n.Publish(event.Type(), m)
}

func (n *testNatsClient) Send(ctx context.Context, event cloudevents.Event) error {
	return n.SendEvent(event)
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

func verifySequenceEndsUpInState(t *testing.T, projectName string, context *apimodels.EventContext, timeout time.Duration, desiredStates []string) {
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, desiredStates)

	require.Eventually(t, func() bool {
		states, err := getStates(projectName, context)
		if err != nil {
			return false
		}

		for _, state := range states.States {
			if doesSequenceHaveOneOfTheDesiredStates(state, context, desiredStates) {
				return true
			}
		}
		return false
	}, timeout, 100*time.Millisecond)
}

func getStates(projectName string, context *apimodels.EventContext) (*apimodels.SequenceStates, error) {
	c := http.Client{}

	var reqURL string
	if context != nil {
		reqURL = "http://localhost:8080/v1/sequence/" + projectName + "?keptnContext=" + *context.KeptnContext
	} else {
		reqURL = "http://localhost:8080/v1/sequence/" + projectName
	}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	states := &apimodels.SequenceStates{}

	err = json.Unmarshal(respBytes, states)
	if err != nil {
		return nil, err
	}
	return states, nil
}

func doesSequenceHaveOneOfTheDesiredStates(state apimodels.SequenceState, context *apimodels.EventContext, desiredStates []string) bool {
	if state.Shkeptncontext == *context.KeptnContext {
		for _, desiredState := range desiredStates {
			if state.State == desiredState {
				return true
			}
		}
	}
	return false
}

func getStageOfState(state apimodels.SequenceState, stageName string) *apimodels.SequenceStateStage {
	for index, stage := range state.Stages {
		if stage.Name == stageName {
			return &state.Stages[index]
		}
	}
	return nil
}

func controlSequence(t *testing.T, projectName, keptnContextID string, cmd apimodels.SequenceControlState) {
	command := apimodels.SequenceControlCommand{
		State: cmd,
	}

	mCmd, _ := json.Marshal(command)

	c := http.Client{}
	_, err := c.Post(fmt.Sprintf("http://localhost:8080/v1/sequence/%s/%s/control", projectName, keptnContextID), "application/json", bytes.NewBuffer(mCmd))
	require.Nil(t, err)
}

func setupFakeConfigurationService() {
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
			encodedSLO := base64.StdEncoding.EncodeToString([]byte(qualityGatesShortSLOFileContent))
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
