package go_tests

import (
	"fmt"
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

	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err := GetState(projectName)
		if err != nil {
			return false
		}
		if !IsEqual(t, http.StatusOK, resp.Response().StatusCode, "resp.Response().StatusCode") {
			return false
		}
		if !IsEqual(t, int64(1), states.TotalCount, "states.TotalCount") {
			return false
		}
		if !IsEqual(t, 1, len(states.States), "len(states.States)") {
			return false
		}

		state := states.States[0]

		if !IsEqual(t, projectName, state.Project, "state.Project") {
			return false
		}
		if !IsEqual(t, keptnContextID, state.Shkeptncontext, "state.Shkeptncontext") {
			return false
		}
		if !IsEqual(t, scmodels.SequenceStartedState, state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 1, len(state.Stages), "len(state.Stages)") {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, keptnv2.GetTriggeredEventType("mytask"), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 100*time.Second, 2*time.Second)

	myTaskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	// send .started event
	_, err = keptn.SendTaskStartedEvent(nil, source)

	// abort sequence
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.AbortSequence,
		Stage: "dev",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// send .finished event
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)

	// verify state -> should be finished
	require.Eventually(t, func() bool {
		states, resp, err := GetState(projectName)
		if err != nil {
			return false
		}
		if !IsEqual(t, http.StatusOK, resp.Response().StatusCode, "resp.Response().StatusCode") {
			return false
		}
		state := states.States[0]

		if !IsEqual(t, "finished", state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 1, len(state.Stages), "len(state.Stages)") {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, stageName, stage.Name, "stage.Name") {
			return false
		}

		if !IsEqual(t, keptnv2.GetFinishedEventType(stageName+"."+sequencename), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 100*time.Second, 2*time.Second)

	require.Nil(t, err)

}
