package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"strings"
)

var errProjectAlreadyExists = errors.New("project already exists")

func newProjectManager() (*projectManager, error) {
	base, err := newAPIBase()
	if err != nil {
		return nil, err
	}
	return &projectManager{
		apiBase: base,
		eventRepo: &db.MongoDBEventsRepo{
			Logger: base.logger,
		},
		taskSequenceRepo: &db.TaskSequenceMongoDBRepo{
			Logger: base.logger,
		},
	}, nil
}

type projectManager struct {
	*apiBase
	eventRepo        db.EventRepo
	taskSequenceRepo db.TaskSequenceRepo
}

type gitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

func (pm *projectManager) deleteProject(projectName string) (*operations.DeleteProjectResponse, error) {
	pm.logger.Info("Deleting project " + projectName)
	result := &operations.DeleteProjectResponse{
		Message: "",
	}

	projectToDelete := keptnapimodels.Project{
		ProjectName: projectName,
	}

	project, errObj := pm.projectAPI.GetProject(projectToDelete)
	if errObj != nil {
		result.Message = result.Message + fmt.Sprintf("Project %s cannot be retrieved anymore. Any Git upstream of the project will not be deleted.\n", projectName)
	} else if project != nil && project.GitRemoteURI != "" {
		result.Message = result.Message + fmt.Sprintf("The Git upstream of the project will not be deleted: %s\n", project.GitRemoteURI)
	}

	// check for an upstream repo secret
	secret, err := pm.secretStore.GetSecret(getUpstreamRepoCredsSecretName(projectName))
	if err != nil {
		pm.logger.Error("could not delete git upstream credentials secret: " + err.Error())
	}
	if secret != nil {
		// try to delete the secret
		if err := pm.secretStore.DeleteSecret(getUpstreamRepoCredsSecretName(projectName)); err != nil {
			// if anything goes wrong, log the error, but continue with deleting the remaining project resources
			pm.logger.Error("could not delete git upstream credentials secret: " + err.Error())
			result.Message = result.Message + "WARNING: Could not delete secret containing the git upstream repo credentials. \n"
			result.Message = result.Message + fmt.Sprintf("Please make sure to delete the secret manually by executing 'kubectl delete secret %s -n %s' \n", getUpstreamRepoCredsSecretName(projectName), common.GetKeptnNamespace())
		}
	}

	if _, errObj := pm.projectAPI.DeleteProject(keptnapimodels.Project{
		ProjectName: projectName,
	}); errObj != nil {
		return nil, pm.logAndReturnError(fmt.Sprintf("could not delete project: %s", *errObj.Message))
	}

	result.Message = result.Message + pm.getDeleteInfoMessage(projectName)

	if err := pm.taskSequenceRepo.DeleteTaskSequenceCollection(projectName); err != nil {
		pm.logger.Error("could not delete task sequence collection: " + err.Error())
	}

	if err := pm.eventRepo.DeleteEventCollections(projectName); err != nil {
		pm.logger.Error("could not delete event collections: " + err.Error())
	}

	if err := pm.sendProjectDeleteSuccessFinishedEvent(uuid.New().String(), projectName); err != nil {
		pm.logger.Error("could not send project.delete.finished event: " + err.Error())
	}

	return result, nil
}

func getShipyardNotAvailableError(project string) string {
	return fmt.Sprintf("Shipyard of project %s cannot be retrieved anymore. "+
		"After deleting the project, the namespaces containing the services are still available. "+
		"This may cause problems if a project with the same name is created later.", project)
}

