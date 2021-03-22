package backend

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/secret-service/pkg/common"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/keptn/keptn/secret-service/pkg/repository/fake"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
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

/**
CREATE SECREAT TESTS
*/
func TestCreateSecrets(t *testing.T) {

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

	nextSecret := createTestSecret("my-secret-2", "my-scope")
	err = backend.CreateSecret(nextSecret)
	assert.Nil(t, err)

	k8sRole1, err = kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get("my-scope-read-secrets", metav1.GetOptions{})
	assert.Equal(t, []string{"my-secret", "my-secret-2"}, k8sRole1.Rules[0].ResourceNames)
	k8sRole2, err = kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get("my-scope-manage-secrets", metav1.GetOptions{})
	assert.Equal(t, []string{"my-secret", "my-secret-2"}, k8sRole2.Rules[0].ResourceNames)

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

	kubernetes.Fake.PrependReactor("create", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
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

	kubernetes.Fake.PrependReactor("create", "roles", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
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

/**
DELETE SECRET TESTS
*/
func TestDeleteK8sSecret(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	kubernetes.Fake.PrependReactor("get", "roles", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Role{
			Rules: []v1.PolicyRule{
				{
					Resources:     []string{"secrets"},
					ResourceNames: []string{"my-other-secret", "my-secret"},
				},
			},
		}, nil
	})

	kubernetes.Fake.PrependReactor("patch", "roles", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})

	kubernetes.Fake.PrependReactor("delete", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})

	err := backend.DeleteSecret(model.Secret{
		Name:  "my-secret",
		Scope: "my-scope",
	})

	actions := kubernetes.Fake.Actions()
	_ = actions

	assert.Nil(t, err)
	assert.True(t, kubernetes.Fake.Actions()[0].Matches("delete", "secrets"))
	assert.Equal(t, schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}, kubernetes.Fake.Actions()[0].GetResource())
	assert.True(t, kubernetes.Fake.Actions()[1].Matches("get", "roles"))
	assert.Equal(t, schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}, kubernetes.Fake.Actions()[1].GetResource())
	assert.Equal(t, `[{"op":"replace","path":"/rules/0/resourceNames","value":["my-other-secret"]}]`, string(kubernetes.Fake.Actions()[2].(k8stesting.PatchAction).GetPatch()))
	assert.True(t, kubernetes.Fake.Actions()[3].Matches("get", "roles"))
	assert.Equal(t, schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}, kubernetes.Fake.Actions()[3].GetResource())
	assert.Equal(t, `[{"op":"replace","path":"/rules/0/resourceNames","value":["my-other-secret"]}]`, string(kubernetes.Fake.Actions()[4].(k8stesting.PatchAction).GetPatch()))

}

func TestDeleteK8sSecret_SecretNotFound(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	err := backend.DeleteSecret(model.Secret{
		Name:  "my-secret",
		Scope: "my-scope",
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrSecretNotFound, err)
}

/**
UPDATE SECRET TESTS
*/
func TestUpdateSecret(t *testing.T) {

	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	kubernetes.Fake.PrependReactor("update", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, nil
	})

	secret := createTestSecret("my-secret", "my-scope")
	err := backend.UpdateSecret(secret)
	assert.Nil(t, err)

}

func TestUpdateSecret_SecretNotFound(t *testing.T) {

	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	secret := createTestSecret("my-secret", "my-scope")
	err := backend.UpdateSecret(secret)
	assert.NotNil(t, err)
	assert.Equal(t, ErrSecretNotFound, err)
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
