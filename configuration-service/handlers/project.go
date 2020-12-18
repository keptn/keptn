package handlers

import (
	"fmt"
	"os"
	"time"

	keptn "github.com/keptn/go-utils/pkg/lib"
	k8sutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
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
	var payload = &models.ExpandedProjects{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Projects:    []*models.ExpandedProject{},
	}

	mv := common.GetProjectsMaterializedView()

	allProjects, err := mv.GetProjects()

	if err != nil || allProjects == nil {
		return project.NewGetProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	paginationInfo := common.Paginate(len(allProjects), params.PageSize, params.NextPageKey)

	totalCount := len(allProjects)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, project := range allProjects[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Projects = append(payload.Projects, project)
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return project.NewGetProjectOK().WithPayload(payload)
}

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams) middleware.Responder {
	credentialsCreated := false
	logger := keptn.NewLogger("", "", "configuration-service")
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
	if params.Project.GitUser != "" && params.Project.GitToken != "" && params.Project.GitRemoteURI != "" {
		// try to clone the repo
		var err error
		credentials := common.GitCredentials{
			User:      params.Project.GitUser,
			Token:     params.Project.GitToken,
			RemoteURI: params.Project.GitRemoteURI,
		}
		initializedGit, err = common.CloneRepo(params.Project.ProjectName, credentials)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not clone git repository during creating project %s", params.Project.ProjectName))
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not clone git repository")})
		}

		// store credentials (e.g., as a kubernetes secret)
		err = common.StoreGitCredentials(params.Project.ProjectName, params.Project.GitUser, params.Project.GitToken, params.Project.GitRemoteURI)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not store git credentials during creating project %s", params.Project.ProjectName))
			logger.Error(err.Error())
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not store git credentials")})
		}
		credentialsCreated = true
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
		// Cleanup credentials before we exit
		if credentialsCreated {
			err = common.DeleteCredentials(params.Project.ProjectName)
			if err != nil {
				logger.Error(fmt.Sprintf("Could not delete credentials during creating project %s", params.Project.ProjectName))
				logger.Error(err.Error())
			}
		}

		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not store project metadata")})
	}

	err = common.StageAndCommitAll(params.Project.ProjectName, "Added metadata.yaml", initializedGit)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit metadata.yaml during creating project %s", params.Project.ProjectName))
		logger.Error(err.Error())
		// Cleanup credentials before we exit
		if credentialsCreated {
			err = common.DeleteCredentials(params.Project.ProjectName)
			if err != nil {
				logger.Error(fmt.Sprintf("Could not delete credentials during creating project %s", params.Project.ProjectName))
				logger.Error(err.Error())
			}
		}
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}

	// make sure that master branch exists
	logger.Info("making sure that master branch exists")
	err = common.EnsureBranchAvailability(params.Project.ProjectName, "master")
	if err != nil {
		logger.Info("Could not ensure that master branch is available for project " + params.Project.ProjectName + ": " + err.Error())
		return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not ensure that master branch is available for project " + params.Project.ProjectName)})
	}
	logger.Info("master branch has been stored")

	mv := common.GetProjectsMaterializedView()

	err = mv.CreateProject(params.Project)

	if err != nil {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}
	return project.NewPostProjectNoContent()
}

// GetProjectProjectNameHandlerFunc gets a project by its name
func GetProjectProjectNameHandlerFunc(params project.GetProjectProjectNameParams) middleware.Responder {

	mv := common.GetProjectsMaterializedView()

	prj, err := mv.GetProject(params.ProjectName)

	if err != nil {
		return project.NewGetProjectProjectNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	if prj == nil {
		return project.NewGetProjectProjectNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	return project.NewGetProjectProjectNameOK().WithPayload(prj)
}

// PutProjectProjectNameHandlerFunc updates a project
func PutProjectProjectNameHandlerFunc(params project.PutProjectProjectNameParams) middleware.Responder {
	logger := keptn.NewLogger("", "", "configuration-service")
	// check if the project already exists
	if common.ProjectExists(params.Project.ProjectName) {
		common.LockProject(params.Project.ProjectName)
		defer common.UnlockProject(params.Project.ProjectName)

		logger.Debug("Updating project " + params.ProjectName)

		if params.Project.GitUser != "" && params.Project.GitToken != "" && params.Project.GitRemoteURI != "" {
			logger.Debug("Storing Git credentials for project " + params.ProjectName)

			// store credentials (e.g., as a kubernetes secret)
			err := common.StoreGitCredentials(params.Project.ProjectName, params.Project.GitUser, params.Project.GitToken, params.Project.GitRemoteURI)
			if err != nil {
				logger.Error(fmt.Sprintf("Could not store git credentials during creating project %s: %v", params.Project.ProjectName, err))
				return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not store git credentials")})
			}

			logger.Debug("Add Git origin and push changes for project " + params.ProjectName)
			err = common.AddOrigin(params.Project.ProjectName)
			if err != nil {
				logger.Error(fmt.Sprintf("Could not add upstream repository while updating project %s: %v", params.Project.ProjectName, err))
				// Cleanup credentials before we exit
				if credsErr := common.DeleteCredentials(params.Project.ProjectName); credsErr != nil {
					logger.Error(fmt.Sprintf("Could not delete credentials while updating project %s: %v", params.Project.ProjectName, credsErr))
				}
				return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
			}
			mv := common.GetProjectsMaterializedView()
			// try to update the materialized view. If this fails, it should not prevent the further execution
			err = mv.UpdateUpstreamInfo(params.Project.ProjectName, params.Project.GitRemoteURI, params.Project.GitUser)
			if err != nil {
				logger.Info("Could not update project materialized view " + params.ProjectName + ": " + err.Error())
				return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not update materialized view of project " + params.ProjectName)})
			}
			// make sure that master branch exists
			err = common.EnsureBranchAvailability(params.ProjectName, "master")
			if err != nil {
				logger.Info("Could not ensure that master branch is available for project " + params.ProjectName + ": " + err.Error())
				return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not ensure that master branch is available for project " + params.ProjectName)})
			}
		} else {
			logger.Info("Project " + params.ProjectName + " not updated as Git credentials were missing.")
			return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project not updated as Git credentials were missing")})
		}

	} else {
		return project.NewPostProjectBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}
	return project.NewPutProjectProjectNameNoContent()
}

// DeleteProjectProjectNameHandlerFunc deletes a project
func DeleteProjectProjectNameHandlerFunc(params project.DeleteProjectProjectNameParams) middleware.Responder {
	logger := keptn.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	logger.Debug("Deleting project " + params.ProjectName)

	err := os.RemoveAll(config.ConfigDir + "/" + params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not delete directory during deleting project %s", params.ProjectName))
		logger.Error(err.Error())
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete project")})
	}

	creds, _ := common.GetCredentials(params.ProjectName)
	if creds != nil {
		err = common.DeleteCredentials(params.ProjectName)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not delete credentials during deleting project %s", params.ProjectName))
			logger.Error(err.Error())
			return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete upstream credentials")})
		}
	}

	mv := common.GetProjectsMaterializedView()
	err = mv.DeleteProject(params.ProjectName)
	if err != nil {
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	logger.Debug("Project " + params.ProjectName + " has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
