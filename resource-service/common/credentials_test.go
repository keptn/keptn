package common

import (
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"os"
	"testing"
)

func TestK8sCredentialReader_ReadSecret(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := NewK8sCredentialReader(fake.NewSimpleClientset(
		getK8sSecret(),
	))

	secret, err := secretReader.GetCredentials("my-project")

	require.Nil(t, err)
	require.Equal(t, &GitCredentials{
		User:      "user",
		Token:     "token",
		RemoteURI: "uri",
	}, secret)

	secret, err = secretReader.GetCredentials("my-other-project")

	require.ErrorIs(t, err, errors.ErrCredentialsNotFound)
	require.Nil(t, secret)
}

func getK8sSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "git-credentials-my-project",
			Namespace: "keptn",
		},
		Data: map[string][]byte{
			"git-credentials": []byte(`{"user":"user","token":"token","remoteURI":"uri"}`)},
		Type: corev1.SecretTypeOpaque,
	}
}
