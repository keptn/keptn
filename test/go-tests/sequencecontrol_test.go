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
            - name: "mytask"
            - name: "mynexttask"`

func TestAbortSequence(t *testing.T) {

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

	myTaskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)

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

func TestPauseAndResumeSequence(t *testing.T) {
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

	myTaskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	t.Log("sending task started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)

	t.Log("pausing sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.PauseSequence,
		Stage: "",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("sending task finished event")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{&keptnContextID}, 2*time.Minute, []string{scmodels.SequencePaused})

	t.Log("verifying that the next task has not being triggered")
	time.Sleep(10 * time.Second) //sorry, but I don't know how to verify it without a waiting
	nextTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("mynextTask"))
	require.Nil(t, err)
	require.Nil(t, nextTriggeredEvent)

	t.Log("resuming sequence")
	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.ResumeSequence,
		Stage: "",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	t.Log("verifying that the next task has being triggered")
	require.Eventually(t, func() bool {
		nextTaskTriggered, _ := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("mynexttask"))
		return nextTaskTriggered != nil
	}, 20*time.Second, 2*time.Second)

}
