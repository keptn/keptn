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
	"testing"
)

const sequenceQueueShipyard = `--- 
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
              name: mytask
    - 
      name: staging
      sequences: 
        - 
          name: delivery
          tasks: 
            - 
              name: mytask`

func Test_SequenceQueue(t *testing.T) {
	projectName := "sequence-queue"
	serviceName := "my-service"
	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceQueueShipyard)
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

	// trigger the task sequence
	context := triggerSequence(t, projectName, serviceName, "dev", "delivery")

	// wait for the recreated state to be available
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, scmodels.SequenceStartedState)
	VerifySequenceEndsUpInState(t, projectName, context, scmodels.SequenceStartedState)
	t.Log("received the expected state!")

	// trigger a second sequence - this one should stay in 'triggered' state until the previous sequence is finished
	secondContext := triggerSequence(t, projectName, serviceName, "dev", "delivery")

	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, scmodels.SequenceTriggeredState)
	VerifySequenceEndsUpInState(t, projectName, secondContext, scmodels.SequenceTriggeredState)
	t.Log("received the expected state!")

	// check if mytask.triggered has been sent for first sequence - this one should be available
	triggeredEventOfFirstSequence, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfFirstSequence)

	// check if mytask.triggered has been sent for second sequence - this one should NOT be available
	triggeredEventOfSecondSequence, err := GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.Nil(t, triggeredEventOfSecondSequence)

	// send .started and .finished event for task of first sequence
	cloudEvent := keptnv2.ToCloudEvent(*triggeredEventOfFirstSequence)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
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

	// now that all tasks for the first sequence have been executed, the second sequence should eventually have the status 'started'
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, scmodels.SequenceStartedState)
	VerifySequenceEndsUpInState(t, projectName, secondContext, scmodels.SequenceStartedState)
	t.Log("received the expected state!")

	// check if mytask.triggered has been sent for second sequence - now it should be available
	triggeredEventOfSecondSequence, err = GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)

	// test if sequences are triggered correctly after a timeout
	err = setShipyardControllerTaskTimeout(t, "10s")
	defer func() {
		// increase the timeout value again
		err = setShipyardControllerTaskTimeout(t, "20m")
		require.Nil(t, err)
	}()
	require.Nil(t, err)

	// trigger the first task sequence - this should time out
	context = triggerSequence(t, projectName, serviceName, "staging", "delivery")
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, scmodels.TimedOut)
	VerifySequenceEndsUpInState(t, projectName, context, scmodels.TimedOut)
	t.Log("received the expected state!")

	// now trigger the second sequence - this should start and a .triggered event for mytask should be sent
	secondContext = triggerSequence(t, projectName, serviceName, "staging", "delivery")
	t.Logf("waiting for state with keptnContext %s to have the status %s", *secondContext.KeptnContext, scmodels.SequenceStartedState)
	VerifySequenceEndsUpInState(t, projectName, secondContext, scmodels.SequenceStartedState)
	triggeredEventOfSecondSequence, err = GetLatestEventOfType(*secondContext.KeptnContext, projectName, "staging", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)
}

func triggerSequence(t *testing.T, projectName, serviceName, stageName, sequenceName string) *models.EventContext {
	source := "golang-test"
	eventType := keptnv2.GetTriggeredEventType(stageName + "." + sequenceName)
	t.Log("starting task sequence")
	resp, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.DeploymentTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
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
	return context
}
