package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/osutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const selfHealingShipyard = `apiVersion: "spec.keptn.sh/0.2.2"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "production"
      sequences:
        - name: "remediation"
          triggeredOn:
            - event: "production.remediation.finished"
              selector:
                match:
                  evaluation.result: "fail"
          tasks:
            - name: "get-action"
            - name: "action"
            - name: "evaluation"
              triggeredAfter: "1m"`

const remediationFileContent = `apiVersion: spec.keptn.sh/0.1.4
kind: Remediation
metadata:
  name: service-remediation
spec:
  remediations:
    - problemType: Response time degradation
      actionsOnOpen:
      - action: toggle-feature
        name: toggle-feature
        description: Toggle feature flag EnablePromotion to OFF
        value:
          EnablePromotion: "off"`

const defaultUnleashServiceVersion = "master"
const unleashServiceEnvVar = "UNLEASH_SERVICE_VERSION"

type RemediationTriggered struct {
	keptnv2.EventData
	Problem keptnv2.ProblemDetails `json:"problem"`
}

func Test_SelfHealing(t *testing.T) {
	projectName := "self-healing"
	serviceName := "my-service"
	shipyardFilePath, err := CreateTmpShipyardFile(selfHealingShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", GetKeptnNameSpaceFromEnv(), projectName))
	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))
	require.Nil(t, err)

	// trigger a remediation - this should fail because no remediation.yaml is available yet
	t.Log("triggering a remediation with no remediation.yaml")
	remediationFinishedEvent := performRemediation(t, projectName, serviceName)

	require.Equal(t, "shipyard-controller", *remediationFinishedEvent.Source)
	finishedEventData := &keptnv2.EventData{}
	err = keptnv2.Decode(remediationFinishedEvent.Data, finishedEventData)
	require.Nil(t, err)
	require.Equal(t, "shipyard-controller", *remediationFinishedEvent.Source)
	require.Equal(t, keptnv2.StatusErrored, finishedEventData.Status)
	require.Equal(t, keptnv2.ResultFailed, finishedEventData.Result)

	t.Log("adding remediation.yaml file")
	remediationFilePath, err := CreateTmpFile("remediation-*.yaml", remediationFileContent)
	defer os.Remove(remediationFilePath)
	require.Nil(t, err)
	require.NotEmpty(t, remediationFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=remediation.yaml", projectName, "production", serviceName, remediationFilePath))
	require.Nil(t, err)

	t.Log("Installing unleash-service")
	unleashServiceVersion := osutils.GetOSEnvOrDefault(unleashServiceEnvVar, defaultUnleashServiceVersion)
	_, err = ExecuteCommand(fmt.Sprintf("kubectl apply -f \"https://raw.githubusercontent.com/keptn-contrib/unleash-service/%s/deploy/service.yaml\" -n %s", unleashServiceVersion, GetKeptnNameSpaceFromEnv()))
	require.Nil(t, err)

	err = WaitForPodOfDeployment("unleash-service")
	require.Nil(t, err)

	t.Log("remediation.yaml and unleash-service are ready. let's trigger another remediation")
	remediationFinishedEvent = performRemediation(t, projectName, serviceName)

	// inspect the remediation.finished event again
	finishedEventData = &keptnv2.EventData{}
	err = keptnv2.Decode(remediationFinishedEvent.Data, finishedEventData)
	require.Nil(t, err)
	require.Equal(t, "shipyard-controller", *remediationFinishedEvent.Source)
	require.Equal(t, keptnv2.StatusErrored, finishedEventData.Status)
	require.Equal(t, keptnv2.ResultFailed, finishedEventData.Result)

	t.Log("verifying if action.triggered event has been sent")
	var actionTriggered *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(remediationFinishedEvent.Shkeptncontext, projectName, "production", keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName))
		if err != nil || event == nil {
			return false
		}
		actionTriggered = event
		return true
	}, 1*time.Minute, 10*time.Second)

	triggeredEventData := &keptnv2.ActionTriggeredEventData{}
	err = keptnv2.Decode(actionTriggered.Data, triggeredEventData)
	require.Equal(t, keptnv2.ActionInfo{
		Name:        "toggle-feature",
		Action:      "toggle-feature",
		Description: "Toggle feature flag EnablePromotion to OFF",
		Value: map[string]interface{}{
			"EnablePromotion": "off",
		},
	}, triggeredEventData.Action)

	t.Log("verifying if action.finished event has been sent")
	var actionFinished *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(remediationFinishedEvent.Shkeptncontext, projectName, "production", keptnv2.GetFinishedEventType(keptnv2.ActionTaskName))
		if err != nil || event == nil {
			return false
		}
		actionFinished = event
		return true
	}, 1*time.Minute, 10*time.Second)

	finishedEventData = &keptnv2.EventData{}
	err = keptnv2.Decode(actionFinished.Data, finishedEventData)
	require.Equal(t, keptnv2.StatusErrored, finishedEventData.Status)
}

func performRemediation(t *testing.T, projectName string, serviceName string) *models.KeptnContextExtendedCE {
	keptnContext, err := TriggerSequence(projectName, serviceName, "production", "remediation", &RemediationTriggered{
		Problem: keptnv2.ProblemDetails{
			RootCause:    "Response time degradation",
			ProblemTitle: "My Problem",
		},
	})

	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	t.Log("waiting for remediation.finished event to be available")
	var remediationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(keptnContext, projectName, "production", keptnv2.GetFinishedEventType("production.remediation"))
		if err != nil || event == nil {
			return false
		}
		remediationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)
	return remediationFinishedEvent
}
