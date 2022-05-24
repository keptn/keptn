package lib

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestDeniedURLS(t *testing.T) {
	kubeEnvs := map[string]string{"KUBERNETES_SERVICE_HOST": "1.2.3.4", "KUBERNETES_SERVICE_PORT": "9876"}
	urls := GetDeniedURLs(kubeEnvs)

	expected := []string{"1.2.3.4", "kubernetes:9876", "kubernetes.default:9876", "kubernetes.default.svc:9876", "kubernetes.default.svc.cluster.local:9876", "1.2.3.4:9876"}

	require.Equal(t, 6, len(urls))
	require.Equal(t, expected, urls)
}

func TestDeniedURLSNoEnv(t *testing.T) {
	kubeEnvs := map[string]string{}
	urls := GetDeniedURLs(kubeEnvs)
	t.Logf("Current denylist: %s", urls)
	expected := 0
	require.Equal(t, expected, len(urls))
}

func TestCannotGetConfigMap(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("get", "configmap", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("cannot get configmap")
	})
	denyListProvider := denyListProvider{
		getDeniedURLs: func(env map[string]string) []string {
			return []string{"1.2.3.4", "kubernetes:9876"}
		},
		kubeClient: client,
	}

	got := denyListProvider.Get()
	require.Equal(t, []string{"1.2.3.4", "kubernetes:9876"}, got)

}

func TestGetDenyList(t *testing.T) {
	denyListString := "some\nurl\nip"
	tests := []struct {
		name             string
		denyListProvider DenyListProvider
		want             []string
	}{
		{
			name: "valid empty configmap",
			denyListProvider: denyListProvider{
				getDeniedURLs: func(env map[string]string) []string {
					return []string{"1.2.3.4", "kubernetes:9876"}
				},
				kubeClient: fake.NewSimpleClientset(
					&corev1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name: "keptn-webhook-config",
						},
						Data: map[string]string{
							"denyList": ""},
					}),
			},
			want: []string{"1.2.3.4", "kubernetes:9876"},
		},
		{
			name: "valid",
			denyListProvider: denyListProvider{
				getDeniedURLs: func(env map[string]string) []string {
					return []string{"1.2.3.4", "kubernetes:9876"}
				},
				kubeClient: fake.NewSimpleClientset(
					&corev1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name: "keptn-webhook-config",
						},
						Data: map[string]string{
							"denyList": denyListString},
					}),
			},
			want: []string{"1.2.3.4", "kubernetes:9876", "some", "url", "ip"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.denyListProvider.Get()
			require.Equal(t, tt.want, got)
		})
	}
}
