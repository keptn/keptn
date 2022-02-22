package go_tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
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

const commitIDDeliveryShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-podtato-ohead"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"
        - name: "delivery-direct"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"`

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

func createResourceWithCommit(t *testing.T, projectName, stage, serviceName, content, uri string) string {
	ctx, closeInternalKeptnAPI := context.WithCancel(context.Background())
	defer closeInternalKeptnAPI()
	internalKeptnAPI, err := GetInternalKeptnAPI(ctx, "service/configuration-service", "8889", "8080")
	require.Nil(t, err)
	resp, err := internalKeptnAPI.Put(basePath+"/"+projectName+"/stage/"+stage+"/service/"+serviceName+"/resource", models.Resources{
		Resources: []*models.Resource{
			{
				ResourceContent: content,
				ResourceURI:     strutils.Stringp(uri),
			},
		},
	}, 3)
	require.Nil(t, err)

	t.Logf("Received response %s", resp.String())
	require.Equal(t, 200, resp.Response().StatusCode)

	response := struct {
		CommitID string `json:"commitID"`
	}{}
	resp.ToJSON(&response)
	t.Log("Saved with commitID", response.CommitID)
	return response.CommitID
}

func Test_DeliveryGitCommitID(t *testing.T) {
	projectName := "commit-id-delivery6"
	serviceName := "helloservice"
	stageName := "dev"
	commitID := ""
	uri := "helm/helloservice.tgz"
	//newContent := "/Td6WFoAAATm1rRGAgAhARYAAAB0L+Wj4Bn/An9dADQZSe6OGhIDB9qF4IygZ+qL2Y/tbUkfy7ELAbc+GgCDeMGc4lucq0L6259TFNKlQbo5gkrbuoF2M40Ju7dl5hYdSZ+DQ4q1cs82f5YnJd2U26sLAa38OqBNm6+2Z7Gc+BUb03UkWqexMUxFzpcl2/fM4s1a+3Fo5609/ZDrg8vOD2OTUi0/4VHiGbiUgEvaAgpXGxoTmrD06LHtEnuL/EMjrAWN28xvHQXz5+Ztyp0h/tJNrp6Fk1qwHiXehrivpEVx1M+jzsgVPn0m0VLw37I0rRs87KCq/9g2ZycpPUjNwQl1H81uxu1cWqG1HWBoScv3IZf+tFfM7CscviDEFVO2MaZt7yNi3QOP0hDNv+BtVuSPgOrRWbeM21n6X/DTWXWgWq75LnGVOXDsdpzStat294+NCVTEx/reoyEDSBUW5PbECDDOgOyfovobDPEBbmDnka/ggWazn/P6Pe7UKbG1NaNJsrV5vJA1wqWmaU9oid57PmCGWDAcJp/AFIAihM6Lxa8TUh9dR91QfRMS1AsCIW3jyaefJdEoW13zNMJ1JSoqHJVHJPQgS7tDOVmf9lvr8KSJ61DVo6Wfo4Tz/yMiVcdxo1xLYc+24iS/w3RQa3Of1dUxm+S3IN+fRu15k+39vC2PEWBernL0sXxuZwJsULyU5jSKhLvQPgjoZuQyJ2NzLoO5CB9ahKUct7pH+X17L96e7Ul6KDdyQGOcFM1ZfJ1lKEndv5Fm/zhPN+t+m0PCdVvR4f70kAcqEHQNnpYh0sC+bJbJP1FSLVEWscMIux/5xHk2+Y/i6YAYcKSMqVY+Iw/nHkrRux9mpa31uLySSu1Xx6iCUx2udPOdAAAAQ/wZhn0S0ksAAZsFgDQAAFpGL7KxxGf7AgAAAAAEWVo="
	content := "H4sIAAAAAAAAA+2YzY7aMBDHOecp5gUIMQHSza3alVqpPSBttXeTTIlbJ45shwqt9t3rQMgmfIgeomVXnR8HlBl/jLHnPw7+JEMplUG9QT25z7i2/pbncjQggWOxmNXfLJoH3e8dLJyO2Mx9WBRGARsFbBYt2AiCIYO4RGUs1wAjlVYr/utyu2v+DwovxRNqI1QRw4Z5KZpEi9Lunj/DV5Q5JPWhgJ9Kw7dqhbpAi8YreI4xtGdHJOhtDuMEPvMD79YrI/4Fv5f/FvNScre/kyHnqHM8iuYX899xlP/zIFyMYD5kEJf4z/P/0v6nWEq1zbEYoBxc0X8WzsL+/k8DFoSk/29BV/95WZqJKwK/RZHG8NCeAC9Hy1NueewBnBF+U2JSu7TrIRJuYnh+Bv+JywqN3xjvVVVYeHlxzQxKTKzSdReAnNsk+85XKM3eAHUcRzMAHE5m06kTUI3s9T8/gpu4ibMmUYXlonBLP1jGzdL2qdAOJXK+xt6KdpbdSjotlpWUS+VWunV1U/7hW9P6S6VtJ7bx69xL54nhzuVA65VigwUas9Rqha+dADJryy9ouyY3NLdZDJO+7XRQF2MhrODyASXfPqILIHW7FHYalKiFSs+6rMhRVbb1zVufRp6Kdx/trVPsXXNJ/5u8GeRd4Nr9PwyP638Uzhek/2/B0f1/L/2PjWhe0f2+7p5q7kFvOwI4bhL+0yHdm2GdXDQGtxtrtCfKWGplVaJkDD/ulydF5HTuW/+uH4V+/m/2FW7gPwCuvv+f3v+nbEr5/xY015t1lmhfqEmpUsutGmeush8e9mcj3tSv9VOve5+Lgd06foIgCIIgCIIgCIIgCIIgCOI8fwFteWwzACgAAA=="

	t.Logf("Creating a new project %s with a Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(commitIDDeliveryShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", serviceName, projectName)
	require.Nil(t, err)

	//first part

	t.Logf("Adding resource for service %s in project %s", serviceName, projectName)
	commitID = createResourceWithCommit(t, projectName, stageName, serviceName, content, uri)
	t.Logf("commitID is %s", commitID)

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("./../../cli/cli trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s --git-commit-id=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery", commitID)
	require.Nil(t, err)

	// wait for the evaluation.triggered event to be available and check it
	var getDeliveryTriggeredEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.triggered event is available")
		event, err := GetLatestEventOfType("", projectName, stageName, keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		getDeliveryTriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	t.Log("got triggered event, checking commitID")
	keptnContext := getDeliveryTriggeredEvent.Shkeptncontext

	require.Equal(t, commitID, getDeliveryTriggeredEvent.GitCommitID)
	t.Log("commitID is present and correct")

	// wait for the evaluation.finished event to be available and check it
	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 2*time.Minute, 10*time.Second)

	t.Log("got finished event, checking commitid")
	require.Equal(t, commitID, evaluationFinishedEvent.GitCommitID)
	t.Log("commitID is present and correct")

	// second part

	t.Logf("Updating invalid resource for service %s in project %s", serviceName, projectName)
	commitID1 := createResourceWithCommit(t, projectName, stageName, serviceName, content, uri)
	t.Logf("new commitID is %s", commitID1)

	t.Logf("Trigger another delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("./../../cli/cli trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s --git-commit-id=%s", projectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery", commitID)
	require.Nil(t, err)

	// wait for the evaluation.triggered event to be available and check it
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.triggered event is available")
		event, err := GetLatestEventOfType("", projectName, stageName, keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		getDeliveryTriggeredEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)

	t.Log("got triggered event, checking commitid")
	keptnContext = getDeliveryTriggeredEvent.Shkeptncontext

	require.Equal(t, commitID, getDeliveryTriggeredEvent.GitCommitID)
	t.Log("commitID is present and correct")

	// wait for the evaluation.finished event to be available and check it
	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
		if err != nil || event == nil {
			return false
		}
		evaluationFinishedEvent = event
		return true
	}, 2*time.Minute, 10*time.Second)

	t.Log("got finished event, checking commitid")
	require.Equal(t, commitID, evaluationFinishedEvent.GitCommitID)
	t.Log("commitID is present and correct")

}

