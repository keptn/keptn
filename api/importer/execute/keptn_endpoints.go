package execute

import (
	"os"

	"github.com/keptn/keptn/api/utils"
)

const defaultControlPlaneEndpoint = "http://shipyard-controller:8080"
const defaultConfigurationServiceEndpoint = "http://configuration-service:8080"
const defaultSecretServiceEndpoint = "http://secret-service:8080"

// StaticKeptnEndpointProvider is a completely static implementation of the KeptnEndpointProvider interface.
// This is meant to be used when we are running in the same namespace as the other keptn services
type StaticKeptnEndpointProvider struct{}

// GetControlPlaneEndpoint returns the default shipyard controller expectedControlPlaneEndpoint
func (_ StaticKeptnEndpointProvider) GetControlPlaneEndpoint() string {
	return utils.SanitizeURL(defaultControlPlaneEndpoint)
}

func NewKeptnEndpointProviderFromEnv() *configurableKeptnEndpointProvider {
	kep := new(configurableKeptnEndpointProvider)

	const controlPlaneServiceEnvVar = "CONTROLPLANE_URI"
	kep.controlPlane = utils.SanitizeURL(
		getEnvVariablewithDefault(
			controlPlaneServiceEnvVar, defaultControlPlaneEndpoint,
		),
	)

	const configurationServiceEnvVar = "CONFIGURATION_URI"
	kep.configurationService = utils.SanitizeURL(
		getEnvVariablewithDefault(
			configurationServiceEnvVar, defaultConfigurationServiceEndpoint,
		),
	)

	const secretServiceEnvVar = "SECRET_SERVICE_URI"
	kep.secretService = utils.SanitizeURL(
		getEnvVariablewithDefault(
			secretServiceEnvVar, defaultSecretServiceEndpoint,
		),
	)

	return kep
}

type configurableKeptnEndpointProvider struct {
	controlPlane         string
	configurationService string
	secretService        string
}

func (kep *configurableKeptnEndpointProvider) GetControlPlaneEndpoint() string {
	return kep.controlPlane
}

func (kep *configurableKeptnEndpointProvider) GetSecretsServiceEndpoint() string {
	return kep.secretService
}

func (kep *configurableKeptnEndpointProvider) GetConfigurationServiceEndpoint() string {
	return kep.configurationService
}

func getEnvVariablewithDefault(envVarName string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVarName)
	if !ok {
		return defaultValue
	}
	return envVarValue
}
