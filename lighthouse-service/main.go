package main

import (
	"context"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/eventsource"
	"github.com/keptn/keptn/cp-connector/pkg/logforwarder"
	"github.com/keptn/keptn/cp-connector/pkg/subscriptionsource"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-common/api"
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
	log := logger.New()
	logLevel, err := logger.ParseLevel(env.LogLevel)
	if err != nil {
		log.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		log.SetLevel(logger.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}

	api, err := api.NewInternal(nil)
	if err != nil {
		log.Fatal(err)
	}

	subscriptionSource := subscriptionsource.New(api.UniformV1())
	subscriptionsource.WithLogger(log)

	natsConnector := nats.NewFromEnv()
	nats.WithLogger(log)

	eventSource := eventsource.New(natsConnector)
	eventsource.WithLogger(log)

	logForwarder := logforwarder.New(api.LogsV1())
	logforwarder.WithLogger(log)

	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder)
	controlplane.WithLogger(log)

	go func() {
		keptnapi.RunHealthEndpoint("8080", keptnapi.WithReadinessConditionFunc(func() bool {
			return controlPlane.IsRegistered()
		}))
	}()

	ctx, wg := getGracefulContext()
	err = controlPlane.Register(ctx, LighthouseService{env})
	if err != nil {
		log.Fatal(err)
	}

	// this segment will be reached once the context has been cancelled - i.e. due to receiving the SIGTERM signal
	log.Info("Waiting for evaluation event handlers to finish")
	// add additional waiting time to ensure the waitGroup has been increased for all events that have been received between receiving SIGTERM and this point
	<-time.After(5 * time.Second)
	wg.Wait()
	logger.Info("All evaluation handlers finished - exiting")
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
		Name: l.env.K8SDeploymentComponent,
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

func getGracefulContext() (context.Context, *sync.WaitGroup) {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), event_handler.GracefulShutdownKey, wg)))

	go func() {
		<-ch
		logger.Info("Container termination triggered, starting graceful shutdown and cancelling context")
		cancel()
	}()

	return ctx, wg
}
