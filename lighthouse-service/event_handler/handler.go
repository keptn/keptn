package event_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

// converts a TestsFinishedEventType into a StartEvaluationEventType
func convertTestsFinishedToStartEvaluationEvent(previousEvent cloudevents.Event) cloudevents.Event {
	var shkeptncontext string
	endtime := previousEvent.Time().String()

	var data map[string]string

	previousEvent.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
	previousEvent.Context.ExtensionAs("time", &endtime)
	// eventData := &keptnevents.TestsFinishedEventData{}
	_ = previousEvent.DataAs(&data)

	getSLIEvent := keptnevents.StartEvaluationEventData{
		Project:      data["project"],
		Service:      data["service"],
		Stage:        data["stage"],
		Start:        data["startedat"],
		End:          endtime,
		TestStrategy: data["teststrategy"],
	}

	source, _ := url.Parse("lighthouse-service")
	contentType := "application/json"

	getSLIEventJson, _ := json.Marshal(getSLIEvent)

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.StartEvaluationEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: getSLIEventJson,
	}

	return event
}

func NewEventHandler(event cloudevents.Event, logger *keptnutils.Logger) (EvaluationEventHandler, error) {
	switch event.Type() {
	case keptnevents.TestFinishedEventType_0_5_0_Compatible:
		return &StartEvaluationHandler{Logger: logger, Event: convertTestsFinishedToStartEvaluationEvent(event)}, nil // backwards compatibility to Keptn versions <= 0.5.x
	case keptnevents.TestsFinishedEventType:
		return &StartEvaluationHandler{Logger: logger, Event: convertTestsFinishedToStartEvaluationEvent(event)}, nil
	case keptnevents.StartEvaluationEventType:
		return &StartEvaluationHandler{Logger: logger, Event: event}, nil // new event type in Keptn versions >= 0.6
	case keptnevents.InternalGetSLIDoneEventType:
		return &EvaluateSLIHandler{Logger: logger, Event: event, HTTPClient: &http.Client{}}, nil
	default:
		return nil, errors.New("received unknown event type")
	}
}
