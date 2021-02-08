package handler

import (
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnfake "github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_getStartEndTime(t *testing.T) {
	type args struct {
		startDatePoint string
		endDatePoint   string
		timeframe      string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		{
			name: "start and end date provided - return those",
			args: args{
				startDatePoint: time.Now().Round(time.Minute).UTC().Format(dateLayout),
				endDatePoint:   time.Now().Add(5 * time.Minute).UTC().Round(time.Minute).Format(dateLayout),
				timeframe:      "",
			},
			want:    time.Now().Round(time.Minute).UTC(),
			want1:   time.Now().Add(5 * time.Minute).Round(time.Minute).UTC(),
			wantErr: false,
		},
		{
			name: "start and timeframe - return startdate and startdate + timeframe",
			args: args{
				startDatePoint: time.Now().Round(time.Minute).UTC().Format(dateLayout),
				endDatePoint:   "",
				timeframe:      "10m",
			},
			want:    time.Now().Round(time.Minute).UTC(),
			want1:   time.Now().Add(10 * time.Minute).Round(time.Minute).UTC(),
			wantErr: false,
		},
		{
			name: "only timeframe provided - return time.Now - timeframe and time.Now",
			args: args{
				startDatePoint: "",
				endDatePoint:   "",
				timeframe:      "10m",
			},
			want:    time.Now().UTC().Add(-10 * time.Minute).Round(time.Minute).UTC(),
			want1:   time.Now().UTC().Round(time.Minute).UTC(),
			wantErr: false,
		},
		{
			name: "startDate > endDate provided - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(dateLayout),
				endDatePoint:   time.Now().UTC().Format(dateLayout),
				timeframe:      "",
			},
			wantErr: true,
		},
		{
			name: "startDate, endDate and timeframe provided - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(dateLayout),
				endDatePoint:   time.Now().UTC().Format(dateLayout),
				timeframe:      "5m",
			},
			wantErr: true,
		},
		{
			name: "startDate provided, but no endDate or timeframe - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(dateLayout),
				endDatePoint:   "",
				timeframe:      "",
			},
			wantErr: true,
		},
		{
			name: "endDate provided, but no startDate or timeframe - return error",
			args: args{
				startDatePoint: "",
				endDatePoint:   time.Now().Add(1 * time.Minute).UTC().Format(dateLayout),
				timeframe:      "",
			},
			wantErr: true,
		},
		{
			name: "invalid timeframe string - return error",
			args: args{
				startDatePoint: "",
				endDatePoint:   "",
				timeframe:      "xyz",
			},
			wantErr: true,
		},
		{
			name: "invalid timeframe string - return error",
			args: args{
				startDatePoint: "",
				endDatePoint:   "",
				timeframe:      "xym",
			},
			wantErr: true,
		},
		{
			name: "invalid startDate string - return error",
			args: args{
				startDatePoint: "abc",
				endDatePoint:   "",
				timeframe:      "5m",
			},
			wantErr: true,
		},
		{
			name: "invalid endDate string - return error",
			args: args{
				startDatePoint: time.Now().Add(1 * time.Minute).UTC().Format(dateLayout),
				endDatePoint:   "abc",
				timeframe:      "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := getStartEndTime(tt.args.startDatePoint, tt.args.endDatePoint, tt.args.timeframe)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStartEndTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				// if we wanted an error and got it, we don't care about the other return values
				return
			}
			roundedStart := gotStart.Round(time.Minute)
			roundedEnd := gotEnd.Round(time.Minute)
			if roundedStart != tt.want {
				t.Errorf("getStartEndTime() start = %v, want %v", roundedStart, tt.want)
			}
			if roundedEnd != tt.want1 {
				t.Errorf("getStartEndTime() end = %v, want %v", roundedEnd, tt.want1)
			}
		})
	}
}

func TestEvaluationManager_CreateEvaluation(t *testing.T) {
	type fields struct {
		EventSender *keptnfake.EventSender
		ServiceAPI  IServiceAPI
		Logger      keptn.LoggerInterface
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
				ServiceAPI: &fake.IServiceAPIMock{GetServiceFunc: func(project string, stage string, service string) (*apimodels.Service, error) {
					return &apimodels.Service{}, nil
				}},
				Logger: keptn.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				stage:   "test-stage",
				service: "test-service",
				params: &operations.CreateEvaluationParams{
					Labels:    map[string]string{"foo": "bar"},
					Start:     "2020-01-02T15:00:00",
					End:       "2020-01-02T16:00:00",
					Timeframe: "",
				},
			},
			wantResponse: true,
			wantErr:      nil,
			wantEvents: []eventTypeWithPayload{
				{
					eventType: keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
					payload: &keptnv2.EvaluationTriggeredEventData{
						EventData: keptnv2.EventData{
							Project: "test-project",
							Stage:   "test-stage",
							Service: "test-service",
							Labels:  map[string]string{"foo": "bar"},
						},
						Evaluation: keptnv2.Evaluation{
							Start: "2020-01-02 15:00:00 +0000 UTC",
							End:   "2020-01-02 16:00:00 +0000 UTC",
						},
					},
				},
			},
		},
		{
			name: "invalid timeframe",
			fields: fields{
				EventSender: &keptnfake.EventSender{},
				ServiceAPI: &fake.IServiceAPIMock{GetServiceFunc: func(project string, stage string, service string) (*apimodels.Service, error) {
					return &apimodels.Service{}, nil
				}},
				Logger: keptn.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				stage:   "test-stage",
				service: "test-service",
				params: &operations.CreateEvaluationParams{
					Labels:    map[string]string{"foo": "bar"},
					Start:     "2030-01-02T15:00:00",
					End:       "2020-01-02T16:00:00",
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
				ServiceAPI: &fake.IServiceAPIMock{GetServiceFunc: func(project string, stage string, service string) (*apimodels.Service, error) {
					return &apimodels.Service{}, nil
				}},
				Logger: keptn.NewLogger("", "", ""),
			},
			args: args{
				project: "test-project",
				stage:   "test-stage",
				service: "test-service",
				params: &operations.CreateEvaluationParams{
					Labels:    map[string]string{"foo": "bar"},
					Start:     "2020-01-02T15:00:00",
					End:       "2020-01-02T16:00:00",
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
				tt.fields.Logger,
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
				assert.Equal(t, event.Context.GetType(), tt.wantEvents[index].eventType)
				receivedPayload := &keptnv2.EvaluationTriggeredEventData{}
				if err := event.DataAs(receivedPayload); err != nil {
					t.Errorf("could not decode received CloudEvent: %s", err.Error())
				}
				assert.EqualValues(t, receivedPayload, tt.wantEvents[index].payload)
			}

		})
	}
}
