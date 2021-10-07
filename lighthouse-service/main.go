package main

import (
	"context"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
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

	go keptnapi.RunHealthEndpoint("10998")
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := getGracefulContext()

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(keptnapi.HealthEndpointHandler))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(ctx, gotEvent))
	ctx.Value("Wg").(*sync.WaitGroup).Wait()
	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	_ = event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptncommon.NewLogger(shkeptncontext, event.Context.GetID(), "lighthouse-service")

	handler, err := event_handler.NewEventHandler(event, logger)

	if err != nil {
		logger.Error("Received unknown event type: " + event.Type())
		return err
	}
	if handler != nil {
		return handler.HandleEvent(ctx)
	}

	return nil
}

func getGracefulContext() context.Context {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), "Wg", wg)))

	go func() {
		<-ch
		log.Fatal("Container termination triggered, waiting for graceful shutdown")
		cancel()
	}()

	return ctx
}
