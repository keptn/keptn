package config

import "strings"

type EnvConfig struct {
	KeptnAPIEndpoint    string `envconfig:"KEPTN_API_ENDPOINT" default:""`
	KeptnAPIToken       string `envconfig:"KEPTN_API_TOKEN" default:""`
	APIProxyPort        int    `envconfig:"API_PROXY_PORT" default:"8081"`
	APIProxyPath        string `envconfig:"API_PROXY_PATH" default:"/"`
	HTTPPollingInterval string `envconfig:"HTTP_POLLING_INTERVAL" default:"10"`
	EventForwardingPath string `envconfig:"EVENT_FORWARDING_PATH" default:"/event"`
	VerifySSL           bool   `envconfig:"HTTP_SSL_VERIFY" default:"true"`
	PubSubURL           string `envconfig:"PUBSUB_URL" default:"nats://keptn-nats-cluster"`
	PubSubTopic         string `envconfig:"PUBSUB_TOPIC" default:""`
	PubSubRecipient     string `envconfig:"PUBSUB_RECIPIENT" default:"http://127.0.0.1"`
	PubSubRecipientPort string `envconfig:"PUBSUB_RECIPIENT_PORT" default:"8080"`
	PubSubRecipientPath string `envconfig:"PUBSUB_RECIPIENT_PATH" default:""`
	ProjectFilter       string `envconfig:"PROJECT_FILTER" default:""`
	StageFilter         string `envconfig:"STAGE_FILTER" default:""`
	ServiceFilter       string `envconfig:"SERVICE_FILTER" default:""`
	DisableRegistration bool   `envconfig:"DISABLE_REGISTRATION" default:"false"`
	Location            string `envconfig:"LOCATION" default:""`
	DistributorVersion  string `envconfig:"DISTRIBUTOR_VERSION" default:""`
	Version             string `envconfig:"VERSION" default:""`
	K8sDeploymentName   string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8sNamespace        string `envconfig:"K8S_NAMESPACE" default:""`
	K8sPodName          string `envconfig:"K8S_POD_NAME" default:""`
	K8sNodeName         string `envconfig:"K8S_NODE_NAME" default:""`
}

func GetPubSubConnectionType(env EnvConfig) ConnectionType {
	if env.KeptnAPIEndpoint == "" {
		// if no Keptn API URL has been defined, this means that run inside the Keptn cluster -> we can subscribe to events directly via NATS
		return ConnectionTypeNATS
	}
	// if a Keptn API URL has been defined, this means that the distributor runs outside of the Keptn cluster -> therefore no NATS connection is possible
	return ConnectionTypeHTTP
}

func GetPubSubRecipientURL(env EnvConfig) string {
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