func (pm *projectManager) getDeleteInfoMessage(project string) string {
	res, err := pm.resourceAPI.GetProjectResource(project, "shipyard.yaml")
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

func (pm *projectManager) updateProject(params *operations.CreateProjectParams) error {
	pm.logger.Info(fmt.Sprintf("checking if project %s exists before updating it", *params.Name))
	project, err := pm.projectAPI.GetProject(keptnapimodels.Project{
		ProjectName: *params.Name,
	})
	if err != nil {
		msg := fmt.Sprintf("Could not check if project %s exists; %s", *params.Name, *err.Message)
		pm.logger.Error(msg)
		return errors.New(msg)
	}
	if project == nil {
		msg := fmt.Sprintf("Project %s does not exist", *params.Name)
		pm.logger.Error(msg)
		return errors.New(msg)
	}
	oldSecret, getSecretErr := pm.getUpstreamRepoCredentials(*params.Name)
	if getSecretErr != nil {
		// log the error but continue
		pm.logger.Error(fmt.Sprintf("could not read previous secret of project %s: %s", *params.Name, getSecretErr.Error()))
	}

	if params.GitRemoteURL != "" && params.GitUser != "" && params.GitToken != "" {
		if err := pm.createUpstreamRepoCredentials(params); err != nil {
			return pm.logAndReturnError(err.Error())
		}
	}

	_, errObj := pm.projectAPI.UpdateConfigurationServiceProject(keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	})

	if errObj != nil {
		msg := fmt.Sprintf("Could not update upstream repository of project %s: %s", *params.Name, *errObj.Message)

		if oldSecret != nil {
			// restore previous secret
			if createErr := pm.createUpstreamRepoCredentials(&operations.CreateProjectParams{
				GitRemoteURL: oldSecret.RemoteURI,
				GitToken:     oldSecret.Token,
				GitUser:      oldSecret.User,
				Name:         params.Name,
			}); createErr != nil {
				pm.logger.Error(fmt.Sprintf("Could not restore previous upstream repo credentials: %s", createErr.Error()))
			} else {
				// restore the upstream on the configuration service
				if _, restoreErrObj := pm.projectAPI.UpdateConfigurationServiceProject(keptnapimodels.Project{
					GitRemoteURI: oldSecret.RemoteURI,
					GitToken:     oldSecret.Token,
					GitUser:      oldSecret.User,
					ProjectName:  *params.Name,
				}); restoreErrObj != nil {
					pm.logger.Error(fmt.Sprintf("Could not restore previous upstream on configuration service: %s", *restoreErrObj.Message))
				}
			}
		} else {
			if delErr := pm.deleteUpstreamRepoCredentials(params); delErr != nil {
				pm.logger.Error(fmt.Sprintf("Could not delete upstream repo credentials: %s", delErr.Error()))
			}
		}
		return pm.logAndReturnError(msg)
	}
	return nil
}

func (pm *projectManager) createProject(params *operations.CreateProjectParams) (bool, error) {
	secretCreated := false
	keptnContext := uuid.New().String()
	// check if the project already exists
	pm.logger.Info(fmt.Sprintf("checking if project %s already exists before creating it", *params.Name))
	project, _ := pm.projectAPI.GetProject(keptnapimodels.Project{
		ProjectName: *params.Name,
	})
	if project != nil {
		pm.logger.Info(fmt.Sprintf("Project %s already exists", *params.Name))
		return secretCreated, errProjectAlreadyExists
	}

	if err := pm.sendProjectCreateStartedEvent(keptnContext, params); err != nil {
		return secretCreated, pm.logAndReturnError(err.Error())
	}

	// if available, create the upstream repository credentials secret.
	// this has to be done before creating the project on the configuration service
	if params.GitRemoteURL != "" && params.GitUser != "" && params.GitToken != "" {
		pm.logger.Info(fmt.Sprintf("Storing upstream repo credentials for project %s", *params.Name))
		if err := pm.createUpstreamRepoCredentials(params); err != nil {
			return secretCreated, pm.logAndReturnError(err.Error())
		}
		pm.logger.Info(fmt.Sprintf("Successfully stored upstream repo credentials for project %s", *params.Name))
		secretCreated = true
	}

	// create the project
	_, errObj := pm.projectAPI.CreateProject(keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	})

	if errObj != nil {
		return secretCreated, pm.logAndReturnError(fmt.Sprintf("could not create project: %s", *errObj.Message))
	}
	pm.logger.Info(fmt.Sprintf("Project %s created", *params.Name))

	decodedShipyard, err := base64.StdEncoding.DecodeString(*params.Shipyard)
	if err != nil {
		// error should not occur at this stage because the shipyard content has been validated at this stage, but let's check anyways
		return secretCreated, pm.logAndReturnError(fmt.Sprintf("could not decode shipyard: " + err.Error()))
	}
	shipyard, err := common.UnmarshalShipyard(string(decodedShipyard))

	// create the stages
	for _, shipyardStage := range shipyard.Spec.Stages {
		if _, errorObj := pm.stagesAPI.CreateStage(*params.Name, shipyardStage.Name); err != nil {
			return secretCreated, pm.logAndReturnError(fmt.Sprintf("Failed to create stage %s: %s", shipyardStage.Name, *errorObj.Message))
		}
		pm.logger.Info(fmt.Sprintf("Stage %s created", shipyardStage.Name))
	}
	pm.logger.Info("created all stages of project " + *params.Name)

	// upload the shipyard file
	uri := "shipyard.yaml"
	_, err = pm.resourceAPI.CreateProjectResources(*params.Name, []*keptnapimodels.Resource{
		{
			ResourceContent: string(decodedShipyard),
			ResourceURI:     &uri,
		},
	})

	if err != nil {
		return secretCreated, pm.logAndReturnError(fmt.Sprintf("could not upload shipyard.yaml: %s", err.Error()))
	}
	pm.logger.Info("uploaded shipyard.yaml of project " + *params.Name)

	if err := pm.sendProjectCreateSuccessFinishedEvent(keptnContext, params); err != nil {
		return secretCreated, pm.logAndReturnError(err.Error())
	}
	return secretCreated, nil
}

