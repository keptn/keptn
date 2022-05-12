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
}
