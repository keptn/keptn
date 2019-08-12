package handlers

import (
	"io/ioutil"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_resource"
)

// GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc get list of resources for the service
func GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResource has not yet been implemented")
}

// GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName
	resourcePath := serviceConfigPath + "/" + params.ResourceURI
	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
	}
	// utils.Debug("", "Checking out "+params.StageName+" branch")
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if !common.FileExists(resourcePath) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project resource not found")})
	}

	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	resourceContent := strfmt.Base64(dat)
	return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIOK().WithPayload(
		&models.Resource{
			ResourceURI:     &params.ResourceURI,
			ResourceContent: resourceContent,
		})
}

// DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc deletes the specified resource
func DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(params service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
}

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
	}
	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

	// utils.Debug("", "Creating new resource(s) in: "+serviceConfigPath+" in stage "+params.StageName)
	// utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		filePath := serviceConfigPath + "/" + *res.ResourceURI
		// don't overwrite existing files
		if !common.FileExists(filePath) {
			// utils.Debug("", "Adding resource: "+filePath)
			common.WriteFile(filePath, res.ResourceContent)
		}
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully added resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc updates a list of resources
func PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
	}
	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

	// utils.Debug("", "Updating resource(s) in: "+serviceConfigPath+" in stage "+params.StageName)
	// utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		filePath := serviceConfigPath + "/" + *res.ResourceURI
		// utils.Debug("", "Updating resource: "+filePath)
		common.WriteFile(filePath, res.ResourceContent)
	}

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resources")
	if err != nil {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully updated resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc updates a specified resource
func PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
	}
	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

	// utils.Debug("", "updating resource(s) in: "+serviceConfigPath+" in stage "+params.StageName)
	// utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	filePath := serviceConfigPath + "/" + params.ResourceURI
	common.WriteFile(filePath, params.Resource.ResourceContent)

	// utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI)
	if err != nil {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	// utils.Debug("", "Successfully updated resource: "+params.ResourceURI)

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}
