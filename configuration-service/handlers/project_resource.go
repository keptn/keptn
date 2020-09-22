package handlers

import (
	"encoding/base64"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"io/ioutil"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"
)

// GetProjectProjectNameResourceHandlerFunc get list of project resources
func GetProjectProjectNameResourceHandlerFunc(params project_resource.GetProjectProjectNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewGetProjectProjectNameResourceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	err := common.CheckoutBranch(params.ProjectName, "master", *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewGetProjectProjectNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	result := common.GetPaginatedResources(projectConfigPath, params.PageSize, params.NextPageKey)

	return project_resource.NewGetProjectProjectNameResourceOK().WithPayload(result)
}

// PutProjectProjectNameResourceHandlerFunc update list of project resources
func PutProjectProjectNameResourceHandlerFunc(params project_resource.PutProjectProjectNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	logger.Debug("Updating resource(s) in: " + projectConfigPath)
	logger.Debug("Checking out master branch")

	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	for _, res := range params.Resources.Resources {
		filePath := projectConfigPath + "/" + *res.ResourceURI
		logger.Debug("Updating resource: " + filePath)
		common.WriteBase64EncodedFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
		if strings.ToLower(*res.ResourceURI) == "shipyard.yaml" {
			mv := common.GetProjectsMaterializedView()
			logger.Debug("updating shipyard.yaml content for project " + params.ProjectName + " in mongoDB table")
			err := mv.UpdateShipyard(params.ProjectName, res.ResourceContent)
			if err != nil {
				logger.Error("Could not update shipyard.yaml content for project " + params.ProjectName + ": " + err.Error())
				return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
			}
		}
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resources", true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = "master"

	return project_resource.NewPutProjectProjectNameResourceCreated().WithPayload(metadata)
}

// PostProjectProjectNameResourceHandlerFunc creates a list of new resources
func PostProjectProjectNameResourceHandlerFunc(params project_resource.PostProjectProjectNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	logger.Debug("Creating new resource(s) in: " + projectConfigPath)
	logger.Debug("Checking out master branch")

	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	for _, res := range params.Resources.Resources {
		filePath := projectConfigPath + "/" + *res.ResourceURI
		logger.Debug("Adding resource: " + filePath)
		common.WriteBase64EncodedFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources", true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to master branch of project %s", params.ProjectName))
		logger.Error(err.Error())

		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully added resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = "master"

	return project_resource.NewPostProjectProjectNameResourceCreated().WithPayload(metadata)
}

// GetProjectProjectNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.GetProjectProjectNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewGetProjectProjectNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	logger.Debug("Checking out master branch")
	err := common.CheckoutBranch(params.ProjectName, "master", *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(fmt.Sprintf("Could check out master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewGetProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI
	if !common.FileExists(resourcePath) {
		return project_resource.NewGetProjectProjectNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project resource not found")})
	}

	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		logger.Error(err.Error())
		return project_resource.NewGetProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read file")})
	}

	resourceContent := base64.StdEncoding.EncodeToString(dat)

	resource := &models.Resource{
		ResourceURI:     &params.ResourceURI,
		ResourceContent: resourceContent,
	}

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = "master"
	resource.Metadata = metadata

	return project_resource.NewGetProjectProjectNameResourceResourceURIOK().WithPayload(resource)
}

// PutProjectProjectNameResourceResourceURIHandlerFunc updates a resource
func PutProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.PutProjectProjectNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	logger.Debug("Creating new resource(s) in: " + projectConfigPath)
	logger.Debug("Checking out branch: master")

	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	filePath := projectConfigPath + "/" + params.ResourceURI
	common.WriteBase64EncodedFile(filePath, params.Resource.ResourceContent)

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI, true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resource: " + params.ResourceURI)

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = "master"

	return project_resource.NewPutProjectProjectNameResourceResourceURICreated().WithPayload(metadata)

}

// DeleteProjectProjectNameResourceResourceURIHandlerFunc deletes a project resource
func DeleteProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI

	err = common.DeleteFile(resourcePath)
	if err != nil {
		logger.Error(err.Error())
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete file")})
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Deleted resources", true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to master branch of project %s", params.ProjectName))
		logger.Error(err.Error())
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully deleted resources")

	metadata := common.GetResourceMetadata(params.ProjectName)

	return project_resource.NewDeleteProjectProjectNameResourceResourceURINoContent().WithPayload(metadata)
}
