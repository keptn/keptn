package fake

import "github.com/keptn/go-utils/pkg/api/models"

type UniformInterfaceMock struct {
	RegisterIntegrationFn func(models.Integration) (string, error)
	PingFn                func(string) (*models.Integration, error)
}

func (m *UniformInterfaceMock) Ping(integrationID string) (*models.Integration, error) {
	if m.PingFn != nil {
		return m.PingFn(integrationID)
	}
	panic("Ping() not implemented")
}
func (m *UniformInterfaceMock) RegisterIntegration(integration models.Integration) (string, error) {
	if m.RegisterIntegrationFn != nil {
		return m.RegisterIntegrationFn(integration)
	}
	panic("RegisterIntegraiton not imiplemented")
}

func (m *UniformInterfaceMock) CreateSubscription(integrationID string, subscription models.EventSubscription) (string, error) {
	panic("implement me")
}

func (m *UniformInterfaceMock) UnregisterIntegration(integrationID string) error {
	panic("implement me")
}

func (m *UniformInterfaceMock) GetRegistrations() ([]*models.Integration, error) {
	panic("implement me")
}