func Test_EvaluationGitCommitID(t *testing.T) {
	projectName := "commit-id-evaluation"
	serviceName := "my-service"
	shipyardFilePath, err := CreateTmpShipyardFile(commitIDShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Log("deleting lighthouse configmap from previous test run")
	_, _ = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-keptn-%s", GetKeptnNameSpaceFromEnv(), projectName))

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	t.Log("Testing the evaluation...")

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("adding an SLI provider")
	_, err = ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", GetKeptnNameSpaceFromEnv(), projectName))
	require.Nil(t, err)

	var evaluationFinishedEvent *models.KeptnContextExtendedCE
	evaluationFinishedPayload := &keptnv2.EvaluationFinishedEventData{}

	//first part

	t.Log("storing SLO file")
	commitID := storeWithCommit(t, projectName, "hardening", serviceName, gitCommitIDSLO, "slo.yaml")
	t.Logf("commitID is %s", commitID)

	t.Log("triggering the evaluation")
	_, evaluationFinishedEvent = performResourceServiceEvaluationTest(t, projectName, serviceName, commitID)

	t.Log("checking the finished event")
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)

	require.Equal(t, evaluationFinishedEvent.GitCommitID, commitID)
	t.Log("commitID is present and correct")

	//second part

	t.Log("storing second invalid SLO file")
	commitID1 := storeWithCommit(t, projectName, "hardening", serviceName, "gitCommitIDSLO", "slo.yaml")
	t.Logf("new commitID is %s", commitID1)

	t.Log("triggering the evaluation again with the old commitID")
	_, evaluationFinishedEvent = performResourceServiceEvaluationTest(t, projectName, serviceName, commitID)

	t.Log("checking the finished event again with the old commitID")
	err = keptnv2.Decode(evaluationFinishedEvent.Data, evaluationFinishedPayload)
	require.Nil(t, err)

	require.Equal(t, evaluationFinishedEvent.GitCommitID, commitID)
	t.Log("commitID is present and correct")
}

