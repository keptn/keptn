package go_tests

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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

const defaultKeptnNamespace = "keptn"

func Test_SequenceStateIntegrationTest(t *testing.T) {
	if os.Getenv("KEPTN_NAMESPACE") == "" {
		os.Setenv("KEPTN_NAMESPACE", defaultKeptnNamespace)
	}
	projectName := "state"
	serviceName := "my-service"
	file, err := createTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer os.Remove(file)

	source := "golang-test"

	uniform := []string{"helm-service", "jmeter-service", "lighthouse-service"}
	if err := scaleDownUniform(uniform); err != nil {
		t.Errorf("scaling down uniform failed: %s", err.Error())
	}

	defer func() {
		if err := scaleUpUniform(uniform); err != nil {
			t.Errorf("could not scale up uniform: " + err.Error())
		}
	}()

	resp, err := apiGETRequest("/controlPlane/v1/project/" + projectName)

	if resp.Response().StatusCode != http.StatusNotFound {
		// delete project if it exists
		_, err = executeCommand(fmt.Sprintf("keptn delete project %s", projectName))
		require.NotNil(t, err)
	}

	output, err := executeCommand(fmt.Sprintf("keptn create project %s --shipyard=./%s", projectName, file))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	output, err = executeCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	states, resp, err := getState(projectName)

	// send a delivery.triggered event
	eventType := keptnv2.GetTriggeredEventType("dev.delivery")
	resp, err = apiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
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
		Shkeptnspecversion: "0.2.0",
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
		if !isEqual(t, "resp.Response().StatusCode", http.StatusOK, resp.Response().StatusCode) {
			return false
		}
		if !isEqual(t, "states.TotalCount", int64(1), states.TotalCount) {
			return false
		}
		if !isEqual(t, "len(states.States)", 1, len(states.States)) {
			return false
		}

		state := states.States[0]

		if !isEqual(t, "state.Project", projectName, state.Project) {
			return false
		}
		if !isEqual(t, "state.Shkeptncontext", *context.KeptnContext, state.Shkeptncontext) {
			return false
		}
		if !isEqual(t, "state.State", "triggered", state.State) {
			return false
		}

		if !isEqual(t, "len(state.Stages)", 1, len(state.Stages)) {
			return false
		}

		stage := state.Stages[0]

		if !isEqual(t, "stage.Name", "dev", stage.Name) {
			return false
		}
		if !isEqual(t, "stage.Image", "carts:test", stage.Image) {
			return false
		}

		if !isEqual(t, "stage.LatestEvent.Type", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), stage.LatestEvent.Type) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get deployment.triggered event
	deploymentTriggeredEvent, err := getLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))
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

		if !isEqual(t, "len(state.Stages)", 1, len(state.Stages)) {
			return false
		}

		stage := state.Stages[0]

		if !isEqual(t, "stage.LatestEvent.Type", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), stage.LatestEvent.Type) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	// get evaluation.triggered event
	evaluationTriggeredEvent, err := getLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName))
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

		if !isEqual(t, "state.State", "triggered", state.State) {
			return false
		}

		if !isEqual(t, "len(state.Stages)", 2, len(state.Stages)) {
			return false
		}

		devStage := state.Stages[0]

		if !isEqual(t, "devStage.LatestEvaluation.Score", 100.0, devStage.LatestEvaluation.Score) {
			return false
		}

		if !isEqual(t, "devStage.LatestEvent.Type", keptnv2.GetFinishedEventType("dev.delivery"), devStage.LatestEvent.Type) {
			return false
		}

		stagingStage := state.Stages[1]

		if !isEqual(t, "stagingStage.LatestEvent.Type", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), stagingStage.LatestEvent.Type) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)

	deploymentTriggeredEvent, err = getLatestEventOfType(*context.KeptnContext, projectName, "staging", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))

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

		if !isEqual(t, "state.State", "finished", state.State) {
			return false
		}

		if !isEqual(t, "len(state.Stages)", 2, len(state.Stages)) {
			return false
		}

		stagingStage := state.Stages[1]

		if !isEqual(t, "stagingStage.LatestEvent.Type", keptnv2.GetFinishedEventType("staging.delivery"), stagingStage.LatestEvent.Type) {
			return false
		}

		return true
	}, 10*time.Second, 2*time.Second)
}

