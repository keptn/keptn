package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	fakek8s "k8s.io/client-go/kubernetes/fake"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const natsTestPort = 8370
const mongoDBVersion = "4.4.9"

const testShipyardFile = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: delivery
      tasks:
      - name: mytask
        properties:  
          strategy: direct
      - name: test
        properties:
          kind: functional
      - name: evaluation 
      - name: release 
    - name: rollback
      tasks:
      - name: rollback
      triggeredOn:
        - event: dev.artifact-delivery.finished
          selector:
            match:
              result: fail
    - name: delivery-with-approval
      tasks:
      - name: approval
        properties:
          pass: manual
          warning: manual
      - name: mytask
    
  - name: hardening
    sequences:
    - name: artifact-delivery
      triggeredOn:
        - event: dev.artifact-delivery.finished
      tasks:
      - name: deployment
        properties: 
          strategy: blue_green_service
      - name: test
        properties:  
          kind: performance
      - name: evaluation
      - name: release

  - name: production
    sequences:
    - name: artifact-delivery 
      triggeredOn:
        - event: hardening.artifact-delivery.finished
      tasks:
      - name: deployment
        properties:
          strategy: blue_green
      - name: release
      
    - name: remediation
      tasks:
      - name: remediation
      - name: evaluation`

const testShipyardFileWithApprovalSequence = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: delivery
      tasks:
      - name: mytask

    - name: delivery-with-approval
      tasks:
      - name: approval
        properties:
          pass: manual
          warning: manual
      - name: mytask`

const sequenceStateParallelStagesShipyard = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: shipyard-parallel-stages
spec:
  stages:
    - name: dev
      sequences:
        - name: delivery
          tasks:
            - name: delivery
    - name: staging-2
      sequences:
        - name: delivery
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: delivery
    - name: staging-1
      sequences:
        - name: delivery
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: delivery`

const sequenceTimeoutWithTriggeredAfterShipyard = `--- 
apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata: 
  name: shipyard-sockshop
