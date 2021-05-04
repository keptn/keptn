package go_tests

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/assert"
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
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"


    - name: "staging"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "test"
              properties:
                teststrategy: "performance"
            - name: "evaluation"
            - name: "release"
        - name: "rollback"
          triggeredOn:
            - event: "staging.delivery.finished"
              selector:
                match:
                  result: "fail"
          tasks:
            - name: "rollback"`

func Test_SequenceStateIntegrationTest(t *testing.T) {
	// TODO
	os.Setenv("KEPTN_NAMESPACE", "keptn")
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

	_, err = executeCommand(fmt.Sprintf("keptn delete project %s", projectName))
	require.Nil(t, err)

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

	// verify state

	assert.Eventually(t, func() bool {
		states, resp, err = getState(projectName)
		if err != nil {
			return false
			return false
		}

		if states.TotalCount != 1 {
			return false
			return false
		}

		if len(states.States) != 1 {
			return false
		}
		require.Equal(t, http.StatusOK, resp.Response().StatusCode)
		require.Empty(t, states.States)
		require.Empty(t, states.NextPageKey)
		require.Empty(t, states.TotalCount)

		return true
	}, 10*time.Second, 2*time.Second)

}

func verifyStateWithRetry(projectName string, retries int, verify func(resp *req.Resp, states *scmodels.SequenceStates, err error) error) error {
	for i := 0; i < retries; i = i + 1 {
		states, resp, err := getState(projectName)

		if verifyErr := verify(resp, states, err); verifyErr == nil {
			return nil
		}
		<-time.After(5 * time.Second)
	}
	return errors.New("could not verify sequence state")
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
			return err
		}
	}
	return nil
}

func scaleUpUniform(deployments []string) error {
	for _, deployment := range deployments {
		if err := keptnutils.ScaleDeployment(false, deployment, os.Getenv("KEPTN_NAMESPACE"), 1); err != nil {
			return err
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
