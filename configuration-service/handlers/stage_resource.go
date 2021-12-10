package handlers

import (
	"encoding/base64"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"path/filepath"

	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage_resource"
)

// GetProjectProjectNameStageStageNameResourceHandlerFunc get list of stage resources
func GetProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.GetProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String(common.StageDoesNotExistErrorMsg)})
	}

	logger.Debug("Checking out " + params.StageName + " branch")
	err := common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not check out %s branch for project %s", params.StageName, params.ProjectName)
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not retrieve stage resources")})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	result := common.GetPaginatedResources(projectConfigPath, params.PageSize, params.NextPageKey)
	return stage_resource.NewGetProjectProjectNameStageStageNameResourceOK().WithPayload(result)
}

// GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc get the specified resource
func GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.GetProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not unescape resource name")})
	}

	stageConfigPath := common.GetStageConfigPath(params.ProjectName, params.StageName)
	resourcePath := stageConfigPath + "/" + unescapedResourceName

	logger.Debug("Checking out " + params.StageName + " branch")
	err = common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not check out %s branch for project %s", params.StageName, params.ProjectName)
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch containing stage config")})
	}

	if !common.FileExists(resourcePath) {
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage resource not found")})
	}

	resourcePath = filepath.Clean(resourcePath)
	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		logger.WithError(err).Error("could not read file")
		return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read file")})
	}

	resourceContent := base64.StdEncoding.EncodeToString(dat)

	resource := &models.Resource{
		ResourceURI:     &params.ResourceURI,
		ResourceContent: resourceContent,
	}

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	resource.Metadata = metadata

	return stage_resource.NewGetProjectProjectNameStageStageNameResourceResourceURIOK().WithPayload(resource)
}

// PostProjectProjectNameStageStageNameResourceHandlerFunc creates list of new resources in a stage
func PostProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.PostProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewPostProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(common.StageDoesNotExistErrorMsg)})
	}

	stageConfigPath := common.GetStageConfigPath(params.ProjectName, params.StageName)

	logger.Debug("Creating new resource(s) in: " + stageConfigPath + " in stage " + params.StageName)

	for _, res := range params.Resources.Resources {
		filePath := stageConfigPath + "/" + *res.ResourceURI
		logger.Debug("Adding resource: " + filePath)
		common.WriteBase64EncodedFile(filePath, res.ResourceContent)
	}

	logger.Debug("Staging Changes")
	commitId, err := common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch for project %s", params.StageName, params.ProjectName)
		return stage_resource.NewPostProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully added resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	if commitId != "" {
		metadata.Version = commitId
	}
	return stage_resource.NewPostProjectProjectNameStageStageNameResourceCreated().WithPayload(metadata)
}

// PutProjectProjectNameStageStageNameResourceHandlerFunc updates list of stage resources
func PutProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.PutProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(common.StageDoesNotExistErrorMsg)})
	}

	stageConfigPath := common.GetStageConfigPath(params.ProjectName, params.StageName)

	logger.Debug("Creating new resource(s) in: " + stageConfigPath + " in stage " + params.StageName)

	for _, res := range params.Resources.Resources {
		filePath := stageConfigPath + "/" + *res.ResourceURI
		common.WriteBase64EncodedFile(filePath, res.ResourceContent)
	}

	logger.Debug("Staging Changes")
	commitId, err := common.StageAndCommitAll(params.ProjectName, "Updated resources")
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch for project %s", params.StageName, params.ProjectName)
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	if commitId != "" {
		metadata.Version = commitId
	}
	return stage_resource.NewPutProjectProjectNameStageStageNameResourceCreated().WithPayload(metadata)
}

// PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc updates the specified stage resource
func PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.PutProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(common.StageDoesNotExistErrorMsg)})
	}

	stageConfigPath := common.GetStageConfigPath(params.ProjectName, params.StageName)

	logger.Debug("Creating new resource(s) in: " + stageConfigPath + " in stage " + params.StageName)

	filePath := stageConfigPath + "/" + params.ResourceURI
	common.WriteBase64EncodedFile(filePath, params.Resource.ResourceContent)

	logger.Debug("Staging Changes")
	commitId, err := common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch for project %s", params.StageName, params.ProjectName)
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resource: " + params.ResourceURI)

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	if commitId != "" {
		metadata.Version = commitId
	}
	return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURICreated().WithPayload(metadata)
}

// DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc deletes the specified stage resource
func DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage_resource.NewDeleteProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(common.StageDoesNotExistErrorMsg)})
	}

	stageConfigPath := common.GetStageConfigPath(params.ProjectName, params.StageName)

	logger.Debug("Creating new resource(s) in: " + stageConfigPath + " in stage " + params.StageName)

	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not unescape resource name")})
	}

	filePath := stageConfigPath + "/" + unescapedResourceName

	err = common.DeleteFile(filePath)
	if err != nil {
		logger.Error(err.Error())
		return stage_resource.NewDeleteProjectProjectNameStageStageNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete file")})
	}

	logger.Debug("Staging Changes")
	commitId, err := common.StageAndCommitAll(params.ProjectName, "Updated resource: "+unescapedResourceName)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch for project %s", params.StageName, params.ProjectName)
		return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debugf("Successfully updated resource: %s" + unescapedResourceName)

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	if commitId != "" {
		metadata.Version = commitId
	}
	return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURICreated().WithPayload(metadata)
}
