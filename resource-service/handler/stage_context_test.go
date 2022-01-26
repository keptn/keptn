package handler

import (
	"errors"
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

type testBranchStageContextFields struct {
	git        *common_mock.IGitMock
	fileSystem *common_mock.IFileSystemMock
}

func TestBranchStageContext_Establish_ProjectContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	bs := NewBranchStageContext(fields.git, fields.fileSystem)

	configPath, err := bs.Establish(models.Project{ProjectName: "my-project"}, nil, nil, common_models.GitContext{})

	require.Nil(t, err)

	require.Equal(t, common.GetProjectConfigPath("my-project"), configPath)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")
}

func TestBranchStageContext_Establish_StageContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	bs := NewBranchStageContext(fields.git, fields.fileSystem)

	configPath, err := bs.Establish(models.Project{ProjectName: "my-project"}, &models.Stage{StageName: "my-stage"}, nil, common_models.GitContext{})

	require.Nil(t, err)

	require.Equal(t, common.GetProjectConfigPath("my-project"), configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "my-stage")
}

func TestBranchStageContext_Establish_ServiceContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	bs := NewBranchStageContext(fields.git, fields.fileSystem)

	configPath, err := bs.Establish(models.Project{ProjectName: "my-project"}, &models.Stage{StageName: "my-stage"}, &models.Service{ServiceName: "my-service"}, common_models.GitContext{})

	require.Nil(t, err)

	require.Equal(t, common.GetServiceConfigPath("my-project", "my-service"), configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "my-stage")
}

func TestBranchStageContext_Establish_CannotGetDefaultBranch(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.git.GetDefaultBranchFunc = func(gitContext common_models.GitContext) (string, error) {
		return "", errors.New("oops")
	}

	bs := NewBranchStageContext(fields.git, fields.fileSystem)

	configPath, err := bs.Establish(models.Project{ProjectName: "my-project"}, nil, nil, common_models.GitContext{})

	require.NotNil(t, err)

	require.Empty(t, configPath)

	require.Empty(t, fields.git.CheckoutBranchCalls())
}

func TestBranchStageContext_Establish_CannotCheckoutBranch(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.git.CheckoutBranchFunc = func(gitContext common_models.GitContext, branch string) error {
		return errors.New("oops")
	}

	bs := NewBranchStageContext(fields.git, fields.fileSystem)

	configPath, err := bs.Establish(models.Project{ProjectName: "my-project"}, nil, nil, common_models.GitContext{})

	require.NotNil(t, err)

	require.Empty(t, configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
}

func TestBranchStageContext_Establish_ServiceDirectoryDoesNotExist(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	bs := NewBranchStageContext(fields.git, fields.fileSystem)

	configPath, err := bs.Establish(models.Project{ProjectName: "my-project"}, &models.Stage{StageName: "my-stage"}, &models.Service{ServiceName: "my-service"}, common_models.GitContext{})

	require.NotNil(t, err)
	require.ErrorIs(t, err, kerrors.ErrServiceNotFound)

	require.Empty(t, configPath)

}

func getTestBranchStageContextFields() testBranchStageContextFields {
	return testBranchStageContextFields{
		git: &common_mock.IGitMock{
			CheckoutBranchFunc: func(gitContext common_models.GitContext, branch string) error {
				return nil
			},
			CloneRepoFunc: func(gitContext common_models.GitContext) (bool, error) {
				return true, nil
			},
			CreateBranchFunc: func(gitContext common_models.GitContext, branch string, sourceBranch string) error {
				return nil
			},
			GetCurrentRevisionFunc: func(gitContext common_models.GitContext) (string, error) {
				return "my-revision", nil
			},
			GetDefaultBranchFunc: func(gitContext common_models.GitContext) (string, error) {
				return "main", nil
			},
			GetFileRevisionFunc: func(gitContext common_models.GitContext, revision string, file string) ([]byte, error) {
				return []byte("file-content"), nil
			},
			ProjectExistsFunc: func(gitContext common_models.GitContext) bool {
				return true
			},
			ProjectRepoExistsFunc: func(projectName string) bool {
				return true
			},
			PullFunc: func(gitContext common_models.GitContext) error {
				return nil
			},
			PushFunc: func(gitContext common_models.GitContext) error {
				return nil
			},
			StageAndCommitAllFunc: func(gitContext common_models.GitContext, message string) (string, error) {
				return "my-revision", nil
			},
		},
		fileSystem: &common_mock.IFileSystemMock{
			DeleteFileFunc: func(path string) error {
				return nil
			},
			FileExistsFunc: func(path string) bool {
				return true
			},
			MakeDirFunc: func(path string) error {
				return nil
			},
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte("file-content"), nil
			},
			WalkPathFunc: func(path string, walkFunc filepath.WalkFunc) error {

				_ = walkFunc(path+"/file1", newFakeFileInfo("file1", false), nil)
				_ = walkFunc(path+"/file2", newFakeFileInfo("file2", false), nil)
				_ = walkFunc(path+"/file3", newFakeFileInfo("file2", false), nil)

				return nil
			},
			WriteBase64EncodedFileFunc: func(path string, content string) error {
				return nil
			},
			WriteFileFunc: func(path string, content []byte) error {
				return nil
			},
			WriteHelmChartFunc: func(path string) error {
				return nil
			},
		},
	}
}
