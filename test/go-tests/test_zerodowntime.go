package go_tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

const zeroDownTimeShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
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
                deploymentstrategy: "direct"`

// Test_ZeroDownTimeTriggerSequence tests whether a sequence is started event though the shipyard controller is down at the moment where the sequence.triggered event is sent to the API.
// This is for testing the at least once delivery guarantee achieved by using JetStream for the shipyard controller (see http://github.com/keptn/keptn/issue/6685).
func Test_ZeroDownTimeTriggerSequence(t *testing.T) {
	projectName := "zero-downtime"
	serviceName := "my-service"
	stageName := "dev"
	sequenceName := "delivery"
	shipyardFile, err := CreateTmpShipyardFile(zeroDownTimeShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(shipyardFile)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", shipyardFile, err)
		}
	}()

	// check if the project 'state' is already available - if not, delete it before creating it again
	// check if the project is already available - if not, delete it before creating it again
	projectName, err = CreateProject(projectName, shipyardFile)
	require.Nil(t, err)

	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// scale down the shipyard controller
	err = ScaleDownUniform([]string{"shipyard-controller"})

	require.Nil(t, err)

	err = WaitForDeploymentToBeScaledDown("shipyard-controller")

	require.Nil(t, err)

	keptnContext, err := TriggerSequence(projectName, serviceName, stageName, sequenceName, nil)

	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	// now that we have sent the triggered event, scale the shipyard controller back up again
	err = ScaleUpUniform([]string{"shipyard-controller"}, 1)

	require.Nil(t, err)

	// eventually, the triggered event should be received by the shipyard controller, and a sequence state should have been created
	var states *models.SequenceStates
	var err2 error
	require.Eventually(t, func() bool {
		states, _, err2 = GetState(projectName)
		if err != nil {
			return false
		}
		if len(states.States) == 0 {
			return false
		}
		return true
	}, 1*time.Minute, 5*time.Second)

	require.Nil(t, err2)

	require.Equal(t, keptnContext, states.States[0].Shkeptncontext)

	// check if the first triggered event for the sequence has been sent out
	require.Eventually(t, func() bool {
		triggeredEvent, err := GetLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetTriggeredEventType("delivery"))
		if err != nil || triggeredEvent == nil {
			return false
		}
		return true
	}, 30*time.Second, 5*time.Second)
}
