package handler

import (
	"errors"
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	handler_mock "github.com/keptn/keptn/resource-service/handler/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testConfigDir = "/data/config/my-project"
const testServiceConfigDir = "/data/config/my-project/my-service"

type testResourceManagerFields struct {
	git              *common_mock.IGitMock
	credentialReader *common_mock.CredentialReaderMock
	fileSystem       *common_mock.IFileSystemMock
	stageContext     *handler_mock.IConfigurationContextMock
}

func TestResourceManager_CreateResources_ProjectResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Equal(t, models.Project{ProjectName: "my-project"}, fields.stageContext.EstablishCalls()[0].Params.Project, 1)
	require.Nil(t, fields.stageContext.EstablishCalls()[0].Params.Stage, 1)
	require.Nil(t, fields.stageContext.EstablishCalls()[0].Params.Service, 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 2)
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path, common.GetProjectConfigPath("my-project")+"/file1")
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[1].Path, common.GetProjectConfigPath("my-project")+"/file2")
}

func TestResourceManager_CreateResources_StageResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Equal(t, models.Project{ProjectName: "my-project"}, fields.stageContext.EstablishCalls()[0].Params.Project, 1)
	require.Equal(t, &models.Stage{StageName: "my-stage"}, fields.stageContext.EstablishCalls()[0].Params.Stage, 1)
	require.Nil(t, fields.stageContext.EstablishCalls()[0].Params.Service, 1)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 2)
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path, common.GetProjectConfigPath("my-project")+"/file1")
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[1].Path, common.GetProjectConfigPath("my-project")+"/file2")
}

func TestResourceManager_CreateResources_ServiceResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Equal(t, models.Project{ProjectName: "my-project"}, fields.stageContext.EstablishCalls()[0].Params.Project, 1)
	require.Equal(t, &models.Stage{StageName: "my-stage"}, fields.stageContext.EstablishCalls()[0].Params.Stage, 1)
	require.Equal(t, &models.Service{ServiceName: "my-service"}, fields.stageContext.EstablishCalls()[0].Params.Service, 1)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 2)
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path, testConfigDir+"/file1")
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[1].Path, testConfigDir+"/file2")
}

func TestResourceManager_CreateResources_ServiceResource_HelmChart(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "helm/service.tgz",
				},
			},
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 1)
	require.Equal(t, fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path, testConfigDir+"/helm/service.tgz")

	require.Len(t, fields.fileSystem.WriteHelmChartCalls(), 1)
	require.Equal(t, fields.fileSystem.WriteHelmChartCalls()[0].Path, testConfigDir+"/helm/service.tgz")
}

func TestResourceManager_CreateResources_ServiceResource_HelmChartWriteFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.WriteHelmChartFunc = func(path string) error {
		return errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "helm/service.tgz",
				},
			},
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 0)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 1)
	require.Equal(t, testConfigDir+"/helm/service.tgz", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)

	require.Len(t, fields.fileSystem.WriteHelmChartCalls(), 1)
	require.Equal(t, testConfigDir+"/helm/service.tgz", fields.fileSystem.WriteHelmChartCalls()[0].Path)
}

func TestResourceManager_CreateResources_ProjectResource_ProjectNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Nil(t, revision)

	require.Empty(t, fields.git.CheckoutBranchCalls())

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileSystem.WriteBase64EncodedFileCalls())
}

func TestResourceManager_CreateResources_ProjectResource_CannotReadCredentials(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.credentialReader.GetCredentialsFunc = func(project string) (*common_models.GitCredentials, error) {
		return nil, errors2.ErrMalformedCredentials
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.ErrorIs(t, err, errors2.ErrMalformedCredentials)

	require.Nil(t, revision)

	require.Empty(t, fields.git.CheckoutBranchCalls())

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileSystem.WriteBase64EncodedFileCalls())
}

