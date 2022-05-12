package v1

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
		want JsonStringEncodedSequenceExecution
	}{
		{
			name: "transform sequence execution",
			args: args{
				se: testSequenceExecution,
			},
			want: testJsonStringEncodedSequenceExecution,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := ModelTransformer{}
			got := mt.TransformToDBModel(tt.args.se)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestModelTransformer_TransformEventToDBModel(t *testing.T) {
	type args struct {
		event models.TaskEvent
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "transform event",
			args: args{
				event: models.TaskEvent{
					EventType: "mytask.triggered",
					Source:    "my-service",
					Result:    "pass",
					Status:    "succeeded",
					Properties: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			want: TaskEvent{
				EventType:         "mytask.triggered",
				Source:            "my-service",
				Result:            "pass",
				Status:            "succeeded",
				EncodedProperties: `{"foo":"bar"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mo := ModelTransformer{}
			got := mo.TransformEventToDBModel(tt.args.event)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestModelTransformer_TransformToSequenceExecution(t *testing.T) {
	type args struct {
		dbItem interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *models.SequenceExecution
		wantErr bool
	}{
		{
			name: "transform v1 schema",
			args: args{
				dbItem: JsonStringEncodedSequenceExecution{
					ID: "1",
					SchemaVersion: SchemaVersion{
						SchemaVersion: SchemaVersionV1,
					},
					Sequence: Sequence{
						Name: "my-sequence",
						Tasks: []Task{
							{
								Name:              "delivery",
								TriggeredAfter:    "1m",
								EncodedProperties: `{"foo":"bar"}`,
							},
						},
					},
					Status: SequenceExecutionStatus{
						State: "started",
					},
					Scope: models.EventScope{
						EventData: keptnv2.EventData{
							Project: "my-project",
							Stage:   "my-stage",
							Service: "my-service",
						},
					},
					EncodedInputProperties: `{"foo":"bar"}`,
				},
			},
			want: &models.SequenceExecution{
				ID:            "1",
				SchemaVersion: SchemaVersionV1,
				Sequence: keptnv2.Sequence{
					Name: "my-sequence",
					Tasks: []keptnv2.Task{
						{
							Name:           "delivery",
							TriggeredAfter: "1m",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				Status: models.SequenceExecutionStatus{
					State:         "started",
					PreviousTasks: []models.TaskExecutionResult{},
					CurrentTask: models.TaskExecutionState{
						Events: []models.TaskEvent{},
					},
				},
				Scope: models.EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"foo": "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "transform previous schema",
			args: args{
				dbItem: &models.SequenceExecution{
					ID: "1",
					Sequence: keptnv2.Sequence{
						Name: "my-sequence",
						Tasks: []keptnv2.Task{
							{
								Name:           "delivery",
								TriggeredAfter: "1m",
								Properties: map[string]interface{}{
									"foo": "bar",
								},
							},
						},
					},
					Status: models.SequenceExecutionStatus{
						State:         "started",
						PreviousTasks: []models.TaskExecutionResult{},
						CurrentTask: models.TaskExecutionState{
							Events: []models.TaskEvent{},
						},
					},
					Scope: models.EventScope{
						EventData: keptnv2.EventData{
							Project: "my-project",
							Stage:   "my-stage",
							Service: "my-service",
						},
					},
					InputProperties: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			want: &models.SequenceExecution{
				ID: "1",
				Sequence: keptnv2.Sequence{
					Name: "my-sequence",
					Tasks: []keptnv2.Task{
						{
							Name:           "delivery",
							TriggeredAfter: "1m",
							Properties: map[string]interface{}{
								"foo": "bar",
							},
						},
					},
				},
				Status: models.SequenceExecutionStatus{
					State:         "started",
					PreviousTasks: []models.TaskExecutionResult{},
					CurrentTask: models.TaskExecutionState{
						Events: []models.TaskEvent{},
					},
				},
				Scope: models.EventScope{
					EventData: keptnv2.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				InputProperties: map[string]interface{}{
					"foo": "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid object",
			args: args{
				dbItem: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid object with schema version 1",
			args: args{
				dbItem: map[string]interface{}{
					"schemaVersion": "1",
					"sequence":      "invalid",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid object with unspecified schema version",
			args: args{
				dbItem: map[string]interface{}{
					"schemaVersion": "0",
					"sequence":      "invalid",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mo := ModelTransformer{}
			got, err := mo.TransformToSequenceExecution(tt.args.dbItem)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
