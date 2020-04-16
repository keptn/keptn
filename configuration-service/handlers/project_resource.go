package handlers

import (
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	utils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"
)

// GetProjectProjectNameResourceHandlerFunc get list of project resources
func GetProjectProjectNameResourceHandlerFunc(params project_resource.GetProjectProjectNameResourceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {

		return project_resource.NewGetProjectProjectNameResourceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	err := common.CheckoutBranch(params.ProjectName, "master", *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewGetProjectProjectNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	result := common.GetPaginatedResources(projectConfigPath, params.PageSize, params.NextPageKey)

	return project_resource.NewGetProjectProjectNameResourceOK().WithPayload(result)
}

// PutProjectProjectNameResourceHandlerFunc update list of project resources
func PutProjectProjectNameResourceHandlerFunc(params project_resource.PutProjectProjectNameResourceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	if !common.ProjectExists(params.ProjectName) {

		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	logger.Debug("Updating resource(s) in: " + projectConfigPath)
	logger.Debug("Checking out master branch")
	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
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
			}
		}
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resources")
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Coult not retrieve latest version")})
	}

	return project_resource.NewPutProjectProjectNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PostProjectProjectNameResourceHandlerFunc creates a list of new resources
func PostProjectProjectNameResourceHandlerFunc(params project_resource.PostProjectProjectNameResourceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {

		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	logger.Debug("Creating new resource(s) in: " + projectConfigPath)
	logger.Debug("Checking out master branch")
	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	for _, res := range params.Resources.Resources {
		filePath := projectConfigPath + "/" + *res.ResourceURI
		logger.Debug("Adding resource: " + filePath)
		common.WriteBase64EncodedFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully added resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not retrieve latest version")})
	}

	return project_resource.NewPostProjectProjectNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// GetProjectProjectNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.GetProjectProjectNameResourceResourceURIParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI
	if !common.ProjectExists(params.ProjectName) {

		return project_resource.NewGetProjectProjectNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	logger.Debug("Checking out master branch")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	err := common.CheckoutBranch(params.ProjectName, "master", *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewGetProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	if !common.FileExists(resourcePath) {

		return project_resource.NewGetProjectProjectNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project resource not found")})
	}

	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewGetProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read file")})
	}

	resourceContent := base64.StdEncoding.EncodeToString(dat)

	return project_resource.NewGetProjectProjectNameResourceResourceURIOK().WithPayload(
		&models.Resource{
			ResourceURI:     &params.ResourceURI,
			ResourceContent: resourceContent,
		})
}

// PutProjectProjectNameResourceResourceURIHandlerFunc updates a resource
func PutProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.PutProjectProjectNameResourceResourceURIParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {

		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	logger.Debug("Creating new resource(s) in: " + projectConfigPath)
	logger.Debug("Checking out branch: master")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	filePath := projectConfigPath + "/" + params.ResourceURI
	common.WriteBase64EncodedFile(filePath, params.Resource.ResourceContent)

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resource: " + params.ResourceURI)

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not retrieve latest version")})
	}

	return project_resource.NewPutProjectProjectNameResourceResourceURICreated().WithPayload(&models.Version{
		Version: newVersion,
	})

}

// DeleteProjectProjectNameResourceResourceURIHandlerFunc deletes a project resource
func DeleteProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {

		return project_resource.NewDeleteProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	err := common.CheckoutBranch(params.ProjectName, "master", false)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	err = common.DeleteFile(resourcePath)
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete file")})
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Deleted resources")
	if err != nil {
		logger.Error(err.Error())

		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully deleted resources")

	return project_resource.NewDeleteProjectProjectNameResourceResourceURINoContent()
}
