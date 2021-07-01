package handlers

import (
	"fmt"
	"github.com/google/martian/log"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"os"
	"time"

	k8sutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/ghodss/yaml"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
)

type projectMetadata struct {
	ProjectName       string
	CreationTimestamp string
}

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
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
	var initializedGit bool
	credentials, err := common.GetCredentials(params.Project.ProjectName)
	if err == nil && credentials != nil {
		// try to clone the repo
		var err error

		initializedGit, err = common.CloneRepo(params.Project.ProjectName, *credentials)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not clone git repository during creating project %s", params.Project.ProjectName))
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not clone git repository")})
		}

	} else {
		// if no remote URI has been specified, create a new repo
		///////////////////////////////////////////////////
		err := os.MkdirAll(projectConfigPath, os.ModePerm)
		if err != nil {
			logger.Error(fmt.Sprintf("Could make directory during creating project %s", params.Project.ProjectName))
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not create project")})
		}

		_, err = k8sutils.ExecuteCommandInDirectory("git", []string{"init"}, projectConfigPath)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not initialize git repository during creating project %s", params.Project.ProjectName))
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
		logger.Error(fmt.Sprintf("Could not write metadata.yaml during creating project %s", params.Project.ProjectName))
		logger.Error(err.Error())

		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not store project metadata")})
	}

	err = common.StageAndCommitAll(params.Project.ProjectName, "Added metadata.yaml", initializedGit)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit metadata.yaml during creating project %s", params.Project.ProjectName))
		logger.Error(err.Error())
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	return project.NewPostProjectNoContent()
}

// PutProjectProjectNameHandlerFunc updates a project
func PutProjectProjectNameHandlerFunc(params project.PutProjectProjectNameParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")

	projectName := params.Project.ProjectName

	if common.ProjectExists(projectName) {
		common.LockProject(projectName)
		defer common.UnlockProject(projectName)

		gitCredentials, err := common.GetCredentials(projectName)
		if err == nil && gitCredentials != nil {
			log.Infof("Storing Git credentials for project %s", projectName)
			if err := common.UpdateOrCreateOrigin(projectName); err != nil {
				logger.Error(fmt.Sprintf("Could not add upstream repository while updating project %s: %v", params.Project.ProjectName, err))
				return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
			}

		}
	} else {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	return project.NewPutProjectProjectNameNoContent()
}

// DeleteProjectProjectNameHandlerFunc deletes a project
func DeleteProjectProjectNameHandlerFunc(params project.DeleteProjectProjectNameParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	logger.Debug("Deleting project " + params.ProjectName)

	err := os.RemoveAll(config.ConfigDir + "/" + params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not delete directory during deleting project %s", params.ProjectName))
		logger.Error(err.Error())
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete project")})
	}

	logger.Debug("Project " + params.ProjectName + " has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
