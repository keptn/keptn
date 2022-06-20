package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptn2 "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
	"gopkg.in/yaml.v3"
)

type projectMetadata struct {
	ProjectName       string
	CreationTimestamp string
}

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams) middleware.Responder {
	projectConfigPath := config.ConfigDir + "/" + params.Project.ProjectName

	// check if the project already exists
	if common.ProjectExists(params.Project.ProjectName) {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Project already exists")})
	}

	common.LockProject(params.Project.ProjectName)
	defer common.UnlockProject(params.Project.ProjectName)

	rollbackFunc := func() {
		logger.Infof("Rollback: try to delete created directory for project %s", params.Project.ProjectName)
		if err := os.RemoveAll(config.ConfigDir + "/" + params.Project.ProjectName); err != nil {
			logger.Errorf("Rollback failed: could not delete created directory for project %s: %s", params.Project.ProjectName, err.Error())
		}
	}

	////////////////////////////////////////////////////
	// clone existing repo
	////////////////////////////////////////////////////
	var initializedGit bool
	credentials, err := common.GetCredentials(params.Project.ProjectName)
	if err == nil && credentials != nil {
		// try to clone the repo
		var err error

		initializedGit, err = common.CloneRepo(params.Project.ProjectName, *credentials)
		if err != nil {
			logger.WithError(err).Errorf("Could not clone git repository during creating project %s", params.Project.ProjectName)
			rollbackFunc()
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not clone git repository")})
		}

	} else {
		// if no remote URI has been specified, create a new repo
		///////////////////////////////////////////////////
		err := os.MkdirAll(projectConfigPath, os.ModePerm)
		if err != nil {
			logger.WithError(err).Errorf("Could make directory during creating project %s", params.Project.ProjectName)
			rollbackFunc()
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not create project")})
		}

		_, err = keptn2.ExecuteCommandInDirectory("git", []string{"init"}, projectConfigPath)
		if err != nil {
			logger.WithError(err).Errorf("Could not initialize git repository during creating project %s", params.Project.ProjectName)
			rollbackFunc()
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not initialize git repo")})
		}
	}
	err = common.ConfigureGitUser(params.Project.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not configure git during creating project %s", params.Project.ProjectName)
		rollbackFunc()
		return project.NewPostProjectDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: swag.String("Could not configure git in project repo")})
	}
	////////////////////////////////////////////////////
	newProjectMetadata := &projectMetadata{
		ProjectName:       params.Project.ProjectName,
		CreationTimestamp: time.Now().String(),
	}

	metadataString, err := yaml.Marshal(newProjectMetadata)

	err = common.WriteFile(projectConfigPath+"/metadata.yaml", metadataString)
	if err != nil {
		logger.WithError(err).Errorf("Could not write metadata.yaml during creating project %s", params.Project.ProjectName)
		rollbackFunc()
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not store project metadata")})
	}

	err = common.StageAndCommitAll(params.Project.ProjectName, "Added metadata.yaml", initializedGit)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit metadata.yaml during creating project %s", params.Project.ProjectName)
		rollbackFunc()
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not commit changes")})
	}
	return project.NewPostProjectNoContent()
}

// PutProjectProjectNameHandlerFunc updates a project
func PutProjectProjectNameHandlerFunc(params project.PutProjectProjectNameParams) middleware.Responder {
	projectName := params.Project.ProjectName

	if common.ProjectExists(projectName) {
		common.LockProject(projectName)
		defer common.UnlockProject(projectName)

		err := common.ConfigureGitUser(params.Project.ProjectName)
		if err != nil {
			logger.WithError(err).Errorf("Could not configure git during creating project %s", params.Project.ProjectName)
			return project.NewPostProjectDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: swag.String("Could not configure git in project repo")})
		}

		gitCredentials, err := common.GetCredentials(projectName)
		if err == nil && gitCredentials != nil {
			logger.Infof("Storing Git credentials for project %s", projectName)
			if err := common.UpdateOrCreateOrigin(projectName); err != nil {
				logger.WithError(err).Errorf("Could not add upstream repository while updating project %s", params.Project.ProjectName)
				// TODO: use git library.
				// until we do not use a propper git library it is hard/not possible to
				// determine the correct error cases, so we need to rely on the output of the command
				if strings.Contains(err.Error(), common.GitURLNotFound) || strings.Contains(err.Error(), common.HostNotFound) {
					logger.Error("Invalid URL detected")
					return project.NewPutProjectProjectNameBadRequest().WithPayload(&models.Error{Code: http.StatusNotFound, Message: swag.String(common.RepositoryNotFoundErrorMsg)})
				}
				if strings.Contains(err.Error(), "Authentication failed") || strings.Contains(err.Error(), common.WrongToken) || strings.Contains(err.Error(), common.GitError) {
					logger.Error("Authentication error detected")
					return project.NewPutProjectProjectNameBadRequest().WithPayload(&models.Error{Code: http.StatusNotFound, Message: swag.String(common.RepositoryNotFoundErrorMsg)})
				}

				return project.NewPutProjectProjectNameDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: swag.String(common.InternalErrorErrMsg)})
			}
		}
	} else {
		return project.NewPutProjectProjectNameDefault(http.StatusBadRequest).WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String(common.ProjectDoesNotExistErrorMsg)})
	}
	return project.NewPutProjectProjectNameNoContent()
}

// DeleteProjectProjectNameHandlerFunc deletes a project
func DeleteProjectProjectNameHandlerFunc(params project.DeleteProjectProjectNameParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	logger.Debug("Deleting project " + params.ProjectName)

	err := os.RemoveAll(config.ConfigDir + "/" + params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not delete directory during deleting project %s", params.ProjectName)
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not delete project")})
	}

	logger.Debug("Project " + params.ProjectName + " has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
