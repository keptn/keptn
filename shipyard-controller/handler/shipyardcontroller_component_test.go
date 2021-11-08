package handler

import (
	"context"
	"errors"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	fakehooks "github.com/keptn/keptn/shipyard-controller/handler/sequencehooks/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

const testShipyardFileWithInvalidVersion = `apiVersion: 0
kind: Shipyard
metadata:
  name: test-shipyard`

const testShipyardFile = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
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

const testShipyardFileWithDuplicateTasks = `apiVersion: spec.keptn.sh/0.2.2
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
      - name: deployment
      - name: evaluation`

const mongoDBVersion = "4.4.9"

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

	return func() { mongoServer.Stop() }
}

//Scenario 1: Complete task sequence execution + triggering of next task sequence. Events are received in order
func Test_shipyardController_Scenario1(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	done := false

	// STEP 1
	// send dev.artifact-delivery.triggered event
	sequenceTriggeredEvent := getArtifactDeliveryTriggeredEvent("dev")
	err := sc.HandleIncomingEvent(sequenceTriggeredEvent, true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event dispatcher -> should contain deployment.triggered event with properties: [deployment]
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())

	deploymentEvent := &keptnv2.DeploymentTriggeredEventData{}
	err = verifyEvent.Event.DataAs(deploymentEvent)
	require.Nil(t, err)
	require.Equal(t, 1, len(deploymentEvent.Deployment.DeploymentURIsPublic))
	require.Equal(t, "direct", deploymentEvent.Deployment.DeploymentStrategy)
	require.Equal(t, "carts", deploymentEvent.ConfigurationChange.Values["image"])

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.TestTaskName,
		"",
	)
	if done {
		return
	}
	require.Equal(t, 2, len(mockDispatcher.AddCalls()))
	verifyEvent = mockDispatcher.AddCalls()[1].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), verifyEvent.Event.Type())

	taskEvent := &keptnv2.TestTriggeredEventData{}
	err = verifyEvent.Event.DataAs(taskEvent)
	require.Nil(t, err)
	require.Equal(t, 3, len(taskEvent.Deployment.DeploymentURIsPublic))
	require.Equal(t, 2, len(taskEvent.Deployment.DeploymentURIsLocal))

	// also check if the payload of the .triggered event that started the sequence is present
	deploymentEvent = &keptnv2.DeploymentTriggeredEventData{}
	err = verifyEvent.Event.DataAs(deploymentEvent)
	require.Equal(t, "carts", deploymentEvent.ConfigurationChange.Values["image"])

	// STEP 4
	// send test.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.TestTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 5
	// send test.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getTestTaskFinishedEvent("dev", triggeredID),
		keptnv2.TestTaskName,
		keptnv2.EvaluationTaskName,
		"",
	)
	if done {
		return
	}
	require.Equal(t, 3, len(mockDispatcher.AddCalls()))
	verifyEvent = mockDispatcher.AddCalls()[2].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), verifyEvent.Event.Type())

	evaluationEvent := &keptnv2.EvaluationTriggeredEventData{}
	err = verifyEvent.Event.DataAs(evaluationEvent)
	require.Nil(t, err)
	require.Equal(t, 1, len(evaluationEvent.Deployment.DeploymentNames))
	require.Equal(t, "start", evaluationEvent.Test.Start)
	require.Equal(t, "end", evaluationEvent.Test.End)

	// STEP 6
	// send evaluation.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.EvaluationTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 7
	// send evaluation.finished event -> result = warning should not abort the task sequence
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getEvaluationTaskFinishedEvent("dev", triggeredID, keptnv2.ResultWarning), keptnv2.EvaluationTaskName, keptnv2.ReleaseTaskName, "")
	if done {
		return
	}
	require.Equal(t, 4, len(mockDispatcher.AddCalls()))
	verifyEvent = mockDispatcher.AddCalls()[3].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.ReleaseTaskName), verifyEvent.Event.Type())

	// STEP 8
	// send release.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.ReleaseTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 9
	// send release.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getReleaseTaskFinishedEvent("dev", triggeredID), keptnv2.ReleaseTaskName, keptnv2.DeploymentTaskName, "hardening")
	if done {
		return
	}

	require.Equal(t, 7, len(mockDispatcher.AddCalls()))

	// verify dev.artifact-delivery.finished event
	sequenceFinishedEvent := mockDispatcher.AddCalls()[4].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), sequenceFinishedEvent.Event.Type())
	require.Equal(t, sequenceTriggeredEvent.ID, sequenceFinishedEvent.Event.Extensions()["triggeredid"])

	// verify hardening.artifact-delivery.triggered event
	nextSequenceTriggeredEvent := mockDispatcher.AddCalls()[5].Event
	require.Equal(t, keptnv2.GetTriggeredEventType("hardening.artifact-delivery"), nextSequenceTriggeredEvent.Event.Type())

	sequenceTriggeredDataMap := map[string]interface{}{}
	err = nextSequenceTriggeredEvent.Event.DataAs(&sequenceTriggeredDataMap)
	require.Nil(t, err)
	require.NotNil(t, sequenceTriggeredDataMap["configurationChange"])
	require.NotNil(t, sequenceTriggeredDataMap["deployment"])

	// verify deployment.triggered event for hardening stage
	verifyEvent = mockDispatcher.AddCalls()[6].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())
	deploymentEvent = &keptnv2.DeploymentTriggeredEventData{}
	err = verifyEvent.Event.DataAs(deploymentEvent)
	require.Nil(t, err)
	require.Equal(t, "hardening", deploymentEvent.Stage)
	require.Equal(t, "carts", deploymentEvent.ConfigurationChange.Values["image"])

	// verify that data from .finished events of the previous stage are included
	deploymentTriggeredDataMap := map[string]interface{}{}
	err = verifyEvent.Event.DataAs(&deploymentTriggeredDataMap)
	require.Nil(t, err)
	require.NotNil(t, deploymentTriggeredDataMap["test"])

	finishedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Stage: common.Stringp("dev"),
	}, common.FinishedEvent)

	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), "dev")
	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.TestTaskName), "dev")
	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), "dev")
	fake.ShouldNotContainEvent(t, finishedEvents, keptnv2.GetFinishedEventType(keptnv2.ReleaseTaskName), "dev")

	// STEP 9.1
	// send deployment.started event 1 with ID 1
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "hardening", "test-source-1")
	if done {
		return
	}

	// STEP 9.2
	// send deployment.started event 2 with ID 2
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "hardening", "test-source-2")
	if done {
		return
	}

	// STEP 10.1
	// send deployment.finished event 1 with ID 1
	done = sendAndVerifyPartialFinishedEvent(t, sc, getDeploymentFinishedEvent("hardening", triggeredID, "test-source-1", keptnv2.ResultPass), keptnv2.DeploymentTaskName, keptnv2.ReleaseTaskName, "")
	if done {
		return
	}
	// number of calls for dispatcher should not have increased before both finished events are received
	require.Equal(t, 7, len(mockDispatcher.AddCalls()))

	// STEP 10.2
	// send deployment.finished event 1 with ID 1
	triggeredID, done = sendAndVerifyFinishedEvent(t, sc, getDeploymentFinishedEvent("hardening", triggeredID, "test-source-2", keptnv2.ResultPass), keptnv2.DeploymentTaskName, keptnv2.TestTaskName, "")
	if done {
		return
	}
	require.Equal(t, 8, len(mockDispatcher.AddCalls()))
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), mockDispatcher.AddCalls()[7].Event.Event.Type())
}

