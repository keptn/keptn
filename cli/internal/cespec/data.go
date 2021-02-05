package cespec

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"strings"
)

func ce(ceType string, data interface{}) *cloudevents.Event {
	ce := cloudevents.NewEvent()
	ce.SetID("c4d3a334-6cb9-4e8c-a372-7e0b45942f53")
	ce.SetType(ceType)
	ce.SetSource("source-service")
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	ce.SetData(cloudevents.ApplicationJSON, data)
	ce.SetExtension("shkeptncontext", "a3e5f16d-8888-4720-82c7-6995062905c1")
	if !strings.HasSuffix(ceType, ".triggered") {
		ce.SetExtension("triggeredid", "3f9640b6-1d2a-4f11-95f5-23259f1d82d6")
	}
	return &ce
}

var commonEventData = keptnv2.EventData{
	Project: "sockshop",
	Stage:   "dev",
	Service: "carts",
	Labels:  map[string]string{"label-key": "label-value"},
	Status:  keptnv2.StatusSucceeded,
	Result:  keptnv2.ResultPass,
	Message: "a message",
}

var projectCreateData = keptnv2.ProjectCreateData{
	ProjectName:  "sockshop",
	GitRemoteURL: "https://github.com/project/repository",
	Shipyard:     `c3RhZ2VzOg0KICAtIG5hbWU6ICJkZXYiDQogICAgZGVwbG95bWVudF9zdHJhdGVneTogImRpcmVjdCINCiAgICB0ZXN0X3N0cmF0ZWd5OiAiZnVuY3Rpb25hbCINCiAgLSBuYW1lOiAic3RhZ2luZyINCiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIg0KICAgIHRlc3Rfc3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSINCiAgLSBuYW1lOiAicHJvZHVjdGlvbiINCiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIg0KICAgIHJlbWVkaWF0aW9uX3N0cmF0ZWd5OiAiYXV0b21hdGVkIg0K`,
}

var projectCreateStartedEventData = keptnv2.ProjectCreateStartedEventData{
	EventData: commonEventData,
}

var projectCreateFinishedEventData = keptnv2.ProjectCreateFinishedEventData{
	EventData: commonEventData,
	CreatedProject: keptnv2.ProjectCreateData{
		ProjectName:  "sockshop",
		GitRemoteURL: "https://github.com/project/repository",
		Shipyard:     "c3RhZ2VzOg0KICAtIG5hbWU6ICJkZXYiDQogICAgZGVwbG95bWVudF9zdHJhdGVneTogImRpcmVjdCINCiAgICB0ZXN0X3N0cmF0ZWd5OiAiZnVuY3Rpb25hbCINCiAgLSBuYW1lOiAic3RhZ2luZyINCiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIg0KICAgIHRlc3Rfc3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSINCiAgLSBuYW1lOiAicHJvZHVjdGlvbiINCiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIg0KICAgIHJlbWVkaWF0aW9uX3N0cmF0ZWd5OiAiYXV0b21hdGVkIg0K",
	},
}

var serviceCreateStartedEventData = keptnv2.ServiceCreateStartedEventData{
	EventData: commonEventData,
}

var serviceCreateStatusChangesData = keptnv2.ServiceCreateStatusChangedEventData{
	EventData: commonEventData,
}

var serviceCreateFinishedEventData = keptnv2.ServiceCreateFinishedEventData{
	EventData: commonEventData,
	Helm: keptnv2.Helm{
		Chart: `c3RhZ2VzOg0KICAtIG5hbWU6ICJkZXYiDQogICAgZGVwbG95bWVudF9zdHJhdGVneTogImRpcmVjdCINCiAgICB0ZXN0X3N0cmF0ZWd5OiAiZnVuY3Rpb25hbCINCiAgLSBuYW1lOiAic3RhZ2luZyINCiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIg0KICAgIHRlc3Rfc3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSINCiAgLSBuYW1lOiAicHJvZHVjdGlvbiINCiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIg0KICAgIHJlbWVkaWF0aW9uX3N0cmF0ZWd5OiAiYXV0b21hdGVkIg0K"`,
	},
}

