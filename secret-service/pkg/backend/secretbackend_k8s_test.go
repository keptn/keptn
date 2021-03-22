package backend

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/secret-service/pkg/common"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/keptn/keptn/secret-service/pkg/repository/fake"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
	"testing"
)

func FakeNamespaceProvider() common.StringSupplier {
	return func() string {
		return "keptn_namespace"
	}
}

func TestCreateK8sSecretBackend(t *testing.T) {
	backend := NewK8sSecretBackend(k8sfake.NewSimpleClientset(), &fake.ScopesRepositoryMock{})
	assert.NotNil(t, backend)
}

func TestCreateSecret(t *testing.T) {

	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	secret := createTestSecret("my-secret", "my-scope")
	err := backend.CreateSecret(secret)
	assert.Nil(t, err)

	k8sSecret, err := kubernetes.CoreV1().Secrets(FakeNamespaceProvider()()).Get("my-secret", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, k8sSecret)
	assert.Equal(t, "my-secret", k8sSecret.Name)
	assert.Equal(t, map[string]string(secret.Data), k8sSecret.StringData)
	assert.Equal(t, FakeNamespaceProvider()(), k8sSecret.Namespace)

	k8sRole1, err := kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get("my-scope-read-secrets", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, k8sRole1)
	assert.Equal(t, k8sRole1.Rules[0].Resources[0], "secrets")
	assert.Equal(t, k8sRole1.Rules[0].ResourceNames[0], "my-secret")
	assert.Equal(t, k8sRole1.Rules[0].Verbs, []string{"read"})
	assert.Equal(t, k8sRole1.Rules[0].APIGroups, []string{""}) // at least on api group must be present

	k8sRole2, err := kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get("my-scope-manage-secrets", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, k8sRole2)
	assert.Equal(t, k8sRole2.Rules[0].Resources[0], "secrets")
	assert.Equal(t, k8sRole2.Rules[0].ResourceNames[0], "my-secret")
	assert.Equal(t, k8sRole2.Rules[0].Verbs, []string{"create", "read", "update"})
	assert.Equal(t, k8sRole1.Rules[0].APIGroups, []string{""}) // at least on api group must be present
}

func TestCreateSecret_FetchingScopesFails(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return model.Scopes{}, fmt.Errorf("error fetching scopes") }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	secret := createTestSecret("my-secret", "my-scope")
	err := backend.CreateSecret(secret)
	assert.NotNil(t, err)
}

func TestCreateSecret_K8sSecretCreationFails(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	secret := createTestSecret("my-secret", "my-scope")

	kubernetes.Fake.PrependReactor("create", "secrets", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("Error creating kubernetes secret")
	})

	err := backend.CreateSecret(secret)
	assert.NotNil(t, err)
}

func TestCreateSecret_K8sRolesCreationFails(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	secret := createTestSecret("my-secret", "my-scope")

	kubernetes.Fake.PrependReactor("create", "roles", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("Error creating kubernetes roles")
	})

	err := backend.CreateSecret(secret)
	assert.NotNil(t, err)
}

func TestCreateSecret_NoMatchingScopeConfigured(t *testing.T) {

	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	secret := createTestSecret("my-secret", "my-other-scope")
	err := backend.CreateSecret(secret)
	assert.NotNil(t, err)
}

func createTestSecret(name, scope string) model.Secret {
	secret := model.Secret{
		Name:  name,
		Scope: scope,
		Data:  map[string]string{"password": "keptn"},
	}
	return secret
}

func createTestScopes() model.Scopes {
	scopes := model.Scopes{
		Scopes: map[string]model.Scope{
			"my-scope": {
				Capabilities: map[string]model.Capability{
					"my-scope-read-secrets": {
						Permissions: []string{"read"},
					},
					"my-scope-manage-secrets": {
						Permissions: []string{"create", "read", "update"},
					},
				},
			},
		},
	}
	return scopes
}
