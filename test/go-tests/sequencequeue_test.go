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
	"sync"
	"testing"
	"time"
)

const sequenceQueueShipyard = `--- 
apiVersion: spec.keptn.sh/0.2.2
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
          name: delivery-with-approval
          tasks: 
            - 
              name: approval
              properties: 
                pass: manual
                warning: manual
            - 
              name: mytask
    - 
      name: staging
      sequences: 
        - 
          name: delivery
          tasks: 
            - 
              name: mytask
    - 
      name: qg
      sequences: 
        - 
          name: evaluation
          tasks: 
            - 
              name: approval
              properties: 
                pass: automatic
                warning: automatic
            - 
              name: mytask
            - 
              name: evaluation`

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

	// ------------------------------------
	// Scenario 1: make sure sequences are queued correctly
	// ------------------------------------
	// trigger the task sequence
	context := triggerSequence(t, projectName, serviceName, "dev", "delivery")

	// wait for the sequence state to be available
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{scmodels.SequenceStartedState})
	t.Log("received the expected state!")

	// trigger a second sequence - this one should stay in 'triggered' state until the previous sequence is finished
	secondContext := triggerSequence(t, projectName, serviceName, "dev", "delivery")

	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{scmodels.SequenceTriggeredState})
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
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{scmodels.SequenceStartedState})
	t.Log("received the expected state!")

	// check if mytask.triggered has been sent for second sequence - now it should be available
	triggeredEventOfSecondSequence, err = GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)

	// ------------------------------------
	// Scenario 2: test if sequences are triggered correctly after a timeout
	// ------------------------------------

	err = setShipyardControllerTaskTimeout(t, "10s")
	defer func() {
		// increase the timeout value again
		err = setShipyardControllerTaskTimeout(t, "20m")
		require.Nil(t, err)
	}()
	require.Nil(t, err)

	// wait a bit to make sure the .triggered event is received by the new instance of the shipyard controller
	<-time.After(10 * time.Second)

	// trigger the first task sequence - this should time out
	context = triggerSequence(t, projectName, serviceName, "staging", "delivery")
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{scmodels.TimedOut})
	t.Log("received the expected state!")

	// now trigger the second sequence - this should start and a .triggered event for mytask should be sent
	secondContext = triggerSequence(t, projectName, serviceName, "staging", "delivery")
	t.Logf("waiting for state with keptnContext %s to have the status %s", *secondContext.KeptnContext, scmodels.SequenceStartedState)
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{scmodels.SequenceStartedState})
	triggeredEventOfSecondSequence, err = GetLatestEventOfType(*secondContext.KeptnContext, projectName, "staging", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)

	// ------------------------------------
	// Scenario 3: special case approval: once an approval has been triggered, another sequence should take over
	// ------------------------------------

	require.Nil(t, err)
	// increase the task timeout again
	err = setShipyardControllerTaskTimeout(t, "20m")
	require.Nil(t, err)

	// wait a bit to make sure the .triggered event is received by the new instance of the shipyard controller
	<-time.After(10 * time.Second)

	t.Log("starting delivery-with-approval sequence")
	context = triggerSequence(t, projectName, serviceName, "dev", "delivery-with-approval")
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{scmodels.SequenceStartedState})

	// check if approval.triggered has been sent for sequence - now it should be available
	approvalTriggeredEvent, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName))
	require.Nil(t, err)
	require.NotNil(t, approvalTriggeredEvent)

	// send the approval.started event to make sure the sequence will not be cancelled due to a timeout
	approvalTriggeredCE := keptnv2.ToCloudEvent(*approvalTriggeredEvent)
	keptnHandler, err := keptnv2.NewKeptn(&approvalTriggeredCE, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	// now let's trigger the other sequence
	secondContext = triggerSequence(t, projectName, serviceName, "dev", "delivery")
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{scmodels.SequenceStartedState})

	// check if approval.triggered has been sent for sequence - now it should be available
	myTaskTriggeredEvent, err := GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	myTaskCE := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)
	secondKeptnHandler, err := keptnv2.NewKeptn(&myTaskCE, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	// send the mytask.started event
	_, err = secondKeptnHandler.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// now let's send the approval.finished event - the next task should now be queued until the other sequence has been finished
	_, err = keptnHandler.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultPass, Status: keptnv2.StatusSucceeded}, source)
	require.Nil(t, err)

	// wait a bit to make sure mytask.triggered of first sequence is not sent
	<-time.After(10 * time.Second)
	myTaskTriggeredEventOfFirstSequence, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.Nil(t, myTaskTriggeredEventOfFirstSequence)

	// now let's finish mytask of the second sequence
	_, err = secondKeptnHandler.SendTaskFinishedEvent(&keptnv2.EventData{Status: keptnv2.StatusSucceeded, Result: keptnv2.ResultPass}, source)
	require.Nil(t, err)

	// this should have completed the task sequence
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{scmodels.SequenceFinished})

	// now the mytask.triggered event for the second sequence should eventually become available
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
		if err != nil || event == nil {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)

	// ----------------------------
	// Scenario 4: start a couple of task sequences and verify their completion
	// ----------------------------

	nrOfSequences := 10
	var wg sync.WaitGroup
	wg.Add(nrOfSequences)

	for i := 0; i < nrOfSequences; i++ {
		go executeSequenceAndVerifyCompletion(t, projectName, serviceName, "qg", &wg, []string{scmodels.SequenceFinished})
	}
	wg.Wait()
}

func Test_SequenceControl_TriggerMultiple(t *testing.T) {
	projectName := "sequence-queue2"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceAbortShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	for i := 0; i < 10; i++ {
		t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
		_, _ = TriggerSequence(projectName, serviceName, stageName, sequencename, nil)
	}

	openTriggeredEvents := &OpenTriggeredEventsResponse{}
	require.Eventually(t, func() bool {
		resp, err := ApiGETRequest("/controlPlane/v1/event/triggered/" + keptnv2.GetTriggeredEventType("task1") + "?project=" + projectName)
		if err != nil {
			return false
		}

		err = resp.ToJSON(openTriggeredEvents)
		if err != nil {
			return false
		}
		// must be exactly one .triggered event
		return len(openTriggeredEvents.Events) == 1
	}, 20*time.Second, 2*time.Second)
	require.Nil(t, err)
}

func executeSequenceAndVerifyCompletion(t *testing.T, projectName, serviceName, stageName string, wg *sync.WaitGroup, allowedStates []string) {
	defer wg.Done()
	context := triggerSequence(t, projectName, serviceName, stageName, "evaluation")
	source := "golang-test"

	var taskTriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(*context.KeptnContext, projectName, "qg", keptnv2.GetTriggeredEventType("mytask"))
		if err != nil || event == nil {
			return false
		}
		taskTriggeredEvent = event
		return true
	}, 5*time.Minute, 10*time.Second)
	require.NotNil(t, taskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	// send started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// send finished event with result=fail
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultPass,
	}, source)
	require.Nil(t, err)
	VerifySequenceEndsUpInState(t, projectName, context, 5*time.Minute, allowedStates)
	t.Logf("Sequence %s has been finished!", *context.KeptnContext)
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
				Result:  keptnv2.ResultPass,
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
