package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
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
	err = CreateProject(projectName, shipyardFilePath, true)
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
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.AbortSequence,
		Stage: "",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceFinished})

	require.Nil(t, err)

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
	err = CreateProject(projectName, shipyardFilePath, true)
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
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.PauseSequence,
		Stage: "",
	})
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
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.ResumeSequence,
		Stage: "",
	})
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
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.PauseSequence,
		Stage: stageName,
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task2 finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	t.Log("verifying that the next task has not being triggered")
	time.Sleep(5 * time.Second) //sorry, but I don't know how to verify it without a waiting
	task3TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("task3"))
	require.Nil(t, err)
	require.Nil(t, task3TriggeredEvent)

	t.Logf("resuming sequence in stage %s", stageName)
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.ResumeSequence,
		Stage: stageName,
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("verifying that the next task has being triggered")
	require.Eventually(t, func() bool {
		task3TriggeredEvent, _ = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("task3"))
		return task3TriggeredEvent != nil
	}, 20*time.Second, 2*time.Second)
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
	err = CreateProject(projectName, shipyardFilePath, true)
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
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.PauseSequence,
		Stage: "",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task1 finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequencePaused})

	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.ResumeSequence,
		Stage: "",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequenceStartedState})

	task2TriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "prod", keptnv2.GetTriggeredEventType("task2"))
	require.Nil(t, err)
	require.NotNil(t, task2TriggeredEvent)
}
