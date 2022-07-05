package execute

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvVariablewithDefault(t *testing.T) {
	const envVar1Name = "ENV_VAR_1"
	const envVar1Value = "VALUE_1"
	const emptyVarName = "EMPTY_VAR"
	const defaultValue1 = "MY_DEFAULT_1"
	os.Setenv(envVar1Name, envVar1Value)
	os.Setenv(emptyVarName, "")
	tests := []struct {
		name          string
		envVarName    string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "Simple env var lookup",
			envVarName:    envVar1Name,
			defaultValue:  defaultValue1,
			expectedValue: envVar1Value,
		},
		{
			name:          "Env var not defined, use default",
			envVarName:    "NotExistingVar",
			defaultValue:  defaultValue1,
			expectedValue: defaultValue1,
		},
		{
			name:          "Empty env var defined, return the empty value",
			envVarName:    emptyVarName,
			defaultValue:  defaultValue1,
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {

				actual := getEnvVariablewithDefault(tt.envVarName, tt.defaultValue)
				assert.Equal(t, tt.expectedValue, actual)
			},
		)
	}
}

func TestKeptnEndpointProviderFromEnv(t *testing.T) {
	tests := []struct {
		name                         string
		env                          map[string]string
		expectedControlPlaneEndpoint string
	}{
		{
			name:                         "No Env setup - get default value",
			env:                          map[string]string{},
			expectedControlPlaneEndpoint: defaultControlPlaneEndpoint,
		},
		{
			name: "Set custom control plane without scheme",
			env: map[string]string{
				"CONTROLPLANE_URI": "somehost:1234",
			},
			expectedControlPlaneEndpoint: "http://somehost:1234",
		},
		{
			name: "Set custom control plane with http scheme",
			env: map[string]string{
				"CONTROLPLANE_URI": "http://somehost:1234",
			},
			expectedControlPlaneEndpoint: "http://somehost:1234",
		},
		{
			name: "Set custom control plane with https scheme",
			env: map[string]string{
				"CONTROLPLANE_URI": "https://somehost:1234",
			},
			expectedControlPlaneEndpoint: "https://somehost:1234",
		},

	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				os.Clearenv()
				for k, v := range tt.env {
					os.Setenv(k, v)
				}
				sut := KeptnEndpointProviderFromEnv()
				assert.Equal(t, tt.expectedControlPlaneEndpoint, sut.GetControlPlaneEndpoint())
			},
		)
	}
}
