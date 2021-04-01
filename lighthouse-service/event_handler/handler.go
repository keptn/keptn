package event_handler

import (
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
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

	configurationServiceEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	if err != nil {
		return nil, err
	}
	resourceHandler := keptnapi.NewResourceHandler(configurationServiceEndpoint.String())
	serviceHandler := keptnapi.NewServiceHandler(configurationServiceEndpoint.String())

	switch event.Type() {
	case keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName):
		return &StartEvaluationHandler{
			Event:             event,
			KeptnHandler:      keptnHandler,
			SLIProviderConfig: K8sSLIProviderConfig{},
			SLOFileRetriever: SLOFileRetriever{
				ResourceHandler: resourceHandler,
				ServiceHandler:  serviceHandler,
			},
		}, nil
	case keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName):
		return &EvaluateSLIHandler{
			Event:        event,
			HTTPClient:   &http.Client{},
			KeptnHandler: keptnHandler,
			SLOFileRetriever: SLOFileRetriever{
				ResourceHandler: resourceHandler,
				ServiceHandler:  serviceHandler,
			},
			EventStore: keptnHandler.EventHandler,
		}, nil
	case keptn.ConfigureMonitoringEventType:
		return &ConfigureMonitoringHandler{Event: event, KeptnHandler: keptnHandler}, nil
	default:
		logger.Info("received unhandled event type")
		return nil, nil
	}
}
