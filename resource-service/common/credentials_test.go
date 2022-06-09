package common

import (
	"fmt"
	"os"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestK8sCredentialReader_ReadSecret(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := NewK8sCredentialReader(fake.NewSimpleClientset(
		getK8sSecret(),
	))

	secret, err := secretReader.GetCredentials("my-project")

	require.Nil(t, err)
	require.Equal(t, &common_models.GitCredentials{
		User: "user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "token",
		},
		RemoteURL: "https://my-repo",
	}, secret)
}

func TestK8sCredentialReader_ReadSecretNotFound(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := NewK8sCredentialReader(fake.NewSimpleClientset())

	secret, err := secretReader.GetCredentials("my-other-project")

	require.ErrorIs(t, err, errors.ErrCredentialsNotFound)
	require.Nil(t, secret)
}

func TestK8sCredentialReader_ReadSecretWrongFormat(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := NewK8sCredentialReader(fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "git-credentials-my-project",
				Namespace: "keptn",
			},
			Data: map[string][]byte{
				"git-credentials": []byte(`invalid`)},
			Type: corev1.SecretTypeOpaque,
		},
	))

	secret, err := secretReader.GetCredentials("my-project")

	require.ErrorIs(t, err, errors.ErrMalformedCredentials)
	require.Nil(t, secret)
}

func TestK8sCredentialReader_ReadSecretNoToken(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := NewK8sCredentialReader(fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "git-credentials-my-project",
				Namespace: "keptn",
			},
			Data: map[string][]byte{
				"git-credentials": []byte(`{"user":"user","remoteURL":"https://some.url","https":{"token":""}}`)},
			Type: corev1.SecretTypeOpaque,
		},
	))

	secret, err := secretReader.GetCredentials("my-project")

	require.ErrorIs(t, err, errors.ErrCredentialsTokenMustNotBeEmpty)
	require.Nil(t, secret)
}

func TestK8sCredentialReader_ReadSecretNoPrivateKey(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")
	secretReader := NewK8sCredentialReader(fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "git-credentials-my-project",
				Namespace: "keptn",
			},
			Data: map[string][]byte{
				"git-credentials": []byte(`{"user":"user","remoteURL":"ssh://some.url", "ssh":{"privateKey":""}}`)},
			Type: corev1.SecretTypeOpaque,
		},
	))

	secret, err := secretReader.GetCredentials("my-project")

	require.ErrorIs(t, err, errors.ErrCredentialsPrivateKeyMustNotBeEmpty)
	require.Nil(t, secret)
}

func TestK8sCredentialReader_ReadSecretError(t *testing.T) {
	_ = os.Setenv("POD_NAMESPACE", "keptn")

	fakeClient := fake.NewSimpleClientset()

	fakeClient.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("oops")
	})
	secretReader := NewK8sCredentialReader(fakeClient)

	secret, err := secretReader.GetCredentials("my-project")

	require.NotNil(t, err)
	require.Nil(t, secret)
}

func getK8sSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "git-credentials-my-project",
			Namespace: "keptn",
		},
		Data: map[string][]byte{
			"git-credentials": []byte(`{"user":"user","remoteURL":"https://my-repo","https":{"token":"token"}}`)},
		Type: corev1.SecretTypeOpaque,
	}
}
