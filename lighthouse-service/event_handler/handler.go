package event_handler

import (
	"net/http"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnv1 "github.com/keptn/go-utils/pkg/lib"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

func NewEventHandler(event cloudevents.Event, logger *keptn.Logger) (EvaluationEventHandler, error) {
	logger.Debug("Received event: " + event.Type())
	serviceName := "lighthouse-service"
	keptnHandler, err := keptnv1.NewKeptn(&event, keptn.KeptnOpts{
		LoggingOptions: &keptn.LoggingOpts{ServiceName: &serviceName},
	})
	if err != nil {
		return nil, err
	}
	switch event.Type() {
	case keptnv1.TestsFinishedEventType:
		return &StartEvaluationHandler{Event: event, KeptnHandler: keptnHandler, SLIProviderConfig: K8sSLIProviderConfig{}}, nil
	case keptnv1.StartEvaluationEventType:
		return &StartEvaluationHandler{Event: event, KeptnHandler: keptnHandler, SLIProviderConfig: K8sSLIProviderConfig{}}, nil // new event type in Keptn versions >= 0.6
	case keptnv1.InternalGetSLIDoneEventType:
		return &EvaluateSLIHandler{Event: event, HTTPClient: &http.Client{}, KeptnHandler: keptnHandler}, nil
	case keptnv1.ConfigureMonitoringEventType:
		return &ConfigureMonitoringHandler{Event: event, KeptnHandler: keptnHandler}, nil
	default:
		logger.Info("received unhandled event type")
		return nil, nil
	}
}
