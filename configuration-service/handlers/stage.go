package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
)

// PostProjectProjectNameStageHandlerFunc creates a new stage
func PostProjectProjectNameStageHandlerFunc(params stage.PostProjectProjectNameStageParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist.")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch for project %s: %s", params.ProjectName, err.Error()))
		return stage.NewPostProjectProjectNameStageDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not create stage.")})
	}
	logger.Info(fmt.Sprintf("creating stage %s from base %s", params.Stage.StageName, defaultBranch))
	err = common.CreateBranch(params.ProjectName, params.Stage.StageName, defaultBranch)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not create %s branch for project %s: %s", params.Stage.StageName, params.ProjectName, err.Error()))
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
