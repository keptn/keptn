package sdk

type envConfig struct {
	Port                    int    `envconfig:"RCV_PORT" default:"8080"`
	Path                    string `envconfig:"RCV_PATH" default:"/"`
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"configuration-service:8080"`
	Location                string `envconfig:"LOCATION" default:""`
	DistributorVersion      string `envconfig:"DISTRIBUTOR_VERSION" default:"0.9.0"` // TODO: set this automatically
	Version                 string `envconfig:"VERSION" default:""`
	K8sDeploymentName       string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace            string `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName              string `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName             string `envconfig:"K8S_NODE_NAME" default:""`
}
