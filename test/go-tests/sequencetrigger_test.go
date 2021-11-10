package go_tests

import (
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const sequenceTriggerShipyard = `apiVersion: "spec.keptn.sh/0.2.2"
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
                  mytask.result: "fail"
          tasks:
            - name: "mytask"
            - name: "othertask"`

func Test_SequenceLoopIntegrationTest(t *testing.T) {
	projectName := "sequence-loop"
	serviceName := "my-service"
	stageName := "dev"
	sequenceName := "delivery"
	shipyardFilePath, err := CreateTmpShipyardFile(sequenceTriggerShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	source := "golang-test"

	// check if the project is already available - if not, delete it before creating it again
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	keptnContextID, err := TriggerSequence(projectName, serviceName, stageName, sequenceName, nil)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContextID)

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
		states, resp, err := GetState(projectName)
		if err != nil {
			return false
		}
		if http.StatusOK != resp.Response().StatusCode {
			return false
		}
		state := states.States[0]
		if !IsEqual(t, state.Project, projectName, "state.Project") {
			return false
		}
		if !IsEqual(t, state.Shkeptncontext, keptnContextID, "state.Shkeptnkontext") {
			return false
		}
		if !IsEqual(t, state.State, scmodels.SequenceStartedState, "state.State") {
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

		if !IsEqual(t, stageName, stage.Name, "stage.Name") {
			return false
		}

		if !IsEqual(t, keptnv2.GetTriggeredEventType("mytask"), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get mytask.triggered event of second iteration
	myTaskTriggeredEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	cloudEvent = keptnv2.ToCloudEvent(*myTaskTriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	// send .started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// send .finished event with result = pass
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultPass,
	}, source)
	require.Nil(t, err)

	// verify state -> now the next task should have been triggered again
	require.Eventually(t, func() bool {
		states, _, err := GetState(projectName)
		if err != nil {
			return false
		}
		stage := states.States[0].Stages[0]

		if !IsEqual(t, stageName, stage.Name, "stage.Name") {
			return false
		}

		if !IsEqual(t, keptnv2.GetTriggeredEventType("othertask"), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get othertask.triggered event of second iteration
	otherTaskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetTriggeredEventType("othertask"))
	require.Nil(t, err)
	require.NotNil(t, otherTaskTriggeredEvent)

	cloudEvent = keptnv2.ToCloudEvent(*otherTaskTriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)
	require.NotNil(t, keptn)

	// send .started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// send .finished event with result = fail
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Result: keptnv2.ResultFailed,
	}, source)
	require.Nil(t, err)

	// verify state -> now the sequence should be finished and not re-triggered again
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

		if !IsEqual(t, keptnv2.GetFinishedEventType("dev.delivery"), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)
}
