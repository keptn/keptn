package config

import (
	"context"
	"crypto/tls"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EnvConfig struct {
	KeptnAPIEndpoint          string        `envconfig:"KEPTN_API_ENDPOINT" default:""`
	KeptnAPIToken             string        `envconfig:"KEPTN_API_TOKEN" default:""`
	APIProxyPort              int           `envconfig:"API_PROXY_PORT" default:"8081"`
	APIProxyPath              string        `envconfig:"API_PROXY_PATH" default:"/"`
	APIProxyMaxPayloadBytesKB int           `envconfig:"API_PROXY_MAX_PAYLOAD_BYTES_KB" default:"64"`
	APIProxyHTTPTimeout       string        `envconfig:"API_PROXY_HTTP_TIMEOUT" default:"30"`
	HTTPPollingInterval       string        `envconfig:"HTTP_POLLING_INTERVAL" default:"10"`
	EventForwardingPath       string        `envconfig:"EVENT_FORWARDING_PATH" default:"/event"`
	VerifySSL                 bool          `envconfig:"HTTP_SSL_VERIFY" default:"true"`
	PubSubURL                 string        `envconfig:"PUBSUB_URL" default:"nats://keptn-nats"`
	PubSubTopic               string        `envconfig:"PUBSUB_TOPIC" default:""`
	PubSubRecipient           string        `envconfig:"PUBSUB_RECIPIENT" default:"http://127.0.0.1"`
	PubSubRecipientPort       string        `envconfig:"PUBSUB_RECIPIENT_PORT" default:"8080"`
	PubSubRecipientPath       string        `envconfig:"PUBSUB_RECIPIENT_PATH" default:""`
	PubSubGroup               string        `envconfig:"PUBSUB_GROUP" default:""`
	ProjectFilter             string        `envconfig:"PROJECT_FILTER" default:""`
	StageFilter               string        `envconfig:"STAGE_FILTER" default:""`
	ServiceFilter             string        `envconfig:"SERVICE_FILTER" default:""`
	DisableRegistration       bool          `envconfig:"DISABLE_REGISTRATION" default:"false"`
	RegistrationInterval      string        `envconfig:"REGISTRATION_INTERVAL" default:"10s"`
	Location                  string        `envconfig:"LOCATION" default:""`
	DistributorVersion        string        `envconfig:"DISTRIBUTOR_VERSION" default:"0.9.0"` // TODO: set this automatically
	Version                   string        `envconfig:"VERSION" default:""`
	K8sDeploymentName         string        `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace              string        `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName                string        `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName               string        `envconfig:"K8S_NODE_NAME" default:""`
	MaxHeartBeatRetries       int           `envconfig:"MAX_HEARTBEAT_RETRIES" default:"10"`
	HeartbeatInterval         time.Duration `envconfig:"HEARTBEAT_INTERVAL" default:"10s"`
	MaxRegistrationRetries    int           `envconfig:"MAX_REGISTRATION_RETRIES" default:"10"`
	OAuthClientID             string        `envconfig:"OAUTH_CLIENT_ID" default:""`
	OAuthClientSecret         string        `envconfig:"OAUTH_CLIENT_SECRET" default:""`
	OAuthScopes               []string      `envconfig:"OAUTH_SCOPES" default:""`
	OAuthDiscovery            string        `envconfig:"OAUTH_DISCOVERY" default:""`
	OauthTokenURL             string        `envconfig:"OAUTH_TOKEN_URL" default:""`
}

func (env *EnvConfig) PubSubConnectionType() ConnectionType {
	if env.KeptnAPIEndpoint == "" {
		// if no Keptn API URL has been defined, this means that run inside the Keptn cluster -> we can subscribe to events directly via NATS
		return ConnectionTypeNATS
	}
	// if a Keptn API URL has been defined, this means that the distributor runs outside of the Keptn cluster -> therefore no NATS connection is possible
	return ConnectionTypeHTTP
}

func (env *EnvConfig) ValidateKeptnAPIEndpointURL() error {
	if env.KeptnAPIEndpoint != "" {
		_, err := url.ParseRequestURI(env.KeptnAPIEndpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *EnvConfig) ValidateRegistrationConstraints() bool {
	if env.DisableRegistration {
		logger.Infof("Registration to Keptn's control plane disabled")
		return false
	}

	if env.K8sNamespace == "" || env.K8sDeploymentName == "" {
		logger.Warn("Skipping Registration because not all mandatory environment variables are set: K8S_NAMESPACE, K8S_DEPLOYMENT_NAME")
		return false
	}

	if isOneOfFilteredServices(env.K8sDeploymentName) {
		logger.Infof("Skipping Registration because service name %s is actively filtered", env.K8sDeploymentName)
		return false
	}

	return true
}

func (env *EnvConfig) ProxyHost(path string) (string, string, string) {
	// if the endpoint is empty, redirect to the internal services
	if env.KeptnAPIEndpoint == "" {
		for key, value := range InClusterAPIProxyMappings {
			if strings.HasPrefix(path, key) {
				split := strings.Split(strings.TrimPrefix(path, "/"), "/")
				join := strings.Join(split[1:], "/")
				return "http", value, join
			}
		}
		return "", "", ""
	}

	parsedKeptnURL, err := url.Parse(env.KeptnAPIEndpoint)
	if err != nil {
		return "", "", ""
	}

	// if the endpoint is not empty, map to the correct api
	for key, value := range ExternalAPIProxyMappings {
		if strings.HasPrefix(path, key) {
			split := strings.Split(strings.TrimPrefix(path, "/"), "/")
			join := strings.Join(split[1:], "/")
			path = value + "/" + join
			path = queryEscapeConfigurationServiceURI(path, value)
			if parsedKeptnURL.Path != "" {
				path = strings.TrimSuffix(parsedKeptnURL.Path, "/") + path
			}
			return parsedKeptnURL.Scheme, parsedKeptnURL.Host, path
		}
	}
	return "", "", ""
}

func (env *EnvConfig) OAuthEnabled() bool {
	clientIDAndSecretSet := env.OAuthClientID != "" && env.OAuthClientSecret != ""
	tokenURLOrDiscoverySet := env.OauthTokenURL != "" || env.OAuthDiscovery != ""
	scopesSet := len(env.OAuthScopes) > 0
	return clientIDAndSecretSet && tokenURLOrDiscoverySet && scopesSet
}

func (env *EnvConfig) HTTPPollingEndpoint() string {
	endpoint := env.KeptnAPIEndpoint
	if endpoint == "" {
		if endpoint == "" {
			return DefaultEventsEndpoint
		}
	} else {
		endpoint = strings.TrimSuffix(env.KeptnAPIEndpoint, "/") + "/controlPlane/v1/event/triggered"
	}

	parsedURL, _ := url.Parse(endpoint)

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	if parsedURL.Path == "" {
		parsedURL.Path = "v1/event/triggered"
	}

	return parsedURL.String()
}

func (env *EnvConfig) PubSubRecipientURL() string {
	recipientService := env.PubSubRecipient

	if !strings.HasPrefix(recipientService, "https://") && !strings.HasPrefix(recipientService, "http://") {
		recipientService = "http://" + recipientService
	}

	path := ""
	if env.PubSubRecipientPath != "" {
		path = "/" + strings.TrimPrefix(env.PubSubRecipientPath, "/")
	}
	return recipientService + ":" + env.PubSubRecipientPort + path
}

func (env *EnvConfig) PubSubTopics() []string {
	if env.PubSubTopic == "" {
		return []string{}
	}
	return strings.Split(env.PubSubTopic, ",")
}

func (env *EnvConfig) HTTPClient() *http.Client {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !env.VerifySSL}, //nolint:gosec
		},
		Timeout: env.GetAPIProxyHTTPTimeout(),
	}

	if env.OAuthEnabled() {
		conf := clientcredentials.Config{
			ClientID:     env.OAuthClientID,
			ClientSecret: env.OAuthClientSecret,
			Scopes:       env.OAuthScopes,
			TokenURL:     env.OauthTokenURL,
		}
		return conf.Client(context.WithValue(context.TODO(), oauth2.HTTPClient, c))
	}
	return c
}

func (env *EnvConfig) GetAPIProxyHTTPTimeout() time.Duration {
	timeout, err := strconv.ParseInt(env.APIProxyHTTPTimeout, 10, 64)
	if err != nil {
		timeout = DefaultAPIProxyHTTPTimeout
	}
	return time.Duration(timeout) * time.Second
}

func (env *EnvConfig) GetAPIProxyMaxBytes() int64 {
	return int64(env.APIProxyMaxPayloadBytesKB << 10)
}

func queryEscapeConfigurationServiceURI(path string, value string) string {
	// special case: configuration service /resource requests with nested resource URIs need to have an escaped '/' - see https://github.com/keptn/keptn/issues/2707
	if value == "/resource-service" {
		splitPath := strings.Split(path, "/resource/")
		if len(splitPath) > 1 {
			path = ""
			for i := 0; i < len(splitPath)-1; i++ {
				path = splitPath[i] + "/resource/"
			}
			path += url.QueryEscape(splitPath[len(splitPath)-1])
		}
	}
	return path
}
func isOneOfFilteredServices(serviceName string) bool {
	switch serviceName {
	case
		"statistics-service",
		"api-service",
		"mongodb-datastore",
		"resource-service",
		"secret-service",
		"shipyard-controller":
		return true
	}
	return false
}
