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
	logger "github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
	"github.com/pkg/errors"
)

type envConfig struct {
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"http://configuration-service:8080"`
	K8SDeploymentName       string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8SDeploymentVersion    string `envconfig:"K8S_DEPLOYMENT_VERSION" default:""`
	K8SDeploymentComponent  string `envconfig:"K8S_DEPLOYMENT_COMPONENT" default:""`
	K8SPodName              string `envconfig:"K8S_POD_NAME" default:""`
	K8SNamespace            string `envconfig:"K8S_NAMESPACE" default:""`
	K8SNodeName             string `envconfig:"K8S_NODE_NAME" default:""`
	LogLevel                string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	logLevel, err := logger.ParseLevel(env.LogLevel)
	if err != nil {
		logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		logger.SetLevel(logger.InfoLevel)
	} else {
		logger.SetLevel(logLevel)
	}

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
	logForwarder := controlplane.NewLogForwarder(api.LogsV1())

	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder)

	go func() {
		keptnapi.RunHealthEndpoint("8080", keptnapi.WithReadinessConditionFunc(func() bool {
			return controlPlane.IsRegistered()
		}))
	}()

	ctx := getGracefulContext()
	err = controlPlane.Register(ctx, LighthouseService{env})
	if err != nil {
		log.Fatal(err)
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
		Name: l.env.K8SDeploymentName,
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
		logger.Info("Container termination triggered, starting graceful shutdown and cancelling context")
		cancel()
		logger.Info("Waiting for event handlers to finish")
		wg.Wait()
		logger.Info("All handlers finished - ready to shut down")
	}()

	return ctx
}
