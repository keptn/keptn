package backend

import (
	"errors"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/keptn/keptn/secret-service/pkg/repository/fake"
	"github.com/stretchr/testify/require"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

/**
GET SCOPE TESTS
*/
func TestGetScopes(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return createTestScopes(), nil }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	scopes, err := backend.GetScopes()
	require.Nil(t, err)

	require.Equal(t, []string{"my-scope"}, scopes)
}

func TestGetScopes_Fails(t *testing.T) {
	kubernetes := k8sfake.NewSimpleClientset()
	scopesRepository := &fake.ScopesRepositoryMock{}
	scopesRepository.ReadFunc = func() (model.Scopes, error) { return model.Scopes{}, errors.New("oops") }

	backend := K8sSecretBackend{
		KubeAPI:                kubernetes,
		KeptnNamespaceProvider: FakeNamespaceProvider(),
		ScopesRepository:       scopesRepository,
	}

	scopes, err := backend.GetScopes()

	require.NotNil(t, err)
	require.Nil(t, scopes)
}
