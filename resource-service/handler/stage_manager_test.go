package handler

import (
	"errors"
	"testing"

	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	handler_mock "github.com/keptn/keptn/resource-service/handler/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
)

const testStageConfigDir = "/data/config/my-project/.keptn-stages/my-stage"

type stageManagerTestFields struct {
	git                  *common_mock.IGitMock
	credentialReader     *common_mock.CredentialReaderMock
	configurationContext *handler_mock.IConfigurationContextMock
	fileSystem           *common_mock.IFileSystemMock
}

func TestStageManager_CreateStage(t *testing.T) {
	params := models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{
				StageName: "my-stage",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &common_models.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestStageManagerFields()
	s := NewStageManager(fields.git, fields.credentialReader)
	err := s.CreateStage(params)

	require.Nil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CreateBranchCalls(), 1)
	require.Equal(t, fields.git.CreateBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CreateBranchCalls()[0].SourceBranch, "main")
	require.Equal(t, fields.git.CreateBranchCalls()[0].Branch, "my-stage")

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
}

func TestStageManager_CreateStage_NoCredentialsFound(t *testing.T) {
	params := models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{
				StageName: "my-stage",
			},
		},
	}

	fields := getTestStageManagerFields()

	fields.credentialReader.GetCredentialsFunc = func(project string) (*common_models.GitCredentials, error) {
		return nil, errors2.ErrCredentialsNotFound
	}
	s := NewStageManager(fields.git, fields.credentialReader)
	err := s.CreateStage(params)

	require.ErrorIs(t, err, errors2.ErrCredentialsNotFound)

	require.Empty(t, fields.git.CreateBranchCalls())
}

func TestStageManager_CreateStage_ProjectDoesNotExist(t *testing.T) {
	params := models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{
				StageName: "my-stage",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &common_models.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestStageManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	s := NewStageManager(fields.git, fields.credentialReader)
	err := s.CreateStage(params)

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.CreateBranchCalls())
}

func TestStageManager_CreateStage_CannotGetDefaultBranch(t *testing.T) {
	params := models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{
				StageName: "my-stage",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &common_models.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestStageManagerFields()

	fields.git.GetDefaultBranchFunc = func(gitContext common_models.GitContext) (string, error) {
		return "", errors.New("oops")
	}

	s := NewStageManager(fields.git, fields.credentialReader)
	err := s.CreateStage(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.CreateBranchCalls())
}

func TestStageManager_CreateStage_CannotCreateBranch(t *testing.T) {
	params := models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{
				StageName: "my-stage",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &common_models.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestStageManagerFields()

	fields.git.CreateBranchFunc = func(gitContext common_models.GitContext, branch string, sourceBranch string) error {
		return errors2.ErrStageAlreadyExists
	}

	s := NewStageManager(fields.git, fields.credentialReader)
	err := s.CreateStage(params)

	require.ErrorIs(t, err, errors2.ErrStageAlreadyExists)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CreateBranchCalls(), 1)
	require.Equal(t, fields.git.CreateBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CreateBranchCalls()[0].SourceBranch, "main")
	require.Equal(t, fields.git.CreateBranchCalls()[0].Branch, "my-stage")
}

func TestStageManager_CreateStage_CannotPushBranch(t *testing.T) {
	params := models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{
				StageName: "my-stage",
			},
		},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User: "my-user",
			HttpsAuth: &common_models.HttpsGitAuth{
				Token: "my-token",
			},
			RemoteURL: "my-remote-uri",
		},
	}

	fields := getTestStageManagerFields()

	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	s := NewStageManager(fields.git, fields.credentialReader)
	err := s.CreateStage(params)

	require.NotNil(t, err)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CreateBranchCalls(), 1)
	require.Equal(t, fields.git.CreateBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CreateBranchCalls()[0].SourceBranch, "main")
	require.Equal(t, fields.git.CreateBranchCalls()[0].Branch, "my-stage")
}

