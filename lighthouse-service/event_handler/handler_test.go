package event_handler

import (
	"context"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-test/deep"

	"github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
)

func TestNewEventHandler(t *testing.T) {
	incomingEvent := cloudevents.NewEvent()
	incomingEvent.SetID("my-id")
	incomingEvent.SetSource("my-source")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fakeSender := func(ce models.KeptnContextExtendedCE) error { return nil }
	ctx = context.WithValue(ctx, types.EventSenderKey, controlplane.EventSender(fakeSender))
	defer cancel()

	keptnHandler, _ := keptnv2.NewKeptn(&incomingEvent, keptncommon.KeptnOpts{})

	type args struct {
		event cloudevents.Event
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
				event: incomingEvent,
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
				event: incomingEvent,
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
				event: incomingEvent,
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
				event: incomingEvent,
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

			got, err := NewEventHandler(ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEventHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			deep.MaxDepth = 5
			if len(deep.Equal(got, tt.want)) > 0 {
				t.Errorf("NewEventHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventSenderWithoutContext(t *testing.T) {
	incomingEvent := cloudevents.NewEvent()
	incomingEvent.SetID("my-id")
	incomingEvent.SetSource("my-source")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx = context.WithValue(ctx, types.EventSenderKey, nil)
	defer cancel()

	_, err := NewEventHandler(ctx, incomingEvent)
	require.Error(t, err)
	require.Equal(t, "could not get eventSender from context", err.Error())

}
