package main

import (
	"context"
	"encoding/json"
	"fmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"keptn/gatekeeper-service/pkg/handler"
	"log"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

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
	go keptnapi.RunHealthEndpoint("10999")
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
	serviceName := "gatekeeper-service"
	keptnHandlerV2, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{ServiceName: &serviceName},
	})

	if err != nil {
		l := keptncommon.NewLogger("", event.Context.GetID(), "gatekeeper-service")
		l.Error("failed to initialize Keptn handler: " + err.Error())
		return
	}

	// TODO: Retrieving the shipyard file will become obsolete because required properties will be located in the event

	shipyard := &keptn.Shipyard{}
	shipyardResource, err := keptnHandlerV2.ResourceHandler.GetProjectResource(keptnHandlerV2.KeptnBase.Event.GetProject(), "shipyard.yaml")
	if err != nil {
		keptnHandlerV2.Logger.Error("failed to retrieve shipyard: " + err.Error())
		return
	}
	err = json.Unmarshal([]byte(shipyardResource.ResourceContent), shipyard)
	if err != nil {
		keptnHandlerV2.Logger.Error("failed to decode shipyard: " + err.Error())
		return
	}

	handlers := []handler.Handler{handler.NewEvaluationDoneEventHandler(keptnHandlerV2),
		handler.NewApprovalTriggeredEventHandler(keptnHandlerV2),
		handler.NewApprovalFinishedEventHandler(keptnHandlerV2)}

	unhandled := true
	for _, handler := range handlers {
		if handler.IsTypeHandled(event) {
			unhandled = false
			handler.Handle(event, keptnHandlerV2, shipyard)
		}
	}

	if unhandled {
		keptnHandlerV2.Logger.Error(fmt.Sprintf("Received unexpected keptn event type %s", event.Type()))
	}
}
