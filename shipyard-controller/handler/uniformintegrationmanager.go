package handler

import (
	"errors"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

var ErrSubscriptionMissing = errors.New("integration must have at least one subscription")

//go:generate moq --skip-ensure -pkg fake -out ./fake/uniformintegrationmanager.go . IUniformIntegrationManager
type IUniformIntegrationManager interface {
	Register(integration models.Integration) error
	Unregister(id string) error
	GetRegistrations(params models.GetUniformIntegrationsParams) ([]models.Integration, error)
	CreateOrUpdateSubscription(integrationID string, subscription models.Subscription) error
	DeleteSubscription(integrationID, subscriptionID string) error
	GetSubscription(integrationID, subscriptionID string) (*models.Subscription, error)
	UpdateLastSeen(integrationID string) (*models.Integration, error)
}

type UniformIntegrationManager struct {
	repo db.UniformRepo
}

func NewUniformIntegrationManager(repo db.UniformRepo) *UniformIntegrationManager {
	return &UniformIntegrationManager{
		repo: repo,
	}
}

func (uim *UniformIntegrationManager) Register(integration models.Integration) error {
	return uim.repo.CreateOrUpdateUniformIntegration(integration)
}

func (uim *UniformIntegrationManager) Unregister(id string) error {
	return uim.repo.DeleteUniformIntegration(id)
}

func (uim *UniformIntegrationManager) GetRegistrations(params models.GetUniformIntegrationsParams) ([]models.Integration, error) {
	return uim.repo.GetUniformIntegrations(params)
}
func (uim *UniformIntegrationManager) CreateOrUpdateSubscription(integrationID string, subscription models.Subscription) error {
	return uim.repo.CreateOrUpdateSubscription(integrationID, subscription)
}

func (uim *UniformIntegrationManager) DeleteSubscription(integrationID, subscriptionID string) error {
	return uim.repo.DeleteSubscription(integrationID, subscriptionID)
}

func (uim *UniformIntegrationManager) GetSubscription(integrationID, subscriptionID string) (*models.Subscription, error) {
	return uim.repo.GetSubscription(integrationID, subscriptionID)
}

func (uim *UniformIntegrationManager) UpdateLastSeen(integrationID string) (*models.Integration, error) {
	return uim.repo.UpdateLastSeen(integrationID)
}
