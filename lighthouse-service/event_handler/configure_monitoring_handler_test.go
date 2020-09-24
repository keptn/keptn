package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestConfigureMonitoringHandler_getSLISourceConfigMap(t *testing.T) {
	type fields struct {
		Logger *keptn.Logger
		Event  cloudevents.Event
	}
	type args struct {
		e *keptnevents.ConfigureMonitoringEventData
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
				e: &keptnevents.ConfigureMonitoringEventData{
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
				e: &keptnevents.ConfigureMonitoringEventData{
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