spec: 
  stages: 
    - 
      name: dev
      sequences: 
        - 
          name: delivery
          tasks: 
            - 
              triggeredAfter: "10s"
              name: unknown`

var configurationService *httptest.Server

func setupNatsServer(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	svr := natsserver.RunServer(&opts)

	connect, _ := nats.Connect(svr.ClientURL())

	js, _ := connect.JetStream()

	js.DeleteStream("keptn")

	return svr
}

func setupLocalMongoDB() func() {
	mongoServer, err := memongo.Start(mongoDBVersion)
	randomDbName := memongo.RandomDatabase()

	os.Setenv("MONGODB_DATABASE", randomDbName)
	os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", fmt.Sprintf("%s/%s", mongoServer.URI(), randomDbName))

	var mongoDBClient *mongo.Client
	mongoDBClient, err = mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		log.Fatalf("Mongo Client setup failed: %s", err)
	}
	err = mongoDBClient.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Mongo Server setup failed: %s", err)
	}

	fmt.Println(fmt.Sprintf("MongoDB Server: %s", mongoServer.URI()))

	return func() { mongoServer.Stop() }
}

func startFakeConfigurationService() func() {
	configurationService = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))

	return configurationService.Close
}

func TestMain(m *testing.M) {

	flag.Parse()
	if testing.Short() {
		return
	}
	natsServer := setupNatsServer(natsTestPort)
	defer startFakeConfigurationService()()
	defer natsServer.Shutdown()
	defer setupLocalMongoDB()()
	fakeClient := fakek8s.NewSimpleClientset()

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
	os.Setenv(envVarDisableLeaderElection, "true")
	os.Setenv(envVarConfigurationSvcEndpoint, configurationService.URL)
	os.Setenv(envVarNatsURL, natsURL)
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv(envVarEventDispatchIntervalSec, "1")
	os.Setenv(envVarSequenceDispatchIntervalSec, "1s")
	os.Setenv(envVarTaskStartedWaitDuration, "10s")

	go _main(fakeClient)
	m.Run()
}

func Test_getDurationFromEnvVar(t *testing.T) {
	type args struct {
		envVarValue string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "get default value",
			args: args{
				envVarValue: "",
			},
			want: 432000 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "10s",
			},
			want: 10 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "2m",
			},
			want: 120 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "1h30m",
			},
			want: 5400 * time.Second,
		},
		{
			name: "get default value because of invalid config",
			args: args{
				envVarValue: "invalid",
			},
			want: 432000 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LOG_TTL", tt.args.envVarValue)
			if got := getDurationFromEnvVar("LOG_TTL", envVarLogsTTLDefault); got != tt.want {
				t.Errorf("getLogTTLDurationInSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__main_SequenceQueue(t *testing.T) {
	projectName := "my-project-queue"
	serviceName := "my-service"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()

	source := "golang-test"

	context := natsClient.triggerSequence(projectName, serviceName, "dev", "delivery")

	t.Logf("wait for the sequence state to be available")
	verifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{apimodels.SequenceStartedState})
	t.Log("received the expected state!")

	t.Logf("trigger a second sequence - this one should stay in 'waiting' state until the previous sequence is finished")
	secondContext := natsClient.triggerSequence(projectName, serviceName, "dev", "delivery")

	t.Logf("checking if the second sequence is in state 'waiting'")
	verifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{apimodels.SequenceWaitingState})
	t.Log("received the expected state!")

	t.Logf("check if mytask.triggered has been sent for first sequence - this one should be available")
	triggeredEventOfFirstSequence := natsClient.getLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.NotNil(t, triggeredEventOfFirstSequence)

	t.Logf("check if mytask.triggered has been sent for second sequence - this one should NOT be available")
	triggeredEventOfSecondSequence := natsClient.getLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, triggeredEventOfSecondSequence)

	t.Logf("send .started and .finished event for task of first sequence")
	cloudEvent := keptnv2.ToCloudEvent(*triggeredEventOfFirstSequence)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)

	t.Logf("send started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	t.Logf("send finished event with result=fail")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultFailed,
	}, source)
	require.Nil(t, err)

	t.Logf("now that all tasks for the first sequence have been executed, the second sequence should eventually have the status 'started'")
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, apimodels.SequenceStartedState)
	verifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{apimodels.SequenceStartedState})
	t.Log("received the expected state!")

	t.Logf("check if mytask.triggered has been sent for second sequence - now it should be available")
	require.Eventually(t, func() bool {
		triggeredEventOfSecondSequence = natsClient.getLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
		return triggeredEventOfSecondSequence != nil
	}, 2*time.Second, 100*time.Millisecond)

}

func Test__main_SequenceQueueWithTimeout(t *testing.T) {
	projectName := "my-project-timeout"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()

	t.Logf("trigger the first task sequence - this should time out")
	firstSequenceContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)

	verifySequenceEndsUpInState(t, projectName, firstSequenceContext, 20*time.Second, []string{apimodels.TimedOut})
	t.Log("received the expected state!")

	t.Logf("now trigger the second sequence - this should start and a .triggered event for mytask should be sent")
	secondContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	t.Logf("waiting for state with keptnContext %s to have the status %s", *secondContext.KeptnContext, apimodels.SequenceStartedState)
	verifySequenceEndsUpInState(t, projectName, secondContext, 10*time.Second, []string{apimodels.SequenceStartedState})
	triggeredEventOfSecondSequence := natsClient.getLatestEventOfType(*secondContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)
}

func Test__main_SequenceQueueApproval(t *testing.T) {
	projectName := "my-project-queue-with-approval"
	stageName := "dev"
	serviceName := "my-service"

	source := "shippy-test"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFileWithApprovalSequence)

	defer tearDown()

	context := natsClient.triggerSequence(projectName, serviceName, stageName, "delivery-with-approval")
	verifySequenceEndsUpInState(t, projectName, context, 10*time.Second, []string{apimodels.SequenceStartedState})

	t.Logf("check if approval.triggered has been sent for sequence - now it should be available")
	approvalTriggeredEvent := natsClient.getLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName))
	require.NotNil(t, approvalTriggeredEvent)

	t.Logf("send the approval.started event to make sure the sequence will not be cancelled due to a timeout")
	approvalTriggeredCE := keptnv2.ToCloudEvent(*approvalTriggeredEvent)
	keptnHandler, err := keptnv2.NewKeptn(&approvalTriggeredCE, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)

	_, err = keptnHandler.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	t.Logf("now let's trigger the other sequence")
	secondContext := natsClient.triggerSequence(projectName, serviceName, stageName, "delivery")
	verifySequenceEndsUpInState(t, projectName, secondContext, 10*time.Second, []string{apimodels.SequenceStartedState})

	t.Logf("check if approval.triggered has been sent for sequence - now it should be available")
	myTaskTriggeredEvent := natsClient.getLatestEventOfType(*secondContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	myTaskCE := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)
	secondKeptnHandler, err := keptnv2.NewKeptn(&myTaskCE, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)

	t.Logf("send the mytask.started event")
	_, err = secondKeptnHandler.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	t.Logf("now let's send the approval.finished event - the next task should now be queued until the other sequence has been finished")
	_, err = keptnHandler.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultPass, Status: keptnv2.StatusSucceeded}, source)
	require.Nil(t, err)

	t.Logf("wait a bit to make sure mytask.triggered of first sequence is not sent")
	<-time.After(3 * time.Second)
	myTaskTriggeredEventOfFirstSequence := natsClient.getLatestEventOfType(*context.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, myTaskTriggeredEventOfFirstSequence)

	t.Logf("now let's finish mytask of the second sequence")
	_, err = secondKeptnHandler.SendTaskFinishedEvent(&keptnv2.EventData{Status: keptnv2.StatusSucceeded, Result: keptnv2.ResultPass}, source)
	require.Nil(t, err)

	t.Logf("this should have completed the task sequence")
	verifySequenceEndsUpInState(t, projectName, secondContext, 10*time.Second, []string{apimodels.SequenceFinished})

	t.Logf("now the mytask.triggered event for the second sequence should eventually become available")
	require.Eventually(t, func() bool {
		event := natsClient.getLatestEventOfType(*context.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
		if event == nil {
			return false
		}
		return true
	}, 1*time.Minute, 100*time.Millisecond)
}

func Test__main_SequenceQueueTriggerMultiple(t *testing.T) {
	projectName := "my-project-queue2"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"
	numSequences := 10

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()

	sequenceContexts := []apimodels.EventContext{}
	for i := 0; i < numSequences; i++ {
		context := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
		t.Logf("triggered sequence %s with context %s", sequencename, *context.KeptnContext)
		sequenceContexts = append(sequenceContexts, *context)
		<-time.After(100 * time.Millisecond)
	}

	var currentActiveSequence apimodels.SequenceState
	for i := 0; i < numSequences; i++ {
		require.Eventually(t, func() bool {
			states, err := getStates(projectName, &sequenceContexts[i])
			if err != nil {
				return false
			}
			for _, state := range states.States {
				if state.State == apimodels.SequenceStartedState {
					// make sure the sequences are started in the chronologically correct order
					if *sequenceContexts[i].KeptnContext != state.Shkeptncontext {
						return false
					}
					currentActiveSequence = state
					t.Logf("received expected active sequence: %s", state.Shkeptncontext)
					return true
				} else {
					t.Logf("sequence does not have expected state: %s", state.State)
				}
			}
			return false
		}, 30*time.Second, 100*time.Millisecond)

		controlSequence(t, projectName, currentActiveSequence.Shkeptncontext, apimodels.AbortSequence)
	}

	require.Nil(t, err)

}

func Test__main_SequenceStateParallelStages(t *testing.T) {
	projectName := "state-parallel-stages"
	serviceName := "my-service"

	source := "golang-test"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, sequenceStateParallelStagesShipyard)

	defer tearDown()
	require.Nil(t, err)

	keptnContext := natsClient.triggerSequence(projectName, serviceName, "dev", "delivery")

	require.NotNil(t, keptnContext)

	// verify state
	var states *apimodels.SequenceStates

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}
		if states == nil || len(states.States) == 0 {
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	require.Equal(t, int64(1), states.TotalCount)
	require.Len(t, states.States, 1)
	state := states.States[0]

	require.Equal(t, apimodels.SequenceStartedState, state.State)

	require.Len(t, state.Stages, 1)
	stage := state.Stages[0]
	require.Equal(t, "dev", stage.Name)
	require.Equal(t, keptnv2.GetTriggeredEventType("delivery"), stage.LatestEvent.Type)

	// get delivery.triggered event
	deliveryTriggeredEvent := natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("delivery"))
	require.NotNil(t, deliveryTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*deliveryTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})

	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}
		if states == nil || len(states.States) == 0 {
			return false
		}
		return states.States[0].Stages[0].LatestEvent.Type == keptnv2.GetStartedEventType("delivery")
	}, 5*time.Second, 100*time.Millisecond)

	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultPass, Status: keptnv2.StatusSucceeded}, source)
	require.Nil(t, err)

	// now the sequences in staging-1 and staging-2 should have been triggered

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}

		state := states.States[0]
		if state.Project != projectName {
			return false
		}

		if len(state.Stages) != 3 {
			return false
		}

		staging1 := getStageOfState(state, "staging-1")
		staging2 := getStageOfState(state, "staging-2")

		if staging1.LatestEvent.Type != keptnv2.GetTriggeredEventType("delivery") {
			return false
		}
		if staging2.LatestEvent.Type != keptnv2.GetTriggeredEventType("delivery") {
			return false
		}
		return true
	}, 3*time.Second, 100*time.Millisecond)

	var staging1TriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		staging1TriggeredEvent = natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, "staging-1", keptnv2.GetTriggeredEventType("delivery"))
		if staging1TriggeredEvent == nil {
			return false
		}
		return true
	}, 30*time.Second, 5*time.Second)

	cloudEvent = keptnv2.ToCloudEvent(*staging1TriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)

	// send started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}

		state := states.States[0]
		if state.Project != projectName {
			return false
		}

		if len(state.Stages) != 3 {
			return false
		}

		staging1 := getStageOfState(state, "staging-1")
		staging2 := getStageOfState(state, "staging-2")

		if staging1.LatestEvent.Type != keptnv2.GetStartedEventType("delivery") {
			return false
		}
		if staging2.LatestEvent.Type != keptnv2.GetTriggeredEventType("delivery") {
			return false
		}
		return true
	}, 3*time.Second, 100*time.Millisecond)

	// send finished event
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultPass,
	}, source)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}

		state := states.States[0]
		if state.Project != projectName {
			return false
		}

		if len(state.Stages) != 3 {
			return false
		}

		staging1 := getStageOfState(state, "staging-1")
		staging2 := getStageOfState(state, "staging-2")

		if staging1.LatestEvent.Type != keptnv2.GetFinishedEventType("staging-1.delivery") {
			return false
		}
		if staging2.LatestEvent.Type != keptnv2.GetTriggeredEventType("delivery") {
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	staging2TriggeredEvent := natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, "staging-2", keptnv2.GetTriggeredEventType("delivery"))
	require.NotNil(t, staging1TriggeredEvent)

	cloudEvent = keptnv2.ToCloudEvent(*staging2TriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)

	// send started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}

		state := states.States[0]
		if state.Project != projectName {
			return false
		}

		if len(state.Stages) != 3 {
			return false
		}

		staging1 := getStageOfState(state, "staging-1")
		staging2 := getStageOfState(state, "staging-2")

		if staging1.LatestEvent.Type != keptnv2.GetFinishedEventType("staging-1.delivery") {
			return false
		}
		if staging2.LatestEvent.Type != keptnv2.GetStartedEventType("delivery") {
			return false
		}
		return true
	}, 3*time.Second, 100*time.Millisecond)

	// now finish the sequence in staging-2
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultPass,
	}, source)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		states, err = getStates(projectName, keptnContext)
		if err != nil {
			return false
		}

		state := states.States[0]
		if state.Project != projectName {
			return false
		}

		if state.State != apimodels.SequenceFinished {
			return false
		}

		if len(state.Stages) != 3 {
			return false
		}

		staging1 := getStageOfState(state, "staging-1")
		staging2 := getStageOfState(state, "staging-2")

		if staging1.LatestEvent.Type != keptnv2.GetFinishedEventType("staging-1.delivery") {
			return false
		}
		if staging2.LatestEvent.Type != keptnv2.GetFinishedEventType("staging-2.delivery") {
			return false
		}
		return true
	}, 3*time.Second, 100*time.Millisecond)
}

func Test__main_SequenceState_RetrieveMultiple(t *testing.T) {
	projectName := "my-project-retrieve-multiple"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	context1 := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, context1)

	context2 := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, context2)

	verifyContextReturnsStates := func(c *apimodels.EventContext, numResults int) {
		require.Eventually(t, func() bool {
			states, err := getStates(projectName, c)
			if err != nil {
				return false
			}
			if states == nil || len(states.States) != numResults {
				return false
			}
			return true
		}, 5*time.Second, 100*time.Millisecond)
	}

	verifyContextReturnsStates(context1, 1)
	verifyContextReturnsStates(context2, 1)

	// filter by providing two contexts
	combinedContext := fmt.Sprintf("%s,%s", *context1.KeptnContext, *context2.KeptnContext)
	verifyContextReturnsStates(&apimodels.EventContext{KeptnContext: &combinedContext}, 2)
}

func Test__main_SequenceState_SequenceNotFound(t *testing.T) {
	projectName := "state-shipyard-unknown-sequence"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "unknown"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	context := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, context)

	var state apimodels.SequenceState
	require.Eventually(t, func() bool {
		states, err := getStates(projectName, context)
		if err != nil {
			return false
		}
		if states == nil || len(states.States) != 1 {
			return false
		}
		state = states.States[0]
		return true
	}, 5*time.Second, 100*time.Millisecond)

	require.Equal(t, apimodels.SequenceFinished, state.State)

	var finishedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		finishedEvent = natsClient.getLatestEventOfType(*context.KeptnContext, projectName, stageName, keptnv2.GetFinishedEventType("dev.unknown"))
		if finishedEvent == nil {
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	eventData := &keptnv2.EventData{}
	err = keptnv2.Decode(finishedEvent.Data, eventData)
	require.Nil(t, err)

	require.Equal(t, keptnv2.StatusErrored, eventData.Status)
	require.Equal(t, keptnv2.ResultFailed, eventData.Result)
}

func Test__main_SequenceState_InvalidShipyard(t *testing.T) {
	projectName := "state-invalid-shipyard"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "unknown"

	shipyardWithInvalidVersion := strings.Replace(testShipyardFile, "spec.keptn.sh/0.2.2", "0.1.7", 1)
	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, shipyardWithInvalidVersion)

	defer tearDown()
	require.Nil(t, err)

	context := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, context)

	var state apimodels.SequenceState
	require.Eventually(t, func() bool {
		states, err := getStates(projectName, context)
		if err != nil {
			return false
		}
		if states == nil || len(states.States) != 1 {
			return false
		}
		state = states.States[0]
		return true
	}, 5*time.Second, 100*time.Millisecond)

	require.Equal(t, apimodels.SequenceFinished, state.State)

	var finishedEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		finishedEvent = natsClient.getLatestEventOfType(*context.KeptnContext, projectName, stageName, keptnv2.GetFinishedEventType("dev.unknown"))
		if finishedEvent == nil {
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	eventData := &keptnv2.EventData{}
	err = keptnv2.Decode(finishedEvent.Data, eventData)
	require.Nil(t, err)

	require.Equal(t, keptnv2.StatusErrored, eventData.Status)
	require.Equal(t, keptnv2.ResultFailed, eventData.Result)
}

func Test__main_SequenceTimeoutDelayedTask(t *testing.T) {
	projectName := "my-project-delayed-task"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, sequenceTimeoutWithTriggeredAfterShipyard)

	defer tearDown()

	require.Nil(t, err)

	// trigger the task sequence
	t.Log("starting task sequence")
	keptnContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, keptnContext)

	// wait 5s and make verify that the sequence has not been timed out
	<-time.After(5 * time.Second)

	// also, the unknown.triggered event should not have been sent yet
	triggeredEvent := natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("unknown"))
	require.Nil(t, err)
	require.Nil(t, triggeredEvent)

	// after some time, the unknown.triggered event should be available
	require.Eventually(t, func() bool {
		triggeredEvent = natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("unknown"))
		if err != nil {
			return false
		}
		if triggeredEvent == nil {
			return false
		}
		return true
	}, 6*time.Second, 100*time.Millisecond)
}

func Test__main_TriggerAndDeleteProject(t *testing.T) {
	projectName := "my-project-queue-trigger-and-delete"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	numServices := 50
	numSequencesPerService := 1

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	triggerSequence := func(serviceName string, wg *sync.WaitGroup) {
		natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(numServices * numSequencesPerService)

	for i := 0; i < numServices; i++ {
		for j := 0; j < numSequencesPerService; j++ {
			serviceName := fmt.Sprintf("service-%d", i)
			go triggerSequence(serviceName, &wg)
		}
	}
	wg.Wait()

	require.Eventually(t, func() bool {
		states, err := getStates(projectName, nil)
		if err != nil {
			return false
		}

		if states.TotalCount != int64(numServices*numSequencesPerService) {
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	require.Eventually(t, func() bool {
		openTriggeredEvents := getOpenTriggeredEvents(t, projectName, "mytask")

		return openTriggeredEvents.TotalCount > 0
	}, 5*time.Second, 100*time.Millisecond)

	c := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/project/"+projectName, nil)
	require.Nil(t, err)

	resp, err := c.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// recreate the project again
	t.Logf("recreating project %s", projectName)

	createProject(t, c, projectName, testShipyardFile)

	require.Nil(t, err)

	// check if there are any open .triggered events for the project
	openTriggeredEvents := getOpenTriggeredEvents(t, projectName, "")

	require.Empty(t, openTriggeredEvents.Events)
}

func Test__main_SequenceControl_Abort(t *testing.T) {
	projectName := "my-project-abort-sequence"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	source := "shipyard-test"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	// trigger the task sequence
	t.Log("starting task sequence")
	keptnContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, keptnContext)

	// verify state
	verifySequenceEndsUpInState(t, projectName, keptnContext, 3*time.Second, []string{apimodels.SequenceStartedState})

	var taskTriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskTriggeredEvent = natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
		return taskTriggeredEvent != nil
	}, 5*time.Second, 100*time.Millisecond)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	//abort the sequence
	controlSequence(t, projectName, *keptnContext.KeptnContext, apimodels.AbortSequence)

	verifySequenceEndsUpInState(t, projectName, keptnContext, 3*time.Second, []string{apimodels.SequenceAborted})

	openTriggeredEvents := getOpenTriggeredEvents(t, projectName, "mytask")

	require.Empty(t, openTriggeredEvents.Events)
	require.Zero(t, openTriggeredEvents.TotalCount)
}

func Test__main_SequenceControl_AbortQueuedSequence(t *testing.T) {
	projectName := "my-project-abort-queued-sequence"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	source := "shipyard-test"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	// trigger the task sequence
	t.Log("starting task sequence")
	keptnContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, keptnContext)

	// verify state
	verifySequenceEndsUpInState(t, projectName, keptnContext, 3*time.Second, []string{apimodels.SequenceStartedState})

	var taskTriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskTriggeredEvent = natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
		return taskTriggeredEvent != nil
	}, 5*time.Second, 100*time.Millisecond)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	// trigger a second sequence which should be put in the queue
	secondContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)

	// verify state
	verifySequenceEndsUpInState(t, projectName, secondContext, 5*time.Second, []string{apimodels.SequenceWaitingState})

	// abort the queued sequence
	t.Log("aborting sequence")
	controlSequence(t, projectName, *secondContext.KeptnContext, apimodels.AbortSequence)

	verifySequenceEndsUpInState(t, projectName, secondContext, 5*time.Second, []string{apimodels.SequenceAborted})
}

func Test__main_SequenceControl_PauseSequence(t *testing.T) {
	projectName := "my-project-abort-paused-sequence"
	stageName := "dev"
	serviceName := "my-service"
	sequencename := "delivery"

	source := "shipyard-test"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	// trigger the task sequence
	t.Log("starting task sequence")
	keptnContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, keptnContext)

	// verify state
	verifySequenceEndsUpInState(t, projectName, keptnContext, 3*time.Second, []string{apimodels.SequenceStartedState})

	var taskTriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskTriggeredEvent = natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
		return taskTriggeredEvent != nil
	}, 1*time.Second, 100*time.Millisecond)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	// pause the sequence
	controlSequence(t, projectName, *keptnContext.KeptnContext, apimodels.PauseSequence)

	// verify state
	verifySequenceEndsUpInState(t, projectName, keptnContext, 5*time.Second, []string{apimodels.SequencePaused})

	// trigger a second sequence which should take over and be started
	secondContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)

	// verify state
	verifySequenceEndsUpInState(t, projectName, secondContext, 5*time.Second, []string{apimodels.SequenceStartedState})
}

func Test__main_SequenceControl_AbortPausedSequenceTaskPartiallyFinished(t *testing.T) {
	projectName := "sequence-abort4"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "delivery"
	source1 := "golang-test-1"
	source2 := "golang-test-2"

	natsClient, tearDown, err := setupTestProject(t, projectName, serviceName, testShipyardFile)

	defer tearDown()
	require.Nil(t, err)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)

	require.NotNil(t, keptnContext)

	// verify state
	verifySequenceEndsUpInState(t, projectName, keptnContext, 5*time.Second, []string{apimodels.SequenceStartedState})

	var taskTriggeredEvent *apimodels.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskTriggeredEvent = natsClient.getLatestEventOfType(*keptnContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
		return taskTriggeredEvent != nil
	}, 5*time.Second, 100*time.Millisecond)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: natsClient})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending two task started events")
	_, err = keptn.SendTaskStartedEvent(nil, source1)
	require.Nil(t, err)
	_, err = keptn.SendTaskStartedEvent(nil, source2)
	require.Nil(t, err)

	// simulate the duration of a task execution
	<-time.After(100 * time.Millisecond)

	t.Logf("send one finished event with result 'fail'")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultFailed, Status: keptnv2.StatusSucceeded}, source1)
	require.Nil(t, err)

	// now trigger another sequence and make sure it is started eventually
	// trigger a second sequence which should be put in the queue
	secondContext := natsClient.triggerSequence(projectName, serviceName, stageName, sequencename)
	require.NotNil(t, secondContext)

	// verify that the second sequence gets the triggered status
	verifySequenceEndsUpInState(t, projectName, secondContext, 5*time.Second, []string{apimodels.SequenceWaitingState})

	// pause the first sequence
	t.Log("pausing sequence")
	controlSequence(t, projectName, *keptnContext.KeptnContext, apimodels.PauseSequence)

	verifySequenceEndsUpInState(t, projectName, keptnContext, 5*time.Second, []string{apimodels.SequencePaused})

	// now abort the first sequence
	t.Log("aborting first sequence")
	controlSequence(t, projectName, *keptnContext.KeptnContext, apimodels.AbortSequence)

	verifySequenceEndsUpInState(t, projectName, keptnContext, 5*time.Second, []string{apimodels.SequenceAborted})

	// now that the first sequence is aborted, the other sequence should eventually be started
	verifySequenceEndsUpInState(t, projectName, secondContext, 5*time.Second, []string{apimodels.SequenceStartedState})

	// also make sure that the triggered event for the first task has been sent
	require.Eventually(t, func() bool {
		taskTriggeredEvent := natsClient.getLatestEventOfType(*secondContext.KeptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
		if taskTriggeredEvent == nil {
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

}

func getOpenTriggeredEvents(t *testing.T, projectName string, taskName string) *apimodels.Events {
	c := http.Client{}
	openTriggeredEvents := &apimodels.Events{}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/event/triggered/"+keptnv2.GetTriggeredEventType(taskName)+"?project="+projectName, nil)
	require.Nil(t, err)

	resp, err := c.Do(req)
	require.Nil(t, err)

	respBytes, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	err = json.Unmarshal(respBytes, openTriggeredEvents)
	require.Nil(t, err)
	return openTriggeredEvents
}

func setupTestProject(t *testing.T, projectName, serviceName, shipyardContent string) (*testNatsClient, func(), error) {
	setupFakeConfigurationService(shipyardContent)

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
	natsClient, err := newTestNatsClient(natsURL, t)
	require.Nil(t, err)

	c := http.Client{
		Timeout: 2 * time.Second,
	}

	createProject(t, c, projectName, shipyardContent)

	service := models.CreateServiceParams{
		ServiceName: &serviceName,
	}

	marshal, err := json.Marshal(service)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/project/"+projectName+"/service", bytes.NewBuffer(marshal))
	require.Nil(t, err)

	resp, err := c.Do(req)

	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	tearDown := func() {
		natsClient.Close()
	}
	return natsClient, tearDown, err
}

func createProject(t *testing.T, c http.Client, projectName string, shipyardContent string) {
	encodedShipyardContent := base64.StdEncoding.EncodeToString([]byte(shipyardContent))
	createProjectObj := models.CreateProjectParams{
		Name:     &projectName,
		Shipyard: &encodedShipyardContent,
	}

	marshal, err := json.Marshal(createProjectObj)

	require.Nil(t, err)

	require.Eventually(t, func() bool {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/v1/project", bytes.NewBuffer(marshal))
		if err != nil {
			return false
		}

		resp, err := c.Do(req)

		if err != nil {
			return false
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return false
		}
		return true
	}, 10*time.Second, 100*time.Millisecond)
}

func setupFakeConfigurationService(shipyardContent string) {
	configurationService.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if strings.Contains(r.RequestURI, "/shipyard.yaml") {
			w.WriteHeader(200)
			encodedShipyard := base64.StdEncoding.EncodeToString([]byte(shipyardContent))
			res := apimodels.Resource{
				ResourceContent: encodedShipyard,
			}

			marshal, _ := json.Marshal(res)
			w.Write(marshal)

			return
		} else if strings.Contains(r.RequestURI, "/metadata.yaml") {
			res := apimodels.Resource{
				Metadata: &apimodels.Version{
					Version: "my-commit-id",
				},
			}

			marshal, _ := json.Marshal(res)
			w.Write(marshal)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
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
	for _, ev := range n.receivedEvents {
		if ev.Shkeptncontext == keptnContext && *ev.Type == eventType {
			ed := &keptnv2.EventData{}
			err := keptnv2.Decode(ev.Data, ed)
			require.Nil(n.t, err)
			if ed.Project == projectName && ed.Stage == stage {
				return &ev
			}
		}
	}
	return nil
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
