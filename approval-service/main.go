package main

import (
	"context"
	"keptn/approval-service/pkg/handler"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/sirupsen/logrus"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const envVarLogLevel = "LOG_LEVEL"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

func main() {
	logger.SetLevel(logger.InfoLevel)

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
		log.Fatalf("Failed to process env var: %s", err)
	}

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

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	ctx.Value(gracefulShutdownKey).(*sync.WaitGroup).Add(1)
	val := ctx.Value(gracefulShutdownKey)
	if val != nil {
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Add(1)
		}
	}
	go switchEvent(ctx, event)
	return nil
}

func switchEvent(ctx context.Context, event cloudevents.Event) {
	defer func() {
		logger.Info("Terminating Evaluate-SLI handler")
		val := ctx.Value(gracefulShutdownKey)
		if val == nil {
			return
		}
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Done()
		}
	}()
	keptnHandlerV2, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{})

	if err != nil {
		logger.WithError(err).Error("failed to initialize Keptn handler")
		return
	}

	handlers := []handler.Handler{
		handler.NewApprovalTriggeredEventHandler(keptnHandlerV2),
	}

	unhandled := true
	for _, currHandler := range handlers {
		if currHandler.IsTypeHandled(event) {
			unhandled = false
			currHandler.Handle(event, keptnHandlerV2)
		}
	}

	if unhandled {
		logger.Debugf("Received unexpected keptn event type %s", event.Type())
	}
}

func getGracefulContext() context.Context {

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), gracefulShutdownKey, wg))
	ctx = cloudevents.WithEncodingStructured(ctx)
	go func() {
		<-ch
		logger.Fatal("Container termination triggered, starting graceful shutdown")
		wg.Wait()
		logger.Fatal("cancelling context")
		cancel()
	}()
	return ctx
}
