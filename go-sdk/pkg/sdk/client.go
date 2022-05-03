package sdk

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
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

func NewControlPlaneFromEnv() controlplane.ControlPlane {
	// TODO
	return nil
}