//Scenario 2: Partial task sequence execution + triggering of next task sequence. Events are received out of order
func Test_shipyardController_Scenario2(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	done := false

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event dispatcher -> should contain deployment.triggered event with properties: [deployment]
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	go func() {
		time.After(2 * time.Second)
		_ = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	}()

	// STEP 3
	// send deployment.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.TestTaskName,
		"",
	)
	if done {
		return
	}
	require.Equal(t, 2, len(mockDispatcher.AddCalls()))
	verifyEvent = mockDispatcher.AddCalls()[1].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), verifyEvent.Event.Type())

	taskEvent := &keptnv2.TestTriggeredEventData{}
	err = verifyEvent.Event.DataAs(taskEvent)
	require.Nil(t, err)
	require.Equal(t, 3, len(taskEvent.Deployment.DeploymentURIsPublic))
	require.Equal(t, 2, len(taskEvent.Deployment.DeploymentURIsLocal))

}

//Scenario 3: Received .finished event with status "errored" should abort task sequence and send .finished event with status "errored"
func Test_shipyardController_Scenario3(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")

	done := false

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event dispatcher -> should contain deployment.triggered event with properties: [deployment]
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	go func() {
		time.After(2 * time.Second)
		_ = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	}()

	// STEP 3
	// send deployment.finished event
	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getErroredDeploymentFinishedEvent("dev", triggeredID, "test-source"),
		keptnv2.DeploymentTaskName,
		"",
	)
	if done {
		return
	}

	// check for dev.artifact-delivery.finished event
	require.Equal(t, 4, len(mockDispatcher.AddCalls()))
	taskSequenceCompletionEvent := mockDispatcher.AddCalls()[1].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), taskSequenceCompletionEvent.Event.Type())

	eventData := &keptnv2.EventData{}
	err = taskSequenceCompletionEvent.Event.DataAs(eventData)
	require.Nil(t, err)
	require.Equal(t, keptnv2.StatusErrored, eventData.Status)
	require.Equal(t, keptnv2.ResultFailed, eventData.Result)

	require.Equal(t, keptnv2.GetTriggeredEventType("dev.rollback"), mockDispatcher.AddCalls()[2].Event.Event.Type())
	require.Equal(t, keptnv2.GetTriggeredEventType("rollback"), mockDispatcher.AddCalls()[3].Event.Event.Type())

}

