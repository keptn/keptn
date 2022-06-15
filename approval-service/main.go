package main

import (
	"context"
	"github.com/keptn/keptn/cp-connector/pkg/eventsource"
	"github.com/keptn/keptn/cp-connector/pkg/logforwarder"
	"github.com/keptn/keptn/cp-connector/pkg/subscriptionsource"
	"keptn/approval-service/pkg/handler"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-common/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
)

type envConfig struct {
	K8SDeploymentName      string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8SDeploymentVersion   string `envconfig:"K8S_DEPLOYMENT_VERSION" default:""`
	K8SDeploymentComponent string `envconfig:"K8S_DEPLOYMENT_COMPONENT" default:""`
	K8SPodName             string `envconfig:"K8S_POD_NAME" default:""`
	K8SNamespace           string `envconfig:"K8S_NAMESPACE" default:""`
	K8SNodeName            string `envconfig:"K8S_NODE_NAME" default:""`
	LogLevel               string `envconfig:"LOG_LEVEL" default:"info"`
}

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

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

	subscriptionSource := subscriptionsource.New(api.UniformV1())
	natsConnector := nats.NewFromEnv()

	eventSource := eventsource.New(natsConnector)
	logForwarder := logforwarder.New(api.LogsV1())

	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder)

	go func() {
		keptnapi.RunHealthEndpoint("8080", keptnapi.WithReadinessConditionFunc(func() bool {
			return controlPlane.IsRegistered()
		}))
	}()

	ctx, wg := getGracefulContext()
	err = controlPlane.Register(ctx, ApprovalService{env})
	if err != nil {
		log.Fatal(err)
	}

	// this segment will be reached once the context has been cancelled - i.e. due to receiving the SIGTERM signal
	logger.Info("Waiting for approval event handlers to finish")
	// add additional waiting time to ensure the waitGroup has been increased for all events that have been received between receiving SIGTERM and this point
	<-time.After(5 * time.Second)
	wg.Wait()
	logger.Info("All approval handlers finished - exiting")
}

type ApprovalService struct {
	env envConfig
}

func (as ApprovalService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	val := ctx.Value(gracefulShutdownKey)
	if val != nil {
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Add(1)
		}
	}

	defer func() {
		logger.Info("Terminating approval handler")
		val := ctx.Value(gracefulShutdownKey)
		if val == nil {
			return
		}
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Done()
		}
	}()

	ce := v0_2_0.ToCloudEvent(event)

	approvalHandler, err := handler.NewApprovalTriggeredEventHandler(ctx, ce)

	if err != nil {
		log.Println(err.Error())
		return errors.Wrap(controlplane.ErrEventHandleFatal, err.Error())
	}

	if approvalHandler != nil {
		if approvalHandler.IsTypeHandled(ce) {
			return approvalHandler.Handle(ce)
		}
		logger.Debugf("Received unexpected keptn event type %s", ce.Type())
	}

	return nil
}

func (l ApprovalService) RegistrationData() controlplane.RegistrationData {
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
				Event:  "sh.keptn.event.approval.triggered",
				Filter: models.EventSubscriptionFilter{},
			},
		},
	}
}

func getGracefulContext() (context.Context, *sync.WaitGroup) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), gracefulShutdownKey, wg))
	ctx = cloudevents.WithEncodingStructured(ctx)
	go func() {
		<-ch
		logger.Info("Container termination triggered, starting graceful shutdown and cancelling context")
		cancel()
	}()
	return ctx, wg
}
