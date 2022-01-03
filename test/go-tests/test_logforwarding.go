package go_tests

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const logForwardingTestShipyard = `--- 
apiVersion: spec.keptn.sh/0.2.3
kind: Shipyard
metadata: 
  name: shipyard-log-forwarding
spec: 
  stages: 
    - 
      name: dev
      sequences: 
        - 
          name: evaluation
          tasks: 
            - name: evaluation
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

		if !IsEqual(t, "finished", state.State, "state.State") {
			return false
		}
		return true
	}, 100*time.Second, 2*time.Second)

	// retrieve the integration for the lighthouse service
	integrations, _, err := getIntegrations()
	require.Nil(t, err)

	integrationID := ""
	for _, integration := range integrations {
		if integration.Name == "lighthouse-service" {
			integrationID = integration.ID
		}
	}

	require.NotEmpty(t, t, integrationID)

	var contextLogEntry *models.LogEntry
	require.Eventually(t, func() bool {
		logs, _, err := getLogs(integrationID)
		if len(logs.Logs) == 0 || err != nil {
			t.Log("error logs of lighthouse service not available yet... retrying in 10s")
			return false
		}
		t.Log("received logs of lighthouse service")
		for _, log := range logs.Logs {
			if log.KeptnContext == keptnContextID {
				contextLogEntry = &log
				return true
			}
		}
		return false
	}, 100*time.Second, 10*time.Second)

	// check if log entry for our task sequence context is available

	require.NotEmpty(t, contextLogEntry)
	t.Logf("received expected error log entry: %v", contextLogEntry)
}

func getIntegrations() ([]*models.Integration, *req.Resp, error) {
	integrations := []*models.Integration{}

	resp, err := ApiGETRequest("/controlPlane/v1/uniform/registration", 3)
	if err != nil {
		return nil, nil, err
	}
	err = resp.ToJSON(&integrations)
	if err != nil {
		return nil, nil, err
	}
	return integrations, resp, nil
}

func getLogs(integrationID string) (*models.GetLogsResponse, *req.Resp, error) {
	logs := &models.GetLogsResponse{}

	resp, err := ApiGETRequest("/controlPlane/v1/log?integrationId="+integrationID, 3)
	if err != nil {
		return nil, nil, err
	}
	err = resp.ToJSON(&logs)
	if err != nil {
		return nil, nil, err
	}
	return logs, resp, nil
}