//var approvalTriggeredEventData = keptnv2.ApprovalTriggeredEventData{
//	EventData: commonEventData,
//	Approval: keptnv2.Approval{
//		Pass:    keptnv2.ApprovalAutomatic,
//		Warning: keptnv2.ApprovalManual,
//	},
//}

var approvalStartedEventData = keptnv2.ApprovalStartedEventData{
	EventData: commonEventData,
}

var approvalStatusChangedEventData = keptnv2.ApprovalStatusChangedEventData{
	EventData: commonEventData,
}

var approvalFinishedEventData = keptnv2.ApprovalFinishedEventData{
	EventData: commonEventData,
}

var deploymentTriggeredEventData = keptnv2.DeploymentTriggeredEventData{
	EventData: commonEventData,
	ConfigurationChange: keptnv2.ConfigurationChange{
		Values: map[string]interface{}{"key": "value"},
	},
	Deployment: keptnv2.DeploymentWithStrategy{
		DeploymentStrategy: keptn.Direct.String(),
	},
}

var deploymentStartedEventData = keptnv2.DeploymentStartedEventData{
	EventData: commonEventData,
}

var deploymentStatusChangedEventData = keptnv2.DeploymentStatusChangedEventData{
	EventData: commonEventData,
}

var deploymentFinishedEventData = keptnv2.DeploymentFinishedEventData{
	EventData: commonEventData,
	Deployment: keptnv2.DeploymentData{
		DeploymentStrategy:   keptn.Direct.String(),
		DeploymentURIsLocal:  []string{"http://carts.sockshop-staging.svc.cluster.local"},
		DeploymentURIsPublic: []string{"http://carts.sockshot.local:80"},
		DeploymentNames:      []string{"deployment-1"},
		GitCommit:            "ca82a6dff817gc66f44342007202690a93763949",
	},
}

//var testTriggeredEventData = keptnv2.TestTriggeredEventData{
//	EventData: commonEventData,
//	Test: keptnv2.TestTriggeredDetails{
//		TestStrategy: "functional",
//	},
//	Deployment: keptnv2.TestTriggeredDeploymentDetails{
//		DeploymentURIsLocal:  []string{"http://carts.sockshop-staging.svc.cluster.local"},
//		DeploymentURIsPublic: []string{"http://carts.sockshot.local:80"},
//	},
//}

var testStartedEventData = keptnv2.TestStartedEventData{
	EventData: commonEventData,
}

var testStatusChangedEventData = keptnv2.TestStatusChangedEventData{
	EventData: commonEventData,
}

var testTestFinishedEventData = keptnv2.TestFinishedEventData{
	EventData: commonEventData,
	Test: struct {
		Start     string `json:"start"`
		End       string `json:"end"`
		GitCommit string `json:"gitCommit"`
	}{
		Start:     "2019-10-20T07:57:27.152330783Z",
		End:       "2019-10-20T08:57:27.152330783Z",
		GitCommit: "ca82a6dff817gc66f44342007202690a93763949",
	},
}

var evaluationTriggeredEventData = keptnv2.EvaluationTriggeredEventData{
	EventData: commonEventData,
	Test: struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Start: "2019-10-20T06:57:27.152330783Z",
		End:   "2019-10-20T07:57:27.152330783Z",
	},
	Evaluation: struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Start: "2019-10-20T07:57:27.152330783Z",
		End:   "2019-10-20T08:57:27.152330783Z",
	},
	Deployment: struct {
		DeploymentNames []string `json:"deploymentNames"`
	}{
		DeploymentNames: []string{"deployment-1"},
	},
}

var evaluationStartedEventData = keptnv2.EvaluationStartedEventData{
	EventData: commonEventData,
}

