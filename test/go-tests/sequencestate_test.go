package go_tests

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
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
	file, err := createTmpShipyardFile(sequenceStateShipyard)
	require.Nil(t, err)
	defer os.Remove(file)

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

	states := &scmodels.SequenceStates{}

	resp, err := apiGETRequest("/controlPlane/v1/state/" + projectName)

	require.Nil(t, err)

	err = resp.ToJSON(states)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().Status)
	require.Empty(t, states.States)
	require.Equal(t, 0, states.NextPageKey)
	require.Equal(t, 0, states.TotalCount)
}

func apiGETRequest(path string) (*req.Resp, error) {
	apiToken, err := keptnutils.GetKeptnAPITokenFromSecret(false, os.Getenv("KEPTN_NAMESPACE"), "keptn-api-token")
	if err != nil {
		return nil, err
	}
	keptnAPIURL := os.Getenv("KEPTN_ENDPOINT")
	if keptnAPIURL == "" {
		serviceIP, err := keptnutils.GetKeptnEndpointFromService(false, os.Getenv("KEPTN_NAMESPACE"), "api-gateway-nginx")
		if err != nil {
			return nil, err
		}
		keptnAPIURL = "http://" + serviceIP + "/api"
	}

	authHeader := req.Header{
		"Accept":        "application/json",
		"x-token": apiToken,
	}

	r, err := req.Get(keptnAPIURL + path, authHeader)
	if err != nil {
		return nil, err
	}

	return r, nil

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
