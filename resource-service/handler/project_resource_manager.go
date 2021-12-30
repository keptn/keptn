package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
	"net/url"
)

type ProjectResourceManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
	fileWriter       common.IFileSystem
}

func NewProjectResourceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem) *ProjectResourceManager {
	projectResourceManager := &ProjectResourceManager{
		git:              git,
		credentialReader: credentialReader,
		fileWriter:       fileWriter,
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

	for _, res := range params.Resources {
		filePath := projectConfigPath + "/" + res.ResourceURI
		logger.Debug("Adding resource: " + filePath)
		if err := p.fileWriter.WriteBase64EncodedFile(projectConfigPath+"/"+res.ResourceURI, string(res.ResourceContent)); err != nil {
			return nil, err
		}
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Added resources")
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
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

	for _, res := range params.Resources {
		filePath := projectConfigPath + "/" + res.ResourceURI
		err := p.fileWriter.WriteBase64EncodedFile(filePath, string(res.ResourceContent))
		if err != nil {
			return nil, err
		}
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Added resources")
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
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

	fileContent, err := p.fileWriter.ReadFile(resourcePath)
	if err != nil {
		return nil, err
	}

	resourceContent := base64.StdEncoding.EncodeToString(fileContent)

	currentRevision, err := p.git.GetCurrentRevision(*gitContext)
	if err != nil {
		return nil, err
	}

	return &models.GetResourceResponse{
		Resource: models.Resource{
			ResourceURI:     params.ResourceURI,
			ResourceContent: models.ResourceContent(resourceContent),
		},
		Metadata: models.Version{
			UpstreamURL: gitContext.Credentials.RemoteURI,
			Version:     currentRevision,
		},
	}, nil
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

	if err := p.fileWriter.WriteBase64EncodedFile(resourcePath, string(params.ResourceContent)); err != nil {
		return nil, err
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Updated resource "+params.ResourceURI)
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
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

	if err := p.fileWriter.DeleteFile(resourcePath); err != nil {
		return nil, err
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Deleted resource "+params.ResourceURI)
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
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
