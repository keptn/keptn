package restapi

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_getEnvConfig(t *testing.T) {
	_ = os.Setenv("MAX_AUTH_REQUESTS_PER_SECOND", "0.5")
	_ = os.Setenv("MAX_AUTH_REQUESTS_BURST", "1")
	config, err := getEnvConfig()
	require.Nil(t, err)
	require.Equal(t, 0.5, config.AuthRequestsPerSecond)
	require.Equal(t, 1, config.AuthRequestMaxBurst)
}

func Test_getEnvConfigUseDefaults(t *testing.T) {
	config, err := getEnvConfig()
	require.Nil(t, err)
	require.Equal(t, float64(1), config.AuthRequestsPerSecond)
	require.Equal(t, 2, config.AuthRequestMaxBurst)
}
