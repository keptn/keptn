package db_mock

import "github.com/keptn/keptn/shipyard-controller/models"

type SecretCredentialsRepoMock struct {
	UpdateSecretFunc func(project *models.ExpandedProjectOld) error
}

func (r SecretCredentialsRepoMock) UpdateSecret(project *models.ExpandedProjectOld) error {
	if r.UpdateSecretFunc != nil {
		return r.UpdateSecretFunc(project)
	}
	panic("implement me")
}
