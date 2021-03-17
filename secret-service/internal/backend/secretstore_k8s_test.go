package backend

import (
	"errors"
	"github.com/keptn/keptn/secret-service/internal/common"
	"github.com/keptn/keptn/secret-service/internal/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
	"testing"
)

func FakeNamespaceProvider() common.StringSupplier {
	return func() string {
		return "keptn_namespace"
	}
}

func TestCreateK8sSecretStore(t *testing.T) {
	store := NewK8sSecretStore(fake.NewSimpleClientset())
	assert.NotNil(t, store)
}

func TestCreateSecret(t *testing.T) {

	kubernetes := fake.NewSimpleClientset()
	secretStore := K8sSecretStore{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
	}

	secret := model.Secret{
		Name:  "my-secret",
		Scope: "my-scope",
		Data:  map[string]string{"password": "keptn"},
	}

	err := secretStore.CreateSecret(secret)
	assert.Nil(t, err)
	kubernetesSecret, err := kubernetes.CoreV1().Secrets(FakeNamespaceProvider()()).Get("my-secret", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, kubernetesSecret)
	assert.Equal(t, "my-secret", kubernetesSecret.Name)
	assert.Equal(t, map[string]string(secret.Data), kubernetesSecret.StringData)
	assert.Equal(t, FakeNamespaceProvider()(), kubernetesSecret.Namespace)
}

func TestCreateSecret_KubernetesSecretCreationFails(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	secretStore := K8sSecretStore{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
	}

	secret := model.Secret{
		Name:  "my-secret",
		Scope: "my-scope",
		Data:  map[string]string{"password": "keptn"},
	}

	kubernetes.Fake.PrependReactor("create", "secrets", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{}, errors.New("Error creating kubernetes secret")
	})

	err := secretStore.CreateSecret(secret)
	assert.NotNil(t, err)

}
