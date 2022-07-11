package go_tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
)

const qualityGatesShipyard = `--- 
apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata: 
  name: shipyard-quality-gates
spec: 
  stages: 
    - 
      name: hardening`

const qualityGatesShortSLOFileContent = `---
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
  pass: "100%"
  warning: "65%"`

const qualityGatesSLOFileContent = `---
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

const invalidSLOFileContent = "invalid"

func Test_QualityGates(t *testing.T) {
	t.Parallel()
	projectName := "quality-gates"
	serviceName := "my-service"
	shipyardFilePath, err := CreateTmpShipyardFile(qualityGatesShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	source := "golang-test"

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Log("deleting lighthouse configmap from previous test run")
	_, _ = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", GetKeptnNameSpaceFromEnv(), projectName))

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Log("triggering evaluation for wrong project")
	cliResp, err := ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", "wrong-project", "hardening", serviceName))
	require.NotNil(t, err)
	require.Contains(t, cliResp, "project not found")

	t.Log("triggering evaluation for wrong stage")
	cliResp, err = ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", projectName, "wrong-stage", serviceName))
	require.NotNil(t, err)
	require.Contains(t, cliResp, "stage not found")

	t.Log("triggering evaluation for wrong service")
	cliResp, err = ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", projectName, "hardening", "wrong-service"))
	require.NotNil(t, err)
	require.Contains(t, cliResp, "service not found")

	t.Log("triggering evaluation for existing project/stage/service with no SLO file and no SLI provider")
	keptnContext, err := TriggerEvaluation(projectName, "hardening", serviceName)

	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	t.Log("waiting for hardening.evaluation.finished event...")
	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 2*time.Minute, 10*time.Second)

	t.Log("got hardening.evaluation.finished event")
	require.NotNil(t, evaluationFinishedEvent)
	require.Equal(t, "lighthouse-service", *evaluationFinishedEvent.Source)
	require.Equal(t, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), *evaluationFinishedEvent.Type)
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, "pass", evaluationFinishedPayload.Evaluation.Result)
	require.Equal(t, float64(0), evaluationFinishedPayload.Evaluation.Score)
	require.Equal(t, "", evaluationFinishedPayload.Evaluation.SLOFileContent)
	require.Equal(t, []*keptnv2.SLIEvaluationResult([]*keptnv2.SLIEvaluationResult(nil)), evaluationFinishedPayload.Evaluation.IndicatorResults)
	require.Equal(t, []string([]string(nil)), evaluationFinishedPayload.Evaluation.ComparedEvents)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultPass, evaluationFinishedPayload.EventData.Result)

	t.Log("hardening.evaluation.finished event is valid")

	//// now let's add an SLI provider
	t.Log("adding SLI provider")
	_, err = ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", GetKeptnNameSpaceFromEnv(), projectName))
	require.Nil(t, err)

	// ...and an SLO file - but an invalid one :(
	t.Log("adding invalid SLO file")
	sloFilePath, err := CreateTmpFile("slo-*.yaml", invalidSLOFileContent)
	require.Nil(t, err)
	defer os.Remove(sloFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=slo.yaml", projectName, "hardening", serviceName, sloFilePath))
	require.Nil(t, err)
	checkCommit := false
	if image, err := GetImageOfDeploymentContainer("resource-service", "resource-service"); err == nil && image != "" {
		checkCommit = true
	}

	t.Log("triggering the evaluation again")
	keptnContext, err = TriggerEvaluation(projectName, "hardening", serviceName)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	// wait for the evaluation.finished event to be available and evaluate it
	t.Log("waiting for hardening.evaluation.finished event...")
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)
	t.Log("got hardening.evaluation.finished event")

	require.NotNil(t, evaluationFinishedEvent)
	require.Equal(t, "lighthouse-service", *evaluationFinishedEvent.Source)
	require.Equal(t, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), *evaluationFinishedEvent.Type)
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, "fail", evaluationFinishedPayload.Evaluation.Result)
	require.Equal(t, float64(0), evaluationFinishedPayload.Evaluation.Score)
	require.Equal(t, "", evaluationFinishedPayload.Evaluation.SLOFileContent)
	require.Equal(t, []*keptnv2.SLIEvaluationResult([]*keptnv2.SLIEvaluationResult(nil)), evaluationFinishedPayload.Evaluation.IndicatorResults)
	require.Equal(t, []string([]string(nil)), evaluationFinishedPayload.Evaluation.ComparedEvents)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusErrored, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.EventData.Result)
	require.NotEmpty(t, evaluationFinishedPayload.Message)

	t.Log("hardening.evaluation.finished event is valid")

	t.Log("adding invalid SLO file")
	sloFilePath, err = CreateTmpFile("slo-*.yaml", qualityGatesSLOFileContent)
	require.Nil(t, err)
	defer os.Remove(sloFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=slo.yaml", projectName, "hardening", serviceName, sloFilePath))
	require.Nil(t, err)

	t.Log("triggering the evaluation again (this time valid)")
	keptnContext, evaluationFinishedEvent = PerformResourceServiceTest(t, projectName, serviceName, checkCommit)

	t.Log("got hardening.evaluation.finished event")
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultPass, evaluationFinishedPayload.EventData.Result)
	require.NotEmpty(t, evaluationFinishedPayload.Message)

	require.Equal(t, "pass", evaluationFinishedPayload.Evaluation.Result)
	require.Equal(t, float64(100), evaluationFinishedPayload.Evaluation.Score)

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

	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 1,
		Value: &keptnv2.SLIResult{
			Metric:  "throughput",
			Value:   200,
			Success: true,
			Message: "",
		},
		DisplayName: "",
		PassTargets: []*keptnv2.SLITarget{
			{
				Criteria:    "<=+100%",
				TargetValue: 0,
				Violated:    false,
			},
			{
				Criteria:    ">=-80%",
				TargetValue: 0,
				Violated:    false,
			},
		},
		WarningTargets: nil,
		Status:         "pass",
		KeySLI:         false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[1])

	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 0,
		Value: &keptnv2.SLIResult{
			Metric:        "error_rate",
			Value:         0,
			ComparedValue: 0,
			Success:       true,
			Message:       "",
		},
		DisplayName:    "",
		PassTargets:    nil,
		WarningTargets: nil,
		Status:         "info",
		KeySLI:         false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[2])

	t.Log("hardening.evaluation.finished event is valid")

	firstEvaluationFinishedID := evaluationFinishedEvent.ID

	t.Log("invalidate the previous evaluation")
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

	t.Log("triggering the evaluation again (this time valid)")
	// do another evaluation - the resulting .finished event should not contain the first .finished event (which has been invalidated) in the list of compared evaluation results
	keptnContext, evaluationFinishedEvent = PerformResourceServiceTest(t, projectName, serviceName, checkCommit)
	t.Log("got hardening.evaluation.finished event")

	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.NotContains(t, evaluationFinishedPayload.Evaluation.ComparedEvents, firstEvaluationFinishedID)
	require.Len(t, evaluationFinishedPayload.Evaluation.ComparedEvents, 1)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultPass, evaluationFinishedPayload.EventData.Result)
	require.NotEmpty(t, evaluationFinishedPayload.Message)

	require.Equal(t, "pass", evaluationFinishedPayload.Evaluation.Result)
	require.Equal(t, float64(100), evaluationFinishedPayload.Evaluation.Score)

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

	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 1,
		Value: &keptnv2.SLIResult{
			Metric:  "throughput",
			Value:   200,
			Success: true,
			Message: "",
		},
		DisplayName: "",
		PassTargets: []*keptnv2.SLITarget{
			{
				Criteria:    "<=+100%",
				TargetValue: 0,
				Violated:    false,
			},
			{
				Criteria:    ">=-80%",
				TargetValue: 0,
				Violated:    false,
			},
		},
		WarningTargets: nil,
		Status:         "pass",
		KeySLI:         false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[1])

	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 0,
		Value: &keptnv2.SLIResult{
			Metric:        "error_rate",
			Value:         0,
			ComparedValue: 0,
			Success:       true,
			Message:       "",
		},
		DisplayName:    "",
		PassTargets:    nil,
		WarningTargets: nil,
		Status:         "info",
		KeySLI:         false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[2])

	t.Log("hardening.evaluation.finished event is valid")

	secondEvaluationFinishedID := evaluationFinishedEvent.ID

	// do another evaluation - the resulting .finished event should contain the second .finished event in the list of compared evaluation results
	t.Log("triggering the evaluation again (this time it will be compared with the second evaluation)")
	keptnContext, evaluationFinishedEvent = PerformResourceServiceTest(t, projectName, serviceName, checkCommit)
	t.Log("got hardening.evaluation.finished event")

	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Contains(t, evaluationFinishedPayload.Evaluation.ComparedEvents, secondEvaluationFinishedID)
	require.Len(t, evaluationFinishedPayload.Evaluation.ComparedEvents, 1)
	require.Equal(t, projectName, evaluationFinishedPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationFinishedPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationFinishedPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusSucceeded, evaluationFinishedPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultPass, evaluationFinishedPayload.EventData.Result)
	require.NotEmpty(t, evaluationFinishedPayload.Message)

	require.Equal(t, "pass", evaluationFinishedPayload.Evaluation.Result)
	require.Equal(t, float64(100), evaluationFinishedPayload.Evaluation.Score)

	require.Len(t, evaluationFinishedPayload.Evaluation.IndicatorResults, 3)
	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 1,
		Value: &keptnv2.SLIResult{
			Metric:        "response_time_p95",
			Value:         200,
			ComparedValue: 200,
			Success:       true,
			Message:       "",
		},
		DisplayName: "",
		PassTargets: []*keptnv2.SLITarget{
			{
				Criteria:    "<=+75%",
				TargetValue: 350,
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
				TargetValue: 400,
				Violated:    false,
			},
		},
		Status: "pass",
		KeySLI: false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[0])

	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 1,
		Value: &keptnv2.SLIResult{
			Metric:        "throughput",
			Value:         200,
			ComparedValue: 200,
			Success:       true,
			Message:       "",
		},
		DisplayName: "",
		PassTargets: []*keptnv2.SLITarget{
			{
				Criteria:    "<=+100%",
				TargetValue: 400,
				Violated:    false,
			},
			{
				Criteria:    ">=-80%",
				TargetValue: 40,
				Violated:    false,
			},
		},
		WarningTargets: nil,
		Status:         "pass",
		KeySLI:         false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[1])

	require.Equal(t, &keptnv2.SLIEvaluationResult{
		Score: 0,
		Value: &keptnv2.SLIResult{
			Metric:        "error_rate",
			Value:         0,
			ComparedValue: 0,
			Success:       true,
			Message:       "",
		},
		DisplayName:    "",
		PassTargets:    nil,
		WarningTargets: nil,
		Status:         "info",
		KeySLI:         false,
	}, evaluationFinishedPayload.Evaluation.IndicatorResults[2])

	t.Log("hardening.evaluation.finished event is valid")

	t.Log("retrieving project data")
	project, err := GetProject(projectName)
	require.Nil(t, err)

	t.Log("testing the retrieved project data")
	require.NotEmpty(t, project.Stages)
	require.NotEmpty(t, project.Stages[0].Services)
	require.NotEmpty(t, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)])
	require.Equal(t, keptnContext, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)].KeptnContext)
	require.NotEmpty(t, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName)])
	require.Equal(t, keptnContext, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName)].KeptnContext)
	require.NotEmpty(t, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)])
	require.Equal(t, keptnContext, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)].KeptnContext)

}

// NOTE: test is not used due to missing logic to handle stuck sequences
func Test_QualityGates_NoSLIAnswer(t *testing.T) {
	projectName := "quality-gates-no-sli-answer"
	serviceName := "my-service"

	projectName, keptnContext, _ := qualityGatesGenericTestStart(t, projectName, serviceName)

	t.Logf("Sleeping for 5min...")
	time.Sleep(5 * time.Minute)
	t.Logf("Continue to work...")

	t.Log("Verify sequence ends up in timedOut state")
	sequenceStates, _, err := GetState(projectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &keptnContext}, 2*time.Minute, []string{"timedOut"})
}

// NOTE: test is not used due to missing logic to handle stuck sequences
func Test_QualityGates_SLIStartedEventSend(t *testing.T) {
	source := "golang-test"
	projectName := "quality-gates-no-finished"
	serviceName := "my-service"

	projectName, keptnContext, triggeredID := qualityGatesGenericTestStart(t, projectName, serviceName)

	t.Log("sending get-sli.started event")
	_, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 5min...")
	time.Sleep(5 * time.Minute)
	t.Logf("Continue to work...")

	t.Log("Verify sequence ends up in timedOut state")
	sequenceStates, _, err := GetState(projectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &keptnContext}, 2*time.Minute, []string{"timedOut"})
}

func Test_QualityGates_SLIWrongFinishedPayloadSend(t *testing.T) {
	t.Parallel()
	source := "golang-test"
	projectName := "quality-gates-invalid-finish"
	serviceName := "my-service"

	projectName, keptnContext, triggeredID := qualityGatesGenericTestStart(t, projectName, serviceName)

	t.Log("sending get-sli.started event")
	_, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 15 sec...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Log("sending invalid get-sli.finished event")
	_, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Labels:  nil,
				Status:  "some-status",
				Result:  "strange-one",
				Message: "",
			},
			GetSLI: keptnv2.GetSLIFinished{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 15 sec...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Log("Verify sequence ends up in finished state")
	sequenceStates, _, err := GetState(projectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &keptnContext}, 2*time.Minute, []string{models.SequenceFinished})
}

func Test_QualityGates_AbortedFinishedPayloadSend(t *testing.T) {
	t.Parallel()
	source := "golang-test"
	projectName := "quality-gates-aborted-finish"
	serviceName := "my-service"

	projectName, keptnContext, triggeredID := qualityGatesGenericTestStart(t, projectName, serviceName)

	t.Log("sending get-sli.started event")
	_, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 15 sec...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Log("sending invalid get-sli.finished event")
	_, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusAborted,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
			GetSLI: keptnv2.GetSLIFinished{
				End:   "2022-01-26T10:10:53.931Z",
				Start: "2022-01-26T10:05:53.931Z",
				IndicatorValues: []*keptnv2.SLIResult{
					{
						Metric:        "response_time_p95",
						Value:         200,
						ComparedValue: 0,
						Success:       true,
						Message:       "",
					},
					{
						Metric:        "throughput",
						Value:         200,
						Success:       true,
						ComparedValue: 0,
						Message:       "",
					},
					{
						Metric:        "error_rate",
						Value:         0,
						ComparedValue: 0,
						Success:       true,
						Message:       "",
					},
				},
			},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 15 sec...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Log("Verify sequence ends up in finished state")
	sequenceStates, _, err := GetState(projectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	require.Equal(t, "errored", sequenceStates.States[0].Stages[0].State)
	require.Equal(t, "fail", sequenceStates.States[0].Stages[0].LatestEvaluation.Result)
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &keptnContext}, 2*time.Minute, []string{models.SequenceFinished})
}

func Test_QualityGates_ErroredFinishedPayloadSend(t *testing.T) {
	t.Parallel()
	source := "golang-test"
	projectName := "quality-gates-errored-finish"
	serviceName := "my-service"

	projectName, keptnContext, triggeredID := qualityGatesGenericTestStart(t, projectName, serviceName)

	t.Log("sending get-sli.started event")
	_, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 15 sec...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Log("sending invalid get-sli.finished event")
	_, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusErrored,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
			GetSLI: keptnv2.GetSLIFinished{
				End:   "2022-01-26T10:10:53.931Z",
				Start: "2022-01-26T10:05:53.931Z",
				IndicatorValues: []*keptnv2.SLIResult{
					{
						Metric:        "response_time_p95",
						Value:         200,
						ComparedValue: 0,
						Success:       true,
						Message:       "",
					},
					{
						Metric:        "throughput",
						Value:         200,
						Success:       true,
						ComparedValue: 0,
						Message:       "",
					},
					{
						Metric:        "error_rate",
						Value:         0,
						ComparedValue: 0,
						Success:       true,
						Message:       "",
					},
				},
			},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        triggeredID,
		GitCommitID:        "",
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Logf("Sleeping for 15 sec...")
	time.Sleep(15 * time.Second)
	t.Logf("Continue to work...")

	t.Log("Verify sequence ends up in finished state")
	sequenceStates, _, err := GetState(projectName)
	require.Nil(t, err)
	require.NotEmpty(t, sequenceStates.States)
	require.Equal(t, "errored", sequenceStates.States[0].Stages[0].State)
	require.Equal(t, "fail", sequenceStates.States[0].Stages[0].LatestEvaluation.Result)
	VerifySequenceEndsUpInState(t, projectName, &models.EventContext{KeptnContext: &keptnContext}, 2*time.Minute, []string{models.SequenceFinished})
}

func qualityGatesGenericTestStart(t *testing.T, projectName string, serviceName string) (string, string, string) {
	shipyardFilePath, err := CreateTmpShipyardFile(commitIDShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)
	keptnContext := ""

	t.Log("deleting lighthouse configmap from previous test run")
	ExecuteCommandf("kubectl delete configmap -n %s lighthouse-config-keptn-%s", GetKeptnNameSpaceFromEnv(), projectName)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommandf("keptn create service %s --project=%s", serviceName, projectName)

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("adding an SLI provider")
	_, err = ExecuteCommandf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", GetKeptnNameSpaceFromEnv(), projectName)
	require.Nil(t, err)

	t.Log("storing SLO file")
	_ = storeWithCommit(t, projectName, "hardening", serviceName, qualityGatesShortSLOFileContent, "slo.yaml")

	t.Log("sent hardening.evaluation.triggered")
	_, err = ExecuteCommandf("keptn trigger evaluation --project=%s --stage=hardening --service=%s --start=2022-01-26T10:05:53.931Z --end=2022-01-26T10:10:53.931Z", projectName, serviceName)
	require.Nil(t, err)

	var evaluationTriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if hardening.evaluation.triggered event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetTriggeredEventType("hardening."+keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationTriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	keptnContext = evaluationTriggeredEvent.Shkeptncontext
	t.Logf("Shkeptncontext is %s", keptnContext)

	t.Log("got hardening.evaluation.triggered event")

	t.Log("validating hardening.evaluation.triggered event")
	evaluationTriggeredPayload := &keptnv2.EvaluationTriggeredEventData{}
	err = keptnv2.Decode(evaluationTriggeredEvent.Data, evaluationTriggeredPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationTriggeredEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationTriggeredPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationTriggeredPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationTriggeredPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), evaluationTriggeredPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), evaluationTriggeredPayload.EventData.Result)
	require.Empty(t, evaluationTriggeredPayload.EventData.Message)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.Start)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.End)

	t.Log("hardening.evaluation.triggered event is valid")

	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.triggered event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationTriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	t.Log("got evaluation.triggered event")

	t.Log("validating evaluation.triggered event")
	evaluationTriggeredPayload = &keptnv2.EvaluationTriggeredEventData{}
	err = keptnv2.Decode(evaluationTriggeredEvent.Data, evaluationTriggeredPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationTriggeredEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationTriggeredPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationTriggeredPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationTriggeredPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), evaluationTriggeredPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), evaluationTriggeredPayload.EventData.Result)
	require.Empty(t, evaluationTriggeredPayload.EventData.Message)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.Start)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.End)

	t.Log("evaluation.triggered event is valid")

	var getSLITriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if get-sli.triggered event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
		if err != nil || event == nil {
			return false
		}
		getSLITriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	t.Log("got get-sli.triggered event: ", getSLITriggeredEvent)

	t.Log("validating get-sli.triggered event")
	getSLIPayload := &keptnv2.GetSLITriggeredEventData{}
	err = keptnv2.Decode(getSLITriggeredEvent.Data, getSLIPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, getSLITriggeredEvent.Shkeptncontext)
	require.Equal(t, "my-sli-provider", getSLIPayload.GetSLI.SLIProvider)
	require.NotEmpty(t, getSLIPayload.GetSLI.Start)
	require.NotEmpty(t, getSLIPayload.GetSLI.End)
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "response_time_p95")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "throughput")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "error_rate")
	require.Equal(t, projectName, getSLIPayload.EventData.Project)
	require.Equal(t, "hardening", getSLIPayload.EventData.Stage)
	require.Equal(t, serviceName, getSLIPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), getSLIPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), getSLIPayload.EventData.Result)
	require.Empty(t, getSLIPayload.EventData.Message)

	return projectName, keptnContext, getSLITriggeredEvent.ID
}

func PerformResourceServiceTest(t *testing.T, projectName string, serviceName string, checkCommit bool) (string, *models.KeptnContextExtendedCE) {
	commitID := ""
	commitID1 := ""
	if checkCommit {
		commitID1 = storeWithCommit(t, projectName, "hardening",
			serviceName, qualityGatesShortSLOFileContent, "slo.yaml")

		commitID = storeWithCommit(t, projectName, "hardening",
			serviceName, qualityGatesSLOFileContent, "slo.yaml")
	}
	keptnContext := ""
	source := "golang-test"

	t.Log("sent hardening.evaluation.triggered event with commitid= ", commitID)
	resp, err := ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: keptnv2.EvaluationTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
			},
			Test: keptnv2.Test{},
			Evaluation: keptnv2.Evaluation{
				End:       "2022-01-26T10:10:53.931Z",
				Start:     "2022-01-26T10:05:53.931Z",
				Timeframe: "",
			},
			Deployment: keptnv2.Deployment{},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             strutils.Stringp("fakeshipyard"),
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        "",
		GitCommitID:        commitID,
		Type:               strutils.Stringp(keptnv2.GetTriggeredEventType("hardening." + keptnv2.EvaluationTaskName)),
	}, 3)

	require.Nil(t, err)
	body := resp.String()
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)
	require.NotEmpty(t, body)
	kc := struct {
		KeptnContext *string `json:"keptnContext"`
	}{}
	resp.ToJSON(&kc)
	keptnContext = *kc.KeptnContext

	var evaluationTriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if hardening.evaluation.triggered event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetTriggeredEventType("hardening."+keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationTriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	keptnContext = evaluationTriggeredEvent.Shkeptncontext
	t.Logf("Shkeptncontext is %s", keptnContext)

	t.Log("got hardening.evaluation.triggered event")

	t.Log("validating hardening.evaluation.triggered event")
	evaluationTriggeredPayload := &keptnv2.EvaluationTriggeredEventData{}
	err = keptnv2.Decode(evaluationTriggeredEvent.Data, evaluationTriggeredPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationTriggeredEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationTriggeredPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationTriggeredPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationTriggeredPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), evaluationTriggeredPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), evaluationTriggeredPayload.EventData.Result)
	require.Empty(t, evaluationTriggeredPayload.EventData.Message)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.Start)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.End)

	t.Log("hardening.evaluation.triggered event is valid")

	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.triggered event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationTriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	t.Log("got evaluation.triggered event")

	t.Log("validating evaluation.triggered event")
	evaluationTriggeredPayload = &keptnv2.EvaluationTriggeredEventData{}
	err = keptnv2.Decode(evaluationTriggeredEvent.Data, evaluationTriggeredPayload)
	require.Nil(t, err)
	require.Equal(t, keptnContext, evaluationTriggeredEvent.Shkeptncontext)
	require.Equal(t, projectName, evaluationTriggeredPayload.EventData.Project)
	require.Equal(t, "hardening", evaluationTriggeredPayload.EventData.Stage)
	require.Equal(t, serviceName, evaluationTriggeredPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), evaluationTriggeredPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), evaluationTriggeredPayload.EventData.Result)
	require.Empty(t, evaluationTriggeredPayload.EventData.Message)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.Start)
	require.NotEmpty(t, evaluationTriggeredPayload.Evaluation.End)

	t.Log("evaluation.triggered event is valid")

	var getSLITriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if ", keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName), "for context ", keptnContext, " event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName))
		if err != nil || event == nil {
			return false
		}
		getSLITriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)
	t.Log("got get-sli.triggered event, checking commitid")

	if checkCommit {
		require.Equal(t, commitID, getSLITriggeredEvent.GitCommitID)
	}

	t.Log("validating get-sli.triggered event")
	getSLIPayload := &keptnv2.GetSLITriggeredEventData{}
	err = keptnv2.Decode(getSLITriggeredEvent.Data, getSLIPayload)
	require.Nil(t, err)
	require.Equal(t, "my-sli-provider", getSLIPayload.GetSLI.SLIProvider)
	require.NotEmpty(t, getSLIPayload.GetSLI.Start)
	require.NotEmpty(t, getSLIPayload.GetSLI.End)
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "response_time_p95")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "throughput")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "error_rate")
	require.Equal(t, projectName, getSLIPayload.EventData.Project)
	require.Equal(t, "hardening", getSLIPayload.EventData.Stage)
	require.Equal(t, serviceName, getSLIPayload.EventData.Service)
	require.Equal(t, keptnv2.StatusType(""), getSLIPayload.EventData.Status)
	require.Equal(t, keptnv2.ResultType(""), getSLIPayload.EventData.Result)
	require.Empty(t, getSLIPayload.EventData.Message)

	t.Log("get-sli.triggered event is valid")

	t.Log("sending get-sli.started event")
	resp, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIStartedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
		},
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        getSLITriggeredEvent.ID,
		GitCommitID:        commitID1,
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)

	t.Log("sending get-sli.finished event (valid)")
	resp, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype: "application/json",
		Data: &keptnv2.GetSLIFinishedEventData{
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   "hardening",
				Service: serviceName,
				Labels:  nil,
				Status:  keptnv2.StatusSucceeded,
				Result:  keptnv2.ResultPass,
				Message: "",
			},
			GetSLI: keptnv2.GetSLIFinished{
				Start: getSLIPayload.GetSLI.Start,
				End:   getSLIPayload.GetSLI.End,
				IndicatorValues: []*keptnv2.SLIResult{
					{
						Metric:        "response_time_p95",
						Value:         200,
						ComparedValue: 0,
						Success:       true,
						Message:       "",
					},
					{
						Metric:        "throughput",
						Value:         200,
						Success:       true,
						ComparedValue: 0,
						Message:       "",
					},
					{
						Metric:        "error_rate",
						Value:         0,
						ComparedValue: 0,
						Success:       true,
						Message:       "",
					},
				},
			},
		},
		Extensions:         nil,
		ID:                 uuid.NewString(),
		Shkeptncontext:     keptnContext,
		Shkeptnspecversion: KeptnSpecVersion,
		Source:             &source,
		Specversion:        "1.0",
		Time:               time.Now(),
		Triggeredid:        getSLITriggeredEvent.ID,
		GitCommitID:        commitID1,
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}, 3)
	require.Nil(t, err)

	// wait for the evaluation.finished event to be available and evaluate it
	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if hardening.evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)
	if checkCommit {
		//lighthouse should have used the new commitID to calculate the final result
		require.Equal(t, evaluationFinishedEvent.GitCommitID, commitID1)
	}
	return keptnContext, evaluationFinishedEvent
}

func TriggerEvaluation(projectName, stageName, serviceName string) (string, error) {
	cliResp, err := ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", projectName, stageName, serviceName))

	if err != nil {
		return "", err
	}
	var keptnContext string
	split := strings.Split(cliResp, "\n")
	for _, line := range split {
		if strings.Contains(line, "ID of") {
			splitLine := strings.Split(line, ":")
			if len(splitLine) == 2 {
				keptnContext = strings.TrimSpace(splitLine[1])
			}
		}
	}
	return keptnContext, err
}
