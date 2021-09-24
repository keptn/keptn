package lib_test

import (
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCmdCurlExecutor_InvalidCommand(t *testing.T) {
	executor := lib.NewCmdCurlExecutor()

	output, err := executor.Curl("invalid command")

	require.NotNil(t, err)
	require.Empty(t, output)
}

func TestNewCmdCurlExecutor_UnAllowedURL(t *testing.T) {
	executor := lib.NewCmdCurlExecutor(lib.WithUnAllowedURLs([]string{"kube-api"}))

	output, err := executor.Curl("curl http://kube-api")

	require.NotNil(t, err)
	require.Empty(t, output)
}

func TestNewCmdCurlExecutor_EmptyCommand(t *testing.T) {
	executor := lib.NewCmdCurlExecutor()

	output, err := executor.Curl("")

	require.NotNil(t, err)
	require.Empty(t, output)
}
