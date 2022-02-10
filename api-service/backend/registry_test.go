package backend_test

import (
	"testing"

	"github.com/keptn/keptn/api-service/backend"
	"github.com/keptn/keptn/api-service/backend/fake"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {

	backend.Register("a", func() backend.SecretBackend {
		return &fake.SecretBackendMock{}
	})

	backend.Register("b", func() backend.SecretBackend {
		return &fake.SecretBackendMock{}
	})

	backends := backend.GetRegisteredBackends()
	assert.Contains(t, backends, backend.SecretBackendTypeK8s)

}
