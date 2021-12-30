package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"net/url"
)

type ProjectResourceManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
	fileWriter       common.IFileSystem
	resourceEngine   *ResourceEngine
}

func NewProjectResourceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem) *ProjectResourceManager {
	resourceEngine := NewResourceEngine(git, fileWriter)
	projectResourceManager := &ProjectResourceManager{
		git:              git,
		credentialReader: credentialReader,
		fileWriter:       fileWriter,
		resourceEngine:   resourceEngine,
	}
	return projectResourceManager
}

func (p ProjectResourceManager) CreateResources(params models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, err := p.establishProjectContext(params.Project)
	if err != nil {
		return nil, err
	}

	projectConfigPath := common.GetProjectConfigPath(params.ProjectName)

	return p.resourceEngine.writeResources(gitContext, params.Resources, projectConfigPath)
}

func (p ProjectResourceManager) GetResources(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	_, err := p.establishProjectContext(params.Project)
	if err != nil {
		return nil, err
	}

	projectConfigPath := common.GetProjectConfigPath(params.ProjectName)

	result, err := common.GetPaginatedResources(projectConfigPath, params.PageSize, params.NextPageKey, p.fileWriter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p ProjectResourceManager) UpdateResources(params models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, err := p.establishProjectContext(params.Project)
	if err != nil {
		return nil, err
	}

	projectConfigPath := common.GetProjectConfigPath(params.ProjectName)

	return p.resourceEngine.writeResources(gitContext, params.Resources, projectConfigPath)
}

func (p ProjectResourceManager) GetResource(params models.GetResourceParams) (*models.GetResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, err := p.establishProjectContext(params.Project)
	if err != nil {
		return nil, err
	}

	projectConfigPath := common.GetProjectConfigPath(params.ProjectName)

	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return nil, errors.ErrResourceInvalidResourceURI
	}

	resourcePath := projectConfigPath + "/" + unescapedResourceName
	if !p.fileWriter.FileExists(resourcePath) {
		return nil, errors.ErrResourceNotFound
	}

	return p.resourceEngine.readResource(gitContext, params, resourcePath)
}

func (p ProjectResourceManager) UpdateResource(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, err := p.establishProjectContext(params.Project)
	if err != nil {
		return nil, err
	}

	projectConfigPath := common.GetProjectConfigPath(params.ProjectName)

	resourcePath := projectConfigPath + "/" + params.ResourceURI

	return p.resourceEngine.writeResource(gitContext, resourcePath, string(params.ResourceContent))
}

func (p ProjectResourceManager) DeleteResource(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, err := p.establishProjectContext(params.Project)
	if err != nil {
		return nil, err
	}

	projectConfigPath := common.GetProjectConfigPath(params.ProjectName)

	resourcePath := projectConfigPath + "/" + params.ResourceURI

	return p.resourceEngine.deleteResource(gitContext, resourcePath)
}

func (p ProjectResourceManager) establishProjectContext(project models.Project) (*common.GitContext, error) {
	credentials, err := p.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return nil, fmt.Errorf("could not read credentials for project %s: %w", project.ProjectName, err)
	}

	gitContext := common.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !p.git.ProjectExists(gitContext) {
		return nil, errors.ErrProjectNotFound
	}

	defaultBranch, err := p.git.GetDefaultBranch(gitContext)
	if err != nil {
		return nil, fmt.Errorf("could not determine default branch of project %s: %w", project.ProjectName, err)
	}

	if err := p.git.CheckoutBranch(gitContext, defaultBranch); err != nil {
		return nil, fmt.Errorf("could not check out branch %s of project %s: %w", defaultBranch, project.ProjectName, err)
	}

	return &gitContext, nil
}
