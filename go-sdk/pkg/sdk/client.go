package sdk

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"log"
)

func NewHTTPClientFromEnv() cloudevents.Client {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(env.Port), cloudevents.WithPath(env.Path), cloudevents.WithGetHandlerFunc(api.HealthEndpointHandler))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	return c
}

func NewResourceHandlerFromEnv() *api.ResourceHandler {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return api.NewResourceHandler(env.ConfigurationServiceURL)
}

func NewCPFromEnv() (*controlplane.ControlPlane, controlplane.EventSender, error) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	natsConnector, err := nats.Connect(env.EventBrokerURL)
	if err != nil {
		log.Fatal(err)
	}
	eventSource := controlplane.NewNATSEventSource(natsConnector)
	eventSender := eventSource.Sender()
	controlPlane := controlplane.New(controlplane.NewFixedSubscriptionSource(controlplane.WithFixedSubscriptions(models.EventSubscription{Event: "sh.keptn.>"})), eventSource)
	return controlPlane, eventSender, nil
}

func NewRegistrationDataFromEnv() controlplane.RegistrationData {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return controlplane.RegistrationData{
		Name: "local-service",
		MetaData: models.MetaData{
			Hostname:           "localhost",
			IntegrationVersion: env.Version,
			Location:           env.Location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      env.K8sNamespace,
				PodName:        env.K8sPodName,
				DeploymentName: env.K8sDeploymentName,
			},
		},
		Subscriptions: []models.EventSubscription{},
	}
}
