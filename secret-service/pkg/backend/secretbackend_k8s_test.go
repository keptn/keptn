package backend

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/keptn/secret-service/pkg/common"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/keptn/keptn/secret-service/pkg/repository/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
CREATE SECRET TESTS
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

	k8sSecret, err := kubernetes.CoreV1().Secrets(FakeNamespaceProvider()()).Get(context.TODO(), "my-secret", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, k8sSecret)
	assert.Equal(t, "my-secret", k8sSecret.Name)
	assert.Equal(t, map[string]string(secret.Data), k8sSecret.StringData)
	assert.Equal(t, FakeNamespaceProvider()(), k8sSecret.Namespace)
	assert.Equal(t, SecretServiceName, k8sSecret.Labels["app.kubernetes.io/managed-by"])

	k8sRole1, err := kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get(context.TODO(), "my-scope-read-secrets", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, k8sRole1)
	assert.Equal(t, k8sRole1.Rules[0].Resources[0], "secrets")
	assert.Equal(t, k8sRole1.Rules[0].ResourceNames[0], "my-secret")
	assert.Equal(t, k8sRole1.Rules[0].Verbs, []string{"read"})
	assert.Equal(t, k8sRole1.Rules[0].APIGroups, []string{""}) // at least on api group must be present

	k8sRole2, err := kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get(context.TODO(), "my-scope-manage-secrets", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, k8sRole2)
	assert.Equal(t, k8sRole2.Rules[0].Resources[0], "secrets")
	assert.Equal(t, k8sRole2.Rules[0].ResourceNames[0], "my-secret")
	assert.Equal(t, k8sRole2.Rules[0].Verbs, []string{"create", "read", "update"})
	assert.Equal(t, k8sRole1.Rules[0].APIGroups, []string{""}) // at least on api group must be present

	k8sRoleBinding, _ := kubernetes.RbacV1().RoleBindings(FakeNamespaceProvider()()).Get(context.TODO(), "my-scope-rolebinding", metav1.GetOptions{})
	assert.Equal(t, "my-scope", k8sRoleBinding.Subjects[0].Name)

	nextSecret := createTestSecret("my-secret-2", "my-scope")
	err = backend.CreateSecret(nextSecret)
	assert.Nil(t, err)

	k8sRole1, err = kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get(context.TODO(), "my-scope-read-secrets", metav1.GetOptions{})
	assert.Equal(t, []string{"my-secret", "my-secret-2"}, k8sRole1.Rules[0].ResourceNames)
	k8sRole2, err = kubernetes.RbacV1().Roles(FakeNamespaceProvider()()).Get(context.TODO(), "my-scope-manage-secrets", metav1.GetOptions{})
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
GET SECRET TESTS
*/
func TestGetSecret(t *testing.T) {
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

	secrets, err := backend.GetSecrets()
	require.Nil(t, err)

	require.Equal(t, []model.GetSecretResponseItem{
		{
			SecretMetadata: model.SecretMetadata{
				Name:  "my-secret",
				Scope: "my-scope",
			},
			Keys: []string{"password"},
		},
	}, secrets)
}

func TestGetSecret_Fails(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	kubernetes.Fake.PrependReactor("list", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("oops")
	})

	secrets, err := backend.GetSecrets()

	require.NotNil(t, err)
	require.Nil(t, secrets)
}

func TestDeleteK8sSecret(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "keptn_namespace",
			Labels:    map[string]string{"app.kubernetes.io/scope": "my-scope"},
		},
	},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-other-secret",
				Namespace: "keptn_namespace",
				Labels:    map[string]string{"app.kubernetes.io/scope": "my-scope"},
			},
		},
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-scope-read-secrets",
				Namespace: "keptn_namespace",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Resources:     []string{"secrets"},
					ResourceNames: []string{"my-other-secret", "my-secret"},
				},
			},
		},
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-scope-manage-secrets",
				Namespace: "keptn_namespace",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Resources:     []string{"secrets"},
					ResourceNames: []string{"my-other-secret", "my-secret"},
				},
			},
		})
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	err := backend.DeleteSecret(model.Secret{
		SecretMetadata: model.SecretMetadata{
			Name:  "my-secret",
			Scope: "my-scope",
		},
	})

	assert.Nil(t, err)
	assert.True(t, kubernetes.Fake.Actions()[0].Matches("delete", "secrets"))
	assert.True(t, kubernetes.Fake.Actions()[1].Matches("list", "secrets"))
	assert.True(t, kubernetes.Fake.Actions()[2].Matches("get", "roles"))
	assert.True(t, kubernetes.Fake.Actions()[3].Matches("update", "roles"))
	assert.True(t, kubernetes.Fake.Actions()[4].Matches("get", "roles"))
	assert.True(t, kubernetes.Fake.Actions()[5].Matches("update", "roles"))
}

func TestDeleteLastK8sSecret(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "keptn_namespace",
			Labels:    map[string]string{"app.kubernetes.io/scope": "my-scope"},
		},
	},
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-scope-read-secrets",
				Namespace: "keptn_namespace",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Resources:     []string{"secrets"},
					ResourceNames: []string{"my-secret"},
				},
			},
		},
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-scope-manage-secrets",
				Namespace: "keptn_namespace",
			},
			Rules: []rbacv1.PolicyRule{
				{
					Resources:     []string{"secrets"},
					ResourceNames: []string{"my-secret"},
				},
			},
		})
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	err := backend.DeleteSecret(model.Secret{
		SecretMetadata: model.SecretMetadata{
			Name:  "my-secret",
			Scope: "my-scope",
		},
	})

	appliedActions := kubernetes.Fake.Actions()
	_ = appliedActions
	assert.Nil(t, err)
	assert.True(t, kubernetes.Fake.Actions()[0].Matches("delete", "secrets"))
	assert.True(t, kubernetes.Fake.Actions()[1].Matches("list", "secrets"))
	assert.True(t, kubernetes.Fake.Actions()[2].Matches("delete-collection", "roles"))
	assert.True(t, kubernetes.Fake.Actions()[3].Matches("delete-collection", "rolebindings"))

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
		SecretMetadata: model.SecretMetadata{
			Name:  "my-secret",
			Scope: "my-scope",
		},
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
		SecretMetadata: model.SecretMetadata{
			Name:  name,
			Scope: scope,
		},
		Data: map[string]string{"password": "keptn"},
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
