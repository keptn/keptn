package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
)

type ProjectResourceManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
	fileWriter       common.IFileWriter
}

func NewProjectResourceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileWriter) *ProjectResourceManager {
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
	panic("implement me")
}

func (p ProjectResourceManager) UpdateResources(params models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
	panic("implement me")
}

func (p ProjectResourceManager) GetResource(params models.GetResourceParams) (*models.GetResourceResponse, error) {
	panic("implement me")
}

func (p ProjectResourceManager) UpdateResource(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
	panic("implement me")
}

func (p ProjectResourceManager) DeleteResource(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
	panic("implement me")
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
		return nil, common.ErrProjectNotFound
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
