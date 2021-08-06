package db

import "github.com/keptn/keptn/shipyard-controller/models"

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/uniformrepo_mock.go . UniformRepo
type UniformRepo interface {
	GetUniformIntegrations(filter models.GetUniformIntegrationsParams) ([]models.Integration, error)
	DeleteUniformIntegration(id string) error
	CreateOrUpdateUniformIntegration(integration models.Integration) error
}