func TestResourceManager_CreateResources_ProjectResource_CannotEstablishContext(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.stageContext.EstablishFunc = func(params common_models.ConfigurationContextParams) (string, error) {
		return "", errors.New("oops")
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.CreateResources(models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		CreateResourcesPayload: models.CreateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileSystem.WriteBase64EncodedFileCalls())
}

func TestResourceManager_UpdateResources_ProjectResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResources(models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		UpdateResourcesPayload: models.UpdateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 2)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
	require.Equal(t, testConfigDir+"/file2", fields.fileSystem.WriteBase64EncodedFileCalls()[1].Path)
}

func TestResourceManager_UpdateResources_ProjectResource_ProjectNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResources(models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		UpdateResourcesPayload: models.UpdateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Nil(t, revision)

	require.Empty(t, fields.git.CheckoutBranchCalls())

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileSystem.WriteBase64EncodedFileCalls())
}

func TestResourceManager_UpdateResources_ProjectResource_WritingFileFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.WriteBase64EncodedFileFunc = func(path string, content string) error {
		return errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResources(models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		UpdateResourcesPayload: models.UpdateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
}

func TestResourceManager_UpdateResources_ProjectResource_CommitFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResources(models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		UpdateResourcesPayload: models.UpdateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 2)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
	require.Equal(t, testConfigDir+"/file2", fields.fileSystem.WriteBase64EncodedFileCalls()[1].Path)
}

func TestResourceManager_UpdateResources_ProjectResource_PullFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.PullFunc = func(gitContext common_models.GitContext) error {
		return errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResources(models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		UpdateResourcesPayload: models.UpdateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.git.PullCalls(), 1)
	require.Len(t, fields.git.StageAndCommitAllCalls(), 0)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 0)
}

func TestResourceManager_UpdateResources_ProjectResource_CommitFailsOnFirstTry(t *testing.T) {
	fields := getTestResourceManagerFields()

	i := 0
	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		if i == 0 {
			i++
			return "", errors2.ErrNonFastForwardUpdate
		}
		return "my-revision", nil
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResources(models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		UpdateResourcesPayload: models.UpdateResourcesPayload{
			Resources: []models.Resource{
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file1",
				},
				{
					ResourceContent: "c3RyaW5n",
					ResourceURI:     "file2",
				},
			},
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 2)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 4)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
	require.Equal(t, testConfigDir+"/file2", fields.fileSystem.WriteBase64EncodedFileCalls()[1].Path)
}

func TestResourceManager_UpdateResource_ProjectResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResource(models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		UpdateResourcePayload: models.UpdateResourcePayload{
			ResourceContent: "c3RyaW5n",
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
}

func TestResourceManager_UpdateResource_ProjectResource_ProjectNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResource(models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		UpdateResourcePayload: models.UpdateResourcePayload{
			ResourceContent: "c3RyaW5n",
		},
	})

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Nil(t, revision)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Empty(t, fields.fileSystem.WriteBase64EncodedFileCalls())
}

