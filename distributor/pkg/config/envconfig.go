package config

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
}
