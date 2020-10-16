package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
	"net/url"
	"strings"
)

type gitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags Projects
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    operations.CreateProjectParams     true        "Project"
// @Success 200 {object} operations.CreateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [post]
func CreateProject(c *gin.Context) {
	// validate the input
	createProjectParams := &operations.CreateProjectParams{}
	if err := c.ShouldBindJSON(createProjectParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateCreateProjectParams(createProjectParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Could not validate payload: " + err.Error()),
		})
		return
	}

	pm, err := newProjectManager()
	if err != nil {

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Message: stringp("Could not process request: " + err.Error()),
		})
		return
	}

	if err := pm.createProject(createProjectParams); err != nil {
		if err == errProjectAlreadyExists {
			c.JSON(http.StatusConflict, models.Error{
				Code:    http.StatusConflict,
				Message: stringp(err.Error()),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project
// @Tags Projects
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project name"
// @Success 200 {object} operations.DeleteProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project [delete]
func DeleteProject(c *gin.Context) {
	projectName := c.Param("project")

	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}
}

var errProjectAlreadyExists = errors.New("project already exists")

type projectManager struct {
	projectAPI  *keptnapi.ProjectHandler
	stagesAPI   *keptnapi.StageHandler
	resourceAPI *keptnapi.ResourceHandler
	logger      keptncommon.LoggerInterface
	secretStore common.SecretStore
}

func newProjectManager() (*projectManager, error) {
	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	if err != nil {
		return nil, fmt.Errorf("could not get configuration-service URL: %s", err.Error())
	}
	secretStore, err := common.NewK8sSecretStore()
	if err != nil {
		return nil, fmt.Errorf("could not initilize secret store: " + err.Error())
	}
	return &projectManager{
		projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
		stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
		resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
		logger:      keptncommon.NewLogger("", "", "shipyard-controller"),
		secretStore: secretStore,
	}, nil
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

	if _, errObj := pm.projectAPI.DeleteProject(keptnapimodels.Project{
		ProjectName: projectName,
	}); errObj != nil {
		return nil, pm.logAndReturnError(fmt.Sprintf("could not delete project: %s", *errObj.Message))
	}

	result.Message = result.Message + pm.getDeleteInfoMessage(projectName)

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

func (pm *projectManager) createProject(params *operations.CreateProjectParams) error {
	// check if the project already exists
	pm.logger.Info(fmt.Sprintf("checking if project %s already exists before creating it", *params.Name))
	project, _ := pm.projectAPI.GetProject(keptnapimodels.Project{
		ProjectName: *params.Name,
	})
	if project != nil {
		pm.logger.Info(fmt.Sprintf("Project %s already exists", *params.Name))
		return errProjectAlreadyExists
	}

	if params.GitRemoteURL != "" && params.GitUser != "" && params.GitToken != "" {
		if err := pm.createUpstreamRepoCredentials(params); err != nil {
			return pm.logAndReturnError(err.Error())
		}
	}

	// create the project
	_, errObj := pm.projectAPI.CreateProject(keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	})

	if errObj != nil {
		return pm.logAndReturnError(fmt.Sprintf("could not create project: %s", *errObj.Message))
	}
	pm.logger.Info(fmt.Sprintf("Project %s created", *params.Name))

	decodedShipyard, err := base64.StdEncoding.DecodeString(*params.Shipyard)
	if err != nil {
		// error should not occur at this stage because the shipyard content has been validated at this stage, but let's check anyways
		return pm.logAndReturnError(fmt.Sprintf("could not decode shipyard: " + err.Error()))
	}
	shipyard, err := common.UnmarshalShipyard(string(decodedShipyard))

	// create the stages
	for _, shipyardStage := range shipyard.Spec.Stages {
		if _, errorObj := pm.stagesAPI.CreateStage(*params.Name, shipyardStage.Name); err != nil {
			return pm.logAndReturnError(fmt.Sprintf("Failed to create stage %s: %s", shipyardStage.Name, errorObj.Message))
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
		return pm.logAndReturnError(fmt.Sprintf("could not upload shipyard.yaml: %s", err.Error()))
	}
	pm.logger.Info("uploaded shipyard.yaml of project " + *params.Name)

	finishedEventPayload := keptnv2.CreateProjectFinishedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
		Project: keptnv2.CreateProjectData{},
	}

	source, _ := url.Parse("shipyard-controller")
	eventType := keptnv2.GetFinishedEventType(keptnv2.CreateProjectTaskName)
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", uuid.New().String())
	event.SetData(cloudevents.ApplicationJSON, finishedEventPayload)

	if err := common.SendEvent(event); err != nil {
		return pm.logAndReturnError("could not send create.project.finished event: " + err.Error())
	}
	return nil
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
	if err := pm.secretStore.CreateSecret("git-credentials-"+*params.Name, map[string][]byte{
		"git-credentials": credsEncoded,
	}); err != nil {
		return fmt.Errorf("could not store git credentials: %s", err.Error())
	}
	pm.logger.Info("stored git credentials for project " + *params.Name)
	return nil
}

func (pm *projectManager) logAndReturnError(msg string) error {
	pm.logger.Error(msg)
	return errors.New(msg)
}

func validateCreateProjectParams(createProjectParams *operations.CreateProjectParams) error {

	if createProjectParams.Name == nil || *createProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*createProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}
	if createProjectParams.Shipyard == nil || *createProjectParams.Shipyard == "" {
		return errors.New("shipyard must contain a valid shipyard spec encoded in base64")
	}
	shipyard := &keptnv2.Shipyard{}
	decodeString, err := base64.StdEncoding.DecodeString(*createProjectParams.Shipyard)
	if err != nil {
		return errors.New("could not decode shipyard content using base64 decoder: " + err.Error())
	}

	err = yaml.Unmarshal(decodeString, shipyard)
	if err != nil {
		return fmt.Errorf("could not unmarshal provided shipyard content: %s", err.Error())
	}

	if err := common.ValidateShipyardVersion(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateShipyardStages(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	return nil
}
