package handler

import (
	logger "github.com/sirupsen/logrus"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type Handler interface {
	IsTypeHandled(event cloudevents.Event) bool
	Handle(event cloudevents.Event, keptnHandler *keptnv2.Keptn)
}

func sendEvents(keptnHandler *keptnv2.Keptn, events []cloudevents.Event) {
	for _, outgoingEvent := range events {
		err := keptnHandler.SendCloudEvent(outgoingEvent)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func getCloudEvent(data interface{}, ceType string, shkeptncontext string, triggeredID string) *cloudevents.Event {
	source, _ := url.Parse("approval-service")

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
	err := event.SetData(cloudevents.ApplicationJSON, data)
	if err != nil {
		logger.Errorf("Could not set event data: %v", err)
	}

	return &event
}
