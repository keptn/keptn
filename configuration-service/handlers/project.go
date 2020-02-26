package handlers

import (
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
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
			if f.IsDir() && common.FileExists(config.ConfigDir+"/"+f.Name()+"/metadata.yaml") {
				stages, _ := getStages(stage.GetProjectProjectNameStageParams{ProjectName: f.Name(), DisableUpstreamSync: params.DisableUpstreamSync})
				var project = &models.Project{ProjectName: f.Name(), Stages: stages}
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
	credentialsCreated := false
	logger := utils.NewLogger("", "", "configuration-service")
	projectConfigPath := config.ConfigDir + "/" + params.Project.ProjectName

	// check if the project already exists
	if common.ProjectExists(params.Project.ProjectName) {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project already exists")})
	}

	common.LockProject(params.Project.ProjectName)
	defer common.UnlockProject(params.Project.ProjectName)
	////////////////////////////////////////////////////
	// clone existing repo
	////////////////////////////////////////////////////
	if params.Project.GitUser != "" && params.Project.GitToken != "" && params.Project.GitRemoteURI != "" {
		// try to clone the repo
		err := common.CloneRepo(params.Project.ProjectName, params.Project.GitUser, params.Project.GitToken, params.Project.GitRemoteURI)
		if err != nil {
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not clone git repository")})
		}

		// store credentials (e.g., as a kubernetes secret)
		err = common.StoreGitCredentials(params.Project.ProjectName, params.Project.GitUser, params.Project.GitToken, params.Project.GitRemoteURI)
		if err != nil {
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not store git credentials")})
		}
		credentialsCreated = true
	} else {
		// if no remote URI has been specified, create a new repo
		///////////////////////////////////////////////////
		err := os.MkdirAll(projectConfigPath, os.ModePerm)
		if err != nil {
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not create project")})
		}

		_, err = utils.ExecuteCommandInDirectory("git", []string{"init"}, projectConfigPath)
		if err != nil {
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not initialize git repo")})
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
		logger.Error(err.Error())
		// Cleanup credentials before we exit
		if credentialsCreated {
			err = common.DeleteCredentials(params.Project.ProjectName)

			if err != nil {
				logger.Error(err.Error())
			}
		}

		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not store project metadata")})
	}

	err = common.StageAndCommitAll(params.Project.ProjectName, "Added metadata.yaml")
	if err != nil {
		logger.Error(err.Error())
		// Cleanup credentials before we exit
		if credentialsCreated {
			err = common.DeleteCredentials(params.Project.ProjectName)

			if err != nil {
				logger.Error(err.Error())
			}
		}
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	return project.NewPostProjectNoContent()
}

// GetProjectProjectNameHandlerFunc gets a project by its name
func GetProjectProjectNameHandlerFunc(params project.GetProjectProjectNameParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return project.NewGetProjectProjectNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
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
	logger := utils.NewLogger("", "", "configuration-service")
	logger.Debug("Deleting project " + params.ProjectName)
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)
	err := os.RemoveAll(config.ConfigDir + "/" + params.ProjectName)
	if err != nil {
		logger.Error(err.Error())
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete project")})
	}
	creds, _ := common.GetCredentials(params.ProjectName)
	if creds != nil {
		err = common.DeleteCredentials(params.ProjectName)
		if err != nil {
			logger.Error(err.Error())
			return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete upstream credentials")})
		}
	}

	logger.Debug("Project " + params.ProjectName + " has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
