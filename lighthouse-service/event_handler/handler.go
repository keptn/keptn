package event_handler

import (
	"errors"
	"net/http"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

func NewEventHandler(event cloudevents.Event, logger *keptnutils.Logger) (EvaluationEventHandler, error) {
	logger.Debug("Received event: " + event.Type())
	switch event.Type() {
	case keptnevents.TestsFinishedEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event}, nil
	case keptnevents.StartEvaluationEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event}, nil // new event type in Keptn versions >= 0.6
	case keptnevents.InternalGetSLIDoneEventType:
		return &EvaluateSLIHandler{Logger: logger, Event: event, HTTPClient: &http.Client{}}, nil
	case keptnevents.ConfigureMonitoringEventType:
		return &ConfigureMonitoringHandler{Logger: logger, Event: event}, nil
	default:
		return nil, errors.New("received unknown event type")
	}
}
