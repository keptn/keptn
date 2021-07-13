package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	scmodels "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const sequenceTimeoutShipyard = `--- 
apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata: 
  name: shipyard-sockshop
spec: 
  stages: 
    - 
      name: dev
      sequences: 
        - 
          name: delivery
          tasks: 
            - 
              name: unknown`

func Test_SequenceTimeout(t *testing.T) {
	projectName := "sequence-timeout"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceTimeoutShipyard)
	require.Nil(t, err)
	defer os.Remove(sequenceStateShipyardFilePath)

	source := "golang-test"

	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, sequenceStateShipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	setShipyardControllerTimeout(t, "10s")
	defer func() {
		setShipyardControllerTimeout(t, "20m")
	}()

	eventType := keptnv2.GetTriggeredEventType("dev.delivery")

	// trigger the task sequence
	t.Log("starting task sequence")
	resp, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
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

	// wait for the recreated state to be available
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, scmodels.TimedOut)
	require.Eventually(t, func() bool {
		states, _, err := getState(projectName)
		if err != nil {
			return false
		}
		for _, state := range states.States {
			if state.Shkeptncontext == *context.KeptnContext && state.State == scmodels.TimedOut {
				return true
			}
		}
		return false
	}, 2*time.Minute, 10*time.Second)
	t.Log("received the expected state!")
}

func setShipyardControllerTimeout(t *testing.T, timeoutValue string) {
	t.Log("setting TASK_STARTED_WAIT_DURATION of shipyard controller to 10s")
	// temporarily set the timeout value to a lower value
	_, err := ExecuteCommand(fmt.Sprintf("kubectl -n %s set env deployment shipyard-controller TASK_STARTED_WAIT_DURATION=%s", GetKeptnNameSpaceFromEnv(), timeoutValue))
	require.Nil(t, err)

	t.Log("restarting shipyard controller pod")
	err = RestartPod("shipyard-controller")
	require.Nil(t, err)

	// wait a bit to make sure we are waiting for the correct pod to be started
	<-time.After(10 * time.Second)
	t.Log("waiting for shipyard controller pod to be ready again")
	err = WaitForPodOfDeployment("shipyard-controller")
	require.Nil(t, err)
}
