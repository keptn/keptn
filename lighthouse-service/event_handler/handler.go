package event_handler

import (
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

func NewEventHandler(event cloudevents.Event, logger *keptncommon.Logger) (EvaluationEventHandler, error) {
	logger.Debug("Received event: " + event.Type())
	serviceName := "lighthouse-service"

	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{ServiceName: &serviceName},
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
		logger.Info("received unhandled event type")
		return nil, nil
	}
}
