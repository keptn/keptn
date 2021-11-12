package main

import (
	"context"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/sirupsen/logrus"

	"github.com/keptn/keptn/lighthouse-service/event_handler"
)

const envVarLogLevel = "LOG_LEVEL"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	if os.Getenv(envVarLogLevel) != "" {
		logLevel, err := logger.ParseLevel(os.Getenv(envVarLogLevel))
		if err != nil {
			logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			logger.SetLevel(logLevel)
		}
	}

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logger.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := getGracefulContext()

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(keptnapi.HealthEndpointHandler))
	if err != nil {
		logger.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		logger.Fatalf("failed to create client, %v", err)
	}
	logger.Fatal(c.StartReceiver(ctx, gotEvent))

	ctx.Value(event_handler.GracefulShutdownKey).(*sync.WaitGroup).Wait()
	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	_ = event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	handler, err := event_handler.NewEventHandler(event)

	if err != nil {
		logger.Error("Received unknown event type: " + event.Type())
		return err
	}
	if handler != nil {
		return handler.HandleEvent(ctx)
	}

	return nil
}

// storing wait group into context to sync before shutdown
func getGracefulContext() context.Context {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), event_handler.GracefulShutdownKey, wg)))

	go func() {
		<-ch
		logger.Fatal("Container termination triggered, waiting for graceful shutdown")
		cancel()
	}()

	return ctx
}