func getLatestEventOfType(keptnContext, projectName, stage, eventType string) (*models.KeptnContextExtendedCE, error) {
	resp, err := apiGETRequest("/mongodb-datastore/event?project=" + projectName + "&keptnContext=" + keptnContext + "&stage=" + stage + "&type=" + eventType)
	if err != nil {
		return nil, err
	}
	events := &models.Events{}
	if err := resp.ToJSON(events); err != nil {
		return nil, err
	}
	if len(events.Events) > 0 {
		return events.Events[0], nil
	}
	return nil, nil
}

func isEqual(t *testing.T, property string, expected, actual interface{}) bool {
	if expected != actual {
		t.Logf("%s: expected %v, got %v", property, expected, actual)
		return false
	}
	return true
}

type APIEventSender struct {
}

func (sender *APIEventSender) SendEvent(event cloudevents.Event) error {
	_, err := apiPOSTRequest("/v1/event", event)
	return err
}

func getState(projectName string) (*scmodels.SequenceStates, *req.Resp, error) {
	states := &scmodels.SequenceStates{}

	resp, err := apiGETRequest("/controlPlane/v1/state/" + projectName)
	err = resp.ToJSON(states)

	return states, resp, err
}

func apiGETRequest(path string) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := getApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := req.Header{
		"Accept":  "application/json",
		"x-token": apiToken,
	}

	r, err := req.Get(keptnAPIURL+path, authHeader)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func apiPOSTRequest(path string, payload interface{}) (*req.Resp, error) {
	apiToken, keptnAPIURL, err := getApiCredentials()
	if err != nil {
		return nil, err
	}

	authHeader := req.Header{
		"Accept":  "application/json",
		"x-token": apiToken,
	}

	r, err := req.Post(keptnAPIURL+path, authHeader, req.BodyJSON(payload))
	if err != nil {
		return nil, err
	}

	return r, nil
}

func getApiCredentials() (string, string, error) {
	apiToken, err := keptnutils.GetKeptnAPITokenFromSecret(false, os.Getenv("KEPTN_NAMESPACE"), "keptn-api-token")
	if err != nil {
		return "", "", err
	}
	keptnAPIURL := os.Getenv("KEPTN_ENDPOINT")
	if keptnAPIURL == "" {
		serviceIP, err := keptnutils.GetKeptnEndpointFromService(false, os.Getenv("KEPTN_NAMESPACE"), "api-gateway-nginx")
		if err != nil {
			return "", "", err
		}
		keptnAPIURL = "http://" + serviceIP + "/api"
	}
	return apiToken, keptnAPIURL, nil
}

func scaleDownUniform(deployments []string) error {
	for _, deployment := range deployments {
		if err := keptnutils.ScaleDeployment(false, deployment, os.Getenv("KEPTN_NAMESPACE"), 0); err != nil {
			// log the error but continue
			fmt.Println("could not scale down deployment: " + err.Error())
		}
	}
	return nil
}

func scaleUpUniform(deployments []string) error {
	for _, deployment := range deployments {
		if err := keptnutils.ScaleDeployment(false, deployment, os.Getenv("KEPTN_NAMESPACE"), 1); err != nil {
			// log the error but continue
			fmt.Println("could not scale up deployment: " + err.Error())
		}
	}
	return nil
}

func createTmpShipyardFile(shipyardContent string) (string, error) {
	file, err := ioutil.TempFile(".", "shipyard-*.yaml")
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file.Name(), []byte(shipyardContent), os.ModeAppend); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

func executeCommand(cmd string) (string, error) {
	split := strings.Split(cmd, " ")
	if len(split) == 0 {
		return "", errors.New("invalid command")
	}
	return keptnutils.ExecuteCommand(split[0], split[1:])
}
