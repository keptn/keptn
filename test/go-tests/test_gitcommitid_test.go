package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

const commitIDShipyard = `--- 
apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata: 
  name: shipyard-quality-gates
spec: 
  stages: 
    - 
      name: hardening`

const gitCommitIDSLO = `---
spec_version: "0.1.1"
comparison:
  aggregate_function: "avg"
  compare_with: "single_result"
  include_result_with_score: "pass"
  number_of_comparison_results: 1
filter:
objectives:
  - sli: "response_time_p95"
    key_sli: false
    pass:             # pass if (relative change <= 75% AND absolute value is < 75ms)
      - criteria:
          - "<=+75%"  # relative values require a prefixed sign (plus or minus)
          - "<800"     # absolute values only require a logical operator
    warning:          # if the response time is below 200ms, the result should be a warning
      - criteria:
          - "<=1000"
          - "<=+100%"
    weight: 1
  - sli: "throughput"
    pass:
      - criteria:
          - "<=+100%"
          - ">=-80%"
  - sli: "error_rate"
total_score:
  pass: "90%"
  warning: "75%"`

func Test_GitCommitID(t *testing.T) {

	projectName := "commit-id"
	serviceName := "my-service"
	shipyardFilePath, err := CreateTmpShipyardFile(commitIDShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	source := "golang-test"

	t.Log("Checking for existence of configMap")
	resp, err := ExecuteCommand(fmt.Sprintf("kubectl get configmap -n %s lighthouse-config-keptn-%s", GetKeptnNameSpaceFromEnv(), projectName))
	if !strings.Contains(resp, "not found") {
		t.Log("ConfigMap exists, deleting...")
		_, err = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-keptn-%s", GetKeptnNameSpaceFromEnv(), projectName))
		require.Nil(t, err)
	}

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("adding an SLI provider")
	_, err = ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", GetKeptnNameSpaceFromEnv(), projectName))
	require.Nil(t, err)

	t.Logf("adding SLO file")
	sloFilePath, err := CreateTmpFile("slo-*.yaml", gitCommitIDSLO)
	require.Nil(t, err)
	defer os.Remove(sloFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --all-stages --service=%s --resource=%s --resourceUri=slo.yaml", projectName, serviceName, sloFilePath))
	require.Nil(t, err)

	// t.Log("triggering evaluation")
	// keptnContext, err := triggerEvaluation(projectName, "hardening", serviceName)
	// require.Nil(t, err)
	// require.NotEmpty(t, keptnContext)

	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	// // wait for the evaluation.finished event to be available and evaluate it
	// require.Eventually(t, func() bool {
	// 	t.Log("checking if evaluation.finished event is available")
	// 	event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
	// 	if err != nil || event == nil {
	// 		return false
	// 	}
	// 	evaluationFinishedEvent = event
	// 	return true
	// }, 1*time.Minute, 10*time.Second)

	// err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	// require.Nil(t, err)

	// require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.Result)
	// require.NotEmpty(t, evaluationFinishedPayload.Message)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=slo.yaml", projectName, "hardening", serviceName, sloFilePath))
	require.Nil(t, err)

	t.Log("triggering the evaluation again")
	_, evaluationFinishedEvent = performResourceServiceTest(t, projectName, serviceName, true)

	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)

	require.Len(t, evaluationFinishedPayload.Evaluation.IndicatorResults, 3)
	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 1,
		Value: &keptnv2.SLIResult{
			Metric:  "response_time_p95",
			Value:   200,
			Success: true,
			Message: "",
		},
		DisplayName: "",
		PassTargets: []*keptnv2.SLITarget{
			{
				Criteria:    "<=+75%",
				TargetValue: 0,
				Violated:    false,
			},
			{
				Criteria:    "<800",
				TargetValue: 800,
				Violated:    false,
			},
		},
		WarningTargets: []*keptnv2.SLITarget{
			{
				Criteria:    "<=1000",
				TargetValue: 1000,
				Violated:    false,
			},
			{
				Criteria:    "<=+100%",
				TargetValue: 0,
				Violated:    false,
			},
		},
		Status: "pass",
		KeySLI: false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[0])

	// send an evaluation.finished event for this evaluation
	evaluationInvalidatedEventType := "sh.keptn.event.evaluation.invalidated"
	_, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.EventData{
			Project: projectName,
			Stage:   "hardening",
			Service: serviceName,
		},
		ID:                 uuid.NewString(),
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Shkeptncontext:     evaluationFinishedEvent.Shkeptncontext,
		Triggeredid:        evaluationFinishedEvent.Triggeredid,
		Type:               &evaluationInvalidatedEventType,
	}, 3)
	require.Nil(t, err)

}
