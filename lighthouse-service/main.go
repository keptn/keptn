package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
	"github.com/pkg/errors"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port                    int    `envconfig:"RCV_PORT" default:"8080"`
	Path                    string `envconfig:"RCV_PATH" default:"/"`
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"http://configuration-service:8080"`
	K8SDeploymentName       string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8SDeploymentVersion    string `envconfig:"K8S_DEPLOYMENT_VERSION" default:""`
	K8SDeploymentComponent  string `envconfig:"K8S_DEPLOYMENT_COMPONENT" default:""`
	K8SPodName              string `envconfig:"K8S_POD_NAME" default:""`
	K8SNamespace            string `envconfig:"K8S_NAMESPACE" default:""`
	K8SNodeName             string `envconfig:"K8S_NODE_NAME" default:""`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	go func() {
		keptnapi.RunHealthEndpoint("8080")
	}()

	api, err := api.NewInternal(nil)
	if err != nil {
		log.Fatal(err)
	}

	subscriptionSource := controlplane.NewUniformSubscriptionSource(api.UniformV1())
	natsConnector, err := nats.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	eventSource := controlplane.NewNATSEventSource(natsConnector)

	controlPlane := controlplane.New(subscriptionSource, eventSource)
	ctx := getGracefulContext()
	err = controlPlane.Register(ctx, LighthouseService{env})
	if err != nil {
		log.Fatal(err)
	}

	val := ctx.Value(event_handler.GracefulShutdownKey)
	if val != nil {
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Wait()
		}
	}
}

type LighthouseService struct {
	env envConfig
}

func (l LighthouseService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	ce := v0_2_0.ToCloudEvent(event)
	handler, err := event_handler.NewEventHandler(ctx, ce)

	if err != nil {
		log.Println(err.Error())
		return errors.Wrap(controlplane.ErrEventHandleFatal, err.Error())
	}
	if handler != nil {
		return handler.HandleEvent(ctx)
	}

	return nil
}

func (l LighthouseService) RegistrationData() controlplane.RegistrationData {
	return controlplane.RegistrationData{
		Name: l.env.K8SPodName,
		MetaData: models.MetaData{
			Hostname:           l.env.K8SNodeName,
			IntegrationVersion: l.env.K8SDeploymentVersion,
			Location:           l.env.K8SDeploymentComponent,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      l.env.K8SNamespace,
				PodName:        l.env.K8SPodName,
				DeploymentName: l.env.K8SDeploymentName,
			},
		},
		Subscriptions: []models.EventSubscription{
			{
				Event:  "sh.keptn.event.evaluation.triggered",
				Filter: models.EventSubscriptionFilter{},
			},
			{
				Event:  "sh.keptn.event.get-sli.finished",
				Filter: models.EventSubscriptionFilter{},
			},
			{
				Event:  "sh.keptn.event.monitoring.configure",
				Filter: models.EventSubscriptionFilter{},
			},
		},
	}
}

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
}
