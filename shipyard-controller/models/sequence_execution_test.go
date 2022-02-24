package models

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSequenceExecution_GetNextTriggeredEventData(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        keptnv2.Sequence
		Status          SequenceExecutionStatus
		Scope           EventScope
		InputProperties map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "get initial triggered event - no input data",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "mytask",
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: nil,
					CurrentTask:   TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: nil,
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
			},
		},
		{
			name: "get initial triggered event - with input data",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "mytask",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
						{
							Name: "my-second-task",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: nil,
					CurrentTask:   TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"configurationChange": map[string]interface{}{
						"image": "1.0",
					},
				},
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
				"configurationChange": map[string]interface{}{
					"image": "1.0",
				},
				"mytask": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		{
			name: "get next triggered event - with input data and completed tasks",
			fields: fields{
				Sequence: keptnv2.Sequence{
					Name: "delivery",
					Tasks: []keptnv2.Task{
						{
							Name: "mytask",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
						{
							Name: "my-second-task",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				Status: SequenceExecutionStatus{
					PreviousTasks: []TaskExecutionResult{
						{
							Name:   "mytask",
							Result: keptnv2.ResultPass,
							Status: keptnv2.StatusSucceeded,
							Properties: map[string]interface{}{
								"bar": "foo",
							},
						},
					},
					CurrentTask: TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"configurationChange": map[string]interface{}{
						"image": "1.0",
					},
				},
			},
			want: map[string]interface{}{
				"project": "my-project",
				"stage":   "my-stage",
				"service": "my-service",
				"configurationChange": map[string]interface{}{
					"image": "1.0",
				},
				"mytask": map[string]interface{}{
					"bar": "foo",
				},
				"my-second-task": map[string]interface{}{
					"foo": "bar",
				},
				"result": keptnv2.ResultPass,
				"status": keptnv2.StatusSucceeded,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				ID:              tt.fields.ID,
				Sequence:        tt.fields.Sequence,
				Status:          tt.fields.Status,
				Scope:           tt.fields.Scope,
				InputProperties: tt.fields.InputProperties,
			}
			got := e.GetNextTriggeredEventData()

			require.Equal(t, tt.want, got)
		})
	}
}

func TestSequenceExecution_CompleteCurrentTask(t *testing.T) {
	type fields struct {
		Status SequenceExecutionStatus
	}
	tests := []struct {
		name              string
		fields            fields
		wantResult        keptnv2.ResultType
		wantStatus        keptnv2.StatusType
		wantPreviousTasks []TaskExecutionResult
	}{
		{
			name: "successful task with one executor",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultPass,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultPass,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
					},
				},
			},
		},
		{
			name: "successful task with multiple executors",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-second-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"otherProperty": "otherValue",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultPass,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultPass,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
						"otherProperty": "otherValue",
					},
				},
			},
		},
		{
			name: "multiple executors - one of them failed",
			fields: fields{
				Status: SequenceExecutionStatus{
					CurrentTask: TaskExecutionState{
						Name:        "deployment",
						TriggeredID: "my-triggered-id",
						Events: []TaskEvent{
							{
								EventType: "deployment.started",
								Source:    "my-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-service",
								Result:    keptnv2.ResultPass,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"deploymentURI": "my-deployment-uri",
								},
							},
							{
								EventType: "deployment.started",
								Source:    "my-second-service",
								Result:    "",
								Status:    "",
							},
							{
								EventType: "deployment.finished",
								Source:    "my-second-service",
								Result:    keptnv2.ResultFailed,
								Status:    keptnv2.StatusSucceeded,
								Properties: map[string]interface{}{
									"otherProperty": "otherValue",
								},
							},
						},
					},
				},
			},
			wantResult: keptnv2.ResultFailed,
			wantStatus: keptnv2.StatusSucceeded,
			wantPreviousTasks: []TaskExecutionResult{
				{
					Name:        "deployment",
					TriggeredID: "my-triggered-id",
					Result:      keptnv2.ResultFailed,
					Status:      keptnv2.StatusSucceeded,
					Properties: map[string]interface{}{
						"deploymentURI": "my-deployment-uri",
						"otherProperty": "otherValue",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SequenceExecution{
				Status: tt.fields.Status,
			}
			result, status := e.CompleteCurrentTask()

			require.Equal(t, tt.wantResult, result)
			require.Equal(t, tt.wantStatus, status)

			require.Equal(t, tt.wantPreviousTasks, e.Status.PreviousTasks)
		})
	}
}
