package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

type ConfigureMonitoringHandler struct {
	Event        cloudevents.Event
	KeptnHandler *keptnutils.Keptn
}

var namespace = os.Getenv("POD_NAMESPACE")

func (eh *ConfigureMonitoringHandler) HandleEvent() error {

	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnevents.ConfigureMonitoringEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	configMap := eh.getSLISourceConfigMap(e)

	kubeAPI, err := getKubeAPI()
	if err != nil {
		return err
	}

	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not create Kube API")
		return err
	}
	_, err = kubeAPI.CoreV1().ConfigMaps(namespace).Create(configMap)

	if err != nil {
		_, err = kubeAPI.CoreV1().ConfigMaps(namespace).Update(configMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func getKubeAPI() (*kubernetes.Clientset, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	kubeAPI, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeAPI, nil
}

func (eh *ConfigureMonitoringHandler) getSLISourceConfigMap(e *keptnevents.ConfigureMonitoringEventData) *v1.ConfigMap {
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
