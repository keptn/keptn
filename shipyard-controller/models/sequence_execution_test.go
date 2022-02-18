package models

import (
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSequenceExecution_GetNextTriggeredEventData(t *testing.T) {
	type fields struct {
		ID              string
		Sequence        v0_2_0.Sequence
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
				Sequence: v0_2_0.Sequence{
					Name: "delivery",
					Tasks: []v0_2_0.Task{
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
					EventData: v0_2_0.EventData{
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
				Sequence: v0_2_0.Sequence{
					Name: "delivery",
					Tasks: []v0_2_0.Task{
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
					EventData: v0_2_0.EventData{
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
				Sequence: v0_2_0.Sequence{
					Name: "delivery",
					Tasks: []v0_2_0.Task{
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
							Result: "pass",
							Status: "success",
							Properties: map[string]interface{}{
								"bar": "foo",
							},
						},
					},
					CurrentTask: TaskExecutionState{},
				},
				Scope: EventScope{
					EventData: v0_2_0.EventData{
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
				"result": "pass",
				"status": "success",
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
