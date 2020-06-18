package event_handler

import (
	"errors"
	"net/http"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/keptn/go-utils/pkg/lib"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

func NewEventHandler(event cloudevents.Event, logger *keptn.Logger) (EvaluationEventHandler, error) {
	logger.Debug("Received event: " + event.Type())
	serviceName := "lighthouse-service"
	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{
		LoggingOptions: &keptn.LoggingOpts{ServiceName: &serviceName},
	})
	if err != nil {
		return nil, err
	}
	switch event.Type() {
	case keptn.TestsFinishedEventType:
		return &StartEvaluationHandler{Event: event, KeptnHandler: keptnHandler}, nil
	case keptn.StartEvaluationEventType:
		return &StartEvaluationHandler{Event: event, KeptnHandler: keptnHandler}, nil // new event type in Keptn versions >= 0.6
	case keptn.InternalGetSLIDoneEventType:
		return &EvaluateSLIHandler{Event: event, HTTPClient: &http.Client{}, KeptnHandler: keptnHandler}, nil
	case keptn.ConfigureMonitoringEventType:
		return &ConfigureMonitoringHandler{Event: event, KeptnHandler: keptnHandler}, nil
	default:
		return nil, errors.New("received unknown event type")
	}
}
