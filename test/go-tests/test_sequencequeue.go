package go_tests

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	//models "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
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

const sequenceQueueShipyard2 = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "task1"
            - name: "task2"
            - name: "task3"`

func Test_SequenceQueue(t *testing.T) {
	projectName := "sequence-queue"
	serviceName := "my-service"

	sequenceStateShipyardFilePath, err := CreateTmpShipyardFile(sequenceQueueShipyard)
	require.Nil(t, err)
	defer os.Remove(sequenceStateShipyardFilePath)

	source := "golang-test"

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, sequenceStateShipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// ------------------------------------
	// Scenario 1: make sure sequences are queued correctly
	// ------------------------------------
	t.Logf("Scenario 1")
	t.Logf("trigger the task sequence")
	context := triggerSequence(t, projectName, serviceName, "dev", "delivery")

	t.Logf("wait for the sequence state to be available")
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{models.SequenceStartedState})
	t.Log("received the expected state!")

	t.Logf("trigger a second sequence - this one should stay in 'waiting' state until the previous sequence is finished")
	secondContext := triggerSequence(t, projectName, serviceName, "dev", "delivery")

	t.Logf("checking if the second sequence is in state 'waiting'")
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{models.SequenceWaitingState})
	t.Log("received the expected state!")

	t.Logf("check if mytask.triggered has been sent for first sequence - this one should be available")
	triggeredEventOfFirstSequence, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfFirstSequence)

	t.Logf("check if mytask.triggered has been sent for second sequence - this one should NOT be available")
	triggeredEventOfSecondSequence, err := GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.Nil(t, triggeredEventOfSecondSequence)

	t.Logf("send .started and .finished event for task of first sequence")
	cloudEvent := keptnv2.ToCloudEvent(*triggeredEventOfFirstSequence)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	t.Logf("send started event")
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	t.Logf("send finished event with result=fail")
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultFailed,
	}, source)
	require.Nil(t, err)

	t.Logf("now that all tasks for the first sequence have been executed, the second sequence should eventually have the status 'started'")
	t.Logf("waiting for state with keptnContext %s to have the status %s", *context.KeptnContext, models.SequenceStartedState)
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{models.SequenceStartedState})
	t.Log("received the expected state!")

	t.Logf("check if mytask.triggered has been sent for second sequence - now it should be available")
	triggeredEventOfSecondSequence, err = GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)

	// ------------------------------------
	// Scenario 2: test if sequences are triggered correctly after a timeout
	// ------------------------------------

	t.Logf("Scenario 2")
	err = setShipyardControllerTaskTimeout(t, "10s")
	defer func() {
		t.Logf("increase the timeout value again")
		err = setShipyardControllerTaskTimeout(t, "20m")
		require.Nil(t, err)
	}()
	require.Nil(t, err)

	t.Logf("wait a bit to make sure the .triggered event is received by the new instance of the shipyard controller")
	<-time.After(10 * time.Second)

	t.Logf("trigger the first task sequence - this should time out")
	context = triggerSequence(t, projectName, serviceName, "staging", "delivery")
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{models.TimedOut})
	t.Log("received the expected state!")

	t.Logf("now trigger the second sequence - this should start and a .triggered event for mytask should be sent")
	secondContext = triggerSequence(t, projectName, serviceName, "staging", "delivery")
	t.Logf("waiting for state with keptnContext %s to have the status %s", *secondContext.KeptnContext, models.SequenceStartedState)
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{models.SequenceStartedState})
	triggeredEventOfSecondSequence, err = GetLatestEventOfType(*secondContext.KeptnContext, projectName, "staging", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, triggeredEventOfSecondSequence)

	// ------------------------------------
	// Scenario 3: special case approval: once an approval has been triggered, another sequence should take over
	// ------------------------------------

	t.Logf("Scenario 3")
	require.Nil(t, err)
	t.Logf("increase the task timeout again")
	err = setShipyardControllerTaskTimeout(t, "20m")
	require.Nil(t, err)

	t.Logf("wait a bit to make sure the .triggered event is received by the new instance of the shipyard controller")
	<-time.After(10 * time.Second)

	t.Log("starting delivery-with-approval sequence")
	context = triggerSequence(t, projectName, serviceName, "dev", "delivery-with-approval")
	VerifySequenceEndsUpInState(t, projectName, context, 2*time.Minute, []string{models.SequenceStartedState})

	t.Logf("check if approval.triggered has been sent for sequence - now it should be available")
	approvalTriggeredEvent, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName))
	require.Nil(t, err)
	require.NotNil(t, approvalTriggeredEvent)

	t.Logf("send the approval.started event to make sure the sequence will not be cancelled due to a timeout")
	approvalTriggeredCE := keptnv2.ToCloudEvent(*approvalTriggeredEvent)
	keptnHandler, err := keptnv2.NewKeptn(&approvalTriggeredCE, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	t.Logf("now let's trigger the other sequence")
	secondContext = triggerSequence(t, projectName, serviceName, "dev", "delivery")
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{models.SequenceStartedState})

	t.Logf("check if approval.triggered has been sent for sequence - now it should be available")
	myTaskTriggeredEvent, err := GetLatestEventOfType(*secondContext.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.NotNil(t, myTaskTriggeredEvent)

	myTaskCE := keptnv2.ToCloudEvent(*myTaskTriggeredEvent)
	secondKeptnHandler, err := keptnv2.NewKeptn(&myTaskCE, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	t.Logf("send the mytask.started event")
	_, err = secondKeptnHandler.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	t.Logf("now let's send the approval.finished event - the next task should now be queued until the other sequence has been finished")
	_, err = keptnHandler.SendTaskFinishedEvent(&keptnv2.EventData{Result: keptnv2.ResultPass, Status: keptnv2.StatusSucceeded}, source)
	require.Nil(t, err)

	t.Logf("wait a bit to make sure mytask.triggered of first sequence is not sent")
	<-time.After(10 * time.Second)
	myTaskTriggeredEventOfFirstSequence, err := GetLatestEventOfType(*context.KeptnContext, projectName, "dev", keptnv2.GetTriggeredEventType("mytask"))
	require.Nil(t, err)
	require.Nil(t, myTaskTriggeredEventOfFirstSequence)

	t.Logf("now let's finish mytask of the second sequence")
	_, err = secondKeptnHandler.SendTaskFinishedEvent(&keptnv2.EventData{Status: keptnv2.StatusSucceeded, Result: keptnv2.ResultPass}, source)
	require.Nil(t, err)

	t.Logf("this should have completed the task sequence")
	VerifySequenceEndsUpInState(t, projectName, secondContext, 2*time.Minute, []string{models.SequenceFinished})

	t.Logf("now the mytask.triggered event for the second sequence should eventually become available")
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

	t.Logf("Scenario 4")
	nrOfSequences := 10
	var wg sync.WaitGroup
	wg.Add(nrOfSequences)

	for i := 0; i < nrOfSequences; i++ {
		go executeSequenceAndVerifyCompletion(t, projectName, serviceName, "qg", &wg, []string{models.SequenceFinished})
	}
	wg.Wait()
}

func Test_SequenceQueue_TriggerMultiple(t *testing.T) {
	projectName := "sequence-queue2"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"

	numSequences := 10

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceQueueShipyard2)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	sequenceContexts := []string{}
	for i := 0; i < numSequences; i++ {
		contextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)
		t.Logf("triggered sequence %s with context %s", sequencename, contextID)
		sequenceContexts = append(sequenceContexts, contextID)
		if i == 0 {
			// after triggering the first sequence, wait a few seconds to make sure this one is the first to be executed
			// all other sequences should be sorted correctly internally
			<-time.After(10 * time.Second)
		} else {
			<-time.After(2 * time.Second)
		}
	}
	verifyNumberOfOpenTriggeredEvents(t, projectName, 1)

	var currentActiveSequence models.SequenceState
	for i := 0; i < numSequences; i++ {
		require.Eventually(t, func() bool {
			states, _, err := GetState(projectName)
			if err != nil {
				return false
			}
			for _, state := range states.States {
				if state.State == models.SequenceStartedState {
					// make sure the sequences are started in the chronologically correct order
					if sequenceContexts[i] != state.Shkeptncontext {
						return false
					}
					currentActiveSequence = state
					t.Logf("received expected active sequence: %s", state.Shkeptncontext)
					return true
				}
			}
			return false
		}, 15*time.Second, 2*time.Second)

		_, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, currentActiveSequence.Shkeptncontext), models.SequenceControlCommand{
			State: models.AbortSequence,
			Stage: "",
		}, 3)
		require.Nil(t, err)
		if i == numSequences-1 {
			verifyNumberOfOpenTriggeredEvents(t, projectName, 0)
		} else {
			verifyNumberOfOpenTriggeredEvents(t, projectName, 1)
		}
	}

	require.Nil(t, err)
}

func Test_SequenceQueue_TriggerAndDeleteProject(t *testing.T) {
	projectName := "sequence-queue3ses"

	stageName := "dev"
	sequencename := "mysequence"

	numServices := 50
	numSequencesPerService := 1

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceQueueShipyard2)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	nProjectName, err := CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	for i := 0; i < numServices; i++ {
		serviceName := fmt.Sprintf("service-%d", i)
		output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, nProjectName))
		require.Nil(t, err)
		require.Contains(t, output, "created successfully")
	}

	triggerSequence := func(serviceName string, wg *sync.WaitGroup) {
		_, err := TriggerSequence(nProjectName, serviceName, stageName, sequencename, nil)
		require.Nil(t, err)
		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(numServices * numSequencesPerService)

	for i := 0; i < numServices; i++ {
		for j := 0; j < numSequencesPerService; j++ {
			serviceName := fmt.Sprintf("service-%d", i)
			go triggerSequence(serviceName, &wg)
		}
	}
	wg.Wait()

	// after all sequences have been triggered, delete the project
	//_, err = ExecuteCommand(fmt.Sprintf("keptn delete project %s", projectName))
	_, err = ApiDELETERequest("/controlPlane/v1/project/"+nProjectName, 3)

	require.Nil(t, err)

	// recreate the project again
	t.Logf("recreating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	// check if there are any open .triggered events for the project
	openTriggeredEvents := &OpenTriggeredEventsResponse{}
	resp, err := ApiGETRequest("/controlPlane/v1/event/triggered/"+keptnv2.GetTriggeredEventType("task1")+"?project="+projectName, 3)
	require.Nil(t, err)

	err = resp.ToJSON(openTriggeredEvents)
	require.Nil(t, err)

	require.Empty(t, openTriggeredEvents.Events)
}

func verifyNumberOfOpenTriggeredEvents(t *testing.T, projectName string, numberOfEvents int) {
	openTriggeredEvents := &OpenTriggeredEventsResponse{}
	require.Eventually(t, func() bool {
		resp, err := ApiGETRequest("/controlPlane/v1/event/triggered/"+keptnv2.GetTriggeredEventType("task1")+"?project="+projectName, 3)
		if err != nil {
			return false
		}

		err = resp.ToJSON(openTriggeredEvents)
		if err != nil {
			return false
		}
		// must be exactly one .triggered event
		t.Logf("received %d events, expected %d", len(openTriggeredEvents.Events), numberOfEvents)
		return len(openTriggeredEvents.Events) == numberOfEvents
	}, 20*time.Second, 2*time.Second)
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
	}, 15*time.Minute, 10*time.Second)
	require.NotNil(t, taskTriggeredEvent)

	cloudEvent := keptnv2.ToCloudEvent(*taskTriggeredEvent)

	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	// send started event
	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// verify that the started event has made its way to shippy
	require.Eventually(t, func() bool {
		state, _, err := GetStateByContext(projectName, *context.KeptnContext)
		if err != nil {
			return false
		}
		if len(state.States) == 0 || len(state.States[0].Stages) == 0 {
			return false
		}
		return state.States[0].Stages[0].LatestEvent.Type == keptnv2.GetStartedEventType("mytask")
	}, 1*time.Minute, 5*time.Second)

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
	}, 3)
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
