package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

const upgradeShipyardMockResponseContent = `stages:
  - name: "dev"
    deployment_strategy: "direct"
    test_strategy: "functional"
  - name: "staging"
    approval_strategy: 
      pass: "automatic"
      warning: "manual"
    deployment_strategy: "blue_green_service"
    test_strategy: "performance"
  - name: "production"
    approval_strategy: 
      pass: "automatic"
      warning: "manual"
    deployment_strategy: "blue_green_service"
    remediation_strategy: "automated"`

const upgradeShipyardResourceMockResponse = `{
      "resourceContent": "c3RhZ2VzOgogIC0gbmFtZTogImRldiIKICAgIGRlcGxveW1lbnRfc3RyYXRlZ3k6ICJkaXJlY3QiCiAgICB0ZXN0X3N0cmF0ZWd5OiAiZnVuY3Rpb25hbCIKICAtIG5hbWU6ICJzdGFnaW5nIgogICAgYXBwcm92YWxfc3RyYXRlZ3k6IAogICAgICBwYXNzOiAiYXV0b21hdGljIgogICAgICB3YXJuaW5nOiAibWFudWFsIgogICAgZGVwbG95bWVudF9zdHJhdGVneTogImJsdWVfZ3JlZW5fc2VydmljZSIKICAgIHRlc3Rfc3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSIKICAtIG5hbWU6ICJwcm9kdWN0aW9uIgogICAgYXBwcm92YWxfc3RyYXRlZ3k6IAogICAgICBwYXNzOiAiYXV0b21hdGljIgogICAgICB3YXJuaW5nOiAibWFudWFsIgogICAgZGVwbG95bWVudF9zdHJhdGVneTogImJsdWVfZ3JlZW5fc2VydmljZSIKICAgIHJlbWVkaWF0aW9uX3N0cmF0ZWd5OiAiYXV0b21hdGVkIg==",
      "resourceURI": "shipyard.yaml"
}`

func Test_UpgradeProjectShipyard(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	receivedUpgradedShipyard := make(chan bool)
	mocking = true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if r.Method == http.MethodGet && strings.Contains(r.RequestURI, "shipyard.yaml") {
				w.Write([]byte(upgradeShipyardResourceMockResponse))
				return
			} else if r.Method == http.MethodPut && strings.Contains(r.RequestURI, "resource") {
				defer r.Body.Close()
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("could not read received event payload: %s", err.Error())
				}
				resource := &apimodels.Resource{}
				if err := json.Unmarshal(bytes, resource); err != nil {
					t.Errorf("could not decode received resource: %s", err.Error())
				}
				if *resource.ResourceURI != "shipyard.yaml" {
					t.Errorf("did not receive upgraded shipyard: %s", err.Error())
				}
				v := &apimodels.Version{}
				marshal, _ := json.Marshal(v)
				w.Write(marshal)
				go func() {
					receivedUpgradedShipyard <- true
				}()
			}
			return
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	cmd := fmt.Sprintf("upgrade project sockshop --shipyard -y")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}

	select {
	case <-receivedUpgradedShipyard:
		t.Log("shipyard has been updated successfully")
		break
	case <-time.After(5 * time.Second):
		t.Error("shipyard has not been updated")
	}
}

