package go_tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

const sequenceAbortShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "task1"
            - name: "task2"
            - name: "task3"
    - name: "prod"
      sequences:
        - name: "mysequence"
          triggeredOn:
            - event: "dev.mysequence.finished"
          tasks:
            - name: "task4"
            - name: "task5"`

const shipyardWithMultipleStages = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "task1"
    - name: "prod"
      sequences:
        - name: "mysequence"
          triggeredOn:
            - event: "dev.mysequence.finished"
          tasks:
            - name: "task2"`

const shipyardWithParallelStages = `--- 
apiVersion: spec.keptn.sh/0.2.3
kind: Shipyard
metadata: 
  name: shipyard-echo-service
spec: 
  stages: 
    - name: dev
      sequences: 
        - name: mysequence
          tasks: 
            - name: task1
    - name: prod-a
      sequences: 
        - name: mysequence
          tasks: 
            - name: task2
          triggeredOn: 
            - event: "dev.mysequence.finished"
    - name: prod-b
      sequences: 
        - name: mysequence
          tasks: 
            - name: task2
          triggeredOn: 
            - event: "dev.mysequence.finished"`

func Test_SequenceControl_Abort(t *testing.T) {
	projectName := "sequence-abort"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source := "golang-test"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceAbortShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	taskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
	require.Nil(t, err)
	require.NotNil(t, taskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	t.Log("aborting sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.AbortSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	resp, err = ApiGETRequest("/controlPlane/v1/event/triggered/"+keptnv2.GetTriggeredEventType("task1")+"?project="+projectName, 3)
	require.Nil(t, err)

	openTriggeredEvents := &OpenTriggeredEventsResponse{}
	err = resp.ToJSON(openTriggeredEvents)
	require.Nil(t, err)

	require.Empty(t, openTriggeredEvents.Events)

	t.Log("sending task finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)
	require.Nil(t, err)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceAborted})

}

func Test_SequenceControl_AbortQueuedSequence(t *testing.T) {
	projectName := "sequence-abort2"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source := "golang-test"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceAbortShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	taskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
	require.Nil(t, err)
	require.NotNil(t, taskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	// trigger a second sequence which should be put in the queue
	secondContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&secondContextID}, 5*time.Minute, []string{scmodels.SequenceWaitingState})

	// abort the queued sequence
	t.Log("aborting sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, secondContextID), scmodels.SequenceControlCommand{
		State: scmodels.AbortSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&secondContextID}, 2*time.Minute, []string{scmodels.SequenceAborted})
}

