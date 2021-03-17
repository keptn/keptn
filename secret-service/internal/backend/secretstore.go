package backend

import "github.com/keptn/keptn/secret-service/internal/model"

//go:generate moq -pkg fake -out ./fake/secretstore_mock.go . SecretStore
type SecretStore interface {
	CreateSecret(model.Secret) error
	UpdateSecret(model.Secret) error
	DeleteSecret(model.Secret) error
}
