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

const sequenceTimeoutWithTriggeredAfterShipyard = `--- 
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
              triggeredAfter: "1m"
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
	}, 3)
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
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{scmodels.TimedOut})
	t.Log("received the expected state!")
}

func Test_SequenceTimeoutDelayedTask(t *testing.T) {
	projectName := "sequence-timeout-delay"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceTimeoutWithTriggeredAfterShipyard)
	require.Nil(t, err)
	defer os.Remove(sequenceStateShipyardFilePath)

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

	// trigger the task sequence
	t.Log("starting task sequence")
	keptnContextID, err := TriggerSequence(projectName, serviceName, "dev", "delivery", nil)
	require.Nil(t, err)

	// wait a minute and make verify that the sequence has not been timed out
	<-time.After(30 * time.Second)

	// also, the unknown.triggered event should not have been sent yet
	triggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("unknown"))
	require.Nil(t, err)
	require.Nil(t, triggeredEvent)

	states, _, err := GetState(projectName)
	require.Nil(t, err)

	require.Len(t, states.States, 1)
	state := states.States[0]

	require.Equal(t, scmodels.SequenceStartedState, state.State)

	// after some time, the unknown.triggered event should be available
	require.Eventually(t, func() bool {
		triggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, "dev", keptnv2.GetTriggeredEventType("unknown"))
		if err != nil {
			return false
		}
		if triggeredEvent == nil {
			return false
		}
		return true
	}, 65*time.Second, 5*time.Second)

}

func setShipyardControllerTaskTimeout(t *testing.T, timeoutValue string) error {
	return SetShipyardControllerEnvVar(t, "TASK_STARTED_WAIT_DURATION", timeoutValue)
}
