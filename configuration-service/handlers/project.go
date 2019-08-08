package handlers

import (
	"io/ioutil"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
)

// GetProjectHandlerFunc gets a list of projects
func GetProjectHandlerFunc(params project.GetProjectParams) middleware.Responder {

	var payload = &models.Projects{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Projects:    []*models.Project{},
	}

	files, err := ioutil.ReadDir(config.ConfigDir)
	if err != nil {
		return project.NewGetProjectOK().WithPayload(payload)
	}

	paginationInfo := common.Paginate(len(files), params.PageSize, params.NextPageKey)

	totalCount := len(files)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, f := range files[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			if f.IsDir() {
				payload.Projects = append(payload.Projects, &models.Project{ProjectName: f.Name()})
			}
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return project.NewGetProjectOK().WithPayload(payload)
}

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams) middleware.Responder {

	projectConfigPath := config.ConfigDir + "/" + params.Project.ProjectName
	err := os.MkdirAll(projectConfigPath, os.ModePerm)
	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	out, err := utils.ExecuteCommandInDirectory("git", []string{"init"}, projectConfigPath)
	utils.Debug("", "Init git result: "+out)
	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	return project.NewPostProjectNoContent()
}

// GetProjectProjectNameHandlerFunc gets a project by its name
func GetProjectProjectNameHandlerFunc(params project.GetProjectProjectNameParams) middleware.Responder {
	return middleware.NotImplemented("operation project.GetProjectProjectName has not yet been implemented")
}

// PutProjectProjectNameHandlerFunc updates a project
func PutProjectProjectNameHandlerFunc(params project.PutProjectProjectNameParams) middleware.Responder {
	return middleware.NotImplemented("operation project.PutProjectProjectName has not yet been implemented")
}

// DeleteProjectProjectNameHandlerFunc deletes a project
func DeleteProjectProjectNameHandlerFunc(params project.DeleteProjectProjectNameParams) middleware.Responder {
	utils.Debug("", "Deleting project "+params.ProjectName)
	err := os.RemoveAll(config.ConfigDir + "/" + params.ProjectName)
	if err != nil {
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}
	utils.Debug("", "Project "+params.ProjectName+" has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
