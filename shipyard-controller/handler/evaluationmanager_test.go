package handler

import (
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnfake "github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluationManager_CreateEvaluation(t *testing.T) {
	type fields struct {
		EventSender *keptnfake.EventSender
		ServiceAPI  db.ServicesDbOperations
	}
	type args struct {
		project string
		stage   string
		service string
		params  *operations.CreateEvaluationParams
	}
	type eventTypeWithPayload struct {
		eventType string
		payload   *keptnv2.EvaluationTriggeredEventData
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse bool
		wantErr      *models.Error
		wantEvents   []eventTypeWithPayload
	}{
		{
			name: "everything ok - send evaluation.triggered event",
			fields: fields{
				EventSender: &keptnfake.EventSender{},
				ServiceAPI: &db_mock.ServicesDbOperationsMock{GetServiceFunc: func(project string, stage string, service string) (*models.ExpandedService, error) {
					return &models.ExpandedService{}, nil
				}},
			},
			args: args{
				project: "test-project",
				stage:   "test-stage",
				service: "test-service",
				params: &operations.CreateEvaluationParams{
					Labels:    map[string]string{"foo": "bar"},
					Start:     "2020-01-02T15:00:00.000Z",
					End:       "2020-01-02T16:00:00.000Z",
					Timeframe: "",
				},
			},
			wantResponse: true,
			wantErr:      nil,
			wantEvents: []eventTypeWithPayload{
				{
					eventType: keptnv2.GetTriggeredEventType("test-stage." + keptnv2.EvaluationTaskName),
					payload: &keptnv2.EvaluationTriggeredEventData{
						EventData: keptnv2.EventData{
							Project: "test-project",
							Stage:   "test-stage",
							Service: "test-service",
							Labels:  map[string]string{"foo": "bar"},
						},
						Evaluation: keptnv2.Evaluation{
							Start: "2020-01-02T15:00:00.000Z",
							End:   "2020-01-02T16:00:00.000Z",
						},
					},
				},
			},
		},
		{
			name: "invalid timeframe",
			fields: fields{
				EventSender: &keptnfake.EventSender{},
				ServiceAPI: &db_mock.ServicesDbOperationsMock{GetServiceFunc: func(project string, stage string, service string) (*models.ExpandedService, error) {
					return &models.ExpandedService{}, nil
				}},
			},
			args: args{
				project: "test-project",
				stage:   "test-stage",
				service: "test-service",
				params: &operations.CreateEvaluationParams{
					Labels:    map[string]string{"foo": "bar"},
					Start:     "2030-01-02T15:00:00.000Z",
					End:       "2020-01-02T16:00:00.000Z",
					Timeframe: "",
				},
			},
			wantResponse: false,
			wantErr:      &models.Error{Code: evaluationErrInvalidTimeframe},
			wantEvents:   []eventTypeWithPayload{},
		},
		{
			name: "sending event failed",
			fields: fields{
				EventSender: &keptnfake.EventSender{
					Reactors: map[string]func(event cloudevents.Event) error{
						"*": func(event cloudevents.Event) error {
							return errors.New("")
						},
					},
				},
				ServiceAPI: &db_mock.ServicesDbOperationsMock{GetServiceFunc: func(project string, stage string, service string) (*models.ExpandedService, error) {
					return &models.ExpandedService{}, nil
				}},
			},
			args: args{
				project: "test-project",
				stage:   "test-stage",
				service: "test-service",
				params: &operations.CreateEvaluationParams{
					Labels:    map[string]string{"foo": "bar"},
					Start:     "2020-01-02T15:00:00.000Z",
					End:       "2020-01-02T16:00:00.000Z",
					Timeframe: "",
				},
			},
			wantResponse: false,
			wantErr:      &models.Error{Code: evaluationErrSendEventFailed},
			wantEvents:   []eventTypeWithPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em, err := NewEvaluationManager(
				tt.fields.EventSender,
				tt.fields.ServiceAPI,
			)
			if err != nil {
				t.Error(err.Error())
			}

			gotContext, gotErr := em.CreateEvaluation(tt.args.project, tt.args.stage, tt.args.service, tt.args.params)

			if tt.wantErr != nil {
				if gotErr == nil {
					t.Error("CreateEvaluation() - expected error but did not get any")
				}
				assert.Equal(t, gotErr.Code, tt.wantErr.Code)
			} else if gotErr != nil {
				t.Errorf("CreateEvaluation() - expected no error but got %v", gotErr)
			}

			if tt.wantResponse {
				assert.NotEmpty(t, gotContext.KeptnContext)
			}

			if len(tt.fields.EventSender.SentEvents) != len(tt.wantEvents) {
				t.Error("Did not receive expected number of events")
			}

			for index, event := range tt.fields.EventSender.SentEvents {
				assert.Equal(t, tt.wantEvents[index].eventType, event.Context.GetType())
				receivedPayload := &keptnv2.EvaluationTriggeredEventData{}
				if err := event.DataAs(receivedPayload); err != nil {
					t.Errorf("could not decode received CloudEvent: %s", err.Error())
				}
				assert.EqualValues(t, receivedPayload, tt.wantEvents[index].payload)
			}

		})
	}
}
