package handlers

import (
	"github.com/ghodss/yaml"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptnmodels "github.com/keptn/go-utils/pkg/lib"
	utils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
	"io/ioutil"
)

func getStages(params stage.GetProjectProjectNameStageParams) ([]*models.Stage, errors.Error) {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return nil, errors.New(404, "Project does not exist.")
	}

	err := common.CheckoutBranch(params.ProjectName, "master", *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(500, "Could not retrieve stages.")
	}

	shipyardPath := config.ConfigDir + "/" + params.ProjectName + "/shipyard.yaml"

	if !common.FileExists(shipyardPath) {
		return nil, errors.New(500, "Could not retrieve stages.")
	}

	dat, err := ioutil.ReadFile(shipyardPath)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(500, "Could not read shipyard file.")
	}

	shipyard := &keptnmodels.Shipyard{}

	err = yaml.Unmarshal(dat, shipyard)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(500, "Could not read shipyard file.")
	}

	var result []*models.Stage
	for _, stage := range shipyard.Stages {
		stage := &models.Stage{
			StageName: stage.Name,
		}
		result = append(result, stage)
	}

	return result, nil
}

// PostProjectProjectNameStageHandlerFunc creates a new stage
func PostProjectProjectNameStageHandlerFunc(params stage.PostProjectProjectNameStageParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist.")})
	}
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	err := common.CreateBranch(params.ProjectName, params.Stage.StageName, "master")
	if err != nil {
		logger.Error(err.Error())
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not create stage.")})
	}

	mv := common.GetProjectsMaterializedView()
	mv.CreateStage(params.ProjectName, params.Stage.StageName)

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
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	result, err := getStages(params)

	if err != nil {
		if err.Code() == 404 {
			return stage.NewGetProjectProjectNameStageNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String(err.Error())})
		} else {
			return stage.NewGetProjectProjectNameStageDefault(int(err.Code())).WithPayload(&models.Error{Code: int64(err.Code()), Message: swag.String(err.Error())})
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
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	if !common.ProjectExists(params.ProjectName) {
		return stage.NewGetProjectProjectNameStageStageNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	if !common.StageExists(params.ProjectName, params.StageName, *params.DisableUpstreamSync) {
		return stage.NewGetProjectProjectNameStageStageNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
	}

	var stageResponse = &models.Stage{
		StageName: params.StageName,
	}
	return stage.NewGetProjectProjectNameStageStageNameOK().WithPayload(stageResponse)
}
