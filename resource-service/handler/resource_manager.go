package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
)

//IResourceManager provides an interface for resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/resource_manager_mock.go . IResourceManager
type IResourceManager interface {
	CreateResources(params models.CreateResourcesParams) (*models.WriteResourceResponse, error)
	GetResources(params models.GetResourcesParams) (*models.GetResourcesResponse, error)
	UpdateResources(params models.UpdateResourcesParams) (*models.WriteResourceResponse, error)
	GetResource(params models.GetResourceParams) (*models.GetResourceResponse, error)
	UpdateResource(params models.UpdateResourceParams) (*models.WriteResourceResponse, error)
	DeleteResource(params models.DeleteResourceParams) (*models.WriteResourceResponse, error)
}

type ResourceManager struct {
	git                  common.IGit
	credentialReader     common.CredentialReader
	fileSystem           common.IFileSystem
	configurationContext IConfigurationContext
}

func NewResourceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem, stageContext IConfigurationContext) *ResourceManager {
	projectResourceManager := &ResourceManager{
		git:                  git,
		credentialReader:     credentialReader,
		fileSystem:           fileWriter,
		configurationContext: stageContext,
	}
	return projectResourceManager
}

func (p ResourceManager) CreateResources(params models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}

	return p.writeAndCommitResources(gitContext, params.Resources, configPath)
}

func (p ResourceManager) GetResources(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}
	revision, err := p.git.GetCurrentRevision(*gitContext)
	if err != nil {
		return nil, err
	}
	metadata := models.Version{
		UpstreamURL: gitContext.Credentials.RemoteURI,
		Version:     revision,
	}

	result, err := GetPaginatedResources(configPath, params.PageSize, params.NextPageKey, p.fileSystem, metadata)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p ResourceManager) UpdateResources(params models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}

	return p.writeAndCommitResources(gitContext, params.Resources, configPath)
}

func (p ResourceManager) GetResource(params models.GetResourceParams) (*models.GetResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}

	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return nil, kerrors.ErrResourceInvalidResourceURI
	}

	return p.readResource(gitContext, params, configPath, unescapedResourceName)
}

func (p ResourceManager) UpdateResource(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}

	resourcePath := configPath + "/" + params.ResourceURI

	return p.writeAndCommitResource(gitContext, resourcePath, string(params.ResourceContent))
}

func (p ResourceManager) DeleteResource(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}

	unescapedResource, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return nil, err
	}

	resourcePath := configPath + "/" + unescapedResource

	var resultErr error
	var resultCommit *models.WriteResourceResponse
	_ = retry.Retry(func() error {
		err = p.git.Pull(*gitContext)
		if err != nil {
			resultErr = err
			return nil
		}
		response, err := p.deleteResource(gitContext, resourcePath)
		if err != nil {
			if errors.Is(err, kerrors.ErrNonFastForwardUpdate) || errors.Is(err, kerrors.ErrForceNeeded) {
				return err
			}
			resultErr = err
			return nil
		}
		resultCommit = response
		resultErr = err
		return nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))
	return resultCommit, resultErr
}

