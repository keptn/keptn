package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"time"
)

const shipyardVersion = "spec.keptn.sh/0.2.0"

//go:generate moq -pkg fake -skip-ensure -out ./fake/projectmanager.go . IProjectManager
type IProjectManager interface {
	Get() ([]*models.ExpandedProject, error)
	GetByName(projectName string) (*models.ExpandedProject, error)
	Create(params *models.CreateProjectParams) (error, common.RollbackFunc)
	Update(params *models.UpdateProjectParams) (error, common.RollbackFunc)
	Delete(projectName string) (string, error)
}

type ProjectManager struct {
	ConfigurationStore      common.ConfigurationStore
	SecretStore             common.SecretStore
	ProjectMaterializedView db.ProjectMVRepo
	TaskSequenceRepository  db.TaskSequenceRepo
	EventRepository         db.EventRepo
	SequenceQueueRepo       db.SequenceQueueRepo
	EventQueueRepo          db.EventQueueRepo
}

var nilRollback = func() error {
	return nil
}

func NewProjectManager(
	configurationStore common.ConfigurationStore,
	secretStore common.SecretStore,
	projectMVrepo db.ProjectMVRepo,
	taskSequenceRepo db.TaskSequenceRepo,
	eventRepo db.EventRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	eventQueueRepo db.EventQueueRepo) *ProjectManager {
	projectUpdater := &ProjectManager{
		ConfigurationStore:      configurationStore,
		SecretStore:             secretStore,
		ProjectMaterializedView: projectMVrepo,
		TaskSequenceRepository:  taskSequenceRepo,
		EventRepository:         eventRepo,
		SequenceQueueRepo:       sequenceQueueRepo,
		EventQueueRepo:          eventQueueRepo,
	}
	return projectUpdater
}

func (pm *ProjectManager) Get() ([]*models.ExpandedProject, error) {
	allProjects, err := pm.ProjectMaterializedView.GetProjects()
	if err != nil {
		return nil, err
	}
	return allProjects, nil
}

func (pm *ProjectManager) GetByName(projectName string) (*models.ExpandedProject, error) {
	project, err := pm.ProjectMaterializedView.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return project, err
}

func (pm *ProjectManager) Create(params *models.CreateProjectParams) (error, common.RollbackFunc) {

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

	rollbackFunc := func() error {
		log.Infof("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name)
		if err := pm.deleteGITRepositorySecret(*params.Name); err != nil {
			log.Errorf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error())
			return err
		}
		return nil
	}

	if err != nil {
		log.Errorf("Error occurred while creating project in configuration service: %s", err.Error())
		return err, rollbackFunc
	}
	log.Infof("Created project in configuration service: %s", *params.Name)

	// extend the rollback func to also delete the project in case anything goes wrong afterwards
	rollbackFunc = func() error {
		log.Infof("Rollback: Try to delete project %s from configuration service", *params.Name)
		if err := pm.ConfigurationStore.DeleteProject(*params.Name); err != nil {
			log.Errorf("Rollback failed: Unable to delete project %s from configuration service: %s", *params.Name, err.Error())
			return err
		}
		log.Infof("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name)
		if err := pm.deleteGITRepositorySecret(*params.Name); err != nil {
			log.Errorf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error())
			return err
		}
		return nil
	}

	decodedShipyard, _ := base64.StdEncoding.DecodeString(*params.Shipyard)
	shipyard, _ := common.UnmarshalShipyard(string(decodedShipyard))
	for _, shipyardStage := range shipyard.Spec.Stages {
		if err := pm.ConfigurationStore.CreateStage(*params.Name, shipyardStage.Name); err != nil {
			return err, rollbackFunc
		}
		log.Infof("Stage %s created", shipyardStage.Name)
	}
	log.Infof("Created all stages of project %s", *params.Name)

	uri := "shipyard.yaml"
	projectResource := []*apimodels.Resource{
		{
			ResourceContent: string(decodedShipyard),
			ResourceURI:     &uri,
		},
	}
	if err := pm.ConfigurationStore.CreateProjectShipyard(*params.Name, projectResource); err != nil {
		log.Errorf("Error occurred while uploading shipyard resource to configuration service: %s", err.Error())
		return err, rollbackFunc
	}

	if err := pm.createProjectInRepository(params, decodedShipyard, shipyard); err != nil {
		return err, rollbackFunc
	}

	// make sure mongodb collections from previous project with the same name are emptied
	pm.deleteProjectSequenceCollections(*params.Name)

	return nil, nilRollback
}

