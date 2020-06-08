package handler

import (
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

const PassResult = "pass"
const WarningResult = "warning"
const FailResult = "fail"
const TestStrategyRealUser = "real-user"
const DeploymentStrategyBlueGreen = "blue_green_service"

const SucceededResult = "succeeded"

type Handler interface {
	IsTypeHandled(event cloudevents.Event) bool
	Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn,
		shipyard *keptnevents.Shipyard)
}

func sendEvents(keptnHandler *keptnevents.Keptn, events []cloudevents.Event, l *keptnevents.Logger) {
	for _, outgoingEvent := range events {
		err := keptnHandler.SendCloudEvent(outgoingEvent)
		if err != nil {
			l.Error(err.Error())
		}
	}
}

func getCloudEvent(data interface{}, ceType string, shkeptncontext string) *cloudevents.Event {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	return &cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        ceType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: data,
	}
}

func getConfigurationChangeEventForCanary(project, service, nextStage, image, shkeptncontext string, labels map[string]string) *cloudevents.Event {

	valuesCanary := make(map[string]interface{})
	valuesCanary["image"] = image
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:      project,
		Service:      service,
		Stage:        nextStage,
		ValuesCanary: valuesCanary,
		Canary:       &keptnevents.Canary{Action: keptnevents.Set, Value: 100},
		Labels:       labels,
	}

	return getCloudEvent(configChangedEvent, keptnevents.ConfigurationChangeEventType, shkeptncontext)
}