var evaluationStatusChangedEventData = keptnv2.EvaluationStatusChangedEventData{
	EventData: commonEventData,
}

var evaluationFinishedEventData = keptnv2.EvaluationFinishedEventData{
	EventData: commonEventData,
	Evaluation: keptnv2.EvaluationDetails{
		TimeStart:      "2019-11-18T11:21:06Z",
		TimeEnd:        "2019-11-18T11:29:36Z",
		Result:         "fail",
		Score:          0,
		SLOFileContent: "LS0tDQpzcGVjX3ZlcnNpb246ICcxLjAnDQpjb21wYXJpc29uOg0KICBjb21wYXJlX3dpdGg6ICJzaW5nbGVfcmVzdWx0Ig0KICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyINCiAgYWdncmVnYXRlX2Z1bmN0aW9uOiBhdmcNCm9iamVjdGl2ZXM6DQogIC0gc2xpOiByZXNwb25zZV90aW1lX3A5NQ0KICAgIHBhc3M6ICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNTAwKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcg0KICAgIHdhcm5pbmc6ICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogNzUl",
		IndicatorResults: []*keptnv2.SLIEvaluationResult{&keptnv2.SLIEvaluationResult{
			Score: 0,
			Value: &keptnv2.SLIResult{
				Metric:  "response_time_p95",
				Value:   1002.6278552658177,
				Success: true,
				Message: "a message",
			},
			Targets: []*keptnv2.SLITarget{&keptnv2.SLITarget{
				Criteria:    "<=+10%",
				TargetValue: 600,
				Violated:    true,
			}},
			Status: "failed",
		},
		},
		ComparedEvents: []string{"event-id-1", "event-id-2"},
		GitCommit:      "",
	},
}

var releaseTriggeredEventData = keptnv2.ReleaseTriggeredEventData{
	EventData: commonEventData,
	Deployment: keptnv2.DeploymentWithStrategy{
		DeploymentStrategy: keptn.Duplicate.String(),
	},
}

var releaseStartedEventData = keptnv2.ReleaseStartedEventData{
	EventData: commonEventData,
}

var releaseStatusChangedEventData = keptnv2.ReleaseStatusChangedEventData{
	EventData: commonEventData,
}

var releaseFinishedEventData = keptnv2.ReleaseFinishedEventData{
	EventData: commonEventData,
	Release: keptnv2.ReleaseData{
		GitCommit: "ca82a6dff817gc66f44342007202690a93763949",
	},
}

var remediationTriggeredEventData = keptnv2.RemediationTriggeredEventData{
	EventData: commonEventData,
	Problem: keptnv2.ProblemDetails{
		State:          "OPEN",
		ProblemID:      "ab81-941c-f198",
		ProblemTitle:   "Response time degradation",
		ProblemDetails: json.RawMessage{},
		PID:            "P23",
		ProblemURL:     "https://.../#problems/problemdetails;pid=93a5-3fas-a09d-8ckf",
		ImpactedEntity: "carts-primary",
		Tags:           "tag",
	},
}

var remediationStartedEventData = keptnv2.RemediationStartedEventData{
	EventData: commonEventData,
}

var remediationStatusChangedEventData = keptnv2.RemediationStatusChangedEventData{
	EventData: commonEventData,
	Remediation: keptnv2.Remediation{
		ActionIndex: 1,
		ActionName:  "trigger-runbook",
	},
}

var remediationFinishedEventData = keptnv2.RemediationFinishedEventData{
	EventData: commonEventData,
}

