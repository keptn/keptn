package handler

import (
	"reflect"

	keptnevents "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go"
)

const shkeptncontext = "1234-4567"
const eventID = "2345-6789"

func compareEventContext(in1 cloudevents.Event, in2 cloudevents.Event) bool {

	return in1.Type() == in2.Type() && in1.Source() == in2.Source() &&
		in1.DataContentType() == in2.DataContentType() && reflect.DeepEqual(in1.Extensions(), in2.Extensions()) &&
		reflect.DeepEqual(in1.Data, in2.Data)
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

func getShipyardWithApproval(approvalStrategyForPass keptnevents.ApprovalStrategy,
	approvalStrategyForWarning keptnevents.ApprovalStrategy) keptnevents.Shipyard {

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
				Name: "production",
				ApprovalStrategy: &keptnevents.ApprovalStrategyStruct{
					Pass:    approvalStrategyForPass,
					Warning: approvalStrategyForWarning,
				},
				DeploymentStrategy:  "blue_green_service",
				TestStrategy:        "",
				RemediationStrategy: "",
			},
		},
	}
}

func getShipyardWithoutDeploymentStrategy(approvalStrategyForPass keptnevents.ApprovalStrategy,
	approvalStrategyForWarning keptnevents.ApprovalStrategy) keptnevents.Shipyard {

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
				ApprovalStrategy: &keptnevents.ApprovalStrategyStruct{
					Pass:    approvalStrategyForPass,
					Warning: approvalStrategyForWarning,
				},
			},
			{
				Name:                "production",
				DeploymentStrategy:  "",
				TestStrategy:        "",
				RemediationStrategy: "",
				ApprovalStrategy:    nil,
			},
		},
	}
}

func getApprovalTriggeredTestData(evaluationResult string) keptnevents.ApprovalTriggeredEventData {

	return keptnevents.ApprovalTriggeredEventData{
		Project:            "sockshop",
		Service:            "carts",
		Stage:              "production",
		TestStrategy:       getPtr("performance"),
		DeploymentStrategy: getPtr("blue_green_service"),
		Tag:                "0.11.1",
		Image:              "docker.io/keptnexamples/carts",
		Labels: map[string]string{
			"l1": "lValue",
		},
		Result: evaluationResult,
	}
}

func getApprovalFinishedTestData(result, status string) keptnevents.ApprovalFinishedEventData {

	return keptnevents.ApprovalFinishedEventData{
		Project:            "sockshop",
		Service:            "carts",
		Stage:              "production",
		TestStrategy:       getPtr("performance"),
		DeploymentStrategy: getPtr("blue_green_service"),
		Tag:                "0.11.1",
		Image:              "docker.io/keptnexamples/carts",
		Labels: map[string]string{
			"l1": "lValue",
		},
		Approval: keptnevents.ApprovalData{
			Result: result,
			Status: status,
		},
	}
}