func Test_transformShipyard(t *testing.T) {
	type args struct {
		shipyard *keptn.Shipyard
	}
	tests := []struct {
		name string
		args args
		want *keptnv2.Shipyard
	}{
		{
			name: "transform shipyard",
			args: args{
				shipyard: &keptn.Shipyard{
					Stages: getTestV1ShipyardStages(),
				},
			},
			want: &keptnv2.Shipyard{
				ApiVersion: "spec.keptn.sh/0.2.0",
				Kind:       "Shipyard",
				Metadata:   keptnv2.Metadata{},
				Spec: keptnv2.ShipyardSpec{
					Stages: []keptnv2.Stage{
						{
							Name: "dev",
							Sequences: []keptnv2.Sequence{
								{
									Name:     "artifact-delivery",
									Triggers: []string{},
									Tasks: []keptnv2.Task{
										{
											Name: "deployment",
											Properties: map[string]string{
												"deploymentstrategy": "direct",
											},
										},
										{
											Name: "test",
											Properties: map[string]string{
												"teststrategy": "functional",
											},
										},
										{
											Name: "evaluation",
										},
										{
											Name: "approval",
											Properties: map[string]string{
												"pass":    "automatic",
												"warning": "manual",
											},
										},
										{
											Name: "release",
										},
									},
								},
								{
									Name:     "artifact-delivery-direct",
									Triggers: []string{},
									Tasks: []keptnv2.Task{
										{
											Name: "deployment",
											Properties: map[string]string{
												"deploymentstrategy": "direct",
											},
										},
										{
											Name: "test",
											Properties: map[string]string{
												"teststrategy": "functional",
											},
										},
										{
											Name: "evaluation",
										},
										{
											Name: "approval",
											Properties: map[string]string{
												"pass":    "automatic",
												"warning": "manual",
											},
										},
										{
											Name: "release",
										},
									},
								},
							},
						},
						{
							Name: "staging",
							Sequences: []keptnv2.Sequence{
								{
									Name:     "artifact-delivery",
									Triggers: []string{"dev.artifact-delivery.finished"},
									Tasks: []keptnv2.Task{
										{
											Name: "deployment",
											Properties: map[string]string{
												"deploymentstrategy": "blue_green_service",
											},
										},
										{
											Name: "test",
											Properties: map[string]string{
												"teststrategy": "performance",
											},
										},
										{
											Name: "evaluation",
										},
										{
											Name: "approval",
											Properties: map[string]string{
												"pass":    "automatic",
												"warning": "manual",
											},
										},
										{
											Name: "release",
										},
									},
								},
								{
									Name:     "artifact-delivery-direct",
									Triggers: []string{"dev.artifact-delivery-direct.finished"},
									Tasks: []keptnv2.Task{
										{
											Name: "deployment",
											Properties: map[string]string{
												"deploymentstrategy": "direct",
											},
										},
										{
											Name: "test",
											Properties: map[string]string{
												"teststrategy": "performance",
											},
										},
										{
											Name: "evaluation",
										},
										{
											Name: "approval",
											Properties: map[string]string{
												"pass":    "automatic",
												"warning": "manual",
											},
										},
										{
											Name: "release",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformShipyard(tt.args.shipyard)

			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Error("returned value did not  match expected")
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}

func getTestV1ShipyardStages() []struct {
	Name                string                        `json:"name" yaml:"name"`
	DeploymentStrategy  string                        `json:"deployment_strategy" yaml:"deployment_strategy"`
	TestStrategy        string                        `json:"test_strategy,omitempty" yaml:"test_strategy"`
	RemediationStrategy string                        `json:"remediation_strategy,omitempty" yaml:"remediation_strategy"`
	ApprovalStrategy    *keptn.ApprovalStrategyStruct `json:"approval_strategy,omitempty" yaml:"approval_strategy"`
} {
	return []struct {
		Name                string                        `json:"name" yaml:"name"`
		DeploymentStrategy  string                        `json:"deployment_strategy" yaml:"deployment_strategy"`
		TestStrategy        string                        `json:"test_strategy,omitempty" yaml:"test_strategy"`
		RemediationStrategy string                        `json:"remediation_strategy,omitempty" yaml:"remediation_strategy"`
		ApprovalStrategy    *keptn.ApprovalStrategyStruct `json:"approval_strategy,omitempty" yaml:"approval_strategy"`
	}{
		{
			Name:               "dev",
			DeploymentStrategy: "direct",
			TestStrategy:       "functional",
			ApprovalStrategy: &keptn.ApprovalStrategyStruct{
				Pass:    keptn.Automatic,
				Warning: keptn.Manual,
			},
		},
		{
			Name:               "staging",
			DeploymentStrategy: "blue_green_service",
			TestStrategy:       "performance",
			ApprovalStrategy: &keptn.ApprovalStrategyStruct{
				Pass:    keptn.Automatic,
				Warning: keptn.Manual,
			},
		},
	}
}
