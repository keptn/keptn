package lib_test

import (
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"os"
	"testing"
)

func TestK8sSecretReater_ReadSecret(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := lib.NewK8sSecretReader(fake.NewSimpleClientset(
		getK8sSecret(map[string]string{
			"app.kubernetes.io/managed-by": "keptn-secret-service",
		}),
	))

	secret, err := secretReader.ReadSecret("my-secret", "foo")

	require.Nil(t, err)
	require.Equal(t, "bar", secret)

	secret, err = secretReader.ReadSecret("my-missing-secret", "foo")

	require.NotNil(t, err)
	require.Empty(t, secret)
}

func TestK8sSecretReater_ReadSecretWithInvalidScope(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := lib.NewK8sSecretReader(fake.NewSimpleClientset(
		getK8sSecret(map[string]string{}),
	))

	secret, err := secretReader.ReadSecret("my-secret", "foo")

	require.NotNil(t, err)
	require.Equal(t, "", secret)
}

func getK8sSecret(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "keptn",
			Labels:    labels,
		},
		Data: map[string][]byte{
			"foo": []byte("bar"),
		},
		Type: corev1.SecretTypeOpaque,
	}
}