//Scenario 4: Received .finished event with result "fail" - stop task sequence
func Test_shipyardController_Scenario4(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")

	done := false

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event dispatcher -> should contain deployment.triggered event with properties: [deployment]
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.TestTaskName,
		"",
	)
	if done {
		return
	}
	require.Equal(t, 2, len(mockDispatcher.AddCalls()))
	verifyEvent = mockDispatcher.AddCalls()[1].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), verifyEvent.Event.Type())

	// STEP 4
	// send test.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.TestTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 5
	// send test.finished event
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getTestTaskFinishedEvent("dev", triggeredID),
		keptnv2.TestTaskName,
		keptnv2.EvaluationTaskName,
		"",
	)
	if done {
		return
	}
	require.Equal(t, 3, len(mockDispatcher.AddCalls()))
	verifyEvent = mockDispatcher.AddCalls()[2].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), verifyEvent.Event.Type())
	// STEP 6
	// send evaluation.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.EvaluationTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 7
	// send evaluation.finished event with result=fail

	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getEvaluationTaskFinishedEvent("dev", triggeredID, keptnv2.ResultFailed),
		keptnv2.EvaluationTaskName,
		"",
	)
	if done {
		return
	}

	// check for dev.artifact-delivery.finished event
	require.Equal(t, 6, len(mockDispatcher.AddCalls()))
	taskSequenceCompletionEvent := mockDispatcher.AddCalls()[3].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), taskSequenceCompletionEvent.Event.Type())

	eventData := &keptnv2.EventData{}
	err = taskSequenceCompletionEvent.Event.DataAs(eventData)
	require.Nil(t, err)
	require.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	require.Equal(t, keptnv2.ResultFailed, eventData.Result)

	require.Equal(t, keptnv2.GetTriggeredEventType("dev.rollback"), mockDispatcher.AddCalls()[4].Event.Event.Type())
	require.Equal(t, keptnv2.GetTriggeredEventType("rollback"), mockDispatcher.AddCalls()[5].Event.Event.Type())
}

//Scenario 4a: Handling multiple finished events, one has result==failed, ==> task sequence is stopped
func Test_shipyardController_Scenario4a(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 1 with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")

	done := false

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event dispatcher -> should contain deployment.triggered event with properties: [deployment]
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started events
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "another-test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	err = sendFinishedEvent(sc, getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultFailed))
	require.Nil(t, err)

	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "another-test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		"",
	)

	// check for dev.artifact-delivery.finished event
	require.Equal(t, 4, len(mockDispatcher.AddCalls()))
	taskSequenceCompletionEvent := mockDispatcher.AddCalls()[1].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), taskSequenceCompletionEvent.Event.Type())

	eventData := &keptnv2.EventData{}
	err = taskSequenceCompletionEvent.Event.DataAs(eventData)
	require.Nil(t, err)
	require.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	require.Equal(t, keptnv2.ResultFailed, eventData.Result)

	require.Equal(t, keptnv2.GetTriggeredEventType("dev.rollback"), mockDispatcher.AddCalls()[2].Event.Event.Type())
	require.Equal(t, keptnv2.GetTriggeredEventType("rollback"), mockDispatcher.AddCalls()[3].Event.Event.Type())
}

