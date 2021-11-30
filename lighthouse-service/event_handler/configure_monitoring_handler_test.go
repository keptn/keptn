package event_handler

import (
	"context"
	"errors"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"reflect"
	"sync"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

func TestConfigureMonitoringHandler_getSLISourceConfigMap(t *testing.T) {
	type args struct {
		e *keptn.ConfigureMonitoringEventData
	}
	tests := []struct {
		name string
		args args
		want *v1.ConfigMap
	}{
		{
			name: "configure for prometheus monitoring",
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
			if got := getSLISourceConfigMap(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSLISourceConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigureMonitoringHandler_HandleEvent_ConfigMapDoesntExistYet(t *testing.T) {
	ce := cloudevents.NewEvent()
	wg := &sync.WaitGroup{}
	ctx := cloudevents.WithEncodingStructured(context.WithValue(context.Background(), GracefulShutdownKey, wg))
	configureMonitoringData := &keptnevents.ConfigureMonitoringEventData{
		Project: "my-project",
		Service: "my-service",
		Type:    "my-sli-provider",
	}
	ce.SetType(keptnv2.GetTriggeredEventType(keptnv2.ConfigureMonitoringTaskName))
	ce.SetData(cloudevents.ApplicationJSON, configureMonitoringData)

	logger, _ := test.NewNullLogger()

	fakeK8sClient := fake.NewSimpleClientset()

	handler, err := NewConfigureMonitoringHandler(ce, logger, WithK8sClient(fakeK8sClient))
	require.Nil(t, err)

	err = handler.HandleEvent(ctx)
	require.Nil(t, err)
	require.Len(t, fakeK8sClient.Actions(), 1)
	require.Equal(t, "create", fakeK8sClient.Actions()[0].GetVerb())
}

func TestConfigureMonitoringHandler_HandleEvent_ConfigMapDoesntExistYetAndCreateFails(t *testing.T) {
	ce := cloudevents.NewEvent()
	wg := &sync.WaitGroup{}
	ctx := cloudevents.WithEncodingStructured(context.WithValue(context.Background(), GracefulShutdownKey, wg))
	configureMonitoringData := &keptnevents.ConfigureMonitoringEventData{
		Project: "my-project",
		Service: "my-service",
		Type:    "my-sli-provider",
	}
	ce.SetType(keptnv2.GetTriggeredEventType(keptnv2.ConfigureMonitoringTaskName))
	ce.SetData(cloudevents.ApplicationJSON, configureMonitoringData)

	logger, hook := test.NewNullLogger()

	fakeK8sClient := fake.NewSimpleClientset()

	fakeK8sClient.PrependReactor("create", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("oops")
	})

	handler, err := NewConfigureMonitoringHandler(ce, logger, WithK8sClient(fakeK8sClient))
	require.Nil(t, err)

	err = handler.HandleEvent(ctx)
	require.NotNil(t, err)
	require.Len(t, fakeK8sClient.Actions(), 1)
	require.Equal(t, "create", fakeK8sClient.Actions()[0].GetVerb())

	require.NotNil(t, hook)
	require.NotEmpty(t, hook.Entries)
	require.Contains(t, hook.LastEntry().Message, "could not create")
	require.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
}

func TestConfigureMonitoringHandler_HandleEvent_ConfigMapAlreadyExists(t *testing.T) {
	ce := cloudevents.NewEvent()
	wg := &sync.WaitGroup{}
	ctx := cloudevents.WithEncodingStructured(context.WithValue(context.Background(), GracefulShutdownKey, wg))
	configureMonitoringData := &keptnevents.ConfigureMonitoringEventData{
		Project: "my-project",
		Service: "my-service",
		Type:    "my-sli-provider",
	}
	ce.SetType(keptnv2.GetTriggeredEventType(keptnv2.ConfigureMonitoringTaskName))
	ce.SetData(cloudevents.ApplicationJSON, configureMonitoringData)

	logger, _ := test.NewNullLogger()

	// initialize the fake k8s client with an already existing configmap for the project
	fakeK8sClient := fake.NewSimpleClientset(getSLISourceConfigMap(configureMonitoringData))

	handler, err := NewConfigureMonitoringHandler(ce, logger, WithK8sClient(fakeK8sClient))
	require.Nil(t, err)

	err = handler.HandleEvent(ctx)
	require.Nil(t, err)
	require.Len(t, fakeK8sClient.Actions(), 2)
	require.Equal(t, "create", fakeK8sClient.Actions()[0].GetVerb())
	require.Equal(t, "update", fakeK8sClient.Actions()[1].GetVerb())
}

func TestConfigureMonitoringHandler_HandleEvent_ConfigMapAlreadyExistsUpdateFails(t *testing.T) {
	ce := cloudevents.NewEvent()
	wg := &sync.WaitGroup{}
	ctx := cloudevents.WithEncodingStructured(context.WithValue(context.Background(), GracefulShutdownKey, wg))
	configureMonitoringData := &keptnevents.ConfigureMonitoringEventData{
		Project: "my-project",
		Service: "my-service",
		Type:    "my-sli-provider",
	}
	ce.SetType(keptnv2.GetTriggeredEventType(keptnv2.ConfigureMonitoringTaskName))
	ce.SetData(cloudevents.ApplicationJSON, configureMonitoringData)

	logger, hook := test.NewNullLogger()

	// initialize the fake k8s client with an already existing configmap for the project
	fakeK8sClient := fake.NewSimpleClientset(getSLISourceConfigMap(configureMonitoringData))

	fakeK8sClient.PrependReactor("update", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("oops")
	})

	handler, err := NewConfigureMonitoringHandler(ce, logger, WithK8sClient(fakeK8sClient))
	require.Nil(t, err)

	err = handler.HandleEvent(ctx)
	require.NotNil(t, err)
	require.Len(t, fakeK8sClient.Actions(), 2)
	require.Equal(t, "create", fakeK8sClient.Actions()[0].GetVerb())
	require.Equal(t, "update", fakeK8sClient.Actions()[1].GetVerb())

	require.NotNil(t, hook)
	require.NotEmpty(t, hook.Entries)
	require.Contains(t, hook.LastEntry().Message, "could not update")
	require.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
}
