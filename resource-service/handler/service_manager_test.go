package handler

import (
	"errors"
	"strings"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"

	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	handler_mock "github.com/keptn/keptn/resource-service/handler/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type serviceManagerTestFields struct {
	git                  *common_mock.IGitMock
	credentialReader     *common_mock.CredentialReaderMock
	fileWriter           *common_mock.IFileSystemMock
	configurationContext *handler_mock.IConfigurationContextMock
}

func TestServiceManager_CreateService(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.Nil(t, err)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
	require.Equal(t, fields.git.StageAndCommitAllCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.fileWriter.WriteFileCalls(), 1)
	require.Equal(t, fields.fileWriter.WriteFileCalls()[0].Path, testServiceConfigDir+"/metadata.yaml")

	md := &common.ServiceMetadata{}
	err = yaml.Unmarshal(fields.fileWriter.WriteFileCalls()[0].Content, md)

	require.Nil(t, err)
	require.Equal(t, md.ServiceName, params.ServiceName)
}

func TestServiceManager_CreateService_CannotReadCredentials(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	fields := getTestServiceManagerFields()

	fields.credentialReader.GetCredentialsFunc = func(project string) (*common_models.GitCredentials, error) {
		return nil, errors2.ErrCredentialsNotFound
	}
	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.ErrorIs(t, err, errors2.ErrCredentialsNotFound)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileWriter.WriteFileCalls())
}

func TestServiceManager_CreateService_ProjectNotFound(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileWriter.WriteFileCalls())
}

func TestServiceManager_CreateService_StageNotFound(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	fields.configurationContext.EstablishFunc = func(params common_models.ConfigurationContextParams) (string, error) {
		return "", errors2.ErrStageNotFound
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.ErrorIs(t, err, errors2.ErrStageNotFound)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileWriter.WriteFileCalls())
}

func TestServiceManager_CreateService_ServiceAlreadyExists(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.ErrorIs(t, err, errors2.ErrServiceAlreadyExists)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileWriter.WriteFileCalls())
}

func TestServiceManager_CreateService_CannotCreateDirectory(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()
	fields.fileWriter.MakeDirFunc = func(path string) error {
		return errors.New("oops")
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileWriter.WriteFileCalls())
}

func TestServiceManager_CreateService_CannotCreateMetadata(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()
	fields.fileWriter.WriteFileFunc = func(path string, content []byte) error {
		return errors.New("oops")
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Len(t, fields.fileWriter.WriteFileCalls(), 1)
}

func TestServiceManager_CreateService_CannotCommit(t *testing.T) {
	params := models.CreateServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		CreateServicePayload: models.CreateServicePayload{
			Service: models.Service{
				ServiceName: "my-service",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()
	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.CreateService(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
	require.Equal(t, fields.git.StageAndCommitAllCalls()[0].GitContext, expectedGitContext)
}

func TestServiceManager_DeleteService(t *testing.T) {
	params := models.DeleteServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		Service: models.Service{ServiceName: "my-service"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		if strings.Contains(path, "my-service") {
			return true
		}
		return false
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.DeleteService(params)

	require.Nil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
	require.Equal(t, fields.git.StageAndCommitAllCalls()[0].GitContext, expectedGitContext)
}

func TestServiceManager_DeleteService_ProjectDoesNotExist(t *testing.T) {
	params := models.DeleteServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		Service: models.Service{ServiceName: "my-service"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		if strings.Contains(path, "my-service") {
			return true
		}
		return false
	}

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.DeleteService(params)

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.CheckoutBranchCalls())
	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestServiceManager_DeleteService_ServiceDoesNotExist(t *testing.T) {
	params := models.DeleteServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		Service: models.Service{ServiceName: "my-service"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return false
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.DeleteService(params)

	require.ErrorIs(t, err, errors2.ErrServiceNotFound)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestServiceManager_DeleteService_DeleteDirectoryFails(t *testing.T) {
	params := models.DeleteServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		Service: models.Service{ServiceName: "my-service"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		if strings.Contains(path, "my-service") {
			return true
		}
		return false
	}
	fields.fileWriter.DeleteFileFunc = func(path string) error {
		return errors.New("oops")
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.DeleteService(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestServiceManager_DeleteService_CannotCommit(t *testing.T) {
	params := models.DeleteServiceParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
		Service: models.Service{ServiceName: "my-service"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &apimodels.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestServiceManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		if strings.Contains(path, "my-service") {
			return true
		}
		return false
	}
	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	p := NewServiceManager(fields.git, fields.credentialReader, fields.fileWriter, fields.configurationContext)
	err := p.DeleteService(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.configurationContext.EstablishCalls(), 1)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
}

func getTestServiceManagerFields() serviceManagerTestFields {
	return serviceManagerTestFields{
		git: &common_mock.IGitMock{
			ResetHardFunc:         func(gitContext common_models.GitContext) error { return nil },
			PullFunc:              func(gitContext common_models.GitContext) error { return nil },
			ProjectExistsFunc:     func(gitContext common_models.GitContext) bool { return true },
			ProjectRepoExistsFunc: func(projectName string) bool { return true },
			CloneRepoFunc:         func(gitContext common_models.GitContext) (bool, error) { return true, nil },
			StageAndCommitAllFunc: func(gitContext common_models.GitContext, message string) (string, error) { return "", nil },
			GetDefaultBranchFunc:  func(gitContext common_models.GitContext) (string, error) { return "main", nil },
			CheckoutBranchFunc:    func(gitContext common_models.GitContext, branch string) error { return nil },
		},
		credentialReader: &common_mock.CredentialReaderMock{
			GetCredentialsFunc: func(project string) (*common_models.GitCredentials, error) {
				return &common_models.GitCredentials{
					User: "my-user",
					HttpsAuth: &apimodels.HttpsGitAuth{
						Token: "my-token",
					},
					RemoteURL: "my-remote-uri",
				}, nil
			},
		},
		fileWriter: &common_mock.IFileSystemMock{
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
		configurationContext: &handler_mock.IConfigurationContextMock{
			EstablishFunc: func(params common_models.ConfigurationContextParams) (string, error) {
				return testServiceConfigDir, nil
			},
		},
	}
}