func TestResourceManager_UpdateResource_ProjectResource_WritingFileFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.WriteBase64EncodedFileFunc = func(path string, content string) error {
		return errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResource(models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		UpdateResourcePayload: models.UpdateResourcePayload{
			ResourceContent: "c3RyaW5n",
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
}

func TestResourceManager_UpdateResource_ProjectResource_CommitFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		return "", errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResource(models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		UpdateResourcePayload: models.UpdateResourcePayload{
			ResourceContent: "c3RyaW5n",
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
}

func TestResourceManager_UpdateResource_ProjectResource_PullFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.PullFunc = func(gitContext common_models.GitContext) error {
		return errors.New("oops")
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResource(models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		UpdateResourcePayload: models.UpdateResourcePayload{
			ResourceContent: "c3RyaW5n",
		},
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.git.PullCalls(), 1)
	require.Len(t, fields.git.StageAndCommitAllCalls(), 0)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 0)
}

func TestResourceManager_UpdateResource_ProjectResource_CommitFailsOnFirstTry(t *testing.T) {
	fields := getTestResourceManagerFields()

	i := 0
	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		if i == 0 {
			i++
			return "", errors2.ErrNonFastForwardUpdate
		}
		return "my-revision", nil
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.UpdateResource(models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		UpdateResourcePayload: models.UpdateResourcePayload{
			ResourceContent: "c3RyaW5n",
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 2)

	require.Len(t, fields.fileSystem.WriteBase64EncodedFileCalls(), 2)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.WriteBase64EncodedFileCalls()[0].Path)
}

func TestResourceManager_DeleteResource_ProjectResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.DeleteResource(models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 1)

	require.Len(t, fields.fileSystem.DeleteFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.DeleteFileCalls()[0].Path)
}

func TestResourceManager_DeleteResource_ProjectResource_ProjectNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.DeleteResource(models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.ErrorIs(t, err, errors2.ErrProjectNotFound)

	require.Nil(t, revision)

	require.Empty(t, fields.stageContext.EstablishCalls())
	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileSystem.DeleteFileCalls())
}

