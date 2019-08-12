package handlers

import (
	"io/ioutil"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"
)

// GetProjectProjectNameResourceHandlerFunc get list of project resources
func GetProjectProjectNameResourceHandlerFunc(params project_resource.GetProjectProjectNameResourceParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewGetProjectProjectNameResourceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project does not exist")})
	}

	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewGetProjectProjectNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	result := common.GetPaginatedResources(projectConfigPath, params.PageSize, params.NextPageKey)
	return project_resource.NewGetProjectProjectNameResourceOK().WithPayload(result)
}

// PutProjectProjectNameResourceHandlerFunc update list of project resources
func PutProjectProjectNameResourceHandlerFunc(params project_resource.PutProjectProjectNameResourceParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	// utils.Debug("", "Updatingresource(s) in: "+projectConfigPath)
	// utils.Debug("", "Checking out master branch")
	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		// filePath := projectConfigPath + "/" + *res.ResourceURI
		// utils.Debug("", "Updating resource: "+filePath)
		common.WriteFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resources")
	if err != nil {
		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully updated resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return project_resource.NewPutProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return project_resource.NewPutProjectProjectNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PostProjectProjectNameResourceHandlerFunc creates a list of new resources
func PostProjectProjectNameResourceHandlerFunc(params project_resource.PostProjectProjectNameResourceParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	// utils.Debug("", "Creating new resource(s) in: "+projectConfigPath)
	// utils.Debug("", "Checking out master branch")
	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		filePath := projectConfigPath + "/" + *res.ResourceURI
		if !common.FileExists(filePath) {
			// utils.Debug("", "Adding resource: "+filePath)
			common.WriteFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
		}
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully added resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return project_resource.NewPostProjectProjectNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// GetProjectProjectNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.GetProjectProjectNameResourceResourceURIParams) middleware.Responder {
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewGetProjectProjectNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	// utils.Debug("", "Checking out master branch")
	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewGetProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if !common.FileExists(resourcePath) {
		return project_resource.NewGetProjectProjectNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project resource not found")})
	}

	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return project_resource.NewGetProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	resourceContent := strfmt.Base64(dat)
	return project_resource.NewGetProjectProjectNameResourceResourceURIOK().WithPayload(
		&models.Resource{
			ResourceURI:     &params.ResourceURI,
			ResourceContent: resourceContent,
		})
}

// PutProjectProjectNameResourceResourceURIHandlerFunc updates a resource
func PutProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.PutProjectProjectNameResourceResourceURIParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	// utils.Debug("", "Creating new resource(s) in: "+projectConfigPath)
	// utils.Debug("", "Checking out branch: master")
	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	filePath := projectConfigPath + "/" + params.ResourceURI
	common.WriteFile(filePath, params.Resource.ResourceContent)

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI)
	if err != nil {
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully updated resource: "+params.ResourceURI)

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return project_resource.NewPutProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return project_resource.NewPutProjectProjectNameResourceResourceURICreated().WithPayload(&models.Version{
		Version: newVersion,
	})

}

// DeleteProjectProjectNameResourceResourceURIHandlerFunc deletes a project resource
func DeleteProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	resourcePath := projectConfigPath + "/" + params.ResourceURI
	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}
	err = common.DeleteFile(resourcePath)
	if err != nil {
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Deleted resources")
	if err != nil {
		return project_resource.NewDeleteProjectProjectNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully deleted resources")

	return project_resource.NewDeleteProjectProjectNameResourceResourceURINoContent()
}
