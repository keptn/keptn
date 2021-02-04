package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"strconv"
	"strings"
	"time"
)

var ErrProjectAlreadyExists = errors.New("project already exists")

type ProjectManager struct {
	Logger                  keptncommon.LoggerInterface
	ConfigurationStore      common.ConfigurationStore
	SecretStore             common.SecretStore
	ProjectMaterializedView db.ProjectsDBOperations
	TaskSequenceRepository  db.TaskSequenceRepo
	EventRpository          db.EventRepo
}

type rollbackfunc func() error

var nilRollback = func() error {
	return nil
}

func NewProjectManager(
	configurationStore common.ConfigurationStore,
	secretStore common.SecretStore,
	dbProjectsOperations db.ProjectsDBOperations,
	taskSequenceRepo db.TaskSequenceRepo,
	eventRepo db.EventRepo) *ProjectManager {

	projectUpdater := &ProjectManager{
		ConfigurationStore:      configurationStore,
		SecretStore:             secretStore,
		ProjectMaterializedView: dbProjectsOperations,
		TaskSequenceRepository:  taskSequenceRepo,
		EventRpository:          eventRepo,
		Logger:                  keptncommon.NewLogger("", "", "shipyard-controller"),
	}
	return projectUpdater
}

func (pm *ProjectManager) Get() ([]*models.ExpandedProject, error) {
	pm.Logger.Info("Getting all projects")
	allProjects, err := pm.ProjectMaterializedView.GetProjects()
	if err != nil {
		return nil, err
	}
	return allProjects, nil
}

func (pm *ProjectManager) GetByName(projectName string) (*models.ExpandedProject, error) {
	pm.Logger.Info("Getting project with name " + projectName)
	project, err := pm.ProjectMaterializedView.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	return project, err

}

