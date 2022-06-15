package db_mock

import "github.com/keptn/keptn/shipyard-controller/models"

type MongoDBProjectCredentialsRepoMock struct {
	GetOldCredentialsProjectsFunc func() ([]*models.ExpandedProjectOld, error)
	UpdateProjectFunc             func(project *models.ExpandedProjectOld) error
}

func (r MongoDBProjectCredentialsRepoMock) GetOldCredentialsProjects() ([]*models.ExpandedProjectOld, error) {
	if r.GetOldCredentialsProjectsFunc != nil {
		return r.GetOldCredentialsProjectsFunc()
	}
	panic("implement me")
}

func (r MongoDBProjectCredentialsRepoMock) UpdateProject(project *models.ExpandedProjectOld) error {
	if r.UpdateProjectFunc != nil {
		return r.UpdateProjectFunc(project)
	}
	panic("implement me")
}
