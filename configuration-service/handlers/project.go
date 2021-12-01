package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/martian/log"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams) middleware.Responder {
	// check if the project already exists
	//if common.ProjectExists(params.Project.ProjectName) {
	//	return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Project already exists")})
	//}

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
	credentials, err := common.GetCredentials(params.Project.ProjectName)
	if err != nil || credentials == nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("No upstream credentials found")})
	}

	client := common.NewGitClient()

	// try to clone the repo
	err = client.CloneRepo(params.Project.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not clone git repository during creating project %s", params.Project.ProjectName)
		rollbackFunc()
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not clone git repository")})
	}

	////////////////////////////////////////////////////

	return project.NewPostProjectNoContent()
}

// PutProjectProjectNameHandlerFunc updates a project
func PutProjectProjectNameHandlerFunc(params project.PutProjectProjectNameParams) middleware.Responder {
	projectName := params.Project.ProjectName

	if common.ProjectExists(projectName) {
		common.LockProject(projectName)
		defer common.UnlockProject(projectName)

		gitCredentials, err := common.GetCredentials(projectName)
		if err == nil && gitCredentials != nil {
			log.Infof("Storing Git credentials for project %s", projectName)
			if err := common.UpdateOrCreateOrigin(projectName); err != nil {
				logger.WithError(err).Errorf("Could not add upstream repository while updating project %s", params.Project.ProjectName)
				// TODO: use git library.
				// until we do not use a propper git library it is hard/not possible to
				// determine the correct error cases, so we need to rely on the output of the command
				if strings.Contains(err.Error(), common.GitURLNotFound) || strings.Contains(err.Error(), common.HostNotFound) {
					logger.Error("Invalid URL detected")
					return project.NewPutProjectProjectNameBadRequest().WithPayload(&models.Error{Code: http.StatusNotFound, Message: swag.String(err.Error())})
				}
				if strings.Contains(err.Error(), "Authentication failed") || strings.Contains(err.Error(), common.WrongToken) || strings.Contains(err.Error(), common.GitError) {
					logger.Error("Authentication error detected")
					return project.NewPutProjectProjectNameBadRequest().WithPayload(&models.Error{Code: http.StatusFailedDependency, Message: swag.String(err.Error())})
				}

				return project.NewPutProjectProjectNameDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: swag.String(err.Error())})
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
