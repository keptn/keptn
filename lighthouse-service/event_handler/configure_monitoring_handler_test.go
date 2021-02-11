package event_handler

import (
	"reflect"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
)

func TestConfigureMonitoringHandler_getSLISourceConfigMap(t *testing.T) {
	type fields struct {
		Logger *keptncommon.Logger
		Event  cloudevents.Event
	}
	type args struct {
		e *keptn.ConfigureMonitoringEventData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1.ConfigMap
	}{
		{
			name: "configure for prometheus monitoring",
			fields: fields{
				Logger: nil,
				Event:  cloudevents.Event{},
			},
			args: args{
				e: &keptn.ConfigureMonitoringEventData{
					Type:    "prometheus",
					Project: "sockshop",
					Service: "",
				},
			},
			want: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "lighthouse-config-sockshop",
					Namespace: "keptn",
				},
				Data: map[string]string{
					"sli-provider": "prometheus",
				},
			},
		},
		{
			name: "configure for dynatrace monitoring",
			fields: fields{
				Logger: nil,
				Event:  cloudevents.Event{},
			},
			args: args{
				e: &keptn.ConfigureMonitoringEventData{
					Type:    "dynatrace",
					Project: "sockshop",
					Service: "",
				},
			},
			want: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "lighthouse-config-sockshop",
					Namespace: "keptn",
				},
				Data: map[string]string{
					"sli-provider": "dynatrace",
				},
			},
		},
	}

	namespace = "keptn"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eh := &ConfigureMonitoringHandler{
				KeptnHandler: nil,
				Event:        tt.fields.Event,
			}
			if got := eh.getSLISourceConfigMap(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSLISourceConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
