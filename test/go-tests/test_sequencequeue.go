package go_tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	//models "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

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

	//var currentActiveSequence models.SequenceState
	//for i := 0; i < numSequences; i++ {
	//	require.Eventually(t, func() bool {
	//		states, _, err := GetState(projectName)
	//		if err != nil {
	//			return false
	//		}
	//		for _, state := range states.States {
	//			if state.State == models.SequenceStartedState {
	//				// make sure the sequences are started in the chronologically correct order
	//				if sequenceContexts[i] != state.Shkeptncontext {
	//					return false
	//				}
	//				currentActiveSequence = state
	//				t.Logf("received expected active sequence: %s", state.Shkeptncontext)
	//				return true
	//			}
	//		}
	//		return false
	//	}, 15*time.Second, 2*time.Second)
	//
	//	_, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, currentActiveSequence.Shkeptncontext), models.SequenceControlCommand{
	//		State: models.AbortSequence,
	//		Stage: "",
	//	}, 3)
	//	require.Nil(t, err)
	//	if i == numSequences-1 {
	//		verifyNumberOfOpenTriggeredEvents(t, projectName, 0)
	//	} else {
	//		verifyNumberOfOpenTriggeredEvents(t, projectName, 1)
	//	}
	//}
	//
	//require.Nil(t, err)
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
