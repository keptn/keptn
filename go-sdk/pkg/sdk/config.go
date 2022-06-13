package sdk

import (
	"log"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	APIProxyHTTPTimeout     string   `envconfig:"API_PROXY_HTTP_TIMEOUT" default:"30"`
	ConfigurationServiceURL string   `envconfig:"CONFIGURATION_SERVICE" default:"configuration-service:8080"`
	EventBrokerURL          string   `envconfig:"EVENTBROKER" default:"nats://keptn-nats"`
	PubSubTopic             string   `envconfig:"PUBSUB_TOPIC" default:""`
	HealthEndpointPort      string   `envconfig:"HEALTH_ENDPOINT_PORT" default:"8080"`
	HealthEndpointEnabled   bool     `envconfig:"HEALTH_ENDPOINT_ENABLED" default:"true"`
	InternalAPIEnabled      bool     `envconfig:"INTERNAL_API_ENABLED" default:"true"`
	KeptnAPIEndpoint        string   `envconfig:"KEPTN_API_ENDPOINT" default:""`
	KeptnAPIToken           string   `envconfig:"KEPTN_API_TOKEN" default:""`
	Location                string   `envconfig:"LOCATION" default:"control-plane"`
	K8sDeploymentVersion    string   `envconfig:"K8S_DEPLOYMENT_VERSION" default:""`
	K8sDeploymentName       string   `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace            string   `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName              string   `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName             string   `envconfig:"K8S_NODE_NAME" default:""`
	OAuthClientID           string   `envconfig:"OAUTH_CLIENT_ID" default:""`
	OAuthClientSecret       string   `envconfig:"OAUTH_CLIENT_SECRET" default:""`
	OAuthScopes             []string `envconfig:"OAUTH_SCOPES" default:""`
	OAuthDiscovery          string   `envconfig:"OAUTH_DISCOVERY" default:""`
	OauthTokenURL           string   `envconfig:"OAUTH_TOKEN_URL" default:""`
	VerifySSL               bool     `envconfig:"HTTP_SSL_VERIFY" default:"true"`
}

const (
	DefaultAPIProxyHTTPTimeout = 30
)

func newEnvConfig() envConfig {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return env
}

func (env *envConfig) OAuthEnabled() bool {
	clientIDAndSecretSet := env.OAuthClientID != "" && env.OAuthClientSecret != ""
	tokenURLOrDiscoverySet := env.OauthTokenURL != "" || env.OAuthDiscovery != ""
	scopesSet := len(env.OAuthScopes) > 0
	return clientIDAndSecretSet && tokenURLOrDiscoverySet && scopesSet
}

func (env *envConfig) GetAPIProxyHTTPTimeout() time.Duration {
	timeout, err := strconv.ParseInt(env.APIProxyHTTPTimeout, 10, 64)
	if err != nil {
		timeout = DefaultAPIProxyHTTPTimeout
	}
	return time.Duration(timeout) * time.Second
}
