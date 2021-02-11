package handler

import (
	"reflect"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const (
	shkeptncontext = "1234-4567"
	eventID        = "2345-6789"
)

func compareEventContext(in1 cloudevents.Event, in2 cloudevents.Event) bool {
	return in1.Type() == in2.Type() && in1.Source() == in2.Source() &&
		in1.DataContentType() == in2.DataContentType() && reflect.DeepEqual(in1.Extensions(), in2.Extensions()) &&
		reflect.DeepEqual(in1.Data(), in2.Data())
}

func getPtr(s string) *string {
	return &s
}

func getShipyardWithoutApproval() keptnevents.Shipyard {
	return keptnevents.Shipyard{
		Stages: []struct {
			Name                string                              `json:"name" yaml:"name"`
			DeploymentStrategy  string                              `json:"deployment_strategy" yaml:"deployment_strategy"`
			TestStrategy        string                              `json:"test_strategy,omitempty" yaml:"test_strategy"`
			RemediationStrategy string                              `json:"remediation_strategy,omitempty" yaml:"remediation_strategy"`
			ApprovalStrategy    *keptnevents.ApprovalStrategyStruct `json:"approval_strategy,omitempty" yaml:"approval_strategy"`
		}{
			{
				Name:                "dev",
				DeploymentStrategy:  "direct",
				TestStrategy:        "functional",
				RemediationStrategy: "",
				ApprovalStrategy:    nil,
			},
			{
				Name:                "hardening",
				DeploymentStrategy:  "blue_green_service",
				TestStrategy:        "performance",
				RemediationStrategy: "",
				ApprovalStrategy:    nil,
			},
			{
				Name:                "production",
				DeploymentStrategy:  "blue_green_service",
				TestStrategy:        "",
				RemediationStrategy: "",
				ApprovalStrategy:    nil,
			},
		},
	}
}

func getApprovalTriggeredTestData(result keptnv2.ResultType, passStrategy, warningStrategy string) keptnv2.ApprovalTriggeredEventData {
	return keptnv2.ApprovalTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Labels: map[string]string{
				"l1": "lValue",
			},
			Status:  keptnv2.StatusSucceeded,
			Result:  result,
			Message: "",
		},
		Approval: keptnv2.Approval{
			Pass:    passStrategy,
			Warning: warningStrategy,
		},
	}
}

func getApprovalFinishedTestData(result keptnv2.ResultType, status keptnv2.StatusType) keptnv2.ApprovalFinishedEventData {
	return keptnv2.ApprovalFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Labels: map[string]string{
				"l1": "lValue",
			},
			Status:  status,
			Result:  result,
			Message: "",
		},
	}
}

func getApprovalStartedTestData(status keptnv2.StatusType) keptnv2.ApprovalStartedEventData {
	return keptnv2.ApprovalStartedEventData{
		EventData: keptnv2.EventData{
			Project: "sockshop",
			Stage:   "production",
			Service: "carts",
			Labels: map[string]string{
				"l1": "lValue",
			},
			Status:  status,
			Message: "",
		},
	}
}
