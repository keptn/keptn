package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/keptn/go-utils/pkg/lib"
	"net/http"
	"reflect"
	"testing"
)

func TestNewEventHandler(t *testing.T) {
	type args struct {
		event  cloudevents.Event
		logger *keptn.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    EvaluationEventHandler
		wantErr bool
	}{
		{
			name: "tests-finished -> start-evaluation handler",
			args: args{
				event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.TestsFinishedEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
				logger: nil,
			},
			want: &StartEvaluationHandler{
				Logger: nil,
				Event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.TestsFinishedEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
			},
			wantErr: false,
		},
		{
			name: "start-evaluation -> start-evaluation handler",
			args: args{
				event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.StartEvaluationEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
				logger: nil,
			},
			want: &StartEvaluationHandler{
				Logger: nil,
				Event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.StartEvaluationEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
			},
			wantErr: false,
		},
		{
			name: "get-sli.done -> evaluate-sli handler",
			args: args{
				event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.InternalGetSLIDoneEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
				logger: nil,
			},
			want: &EvaluateSLIHandler{
				Logger: nil,
				Event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.InternalGetSLIDoneEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
				HTTPClient: &http.Client{},
			},
			wantErr: false,
		},
		{
			name: "configure-monitoring -> configure monitoring handler",
			args: args{
				event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.ConfigureMonitoringEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
				logger: nil,
			},
			want: &ConfigureMonitoringHandler{
				Logger: nil,
				Event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: keptn.ConfigureMonitoringEventType,
					},
					Data:        nil,
					DataEncoded: false,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid event type -> error",
			args: args{
				event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						Type: "hulumulu",
					},
					Data:        nil,
					DataEncoded: false,
				},
				logger: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEventHandler(tt.args.event, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEventHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