func (pm *ProjectManager) Update(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
	// old secret for rollback
	oldSecret, err := pm.getGITRepositorySecret(*params.Name)
	if err != nil {
		return err, nilRollback
	}

	// old project for rollback
	oldProject, err := pm.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil {
		return err, nilRollback
	} else if oldProject == nil {
		return ErrProjectNotFound, nilRollback
	}

	if params.GitUser != "" && params.GitToken != "" && params.GitRemoteURL != "" {
		// try to update git repository secret
		err = pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
			User:      params.GitUser,
			Token:     params.GitToken,
			RemoteURI: params.GitRemoteURL,
		})

		// no roll back needed since updating the git repository secret was the first operation
		if err != nil {
			return err, nilRollback
		}
	}

	// new project content in configuration service
	projectToUpdate := apimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	// project content in configuration service to rollback
	projectToRollback := apimodels.Project{
		CreationDate:    oldProject.CreationDate,
		GitRemoteURI:    oldProject.GitRemoteURI,
		GitUser:         oldProject.GitUser,
		ProjectName:     oldProject.ProjectName,
		ShipyardVersion: oldProject.ShipyardVersion,
	}

	// try to update the project information in configuration service
	err = pm.ConfigurationStore.UpdateProject(projectToUpdate)

	if err != nil {
		return err, func() error {
			// try to rollback already updated git repository secret
			if err := pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
				User:      oldSecret.User,
				Token:     oldSecret.Token,
				RemoteURI: oldSecret.RemoteURI,
			}); err != nil {
				return err
			}
			// try to rollback already updated project in configuration store
			return pm.ConfigurationStore.UpdateProject(projectToRollback)
		}
	}

	// try to update shipyard project resource
	if params.Shipyard != "" {
		shipyardResource := apimodels.Resource{
			ResourceContent: params.Shipyard,
			ResourceURI:     common.Stringp("shipyard.yaml"),
		}
		err = pm.ConfigurationStore.UpdateProjectResource(*params.Name, &shipyardResource)
		if err != nil {
			return err, func() error {
				// try to rollback already updated git repository secret
				if err = pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
					User:      oldSecret.User,
					Token:     oldSecret.Token,
					RemoteURI: oldSecret.RemoteURI}); err != nil {
					return err
				}
				// try to rollback already updated project in configuration store
				return pm.ConfigurationStore.UpdateProject(projectToRollback)
			}
		}
	}

	// copy by value
	updateProject := *oldProject
	updateProject.GitUser = params.GitUser
	updateProject.GitRemoteURI = params.GitRemoteURL
	if params.Shipyard != "" {
		updateProject.Shipyard = params.Shipyard
	}

	// try to update project information in database
	err = pm.ProjectMaterializedView.UpdateProject(&updateProject)
	if err != nil {
		return err, func() error {
			// try to rollback already updated project resource in configuration service
			if err = pm.ConfigurationStore.UpdateProjectResource(*params.Name, &apimodels.Resource{
				ResourceContent: oldProject.Shipyard,
				ResourceURI:     common.Stringp("shipyard.yaml")}); err != nil {
				return err
			}

			// try to rollback already updated project information in configuration service
			if err = pm.ConfigurationStore.UpdateProject(projectToRollback); err != nil {
				return err
			}

			// try to rollback already updated git repository secret
			return pm.updateGITRepositorySecret(*params.Name, &gitCredentials{
				User:      oldSecret.User,
				Token:     oldSecret.Token,
				RemoteURI: oldSecret.RemoteURI,
			})
		}
	}

	return nil, nilRollback
}

func (pm *ProjectManager) Delete(projectName string) (string, error) {
	log.Infof("Deleting project %s", projectName)
	var resultMessage strings.Builder

	project, err := pm.ProjectMaterializedView.GetProject(projectName)
	if err != nil {
		resultMessage.WriteString(fmt.Sprintf("Project %s cannot be retrieved anymore. Any Git upstream of the project will not be deleted.\n", projectName))
	} else if project != nil && project.GitRemoteURI != "" {
		resultMessage.WriteString(fmt.Sprintf("The Git upstream of the project will not be deleted: %s\n", project.GitRemoteURI))
	}

	secret, err := pm.SecretStore.GetSecret("git-credentials-" + projectName)
	if err != nil {
		log.Errorf("could not delete git upstream credentials secret: %s", err.Error())
	}
	if secret != nil {
		if err := pm.SecretStore.DeleteSecret("git-credentials-" + projectName); err != nil {
			log.Errorf("could not delete git upstream credentials secret: %s", err.Error())
			resultMessage.WriteString("WARNING: Could not delete secret containing the git upstream repo credentials. \n")
			resultMessage.WriteString(fmt.Sprintf("Please make sure to delete the secret manually by executing 'kubectl delete secret %s -n %s' \n", "git-credentials-"+projectName, common.GetKeptnNamespace()))
		}
	}

	if err := pm.ConfigurationStore.DeleteProject(projectName); err != nil {
		return resultMessage.String(), pm.logAndReturnError(fmt.Sprintf("could not delete project: %s", err.Error()))
	}

	resultMessage.WriteString(pm.getDeleteInfoMessage(projectName))

	if err := pm.ProjectMaterializedView.DeleteProject(projectName); err != nil {
		log.Errorf("could not delete project: %s", err.Error())
	}

	pm.deleteProjectSequenceCollections(projectName)

	return resultMessage.String(), nil
}

func (pm *ProjectManager) deleteProjectSequenceCollections(projectName string) {
	if err := pm.EventRepository.DeleteEventCollections(projectName); err != nil {
		log.Errorf("could not delete task sequence collection: %s", err.Error())
	}

	if err := pm.TaskSequenceRepository.DeleteRepo(projectName); err != nil {
		log.Errorf("could not delete task sequence collection: %s", err.Error())
	}

	if err := pm.SequenceQueueRepo.DeleteQueuedSequences(models.QueueItem{Scope: models.EventScope{
		EventData: keptnv2.EventData{
			Project: projectName,
		},
	}}); err != nil {
		log.Errorf("could not delete queued sequences: %s", err.Error())
	}

	if err := pm.EventQueueRepo.DeleteQueuedEvents(models.EventScope{
		EventData: keptnv2.EventData{
			Project: projectName,
		},
	}); err != nil {
		log.Errorf("could not delete queued events: %s", err.Error())
	}
}

func (pm *ProjectManager) createProjectInRepository(params *models.CreateProjectParams, decodedShipyard []byte, shipyard *keptnv2.Shipyard) error {

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
	log.Infof("deleting git credentials for project %s", projectName)

	if err := pm.SecretStore.DeleteSecret("git-credentials-" + projectName); err != nil {
		return fmt.Errorf("could not delete git credentials: %s", err.Error())
	}
	log.Infof("deleted git credentials for project %s", projectName)
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
	log.Error(msg)
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
