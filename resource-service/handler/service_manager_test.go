package handler

import (
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

type serviceManagerTestFields struct {
	git              *common_mock.IGitMock
	credentialReader *common_mock.CredentialReaderMock
	fileWriter       *common_mock.IFileWriterMock
}

func TestServiceManager_CreateService(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{"my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common.GitContext{
		Project: "my-project",
		Credentials: &common.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()
	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateService(params)

	require.Nil(t, err)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
	require.Equal(t, fields.git.StageAndCommitAllCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.fileWriter.WriteFileCalls(), 1)

	md := &common.ServiceMetadata{}
	err = yaml.Unmarshal(fields.fileWriter.WriteFileCalls()[0].Content, md)

	require.Nil(t, err)
	require.Equal(t, md.ServiceName, params.ServiceName)
}

func getTestServiceManagerFields() serviceManagerTestFields {
	return serviceManagerTestFields{
		git: &common_mock.IGitMock{
			ProjectExistsFunc: func(gitContext common.GitContext) bool {
				return true
			},
			ProjectRepoExistsFunc: func(projectName string) bool {
				return true
			},
			CloneRepoFunc: func(gitContext common.GitContext) (bool, error) {
				return true, nil
			},
			StageAndCommitAllFunc: func(gitContext common.GitContext, message string) error {
				return nil
			},
			GetDefaultBranchFunc: func(gitContext common.GitContext) (string, error) {
				return "main", nil
			},
			CheckoutBranchFunc: func(gitContext common.GitContext, branch string) error {
				return nil
			},
		},
		credentialReader: &common_mock.CredentialReaderMock{
			GetCredentialsFunc: func(project string) (*common.GitCredentials, error) {
				return &common.GitCredentials{
					User:      "my-user",
					Token:     "my-token",
					RemoteURI: "my-remote-uri",
				}, nil
			},
		},
		fileWriter: &common_mock.IFileWriterMock{
			FileExistsFunc: func(path string) bool {
				return false
			},
			WriteBase64EncodedFileFunc: func(path string, content string) error {
				return nil
			},
			WriteFileFunc: func(path string, content []byte) error {
				return nil
			},
			DeleteFileFunc: func(path string) error {
				return nil
			},
			MakeDirFunc: func(path string) error {
				return nil
			},
		},
	}
}
