package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage_resource"
)

// GetProjectProjectNameStageStageNameResourceHandlerFunc get list of stage resources
func GetProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.GetProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation stage_resource.GetProjectProjectNameStageStageNameResource has not yet been implemented")
}

// GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc get the specified resource
func GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.GetProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation stage_resource.GetProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
}

// PostProjectProjectNameStageStageNameResourceHandlerFunc creates list of new resources in a stage
func PostProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.PostProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName

	utils.Debug("", "Creating new resource(s) in: "+projectConfigPath+" in stage "+params.StageName)
	utils.Debug("", "Checking out branch: "+params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
		return project_resource.NewPostProjectProjectNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	for _, res := range params.Resources.Resources {
		utils.Debug("", "Adding resource: "+projectConfigPath+"/"+*res.ResourceURI)

		common.WriteFile(projectConfigPath+"/"+*res.ResourceURI, res.ResourceContent)
	}

	utils.Debug("", "Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources")
	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	utils.Debug("", "Successfully added resources")

	newVersion, err := common.GetCurrentVersion(params.ProjectName)
	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	return project_resource.NewPostProjectProjectNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PutProjectProjectNameStageStageNameResourceHandlerFunc updates list of stage resources
func PutProjectProjectNameStageStageNameResourceHandlerFunc(params stage_resource.PutProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation stage_resource.PutProjectProjectNameStageStageNameResource has not yet been implemented")
}

// PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc updates the specified stage resource
func PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.PutProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation stage_resource.PutProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
}

// DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc deletes the specified stage resource
func DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(params stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
}
