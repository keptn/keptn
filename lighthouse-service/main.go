package main

import (
	"context"
	"log"
	"net/http"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnObs "github.com/keptn/go-utils/pkg/common/observability"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	serviceName = "lighthouse-service"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	shutdown := keptnObs.InitOTelTraceProvider(serviceName)
	defer shutdown()

	go keptnapi.RunHealthEndpoint("10998")
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP(
		cloudevents.WithPath(env.Path),
		cloudevents.WithPort(env.Port),

		// the middleware will ensure that the traceparent is injected into the context
		// that is passed to the StartReceiver handler func
		// https://github.com/cloudevents/sdk-go/pull/708
		cloudevents.WithMiddleware(func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "receive")
		}),
	)

	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	// the observability service will start a span for each call to `gotEvent`, adding the event data into the span
	c, err := cloudevents.NewClient(p, client.WithObservabilityService(keptnObs.NewOTelObservabilityService()))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	_ = event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptncommon.NewLogger(shkeptncontext, event.Context.GetID(), "lighthouse-service")

	handler, err := event_handler.NewEventHandler(ctx, event, logger)

	if err != nil {
		logger.Error("Received unknown event type: " + event.Type())
		return err
	}
	if handler != nil {
		return handler.HandleEvent()
	}

	return nil
}
