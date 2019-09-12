package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
)

// PostProjectProjectNameStageHandlerFunc creates a new stage
func PostProjectProjectNameStageHandlerFunc(params stage.PostProjectProjectNameStageParams) middleware.Responder {
	common.Lock()
	defer common.UnLock()
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist.")})
	}
	err := common.CreateBranch(params.ProjectName, params.Stage.StageName, "master")
	if err != nil {
		logger.Error(err.Error())
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not create stage.")})
	}
	return stage.NewPostProjectProjectNameStageNoContent()
}

// PutProjectProjectNameStageStageNameHandlerFunc updates a stage
func PutProjectProjectNameStageStageNameHandlerFunc(params stage.PutProjectProjectNameStageStageNameParams) middleware.Responder {
	return middleware.NotImplemented("operation stage.PutProjectProjectNameStageStageName has not yet been implemented")
}

// DeleteProjectProjectNameStageStageNameHandlerFunc deletes a stage
func DeleteProjectProjectNameStageStageNameHandlerFunc(params stage.DeleteProjectProjectNameStageStageNameParams) middleware.Responder {
	return middleware.NotImplemented("operation stage.DeleteProjectProjectNameStageStageName has not yet been implemented")
}

// GetProjectProjectNameStageHandlerFunc gets list of stages for a project
func GetProjectProjectNameStageHandlerFunc(params stage.GetProjectProjectNameStageParams) middleware.Responder {
	common.Lock()
	defer common.UnLock()
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewGetProjectProjectNameStageNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project does not exist.")})
	}
	branches, err := common.GetBranches(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())
		return stage.NewGetProjectProjectNameStageDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not retrieve stages.")})
	}

	result := []*models.Stage{}

	for _, branch := range branches {
		if branch != "master" && branch != "" {
			stage := &models.Stage{
				StageName: branch,
			}
			result = append(result, stage)
		}
	}
	return stage.NewGetProjectProjectNameStageOK().WithPayload(&models.Stages{
		NextPageKey: "",
		PageSize:    float64(len(result)),
		TotalCount:  float64(len(result)),
		Stages:      result,
	})
}

// GetProjectProjectNameStageStageNameHandlerFunc gets the specified stage
func GetProjectProjectNameStageStageNameHandlerFunc(params stage.GetProjectProjectNameStageStageNameParams) middleware.Responder {
	common.Lock()
	defer common.UnLock()
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewGetProjectProjectNameStageStageNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	if !common.StageExists(params.ProjectName, params.StageName) {
		return stage.NewGetProjectProjectNameStageStageNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
	}

	var stageResponse = &models.Stage{
		StageName: params.StageName,
	}
	return stage.NewGetProjectProjectNameStageStageNameOK().WithPayload(stageResponse)
}
