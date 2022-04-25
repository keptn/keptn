package restapi

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getEnvConfig(t *testing.T) {
	defer os.Unsetenv("MAX_AUTH_ENABLED")
	defer os.Unsetenv("MAX_AUTH_REQUESTS_PER_SECOND")
	defer os.Unsetenv("MAX_AUTH_REQUESTS_BURST")
	_ = os.Setenv("MAX_AUTH_ENABLED", "false")
	_ = os.Setenv("MAX_AUTH_REQUESTS_PER_SECOND", "0.5")
	_ = os.Setenv("MAX_AUTH_REQUESTS_BURST", "1")

	config, err := getEnvConfig()
	require.Nil(t, err)
	require.Equal(t, false, config.MaxAuthEnabled)
	require.Equal(t, 0.5, config.MaxAuthRequestsPerSecond)
	require.Equal(t, 1, config.MaxAuthRequestBurst)
}

func Test_getEnvConfigUseDefaults(t *testing.T) {
	config, err := getEnvConfig()
	require.Nil(t, err)
	require.Equal(t, true, config.MaxAuthEnabled)
	require.Equal(t, float64(1), config.MaxAuthRequestsPerSecond)
	require.Equal(t, 2, config.MaxAuthRequestBurst)
}
