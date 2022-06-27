package config

// EnvConfig holds the parsed environment variables
// TODO: add other environment variables supported by Shippy
type EnvConfig struct {
	// ProjectNameMaxSize is the maximum number of characters
	// a Keptn project is allowed to have
	ProjectNameMaxSize int `envconfig:"PROJECT_NAME_MAX_SIZE" default:"200"`
	// ServiceNameMaxSize is the maximum number of characters
	// a Keptn service is allowed to have
	// The limit of 43 characters for a service's name is currently imposed by the helm-service,
	// which, if being used for the CD use case with blue/green deployments generates a helm release called <serviceName>-generated,
	// and helm releases have a maximum length of 53 characters. Therefore, we use this value as a default.
	// If the helm chart generation for blue/green deployments is not needed, and this value is too small, it can be adapted here
	ServiceNameMaxSize int `envconfig:"SERVICE_NAME_MAX_SIZE" default:"43"`
	// AutomaticProvisioningURL is a URL to a REST API to provision
	// git credentials if they are not set by the user
	AutomaticProvisioningURL string `envconfig:"AUTOMATIC_PROVISIONING_URL" default:""`
	// PreStopHookTime is the duration of the preStop hook. The duration defined via this value will be the duration between signaling the
	// termination of the shipyard-controller's pod and the reception of the SIGTERM signal
	PreStopHookTime int `envconfig:"PRE_STOP_HOOK_TIME" default:"5"`
	// ConfigurationSvcEndpoint is the URL of the configuration service
	ConfigurationSvcEndpoint string `envconfig:"CONFIGURATION_SERVICE" default:"http://configuration-service:8080"`
	// EventDispatchIntervalSec is the interval with which the event dispatcher tries to send events
	EventDispatchIntervalSec int `envconfig:"EVENT_DISPATCH_INTERVAL_SEC" default:"10"`
	// SequenceDispatchIntervalSec is the interval with which the sequence dispatcher tries to dispatch sequences
	SequenceDispatchIntervalSec string `envconfig:"SEQUENCE_DISPATCH_INTERVAL_SEC" default:"10s"`
	// TaskStartedWaitDuration is the time the sequence watcher waits before timing out a sequence if there is no .started event for a sent task.triggered event
	TaskStartedWaitDuration string `envconfig:"TASK_STARTED_WAIT_DURATION" default:"10m"`
	// UniformIntegrationTTL is the time after which a uniform integration gets removed from the database if it did not receive a heartbeat signal
	UniformIntegrationTTL string `envconfig:"UNIFORM_INTEGRATION_TTL" default:"1m"`
	// SequenceWatcherInterval is the interval with which the sequence watcher tries to find orphaned tasks
	SequenceWatcherInterval string `envconfig:"SEQUENCE_WATCHER_INTERVAL" default:"1m"`
	// NatsURL is the URL of the nats server
	NatsURL string `envconfig:"NATS_URL" default:"nats://keptn-nats"`
	// LogTTL is the retention period for uniform log entries
	LogTTL string `envconfig:"LOG_TTL" default:"120h"`
	// LogLevel is the log level of the shipyard-controller
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	// DisableLeaderElection allows to disable the leader election
	DisableLeaderElection bool `envconfig:"DISABLE_LEADER_ELECTION" default:"false"`
}
