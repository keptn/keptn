package lib_test

import (
	"testing"

	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/stretchr/testify/require"
)

func TestDeniedURLS(t *testing.T) {
	kubeEnvs := map[string]string{"KUBERNETES_SERVICE_HOST": "1.2.3.4", "KUBERNETES_SERVICE_PORT": "9876"}
	urls := lib.GetDeniedURLs(kubeEnvs)

	expected := []string{"kubernetes", "kubernetes.default", "kubernetes.default.svc", "kubernetes.default.svc.cluster.local", "localhost", "127.0.0.1", "::1", "1.2.3.4", "kubernetes:9876", "kubernetes.default:9876", "kubernetes.default.svc:9876", "kubernetes.default.svc.cluster.local:9876", "1.2.3.4:9876"}

	require.Equal(t, 13, len(urls))
	require.Equal(t, expected, urls)
}
