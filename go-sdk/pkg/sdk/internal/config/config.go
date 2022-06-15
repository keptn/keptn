package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"strconv"
	"time"
)

type EnvConfig struct {
	APIProxyHTTPTimeout     string   `EnvConfig:"API_PROXY_HTTP_TIMEOUT" default:"30"`
	ConfigurationServiceURL string   `EnvConfig:"CONFIGURATION_SERVICE" default:"configuration-service:8080"`
	EventBrokerURL          string   `EnvConfig:"EVENTBROKER" default:"nats://keptn-nats"`
	PubSubTopic             string   `EnvConfig:"PUBSUB_TOPIC" default:""`
	HealthEndpointPort      string   `EnvConfig:"HEALTH_ENDPOINT_PORT" default:"8080"`
	HealthEndpointEnabled   bool     `EnvConfig:"HEALTH_ENDPOINT_ENABLED" default:"true"`
	KeptnAPIEndpoint        string   `EnvConfig:"KEPTN_API_ENDPOINT" default:""`
	KeptnAPIToken           string   `EnvConfig:"KEPTN_API_TOKEN" default:""`
	Location                string   `EnvConfig:"LOCATION" default:"control-plane"`
	K8sDeploymentVersion    string   `EnvConfig:"K8S_DEPLOYMENT_VERSION" default:""`
	K8sDeploymentName       string   `EnvConfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace            string   `EnvConfig:"K8S_NAMESPACE" default:""`
	K8sPodName              string   `EnvConfig:"K8S_POD_NAME" default:""`
	K8sNodeName             string   `EnvConfig:"K8S_NODE_NAME" default:""`
	OAuthClientID           string   `EnvConfig:"OAUTH_CLIENT_ID" default:""`
	OAuthClientSecret       string   `EnvConfig:"OAUTH_CLIENT_SECRET" default:""`
	OAuthScopes             []string `EnvConfig:"OAUTH_SCOPES" default:""`
	OAuthDiscovery          string   `EnvConfig:"OAUTH_DISCOVERY" default:""`
	OauthTokenURL           string   `EnvConfig:"OAUTH_TOKEN_URL" default:""`
	VerifySSL               bool     `EnvConfig:"HTTP_SSL_VERIFY" default:"true"`
}

type ConnectionType string

const (
	DefaultAPIProxyHTTPTimeout                = 30
	ConnectionTypeNATS         ConnectionType = "nats"
	ConnectionTypeHTTP         ConnectionType = "http"
)

func NewEnvConfig() EnvConfig {
	var env EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return env
}

func (env *EnvConfig) OAuthEnabled() bool {
	clientIDAndSecretSet := env.OAuthClientID != "" && env.OAuthClientSecret != ""
	tokenURLOrDiscoverySet := env.OauthTokenURL != "" || env.OAuthDiscovery != ""
	scopesSet := len(env.OAuthScopes) > 0
	return clientIDAndSecretSet && tokenURLOrDiscoverySet && scopesSet
}

func (env *EnvConfig) GetAPIProxyHTTPTimeout() time.Duration {
	timeout, err := strconv.ParseInt(env.APIProxyHTTPTimeout, 10, 64)
	if err != nil {
		timeout = DefaultAPIProxyHTTPTimeout
	}
	return time.Duration(timeout) * time.Second
}

func (env *EnvConfig) PubSubConnectionType() ConnectionType {
	if env.KeptnAPIEndpoint == "" {
		// if no Keptn API URL has been defined, this means that run inside the Keptn cluster -> we can subscribe to events directly via NATS
		return ConnectionTypeNATS
	}
	// if a Keptn API URL has been defined, this means that the distributor runs outside of the Keptn cluster -> therefore no NATS connection is possible
	return ConnectionTypeHTTP
}
