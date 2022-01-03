package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"strings"
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
            - name: "delivery"
              properties:
                deploymentstrategy: "direct"
            - name: "evaluation"


    - name: "staging"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "delivery"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "evaluation"`

func Test_SequenceState(t *testing.T) {
	projectName := "state"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(sequenceStateShipyardFilePath)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", sequenceStateShipyardFilePath, err)
		}
	}()

	source := "golang-test"

	uniform := []string{"lighthouse-service"}

	// scale down the services that are usually involved in the sequence defined in the shipyard above.
	// this way we can control the events sent during this sequence and check whether the state is updated appropriately
	if err := ScaleDownUniform(uniform); err != nil {
		t.Errorf("scaling down uniform failed: %s", err.Error())
	}

	defer func() {
		if err := ScaleUpUniform(uniform, 1); err != nil {
			t.Errorf("could not scale up uniform: " + err.Error())
		}
	}()

	// check if the project 'state' is already available - if not, delete it before creating it again
	// check if the project is already available - if not, delete it before creating it again
	err = CreateProject(projectName, sequenceStateShipyardFilePath, true)
	require.Nil(t, err)

	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	states, resp, err := GetState(projectName)

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
	}, 3)
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
		states, resp, err = GetState(projectName)
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
		if !IsEqual(t, scmodels.SequenceStartedState, state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 1, len(state.Stages), "len(state.Stages)") {
			return false
		}

		stage := state.Stages[0]

		if !IsEqual(t, "dev", stage.Name, "stage.Name") {
			return false
		}

		if !IsEqual(t, keptnv2.GetTriggeredEventType("delivery"), stage.LatestEvent.Type, "stage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get deployment.triggered event
	deliveryTriggeredEvent, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("delivery"))
	require.Nil(t, err)
	require.NotNil(t, deliveryTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*deliveryTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})

	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = GetState(projectName)
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
		if state.State != scmodels.SequenceStartedState {
			return false
		}

		if len(state.Stages) != 1 {
			return false
		}

		stage := state.Stages[0]

		if stage.LatestEvent.Type != keptnv2.GetStartedEventType("delivery") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	_, err = keptn.SendTaskFinishedEvent(nil, source)
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = GetState(projectName)
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
	}, "lighthouse-service")
	require.Nil(t, err)

	// verify state
	require.Eventually(t, func() bool {
		states, resp, err = GetState(projectName)
		if err != nil {
			return false
		}
		state := states.States[0]

		if !IsEqual(t, scmodels.SequenceStartedState, state.State, "state.State") {
			return false
		}

		if !IsEqual(t, 2, len(state.Stages), "len(state.Stages)") {
			return false
		}

		devStage := state.Stages[0]

		if devStage.LatestEvaluation == nil {
			return false
		}

		if !IsEqual(t, 100.0, devStage.LatestEvaluation.Score, "devStage.LatestEvaluation.Score") {
			return false
		}

		if !IsEqual(t, keptnv2.GetFinishedEventType("dev.delivery"), devStage.LatestEvent.Type, "devStage.LatestEvent.Type") {
			return false
		}

		stagingStage := state.Stages[1]

		if !IsEqual(t, keptnv2.GetTriggeredEventType("delivery"), stagingStage.LatestEvent.Type, "stagingStage.LatestEvent.Type") {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	deliveryTriggeredEvent, err = GetLatestEventOfType(*context.KeptnContext, projectName, "staging", keptnv2.GetTriggeredEventType("delivery"))

	require.Nil(t, err)
	require.NotNil(t, deliveryTriggeredEvent)

	cloudEvent = keptnv2.ToCloudEvent(*deliveryTriggeredEvent)

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
		states, resp, err = GetState(projectName)
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

func Test_SequenceState_CannotRetrieveShipyard(t *testing.T) {
	projectName := "state-no-shipyard"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(sequenceStateShipyardFilePath)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", sequenceStateShipyardFilePath, err)
		}
	}()

	err = CreateProject(projectName, sequenceStateShipyardFilePath, true)
	require.Nil(t, err)

	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)

	// delete the shipyard file
	_, err = ApiDELETERequest(fmt.Sprintf("/configuration-service/v1/project/%s/resource/shipyard.yaml", projectName), 3)
	require.Nil(t, err)

	_, err = TriggerSequence(projectName, serviceName, "dev", "evaluation", nil)
	require.Nil(t, err)

	var states *scmodels.SequenceStates
	require.Eventually(t, func() bool {
		states, _, err = GetState(projectName)
		if err != nil {
			return false
		} else if states == nil || len(states.States) == 0 {
			return false
		}
		return true
	}, 20*time.Second, 3*time.Second)

	require.Len(t, states.States, 1)
	require.Equal(t, scmodels.SequenceFinished, states.States[0].State)
}

func Test_SequenceState_InvalidShipyard(t *testing.T) {
	projectName := "state-invalid-shipyard"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(sequenceStateShipyardFilePath)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", sequenceStateShipyardFilePath, err)
		}
	}()

	err = CreateProject(projectName, sequenceStateShipyardFilePath, true)
	require.Nil(t, err)

	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)

	// upload a shipyard with an invalid version
	invalidShipyardString := strings.Replace(sequenceQueueShipyard, "spec.keptn.sh/0.2.2", "0.1.7", 1)

	invalidShipyardFile, err := CreateTmpShipyardFile(invalidShipyardString)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(invalidShipyardFile)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", invalidShipyardFile, err)
		}
	}()

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --resource=%s --resourceUri=shipyard.yaml", projectName, invalidShipyardFile))
	require.Nil(t, err)

	_, err = TriggerSequence(projectName, serviceName, "dev", "evaluation", nil)
	require.Nil(t, err)

	var states *scmodels.SequenceStates
	require.Eventually(t, func() bool {
		states, _, err = GetState(projectName)
		if err != nil {
			return false
		} else if states == nil || len(states.States) == 0 {
			return false
		}
		return true
	}, 20*time.Second, 3*time.Second)

	require.Len(t, states.States, 1)
	require.Equal(t, scmodels.SequenceFinished, states.States[0].State)
}

func Test_SequenceState_SequenceNotFound(t *testing.T) {
	projectName := "state-shipyard-unknown-sequence"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(sequenceStateShipyardFilePath)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", sequenceStateShipyardFilePath, err)
		}
	}()

	err = CreateProject(projectName, sequenceStateShipyardFilePath, true)
	require.Nil(t, err)

	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)

	// start a sequence that is not known
	_, err = TriggerSequence(projectName, serviceName, "dev", "unknown", nil)
	require.Nil(t, err)

	var states *scmodels.SequenceStates
	require.Eventually(t, func() bool {
		states, _, err = GetState(projectName)
		if err != nil {
			return false
		} else if states == nil || len(states.States) == 0 {
			return false
		} else if states.States[0].State != scmodels.SequenceFinished {
			return false
		}
		return true
	}, 20*time.Second, 3*time.Second)
}

func Test_SequenceState_RetrieveMultipleSequence(t *testing.T) {
	projectName := "state-retrieve-multiple"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(sequenceStateShipyardFilePath)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", sequenceStateShipyardFilePath, err)
		}
	}()

	err = CreateProject(projectName, sequenceStateShipyardFilePath, true)
	require.Nil(t, err)

	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)

	// start the first sequence
	context1, err := TriggerSequence(projectName, serviceName, "dev", "delivery", nil)
	require.Nil(t, err)

	var states *scmodels.SequenceStates
	require.Eventually(t, func() bool {
		// filter sequences by providing the context ID
		states, _, err = GetStateByContext(projectName, context1)
		if err != nil {
			return false
		} else if states == nil || len(states.States) != 1 {
			return false
		}
		return true
	}, 20*time.Second, 3*time.Second)

	require.Equal(t, context1, states.States[0].Shkeptncontext)

	// start the first sequence
	context2, err := TriggerSequence(projectName, serviceName, "dev", "delivery", nil)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		// filter sequences by providing the context ID
		states, _, err = GetStateByContext(projectName, context2)
		if err != nil {
			return false
		} else if states == nil || len(states.States) != 1 {
			return false
		}
		return true
	}, 20*time.Second, 3*time.Second)

	require.Equal(t, context2, states.States[0].Shkeptncontext)

	// now let's try to fetch both sequences by providing both context IDs
	states, _, err = GetStateByContext(projectName, fmt.Sprintf("%s,%s", context1, context2))
	require.Nil(t, err)

	require.NotNil(t, states)
	require.Len(t, states.States, 2)
	require.Equal(t, context2, states.States[0].Shkeptncontext)
	require.Equal(t, context1, states.States[1].Shkeptncontext)

}
