package fake

import "github.com/keptn/keptn/shipyard-controller/models"

type IRepositoryProvisionerMock struct {
	ProvideRepositoryFn func(string) (*models.ProvisioningData, error)
	DeleteRepositoryFn  func(string, string) error
}

func (r *IRepositoryProvisionerMock) ProvideRepository(projectName string) (*models.ProvisioningData, error) {
	if r.ProvideRepositoryFn != nil {
		return r.ProvideRepositoryFn(projectName)
	}
	panic("implement me")
}

func (r *IRepositoryProvisionerMock) DeleteRepository(projectName string, namespace string) error {
	if r.DeleteRepositoryFn != nil {
		return r.DeleteRepositoryFn(projectName, namespace)
	}
	panic("implement me")
}
