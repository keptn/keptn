package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/configurationstore"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"github.com/keptn/keptn/shipyard-controller/internal/secretstore"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const shipyardVersion = "spec.keptn.sh/0.2.0"
const errUpdateProject = "failed to update project '%s': %w"

//go:generate moq -pkg fake -skip-ensure -out ./fake/projectmanager.go . IProjectManager
type IProjectManager interface {
	Get() ([]*apimodels.ExpandedProject, error)
	GetByName(projectName string) (*apimodels.ExpandedProject, error)
	Create(params *models.CreateProjectParams, internalOptions models.InternalCreateProjectOptions) (error, common.RollbackFunc)
	Update(params *models.UpdateProjectParams) (error, common.RollbackFunc)
	Delete(projectName string) (string, error)
}

func WithHideAutoProvisionedURL(hideAutoProvisionedURL bool) func(pm *ProjectManager) {
	return func(pm *ProjectManager) {
		pm.hideAutoProvisionedURL = hideAutoProvisionedURL
	}
}

type ProjectManager struct {
	ConfigurationStore      configurationstore.ConfigurationStore
	SecretStore             secretstore.SecretStore
	ProjectMaterializedView db.ProjectMVRepo
	SequenceExecutionRepo   db.SequenceExecutionRepo
	EventRepository         db.EventRepo
	SequenceQueueRepo       db.SequenceQueueRepo
	EventQueueRepo          db.EventQueueRepo
	hideAutoProvisionedURL  bool
}

var nilRollback = func() error {
	return nil
}

func NewProjectManager(
	configurationStore configurationstore.ConfigurationStore,
	secretStore secretstore.SecretStore,
	projectMVrepo db.ProjectMVRepo,
	sequenceExecutionRepo db.SequenceExecutionRepo,
	eventRepo db.EventRepo,
	sequenceQueueRepo db.SequenceQueueRepo,
	eventQueueRepo db.EventQueueRepo,
	opts ...func(pm *ProjectManager)) *ProjectManager {
	projectUpdater := &ProjectManager{
		ConfigurationStore:      configurationStore,
		SecretStore:             secretStore,
		ProjectMaterializedView: projectMVrepo,
		SequenceExecutionRepo:   sequenceExecutionRepo,
		EventRepository:         eventRepo,
		SequenceQueueRepo:       sequenceQueueRepo,
		EventQueueRepo:          eventQueueRepo,
	}

	for _, o := range opts {
		o(projectUpdater)
	}
	return projectUpdater
}

func (pm *ProjectManager) Get() ([]*apimodels.ExpandedProject, error) {
	allProjects, err := pm.ProjectMaterializedView.GetProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range allProjects {
		pm.modifyProjectResponse(project)
	}

	return allProjects, nil
}

func (pm *ProjectManager) GetByName(projectName string) (*apimodels.ExpandedProject, error) {
	project, err := pm.ProjectMaterializedView.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, common.ErrProjectNotFound
	}

	pm.modifyProjectResponse(project)

	return project, err
}

func (pm *ProjectManager) modifyProjectResponse(project *apimodels.ExpandedProject) {
	if project.IsUpstreamAutoProvisioned && pm.hideAutoProvisionedURL {
		project.GitCredentials.User = ""
		project.GitCredentials.RemoteURL = ""
	}
}

