package v1

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

var testSequenceExecution = models.SequenceExecution{
	ID:            "id",
	SchemaVersion: SchemaVersionV1,
	Sequence: keptnv2.Sequence{
		Name: "delivery",
		Tasks: []keptnv2.Task{
			{
				Name: "deployment",
				Properties: map[string]interface{}{
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
}

var testJsonStringEncodedSequenceExecution = JsonStringEncodedSequenceExecution{
	ID: "id",
	SchemaVersion: SchemaVersion{
		SchemaVersion: SchemaVersionV1,
	},
	Sequence: Sequence{
		Name: "delivery",
		Tasks: []Task{
			{
				Name:              "deployment",
				EncodedProperties: `{"deployment.strategy":"direct"}`,
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
				Name:              "deployment",
				TriggeredID:       "tr1",
				Result:            "pass",
				Status:            "succeeded",
				EncodedProperties: `{"foo.bar":"xyz"}`,
			},
			{
				Name:              "evaluation",
				TriggeredID:       "tr2",
				Result:            "pass",
				Status:            "succeeded",
				EncodedProperties: `{"foo":{"bar":"xyz"}}`,
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
					EventType:         keptnv2.GetFinishedEventType("release"),
					Source:            "helm",
					EncodedProperties: `{"release.xyz":"foo"}`,
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
	EncodedInputProperties: `{"foo.bar":"xyz"}`,
}

func TestJsonStringEncodedSequenceExecution_ToSequenceExecution(t *testing.T) {
	tests := []struct {
		name    string
		obj     JsonStringEncodedSequenceExecution
		want    models.SequenceExecution
		wantErr bool
	}{
		{
			name:    "transform back to sequence execution",
			obj:     testJsonStringEncodedSequenceExecution,
			want:    testSequenceExecution,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.ToSequenceExecution()
			require.Equal(t, tt.want, got)
		})
	}
}