func (p ResourceManager) establishContext(project models.Project, stage *models.Stage, service *models.Service) (*common_models.GitContext, string, error) {
	credentials, err := p.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return nil, "", fmt.Errorf(kerrors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !p.git.ProjectExists(gitContext) {
		return nil, "", kerrors.ErrProjectNotFound
	}

	configPath, err := p.configurationContext.Establish(common_models.ConfigurationContextParams{
		Project:                 project,
		Stage:                   stage,
		Service:                 service,
		GitContext:              gitContext,
		CheckConfigDirAvailable: true,
	})
	if err != nil {
		return nil, "", err
	}
	return &gitContext, configPath, nil
}

func (p ResourceManager) readResource(gitContext *common_models.GitContext, params models.GetResourceParams, configPath string, resourceName string) (*models.GetResourceResponse, error) {
	var fileContent []byte
	var revision string
	var err error

	if params.GitCommitID != "" && params.GitCommitID != "\"\"" {
		// if commit ID is set, path needs to be relative to the project directory
		configPath = strings.TrimPrefix(configPath, common.GetProjectConfigPath(params.ProjectName))
		// resource path must not start with "/", otherwise git is not able to resolve the revision
		resourcePath := strings.TrimPrefix(configPath+"/"+resourceName, "/")
		fileContent, err = p.git.GetFileRevision(*gitContext, params.GitCommitID, resourcePath)
		revision = params.GitCommitID
	} else {
		resourcePath := configPath + "/" + resourceName
		if err := p.git.Pull(*gitContext); err != nil {
			return nil, err
		}
		fileContent, err = p.fileSystem.ReadFile(resourcePath)
		if err != nil {
			return nil, err
		}
		revision, err = p.git.GetCurrentRevision(*gitContext)
	}
	if err != nil {
		return nil, err
	}

	resourceContent := base64.StdEncoding.EncodeToString(fileContent)

	return &models.GetResourceResponse{
		Resource: models.Resource{
			ResourceURI:     params.ResourceURI,
			ResourceContent: models.ResourceContent(resourceContent),
		},
		Metadata: models.Version{
			UpstreamURL: gitContext.Credentials.RemoteURI,
			Version:     revision,
		},
	}, nil
}

func (p ResourceManager) writeAndCommitResource(gitContext *common_models.GitContext, resourcePath, resourceContent string) (*models.WriteResourceResponse, error) {

	var resultErr error
	var resultCommit *models.WriteResourceResponse
	_ = retry.Retry(func() error {
		err := p.git.Pull(*gitContext)
		if err != nil {
			resultErr = err
			return nil
		}
		if err := p.storeResource(resourcePath, resourceContent); err != nil {
			resultErr = err
			return nil
		}

		commit, err := p.stageAndCommit(gitContext, "Updated resource")
		if err != nil {
			if errors.Is(err, kerrors.ErrNonFastForwardUpdate) || errors.Is(err, kerrors.ErrForceNeeded) {
				return err
			}
			resultErr = err
			return nil
		}
		resultCommit = commit
		return nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))
	return resultCommit, resultErr
}

func (p ResourceManager) writeAndCommitResources(gitContext *common_models.GitContext, resources []models.Resource, directory string) (*models.WriteResourceResponse, error) {

	var resultErr error
	var resultCommit *models.WriteResourceResponse
	_ = retry.Retry(func() error {
		err := p.git.Pull(*gitContext)
		if err != nil {
			resultErr = err
			return nil
		}
		for _, res := range resources {
			filePath := directory + "/" + res.ResourceURI
			if err := p.storeResource(filePath, string(res.ResourceContent)); err != nil {
				resultErr = err
				return nil
			}
		}

		commit, err := p.stageAndCommit(gitContext, "Updated resource")
		if err != nil {
			if errors.Is(err, kerrors.ErrNonFastForwardUpdate) || errors.Is(err, kerrors.ErrForceNeeded) {
				return err
			}
			resultErr = err
			return nil
		}
		resultCommit = commit
		return nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))
	return resultCommit, resultErr
}

func (p ResourceManager) storeResource(resourcePath, resourceContent string) error {
	if err := p.fileSystem.WriteBase64EncodedFile(resourcePath, resourceContent); err != nil {
		return err
	}
	if common.IsHelmChartPath(resourcePath) {
		if err := p.fileSystem.WriteHelmChart(resourcePath); err != nil {
			return err
		}
	}
	return nil
}

func (p ResourceManager) stageAndCommit(gitContext *common_models.GitContext, message string) (*models.WriteResourceResponse, error) {
	commitID, err := p.git.StageAndCommitAll(*gitContext, message)
	if err != nil {
		return nil, err
	}
	result := &models.WriteResourceResponse{
		CommitID: commitID,
		Metadata: models.Version{
			UpstreamURL: gitContext.Credentials.RemoteURI,
			Version:     commitID,
		},
	}

	return result, nil
}

func (p ResourceManager) deleteResource(gitContext *common_models.GitContext, resourcePath string) (*models.WriteResourceResponse, error) {
	if !p.fileSystem.FileExists(resourcePath) {
		return nil, kerrors.ErrResourceNotFound
	}
	if err := p.fileSystem.DeleteFile(resourcePath); err != nil {
		return nil, err
	}

	return p.stageAndCommit(gitContext, "Deleted resources")
}