func performResourceServiceEvaluationTest(t *testing.T, projectName string, serviceName string, commitID string) (string, *models.KeptnContextExtendedCE) {
	keptnContext := ""
	source := "golang-test"

	t.Log("sent evaluation.hardening.triggered with commitid: ", commitID)
	_, err := ExecuteCommand(fmt.Sprintf("./../../cli/cli trigger evaluation --project=%s --stage=hardening --service=%s --start=2022-01-26T10:05:53.931Z --end=2022-01-26T10:10:53.931Z --git-commit-id=%s", projectName, serviceName, commitID))
	require.Nil(t, err)

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

	t.Log("got SLI triggered event, checking commitid")
	require.Equal(t, commitID, getSLITriggeredEvent.GitCommitID)
	t.Log("commitID is present and correct")

	getSLIPayload := &keptnv2.GetSLITriggeredEventData{}
	err = keptnv2.Decode(getSLITriggeredEvent.Data, getSLIPayload)
	require.Nil(t, err)
	require.Equal(t, "my-sli-provider", getSLIPayload.GetSLI.SLIProvider)
	require.NotEmpty(t, getSLIPayload.GetSLI.Start)
	require.NotEmpty(t, getSLIPayload.GetSLI.End)
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "response_time_p95")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "throughput")
	require.Contains(t, getSLIPayload.GetSLI.Indicators, "error_rate")

	//SLI uses a different commitID
	_, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
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
		GitCommitID:        commitID,
		Type:               strutils.Stringp(keptnv2.GetStartedEventType(keptnv2.GetSLITaskName)),
	}, 3)

	require.Nil(t, err)
	_, err = ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
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
		GitCommitID:        commitID,
		Type:               strutils.Stringp(keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)),
	}, 3)
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