var actionTriggeredEventData = keptnv2.ActionTriggeredEventData{
	EventData: keptnv2.EventData{
		Project: "sockshop",
		Service: "carts",
		Stage:   "dev",
	},
	Action: keptnv2.ActionInfo{
		Name:        "Feature toggeling",
		Action:      "toggle-feature",
		Description: "Toggles a feature flag",
		Value:       map[string]string{"EnableItemCache": "on"},
	},
	Problem: keptnv2.ProblemDetails{
		State:          "OPEN",
		ProblemID:      "762",
		ProblemTitle:   "cpu_usage_sockshop_carts",
		PID:            "93a5-3fas-a09d-8ckf",
		ProblemURL:     "http://problem.url.com",
		ImpactedEntity: "carts-primary",
	},
}

var actionStartedEventData = keptnv2.ActionStartedEventData{
	EventData: commonEventData,
}

var actionFinishedEventData = keptnv2.ActionFinishedEventData{
	EventData: keptnv2.EventData{
		Project: "sockshop",
		Service: "carts",
		Stage:   "dev",
	},
	Action: keptnv2.ActionData{
		GitCommit: "93a5-3fas-a09d-8ckf",
	},
}

//var getSLITriggeredEventData = keptnv2.GetSLITriggeredEventData{
//	EventData: keptnv2.EventData{
//		Project: "sockshop",
//		Stage:   "dev",
//		Service: "carts",
//	},
//	GetSLI: keptnv2.GetSLI{
//		SLIProvider: "dynatrace",
//		Start:       "2019-10-28T15:44:27.152330783Z",
//		End:         "2019-10-28T15:54:27.152330783Z",
//		Indicators:  []string{"throughput", "error_rate", "request_latency_p95"},
//		CustomFilters: []*keptnv2.SLIFilter{{
//			Key:   "dynatraceEntityName",
//			Value: "HealthCheckController",
//		}, {
//			Key:   "tags",
//			Value: "test-subject:true",
//		}},
//	},
//}

var getSLIStartedEventData = keptnv2.GetSLIStartedEventData{
	EventData: commonEventData,
}

var getSLIFinishedEventData = &keptnv2.GetSLIFinishedEventData{
	EventData: keptnv2.EventData{
		Project: "sockshop",
		Service: "carts",
		Stage:   "dev",
	},
	//GetSLI: keptnv2.GetSLIFinished{
	//	Start: "2019-10-20T07:57:27.152330783Z",
	//	End:   "2019-10-22T08:57:27.152330783Z",
	//	IndicatorValues: []*keptnv2.SLIResult{
	//		{
	//			Metric:  "response_time_p50",
	//			Value:   1011.0745528937252,
	//			Success: true,
	//			Message: "",
	//		},
	//	},
	//},
}

var configureMonitoringTriggeredEventData = keptnv2.ConfigureMonitoringTriggeredEventData{
	EventData: keptnv2.EventData{
		Project: "sockshop",
		Service: "carts",
		Stage:   "dev",
	},
	ConfigureMonitoring: keptnv2.ConfigureMonitoringTriggeredParams{
		Type: "dynatrace",
	},
}

var configureMonitoringStartedEventData = keptnv2.ConfigureMonitoringStartedEventData{
	EventData: commonEventData,
}

var configureMonitoringFinishedEventData = keptnv2.ConfigureMonitoringFinishedEventData{
	EventData: commonEventData,
}

var p = `{
      "displayName": "641",
      "endTime": -1,
      "hasRootCause": false,
      "id": "1234_5678V2",
      "impactLevel": "SERVICE",
      "severityLevel": "PERFORMANCE",
      "startTime": 1587624420000,
      "status": "OPEN"
    }`

var problemOpenEventData = keptn.ProblemEventData{
	State:          "OPEN",
	ProblemID:      "ab81-941c-f198",
	ProblemTitle:   "Response Time Degradation",
	ProblemDetails: json.RawMessage(p),
	PID:            "93a5-3fas-a09d-8ckf",
	ProblemURL:     "https://.../#problems/problemdetails;pid=93a5-3fas-a09d-8ckf",
	ImpactedEntity: "carts-primary",
	Project:        "sockshop",
	Stage:          "production",
	Service:        "service",
}
