package handler

import (
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const (
	PassResult                  = "pass"
	WarningResult               = "warning"
	FailResult                  = "fail"
	TestStrategyRealUser        = "real-user"
	DeploymentStrategyBlueGreen = "blue_green_service"
	SucceededResult             = "succeeded"
)

type Handler interface {
	IsTypeHandled(event cloudevents.Event) bool
	Handle(event cloudevents.Event, keptnHandler *keptnv2.Keptn)
}

func sendEvents(keptnHandler *keptnv2.Keptn, events []cloudevents.Event, l keptncommon.LoggerInterface) {
	for _, outgoingEvent := range events {
		err := keptnHandler.SendCloudEvent(outgoingEvent)
		if err != nil {
			l.Error(err.Error())
		}
	}
}

func getCloudEvent(data interface{}, ceType string, shkeptncontext string, triggeredID string) *cloudevents.Event {
	source, _ := url.Parse("gatekeeper-service")

	extensions := map[string]interface{}{"shkeptncontext": shkeptncontext}
	if triggeredID != "" {
		extensions["triggeredid"] = triggeredID
	}

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetType(ceType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetExtension("triggeredid", triggeredID)
	event.SetData(cloudevents.ApplicationJSON, data)

	return &event
}
