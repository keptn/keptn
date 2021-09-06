package event_handler

import (
	"context"
	"net/http"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"go.opentelemetry.io/otel/trace"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	keptnObs "github.com/keptn/go-utils/pkg/common/observability"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type EvaluationEventHandler interface {
	HandleEvent() error
}

func NewEventHandler(ctx context.Context, event cloudevents.Event, logger *keptncommon.Logger) (EvaluationEventHandler, error) {
	logger.Debug("Received event: " + event.Type())
	serviceName := "lighthouse-service"

	// generate a new context from the incoming to avoid it cancelling the routines
	newCtx := getContextWithTraceContext(ctx, event)

	// TODO: We pass a custom httpSender to the keptnhandler
	// because we need to send down the context in order to get proper OTel propagation.
	// Once the methods accept a context parameter, we can remove this code.
	httpSender, err := keptnv2.NewHTTPEventSender("")
	if err != nil {
		return nil, err
	}
	httpSender.Context = newCtx

	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{ServiceName: &serviceName},
		EventSender:    *httpSender,
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

func getContextWithTraceContext(ctx context.Context, event cloudevents.Event) context.Context {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		// the ctx here comes from cloudevents and it has a cancelation
		// so we need to copy the tracecontext to a new context to avoid
		// the goroutines from being cancelled
		return trace.ContextWithSpan(context.Background(), span)
	}

	return keptnObs.ExtractDistributedTracingExtension(context.Background(), event)
}
