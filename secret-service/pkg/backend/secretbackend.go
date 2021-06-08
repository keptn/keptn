package backend

import (
	"github.com/keptn/keptn/secret-service/pkg/model"
)

const DefaultNamespace = "keptn"

//go:generate moq -pkg fake -out ./fake/secretbackend_mock.go . SecretBackend
type SecretBackend interface {
	CreateSecret(model.Secret) error
	UpdateSecret(model.Secret) error
	DeleteSecret(model.Secret) error
	GetSecrets() ([]model.SecretMetadata, error)
}
