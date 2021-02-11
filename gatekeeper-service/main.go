package main

import (
	"context"
	"fmt"
	"keptn/gatekeeper-service/pkg/handler"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
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

	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	go switchEvent(event)
	return nil
}

func switchEvent(event cloudevents.Event) {
	serviceName := "gatekeeper-service"
	keptnHandlerV2, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{ServiceName: &serviceName},
	})

	if err != nil {
		l := keptncommon.NewLogger("", event.Context.GetID(), "gatekeeper-service")
		l.Error("failed to initialize Keptn handler: " + err.Error())
		return
	}

	handlers := []handler.Handler{
		handler.NewApprovalTriggeredEventHandler(keptnHandlerV2),
	}

	unhandled := true
	for _, handler := range handlers {
		if handler.IsTypeHandled(event) {
			unhandled = false
			handler.Handle(event, keptnHandlerV2)
		}
	}

	if unhandled {
		keptnHandlerV2.Logger.Error(fmt.Sprintf("Received unexpected keptn event type %s", event.Type()))
	}
}
