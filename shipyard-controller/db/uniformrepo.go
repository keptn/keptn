package db

import "github.com/keptn/keptn/shipyard-controller/models"

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/uniformrepo_mock.go . UniformRepo
type UniformRepo interface {
	GetUniformIntegrations(filter models.GetUniformIntegrationsParams) ([]models.Integration, error)
	DeleteUniformIntegration(id string) error
	CreateOrUpdateUniformIntegration(integration models.Integration) error
	CreateOrUpdateSubscription(integrationID string, subscription models.Subscription) error
	DeleteSubscription(integrationID, subscriptionID string) error
	GetSubscription(integrationID, subscriptionID string) (*models.Subscription, error)
	GetSubscriptions(integrationID string) ([]models.Subscription, error)
	UpdateLastSeen(integrationID string) (*models.Integration, error)
}
