package event_handler

import (
	"context"
	"github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"os"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigureMonitoringHandler struct {
	Event     cloudevents.Event
	Logger    *logrus.Logger
	K8sClient kubernetes.Interface
}

var namespace = os.Getenv("POD_NAMESPACE")

type ConfigureMonitoringHandlerOption func(cmh *ConfigureMonitoringHandler)

func WithK8sClient(k8sClient kubernetes.Interface) ConfigureMonitoringHandlerOption {
	return func(cmh *ConfigureMonitoringHandler) {
		cmh.K8sClient = k8sClient
	}
}

func NewConfigureMonitoringHandler(event cloudevents.Event, logger *logrus.Logger, opts ...ConfigureMonitoringHandlerOption) (*ConfigureMonitoringHandler, error) {
	cmh := &ConfigureMonitoringHandler{
		Event:  event,
		Logger: logger,
	}

	for _, opt := range opts {
		opt(cmh)
	}

	if cmh.K8sClient == nil {
		defaultK8sClient, err := GetConfig().GetKubeAPI()
		if err != nil {
			return nil, err
		}
		cmh.K8sClient = defaultK8sClient
	}

	return cmh, nil
}

func (eh *ConfigureMonitoringHandler) HandleEvent(ctx context.Context) error {
	ctx.Value(GracefulShutdownKey).(*sync.WaitGroup).Add(1)
	defer func() {
		ctx.Value(GracefulShutdownKey).(*sync.WaitGroup).Done()
	}()

	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnevents.ConfigureMonitoringEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		eh.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	configMap := getSLISourceConfigMap(e)

	_, err = eh.K8sClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})

	if err != nil && k8serrors.IsAlreadyExists(err) {
		_, err = eh.K8sClient.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
		if err != nil {
			eh.Logger.WithError(err).Error("could not update sli-provider ConfigMap")
			return err
		}
	} else if err != nil {
		eh.Logger.WithError(err).Error("could not create sli-provider ConfigMap")
		return err
	}
	return nil
}

func getSLISourceConfigMap(e *keptnevents.ConfigureMonitoringEventData) *v1.ConfigMap {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "lighthouse-config-" + e.Project,
			Namespace: namespace,
		},
		Data: map[string]string{
			"sli-provider": e.Type,
		},
	}
	return configMap
}
