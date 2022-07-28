package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

var configurationService *httptest.Server

const natsTestPort = 8370
const mongoDBVersion = "4.4.9"

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
	defer natsServer.Shutdown()

	go main()
	m.Run()
}

func Test_SLIWrongFinishedPayloadSend(t *testing.T) {
	keptnContext := "context"
	projectName := "quality-gates-invalid-finish"
	serviceName := "my-service"
	stageName := "dev"

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
	natsClient, err := newTestNatsClient(natsURL, t)
	require.Nil(t, err)

	c := http.Client{
		Timeout: 2 * time.Second,
	}

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

	t.Log("sent evaluation.triggered event")
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/event", bytes.NewBuffer(marshal))
	require.Nil(t, err)

	resp, err := c.Do(req)
	require.Nil(t, err)
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	//expect get-sli.triggered
	//expect evaluation.started event
	//send get-sli.started eventTriggerEvaluation
	//send invalid get-sli.finished event
	//expect fail evaluation.finished event

	go func() {
		natsClient.Close()
	}()
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
