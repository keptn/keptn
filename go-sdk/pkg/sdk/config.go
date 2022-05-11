package sdk

type envConfig struct {
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"configuration-service:8080"`
	Location                string `envconfig:"LOCATION" default:"control-plane"`
	K8sDeploymentName       string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace            string `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName              string `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName             string `envconfig:"K8S_NODE_NAME" default:""`
	EventBrokerURL          string `envconfig:"EVENT_BROKER_URL" default:"nats://keptn-nats"`
}