//Scenario 4b: Received .finished event with result "fail" - stop task sequence and trigger next sequence based on result filter
func Test_shipyardController_TriggerOnFail(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")

	done := false

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	// check event dispatcher -> should contain deployment.triggered event with properties: [deployment]
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), verifyEvent.Event.Type())

	// check triggeredEvent Collection -> should contain deployment.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:   common.Stringp("dev"),
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done = fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), "", nil)
	if done {
		return
	}
	triggeredID := triggeredEvents[0].ID

	// STEP 2
	// send deployment.started event
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	done = sendFinishedEventAndVerifyTaskSequenceCompletion(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultFailed),
		keptnv2.DeploymentTaskName,
		"",
	)
	if done {
		return
	}

	// check for dev.artifact-delivery.finished event
	require.Equal(t, 4, len(mockDispatcher.AddCalls()))
	taskSequenceCompletionEvent := mockDispatcher.AddCalls()[1].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), taskSequenceCompletionEvent.Event.Type())

	eventData := &keptnv2.EventData{}
	err = taskSequenceCompletionEvent.Event.DataAs(eventData)
	require.Nil(t, err)
	require.Equal(t, keptnv2.StatusSucceeded, eventData.Status)
	require.Equal(t, keptnv2.ResultFailed, eventData.Result)

	// check for dev.rollback.triggered
	rollbackTriggeredEvent := mockDispatcher.AddCalls()[2].Event
	require.Equal(t, keptnv2.GetTriggeredEventType("dev.rollback"), rollbackTriggeredEvent.Event.Type())

	// check for rollback.triggered
	rollbackTaskTriggeredEvent := mockDispatcher.AddCalls()[3].Event
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.RollbackTaskName), rollbackTaskTriggeredEvent.Event.Type())

	for _, addCall := range mockDispatcher.AddCalls() {
		// hardening.artifact-delivery should not be triggered
		require.NotEqual(t, keptnv2.GetTriggeredEventType("hardening.artifact-delivery"), addCall.Event.Event.Type())
	}
}

//Scenario 5: Received .triggered event for project with invalid shipyard version -> send .finished event with result = fail
func Test_shipyardController_Scenario5(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 5 with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController(testShipyardFileWithInvalidVersion)

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), verifyEvent.Event.Type())

}

func Test_shipyardController_DuplicateTask(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller Scenario 6 (duplicate tasks) with shipyard file %s", testShipyardFileWithDuplicateTasks)
	sc := getTestShipyardController(testShipyardFileWithDuplicateTasks)

	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	triggeredEvent := mockDispatcher.AddCalls()[0].Event
	triggeredKeptnEvent, err := keptnv2.ToKeptnEvent(triggeredEvent.Event)
	require.Equal(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), *triggeredKeptnEvent.Type)

	// STEP 2
	// send deployment.started event
	done := sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredKeptnEvent.ID, "dev", "test-source")
	if done {
		return
	}

	// STEP 3
	// send deployment.finished event
	triggeredID, done := sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredKeptnEvent.ID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.DeploymentTaskName,
		"",
	)
	if done {
		return
	}

	// STEP 4
	// send deployment.started event (for the second deployment task)
	done = sendAndVerifyStartedEvent(t, sc, keptnv2.DeploymentTaskName, triggeredID, "dev", "test-source")
	if done {
		return
	}

	// STEP 5
	// send deployment.finished event for the second deployment task -> now we want an evaluation.triggered event as the next task
	triggeredID, done = sendAndVerifyFinishedEvent(
		t,
		sc,
		getDeploymentFinishedEvent("dev", triggeredID, "test-source", keptnv2.ResultPass),
		keptnv2.DeploymentTaskName,
		keptnv2.EvaluationTaskName,
		"",
	)
}

