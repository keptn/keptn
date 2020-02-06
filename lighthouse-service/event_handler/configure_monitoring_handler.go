package event_handler

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigureMonitoringHandler struct {
	Logger *keptnutils.Logger
	Event  cloudevents.Event
}

func (eh *ConfigureMonitoringHandler) HandleEvent() error {

	var keptnContext string
	_ = eh.Event.ExtensionAs("shkeptncontext", &keptnContext)

	e := &keptnevents.ConfigureMonitoringEventData{}
	err := eh.Event.DataAs(e)
	if err != nil {
		eh.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	configMap := eh.getSLISourceConfigMap(e)

	kubeAPI, err := keptnutils.GetKubeAPI(true)

	if err != nil {
		eh.Logger.Error("Could not create Kube API")
		return err
	}
	_, err = kubeAPI.ConfigMaps("keptn").Create(configMap)

	if err != nil {
		_, err = kubeAPI.ConfigMaps("keptn").Update(configMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (eh *ConfigureMonitoringHandler) getSLISourceConfigMap(e *keptnevents.ConfigureMonitoringEventData) *v1.ConfigMap {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "lighthouse-config-" + e.Project,
			Namespace: "keptn",
		},
		Data: map[string]string{
			"sli-provider": e.Type,
		},
	}
	return configMap
}
