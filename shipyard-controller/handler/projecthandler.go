package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/config"
	"gopkg.in/yaml.v3"

	"net/http"
	"os"
	"sort"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

type ProjectValidator struct {
	ProjectNameMaxSize int
}

func (p ProjectValidator) Validate(params interface{}) error {
	switch t := params.(type) {
	case *models.CreateProjectParams:
		return p.validateCreateProjectParams(t)
	case *models.UpdateProjectParams:
		return p.validateUpdateProjectParams(t)
	default:
		return nil
	}
}
func (p ProjectValidator) validateCreateProjectParams(createProjectParams *models.CreateProjectParams) error {
	if createProjectParams.Name == nil || *createProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if len(*createProjectParams.Name) > p.ProjectNameMaxSize {
		return fmt.Errorf("project name exceeds maximum size of %d characters", p.ProjectNameMaxSize)
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
		return errors.New("could not decode shipyard content")
	}

	err = yaml.Unmarshal(decodeString, shipyard)
	if err != nil {
		return fmt.Errorf("could not unmarshal provided shipyard content")
	}

	if err := common.ValidateShipyardVersion(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateShipyardStages(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateGitRemoteURL(createProjectParams.GitRemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	if createProjectParams.GitPrivateKey != "" && createProjectParams.GitToken != "" {
		return fmt.Errorf("privateKey and token cannot be used together")
	}

	if createProjectParams.GitPrivateKey != "" && createProjectParams.GitProxyURL != "" {
		return fmt.Errorf("privateKey and proxy cannot be used together")
	}

	if createProjectParams.GitPrivateKey != "" {
		decodeString, err = base64.StdEncoding.DecodeString(createProjectParams.GitPrivateKey)
		if err != nil {
			return errors.New("could not decode privateKey content")
		}
	}

	if createProjectParams.GitPrivateKey != "" && createProjectParams.GitPemCertificate != "" {
		return fmt.Errorf("SSH authorization and PEM Certificate be used together")
	}

	if createProjectParams.GitPemCertificate != "" {
		decodeString, err = base64.StdEncoding.DecodeString(createProjectParams.GitPemCertificate)
		if err != nil {
			return errors.New("could not decode PEM Certificate content")
		}
	}

	return nil
}

func (p ProjectValidator) validateUpdateProjectParams(updateProjectParams *models.UpdateProjectParams) error {
	if updateProjectParams.Name == nil || *updateProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*updateProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}

	if updateProjectParams.Shipyard != nil && *updateProjectParams.Shipyard != "" {
		shipyard := &keptnv2.Shipyard{}
		decodeString, err := base64.StdEncoding.DecodeString(*updateProjectParams.Shipyard)
		if err != nil {
			return errors.New("could not decode shipyard content")
		}

		err = yaml.Unmarshal(decodeString, shipyard)
		if err != nil {
			return fmt.Errorf("could not unmarshal provided shipyard content")
		}

		if err := common.ValidateShipyardVersion(shipyard); err != nil {
			return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
		}

		if err := common.ValidateShipyardStages(shipyard); err != nil {
			return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
		}
	}

	if err := common.ValidateGitRemoteURL(updateProjectParams.GitRemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	if updateProjectParams.GitPrivateKey != "" && updateProjectParams.GitToken != "" {
		return fmt.Errorf("privateKey and token cannot be used together")
	}

	if updateProjectParams.GitPrivateKey != "" && updateProjectParams.GitProxyURL != "" {
		return fmt.Errorf("privateKey and proxy cannot be used together")
	}

	if updateProjectParams.GitPrivateKey != "" {
		_, err := base64.StdEncoding.DecodeString(updateProjectParams.GitPrivateKey)
		if err != nil {
			return errors.New("could not decode privateKey content")
		}
	}

	if updateProjectParams.GitPrivateKey != "" && updateProjectParams.GitPemCertificate != "" {
		return fmt.Errorf("SSH authorization and PEM Certificate be used together")
	}

	if updateProjectParams.GitPemCertificate != "" {
		_, err := base64.StdEncoding.DecodeString(updateProjectParams.GitPemCertificate)
		if err != nil {
			return errors.New("could not decode PEM Certificate content")
		}
	}

	return nil
}

type IProjectHandler interface {
	GetAllProjects(context *gin.Context)
	GetProjectByName(context *gin.Context)
	CreateProject(context *gin.Context)
	UpdateProject(context *gin.Context)
	DeleteProject(context *gin.Context)
}

type ProjectHandler struct {
	ProjectManager IProjectManager
	EventSender    common.EventSender
	Env            config.EnvConfig
}

func NewProjectHandler(projectManager IProjectManager, eventSender common.EventSender, env config.EnvConfig) *ProjectHandler {
	return &ProjectHandler{
		ProjectManager: projectManager,
		EventSender:    eventSender,
		Env:            env,
	}
}

// GetAllProjects godoc
// @Summary Get all projects
// @Description Get the list of all projects
// @Tags Projects
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Param   disableUpstreamSync	query		boolean		false	"Disable sync of upstream repo before reading content"
// @Success 200 {object} apimodels.ExpandedProjects	"ok"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [get]
func (ph *ProjectHandler) GetAllProjects(c *gin.Context) {
	params := &models.GetProjectParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	allProjects, err := ph.ProjectManager.Get()
	if err != nil {
		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	sort.Slice(allProjects, func(i, j int) bool {
		return allProjects[i].ProjectName < allProjects[j].ProjectName
	})

	var payload = &apimodels.ExpandedProjects{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Projects:    []*apimodels.ExpandedProject{},
	}

	paginationInfo := common.Paginate(len(allProjects), params.PageSize, params.NextPageKey)
	totalCount := len(allProjects)
	if paginationInfo.NextPageKey < int64(totalCount) {
		payload.Projects = append(payload.Projects, allProjects[paginationInfo.NextPageKey:paginationInfo.EndIndex]...)
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	c.JSON(http.StatusOK, payload)
}

// GetProjectByName godoc
// @Summary Get a project by name
// @Description Get a project by its name
// @Tags Projects
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	project		path	string	true	"The name of the project"
// @Success 200 {object} apimodels.ExpandedProject	"ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal Error)
// @Router /project/{project} [get]
func (ph *ProjectHandler) GetProjectByName(c *gin.Context) {
	projectName := c.Param("project")

	project, err := ph.ProjectManager.GetByName(projectName)
	if err != nil {
		if project == nil && errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ProjectNotFoundMsg, projectName))
			return
		}

		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)

}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags Projects
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    models.CreateProjectParams     true        "Project"
// @Success 201 {object} models.CreateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	keptnContext := uuid.New().String()

	params := &models.CreateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	automaticProvisioningURL := os.Getenv("AUTOMATIC_PROVISIONING_URL")
	if automaticProvisioningURL != "" && createProjectParams.GitRemoteURL == "" {
		values := map[string]string{"name": "John Doe", "occupation": "gardener"}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Errorf(UnableMarshallProvisioningData, err.Error())
			SetFailedDependencyErrorResponse(c, fmt.Sprintf(UnableMarshallProvisioningData, err.Error()))
		}

		_, err = http.Post(automaticProvisioningURL+"/repository", "application/json", bytes.NewBuffer(json_data))

		if err != nil {
			log.Errorf(UnableProvisionInstance, err.Error())
			SetFailedDependencyErrorResponse(c, fmt.Sprintf(UnableProvisionInstance, err.Error()))
		}
	}

	projectValidator := ProjectValidator{ProjectNameMaxSize: ph.Env.ProjectNameMaxSize}
	if err := projectValidator.Validate(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidPayloadMsg, err.Error()))
		return
	}

	common.LockProject(*createProjectParams.Name)
	defer common.UnlockProject(*createProjectParams.Name)

	if err := ph.sendProjectCreateStartedEvent(keptnContext, params); err != nil {
		log.Errorf("could not send project.create.started event: %s", err.Error())
	}

	err, rollback := ph.ProjectManager.Create(params)
	if err != nil {
		if err := ph.sendProjectCreateFailFinishedEvent(keptnContext, params); err != nil {
			log.Errorf("could not send project.create.finished event: %s", err.Error())
		}
		rollback()
		if errors.Is(err, ErrProjectAlreadyExists) {
			SetConflictErrorResponse(c, err.Error())
			return
		}
		SetInternalServerErrorResponse(c, err.Error())
		return
	}
	if err := ph.sendProjectCreateSuccessFinishedEvent(keptnContext, params); err != nil {
		log.Errorf("could not send project.create.finished event: %s", err.Error())
	}

	c.Status(http.StatusCreated)

}