func getTestStageManagerFields() stageManagerTestFields {
	return stageManagerTestFields{
		git: &common_mock.IGitMock{
			ProjectExistsFunc: func(gitContext common_models.GitContext) bool {
				return true
			},
			ProjectRepoExistsFunc: func(projectName string) bool {
				return true
			},
			CloneRepoFunc: func(gitContext common_models.GitContext) (bool, error) {
				return true, nil
			},
			StageAndCommitAllFunc: func(gitContext common_models.GitContext, message string) (string, error) {
				return "", nil
			},
			GetDefaultBranchFunc: func(gitContext common_models.GitContext) (string, error) {
				return "main", nil
			},
			CheckoutBranchFunc: func(gitContext common_models.GitContext, branch string) error {
				return nil
			},
			CreateBranchFunc: func(gitContext common_models.GitContext, branch string, sourceBranch string) error {
				return nil
			},
		},
		credentialReader: &common_mock.CredentialReaderMock{
			GetCredentialsFunc: func(project string) (*common_models.GitCredentials, error) {
				return &common_models.GitCredentials{
					User: "my-user",
					HttpsAuth: &common_models.HttpsGitAuth{
						Token: "my-token",
					},
					RemoteURL: "my-remote-uri",
				}, nil
			},
		},
		configurationContext: &handler_mock.IConfigurationContextMock{EstablishFunc: func(params common_models.ConfigurationContextParams) (string, error) {
			return testStageConfigDir, nil
		}},
		fileSystem: &common_mock.IFileSystemMock{
			DeleteFileFunc: func(path string) error {
				return nil
			},
			FileExistsFunc: func(path string) bool {
				return true
			},
			WriteFileFunc: func(path string, content []byte) error {
				return nil
			},
			MakeDirFunc: func(path string) error {
				return nil
			},
		},
	}
}

func TestDirectoryStageManager_CreateStage(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.Nil(t, err)
}

func TestDirectoryStageManager_CreateStage_CannotEstablishContext(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.configurationContext.EstablishFunc = func(params common_models.ConfigurationContextParams) (string, error) {
		return "", errors.New("oops")
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.NotNil(t, err)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestDirectoryStageManager_CreateStage_CannotGetCredentials(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.credentialReader.GetCredentialsFunc = func(project string) (*common_models.GitCredentials, error) {
		return nil, errors.New("oops")
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.NotNil(t, err)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestDirectoryStageManager_CreateStage_ProjectNotFound(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestDirectoryStageManager_CreateStage_StageAlreadyExists(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return true
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.ErrorIs(t, err, errors2.ErrStageAlreadyExists)
}

func TestDirectoryStageManager_CreateStage_CannotCreateDirectory(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	fields.fileSystem.MakeDirFunc = func(path string) error {
		return errors.New("oops")
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.NotNil(t, err)
}

func TestDirectoryStageManager_CreateStage_CannotWriteMetadata(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	fields.fileSystem.WriteFileFunc = func(path string, content []byte) error {
		return errors.New("oops")
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.NotNil(t, err)
}

func TestDirectoryStageManager_CreateStage_CannotCommit(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.CreateStage(models.CreateStageParams{
		Project: models.Project{ProjectName: "my-project"},
		CreateStagePayload: models.CreateStagePayload{
			Stage: models.Stage{StageName: "my-stage"},
		},
	})

	require.NotNil(t, err)
}

func TestDirectoryStageManager_DeleteStage(t *testing.T) {
	fields := getTestStageManagerFields()

	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.DeleteStage(models.DeleteStageParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
	})

	require.Nil(t, err)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
}

func TestDirectoryStageManager_DeleteStage_CannotEstablishContext(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.configurationContext.EstablishFunc = func(params common_models.ConfigurationContextParams) (string, error) {
		return "", errors.New("oops")
	}
	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.DeleteStage(models.DeleteStageParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
	})

	require.NotNil(t, err)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestDirectoryStageManager_DeleteStage_StageDirectoryNotAvailable(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}
	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.DeleteStage(models.DeleteStageParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
	})

	require.ErrorIs(t, err, errors2.ErrStageNotFound)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestDirectoryStageManager_DeleteStage_CannotDeleteDirectory(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.fileSystem.DeleteFileFunc = func(path string) error {
		return errors.New("oops")
	}
	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.DeleteStage(models.DeleteStageParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
	})

	require.NotNil(t, err)

	require.Empty(t, fields.git.StageAndCommitAllCalls())
}

func TestDirectoryStageManager_DeleteStage_CannotCommitChanges(t *testing.T) {
	fields := getTestStageManagerFields()

	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}
	dm := NewDirectoryStageManager(fields.configurationContext, fields.fileSystem, fields.credentialReader, fields.git)

	err := dm.DeleteStage(models.DeleteStageParams{
		Project: models.Project{ProjectName: "my-project"},
		Stage:   models.Stage{StageName: "my-stage"},
	})

	require.NotNil(t, err)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
}
