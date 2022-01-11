package go_tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
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
	projectName := "quality-gates"
	serviceName := "my-service"
	shipyardFilePath, err := CreateTmpShipyardFile(qualityGatesShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	source := "golang-test"

	_, err = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", GetKeptnNameSpaceFromEnv(), projectName))
	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

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
	keptnContext, err := triggerEvaluation(projectName, "hardening", serviceName)

	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	require.NotNil(t, evaluationFinishedEvent)
	require.Equal(t, "lighthouse-service", *evaluationFinishedEvent.Source)
	require.Equal(t, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), *evaluationFinishedEvent.Type)
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Equal(t, "pass", string(evaluationFinishedPayload.Result))
	require.Equal(t, float64(0), evaluationFinishedPayload.Evaluation.Score)
	require.Equal(t, "", evaluationFinishedPayload.Evaluation.SLOFileContent)

	// now let's add an SLI provider
	_, err = ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", GetKeptnNameSpaceFromEnv(), projectName))
	require.Nil(t, err)

	// ...and an SLO file - but an invalid one :(
	sloFilePath, err := CreateTmpFile("slo-*.yaml", invalidSLOFileContent)
	require.Nil(t, err)
	defer os.Remove(sloFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=slo.yaml", projectName, "hardening", serviceName, sloFilePath))
	require.Nil(t, err)

	t.Log("triggering the evaluation again")
	keptnContext, err = triggerEvaluation(projectName, "hardening", serviceName)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	// wait for the evaluation.finished event to be available and evaluate it
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)

	require.Equal(t, keptnv2.ResultFailed, evaluationFinishedPayload.Result)
	require.NotEmpty(t, evaluationFinishedPayload.Message)

	// ...and an SLO file
	sloFilePath, err = CreateTmpFile("slo-*.yaml", qualityGatesSLOFileContent)
	require.Nil(t, err)
	defer os.Remove(sloFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=slo.yaml", projectName, "hardening", serviceName, sloFilePath))
	require.Nil(t, err)

	t.Log("triggering the evaluation again")
	keptnContext, evaluationFinishedEvent = performEvaluationSequence(t, projectName, serviceName)

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

	firstEvaluationFinishedID := evaluationFinishedEvent.ID

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

	// do another evaluation - the resulting .finished event should not contain the first .finished event (which has been invalidated) in the list of compared evaluation results
	keptnContext, evaluationFinishedEvent = performEvaluationSequence(t, projectName, serviceName)

	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.NotContains(t, evaluationFinishedPayload.Evaluation.ComparedEvents, firstEvaluationFinishedID)
	secondEvaluationFinishedID := evaluationFinishedEvent.ID

	// do another evaluation - the resulting .finished event should contain the second .finished event in the list of compared evaluation results
	keptnContext, evaluationFinishedEvent = performEvaluationSequence(t, projectName, serviceName)

	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)
	require.Contains(t, evaluationFinishedPayload.Evaluation.ComparedEvents, secondEvaluationFinishedID)

	project, err := GetProject(projectName)
	require.Nil(t, err)

	require.NotEmpty(t, project.Stages)
	require.NotEmpty(t, project.Stages[0].Services)
	require.NotEmpty(t, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)])
	require.Equal(t, keptnContext, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)].KeptnContext)
	require.NotEmpty(t, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName)])
	require.Equal(t, keptnContext, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName)].KeptnContext)
	require.NotEmpty(t, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)])
	require.Equal(t, keptnContext, project.Stages[0].Services[0].LastEventTypes[keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)].KeptnContext)
}

func performEvaluationSequence(t *testing.T, projectName string, serviceName string) (string, *models.KeptnContextExtendedCE) {
	keptnContext, err := triggerEvaluation(projectName, "hardening", serviceName)
	source := "golang-test"
	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

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

	getSLIPayload := &keptnv2.GetSLITriggeredEventData{}
	err = keptnv2.Decode(getSLITriggeredEvent.Data, getSLIPayload)
	require.Nil(t, err)
	require.Equal(t, "my-sli-provider", getSLIPayload.GetSLI.SLIProvider)
	require.NotEmpty(t, getSLIPayload.GetSLI.Start)
	require.NotEmpty(t, getSLIPayload.GetSLI.End)

	cloudEvent := keptnv2.ToCloudEvent(*getSLITriggeredEvent)
	keptn, err := keptnv2.NewKeptn(&cloudEvent, keptncommon.KeptnOpts{EventSender: &APIEventSender{}})
	require.Nil(t, err)

	_, err = keptn.SendTaskStartedEvent(nil, source)
	require.Nil(t, err)

	// send finished event with result=fail
	_, err = keptn.SendTaskFinishedEvent(&keptnv2.GetSLIFinishedEventData{
		EventData: keptnv2.EventData{
			Status: keptnv2.StatusSucceeded,
			Result: keptnv2.ResultPass,
		},
		GetSLI: keptnv2.GetSLIFinished{
			Start: getSLIPayload.GetSLI.Start,
			End:   getSLIPayload.GetSLI.End,
			IndicatorValues: []*keptnv2.SLIResult{
				{
					Metric:  "response_time_p95",
					Value:   200,
					Success: true,
					Message: "",
				},
				{
					Metric:  "throughput",
					Value:   200,
					Success: true,
					Message: "",
				},
				{
					Metric:  "error_rate",
					Value:   0,
					Success: true,
					Message: "",
				},
			},
		},
	}, source)
	require.Nil(t, err)

	// wait for the evaluation.finished event to be available and evaluate it
	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, "hardening", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)
	return keptnContext, evaluationFinishedEvent
}

func triggerEvaluation(projectName, stageName, serviceName string) (string, error) {
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
