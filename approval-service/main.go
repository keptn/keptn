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
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
)

const envVarLogLevel = "LOG_LEVEL"

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
	err = controlPlane.Register(ctx, ApprovalService{env})
	if err != nil {
		log.Fatal(err)
	}
}

type ApprovalService struct {
	env envConfig
}

func (as ApprovalService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	ctx.Value(gracefulShutdownKey).(*sync.WaitGroup).Add(1)
	val := ctx.Value(gracefulShutdownKey)
	if val != nil {
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Add(1)
		}
	}

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