func Test_shipyardController_TimeoutSequence(t *testing.T) {
	defer setupLocalMongoDB()()

	sc := getTestShipyardController("")

	fakeTimeoutHook := &fakehooks.ISequenceTimeoutHookMock{OnSequenceTimeoutFunc: func(event models.Event) {}}
	sc.AddSequenceTimeoutHook(fakeTimeoutHook)

	// insert the test data
	_ = sc.eventRepo.InsertEvent("my-project", models.Event{
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ID:             "my-sequence-triggered-id",
		Shkeptncontext: "my-keptn-context-id",
		Type:           common.Stringp(keptnv2.GetTriggeredEventType("my-stage.delivery")),
	}, common.TriggeredEvent)

	_ = sc.eventRepo.InsertEvent("my-project", models.Event{
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ID:             "my-deployment-triggered-id",
		Shkeptncontext: "my-keptn-context-id",
		Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
	}, common.TriggeredEvent)

	sc.taskSequenceRepo.CreateTaskSequenceMapping("my-project", models.TaskSequenceEvent{
		TaskSequenceName: "delivery",
		TriggeredEventID: "my-task-triggered-id",
		Task:             models.Task{},
		Stage:            "my-stage",
		KeptnContext:     "my-keptn-context-id",
	})

	// invoke the CancelSequence function
	err := sc.timeoutSequence(models.SequenceTimeout{
		KeptnContext: "my-keptn-context-id",
		LastEvent: models.Event{
			Data: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			Type:           common.Stringp(keptnv2.GetTriggeredEventType("my-task")),
			ID:             "my-task-triggered-id",
			Shkeptncontext: "my-keptn-context-id",
		},
	})

	require.Nil(t, err)
	require.Len(t, fakeTimeoutHook.OnSequenceTimeoutCalls(), 1)
}

func Test_shipyardController_CancelSequence(t *testing.T) {
	defer setupLocalMongoDB()()
	sc := getTestShipyardController("")

	fakeSequenceFinishedHook := &fakehooks.ISequenceFinishedHookMock{OnSequenceFinishedFunc: func(event models.Event) {}}
	sc.AddSequenceFinishedHook(fakeSequenceFinishedHook)

	// insert the test data
	_ = sc.eventRepo.InsertEvent("my-project", models.Event{
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ID:             "my-sequence-triggered-id",
		Shkeptncontext: "my-keptn-context-id",
		Type:           common.Stringp(keptnv2.GetTriggeredEventType("my-stage.delivery")),
	}, common.TriggeredEvent)

	_ = sc.eventRepo.InsertEvent("my-project", models.Event{
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ID:             "my-deployment-triggered-id",
		Shkeptncontext: "my-keptn-context-id",
		Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
	}, common.TriggeredEvent)

	taskSequenceMapping := models.TaskSequenceEvent{
		TaskSequenceName: "delivery",
		TriggeredEventID: "my-deployment-triggered-id",
		Task:             models.Task{},
		Stage:            "my-stage",
		KeptnContext:     "my-keptn-context-id",
	}
	sc.taskSequenceRepo.CreateTaskSequenceMapping("my-project", taskSequenceMapping)

	// invoke the CancelSequence function
	err := sc.cancelSequence(models.SequenceControl{
		KeptnContext: "my-keptn-context-id",
		Project:      "my-project",
		Stage:        "my-stage",
	})

	require.Nil(t, err)
	require.Len(t, fakeSequenceFinishedHook.OnSequenceFinishedCalls(), 1)
}

func Test_shipyardController_CancelQueuedSequence(t *testing.T) {
	defer setupLocalMongoDB()()

	sc := getTestShipyardController("")
	sequenceDispatcherMock := &fake.ISequenceDispatcherMock{}
	sequenceDispatcherMock.RemoveFunc = func(eventScope models.EventScope) error {
		return nil
	}

	sc.sequenceDispatcher = sequenceDispatcherMock

	fakeSequenceFinishedHook := &fakehooks.ISequenceFinishedHookMock{OnSequenceFinishedFunc: func(event models.Event) {}}
	sc.AddSequenceFinishedHook(fakeSequenceFinishedHook)

	// insert the test data
	_ = sc.eventRepo.InsertEvent("my-project", models.Event{
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ID:             "my-sequence-triggered-id",
		Shkeptncontext: "my-keptn-context-id",
		Type:           common.Stringp(keptnv2.GetTriggeredEventType("my-stage.delivery")),
	}, common.TriggeredEvent)

	// invoke the CancelSequence function
	err := sc.cancelSequence(models.SequenceControl{
		KeptnContext: "my-keptn-context-id",
		Project:      "my-project",
		Stage:        "my-stage",
	})

	require.Nil(t, err)
	require.Len(t, fakeSequenceFinishedHook.OnSequenceFinishedCalls(), 1)
}

