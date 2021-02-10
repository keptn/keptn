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

const shipyardVersion = "spec.keptn.sh/0.2.0"

//go:generate moq -pkg fake -skip-ensure -out ./fake/projectmanager.go . IProjectManager
type IProjectManager interface {
	Get() ([]*models.ExpandedProject, error)
	GetByName(projectName string) (*models.ExpandedProject, error)
	Create(params *operations.CreateProjectParams) (error, common.RollbackFunc)
	Update(params *operations.UpdateProjectParams) (error, common.RollbackFunc)
	Delete(projectName string) (error, string)
}

type ProjectManager struct {
	Logger                  keptncommon.LoggerInterface
	ConfigurationStore      common.ConfigurationStore
	SecretStore             common.SecretStore
	ProjectMaterializedView db.ProjectsDBOperations
	TaskSequenceRepository  db.TaskSequenceRepo
	EventRpository          db.EventRepo
}

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

func (pm *ProjectManager) Create(params *operations.CreateProjectParams) (error, common.RollbackFunc) {

	existingProject, err := pm.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil {
		return err, nilRollback
	}
	if existingProject != nil {
		return ErrProjectAlreadyExists, nilRollback
	}

	err = pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
		User:      params.GitUser,
		Token:     params.GitToken,
		RemoteURI: params.GitRemoteURL,
	})
	if err != nil {
		return err, nilRollback
	}

	err = pm.ConfigurationStore.CreateProject(apimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	})

	if err != nil {
		pm.Logger.Error(fmt.Sprintf("Error occurred while creating project in configuration service: %s", err.Error()))
		return err, func() error {
			pm.Logger.Info(fmt.Sprintf("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name))
			if err := pm.deleteGITRepositorySecret(*params.Name); err != nil {
				pm.Logger.Error(fmt.Sprintf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error()))
				return err
			}
			return nil
		}
	}
	pm.Logger.Info(fmt.Sprintf("Created project in configuration service: %s", *params.Name))

	decodedShipyard, _ := base64.StdEncoding.DecodeString(*params.Shipyard)
	shipyard, err := common.UnmarshalShipyard(string(decodedShipyard))
	for _, shipyardStage := range shipyard.Spec.Stages {
		if err := pm.ConfigurationStore.CreateStage(*params.Name, shipyardStage.Name); err != nil {
			return err, nil
		}
		pm.Logger.Info(fmt.Sprintf("Stage %s created", shipyardStage.Name))
	}
	pm.Logger.Info(fmt.Sprintf("Created all stages of project %s", *params.Name))

	uri := "shipyard.yaml"
	projectResource := []*apimodels.Resource{
		{
			ResourceContent: string(decodedShipyard),
			ResourceURI:     &uri,
		},
	}
	if err := pm.ConfigurationStore.CreateProjectShipyard(*params.Name, projectResource); err != nil {
		pm.Logger.Error(fmt.Sprintf("Error occurred while uploading shipyard resource to configuration service: %s", err.Error()))
		return err, func() error {
			pm.Logger.Info(fmt.Sprintf("Rollback: Try to delete project %s from configuration service", *params.Name))
			if err := pm.ConfigurationStore.DeleteProject(*params.Name); err != nil {
				pm.Logger.Error(fmt.Sprintf("Rollback failed: Unable to delete project %s from configuration service: %s", *params.Name, err.Error()))
				return err
			}
			pm.Logger.Info(fmt.Sprintf("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name))
			if err := pm.deleteGITRepositorySecret(*params.Name); err != nil {
				pm.Logger.Error(fmt.Sprintf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error()))
				return err
			}
			return nil
		}
	}

	if err := pm.createProjectInRepository(params, decodedShipyard, shipyard); err != nil {
		return err, func() error {
			if err := pm.ConfigurationStore.DeleteProject(*params.Name); err != nil {
				return err
			}
			if err := pm.deleteGITRepositorySecret(*params.Name); err != nil {
				return err
			}
			return nil
		}
	}
	return nil, nilRollback

}

func (pm *ProjectManager) Update(params *operations.UpdateProjectParams) (error, common.RollbackFunc) {
	oldSecret, err := pm.getGITRepositorySecret(*params.Name)
	if err != nil {
		return err, nilRollback
	}
	oldProject, err := pm.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil {
		return err, nilRollback
	}

	err = pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
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

	err = pm.ConfigurationStore.UpdateProject(projectToUpdate)
	if err != nil {
		return err, func() error {
			return pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
				User:      oldSecret.User,
				Token:     oldSecret.Token,
				RemoteURI: oldSecret.RemoteURI,
			})
		}
	}

	err = pm.ProjectMaterializedView.UpdateUpstreamInfo(*params.Name, params.GitRemoteURL, params.GitUser)
	if err != nil {
		return err, func() error {

			errConfigStoreRollback := pm.ConfigurationStore.UpdateProject(projectToRollback)
			if errConfigStoreRollback != nil {
				return errConfigStoreRollback
			}
			return pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
				User:      oldSecret.User,
				Token:     oldSecret.Token,
				RemoteURI: oldSecret.RemoteURI,
			})
		}
	}

	return nil, nilRollback
}

