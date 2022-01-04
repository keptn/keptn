package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"net/url"
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
	git              common.IGit
	credentialReader common.CredentialReader
	fileSystem       common.IFileSystem
}

func NewResourceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem) *ResourceManager {
	projectResourceManager := &ResourceManager{
		git:              git,
		credentialReader: credentialReader,
		fileSystem:       fileWriter,
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

	_, configPath, err := p.establishContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return nil, err
	}

	result, err := GetPaginatedResources(configPath, params.PageSize, params.NextPageKey, p.fileSystem)
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
		return nil, errors.ErrResourceInvalidResourceURI
	}

	resourcePath := configPath + "/" + unescapedResourceName

	return p.readResource(gitContext, params, resourcePath)
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

	resourcePath := configPath + "/" + params.ResourceURI

	return p.deleteResource(gitContext, resourcePath)
}

func (p ResourceManager) establishContext(project models.Project, stage *models.Stage, service *models.Service) (*common.GitContext, string, error) {
	credentials, err := p.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return nil, "", fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !p.git.ProjectExists(gitContext) {
		return nil, "", errors.ErrProjectNotFound
	}

	var branch string

	if stage == nil {
		branch, err = p.git.GetDefaultBranch(gitContext)
		if err != nil {
			return nil, "", fmt.Errorf("could not determine default branch of project %s: %w", project.ProjectName, err)
		}
	} else {
		branch = stage.StageName
	}

	if err := p.git.CheckoutBranch(gitContext, branch); err != nil {
		return nil, "", fmt.Errorf("could not check out branch %s of project %s: %w", branch, project.ProjectName, err)
	}

	var configPath string
	if service == nil {
		configPath = common.GetProjectConfigPath(project.ProjectName)
	} else {
		configPath = common.GetServiceConfigPath(project.ProjectName, service.ServiceName)
	}
	return &gitContext, configPath, nil
}

func (p ResourceManager) readResource(gitContext *common.GitContext, params models.GetResourceParams, resourcePath string) (*models.GetResourceResponse, error) {
	var fileContent []byte
	var revision string
	var err error

	if params.GitCommitID != "" {
		fileContent, err = p.git.GetFileRevision(*gitContext, resourcePath, params.GitCommitID, params.ResourceURI)
		revision = params.GitCommitID
	} else {
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

func (p ResourceManager) writeAndCommitResource(gitContext *common.GitContext, resourcePath, resourceContent string) (*models.WriteResourceResponse, error) {
	if err := p.storeResource(resourcePath, resourceContent); err != nil {
		return nil, err
	}

	return p.stageAndCommit(gitContext, "Updated resource")
}

func (p ResourceManager) writeAndCommitResources(gitContext *common.GitContext, resources []models.Resource, directory string) (*models.WriteResourceResponse, error) {
	for _, res := range resources {
		filePath := directory + "/" + res.ResourceURI
		if err := p.storeResource(filePath, string(res.ResourceContent)); err != nil {
			return nil, err
		}
	}

	return p.stageAndCommit(gitContext, "Added resources")
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

func (p ResourceManager) stageAndCommit(gitContext *common.GitContext, message string) (*models.WriteResourceResponse, error) {
	commitID, err := p.git.StageAndCommitAll(*gitContext, message)
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
}

func (p ResourceManager) deleteResource(gitContext *common.GitContext, resourcePath string) (*models.WriteResourceResponse, error) {
	if err := p.fileSystem.DeleteFile(resourcePath); err != nil {
		return nil, err
	}

	return p.stageAndCommit(gitContext, "Deleted resources")
}
