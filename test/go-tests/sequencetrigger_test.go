package go_tests

import (
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const sequenceTriggerShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
			  selector:
				match:
				  mytask.result: fail
          tasks:
            - name: "mytask"
            - name: "evaluation"`

func Test_SequenceTriggerIntegrationTest(t *testing.T) {
	PrepareEnvVars()
	projectName := "sequence-loop"
	serviceName := "my-service"
	stageName := "dev"
	sequenceName := "delivery"
	shipyardFilePath, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	source := "golang-test"

	// check if the project is already available - if not, delete it before creating it again
	err = EnsureProjectExists(projectName, shipyardFilePath)
	require.Nil(t, err)

	keptnContextID, err := TriggerSequence(projectName, serviceName, stageName, sequenceName, nil)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContextID)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err := getState(projectName)
		if err != nil {
			return false
		}
		if !IsEqual(t, "resp.Response().StatusCode", http.StatusOK, resp.Response().StatusCode) {
			return false
		}
		if !IsEqual(t, "states.TotalCount", int64(1), states.TotalCount) {
			return false
		}
		if !IsEqual(t, "len(states.States)", 1, len(states.States)) {
			return false
		}

		state := states.States[0]

		if !IsEqual(t, "state.Project", projectName, state.Project) {
			return false
		}
		if !IsEqual(t, "state.Shkeptncontext", keptnContextID, state.Shkeptncontext) {
			return false
		}
		if !IsEqual(t, "state.State", "triggered", state.State) {
			return false
		}

		if !IsEqual(t, "len(state.Stages)", 1, len(state.Stages)) {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, "stage.LatestEvent.Type", keptnv2.GetTriggeredEventType("mytask"), stage.LatestEvent.Type) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get mytask.triggered event
	myTaskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	// send .started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err := getState(projectName)
		if err != nil {
			return false
		}
		if http.StatusOK != resp.Response().StatusCode {
			return false
		}
		state := states.States[0]
		if state.Project != projectName {
			return false
		}
		if state.Shkeptncontext != keptnContextID {
			return false
		}
		if state.State != "triggered" {
			return false
		}

		if len(state.Stages) != 1 {
			return false
		}

		stage := state.Stages[0]

		if stage.LatestEvent.Type != keptnv2.GetStartedEventType("mytask") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// send .finished event with result = fail
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultFailed,
	}, source)
	require.Nil(t, err)

	// verify state -> the same sequence in the same stage should have been triggered again
	require.Eventually(t, func() bool {
		states, resp, err := getState(projectName)
		if err != nil {
			return false
		}
		if !IsEqual(t, "resp.Response().StatusCode", http.StatusOK, resp.Response().StatusCode) {
			return false
		}
		if !IsEqual(t, "states.TotalCount", int64(1), states.TotalCount) {
			return false
		}
		if !IsEqual(t, "len(states.States)", 1, len(states.States)) {
			return false
		}

		state := states.States[0]

		if !IsEqual(t, "state.Project", projectName, state.Project) {
			return false
		}
		if !IsEqual(t, "state.Shkeptncontext", keptnContextID, state.Shkeptncontext) {
			return false
		}
		if !IsEqual(t, "state.State", "triggered", state.State) {
			return false
		}

		if !IsEqual(t, "len(state.Stages)", 1, len(state.Stages)) {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, "stage.Name", stageName, stage.Name) {
			return false
		}

		if !IsEqual(t, "stage.LatestEvent.Type", keptnv2.GetTriggeredEventType("mytask"), stage.LatestEvent.Type) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)
}
