package handlers

import (
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"os"
	"sort"
	"time"

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

	//sort projects alphabetically
	sort.Slice(allProjects, func(i, j int) bool {
		return allProjects[i].ProjectName < allProjects[j].ProjectName
	})

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
	logger := keptncommon.NewLogger("", "", "configuration-service")
	// check if the project already exists
	if common.ProjectExists(params.Project.ProjectName) {
		common.LockProject(params.Project.ProjectName)
		defer common.UnlockProject(params.Project.ProjectName)

		logger.Debug("Updating project " + params.ProjectName)

		mv := common.GetProjectsMaterializedView()
		logger.Debug("Add or update Git origin and push changes for project " + params.ProjectName)
		projectInfo, err := mv.GetProject(params.Project.ProjectName)
		if err != nil {
			msg := "could not read project information: " + err.Error()
			logger.Error(msg)
			return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(msg)})
		}

		oldRemoteURL := projectInfo.GitRemoteURI
		oldRemoteUser := projectInfo.GitUser
		credentials, err := common.GetCredentials(params.Project.ProjectName)
		if err == nil && credentials != nil {
			logger.Debug("Storing Git credentials for project " + params.ProjectName)

			err = common.UpdateOrCreateOrigin(params.Project.ProjectName)
			if err != nil {
				logger.Error(fmt.Sprintf("Could not add upstream repository while updating project %s: %v", params.Project.ProjectName, err))
				if oldRemoteURL != "" && oldRemoteUser != "" {
					if restoreErr := mv.UpdateUpstreamInfo(params.ProjectName, oldRemoteURL, oldRemoteUser); restoreErr != nil {
						logger.Error(fmt.Sprintf("could not restore upstream info in materializer view to previous values: %s", err.Error()))
					}
				} else if deleteErr := mv.DeleteUpstreamInfo(params.ProjectName); deleteErr != nil {
					logger.Error(fmt.Sprintf("Could not delete upstream info from materialized view: %s", err.Error()))
				}
				return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
			}

			if err := mv.UpdateUpstreamInfo(params.Project.ProjectName, credentials.RemoteURI, credentials.User); err != nil {
				logger.Error(fmt.Sprintf("Could not add upstream repository info for project %s: %v", params.Project.ProjectName, err))
				return project.NewPostProjectDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
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

	mv := common.GetProjectsMaterializedView()
	err = mv.DeleteProject(params.ProjectName)
	if err != nil {
		return project.NewDeleteProjectProjectNameBadRequest().WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}

	logger.Debug("Project " + params.ProjectName + " has been deleted")
	return project.NewDeleteProjectProjectNameNoContent()
}
