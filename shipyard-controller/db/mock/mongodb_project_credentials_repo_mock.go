package db_mock

import "github.com/keptn/keptn/shipyard-controller/models"

type MongoDBSecretCredentialsRepoMock struct {
	UpdateSecretFunc func(project *models.ExpandedProjectOld) error
}

func (r MongoDBSecretCredentialsRepoMock) UpdateSecret(project *models.ExpandedProjectOld) error {
	if r.UpdateSecretFunc != nil {
		return r.UpdateSecretFunc(project)
	}
	panic("implement me")
}
