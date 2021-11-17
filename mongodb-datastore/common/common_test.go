package common

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_transformEvaluationDoneEvent(t *testing.T) {
	type args struct {
		keptnEvent models.KeptnContextExtendedCE
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantEvent models.KeptnContextExtendedCE
	}{
		{
			name: "transform evaluation-done event",
			args: args{
				keptnEvent: models.KeptnContextExtendedCE{
					Event: models.Event{
						Contenttype: "",
						Data: map[string]interface{}{
							"result":  "pass",
							"project": "my-project",
							"stage":   "my-stage",
							"service": "my-service",
							"labels": map[string]interface{}{
								"foo": "bar",
							},
							"evaluationdetails": keptnv2.EvaluationDetails{
								Result: string(keptnv2.ResultPass),
								Score:  10,
							},
						},
						Extensions:  nil,
						ID:          "",
						Source:      "lighthouse-service",
						Specversion: "0.2",
						Time:        models.Time{},
						Type:        Keptn07EvaluationDoneEventType,
					},
					Shkeptncontext: "my-context",
					Triggeredid:    "my-triggeredid",
				},
			},
			wantEvent: models.KeptnContextExtendedCE{
				Event: models.Event{
					Contenttype: "",
					Data: &keptnv2.EvaluationFinishedEventData{
						EventData: keptnv2.EventData{
							Project: "my-project",
							Stage:   "my-stage",
							Service: "my-service",
							Labels: map[string]string{
								"foo": "bar",
							},
							Result: keptnv2.ResultPass,
						},
						Evaluation: keptnv2.EvaluationDetails{
							Result: string(keptnv2.ResultPass),
							Score:  10,
						},
					},
					Extensions:  nil,
					ID:          "",
					Source:      "lighthouse-service",
					Specversion: "1.0",
					Time:        models.Time{},
					Type:        models.Type(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)),
				},
				Shkeptncontext: "my-context",
				Triggeredid:    "my-triggeredid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			convertedEvent, err := TransformEvaluationDoneEvent(tt.args.keptnEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformEvaluationDoneEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.Equal(t, &tt.wantEvent, convertedEvent)
		})
	}
}
