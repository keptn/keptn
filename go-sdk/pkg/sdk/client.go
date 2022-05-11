package sdk

import (
	"github.com/kelseyhightower/envconfig"
	api "github.com/keptn/go-utils/pkg/api/utils"
	api2 "github.com/keptn/keptn/cp-connector/pkg/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"log"
)

func newResourceHandlerFromEnv() *api.ResourceHandler {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return api.NewResourceHandler(env.ConfigurationServiceURL)
}

func newControlPlane() (*controlplane.ControlPlane, controlplane.EventSender) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return newControlPlaneFromEnv(env)
}
func newControlPlaneFromEnv(env envConfig) (*controlplane.ControlPlane, controlplane.EventSender) {
	apiSet, err := api2.NewInternal(nil)
	if err != nil {
		log.Fatal(err)
	}

	natsConnector, err := nats.Connect(env.EventBrokerURL)
	if err != nil {
		log.Fatal(err)
	}

	eventSource := controlplane.NewNATSEventSource(natsConnector)
	eventSender := eventSource.Sender()

	subscriptionSource := controlplane.NewUniformSubscriptionSource(apiSet.UniformV1())
	controlPlane := controlplane.New(subscriptionSource, eventSource)
	return controlPlane, eventSender
}
