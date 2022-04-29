package lib_test

import (
	"fmt"
	"testing"

	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func TestDeniedURLS(t *testing.T) {
	kubeEnvs := map[string]string{"KUBERNETES_SERVICE_HOST": "1.2.3.4", "KUBERNETES_SERVICE_PORT": "9876"}
	urls := lib.GetDeniedURLs(kubeEnvs)

	expected := []string{"1.2.3.4", "kubernetes:9876", "kubernetes.default:9876", "kubernetes.default.svc:9876", "kubernetes.default.svc.cluster.local:9876", "1.2.3.4:9876"}

	require.Equal(t, 6, len(urls))
	require.Equal(t, expected, urls)
}

func TestGetDenyList(t *testing.T) {
	denyListString := "some\nurl\nip"
	tests := []struct {
		name             string
		denyListProvider lib.DenyListProvider
		want             []string
	}{
		{
			name: "cannot get kubeAPI",
			denyListProvider: lib.DenyListProvider{
				GetKubeAPI: func(useInClusterConfig bool) (v1.CoreV1Interface, error) {
					return nil, fmt.Errorf("cannot get kubeAPIClient")
				},
				GetDeniedURLs: func(env map[string]string) []string {
					return []string{"1.2.3.4", "kubernetes:9876"}
				},
			},
			want: []string{"1.2.3.4", "kubernetes:9876"},
		},
		{
			name: "cannot get configmap",
			denyListProvider: lib.DenyListProvider{
				GetKubeAPI: func(useInClusterConfig bool) (v1.CoreV1Interface, error) {
					return &v1.CoreV1Client{}, nil
				},
				GetWebhookDenyListConfigMap: func(kubeAPI v1.CoreV1Interface) (*corev1.ConfigMap, error) {
					return nil, fmt.Errorf("cannot get configmap")
				},
				GetDeniedURLs: func(env map[string]string) []string {
					return []string{"1.2.3.4", "kubernetes:9876"}
				},
			},
			want: []string{"1.2.3.4", "kubernetes:9876"},
		},
		{
			name: "valid empty configmap",
			denyListProvider: lib.DenyListProvider{
				GetKubeAPI: func(useInClusterConfig bool) (v1.CoreV1Interface, error) {
					return &v1.CoreV1Client{}, nil
				},
				GetWebhookDenyListConfigMap: func(kubeAPI v1.CoreV1Interface) (*corev1.ConfigMap, error) {
					return &corev1.ConfigMap{}, nil
				},
				GetDeniedURLs: func(env map[string]string) []string {
					return []string{"1.2.3.4", "kubernetes:9876"}
				},
			},
			want: []string{"1.2.3.4", "kubernetes:9876"},
		},
		{
			name: "valid",
			denyListProvider: lib.DenyListProvider{
				GetKubeAPI: func(useInClusterConfig bool) (v1.CoreV1Interface, error) {
					return &v1.CoreV1Client{}, nil
				},
				GetWebhookDenyListConfigMap: func(kubeAPI v1.CoreV1Interface) (*corev1.ConfigMap, error) {
					return &corev1.ConfigMap{Data: map[string]string{"denyList": denyListString}}, nil
				},
				GetDeniedURLs: func(env map[string]string) []string {
					return []string{"1.2.3.4", "kubernetes:9876"}
				},
			},
			want: []string{"1.2.3.4", "kubernetes:9876", "some", "url", "ip"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.denyListProvider.GetDenyList()
			require.Equal(t, tt.want, got)
		})
	}
}