func (pm *projectManager) sendProjectCreateStartedEvent(keptnContext string, params *operations.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), eventPayload); err != nil {
		return errors.New("could not send create.project.started event: " + err.Error())
	}

	return nil
}

func (pm *projectManager) sendProjectDeleteSuccessFinishedEvent(keptnContext, projectName string) error {
	eventPayload := keptnv2.ProjectDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName), eventPayload); err != nil {
		return errors.New("could not send create.project.finished event: " + err.Error())
	}
	return nil
}

func (pm *projectManager) sendProjectCreateSuccessFinishedEvent(keptnContext string, params *operations.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
		CreatedProject: keptnv2.ProjectCreateData{
			ProjectName:  *params.Name,
			GitRemoteURL: params.GitRemoteURL,
			Shipyard:     *params.Shipyard,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), eventPayload); err != nil {
		return errors.New("could not send create.project.finished event: " + err.Error())
	}
	return nil
}

func (pm *projectManager) getUpstreamRepoCredentials(projectName string) (*gitCredentials, error) {
	secret, err := pm.secretStore.GetSecret(getUpstreamRepoCredsSecretName(projectName))
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

func (pm *projectManager) createUpstreamRepoCredentials(params *operations.CreateProjectParams) error {
	pm.logger.Info("Storing git credentials for project " + *params.Name)
	credentials := &gitCredentials{
		User:      params.GitUser,
		Token:     params.GitToken,
		RemoteURI: params.GitRemoteURL,
	}

	credsEncoded, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	if err := pm.secretStore.UpdateSecret(getUpstreamRepoCredsSecretName(*params.Name), map[string][]byte{
		"git-credentials": credsEncoded,
	}); err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	pm.logger.Info("stored git credentials for project " + *params.Name)
	return nil
}

func (pm *projectManager) deleteUpstreamRepoCredentials(params *operations.CreateProjectParams) error {
	pm.logger.Info("Deleting git credentials for project " + *params.Name)

	if err := pm.secretStore.DeleteSecret(getUpstreamRepoCredsSecretName(*params.Name)); err != nil {
		return fmt.Errorf("could not delete git credentials: %s", err.Error())
	}
	pm.logger.Info("deleted git credentials for project " + *params.Name)
	return nil
}

func getUpstreamRepoCredsSecretName(projectName string) string {
	return "git-credentials-" + projectName
}

func (pm *projectManager) logAndReturnError(msg string) error {
	pm.logger.Error(msg)
	return errors.New(msg)
}
