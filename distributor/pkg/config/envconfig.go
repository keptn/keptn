package config

import (
	logger "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

var Global EnvConfig

type EnvConfig struct {
	KeptnAPIEndpoint     string `envconfig:"KEPTN_API_ENDPOINT" default:""`
	KeptnAPIToken        string `envconfig:"KEPTN_API_TOKEN" default:""`
	APIProxyPort         int    `envconfig:"API_PROXY_PORT" default:"8081"`
	APIProxyPath         string `envconfig:"API_PROXY_PATH" default:"/"`
	HTTPPollingInterval  string `envconfig:"HTTP_POLLING_INTERVAL" default:"10"`
	EventForwardingPath  string `envconfig:"EVENT_FORWARDING_PATH" default:"/event"`
	VerifySSL            bool   `envconfig:"HTTP_SSL_VERIFY" default:"true"`
	PubSubURL            string `envconfig:"PUBSUB_URL" default:"nats://keptn-nats-cluster"`
	PubSubTopic          string `envconfig:"PUBSUB_TOPIC" default:""`
	PubSubRecipient      string `envconfig:"PUBSUB_RECIPIENT" default:"http://127.0.0.1"`
	PubSubRecipientPort  string `envconfig:"PUBSUB_RECIPIENT_PORT" default:"8080"`
	PubSubRecipientPath  string `envconfig:"PUBSUB_RECIPIENT_PATH" default:""`
	ProjectFilter        string `envconfig:"PROJECT_FILTER" default:""`
	StageFilter          string `envconfig:"STAGE_FILTER" default:""`
	ServiceFilter        string `envconfig:"SERVICE_FILTER" default:""`
	DisableRegistration  bool   `envconfig:"DISABLE_REGISTRATION" default:"false"`
	RegistrationInterval string `envconfig:"REGISTRATION_INTERVAL" default:"10s"`
	Location             string `envconfig:"LOCATION" default:""`
	DistributorVersion   string `envconfig:"DISTRIBUTOR_VERSION" default:"0.9.0"` // TODO: set this automatically
	Version              string `envconfig:"VERSION" default:""`
	K8sDeploymentName    string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace         string `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName           string `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName          string `envconfig:"K8S_NODE_NAME" default:""`
}

func GetRegistrationInterval(env EnvConfig) time.Duration {
	duration, err := time.ParseDuration(env.RegistrationInterval)
	if err != nil {
		logger.Warnf("Unable to parse REGISTRATION_INTERVAL environment variable as duration: %s", env.RegistrationInterval)
		return 10 * time.Second
	}
	return duration
}

func GetPubSubConnectionType() ConnectionType {
	if Global.KeptnAPIEndpoint == "" {
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
func (env *EnvConfig) GetProxyHost(path string) (string, string, string) {
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
			// special case: configuration service /resource requests with nested resource URIs need to have an escaped '/' - see https://github.com/keptn/keptn/issues/2707
			if value == "/configuration-service" {
				splitPath := strings.Split(path, "/resource/")
				if len(splitPath) > 1 {
					path = ""
					for i := 0; i < len(splitPath)-1; i++ {
						path = splitPath[i] + "/resource/"
					}
					path += url.QueryEscape(splitPath[len(splitPath)-1])
				}
			}
			if parsedKeptnURL.Path != "" {
				path = strings.TrimSuffix(parsedKeptnURL.Path, "/") + path
			}
			return parsedKeptnURL.Scheme, parsedKeptnURL.Host, path
		}
	}
	return "", "", ""
}

func (env *EnvConfig) GetHTTPPollingEndpoint() string {
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

func (env *EnvConfig) GetPubSubRecipientURL() string {
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
