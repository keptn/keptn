package handlers

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptnmodels "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
	"io/ioutil"
	"sort"
)

func getStages(params stage.GetProjectProjectNameStageParams) ([]*models.Stage, errors.Error) {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return nil, errors.New(404, "Project does not exist.")
	}

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch for project %s: %s", params.ProjectName, err.Error()))
		return nil, errors.New(500, "Could not determine default branch")
	}
	err = common.CheckoutBranch(params.ProjectName, defaultBranch, *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out default branch of project %s: %s", params.ProjectName, err.Error()))
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

	mv := common.GetProjectsMaterializedView()
	err = mv.CreateStage(params.ProjectName, params.Stage.StageName)
	if err != nil {
		return stage.NewPostProjectProjectNameStageBadRequest().WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
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
	var payload = &models.Stages{
		NextPageKey: "",
		PageSize:    0,
		Stages:      []*models.ExpandedStage{},
		TotalCount:  0,
	}
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return stage.NewGetProjectProjectNameStageDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return stage.NewGetProjectProjectNameStageNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	paginationInfo := common.Paginate(len(prj.Stages), params.PageSize, params.NextPageKey)

	allStagesOfProject := prj.Stages

	//sort stages alphabetically
	sort.Slice(allStagesOfProject, func(i, j int) bool {
		return allStagesOfProject[i].StageName < allStagesOfProject[j].StageName
	})

	totalCount := len(allStagesOfProject)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, stg := range allStagesOfProject[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Stages = append(payload.Stages, stg)
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return stage.NewGetProjectProjectNameStageOK().WithPayload(payload)
}

// GetProjectProjectNameStageStageNameHandlerFunc gets the specified stage
func GetProjectProjectNameStageStageNameHandlerFunc(params stage.GetProjectProjectNameStageStageNameParams) middleware.Responder {
	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)
	if err != nil {
		return stage.NewGetProjectProjectNameStageStageNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return stage.NewGetProjectProjectNameStageStageNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	for _, stg := range prj.Stages {
		if stg.StageName == params.StageName {
			return stage.NewGetProjectProjectNameStageStageNameOK().WithPayload(stg)
		}
	}
	return stage.NewGetProjectProjectNameStageStageNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
}