func Test_shipyardController_CancelQueuedSequence_RemoveFromQueueFails(t *testing.T) {
	defer setupLocalMongoDB()()

	sc := getTestShipyardController("")
	sequenceDispatcherMock := &fake.ISequenceDispatcherMock{}
	sequenceDispatcherMock.RemoveFunc = func(eventScope models.EventScope) error {
		return errors.New("oops")
	}

	sc.sequenceDispatcher = sequenceDispatcherMock

	fakeSequenceFinishedHook := &fakehooks.ISequenceFinishedHookMock{OnSequenceFinishedFunc: func(event models.Event) {}}
	sc.AddSequenceFinishedHook(fakeSequenceFinishedHook)

	// insert the test data
	_ = sc.eventRepo.InsertEvent("my-project", models.Event{
		Data: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ID:             "my-sequence-triggered-id",
		Shkeptncontext: "my-keptn-context-id",
		Type:           common.Stringp(keptnv2.GetTriggeredEventType("my-stage.delivery")),
	}, common.TriggeredEvent)

	// invoke the CancelSequence function
	err := sc.cancelSequence(models.SequenceControl{
		KeptnContext: "my-keptn-context-id",
		Project:      "my-project",
		Stage:        "my-stage",
	})

	require.Nil(t, err)
	require.Len(t, fakeSequenceFinishedHook.OnSequenceFinishedCalls(), 1)
}

func Test_SequenceForUnavailableStage(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFile)
	sc := getTestShipyardController("")
	sc.sequenceDispatcher = &fake.ISequenceDispatcherMock{
		AddFunc: func(queueItem models.QueueItem) error {
			return nil
		},
	}

	mockEventDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)
	mockSequenceDispatcher := sc.sequenceDispatcher.(*fake.ISequenceDispatcherMock)

	// STEP 1
	// send unknown.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("unknown"), true)

	require.Nil(t, err)
	require.Len(t, mockEventDispatcher.AddCalls(), 1)
	require.Equal(t, keptnv2.GetFinishedEventType("unknown.artifact-delivery"), mockEventDispatcher.AddCalls()[0].Event.Event.Type())
	require.Empty(t, mockSequenceDispatcher.AddCalls())
}

// Updating event of service fails -> event handling should still happen
func Test_UpdateEventOfServiceFailsFails(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController(testShipyardFileWithInvalidVersion)
	mockDispatcher := sc.eventDispatcher.(*fake.IEventDispatcherMock)

	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(getArtifactDeliveryTriggeredEvent("dev"), true)
	if err != nil {
		t.Errorf("STEP 1 failed: HandleIncomingEvent(dev.artifact-delivery.triggered) returned %v", err)
		return
	}

	require.Eventually(t, func() bool {
		return len(mockDispatcher.AddCalls()) == 1
	}, 10*time.Second, 1*time.Second)
	require.Equal(t, 1, len(mockDispatcher.AddCalls()))
	verifyEvent := mockDispatcher.AddCalls()[0].Event
	require.Equal(t, keptnv2.GetFinishedEventType("dev.artifact-delivery"), verifyEvent.Event.Type())
}

// Scenario 5: Received .triggered event for project with invalid shipyard version -> send .finished event with result = fail
func Test_UpdateServiceShouldNotBeCalledForEmptyService(t *testing.T) {
	defer setupLocalMongoDB()()

	t.Logf("Executing Shipyard Controller with shipyard file %s", testShipyardFileWithInvalidVersion)
	sc := getTestShipyardController("")

	event := getArtifactDeliveryTriggeredEvent("dev")

	event.Data = keptnv2.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "",
	}
	// STEP 1
	// send dev.artifact-delivery.triggered event
	err := sc.HandleIncomingEvent(event, true)

	assert.NotNil(t, err)
}

func getTestShipyardController(shipyardContent string) *shipyardController {
	if shipyardContent == "" {
		shipyardContent = testShipyardFile
	}

	em := &shipyardController{
		projectMvRepo:    db.NewProjectMVRepo(db.NewMongoDBProjectsRepo(db.GetMongoDBConnectionInstance()), db.NewMongoDBEventsRepo(db.GetMongoDBConnectionInstance())),
		eventRepo:        db.NewMongoDBEventsRepo(db.GetMongoDBConnectionInstance()),
		taskSequenceRepo: db.NewTaskSequenceMongoDBRepo(db.GetMongoDBConnectionInstance()),
		eventDispatcher: &fake.IEventDispatcherMock{
			AddFunc: func(event models.DispatcherEvent) error {
				return nil
			},
			RunFunc: func(ctx context.Context) {

			},
		},
		shipyardRetriever: &fake.IShipyardRetrieverMock{
			GetShipyardFunc: func(projectName string) (*keptnv2.Shipyard, error) {
				return common.UnmarshalShipyard(shipyardContent)
			},
			GetCachedShipyardFunc: func(projectName string) (*keptnv2.Shipyard, error) {
				return common.UnmarshalShipyard(shipyardContent)
			},
		},
	}

	em.eventDispatcher.(*fake.IEventDispatcherMock).AddFunc = func(event models.DispatcherEvent) error {
		ev := &models.Event{}
		err := keptnv2.Decode(&event.Event, ev)
		if err != nil {
			return err
		}
		_ = em.HandleIncomingEvent(*ev, true)
		return nil
	}
	return em
}