func TestResourceManager_DeleteResource_ProjectResource_ResourceNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.FileExistsFunc = func(path string) bool {
		if strings.Contains(path, "file1") {
			return false
		}
		return true
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.DeleteResource(models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.ErrorIs(t, err, errors2.ErrResourceNotFound)

	require.Nil(t, revision)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Empty(t, fields.git.StageAndCommitAllCalls())
	require.Empty(t, fields.fileSystem.DeleteFileCalls())
}

func TestResourceManager_DeleteResource_ProjectResource_DeleteFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.DeleteFileFunc = func(path string) error {
		return errors.New("oops")
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.DeleteResource(models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Len(t, fields.fileSystem.DeleteFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.DeleteFileCalls()[0].Path)
}

func TestResourceManager_DeleteResource_Git_ThrowsNonFastForward_Or_ThrowsForceNeeded(t *testing.T) {
	fields := getTestResourceManagerFields()
	fields.git.StageAndCommitAllFunc = func(gitContext common_models.GitContext, message string) (string, error) {
		stageAndCommitCalls := fields.git.StageAndCommitAllCalls()
		if len(stageAndCommitCalls) == 1 {
			return "", errors2.ErrNonFastForwardUpdate
		}
		if len(stageAndCommitCalls) == 2 {
			return "", errors2.ErrForceNeeded
		}

		return "my-revision", nil
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.DeleteResource(models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.Nil(t, err)

	require.Equal(t, &models.WriteResourceResponse{CommitID: "my-revision"}, revision)

	require.Len(t, fields.git.StageAndCommitAllCalls(), 3)
	require.Len(t, fields.fileSystem.DeleteFileCalls(), 3)
	require.Len(t, fields.git.PullCalls(), 3)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.DeleteFileCalls()[0].Path)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.DeleteFileCalls()[1].Path)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.DeleteFileCalls()[2].Path)
}

func TestResourceManager_DeleteResource_ProjectResource_CommitFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.DeleteFileFunc = func(path string) error {
		return errors.New("oops")
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	revision, err := rm.DeleteResource(models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.NotNil(t, err)

	require.Nil(t, revision)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.StageAndCommitAllCalls())

	require.Len(t, fields.fileSystem.DeleteFileCalls(), 1)
	require.Equal(t, testConfigDir+"/file1", fields.fileSystem.DeleteFileCalls()[0].Path)
}

func TestResourceManager_GetResource_ProjectResource(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.Nil(t, err)

	require.Equal(t, &models.GetResourceResponse{
		Resource: models.Resource{
			ResourceContent: "ZmlsZS1jb250ZW50",
			ResourceURI:     "file1",
		},
		Metadata: models.Version{
			UpstreamURL: "remote-url",
			Version:     "my-revision",
		},
	}, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Len(t, fields.git.PullCalls(), 1)

	require.Empty(t, fields.git.GetFileRevisionCalls())
	require.Len(t, fields.fileSystem.ReadFileCalls(), 1)
}

func TestResourceManager_GetResource_ProjectResource_ProvideGitCommitID(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		GetResourceQuery: models.GetResourceQuery{
			GitCommitID: "my-commit-id",
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.GetResourceResponse{
		Resource: models.Resource{
			ResourceContent: "ZmlsZS1jb250ZW50",
			ResourceURI:     "file1",
		},
		Metadata: models.Version{
			UpstreamURL: "remote-url",
			Version:     "my-commit-id",
		},
	}, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Len(t, fields.git.GetFileRevisionCalls(), 1)
	require.Equal(t, "my-commit-id", fields.git.GetFileRevisionCalls()[0].Revision)
	require.Equal(t, "file1", fields.git.GetFileRevisionCalls()[0].File)
}

func TestResourceManager_GetResource_ServiceResource_ProvideGitCommitID(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.stageContext.EstablishFunc = func(params common_models.ConfigurationContextParams) (string, error) {
		return testConfigDir + "/my-service", nil
	}

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		ResourceURI: "file1",
		GetResourceQuery: models.GetResourceQuery{
			GitCommitID: "my-commit-id",
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.GetResourceResponse{
		Resource: models.Resource{
			ResourceContent: "ZmlsZS1jb250ZW50",
			ResourceURI:     "file1",
		},
		Metadata: models.Version{
			UpstreamURL: "remote-url",
			Version:     "my-commit-id",
		},
	}, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Len(t, fields.git.GetFileRevisionCalls(), 1)
	require.Equal(t, "my-commit-id", fields.git.GetFileRevisionCalls()[0].Revision)
	require.Equal(t, "my-service/file1", fields.git.GetFileRevisionCalls()[0].File)
}

func TestResourceManager_GetResource_ProjectResource_PullFails(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.PullFunc = func(gitContext common_models.GitContext) error {
		return errors.New("oops")
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
	})

	require.NotNil(t, err)

	require.Nil(t, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Len(t, fields.git.PullCalls(), 1)

	require.Empty(t, fields.git.GetFileRevisionCalls())
	require.Empty(t, fields.fileSystem.ReadFileCalls())
}

func TestResourceManager_GetResource_ProjectResource_ProjectNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.git.ProjectExistsFunc = func(gitContext common_models.GitContext) bool {
		return false
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "file1",
		GetResourceQuery: models.GetResourceQuery{
			GitCommitID: "my-commit-id",
		},
	})

	require.NotNil(t, err)

	require.Nil(t, result)

	require.Empty(t, fields.git.GetFileRevisionCalls())
}

func TestResourceManager_GetResource_ServiceResource_ServiceNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.stageContext.EstablishFunc = func(params common_models.ConfigurationContextParams) (string, error) {
		return "", errors2.ErrServiceNotFound
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		ResourceURI: "file1",
		GetResourceQuery: models.GetResourceQuery{
			GitCommitID: "my-commit-id",
		},
	})

	require.NotNil(t, err)

	require.Nil(t, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Equal(t, models.Project{ProjectName: "my-project"}, fields.stageContext.EstablishCalls()[0].Params.Project, 1)
	require.Equal(t, &models.Stage{StageName: "my-stage"}, fields.stageContext.EstablishCalls()[0].Params.Stage, 1)
	require.Equal(t, &models.Service{ServiceName: "my-service"}, fields.stageContext.EstablishCalls()[0].Params.Service, 1)
	require.Empty(t, fields.git.GetFileRevisionCalls())
}

func TestResourceManager_GetResource_ProjectResource_ResourceNotFound(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.ReadFileFunc = func(filename string) ([]byte, error) {
		return nil, errors2.ErrResourceNotFound
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		ResourceURI:      "file1",
		GetResourceQuery: models.GetResourceQuery{},
	})

	require.ErrorIs(t, err, errors2.ErrResourceNotFound)

	require.Nil(t, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Empty(t, fields.git.GetFileRevisionCalls())
}

func TestResourceManager_GetResource_ProjectResource_CannotReadFIle(t *testing.T) {
	fields := getTestResourceManagerFields()

	fields.fileSystem.ReadFileFunc = func(filename string) ([]byte, error) {
		return nil, errors.New("oops")
	}
	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
			Stage:   &models.Stage{StageName: "my-stage"},
			Service: &models.Service{ServiceName: "my-service"},
		},
		ResourceURI:      "file1",
		GetResourceQuery: models.GetResourceQuery{},
	})

	require.NotNil(t, err)

	require.Nil(t, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)
	require.Empty(t, fields.git.GetFileRevisionCalls())
}

