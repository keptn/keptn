package go_tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const logForwardingTestShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-log-forwarding"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "evaluation"
          tasks:
            - name: "evaluation"
			  properties:
                timeframe: "invalid"`

func Test_LogForwarding(t *testing.T) {
	projectName := "log-forwarding"
	serviceName := "my-service"
	stageName := "dev"
	sequenceName := "evaluation"
	shipyardFilePath, err := CreateTmpShipyardFile(logForwardingTestShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

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
		states, resp, err := getState(projectName)
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

		if !IsEqual(t, "finished", state.State, "state.State") {
			return false
		}
		return true
	}, 10*time.Second, 2*time.Second)
}
