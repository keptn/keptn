package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"
)

// GetProjectProjectNameResourceHandlerFunc get list of project resources
func GetProjectProjectNameResourceHandlerFunc(params project_resource.GetProjectProjectNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation project_resource.GetProjectProjectNameResource has not yet been implemented")
}

// PutProjectProjectNameResourceHandlerFunc update list of project resources
func PutProjectProjectNameResourceHandlerFunc(params project_resource.PutProjectProjectNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation project_resource.PutProjectProjectNameResource has not yet been implemented")
}

// PostProjectProjectNameResourceHandlerFunc creates a list of new resources
func PostProjectProjectNameResourceHandlerFunc(params project_resource.PostProjectProjectNameResourceParams) middleware.Responder {
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	utils.Debug("", "Creating new resource(s) in: "+projectConfigPath)
	utils.Debug("", "Checking out master branch: ")
	err := common.CheckoutBranch(params.ProjectName, "master")
	if err != nil {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		common.WriteFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
	}

	err = common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		fmt.Print(err.Error())
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return project_resource.NewPostProjectProjectNameResourceCreated()
}

// GetProjectProjectNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.GetProjectProjectNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation project_resource.GetProjectProjectNameResourceResourceURI has not yet been implemented")
}

// PutProjectProjectNameResourceResourceURIHandlerFunc updates a resource
func PutProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.PutProjectProjectNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation project_resource.PutProjectProjectNameResourceResourceURI has not yet been implemented")
}

// DeleteProjectProjectNameResourceResourceURIHandlerFunc deletes a project resource
func DeleteProjectProjectNameResourceResourceURIHandlerFunc(params project_resource.DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation project_resource.DeleteProjectProjectNameResourceResourceURI has not yet been implemented")
}
