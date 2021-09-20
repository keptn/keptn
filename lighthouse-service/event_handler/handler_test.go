package event_handler

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"net/http"
	"os"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-test/deep"

	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

func TestNewEventHandler(t *testing.T) {
	incomingEvent := cloudevents.NewEvent()
	incomingEvent.SetID("my-id")
	incomingEvent.SetSource("my-source")

	serviceName := "lighthouse-service"
	keptnHandler, _ := keptnv2.NewKeptn(&incomingEvent, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{ServiceName: &serviceName},
	})

	type args struct {
		event  cloudevents.Event
		logger *keptncommon.Logger
	}
	tests := []struct {
		name      string
		args      args
		eventType string
		want      EvaluationEventHandler
		wantErr   bool
	}{
		{
			name: "evaluation.triggered -> start-evaluation handler",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
			want: &StartEvaluationHandler{
				Event:             incomingEvent,
				KeptnHandler:      keptnHandler,
				SLIProviderConfig: K8sSLIProviderConfig{},
			},
			wantErr: false,
		},
		{
			name: "get-sli.done -> evaluate-sli handler",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName),
			want: &EvaluateSLIHandler{
				Event:        incomingEvent,
				KeptnHandler: keptnHandler,
				HTTPClient:   &http.Client{},
				EventStore:   keptnHandler.EventHandler,
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
				Event:     incomingEvent,
				Logger:    logrus.New(),
				K8sClient: fake.NewSimpleClientset(),
			},
			wantErr: false,
		},
		{
			name: "invalid event type -> error",
			args: args{
				event:  incomingEvent,
				logger: nil,
			},
			eventType: "nonsense-event",
			want:      nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetConfig().GetKubeAPI = func() (kubernetes.Interface, error) {
				return fake.NewSimpleClientset(), nil
			}
			tt.args.event.SetType(tt.eventType)
			os.Setenv("CONFIGURATION_SERVICE", configurationServiceURL)
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
