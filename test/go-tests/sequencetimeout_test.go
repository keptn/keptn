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

	err = setShipyardControllerTaskTimeout(t, "10s")
	defer func() {
		_ = setShipyardControllerTaskTimeout(t, "20m")
	}()
	require.Nil(t, err)

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
	VerifySequenceEndsUpInState(t, projectName, context, scmodels.TimedOut, 2*time.Minute)
	t.Log("received the expected state!")
}

func setShipyardControllerTaskTimeout(t *testing.T, timeoutValue string) error {
	_, err := ExecuteCommand(fmt.Sprintf("kubectl -n %s set env deployment shipyard-controller TASK_STARTED_WAIT_DURATION=%s", GetKeptnNameSpaceFromEnv(), timeoutValue))
	if err != nil {
		return err
	}

	t.Log("restarting shipyard controller pod")
	err = RestartPod("shipyard-controller")
	if err != nil {
		return err
	}

	// wait 10s to make sure we wait for the updated pod to be ready
	<-time.After(10 * time.Second)
	t.Log("waiting for shipyard controller pod to be ready again")
	return WaitForPodOfDeployment("shipyard-controller")
}
