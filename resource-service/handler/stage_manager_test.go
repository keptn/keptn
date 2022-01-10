package handler

import (
	"errors"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"testing"
)

type stageManagerTestFields struct {
	git              *common_mock.IGitMock
	credentialReader *common_mock.CredentialReaderMock
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
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
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
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
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
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
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
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
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
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
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
					User:      "my-user",
					Token:     "my-token",
					RemoteURI: "my-remote-uri",
				}, nil
			},
		},
	}
}
