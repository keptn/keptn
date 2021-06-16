package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const sequenceStateShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "evaluation"


    - name: "staging"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "evaluation"`

func Test_SequenceStateIntegrationTest(t *testing.T) {
	projectName := "state"
	serviceName := "my-service"
	file, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer os.Remove(file)

	source := "golang-test"

	uniform := []string{"helm-service", "lighthouse-service"}

	// scale down the services that are usually involved in the sequence defined in the shipyard above.
	// this way we can control the events sent during this sequence and check whether the state is updated appropriately
	if err := ScaleDownUniform(uniform); err != nil {
		t.Errorf("scaling down uniform failed: %s", err.Error())
	}

	defer func() {
		if err := ScaleUpUniform(uniform); err != nil {
			t.Errorf("could not scale up uniform: " + err.Error())
		}
	}()

	// check if the project 'state' is already available - if not, delete it before creating it again
	resp, err := ApiGETRequest("/controlPlane/v1/project/" + projectName)

	if resp.Response().StatusCode != http.StatusNotFound {
		// delete project if it exists
		_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", projectName))
		require.NotNil(t, err)
	}

	output, err := ExecuteCommand(fmt.Sprintf("keptn create project %s --shipyard=./%s", projectName, file))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	output, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	states, resp, err := getState(projectName)

	// send a delivery.triggered event
	eventType := keptnv2.GetTriggeredEventType("dev.delivery")

	resp, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "dev",
				Service: serviceName,
			},
			ConfigurationChange: keptnv2.ConfigurationChange{
				Values: map[string]interface{}{"image": "carts:test"},
			},
		},
		ID:                 uuid.NewString(),
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Type:               &eventType,
	})
	require.Nil(t, err)
	body := resp.String()
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)
	require.NotEmpty(t, body)

	context := &models.EventContext{}
	err = resp.ToJSON(context)
	require.Nil(t, err)
	require.NotNil(t, context.KeptnContext)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = getState(projectName)
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
		if !IsEqual(t, *context.KeptnContext, state.Shkeptncontext, "state.Shkeptncontext") {
			return false
		}
		if !IsEqual(t, "triggered", state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 1, len(state.Stages), "len(state.Stages)") {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, "dev", stage.Name, "stage.Name") {
			return false
		}
		if !IsEqual(t, "carts:test", stage.Image, "stage.Image") {
			return false
		}

		if !IsEqual(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get deployment.triggered event
	deploymentTriggeredEvent, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))
	require.Nil(t, err)
	require.NotNil(t, deploymentTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*deploymentTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})

	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = getState(projectName)
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
		if state.Shkeptncontext != *context.KeptnContext {
			return false
		}
		if state.State != "triggered" {
			return false
		}

		if len(state.Stages) != 1 {
			return false
		}

		stage := state.Stages[0]

		if stage.LatestEvent.Type != keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	_, err = keptn.SendTaskFinishedEvent(nil, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = getState(projectName)
		if err != nil {
			return false
		}
		state := states.States[0]

		if !IsEqual(t, 1, len(state.Stages), "len(state.Stages)") {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get evaluation.triggered event
	evaluationTriggeredEvent, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName))
	require.Nil(t, err)
	require.NotNil(t, evaluationTriggeredEvent)

	cloudEvent = keptnv2.ToCloudEvent(*evaluationTriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	// send started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// send finished event with score
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EvaluationFinishedEventData{
		EventData: keptnv2.EventData{
			Status: keptnv2.StatusSucceeded,
			Result: keptnv2.ResultPass,
		},
		Evaluation: keptnv2.EvaluationDetails{
			Score: 100.0,
		},
	}, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = getState(projectName)
		if err != nil {
			return false
		}
		state := states.States[0]

		if !IsEqual(t, "triggered", state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 2, len(state.Stages), "len(state.Stages)") {
			return false
		}

		devStage := state.Stages[0]

		if !IsEqual(t, 100.0, devStage.LatestEvaluation.Score, "devStage.LatestEvaluation.Score") {
			return false
		}

		if !IsEqual(t, keptnv2.GetFinishedEventType("dev.delivery"), devStage.LatestEvent.Type, "devStage.LatestEvent.Type") {
			return false
		}

		stagingStage := state.Stages[1]

		if !IsEqual(t, keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), stagingStage.LatestEvent.Type, "stagingStage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	deploymentTriggeredEvent, err = GetLatestEventOfType(*context.KeptnContext, projectName, "staging", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))

	require.Nil(t, err)
	require.NotNil(t, deploymentTriggeredEvent)

	cloudEvent = keptnv2.ToCloudEvent(*deploymentTriggeredEvent)

	keptn, err = keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	// send started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// send finished event with result=fail
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultFailed,
	}, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = getState(projectName)
		if err != nil {
			return false
		}
		state := states.States[0]

		if !IsEqual(t, "finished", state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 2, len(state.Stages), "len(state.Stages)") {
			return false
		}

		stagingStage := state.Stages[1]

		if !IsEqual(t, keptnv2.GetFinishedEventType("staging.delivery"), stagingStage.LatestEvent.Type, "stagingStage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)
}

func getState(projectName string) (*scmodels.SequenceStates, *req.Resp, error) {
	states := &scmodels.SequenceStates{}

	resp, err := ApiGETRequest("/controlPlane/v1/sequence/" + projectName)
	err = resp.ToJSON(states)

	return states, resp, err
}