func Test_SequenceControl_AbortPausedSequence(t *testing.T) {
	projectName := "sequence-abort3"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source := "golang-test"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceAbortShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	taskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
	require.Nil(t, err)
	require.NotNil(t, taskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	// pause the sequence
	t.Log("pausing sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.PauseSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequencePaused})

	// now trigger another sequence and make sure it is started eventually
	// trigger a second sequence which should be put in the queue
	secondContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// now abort the first sequence
	t.Log("aborting first sequence")
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.AbortSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceAborted})

	// now that the first sequence is aborted, the other sequence should eventually be started
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&secondContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	// also make sure that the triggered event for the first task has been sent
	require.Eventually(t, func() bool {
		taskTriggeredEvent, err := GetLatestEventOfType(secondContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
		if err != nil || taskTriggeredEvent == nil {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)

}

func Test_SequenceControl_AbortPausedSequenceTaskPartiallyFinished(t *testing.T) {
	projectName := "sequence-abort4"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source1 := "golang-test-1"
	source2 := "golang-test-2"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceAbortShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	taskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
	require.Nil(t, err)
	require.NotNil(t, taskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending two task started events")
	_, err = keptn.SendTaskStartedEvent(nil, source1)
	require.Nil(t, err)
	_, err = keptn.SendTaskStartedEvent(nil, source2)
	require.Nil(t, err)

	t.Logf("send one finished event with result 'fail'")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultFailed, Status: keptnv2.StatusSucceeded}, source1)
	require.Nil(t, err)

	// now trigger another sequence and make sure it is started eventually
	// trigger a second sequence which should be put in the queue
	secondContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify that the second sequence gets the triggered status
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&secondContextID}, 2*time.Minute, []string{scmodels.SequenceTriggeredState})

	// pause the sequence
	t.Log("pausing sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.PauseSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequencePaused})

	// now abort the first sequence
	t.Log("aborting first sequence")
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.AbortSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceAborted})

	// now that the first sequence is aborted, the other sequence should eventually be started
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&secondContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	// also make sure that the triggered event for the first task has been sent
	require.Eventually(t, func() bool {
		taskTriggeredEvent, err := GetLatestEventOfType(secondContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
		if err != nil || taskTriggeredEvent == nil {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)

}

func Test_SequenceControl_AbortPausedSequenceMultipleStages(t *testing.T) {
	projectName := "sequence-abort5"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source := "golang-test"

	shipyardFilePath, err := CreateTmpShipyardFile(shipyardWithParallelStages)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state

	var taskTriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskTriggeredEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
		if err != nil || taskTriggeredEvent == nil {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)
	t.Log("sending task finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultPass, Status: keptnv2.StatusSucceeded}, source)
	require.Nil(t, err)

	// wait until sequences in prod-a and prod-b have been started (by retrieving the triggered events of the first task in each stage)
	parallelStagesAreTriggered := func(keptnContextID string) {
		require.Eventually(t, func() bool {
			taskTriggeredEventA, err := GetLatestEventOfType(keptnContextID, projectName, "prod-a", keptnv2.GetTriggeredEventType("task2"))
			if err != nil || taskTriggeredEventA == nil {
				return false
			}
			taskTriggeredEventB, err := GetLatestEventOfType(keptnContextID, projectName, "prod-b", keptnv2.GetTriggeredEventType("task2"))
			if err != nil || taskTriggeredEventB == nil {
				return false
			}
			return true
		}, 1*time.Minute, 10*time.Second)
	}

	parallelStagesAreTriggered(keptnContextID)

	// now trigger another sequence and finish its execution in the first stage
	secondContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	require.Eventually(t, func() bool {
		taskTriggeredEvent, err = GetLatestEventOfType(secondContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
		if err != nil || taskTriggeredEvent == nil {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)

	cloudEvent = keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)
	t.Log("sending task finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultPass, Status: keptnv2.StatusSucceeded}, source)
	require.Nil(t, err)

	// now abort the first sequence
	t.Log("aborting first sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.AbortSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceAborted})

	// now that the first sequence is aborted, the other sequence should start in prod-a and prod-b
	parallelStagesAreTriggered(secondContextID)
}

func Test_SequenceControl_PauseAndResume(t *testing.T) {
	projectName := "sequence-pause-and-resume"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source := "golang-test"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceAbortShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	task1TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
	require.Nil(t, err)
	require.NotNil(t, task1TriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*task1TriggeredEvent)
	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task1 started event")
	keptn.SendTaskStartedEvent(nil, source)

	t.Log("pausing sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.PauseSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task1 finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequencePaused})

	t.Log("verifying that the next task has not being triggered")
	time.Sleep(5 * time.Second) //sorry, but I don't know how to verify it without a waiting
	task2TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("task2"))
	require.Nil(t, err)
	require.Nil(t, task2TriggeredEvent)

	t.Log("resuming sequence")
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.ResumeSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("verifying that the next task has being triggered")
	require.Eventually(t, func() bool {
		task2TriggeredEvent, _ = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task2"))
		return task2TriggeredEvent != nil
	}, 20*time.Second, 2*time.Second)

	cloudEvent = keptnv2.ToCloudEvent(*task2TriggeredEvent)
	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task2 started event")
	keptn.SendTaskStartedEvent(nil, source)

	t.Logf("pausing sequence in stage %s", stageName)
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.PauseSequence,
		Stage: stageName,
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task2 finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	t.Log("verifying that the next task has not been triggered")
	time.Sleep(5 * time.Second) //sorry, but I don't know how to verify it without a waiting
	task3TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("task3"))
	require.Nil(t, err)
	require.Nil(t, task3TriggeredEvent)

	t.Logf("resuming sequence in stage %s", stageName)
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.ResumeSequence,
		Stage: stageName,
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("verifying that the next task has been triggered")
	require.Eventually(t, func() bool {
		task3TriggeredEvent, _ = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task3"))
		return task3TriggeredEvent != nil
	}, 1*time.Minute, 5*time.Second)
}

func Test_SequenceControl_PauseAndResume_2(t *testing.T) {
	projectName := "sequence-pause-and-resume"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	source := "golang-test"

	shipyardFilePath, err := CreateTmpShipyardFile(shipyardWithMultipleStages)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	//TASK 1
	task1TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task1"))
	require.Nil(t, err)
	require.NotNil(t, task1TriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*task1TriggeredEvent)
	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task1 started event")
	keptn.SendTaskStartedEvent(nil, source)

	t.Log("pause sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.PauseSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task1 finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequencePaused})

	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), scmodels.SequenceControlCommand{
		State: scmodels.ResumeSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	require.Eventually(t, func() bool {
		task2TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "prod", keptnv2.GetTriggeredEventType("task2"))
		if err != nil || task2TriggeredEvent == nil {
			return false
		}
		return true
	}, 30*time.Second, 5*time.Second)

}
