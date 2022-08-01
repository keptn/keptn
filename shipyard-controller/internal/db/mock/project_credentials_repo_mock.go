package db_mock

import "github.com/keptn/keptn/shipyard-controller/models"

type ProjectCredentialsRepoMock struct {
	GetOldCredentialsProjectsFunc   func() ([]*models.ExpandedProjectOld, error)
	UpdateProjectFunc               func(project *models.ExpandedProjectOld) error
	CreateOldCredentialsProjectFunc func(project *models.ExpandedProjectOld) error
}

func (r ProjectCredentialsRepoMock) GetOldCredentialsProjects() ([]*models.ExpandedProjectOld, error) {
	if r.GetOldCredentialsProjectsFunc != nil {
		return r.GetOldCredentialsProjectsFunc()
	}
	panic("implement me")
}

func (r ProjectCredentialsRepoMock) UpdateProject(project *models.ExpandedProjectOld) error {
	if r.UpdateProjectFunc != nil {
		return r.UpdateProjectFunc(project)
	}
	panic("implement me")
}

func (r ProjectCredentialsRepoMock) CreateOldCredentialsProject(project *models.ExpandedProjectOld) error {
	if r.CreateOldCredentialsProjectFunc != nil {
		return r.CreateOldCredentialsProjectFunc(project)
	}
	panic("implement me")
}
