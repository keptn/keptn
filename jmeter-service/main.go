package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// JMeterServiceName is the name of the JMeter Keptn Service
const JMeterServiceName = "jmeter-service"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port     int    `envconfig:"RCV_PORT" default:"8080"`
	Path     string `envconfig:"RCV_PATH" default:"/"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logger.Fatalf("Failed to process env var: %s", err)
	}

	logger.SetLevel(logger.InfoLevel)

	if os.Getenv(env.LogLevel) != "" {
		logLevel, err := logger.ParseLevel(os.Getenv(env.LogLevel))
		if err != nil {
			logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			logger.SetLevel(logLevel)
		}
	}

	ctx := getGracefulContext()
	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(keptnapi.HealthEndpointHandler))
	if err != nil {
		logger.Fatalf("Failed to create cloud event client: %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		logger.Fatalf("Failed to create cloud event client: %v", err)
	}

	eventSender, err := keptnv2.NewHTTPEventSender("")
	if err != nil {
		logger.Fatalf("Failed to create event sender: %v", err)
	}

	eventHandler := &EventHandler{testRunner: NewTestRunner(eventSender)}

	logger.Fatal(c.StartReceiver(ctx, eventHandler.handleEvent))
	ctx.Value(gracefulShutdownKey).(*sync.WaitGroup).Wait()
	close(ctx.Value(keptnQuit).(chan os.Signal))

}

func getGracefulContext() context.Context {

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), gracefulShutdownKey, wg))
	ctx = cloudevents.WithEncodingStructured(ctx)
	ctx = context.WithValue(ctx, keptnQuit, ch)
	go func() {
		<-ch
		close(ch)
		logger.Fatal("Container termination triggered, starting graceful shutdown")
		ctx.Value(gracefulShutdownKey).(*sync.WaitGroup).Wait()
		logger.Fatal("cancelling context")
		cancel()
	}()
	return ctx
}