func (pu *ProjectManager) Create(params *operations.CreateProjectParams) (error, rollbackfunc) {

	existingProject, err := pu.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil {
		return err, nilRollback
	}
	if existingProject != nil {
		return ErrProjectAlreadyExists, nilRollback
	}

	err = pu.updateGITRepositorySecret(*params.Name, &gitCredentials{
		User:      params.GitUser,
		Token:     params.GitToken,
		RemoteURI: params.GitRemoteURL,
	})
	if err != nil {
		return err, nilRollback
	}

	err = pu.ConfigurationStore.CreateProject(apimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	})

	if err != nil {
		pu.Logger.Error(fmt.Sprintf("Error occured while creating project in configuration service: %s", err.Error()))
		return err, func() error {
			pu.Logger.Info(fmt.Sprintf("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name))
			if err := pu.deleteGITRepositorySecret(*params.Name); err != nil {
				pu.Logger.Error(fmt.Sprintf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error()))
				return err
			}
			return nil
		}
	}
	pu.Logger.Info(fmt.Sprintf("Created project in configuration service: %s", *params.Name))

	decodedShipyard, _ := base64.StdEncoding.DecodeString(*params.Shipyard)
	shipyard, err := common.UnmarshalShipyard(string(decodedShipyard))
	for _, shipyardStage := range shipyard.Spec.Stages {
		if err := pu.ConfigurationStore.CreateStage(*params.Name, shipyardStage.Name); err != nil {
			return err, nil
		}
		pu.Logger.Info(fmt.Sprintf("Stage %s created", shipyardStage.Name))
	}
	pu.Logger.Info(fmt.Sprintf("Created all stages of project %s", *params.Name))

	uri := "shipyard.yaml"
	projectResource := []*apimodels.Resource{
		{
			ResourceContent: string(decodedShipyard),
			ResourceURI:     &uri,
		},
	}
	if err := pu.ConfigurationStore.CreateProjectShipyard(*params.Name, projectResource); err != nil {
		pu.Logger.Error(fmt.Sprintf("Error occured while uploading shipyard resource to configuraiton service: %s", err.Error()))
		return err, func() error {
			pu.Logger.Info(fmt.Sprintf("Rollback: Try to delete project %s from configuration service", *params.Name))
			if err := pu.ConfigurationStore.DeleteProject(*params.Name); err != nil {
				pu.Logger.Error(fmt.Sprintf("Rollback failed: Unable to delete project %s from configuration service: %s", *params.Name, err.Error()))
				return err
			}
			pu.Logger.Info(fmt.Sprintf("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name))
			if err := pu.deleteGITRepositorySecret(*params.Name); err != nil {
				pu.Logger.Error(fmt.Sprintf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error()))
				return err
			}
			return nil
		}
	}

	if err := pu.createProjectInRepository(params); err != nil {
		return err, func() error {
			if err := pu.ConfigurationStore.DeleteProject(*params.Name); err != nil {
				return err
			}
			if err := pu.deleteGITRepositorySecret(*params.Name); err != nil {
				return err
			}
			return nil
		}
	}
	return nil, nilRollback

}

func (pu *ProjectManager) Update(params *operations.UpdateProjectParams) (error, rollbackfunc) {
	oldSecret, err := pu.getGITRepositorySecret(*params.Name)
	if err != nil {
		return err, nilRollback
	}
	oldProject, err := pu.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil {
		return err, nilRollback
	}

	err = pu.updateGITRepositorySecret(*params.Name, &gitCredentials{
		User:      params.GitUser,
		Token:     params.GitToken,
		RemoteURI: params.GitRemoteURL,
	})
	if err != nil {
		return err, nilRollback
	}

	projectToUpdate := apimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	projectToRollback := apimodels.Project{
		CreationDate:    oldProject.CreationDate,
		GitRemoteURI:    oldProject.GitRemoteURI,
		GitUser:         oldProject.GitUser,
		ProjectName:     oldProject.ProjectName,
		ShipyardVersion: oldProject.ShipyardVersion,
	}

	err = pu.ConfigurationStore.UpdateProject(projectToUpdate)
	if err != nil {
		return err, func() error {
			return pu.updateGITRepositorySecret(*params.Name, &gitCredentials{
				User:      oldSecret.User,
				Token:     oldSecret.Token,
				RemoteURI: oldSecret.RemoteURI,
			})
		}
	}

	err = pu.ProjectMaterializedView.UpdateUpstreamInfo(*params.Name, params.GitRemoteURL, params.GitUser)
	if err != nil {
		return err, func() error {

			errConfigStoreRollback := pu.ConfigurationStore.UpdateProject(projectToRollback)
			if errConfigStoreRollback != nil {
				return errConfigStoreRollback
			}
			return pu.updateGITRepositorySecret(*params.Name, &gitCredentials{
				User:      oldSecret.User,
				Token:     oldSecret.Token,
				RemoteURI: oldSecret.RemoteURI,
			})
		}
	}

	return nil, nilRollback
}

func (pu *ProjectManager) Delete(projectName string) (error, string) {
	pu.Logger.Info(fmt.Sprintf("Deleting project %s", projectName))
	var resultMessage strings.Builder

	project, err := pu.ProjectMaterializedView.GetProject(projectName)
	if err != nil {
		resultMessage.WriteString(fmt.Sprintf("Project %s cannot be retrieved anymore. Any Git upstream of the project will not be deleted.\n", projectName))
	} else if project != nil && project.GitRemoteURI != "" {
		resultMessage.WriteString(fmt.Sprintf("The Git upstream of the project will not be deleted: %s\n", project.GitRemoteURI))
	}

	secret, err := pu.SecretStore.GetSecret("git-credentials-" + projectName)
	if err != nil {
		pu.Logger.Error("could not delete git upstream credentials secret: " + err.Error())
	}
	if secret != nil {
		if err := pu.SecretStore.DeleteSecret("git-credentials-" + projectName); err != nil {
			pu.Logger.Error("could not delete git upstream credentials secret: " + err.Error())
			resultMessage.WriteString("WARNING: Could not delete secret containing the git upstream repo credentials. \n")
			resultMessage.WriteString(fmt.Sprintf("Please make sure to delete the secret manually by executing 'kubectl delete secret %s -n %s' \n", "git-credentials-"+projectName, common.GetKeptnNamespace()))
		}
	}

	if err := pu.ConfigurationStore.DeleteProject(projectName); err != nil {
		return pu.logAndReturnError(fmt.Sprintf("could not delete project: %s", err.Error())), resultMessage.String()
	}

	resultMessage.WriteString(pu.getDeleteInfoMessage(projectName))

	if err := pu.EventRpository.DeleteEventCollections(projectName); err != nil {
		pu.Logger.Error(fmt.Sprintf("could not delete task sequence collection: %s", err.Error()))
	}

	if err := pu.TaskSequenceRepository.DeleteTaskSequenceCollection(projectName); err != nil {
		pu.Logger.Error(fmt.Sprintf("could not delete task equence colleciton: %s", err.Error()))
	}

	if err := pu.ProjectMaterializedView.DeleteProject(projectName); err != nil {
		pu.Logger.Error(fmt.Sprintf("could not delete project: %s", err.Error()))
	}

	return nil, resultMessage.String()

}

func (pu *ProjectManager) createProjectInRepository(params *operations.CreateProjectParams) error {
	p := &apimodels.Project{
		CreationDate: strconv.FormatInt(time.Now().UnixNano(), 10),
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}
	err := pu.ProjectMaterializedView.CreateProject(p)
	if err != nil {
		return err
	}
	return nil
}

func (pu *ProjectManager) getGITRepositorySecret(projectName string) (*gitCredentials, error) {
	secret, err := pu.SecretStore.GetSecret("git-credentials-" + projectName)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	if marshalledSecret, ok := secret["git-credentials"]; ok {
		secretObj := &gitCredentials{}
		if err := json.Unmarshal(marshalledSecret, secretObj); err != nil {
			return nil, err
		}
		return secretObj, nil
	}
	return nil, nil
}

func (pu *ProjectManager) updateGITRepositorySecret(projectName string, credentials *gitCredentials) error {

	credsEncoded, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	if err := pu.SecretStore.UpdateSecret("git-credentials-"+projectName, map[string][]byte{
		"git-credentials": credsEncoded,
	}); err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	return nil
}

func (pu *ProjectManager) deleteGITRepositorySecret(projectName string) error {
	pu.Logger.Info("deleting git credentials for project " + projectName)

	if err := pu.SecretStore.DeleteSecret("git-credentials-" + projectName); err != nil {
		return fmt.Errorf("could not delete git credentials: %s", err.Error())
	}
	pu.Logger.Info("deleted git credentials for project " + projectName)
	return nil

}

func (pu *ProjectManager) getDeleteInfoMessage(project string) string {
	res, err := pu.ConfigurationStore.GetProjectResource(project, "shipyard.yaml")
	if err != nil {
		return getShipyardNotAvailableError(project)
	}

	shipyard := &keptnv2.Shipyard{}
	err = yaml.Unmarshal([]byte(res.ResourceContent), shipyard)
	if err != nil {
		return getShipyardNotAvailableError(project)
	}

	msg := "\n"
	for _, stage := range shipyard.Spec.Stages {
		namespace := project + "-" + stage.Name
		msg += fmt.Sprintf("- A potentially created namespace %s is not managed by Keptn anymore but is not deleted. "+
			"If you would like to delete this namespace, please execute "+
			"'kubectl delete ns %s'\n", namespace, namespace)
	}
	return strings.TrimSpace(msg)
}

func (pu *ProjectManager) logAndReturnError(msg string) error {
	pu.Logger.Error(msg)
	return errors.New(msg)
}

func getShipyardNotAvailableError(project string) string {
	return fmt.Sprintf("Shipyard of project %s cannot be retrieved anymore. "+
		"After deleting the project, the namespaces containing the services are still available. "+
		"This may cause problems if a project with the same name is created later.", project)
}

func toModelProject(project models.ExpandedProject) apimodels.Project {
	return apimodels.Project{
		CreationDate:    project.CreationDate,
		GitRemoteURI:    project.GitRemoteURI,
		GitUser:         project.GitUser,
		ProjectName:     project.ProjectName,
		ShipyardVersion: project.ShipyardVersion,
	}
}

func toExpandedProject(project apimodels.Project) models.ExpandedProject {
	return models.ExpandedProject{
		CreationDate:    project.CreationDate,
		GitRemoteURI:    project.GitRemoteURI,
		GitUser:         project.GitUser,
		ProjectName:     project.ProjectName,
		ShipyardVersion: project.ShipyardVersion,
	}
}

type gitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}