func (pm *ProjectManager) Create(params *models.CreateProjectParams, options models.InternalCreateProjectOptions) (error, common.RollbackFunc) {

	if err := pm.checkForExistingProject(params); err != nil {
		return fmt.Errorf("could not create project '%s': %w", *params.Name, err), nilRollback
	}

	decodedCredentials, err := decodeGitCredentials(*params.GitCredentials)
	if err != nil {
		return fmt.Errorf("could not create project '%s': %w", *params.Name, err), nilRollback
	}
	err = pm.updateGITRepositorySecret(getUpstreamCredentialSecretName(*params.Name), decodedCredentials)
	if err != nil {
		return err, nilRollback
	}

	err = pm.ConfigurationStore.CreateProject(apimodels.Project{
		ProjectName: *params.Name,
	})

	rollbackFunc := func() error {
		log.Infof("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name)
		if err := pm.deleteGITRepositorySecret(getUpstreamCredentialSecretName(*params.Name)); err != nil {
			log.Errorf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error())
			return common.ErrChangesRollback
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
			return common.ErrChangesRollback
		}
		log.Infof("Rollback: Try to delete GIT repository credentials secret for project %s", *params.Name)
		if err := pm.deleteGITRepositorySecret(getUpstreamCredentialSecretName(*params.Name)); err != nil {
			log.Errorf("Rollback failed: Unable to delete GIT repository credentials secret for project %s: %s", *params.Name, err.Error())
			return common.ErrChangesRollback
		}
		return nil
	}

	decodedShipyard, _ := base64.StdEncoding.DecodeString(*params.Shipyard)
	shipyard, _ := common.UnmarshalShipyard(string(decodedShipyard))
	for _, shipyardStage := range shipyard.Spec.Stages {
		if err := pm.ConfigurationStore.CreateStage(*params.Name, shipyardStage.Name); err != nil {
			return fmt.Errorf("failed to create stage '%s' for project '%s'", shipyardStage.Name, *params.Name), rollbackFunc
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
		return fmt.Errorf("failed to upload shipyard resource for project '%s'", *params.Name), rollbackFunc
	}

	if err := pm.createProjectInRepository(params, decodedShipyard, shipyard, options); err != nil {
		log.Errorf("Error occurred creating project in respository: %s", err.Error())
		return fmt.Errorf("failed to create project '%s'", *params.Name), rollbackFunc
	}

	// make sure mongodb collections from previous project with the same name are emptied
	pm.deleteProjectSequenceCollections(*params.Name)

	return nil, nilRollback
}

func (pm *ProjectManager) checkForExistingProject(params *models.CreateProjectParams) error {
	existingProject, err := pm.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil && err != common.ErrProjectNotFound {
		log.Errorf("Error occurred while getting project: %s", err.Error())
		return fmt.Errorf("failed to get information for project '%s': %w", *params.Name, err)
	}
	if existingProject != nil {
		return common.ErrProjectAlreadyExists
	}
	return nil
}

func (pm *ProjectManager) Update(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
	// old secret for rollback
	oldSecret, err := pm.getGITRepositorySecret(*params.Name)
	if err != nil {
		return err, nilRollback
	}

	rollbackSecretCredentials := oldSecret

	// old project for rollback
	oldProject, err := pm.ProjectMaterializedView.GetProject(*params.Name)
	if err != nil {
		log.Errorf("Error occurred while getting project: %s", err.Error())
		return fmt.Errorf("failed to get project: '%s'", *params.Name), nilRollback
	} else if oldProject == nil {
		return common.ErrProjectNotFound, nilRollback
	}

	if params.GitCredentials != nil {
		decodedCredentials, err := decodeGitCredentials(*params.GitCredentials)
		if err != nil {
			return fmt.Errorf("could not update project '%s': %w", *params.Name, err), nilRollback
		}
		// create a temporary secret containing the new git upstream credentials
		err = pm.updateGITRepositorySecret(getTemporaryUpstreamCredentialSecretName(*params.Name), decodedCredentials)

		// no roll back needed since updating the git repository secret was the first operation
		if err != nil {
			return err, nilRollback
		}
	}

	// new project content in resource service
	projectToUpdate := apimodels.Project{
		ProjectName: *params.Name,
	}

	// project content in resource service to rollback
	projectToRollback := apimodels.Project{
		CreationDate:    oldProject.CreationDate,
		ProjectName:     oldProject.ProjectName,
		ShipyardVersion: oldProject.ShipyardVersion,
		GitCredentials:  toInsecureGitCredentials(oldProject.GitCredentials),
	}

	// try to update the project information in configuration service
	err = pm.ConfigurationStore.UpdateProject(projectToUpdate)

	if err != nil {
		log.Errorf("Error occurred while updating the project in configuration store: %s", err.Error())
		return fmt.Errorf(errUpdateProject, projectToUpdate.ProjectName, err), func() error {
			// try to delete the temporary git repository secret holding the new credentials
			if err := pm.deleteGITRepositorySecret(getTemporaryUpstreamCredentialSecretName(*params.Name)); err != nil {
				return common.ErrChangesRollback
			}
			// try to rollback already updated project in configuration store
			return pm.ConfigurationStore.UpdateProject(projectToRollback)
		}
	}

	// if the update was successful, replace the previous git upstream credentials with the new ones
	if params.GitCredentials != nil {
		decodedCredentials, err := decodeGitCredentials(*params.GitCredentials)
		if err != nil {
			return fmt.Errorf("could not update project '%s': %w", *params.Name, err), nilRollback
		}
		// try to update git repository secret
		err = pm.updateGITRepositorySecret(getUpstreamCredentialSecretName(*params.Name), decodedCredentials)
		if err != nil {
			log.Errorf("Error occurred while updating the project in credentials secret: %s", err.Error())
			return fmt.Errorf(errUpdateProject, projectToUpdate.ProjectName, err), func() error {
				// try to rollback already updated project in configuration store
				return pm.ConfigurationStore.UpdateProject(projectToRollback)
			}
		}
		tmpSecretName := getTemporaryUpstreamCredentialSecretName(*params.Name)
		if err := pm.deleteGITRepositorySecret(tmpSecretName); err != nil {
			// log the error in this case, but continue, as deleting the temporary secret is not a blocker for updating the project credentials
			log.Errorf("Could not delete temporary secret %s: %v", tmpSecretName, err)
		}
	}

	var isShipyardPresent = params.Shipyard != nil && *params.Shipyard != ""

	// try to update shipyard project resource
	if isShipyardPresent {
		if err = validateShipyardUpdate(params, oldProject); err != nil {
			return err, nilRollback
		}

		shipyardResource := apimodels.Resource{
			ResourceContent: *params.Shipyard,
			ResourceURI:     common.Stringp("shipyard.yaml"),
		}
		err = pm.ConfigurationStore.UpdateProjectResource(*params.Name, &shipyardResource)
		if err != nil {
			log.Errorf("Error occurred while updating the project shipyard in configuration store: %s", err.Error())
			return fmt.Errorf(errUpdateProject, projectToUpdate.ProjectName, err), func() error {
				// try to rollback already updated git repository secret
				if err = pm.updateGITRepositorySecret(getUpstreamCredentialSecretName(*params.Name), rollbackSecretCredentials); err != nil {
					return common.ErrChangesRollback
				}
				// try to rollback already updated project in configuration store
				return pm.ConfigurationStore.UpdateProject(projectToRollback)
			}
		}
	}

	// copy by value

	updateProject := *oldProject
	if params.GitCredentials != nil {
		updateProject.GitCredentials = toSecureGitCredentials(params.GitCredentials)
	}
	if params.GitCredentials != nil && params.GitCredentials.RemoteURL != "" {
		updateProject.IsUpstreamAutoProvisioned = false
	}
	if isShipyardPresent {
		updateProject.Shipyard = *params.Shipyard
	}

	// try to update project information in database
	err = pm.ProjectMaterializedView.UpdateProject(&updateProject)
	if err != nil {
		log.Errorf("Error occurred while updating the project in materialized view: %s", err.Error())
		return fmt.Errorf(errUpdateProject, projectToUpdate.ProjectName, err), func() error {
			// try to rollback already updated project resource in configuration service
			if err = pm.ConfigurationStore.UpdateProjectResource(*params.Name, &apimodels.Resource{
				ResourceContent: oldProject.Shipyard,
				ResourceURI:     common.Stringp("shipyard.yaml")}); err != nil {
				return common.ErrChangesRollback
			}

			// try to rollback already updated project information in configuration service
			if err = pm.ConfigurationStore.UpdateProject(projectToRollback); err != nil {
				return common.ErrChangesRollback
			}

			// try to rollback already updated git repository secret
			return pm.updateGITRepositorySecret(getUpstreamCredentialSecretName(*params.Name), rollbackSecretCredentials)
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
	} else if project != nil && project.GitCredentials != nil {
		resultMessage.WriteString(fmt.Sprintf("The Git upstream of the project will not be deleted: %s\n", project.GitCredentials.RemoteURL))
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

	resultMessage.WriteString(pm.getDeleteInfoMessage(projectName))

	//  clean up  database
	if err := pm.ProjectMaterializedView.DeleteProject(projectName); err != nil {
		log.Errorf("could not delete project: %s", err.Error())
	}
	pm.deleteProjectSequenceCollections(projectName)

	// attempt deleting from local git
	if err := pm.ConfigurationStore.DeleteProject(projectName); err != nil {
		return resultMessage.String(), pm.logAndReturnError(fmt.Sprintf("could not delete project: %s", err.Error()))
	}
	return resultMessage.String(), nil
}

func (pm *ProjectManager) deleteProjectSequenceCollections(projectName string) {
	if err := pm.EventRepository.DeleteEventCollections(projectName); err != nil {
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

	if err := pm.SequenceExecutionRepo.Clear(projectName); err != nil {
		log.Errorf("could not delete sequence executions: %s", err.Error())
	}
}

func (pm *ProjectManager) createProjectInRepository(params *models.CreateProjectParams, decodedShipyard []byte, shipyard *keptnv2.Shipyard, options models.InternalCreateProjectOptions) error {

	var expandedStages []*apimodels.ExpandedStage

	for _, s := range shipyard.Spec.Stages {
		es := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: s.Name,
		}
		expandedStages = append(expandedStages, es)
	}

	p := &apimodels.ExpandedProject{
		CreationDate:              strconv.FormatInt(time.Now().UnixNano(), 10),
		ProjectName:               *params.Name,
		Shipyard:                  string(decodedShipyard),
		ShipyardVersion:           shipyardVersion,
		GitCredentials:            toSecureGitCredentials(params.GitCredentials),
		Stages:                    expandedStages,
		IsUpstreamAutoProvisioned: options.IsUpstreamAutoProvisioned,
	}

	err := pm.ProjectMaterializedView.CreateProject(p)
	if err != nil {
		return err
	}
	return nil
}

func (pm *ProjectManager) getGITRepositorySecret(projectName string) (*apimodels.GitAuthCredentials, error) {
	secret, err := pm.SecretStore.GetSecret("git-credentials-" + projectName)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	if marshalledSecret, ok := secret["git-credentials"]; ok {
		secretObj := &apimodels.GitAuthCredentials{}
		if err := json.Unmarshal(marshalledSecret, secretObj); err != nil {
			return nil, fmt.Errorf("failed to unmarshal git-credentials secret")
		}
		return secretObj, nil
	}
	return nil, nil
}

func (pm *ProjectManager) updateGITRepositorySecret(secretName string, credentials *apimodels.GitAuthCredentials) error {

	credsEncoded, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	if err := pm.SecretStore.UpdateSecret(secretName, map[string][]byte{
		"git-credentials": credsEncoded,
	}); err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	return nil
}

func (pm *ProjectManager) deleteGITRepositorySecret(secretName string) error {
	log.Infof("deleting git credentials for project %s", secretName)

	if err := pm.SecretStore.DeleteSecret(secretName); err != nil {
		return fmt.Errorf("could not delete git credentials: %s", err.Error())
	}
	log.Infof("deleted git credentials for project %s", secretName)
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

func toModelProject(project apimodels.ExpandedProject) apimodels.Project {
	return apimodels.Project{
		CreationDate:    project.CreationDate,
		GitCredentials:  toInsecureGitCredentials(project.GitCredentials),
		ProjectName:     project.ProjectName,
		ShipyardVersion: project.ShipyardVersion,
	}
}

func toSecureGitCredentials(credentials *apimodels.GitAuthCredentials) *apimodels.GitAuthCredentialsSecure {
	if credentials == nil {
		return nil
	}
	secureCredentials := apimodels.GitAuthCredentialsSecure{
		User:      credentials.User,
		RemoteURL: credentials.RemoteURL,
	}

	if credentials.HttpsAuth != nil {
		httpAuth := apimodels.HttpsGitAuthSecure{
			InsecureSkipTLS: credentials.HttpsAuth.InsecureSkipTLS,
		}
		if credentials.HttpsAuth.Proxy != nil {
			httpAuth.Proxy = &apimodels.ProxyGitAuthSecure{
				URL:    credentials.HttpsAuth.Proxy.URL,
				Scheme: credentials.HttpsAuth.Proxy.Scheme,
				User:   credentials.HttpsAuth.Proxy.User,
			}
		}

		secureCredentials.HttpsAuth = &httpAuth
	}

	return &secureCredentials
}

func toInsecureGitCredentials(credentials *apimodels.GitAuthCredentialsSecure) *apimodels.GitAuthCredentials {
	if credentials == nil {
		return nil
	}

	insecureCredentials := apimodels.GitAuthCredentials{
		User:      credentials.User,
		RemoteURL: credentials.RemoteURL,
	}

	if credentials.HttpsAuth != nil {
		httpAuth := apimodels.HttpsGitAuth{
			InsecureSkipTLS: credentials.HttpsAuth.InsecureSkipTLS,
		}
		if credentials.HttpsAuth.Proxy != nil {
			httpAuth.Proxy = &apimodels.ProxyGitAuth{
				URL:    credentials.HttpsAuth.Proxy.URL,
				Scheme: credentials.HttpsAuth.Proxy.Scheme,
				User:   credentials.HttpsAuth.Proxy.User,
			}
		}

		insecureCredentials.HttpsAuth = &httpAuth
	}

	return &insecureCredentials

}

func validateShipyardStagesUnchanged(oldProject *apimodels.ExpandedProject, newProject *apimodels.ExpandedProject) error {
	if len(newProject.Stages) != len(oldProject.Stages) {
		return fmt.Errorf("unallowed addition/removal of project stages")
	}

	for i, oldStage := range oldProject.Stages {
		// It is more effective to check the names of the stages in two steps.
		// In typical user scenario, the user probably won't want to change the order
		// of the stages, at least it is unlikely. In most of the cases, he will try
		// to edit the name of the stage.
		// Let's consider a check, where the user did not changed the
		// names and the number of stages. If the first condition was not there, for each stage
		// in oldProject.Stages the code needs to jump to another function and cycle through the stages
		// of newProject.Stages -> N/2 string comparisons (assuming N is number of stages)
		// for checking each stage, so in total N*N/2 comparisons.
		// If the condition is there, we will have only 1 comparison for each stage,
		// total N*1 comparisons.
		if oldStage.StageName != newProject.Stages[i].StageName {
			if !stageInArrayOfStages(oldStage.StageName, newProject.Stages) {
				return fmt.Errorf("unallowed rename of project stages")
			}
		}
	}

	return nil
}

func stageInArrayOfStages(comparedStage string, stages []*apimodels.ExpandedStage) bool {
	for _, arrayStage := range stages {
		if arrayStage.StageName == comparedStage {
			return true
		}
	}
	return false
}

func decodeGitCredentials(oldCredentials apimodels.GitAuthCredentials) (*apimodels.GitAuthCredentials, error) {
	credentials := &apimodels.GitAuthCredentials{
		RemoteURL: oldCredentials.RemoteURL,
		User:      oldCredentials.User,
	}
	if oldCredentials.HttpsAuth != nil {
		httpsAuth, err := decodeHttpsAuth(*oldCredentials.HttpsAuth)
		if err != nil {
			return nil, err
		}

		credentials.HttpsAuth = httpsAuth
	}

	if oldCredentials.SshAuth != nil {
		sshAuth, err := decodeSshAuth(*oldCredentials.SshAuth)
		if err != nil {
			return nil, err
		}
		credentials.SshAuth = sshAuth
	}

	return credentials, nil
}

func decodeHttpsAuth(in apimodels.HttpsGitAuth) (*apimodels.HttpsGitAuth, error) {
	httpsAuth := &apimodels.HttpsGitAuth{
		Token:           in.Token,
		InsecureSkipTLS: in.InsecureSkipTLS,
	}
	if in.Certificate != "" {
		decodedPemCertificate, err := base64.StdEncoding.DecodeString(in.Certificate)
		if err != nil {
			return nil, err
		}
		httpsAuth.Certificate = string(decodedPemCertificate)
	}
	if in.Proxy != nil {
		httpsAuth.Proxy = &apimodels.ProxyGitAuth{
			URL:      in.Proxy.URL,
			Scheme:   in.Proxy.Scheme,
			User:     in.Proxy.User,
			Password: in.Proxy.Password,
		}
	}
	return httpsAuth, nil
}

func decodeSshAuth(in apimodels.SshGitAuth) (*apimodels.SshGitAuth, error) {
	sshAuth := &apimodels.SshGitAuth{
		PrivateKeyPass: in.PrivateKeyPass,
	}
	if in.PrivateKey != "" {
		decodedPrivateKey, err := base64.StdEncoding.DecodeString(in.PrivateKey)
		if err != nil {
			return nil, err
		}
		sshAuth.PrivateKey = string(decodedPrivateKey)
	}
	return sshAuth, nil
}

func validateShipyardUpdate(params *models.UpdateProjectParams, oldProject *apimodels.ExpandedProject) error {
	shipyard := &keptnv2.Shipyard{}
	decodedShipyard, _ := base64.StdEncoding.DecodeString(*params.Shipyard)
	_ = yaml.Unmarshal([]byte(decodedShipyard), shipyard)
	var expandedStages []*apimodels.ExpandedStage

	for _, s := range shipyard.Spec.Stages {
		es := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: s.Name,
		}
		expandedStages = append(expandedStages, es)
	}

	newProject := &apimodels.ExpandedProject{
		CreationDate:    strconv.FormatInt(time.Now().UnixNano(), 10),
		GitCredentials:  toSecureGitCredentials(params.GitCredentials),
		ProjectName:     *params.Name,
		Shipyard:        string(decodedShipyard),
		ShipyardVersion: shipyardVersion,
		Stages:          expandedStages,
	}

	err := validateShipyardStagesUnchanged(oldProject, newProject)
	if err != nil {
		return common.ErrInvalidStageChange
	}
	return nil
}

func getUpstreamCredentialSecretName(projectName string) string {
	return fmt.Sprintf("git-credentials-%s", projectName)
}

func getTemporaryUpstreamCredentialSecretName(projectName string) string {
	return fmt.Sprintf("tmp-git-credentials-%s", projectName)
}
