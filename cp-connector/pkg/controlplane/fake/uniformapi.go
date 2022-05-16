package fake

import "github.com/keptn/go-utils/pkg/api/models"

type UniformAPIMock struct {
	RegisterIntegrationFn func(models.Integration) (string, error)
	PingFn                func(string) (*models.Integration, error)
}

func (m *UniformAPIMock) Ping(integrationID string) (*models.Integration, error) {
	if m.PingFn != nil {
		return m.PingFn(integrationID)
	}
	panic("Ping() not implemented")
}
func (m *UniformAPIMock) RegisterIntegration(integration models.Integration) (string, error) {
	if m.RegisterIntegrationFn != nil {
		return m.RegisterIntegrationFn(integration)
	}
	panic("RegisterIntegraiton() not imiplemented")
}

func (m *UniformAPIMock) CreateSubscription(integrationID string, subscription models.EventSubscription) (string, error) {
	panic("implement me")
}

func (m *UniformAPIMock) UnregisterIntegration(integrationID string) error {
	panic("implement me")
}

func (m *UniformAPIMock) GetRegistrations() ([]*models.Integration, error) {
	panic("implement me")
}
