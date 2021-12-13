package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
	logger "github.com/sirupsen/logrus"
)

// PostProjectProjectNameStageHandlerFunc creates a new stage
func PostProjectProjectNameStageHandlerFunc(params stage.PostProjectProjectNameStageParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist.")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not determine default branch for project %s", params.ProjectName)
		return stage.NewPostProjectProjectNameStageDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not create stage.")})
	}
	logger.Info(fmt.Sprintf("creating stage %s from base %s", params.Stage.StageName, defaultBranch))
	err = common.CreateBranch(params.ProjectName, params.Stage.StageName, defaultBranch)
	if err != nil {
		logger.WithError(err).Errorf("Could not create %s branch for project %s", params.Stage.StageName, params.ProjectName)
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
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewDeleteProjectProjectNameStageStageNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist.")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not determine default branch for project %s", params.ProjectName)
		return stage.NewDeleteProjectProjectNameStageStageNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not create stage.")})
	}
	logger.Info(fmt.Sprintf("creating stage %s from base %s", params.StageName, defaultBranch))
	err = common.DeleteBranch(params.ProjectName, params.StageName, defaultBranch)
	if err != nil {
		logger.WithError(err).Errorf("Could not create %s branch for project %s", params.StageName, params.ProjectName)
		return stage.NewDeleteProjectProjectNameStageStageNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not create stage.")})
	}

	return stage.NewPostProjectProjectNameStageNoContent()
}
