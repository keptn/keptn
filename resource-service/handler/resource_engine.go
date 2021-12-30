package handler

import (
	"encoding/base64"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/sirupsen/logrus"
)

type ResourceEngine struct {
	git        common.IGit
	fileSystem common.IFileSystem
}

func NewResourceEngine(git common.IGit, fileSystem common.IFileSystem) *ResourceEngine {
	return &ResourceEngine{
		git:        git,
		fileSystem: fileSystem,
	}
}

func (p ResourceEngine) readResource(gitContext *common.GitContext, params models.GetResourceParams, resourcePath string) (*models.GetResourceResponse, error) {
	var fileContent []byte
	var err error

	if params.GitCommitID != "" {
		fileContent, err = p.git.GetFileRevision(*gitContext, resourcePath, params.GitCommitID, params.ResourceURI)
	} else {
		fileContent, err = p.fileSystem.ReadFile(resourcePath)
	}
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

func (p ResourceEngine) writeResource(gitContext *common.GitContext, resourceContent, resourcePath string) (*models.WriteResourceResponse, error) {
	if err := p.fileSystem.WriteBase64EncodedFile(resourcePath, resourceContent); err != nil {
		return nil, err
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Updated resource")
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
}

func (p ResourceEngine) writeResources(gitContext *common.GitContext, resources []models.Resource, directory string) (*models.WriteResourceResponse, error) {
	for _, res := range resources {
		filePath := directory + "/" + res.ResourceURI
		logrus.Debug("Adding resource: " + filePath)
		if err := p.fileSystem.WriteBase64EncodedFile(directory+"/"+res.ResourceURI, string(res.ResourceContent)); err != nil {
			return nil, err
		}
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Added resources")
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
}

func (p ResourceEngine) deleteResource(gitContext *common.GitContext, resourcePath string) (*models.WriteResourceResponse, error) {
	if err := p.fileSystem.DeleteFile(resourcePath); err != nil {
		return nil, err
	}

	commitID, err := p.git.StageAndCommitAll(*gitContext, "Deleted resource")
	if err != nil {
		return nil, err
	}

	return &models.WriteResourceResponse{CommitID: commitID}, nil
}
