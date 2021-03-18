package backend_test

import (
	"github.com/keptn/keptn/secret-service/internal/backend"
	"github.com/keptn/keptn/secret-service/internal/backend/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Register(t *testing.T) {

	backend.Register("a", func() backend.SecretBackend {
		return &fake.SecretBackendMock{}
	})

	backend.Register("b", func() backend.SecretBackend {
		return &fake.SecretBackendMock{}
	})

	backends := backend.GetRegisteredBackends()
	assert.Contains(t, backends, backend.BackendTypeK8s)

}
