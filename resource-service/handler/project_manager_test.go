package handler

import (
	"errors"
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

type projectManagerTestFields struct {
	git              *common_mock.IGitMock
	credentialReader *common_mock.CredentialReaderMock
	fileWriter       *common_mock.IFileSystemMock
}

func TestProjectManager_CreateProject(t *testing.T) {
	project := models.CreateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()
	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateProject(project)

	require.Nil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.ProjectRepoExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectRepoExistsCalls()[0].ProjectName, expectedGitContext.Project)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)
	require.Equal(t, fields.git.StageAndCommitAllCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.StageAndCommitAllCalls()[0].Message, "initialized project")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.fileWriter.WriteFileCalls(), 1)
	pmd := &common.ProjectMetadata{}
	err = yaml.Unmarshal(fields.fileWriter.WriteFileCalls()[0].Content, pmd)

	require.Nil(t, err)
	require.Equal(t, pmd.ProjectName, project.ProjectName)
}

func TestProjectManager_CreateProject_ProjectAlreadyExists(t *testing.T) {
	project := models.CreateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		if strings.Contains(path, "metadata.yaml") {
			return true
		}
		return false
	}
	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateProject(project)

	require.Equal(t, errors2.ErrProjectAlreadyExists, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.ProjectRepoExistsCalls())

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileWriter.WriteFileCalls())

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")
}

func TestProjectManager_CreateProject_CannotReadCredentials(t *testing.T) {
	project := models.CreateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	fields := getTestProjectManagerFields()

	fields.credentialReader.GetCredentialsFunc = func(project string) (*common_models.GitCredentials, error) {
		return nil, errors2.ErrMalformedCredentials
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateProject(project)

	require.ErrorIs(t, err, errors2.ErrMalformedCredentials)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Empty(t, fields.git.ProjectExistsCalls())
	require.Empty(t, fields.git.ProjectExistsCalls())
	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileWriter.WriteFileCalls())
	require.Empty(t, fields.fileWriter.DeleteFileCalls())
}

func TestProjectManager_CreateProject_ProjectRepoDoesNotExist(t *testing.T) {
	project := models.CreateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}
	fields.git.ProjectRepoExistsFunc = func(projectName string) bool {
		return false
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateProject(project)

	require.Equal(t, errors2.ErrRepositoryNotFound, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileWriter.WriteFileCalls())
}

func TestProjectManager_CreateProject_WritingFileFails(t *testing.T) {
	project := models.CreateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.fileWriter.WriteFileFunc = func(path string, content []byte) error {
		return errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Len(t, fields.fileWriter.WriteFileCalls(), 1)

	require.Len(t, fields.fileWriter.DeleteFileCalls(), 1)
	require.Equal(t, fields.fileWriter.DeleteFileCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName))
}

func TestProjectManager_CreateProject_CommitFails(t *testing.T) {
	project := models.CreateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.CreateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileWriter.WriteFileCalls(), 1)

	require.Len(t, fields.fileWriter.DeleteFileCalls(), 1)
	require.Equal(t, fields.fileWriter.DeleteFileCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName))
}

func TestProjectManager_UpdateProject(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.Nil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")
}

func TestProjectManager_UpdateProject_WithMigration(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
		Migrate: true,
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "metadata.yaml") {
			return []byte(`projectname: "sequence-queue3"`), nil
		}
		return []byte("content"), nil
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.Nil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.git.MigrateProjectCalls(), 1)
}

func TestProjectManager_UpdateProject_WithMigration_CannotPull(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
		Migrate: true,
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "metadata.yaml") {
			return []byte(`projectname: "sequence-queue3"`), nil
		}
		return []byte("content"), nil
	}

	fields.git.PullFunc = func(gitContext common_models.GitContext) error {
		return errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.git.MigrateProjectCalls(), 0)
}

func TestProjectManager_UpdateProject_WithMigration_MigrationFailsOnFirstTry(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
		Migrate: true,
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "metadata.yaml") {
			return []byte(`projectname: "sequence-queue3"`), nil
		}
		return []byte("content"), nil
	}

	nrTries := 0
	fields.git.MigrateProjectFunc = func(gitContext common_models.GitContext, newMetadatacontent []byte) error {
		if nrTries == 0 {
			nrTries++
			return errors.New("oops")
		}
		return nil
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.Nil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.git.MigrateProjectCalls(), 2)
}

func TestProjectManager_UpdateProject_WithMigration_AlreadyMigrated(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
		Migrate: true,
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "metadata.yaml") {
			return []byte(`projectName: "sequence-queue3"
isUsingDirectoryStructure: true`), nil
		}
		return []byte("content"), nil
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.Nil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.git.MigrateProjectCalls(), 0)
}

func TestProjectManager_UpdateProject_WithMigration_InvalidMetadata(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
		Migrate: true,
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "metadata.yaml") {
			// metadata.yaml with wrong indentation
			return []byte(`projectName: "sequence-queue3"
		isUsingDirectoryStructure: true`), nil
		}
		return []byte("content"), nil
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.git.MigrateProjectCalls(), 0)
}

