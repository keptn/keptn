package backend

import (
	"github.com/keptn/keptn/secret-service/pkg/model"
)

const DefaultNamespace = "keptn"

type SecretManager interface {
	CreateSecret(model.Secret) error
	UpdateSecret(model.Secret) error
	DeleteSecret(model.Secret) error
	GetSecrets() ([]model.GetSecretResponseItem, error)
}

type ScopeManager interface {
	GetScopes() ([]string, error)
}

//go:generate moq -pkg fake -out ./fake/secretbackend_mock.go . SecretBackend
type SecretBackend interface {
	SecretManager
	ScopeManager
}
