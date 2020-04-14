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
	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{})
	if err != nil {
		return nil, err
	}
	switch event.Type() {
	case keptn.TestsFinishedEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event, KeptnHandler: keptnHandler}, nil
	case keptn.StartEvaluationEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event, KeptnHandler: keptnHandler}, nil // new event type in Keptn versions >= 0.6
	case keptn.InternalGetSLIDoneEventType:
		return &EvaluateSLIHandler{Logger: logger, Event: event, HTTPClient: &http.Client{}, KeptnHandler: keptnHandler}, nil
	case keptn.ConfigureMonitoringEventType:
		return &ConfigureMonitoringHandler{Logger: logger, Event: event, KeptnHandler: keptnHandler}, nil
	default:
		return nil, errors.New("received unknown event type")
	}
}