func TestProjectManager_UpdateProject_WithMigration_NoMetadata(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
		Migrate: true,
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	nrTries := 0
	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		if nrTries == 0 {
			nrTries++
			return []byte("content"), nil
		}
		if strings.HasSuffix(filename, "metadata.yaml") {
			// metadata.yaml with wrong indentation
			return nil, errors.New("no file :(")
		}
		return []byte("content"), nil
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
	require.Equal(t, fields.fileWriter.FileExistsCalls()[0].Path, common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml")

	require.Len(t, fields.git.MigrateProjectCalls(), 0)
}

func TestProjectManager_UpdateProject_CannotReadCredentials(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	fields := getTestProjectManagerFields()

	fields.credentialReader.GetCredentialsFunc = func(project string) (*common_models.GitCredentials, error) {
		return nil, errors2.ErrMalformedCredentials
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.ErrorIs(t, err, errors2.ErrMalformedCredentials)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Empty(t, fields.git.ProjectExistsCalls())
	require.Empty(t, fields.git.GetDefaultBranchCalls())
	require.Empty(t, fields.git.CheckoutBranchCalls())
	require.Empty(t, fields.fileWriter.FileExistsCalls())
}

func TestProjectManager_UpdateProject_ProjectDoesNotExist(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.GetDefaultBranchCalls())
	require.Empty(t, fields.git.CheckoutBranchCalls())
	require.Empty(t, fields.fileWriter.FileExistsCalls())
}

func TestProjectManager_UpdateProject_ProjectNotInitialized(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}
	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		return nil, errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.GetDefaultBranchCalls())
	require.Empty(t, fields.git.CheckoutBranchCalls())
	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
}

func TestProjectManager_UpdateProject_ProjectNotInitializedEmptyMetadataFile(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}
	fields.fileWriter.ReadFileFunc = func(filename string) ([]byte, error) {
		return []byte(""), nil
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Empty(t, fields.git.GetDefaultBranchCalls())
	require.Empty(t, fields.git.CheckoutBranchCalls())
	require.Len(t, fields.fileWriter.FileExistsCalls(), 1)
}

func TestProjectManager_UpdateProject_CannotGetDefaultBranch(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}
	fields.git.GetDefaultBranchFunc = func(gitContext common_models.GitContext) (string, error) {
		return "", errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)
	require.Empty(t, fields.git.CheckoutBranchCalls())
}

func TestProjectManager_UpdateProject_CheckoutBranchFails(t *testing.T) {
	project := models.UpdateProjectParams{
		Project: models.Project{ProjectName: "my-project"},
	}

	expectedGitContext := common_models.GitContext{
		Project: "my-project",
		Credentials: &common_models.GitCredentials{
			User:      "my-user",
			Token:     "my-token",
			RemoteURI: "my-remote-uri",
		},
	}

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}
	fields.git.CheckoutBranchFunc = func(gitContext common_models.GitContext, branch string) error {
		return errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.UpdateProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.credentialReader.GetCredentialsCalls(), 1)
	require.Equal(t, fields.credentialReader.GetCredentialsCalls()[0].Project, project.ProjectName)

	require.Len(t, fields.git.ProjectExistsCalls(), 1)
	require.Equal(t, fields.git.ProjectExistsCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)
	require.Equal(t, fields.git.GetDefaultBranchCalls()[0].GitContext, expectedGitContext)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].GitContext, expectedGitContext)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")
}

func TestProjectManager_DeleteProject(t *testing.T) {
	project := "my-project"

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.DeleteProject(project)

	require.Nil(t, err)

	require.Len(t, fields.fileWriter.DeleteFileCalls(), 1)
	require.Equal(t, fields.fileWriter.DeleteFileCalls()[0].Path, common.GetProjectConfigPath(project))
}

func TestProjectManager_DeleteProject_CannotDeleteDirectory(t *testing.T) {
	project := "my-project"

	fields := getTestProjectManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return true
	}
	fields.fileWriter.FileExistsFunc = func(path string) bool {
		return true
	}
	fields.fileWriter.DeleteFileFunc = func(path string) error {
		if strings.Contains(path, "metadata") {
			return nil
		}
		return errors.New("oops")
	}

	p := NewProjectManager(fields.git, fields.credentialReader, fields.fileWriter)
	err := p.DeleteProject(project)

	require.NotNil(t, err)

	require.Len(t, fields.fileWriter.DeleteFileCalls(), 1)
}

func getTestProjectManagerFields() projectManagerTestFields {
	return projectManagerTestFields{
		git: &common_mock.IGitMock{
			ResetHardFunc:         func(gitContext common_models.GitContext) error { return nil },
			ProjectExistsFunc:     func(gitContext common_models.GitContext) bool { return true },
			ProjectRepoExistsFunc: func(projectName string) bool { return true },
			CloneRepoFunc:         func(gitContext common_models.GitContext) (bool, error) { return true, nil },
			StageAndCommitAllFunc: func(gitContext common_models.GitContext, message string) (string, error) { return "", nil },
			GetDefaultBranchFunc:  func(gitContext common_models.GitContext) (string, error) { return "main", nil },
			CheckoutBranchFunc:    func(gitContext common_models.GitContext, branch string) error { return nil },
			MigrateProjectFunc:    func(gitContext common_models.GitContext, newMetadatacontent []byte) error { return nil },
			PullFunc:              func(gitContext common_models.GitContext) error { return nil },
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
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte("content"), nil
			},
		},
	}
}