func filterEvents(eventsCollection []models.Event, filter common.EventFilter) ([]models.Event, error) {
	result := []models.Event{}

	for _, event := range eventsCollection {
		scope, _ := models.NewEventScope(event)
		if filter.Type != "" && *event.Type != filter.Type {
			continue
		}
		if filter.Stage != nil && *filter.Stage != scope.Stage {
			continue
		}

		if filter.Service != nil && *filter.Service != scope.Service {
			continue
		}
		if filter.TriggeredID != nil && *filter.TriggeredID != event.Triggeredid {
			continue
		}
		if filter.KeptnContext != nil && *filter.KeptnContext != event.Shkeptncontext {
			continue
		}
		result = append(result, event)
	}
	return result, nil
}

func getDeploymentFinishedEvent(stage string, triggeredID string, source string, result keptnv2.ResultType) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  result,
				Message: "i am a message",
			},
			Deployment: keptnv2.DeploymentFinishedData{
				DeploymentURIsLocal:  []string{"uri-1", "uri-2"},
				DeploymentURIsPublic: []string{"public-uri-1", "public-uri-2"},
				DeploymentNames:      []string{"deployment-1"},
				GitCommit:            "commit-1",
			},
		},
		Extensions:     nil,
		ID:             "deployment-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.deployment.finished"),
	}
}

func getErroredDeploymentFinishedEvent(stage string, triggeredID string, source string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultFailed,
			},
			Deployment: keptnv2.DeploymentFinishedData{
				DeploymentURIsLocal:  []string{"uri-1", "uri-2"},
				DeploymentURIsPublic: []string{"public-uri-1", "public-uri-2"},
				DeploymentNames:      []string{"deployment-1"},
				GitCommit:            "commit-1",
			},
		},
		Extensions:     nil,
		ID:             "deployment-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.deployment.finished"),
	}
}

func getTestTaskFinishedEvent(stage string, triggeredID string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.TestFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
			},
			Test: keptnv2.TestFinishedDetails{
				Start:     "start",
				End:       "end",
				GitCommit: "commit-id",
			},
		},
		Extensions:     nil,
		ID:             "test-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.test.finished"),
	}
}

func getEvaluationTaskFinishedEvent(stage string, triggeredID string, result keptnv2.ResultType) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.EvaluationFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  result,
			},
			Evaluation: keptnv2.EvaluationDetails{
				Result: string(result),
			},
		},
		Extensions:     nil,
		ID:             "evaluation-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.evaluation.finished"),
	}
}

func getReleaseTaskFinishedEvent(stage string, triggeredID string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.ReleaseFinishedEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
			},
		},
		Extensions:     nil,
		ID:             "release-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp("sh.keptn.event.release.finished"),
	}
}

func sendFinishedEvent(sc *shipyardController, finishedEvent models.Event) error {
	return sc.HandleIncomingEvent(finishedEvent, true)
}

func sendAndVerifyFinishedEvent(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextEventType string, nextStage string) (string, bool) {
	err := sc.HandleIncomingEvent(finishedEvent, true)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return "", true
	}

	scope, _ := models.NewEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should not contain <eventType>.triggered event anymore
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: common.Stringp("carts"),
		ID:      &scope.TriggeredID,
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	require.NotContains(t, triggeredEvents, models.Event{
		ID: scope.TriggeredID,
	})

	// check triggeredEvent collection -> should contain <nextEventType>.triggered event
	triggeredEvents, _ = sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(nextEventType),
		Stage:   &nextStage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)

	triggeredID := triggeredEvents[0].ID
	done := fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage, nil)
	if done {
		return "", true
	}

	// check startedEvent collection -> should not contain <eventType>.started event anymore
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(finishedEvent.Triggeredid),
	}, common.StartedEvent)
	done = fake.ShouldNotContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage)
	if done {
		return "", true
	}

	return triggeredID, false
}

