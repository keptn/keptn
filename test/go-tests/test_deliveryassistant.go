package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

const deliveryAssistantShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "combi1"
      sequences:
        - name: "approval"
          tasks:
            - name: "approval"
              properties:
                pass: "automatic"
                warning: "automatic"

    - name: "combi2"
      sequences:
        - name: "approval"
          tasks:
            - name: "approval"
              properties:
                pass: "manual"
                warning: "automatic"

    - name: "combi3"
      sequences:
        - name: "approval"
          tasks:
            - name: "approval"
              properties:
                pass: "automatic"
                warning: "manual"

    - name: "combi4"
      sequences:
        - name: "approval"
          tasks:
            - name: "approval"
              properties:
                pass: "manual"
                warning: "manual"`

func Test_DeliveryAssistant(t *testing.T) {
	projectName := "delivery-assistant"
	serviceName := "my-service"

	shipyardFilePath, err := CreateTmpShipyardFile(deliveryAssistantShipyard)
	require.Nil(t, err)
	defer func() {
		if err := os.Remove(shipyardFilePath); err != nil {
			t.Logf("warning: could not remove shipyard file %s", shipyardFilePath)
		}
	}()

	_, err = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", GetKeptnNameSpaceFromEnv(), projectName))
	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))
	require.Nil(t, err)

	// combi1
	combi1PassContext := triggerApproval(t, projectName, serviceName, "combi1", keptnv2.ResultPass)
	verifyApprovalFinishedEventExistsWithResult(t, combi1PassContext, projectName, "combi1", keptnv2.ResultPass, keptnv2.StatusSucceeded)

	combi1WarningContext := triggerApproval(t, projectName, serviceName, "combi1", keptnv2.ResultWarning)
	verifyApprovalFinishedEventExistsWithResult(t, combi1WarningContext, projectName, "combi1", keptnv2.ResultPass, keptnv2.StatusSucceeded)

	combi1FailContext := triggerApproval(t, projectName, serviceName, "combi1", keptnv2.ResultFailed)
	verifyApprovalFinishedEventExistsWithResult(t, combi1FailContext, projectName, "combi1", keptnv2.ResultFailed, keptnv2.StatusSucceeded)

	// also send an event with no result - in this case, only a .started event should be available
	combi1UnknownContext := triggerApproval(t, projectName, serviceName, "combi1", "")
	verifyApprovalStartedEventExists(t, combi1UnknownContext, projectName, "combi1", keptnv2.StatusSucceeded)
	verifyApprovalFinishedEventDoesNotExist(t, combi1UnknownContext, projectName)
	triggeredEvent := retrieveApprovalTriggeredEvent(t, combi1UnknownContext, projectName, "combi1")
	verifyApprovalUsingCLI(t, combi1UnknownContext, projectName, "combi1", triggeredEvent.ID)

	verifyNoOpenApprovalsLeft(t, projectName, "combi1")

	// combi2
	combi2PassContext := triggerApproval(t, projectName, serviceName, "combi2", keptnv2.ResultPass)
	verifyApprovalFinishedEventDoesNotExist(t, combi2PassContext, projectName)
	triggeredEvent = retrieveApprovalTriggeredEvent(t, combi2PassContext, projectName, "combi2")
	verifyApprovalUsingCLI(t, combi2PassContext, projectName, "combi2", triggeredEvent.ID)

	combi2WarningContext := triggerApproval(t, projectName, serviceName, "combi2", keptnv2.ResultWarning)
	verifyApprovalFinishedEventExistsWithResult(t, combi2WarningContext, projectName, "combi2", keptnv2.ResultPass, keptnv2.StatusSucceeded)

	combi2FailContext := triggerApproval(t, projectName, serviceName, "combi2", keptnv2.ResultFailed)
	verifyApprovalFinishedEventExistsWithResult(t, combi2FailContext, projectName, "combi2", keptnv2.ResultFailed, keptnv2.StatusSucceeded)

	verifyNoOpenApprovalsLeft(t, projectName, "combi2")

	// combi3
	combi3PassContext := triggerApproval(t, projectName, serviceName, "combi3", keptnv2.ResultPass)
	verifyApprovalFinishedEventExistsWithResult(t, combi3PassContext, projectName, "combi3", keptnv2.ResultPass, keptnv2.StatusSucceeded)

	combi3WarningContext := triggerApproval(t, projectName, serviceName, "combi3", keptnv2.ResultWarning)
	triggeredEvent = retrieveApprovalTriggeredEvent(t, combi3WarningContext, projectName, "combi3")
	verifyApprovalUsingCLI(t, combi3WarningContext, projectName, "combi3", triggeredEvent.ID)

	combi3FailContext := triggerApproval(t, projectName, serviceName, "combi3", keptnv2.ResultFailed)
	verifyApprovalFinishedEventExistsWithResult(t, combi3FailContext, projectName, "combi3", keptnv2.ResultFailed, keptnv2.StatusSucceeded)

	verifyNoOpenApprovalsLeft(t, projectName, "combi3")

	// combi4
	combi4PassContext := triggerApproval(t, projectName, serviceName, "combi4", keptnv2.ResultPass)

	verifyApprovalStartedEventExists(t, combi4PassContext, projectName, "combi4", keptnv2.StatusSucceeded)
	triggeredEvent = retrieveApprovalTriggeredEvent(t, combi4PassContext, projectName, "combi4")
	verifyApprovalUsingCLI(t, combi4PassContext, projectName, "combi4", triggeredEvent.ID)

	combi4WarningContext := triggerApproval(t, projectName, serviceName, "combi4", keptnv2.ResultWarning)
	verifyApprovalStartedEventExists(t, combi4PassContext, projectName, "combi4", keptnv2.StatusSucceeded)
	triggeredEvent = retrieveApprovalTriggeredEvent(t, combi4WarningContext, projectName, "combi4")
	verifyApprovalUsingCLI(t, combi4WarningContext, projectName, "combi4", triggeredEvent.ID)

	combi4FailContext := triggerApproval(t, projectName, serviceName, "combi4", keptnv2.ResultFailed)
	verifyApprovalStartedEventExists(t, combi4PassContext, projectName, "combi4", keptnv2.StatusSucceeded)
	verifyApprovalFinishedEventExistsWithResult(t, combi4FailContext, projectName, "combi4", keptnv2.ResultFailed, keptnv2.StatusSucceeded)

	verifyNoOpenApprovalsLeft(t, projectName, "combi4")

}

func verifyNoOpenApprovalsLeft(t *testing.T, projectName, stage string) {
	require.Eventually(t, func() bool {
		t.Logf("checking if no open approvals for stage %s are present anymore", stage)
		approvalCLIOutput, err := ExecuteCommand(fmt.Sprintf("keptn get event approval.triggered --project=%s --stage=%s", projectName, stage))
		if err != nil || !strings.Contains(approvalCLIOutput, "No approval.triggered events have been found") {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)
}

func verifyApprovalUsingCLI(t *testing.T, keptnContext, project, stage, triggeredID string) {
	_, err := ExecuteCommand(fmt.Sprintf("keptn send event approval.finished --id=%s --project=%s --stage=%s", triggeredID, project, stage))
	require.Nil(t, err)
	verifyApprovalFinishedEventExistsWithResult(t, keptnContext, project, stage, keptnv2.ResultPass, keptnv2.StatusSucceeded)
}

func retrieveApprovalTriggeredEvent(t *testing.T, keptnContext, projectName, stage string) *models.KeptnContextExtendedCE {
	var triggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Logf("verifying that approval.triggered event for context %s does exist", keptnContext)
		approvalTriggered, err := GetLatestEventOfType(keptnContext, projectName, stage, keptnv2.GetTriggeredEventType("approval"))
		if err != nil || approvalTriggered == nil {
			return false
		}
		triggeredEvent = approvalTriggered
		return true
	}, 1*time.Minute, 10*time.Second)
	require.NotNil(t, triggeredEvent)
	return triggeredEvent
}

func verifyApprovalFinishedEventExistsWithResult(t *testing.T, keptnContext, projectName, stage string, result keptnv2.ResultType, status keptnv2.StatusType) {
	var finishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Logf("verifying that approval.finished event for context %s does exist", keptnContext)
		approvalFinished, err := GetLatestEventOfType(keptnContext, projectName, stage, keptnv2.GetFinishedEventType("approval"))
		if err != nil || approvalFinished == nil {
			return false
		}
		finishedEvent = approvalFinished
		return true
	}, 1*time.Minute, 10*time.Second)
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(finishedEvent.Data, eventData)
	require.Nil(t, err)
	require.Equal(t, result, eventData.Result)
	require.Equal(t, status, eventData.Status)
}

func verifyApprovalStartedEventExists(t *testing.T, keptnContext, projectName, stage string, status keptnv2.StatusType) {
	var startedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Logf("verifying that approval.finished event for context %s does exist", keptnContext)
		approvalStarted, err := GetLatestEventOfType(keptnContext, projectName, stage, keptnv2.GetStartedEventType("approval"))
		if err != nil || approvalStarted == nil {
			return false
		}
		startedEvent = approvalStarted
		return true
	}, 1*time.Minute, 10*time.Second)
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(startedEvent.Data, eventData)
	require.Nil(t, err)
	require.Equal(t, status, eventData.Status)
}

func verifyApprovalFinishedEventDoesNotExist(t *testing.T, keptnContext string, projectName string) {
	require.Eventually(t, func() bool {
		t.Logf("verifying that approval.finished event for context %s does not exist", keptnContext)
		approvalFinished, err := GetLatestEventOfType(keptnContext, projectName, "combi2", keptnv2.GetFinishedEventType("approval"))
		if err != nil || approvalFinished != nil {
			return false
		}
		return true
	}, 1*time.Minute, 10*time.Second)
}

func triggerApproval(t *testing.T, projectName, serviceName, stageName string, result keptnv2.ResultType) string {
	context, err := TriggerSequence(projectName, serviceName, stageName, "approval", &keptnv2.EventData{
		Result: result,
	})
	require.Nil(t, err)
	require.NotEmpty(t, context)
	return context
}
