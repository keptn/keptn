package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
	"gopkg.in/yaml.v2"
)

type projectMetadata struct {
	ProjectName       string
	CreationTimestamp string
}

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
				var project = &models.Project{ProjectName: f.Name()}
				payload.Projects = append(payload.Projects, project)
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

	////////////////////////////////////////////////////
	// clone existing repo
	////////////////////////////////////////////////////
	if params.Project.GitUser != "" && params.Project.GitToken != "" && params.Project.GitRemoteURI != "" {
		common.StoreGitCredentials(params.Project.ProjectName, params.Project.GitUser, params.Project.GitToken, params.Project.GitRemoteURI)
		common.CloneRepo(params.Project.ProjectName, params.Project.GitUser, params.Project.GitToken, params.Project.GitRemoteURI)
	} else {
		// if no remote URI has been specified, create a new repo
		////////////////////////////////////////////////////
		err := os.MkdirAll(projectConfigPath, os.ModePerm)
		if err != nil {
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
		}

		out, err := utils.ExecuteCommandInDirectory("git", []string{"init"}, projectConfigPath)
		utils.Debug("", "Init git result: "+out)
		if err != nil {
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
		}
	}
	////////////////////////////////////////////////////

	newProjectMetadata := &projectMetadata{
		ProjectName:       params.Project.ProjectName,
		CreationTimestamp: time.Now().String(),
	}

	metadataString, err := yaml.Marshal(newProjectMetadata)
	err = common.WriteFile(projectConfigPath+"/metadata.yaml", metadataString)

	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	_, err = utils.ExecuteCommandInDirectory("git", []string{"add", "."}, projectConfigPath)
	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	_, err = utils.ExecuteCommandInDirectory("git", []string{"commit", "-m", `"added metadata.yaml"`}, projectConfigPath)
	if err != nil {
		fmt.Print(err.Error())
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
	}

	return project.NewPostProjectNoContent()
}

// GetProjectProjectNameHandlerFunc gets a project by its name
func GetProjectProjectNameHandlerFunc(params project.GetProjectProjectNameParams) middleware.Responder {
	var projectResponse = &models.Project{ProjectName: params.ProjectName}
	projectCreds, _ := common.GetCredentials(params.ProjectName)
	if projectCreds != nil {
		projectResponse.GitRemoteURI = projectCreds.RemoteURI
	}
	return project.NewGetProjectProjectNameOK().WithPayload(projectResponse)
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
	creds, _ := common.GetCredentials(params.ProjectName)
	if creds != nil {
		err = common.DeleteCredentials(params.ProjectName)
		if err != nil {
			return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
		}
	}

	utils.Debug("", "Project "+params.ProjectName+" has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
