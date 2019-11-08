package event_handler

import (
	"errors"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"net/http"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

func NewEventHandler(event cloudevents.Event, logger *keptnutils.Logger) (EvaluationEventHandler, error) {
	switch event.Type() {
	case keptnevents.TestsFinishedEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event}, nil
	case keptnevents.StartEvaluationEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event}, nil // new event type in Keptn versions >= 0.6
	case keptnevents.InternalGetSLIDoneEventType:
		return &EvaluateSLIHandler{Logger: logger, Event: event, HTTPClient: &http.Client{}}, nil
	default:
		return nil, errors.New("received unknown event type")
	}
}
