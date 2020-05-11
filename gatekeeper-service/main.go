package main

import (
	"context"
	"fmt"
	"keptn/gatekeeper-service/pkg/handler"
	"log"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"

	keptnevents "github.com/keptn/go-utils/pkg/lib"

	"github.com/kelseyhightower/envconfig"
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
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	go switchEvent(event)
	return nil
}

func switchEvent(event cloudevents.Event) {
	keptnHandler, err := keptnevents.NewKeptn(&event, keptnevents.KeptnOpts{})
	if err != nil {
		l := keptnevents.NewLogger("", event.Context.GetID(), "gatekeeper-service")
		l.Error("failed to initialize Keptn handler: " + err.Error())
		return
	}
	l := keptnevents.NewLogger(keptnHandler.KeptnContext, event.Context.GetID(), "gatekeeper-service")
	shipyard, err := keptnHandler.GetShipyard()
	if err != nil {
		l.Error("failed to retrieve shipyard: " + err.Error())
		return
	}

	handlers := []handler.Handler{handler.NewEvaluationDoneEventHandler(l),
		handler.NewApprovalTriggeredEventHandler(l),
		handler.NewApprovalFinishedEventHandler(l)}

	unhandled := true
	for _, handler := range handlers {
		if handler.IsTypeHandled(event) {
			unhandled = false
			handler.Handle(event, keptnHandler, shipyard)
		}
	}

	if unhandled {
		l.Error(fmt.Sprintf("Received unexpected keptn event type %s", event.Type()))
	}
}
