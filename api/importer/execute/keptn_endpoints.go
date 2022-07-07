package execute

import (
	"os"

	"github.com/keptn/keptn/api/utils"
)

const defaultControlPlaneEndpoint = "http://shipyard-controller:8080"
const defaultSecretServiceEndpoint = "http://secret-service:8080"

// StaticKeptnEndpointProvider is a completely static implementation of the KeptnEndpointProvider interface.
// This is meant to be used when we are running in the same namespace as the other keptn services
type StaticKeptnEndpointProvider struct{}

// GetControlPlaneEndpoint returns the default shipyard controller expectedControlPlaneEndpoint
func (_ StaticKeptnEndpointProvider) GetControlPlaneEndpoint() string {
	return utils.SanitizeURL(defaultControlPlaneEndpoint)
}

// GetSecretsServiceEndpoint returns the default secrets service endpoint
func (_ StaticKeptnEndpointProvider) GetSecretsServiceEndpoint() string {
	return utils.SanitizeURL(defaultSecretServiceEndpoint)
}

type configurableKeptnEndpointProvider struct {
	controlPlane  string
	secretService string
}

func getEnvVariablewithDefault(envVarName string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVarName)
	if !ok {
		return defaultValue
	}
	return envVarValue
}

func KeptnEndpointProviderFromEnv() *configurableKeptnEndpointProvider {

	const controlPlaneServiceEnvVar = "CONTROLPLANE_URI"
	const secretServiceEnvVar = "SECRET_SERVICE_URI"

	kep := new(configurableKeptnEndpointProvider)
	kep.controlPlane = utils.SanitizeURL(
		getEnvVariablewithDefault(
			controlPlaneServiceEnvVar, defaultControlPlaneEndpoint,
		),
	)
	kep.secretService = utils.SanitizeURL(
		getEnvVariablewithDefault(
			secretServiceEnvVar, defaultSecretServiceEndpoint,
		),
	)

	return kep
}

func (kep *configurableKeptnEndpointProvider) GetControlPlaneEndpoint() string {
	return kep.controlPlane
}

func (kep *configurableKeptnEndpointProvider) GetSecretsServiceEndpoint() string {
	return kep.secretService
}
