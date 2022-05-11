package sdk

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type envConfig struct {
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"configuration-service:8080"`
	Location                string `envconfig:"LOCATION" default:"control-plane"`
	Version                 string `envconfig:"VERSION" default:""`
	K8sDeploymentName       string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace            string `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName              string `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName             string `envconfig:"K8S_NODE_NAME" default:""`
	EventBrokerURL          string `envconfig:"EVENT_BROKER_URL" default:"nats://keptn-nats"`
}

func newEnvConfig() envConfig {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return env
}
