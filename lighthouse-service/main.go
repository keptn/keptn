package main

import (
	"context"
	models "github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/lib-cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/lib-cp-connector/pkg/nats"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
	"github.com/pkg/errors"
	"log"
)

const envVarLogLevel = "LOG_LEVEL"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	o := controlplane.ControlPlaneOptions{
		KeptnAPIEndpoint: "",
		KeptnAPIToken:    "",
		NATSEndpoint:     "",
	}

	keptnAPI, err := api.New(o.KeptnAPIEndpoint, api.WithAuthToken(o.KeptnAPIToken))
	if err != nil {
		log.Fatal(err)
	}

	subscriptionSource := controlplane.NewSubscriptionSource(keptnAPI.UniformV1())
	natsConnector, err := nats.Connect(o.NATSEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	eventSource := controlplane.NewNATSEventSource(natsConnector)

	controlPlane := controlplane.New(subscriptionSource, eventSource)
	err = controlPlane.Register(context.TODO(), LighthouseService{})
	if err != nil {
		log.Fatal(err)
	}
}

type LighthouseService struct{}

func (l LighthouseService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {

	ce := v0_2_0.ToCloudEvent(event)
	handler, err := event_handler.NewEventHandler(ctx, ce)

	if err != nil {
		return errors.Wrap(controlplane.ErrEventHandleIgnore, err.Error())
	}
	if handler != nil {
		return handler.HandleEvent(ctx)
	}

	return nil
}

func (l LighthouseService) RegistrationData() controlplane.RegistrationData {
	return controlplane.RegistrationData{
		Name: "lh-wuppi",
		MetaData: models.MetaData{
			Hostname:           "localhost",
			IntegrationVersion: "dev",
			DistributorVersion: "0.15.0",
			Location:           "local",
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      "keptn",
				PodName:        "lh-wuppi",
				DeploymentName: "lh-wuppi",
			},
		},
	}
}

/*func main2() {
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

	val := ctx.Value(event_handler.GracefulShutdownKey)
	if val != nil {
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Wait()
		}
	}
	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {

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
		wg.Wait()
		cancel()
	}()

	return ctx
}*/
