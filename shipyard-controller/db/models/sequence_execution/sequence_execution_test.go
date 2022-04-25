package sequence_execution

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromSequenceExecution(t *testing.T) {
	type args struct {
		se models.SequenceExecution
	}
	tests := []struct {
		name string
		args args
		want SequenceExecution
	}{
		{
			name: "transform sequence execution",
			args: args{
				se: models.SequenceExecution{
					ID: "id",
					Sequence: keptnv2.Sequence{
						Name: "delivery",
						Tasks: []keptnv2.Task{
							{
								Name: "deployment",
								Properties: map[string]string{
									"deployment.strategy": "direct",
								},
							},
							{
								Name: "evaluation",
							},
							{
								Name: "release",
							},
						},
					},
					Status: models.SequenceExecutionStatus{
						State:            "started",
						StateBeforePause: "",
						PreviousTasks: []models.TaskExecutionResult{
							{
								Name:        "deployment",
								TriggeredID: "tr1",
								Result:      "pass",
								Status:      "succeeded",
								Properties: map[string]interface{}{
									"foo.bar": "xyz",
								},
							},
							{
								Name:        "evaluation",
								TriggeredID: "tr2",
								Result:      "pass",
								Status:      "succeeded",
								Properties: map[string]interface{}{
									"foo": map[string]interface{}{
										"bar": "xyz",
									},
								},
							},
						},
						CurrentTask: models.TaskExecutionState{
							Name:        "release",
							TriggeredID: "tr3",
							Events: []models.TaskEvent{
								{
									EventType: keptnv2.GetStartedEventType("release"),
									Source:    "helm",
								},
								{
									EventType: keptnv2.GetFinishedEventType("release"),
									Source:    "helm",
									Properties: map[string]interface{}{
										"release.xyz": "foo",
									},
								},
							},
						},
					},
					Scope: models.EventScope{
						EventData: keptnv2.EventData{
							Project: "my-project",
							Stage:   "my-stage",
							Service: "my-service",
						},
						KeptnContext: "ctx1",
					},
					InputProperties: map[string]interface{}{
						"foo.bar": "xyz",
					},
				},
			},
			want: SequenceExecution{
				ID: "id",
				Sequence: Sequence{
					Name: "delivery",
					Tasks: []Task{
						{
							Name:       "deployment",
							Properties: `{"deployment.strategy":"direct"}`,
						},
						{
							Name: "evaluation",
						},
						{
							Name: "release",
						},
					},
				},
				Status: SequenceExecutionStatus{
					State:            "started",
					StateBeforePause: "",
					PreviousTasks: []TaskExecutionResult{
						{
							Name:        "deployment",
							TriggeredID: "tr1",
							Result:      "pass",
							Status:      "succeeded",
							Properties:  `{"foo.bar":"xyz"}`,
						},
						{
							Name:        "evaluation",
							TriggeredID: "tr2",
							Result:      "pass",
							Status:      "succeeded",
							Properties:  `{"foo":{"bar":"xyz"}}`,
						},
					},
					CurrentTask: TaskExecutionState{
						Name:        "release",
						TriggeredID: "tr3",
						Events: []TaskEvent{
							{
								EventType: keptnv2.GetStartedEventType("release"),
								Source:    "helm",
							},
							{
								EventType:  keptnv2.GetFinishedEventType("release"),
								Source:     "helm",
								Properties: `{"release.xyz":"foo"}`,
							},
						},
					},
				},
				Scope: models.EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
					KeptnContext: "ctx1",
				},
				InputProperties: `{"foo.bar":"xyz"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSequenceExecution(tt.args.se)
			require.Equal(t, tt.want, got)
		})
	}
}