// UpdateProject godoc
// @Summary Updates a project
// @Description Updates project
// @Tags Projects
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    models.UpdateProjectParams     true        "Project"
// @Success 200 {object} models.UpdateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 424 {object} models.Error "Failed Dependency"
// @Failure 404 {object} models.Error "Not Found"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [put]
func (ph *ProjectHandler) UpdateProject(c *gin.Context) {
	// validate the input
	params := &models.UpdateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}
	projectValidator := ProjectValidator{ProjectNameMaxSize: ph.Env.ProjectNameMaxSize}
	if err := projectValidator.Validate(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidPayloadMsg, err.Error()))
		return
	}

	common.LockProject(*params.Name)
	defer common.UnlockProject(*params.Name)

	err, rollback := ph.ProjectManager.Update(params)
	if err != nil {
		rollback()
		if errors.Is(err, common.ErrConfigStoreInvalidToken) {
			SetFailedDependencyErrorResponse(c, err.Error())
			return
		}
		if errors.Is(err, common.ErrConfigStoreUpstreamNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
			return
		}
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
			return
		}
		if errors.Is(err, ErrInvalidStageChange) {
			SetBadRequestErrorResponse(c, err.Error())
			return
		}
		SetInternalServerErrorResponse(c, ErrInternalError.Error())
		return
	}
	c.Status(http.StatusCreated)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project
