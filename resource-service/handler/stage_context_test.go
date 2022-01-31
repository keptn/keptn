package handler

import (
	"errors"
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"testing"
)

type testBranchStageContextFields struct {
	git        *common_mock.IGitMock
	fileSystem *common_mock.IFileSystemMock
}

func TestBranchStageContext_Establish_ProjectContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configPath, err := bs.Establish(params)

	require.Nil(t, err)

	require.Equal(t, common.GetProjectConfigPath("my-project"), configPath)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Empty(t, fields.fileSystem.FileExistsCalls())
}

func TestBranchStageContext_Establish_ProjectContext_ProjectDirectoryNotAvailable(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: true,
	}
	configPath, err := bs.Establish(params)

	require.ErrorIs(t, err, kerrors.ErrProjectNotFound)

	require.Equal(t, "", configPath)

	require.Len(t, fields.git.GetDefaultBranchCalls(), 1)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "main")

	require.Len(t, fields.fileSystem.FileExistsCalls(), 1)
}

func TestBranchStageContext_Establish_StageContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}

	configPath, err := bs.Establish(params)

	require.Nil(t, err)

	require.Equal(t, common.GetProjectConfigPath("my-project"), configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "my-stage")
}

func TestBranchStageContext_Establish_ServiceContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 &models.Service{ServiceName: "my-service"},
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}

	configPath, err := bs.Establish(params)

	require.Nil(t, err)

	require.Equal(t, common.GetServiceConfigPath("my-project", "my-service"), configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "my-stage")
}

func TestBranchStageContext_Establish_ServiceContext_ServiceDirectoryNotAvailable(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 &models.Service{ServiceName: "my-service"},
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: true,
	}
	configPath, err := bs.Establish(params)

	require.ErrorIs(t, err, kerrors.ErrServiceNotFound)

	require.Equal(t, "", configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
	require.Equal(t, fields.git.CheckoutBranchCalls()[0].Branch, "my-stage")

	require.Len(t, fields.fileSystem.FileExistsCalls(), 1)
}

func TestBranchStageContext_Establish_CannotGetDefaultBranch(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.git.GetDefaultBranchFunc = func(gitContext common_models.GitContext) (string, error) {
		return "", errors.New("oops")
	}

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configPath, err := bs.Establish(params)

	require.NotNil(t, err)

	require.Empty(t, configPath)

	require.Empty(t, fields.git.CheckoutBranchCalls())
}

func TestBranchStageContext_Establish_CannotCheckoutBranch(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.git.CheckoutBranchFunc = func(gitContext common_models.GitContext, branch string) error {
		return errors.New("oops")
	}

	bs := NewBranchConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configPath, err := bs.Establish(params)

	require.NotNil(t, err)

	require.Empty(t, configPath)

	require.Len(t, fields.git.CheckoutBranchCalls(), 1)
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
			FileExistsFunc: func(path string) bool {
				return true
			},
		},
	}
}

func TestDirectoryConfigurationContext_Establish_ProjectContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configDir, err := ds.Establish(params)

	require.Nil(t, err)

	require.Equal(t, common.GetProjectConfigPath("my-project"), configDir)
}

func TestDirectoryConfigurationContext_Establish_ProjectContext_ProjectNotFound(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}
	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: true,
	}
	configDir, err := ds.Establish(params)

	require.ErrorIs(t, err, kerrors.ErrProjectNotFound)

	require.Equal(t, "", configDir)
}

func TestDirectoryConfigurationContext_Establish_StageContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configDir, err := ds.Establish(params)

	require.Nil(t, err)

	require.Equal(t, common.GetConfigDir()+"/my-project/.keptn-stages/my-stage", configDir)
}
func TestDirectoryConfigurationContext_Establish_StageContext_StageNotFound(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}
	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: true,
	}
	configDir, err := ds.Establish(params)

	require.ErrorIs(t, err, kerrors.ErrStageNotFound)

	require.Equal(t, "", configDir)
}

func TestDirectoryConfigurationContext_Establish_ServiceContext(t *testing.T) {
	fields := getTestBranchStageContextFields()

	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 &models.Service{ServiceName: "my-service"},
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configDir, err := ds.Establish(params)

	require.Nil(t, err)

	require.Equal(t, common.GetConfigDir()+"/my-project/.keptn-stages/my-stage/my-service", configDir)
}

func TestDirectoryConfigurationContext_Establish_ServiceContext_ServiceNotFound(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		return false
	}
	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   &models.Stage{StageName: "my-stage"},
		Service:                 &models.Service{ServiceName: "my-service"},
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: true,
	}
	configDir, err := ds.Establish(params)

	require.ErrorIs(t, err, kerrors.ErrServiceNotFound)

	require.Equal(t, "", configDir)
}

func TestDirectoryConfigurationContext_Establish_CannotDetermineDefaultBranch(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.git.GetDefaultBranchFunc = func(gitContext common_models.GitContext) (string, error) {
		return "", errors.New("oops")
	}
	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configDir, err := ds.Establish(params)

	require.NotNil(t, err)

	require.Equal(t, "", configDir)
}

func TestDirectoryConfigurationContext_Establish_CannotCheckoutDefaultBranch(t *testing.T) {
	fields := getTestBranchStageContextFields()

	fields.git.CheckoutBranchFunc = func(gitContext common_models.GitContext, branch string) error {
		return errors.New("oops")
	}
	ds := NewDirectoryConfigurationContext(fields.git, fields.fileSystem)

	params := common_models.ConfigurationContextParams{
		Project:                 models.Project{ProjectName: "my-project"},
		Stage:                   nil,
		Service:                 nil,
		GitContext:              common_models.GitContext{},
		CheckConfigDirAvailable: false,
	}
	configDir, err := ds.Establish(params)

	require.NotNil(t, err)

	require.Equal(t, "", configDir)
}