func TestResourceManager_GetResource_ProjectResource_InvalidResourceName(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResource(models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		ResourceURI: "fi%le1",
		GetResourceQuery: models.GetResourceQuery{
			GitCommitID: "my-commit-id",
		},
	})

	require.ErrorIs(t, err, errors2.ErrResourceInvalidResourceURI)

	require.Nil(t, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.GetFileRevisionCalls())
}

func TestResourceManager_GetResources(t *testing.T) {
	fields := getTestResourceManagerFields()

	rm := NewResourceManager(fields.git, fields.credentialReader, fields.fileSystem, fields.stageContext)

	result, err := rm.GetResources(models.GetResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: "my-project"},
		},
		GetResourcesQuery: models.GetResourcesQuery{
			PageSize: 10,
		},
	})

	require.Nil(t, err)

	require.Equal(t, &models.GetResourcesResponse{
		NextPageKey: "0",
		PageSize:    0,
		Resources: []models.GetResourceResponse{
			{
				Resource: models.Resource{
					ResourceContent: "",
					ResourceURI:     "/file1",
				},
				Metadata: models.Version{
					Branch:      "",
					UpstreamURL: "",
					Version:     "",
				},
			},
			{
				Resource: models.Resource{
					ResourceContent: "",
					ResourceURI:     "/file2",
				},
				Metadata: models.Version{
					Branch:      "",
					UpstreamURL: "",
					Version:     "",
				},
			},
			{
				Resource: models.Resource{
					ResourceContent: "",
					ResourceURI:     "/file3",
				},
				Metadata: models.Version{
					Branch:      "",
					UpstreamURL: "",
					Version:     "",
				},
			},
		},
		TotalCount: 3,
	}, result)

	require.Len(t, fields.stageContext.EstablishCalls(), 1)

	require.Empty(t, fields.git.GetFileRevisionCalls())

	require.Len(t, fields.fileSystem.WalkPathCalls(), 1)
}

type fakeFileInfo struct {
	name  string
	isDir bool
}

func newFakeFileInfo(name string, isDir bool) *fakeFileInfo {
	return &fakeFileInfo{name: name, isDir: isDir}
}

func (f fakeFileInfo) Name() string {
	return f.name
}

func (fakeFileInfo) Size() int64 {
	return 100
}

func (fakeFileInfo) Mode() fs.FileMode {
	return os.ModePerm
}

func (fakeFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fakeFileInfo) IsDir() bool {
	return f.isDir
}

func (fakeFileInfo) Sys() interface{} {
	return nil
}

func getTestResourceManagerFields() testResourceManagerFields {
	return testResourceManagerFields{
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
		credentialReader: &common_mock.CredentialReaderMock{
			GetCredentialsFunc: func(project string) (*common_models.GitCredentials, error) {
				return &common_models.GitCredentials{
					User:      "user",
					Token:     "token",
					RemoteURI: "remote-url",
				}, nil
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
		stageContext: &handler_mock.IConfigurationContextMock{
			EstablishFunc: func(params common_models.ConfigurationContextParams) (string, error) {
				return testConfigDir, nil
			},
		},
	}
}