// @Tags Projects
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project name"
// @Success 200 {object} models.DeleteProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project} [delete]
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")

	automaticProvisioningURL := os.Getenv("AUTOMATIC_PROVISIONING_URL")
	if automaticProvisioningURL != "" {
		values := map[string]string{"name": "John Doe", "occupation": "gardener"}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Errorf(UnableMarshallProvisioningData, err.Error())
			SetFailedDependencyErrorResponse(c, fmt.Sprintf(UnableMarshallProvisioningData, err.Error()))
		}

		_, err = http.NewRequest(http.MethodDelete, automaticProvisioningURL+"/repository", bytes.NewBuffer(json_data))

		if err != nil {
			log.Errorf(UnableProvisionDelete, err.Error())
			SetFailedDependencyErrorResponse(c, fmt.Sprintf(UnableProvisionDelete, err.Error()))
		}
	}

	common.LockProject(projectName)
	defer common.UnlockProject(projectName)
	responseMessage, err := ph.ProjectManager.Delete(projectName)
	if err != nil {
		log.Errorf("failed to delete project %s: %s", projectName, err.Error())
		if err := ph.sendProjectDeleteFailFinishedEvent(keptnContext, projectName); err != nil {
			log.Errorf("failed to send finished event: %s", err.Error())
		}
		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	if err := ph.sendProjectDeleteSuccessFinishedEvent(keptnContext, projectName); err != nil {
		log.Errorf("failed to send finished event: %s", err.Error())
	} else {
		log.Debug("Deleted project ", projectName)
	}

	c.JSON(http.StatusOK, models.DeleteProjectResponse{
		Message: responseMessage,
	})
}

func (ph *ProjectHandler) sendProjectCreateStartedEvent(keptnContext string, params *models.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
		},
	}
	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), eventPayload)
	return ph.EventSender.SendEvent(ce)

}

func (ph *ProjectHandler) sendProjectCreateSuccessFinishedEvent(keptnContext string, params *models.CreateProjectParams) error {
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

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), eventPayload)
	return ph.EventSender.SendEvent(ce)
}

func (ph *ProjectHandler) sendProjectCreateFailFinishedEvent(keptnContext string, params *models.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), eventPayload)
	return ph.EventSender.SendEvent(ce)
}

func (ph *ProjectHandler) sendProjectDeleteSuccessFinishedEvent(keptnContext, projectName string) error {
	eventPayload := keptnv2.ProjectDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName), eventPayload)
	return ph.EventSender.SendEvent(ce)
}

func (ph *ProjectHandler) sendProjectDeleteFailFinishedEvent(keptnContext, projectName string) error {
	eventPayload := keptnv2.ProjectDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName), eventPayload)
	return ph.EventSender.SendEvent(ce)
}
