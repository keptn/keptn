package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-test/deep"
	"github.com/keptn/go-utils/pkg/lib"
	"net/http"
	"net/url"
	"testing"
)

func TestNewEventHandler(t *testing.T) {
	incomingEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV02{
			SpecVersion: "0.2",
			Type:        keptn.TestsFinishedEventType,
			Source:      types.URLRef{URL: url.URL{Host: "test"}},
			ID:          "1",
			Time:        nil,
			SchemaURL:   nil,
			ContentType: stringp("application/json"),
			Extensions:  nil,
		},
		Data: []byte(`{
    "project": "sockshop",
    "stage": "staging",
    "service": "carts"
  }`),
		DataEncoded: false,
	}

	serviceName := "lighthouse-service"
	keptnHandler, _ := keptn.NewKeptn(&incomingEvent, keptn.KeptnOpts{
		LoggingOptions: &keptn.LoggingOpts{ServiceName: &serviceName},
	})

	type args struct {
		event  cloudevents.Event
		logger *keptn.Logger
	}
	tests := []struct {
		name      string
		args      args
		eventType string
		want      EvaluationEventHandler
		wantErr   bool
	}{
		{
			name: "tests-finished -> start-evaluation handler",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: keptn.TestsFinishedEventType,
			want: &StartEvaluationHandler{
				Event:        incomingEvent,
				KeptnHandler: keptnHandler,
			},
			wantErr: false,
		},
		{
			name: "start-evaluation -> start-evaluation handler",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: keptn.StartEvaluationEventType,
			want: &StartEvaluationHandler{
				Event:        incomingEvent,
				KeptnHandler: keptnHandler,
			},
			wantErr: false,
		},
		{
			name: "get-sli.done -> evaluate-sli handler",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: keptn.InternalGetSLIDoneEventType,
			want: &EvaluateSLIHandler{
				Event:        incomingEvent,
				KeptnHandler: keptnHandler,
				HTTPClient:   &http.Client{},
			},
			wantErr: false,
		},
		{
			name: "configure-monitoring -> configure monitoring handler",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: keptn.ConfigureMonitoringEventType,
			want: &ConfigureMonitoringHandler{
				Event:        incomingEvent,
				KeptnHandler: keptnHandler,
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
			tt.args.event.SetType(tt.eventType)
			got, err := NewEventHandler(tt.args.event, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEventHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(deep.Equal(got, tt.want)) > 0 {
				t.Errorf("NewEventHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