func (pm *ProjectManager) Delete(projectName string) (error, string) {
	pm.Logger.Info(fmt.Sprintf("Deleting project %s", projectName))
	var resultMessage strings.Builder

	project, err := pm.ProjectMaterializedView.GetProject(projectName)
	if err != nil {
		resultMessage.WriteString(fmt.Sprintf("Project %s cannot be retrieved anymore. Any Git upstream of the project will not be deleted.\n", projectName))
	} else if project != nil && project.GitRemoteURI != "" {
		resultMessage.WriteString(fmt.Sprintf("The Git upstream of the project will not be deleted: %s\n", project.GitRemoteURI))
	}

	secret, err := pm.SecretStore.GetSecret("git-credentials-" + projectName)
	if err != nil {
		pm.Logger.Error("could not delete git upstream credentials secret: " + err.Error())
	}
	if secret != nil {
		if err := pm.SecretStore.DeleteSecret("git-credentials-" + projectName); err != nil {
			pm.Logger.Error("could not delete git upstream credentials secret: " + err.Error())
			resultMessage.WriteString("WARNING: Could not delete secret containing the git upstream repo credentials. \n")
			resultMessage.WriteString(fmt.Sprintf("Please make sure to delete the secret manually by executing 'kubectl delete secret %s -n %s' \n", "git-credentials-"+projectName, common.GetKeptnNamespace()))
		}
	}

	if err := pm.ConfigurationStore.DeleteProject(projectName); err != nil {
		return pm.logAndReturnError(fmt.Sprintf("could not delete project: %s", err.Error())), resultMessage.String()
	}

	resultMessage.WriteString(pm.getDeleteInfoMessage(projectName))

	if err := pm.EventRpository.DeleteEventCollections(projectName); err != nil {
		pm.Logger.Error(fmt.Sprintf("could not delete task sequence collection: %s", err.Error()))
	}

	if err := pm.TaskSequenceRepository.DeleteTaskSequenceCollection(projectName); err != nil {
		pm.Logger.Error(fmt.Sprintf("could not delete task equence collection: %s", err.Error()))
	}

	if err := pm.ProjectMaterializedView.DeleteProject(projectName); err != nil {
		pm.Logger.Error(fmt.Sprintf("could not delete project: %s", err.Error()))
	}

	return nil, resultMessage.String()

}

func (pm *ProjectManager) createProjectInRepository(params *operations.CreateProjectParams, decodedShipyard []byte, shipyard *keptnv2.Shipyard) error {

	var expandedStages []*models.ExpandedStage

	for _, s := range shipyard.Spec.Stages {
		es := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: s.Name,
		}
		expandedStages = append(expandedStages, es)
	}

	p := &models.ExpandedProject{
		CreationDate:    strconv.FormatInt(time.Now().UnixNano(), 10),
		GitRemoteURI:    params.GitRemoteURL,
		GitUser:         params.GitUser,
		ProjectName:     *params.Name,
		Shipyard:        string(decodedShipyard),
		ShipyardVersion: shipyardVersion,
		Stages:          expandedStages,
	}

	err := pm.ProjectMaterializedView.CreateProject(p)
	if err != nil {
		return err
	}
	return nil
}

func (pm *ProjectManager) getGITRepositorySecret(projectName string) (*gitCredentials, error) {
	secret, err := pm.SecretStore.GetSecret("git-credentials-" + projectName)
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

func (pm *ProjectManager) updateGITRepositorySecret(projectName string, credentials *gitCredentials) error {

	credsEncoded, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	if err := pm.SecretStore.UpdateSecret("git-credentials-"+projectName, map[string][]byte{
		"git-credentials": credsEncoded,
	}); err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	return nil
}

func (pm *ProjectManager) deleteGITRepositorySecret(projectName string) error {
	pm.Logger.Info("deleting git credentials for project " + projectName)

	if err := pm.SecretStore.DeleteSecret("git-credentials-" + projectName); err != nil {
		return fmt.Errorf("could not delete git credentials: %s", err.Error())
	}
	pm.Logger.Info("deleted git credentials for project " + projectName)
	return nil

}

func (pm *ProjectManager) getDeleteInfoMessage(project string) string {
	res, err := pm.ConfigurationStore.GetProjectResource(project, "shipyard.yaml")
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

func (pm *ProjectManager) logAndReturnError(msg string) error {
	pm.Logger.Error(msg)
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

type gitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}