func sendFinishedEventAndVerifyTaskSequenceCompletion(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextStage string) bool {
	err := sc.HandleIncomingEvent(finishedEvent, true)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return true
	}

	scope, _ := models.NewEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should not contain <eventType>.triggered event anymore
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done := fake.ShouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage)
	if done {
		return true
	}

	// check startedEvent collection -> should not contain <eventType>.started event anymore
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(finishedEvent.Triggeredid),
	}, common.StartedEvent)
	return fake.ShouldNotContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage)
}

func sendAndVerifyPartialFinishedEvent(t *testing.T, sc *shipyardController, finishedEvent models.Event, eventType, nextEventType string, nextStage string) bool {
	err := sc.HandleIncomingEvent(finishedEvent, true)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s) returned %v", *finishedEvent.Type, err)
		return true
	}

	scope, _ := models.NewEventScope(finishedEvent)
	if nextStage == "" {
		nextStage = scope.Stage
	}
	// check triggeredEvent collection -> should still contain <eventType>.triggered event
	triggeredEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(eventType),
		Stage:   &scope.Stage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)
	done := fake.ShouldContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(eventType), scope.Stage, nil)
	if done {
		return true
	}

	// check triggeredEvent collection -> should not contain <nextEventType>.triggered event
	triggeredEvents, _ = sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:    keptnv2.GetTriggeredEventType(nextEventType),
		Stage:   &nextStage,
		Service: common.Stringp("carts"),
		Source:  common.Stringp("shipyard-controller"),
	}, common.TriggeredEvent)

	done = fake.ShouldNotContainEvent(t, triggeredEvents, keptnv2.GetTriggeredEventType(nextEventType), nextStage)
	if done {
		return true
	}

	// check startedEvent collection -> should still contain one <eventType>.started event
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(eventType),
		Stage:       &scope.Stage,
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(finishedEvent.Triggeredid),
	}, common.StartedEvent)
	if len(startedEvents) != 1 {
		t.Errorf("List of started events does not hold proper number of events. Expected 1 but got %d", len(startedEvents))
		return true
	}
	done = fake.ShouldContainEvent(t, startedEvents, keptnv2.GetStartedEventType(eventType), scope.Stage, nil)
	return done
}

func sendAndVerifyStartedEvent(t *testing.T, sc *shipyardController, taskName string, triggeredID string, stage string, fromSource string) bool {
	err := sc.HandleIncomingEvent(getStartedEvent(stage, triggeredID, taskName, fromSource), true)
	if err != nil {
		t.Errorf("STEP failed: HandleIncomingEvent(%s.started) returned %v", taskName, err)
		return true
	}
	// check startedEvent collection -> should contain <taskName>.started event
	startedEvents, _ := sc.eventRepo.GetEvents("test-project", common.EventFilter{
		Type:        keptnv2.GetStartedEventType(taskName),
		Stage:       common.Stringp(stage),
		Service:     common.Stringp("carts"),
		TriggeredID: common.Stringp(triggeredID),
	}, common.StartedEvent)
	return fake.ShouldContainEvent(t, startedEvents, keptnv2.GetStartedEventType(taskName), stage, nil)
}

func getArtifactDeliveryTriggeredEvent(stage string) models.Event {
	return models.Event{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: "test-project",
				Stage:   stage,
				Service: "carts",
			},
			ConfigurationChange: struct {
				Values map[string]interface{} `json:"values"`
			}{
				Values: map[string]interface{}{
					"image": "carts",
				},
			},
			Deployment: keptnv2.DeploymentTriggeredData{
				DeploymentURIsPublic: []string{"uri"},
				DeploymentStrategy:   "direct",
			},
		},
		Extensions:     nil,
		ID:             "artifact-delivery-triggered-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "",
		Type:           common.Stringp("sh.keptn.event.dev.artifact-delivery.triggered"),
	}
}

func getStartedEvent(stage string, triggeredID string, eventType string, source string) models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           fake.EventScope{Project: "test-project", Stage: stage, Service: "carts"},
		Extensions:     nil,
		ID:             eventType + "-" + source + "-started-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp(source),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    triggeredID,
		Type:           common.Stringp(keptnv2.GetStartedEventType(eventType)),
	}
}
