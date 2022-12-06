package handler

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/config"
	"github.com/keptn/keptn/shipyard-controller/internal/provisioner"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"gopkg.in/yaml.v3"

	"net/http"
	"sort"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

type ProjectValidator struct {
	ProjectNameMaxSize       int
	AutomaticProvisioningURL string
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

	if p.AutomaticProvisioningURL != "" && createProjectParams.GitCredentials == nil {
		return nil
	}

	if createProjectParams.GitCredentials == nil {
		return fmt.Errorf("gitCredentials cannot be empty")
	}

	if err := common.ValidateGitRemoteURL(createProjectParams.GitCredentials.RemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	if createProjectParams.GitCredentials.HttpsAuth != nil && createProjectParams.GitCredentials.SshAuth != nil {
		return fmt.Errorf("SSH and HTTPS authorization cannot be used together")
	}

	if createProjectParams.GitCredentials.SshAuth != nil && createProjectParams.GitCredentials.SshAuth.PrivateKey != "" {
		decodeString, err = base64.StdEncoding.DecodeString(createProjectParams.GitCredentials.SshAuth.PrivateKey)
		if err != nil {
			return errors.New("could not decode privateKey content")
		}
	}

	if createProjectParams.GitCredentials.HttpsAuth != nil && createProjectParams.GitCredentials.HttpsAuth.Certificate != "" {
		decodeString, err = base64.StdEncoding.DecodeString(createProjectParams.GitCredentials.HttpsAuth.Certificate)
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

	if updateProjectParams.GitCredentials == nil {
		return nil
	}

	if err := common.ValidateGitRemoteURL(updateProjectParams.GitCredentials.RemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	if updateProjectParams.GitCredentials.HttpsAuth != nil && updateProjectParams.GitCredentials.SshAuth != nil {
		return fmt.Errorf("SSH and HTTPS authorization cannot be used together")
	}

	if updateProjectParams.GitCredentials.SshAuth != nil && updateProjectParams.GitCredentials.SshAuth.PrivateKey != "" {
		_, err := base64.StdEncoding.DecodeString(updateProjectParams.GitCredentials.SshAuth.PrivateKey)
		if err != nil {
			return errors.New("could not decode privateKey content")
		}
	}

	if updateProjectParams.GitCredentials.HttpsAuth != nil && updateProjectParams.GitCredentials.HttpsAuth.Certificate != "" {
		_, err := base64.StdEncoding.DecodeString(updateProjectParams.GitCredentials.HttpsAuth.Certificate)
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
	ProjectManager        IProjectManager
	EventSender           common.EventSender
	Env                   config.EnvConfig
	RepositoryProvisioner provisioner.IRepositoryProvisioner
	RemoteURLValidator    provisioner.RemoteURLValidator
}

func NewProjectHandler(projectManager IProjectManager, eventSender common.EventSender, env config.EnvConfig, repositoryProvisioner provisioner.IRepositoryProvisioner, remoteURLValidator provisioner.RemoteURLValidator) *ProjectHandler {
	return &ProjectHandler{
		ProjectManager:        projectManager,
		EventSender:           eventSender,
		Env:                   env,
		RepositoryProvisioner: repositoryProvisioner,
		RemoteURLValidator:    remoteURLValidator,
	}
}

// GetAllProjects godoc
// @Summary      Get all projects
// @Description  Get the list of all projects
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}projects:read</span>
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        pageSize             query     int                         false  "The number of items to return"
// @Param        nextPageKey          query     string                      false  "Pointer to the next set of items"
// @Param        disableUpstreamSync  query     boolean                     false  "Disable sync of upstream repo before reading content"
// @Success      200                  {object}  apimodels.ExpandedProjects  "ok"
// @Failure      400                  {object}  models.Error                "Invalid payload"
// @Failure      500                  {object}  models.Error                "Internal error"
// @Router       /project [get]
func (ph *ProjectHandler) GetAllProjects(c *gin.Context) {
	params := &models.GetProjectParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
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
// @Summary      Get a project by name
// @Description  Get a project by its name
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}projects:read</span>
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  path      string                     true  "The name of the project"
// @Success      200      {object}  apimodels.ExpandedProject  "ok"
// @Failure      404      {object}  models.Error               "Not found"
// @Failure      500      {object}  models.Error               "Internal Error)
// @Router       /project/{project} [get]
func (ph *ProjectHandler) GetProjectByName(c *gin.Context) {
	projectName := c.Param("project")

	project, err := ph.ProjectManager.GetByName(projectName)
	if err != nil {
		if project == nil && errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.ProjectNotFoundMsg, projectName))
			return
		}

		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)

}

// CreateProject godoc
// @Summary      Create a new project
// @Description  Create a new project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}projects:write</span>
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  body      models.CreateProjectParams    true  "Project"
// @Success      201      {object}  models.CreateProjectResponse  "ok"
// @Failure      400      {object}  models.Error                  "Invalid payload"
// @Failure      409      {object}  models.Error                  "Conflict"
// @Failure      424      {object}  models.Error                  "Failed dependency"
// @Failure      500      {object}  models.Error                  "Internal error"
// @Router       /project [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	keptnContext := uuid.New().String()

	params := &models.CreateProjectParams{}
	if err := DecodeInputData(c.Request.Body, params); err != nil {
		log.Debugf("bad json %s", err.Error())
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}

	projectValidator := ProjectValidator{ProjectNameMaxSize: ph.Env.ProjectNameMaxSize, AutomaticProvisioningURL: ph.Env.AutomaticProvisioningURL}
	if err := projectValidator.Validate(params); err != nil {
		log.Debugf("invalid project %s", err.Error())
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidPayloadMsg, err.Error()))
		return
	}

	isAutoProvisioned := false
	if provideProvisionedRepository(ph.Env.AutomaticProvisioningURL, params) {
		provisioningData, err := ph.RepositoryProvisioner.ProvideRepository(*params.Name, common.GetKeptnNamespace())
		if err != nil {
			log.Errorf(err.Error())
			SetFailedDependencyErrorResponse(c, err.Error())
			return
		}

		log.Debugf("Provisioner data\nGit URL: %s\nUser: %s\n", provisioningData.GitRemoteURL, provisioningData.GitUser)

		params.GitCredentials = &apimodels.GitAuthCredentials{
			RemoteURL: provisioningData.GitRemoteURL,
			HttpsAuth: &apimodels.HttpsGitAuth{
				InsecureSkipTLS: false,
				Token:           provisioningData.GitToken,
			},
			User: provisioningData.GitUser,
		}
		isAutoProvisioned = true
	} else if err := ph.RemoteURLValidator.Validate(params.GitCredentials.RemoteURL); err != nil {
		log.Debugf("invalid URL %s", err.Error())
		SetUnprocessableEntityResponse(c, fmt.Sprintf(common.InvalidRemoteURLMsg, params.GitCredentials.RemoteURL))
		return
	}

	common.LockProject(*params.Name)
	defer common.UnlockProject(*params.Name)

	if err := ph.sendProjectCreateStartedEvent(keptnContext, params); err != nil {
		log.Errorf("could not send project.create.started event: %s", err.Error())
	}

	err, rollback := ph.ProjectManager.Create(params, models.InternalCreateProjectOptions{IsUpstreamAutoProvisioned: isAutoProvisioned})
	if err != nil {
		if err := ph.sendProjectCreateFailFinishedEvent(keptnContext, params); err != nil {
			log.Errorf("could not send project.create.finished event: %s", err.Error())
		}

		rollback()
		log.Debugf("rolled back %s", err.Error())
		mapError(c, err)
		return
	}
	if err := ph.sendProjectCreateSuccessFinishedEvent(keptnContext, params); err != nil {
		log.Errorf("could not send project.create.finished event: %s", err.Error())
	}

	c.Status(http.StatusCreated)

}

// UpdateProject godoc
// @Summary      Updates a project
// @Description  Updates project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}projects:write</span>
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  body      models.UpdateProjectParams    true  "Project"
// @Success      201      {object}  models.UpdateProjectResponse  "ok"
// @Failure      400      {object}  models.Error                  "Bad Request"
// @Failure      404      {object}  models.Error                  "Not Found"
// @Failure      424      {object}  models.Error                  "Failed Dependency"
// @Failure      500      {object}  models.Error                  "Internal error"
// @Router       /project [put]
func (ph *ProjectHandler) UpdateProject(c *gin.Context) {
	// validate the input
	params := &models.UpdateProjectParams{}
	if err := DecodeInputData(c.Request.Body, params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}
	projectValidator := ProjectValidator{ProjectNameMaxSize: ph.Env.ProjectNameMaxSize, AutomaticProvisioningURL: ph.Env.AutomaticProvisioningURL}
	if err := projectValidator.Validate(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidPayloadMsg, err.Error()))
		return
	}

	if params.GitCredentials != nil && params.GitCredentials.RemoteURL != "" {
		if err := ph.RemoteURLValidator.Validate(params.GitCredentials.RemoteURL); err != nil {
			SetUnprocessableEntityResponse(c, fmt.Sprintf(common.InvalidRemoteURLMsg, params.GitCredentials.RemoteURL))
			return
		}
	}

	common.LockProject(*params.Name)
	defer common.UnlockProject(*params.Name)

	err, rollback := ph.ProjectManager.Update(params)
	if err != nil {
		rollback()
		mapError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

// DeleteProject godoc
// @Summary      Delete a project
// @Description  Delete a project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}projects:delete</span>
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  path      string                        true  "Project name"
// @Success      200      {object}  models.DeleteProjectResponse  "ok"
// @Failure      400      {object}  models.Error                  "Invalid payload"
// @Failure      424      {object}  models.Error                  "Failed Dependency"
// @Failure      500      {object}  models.Error                  "Internal error"
// @Router       /project/{project} [delete]
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")
	namespace := c.Param("namespace")

	common.LockProject(projectName)
	defer common.UnlockProject(projectName)

	automaticProvisioningURL := ph.Env.AutomaticProvisioningURL
	if automaticProvisioningURL != "" {
		err := ph.RepositoryProvisioner.DeleteRepository(projectName, namespace)
		if err != nil {
			// a failure to clean up the provisioned repo should not prevent the project delete
			log.Errorf("Automatic Provisioning error: %s", err.Error())
		}
	}

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
			ProjectName: *params.Name,
			Shipyard:    *params.Shipyard,
		},
	}

	if params.GitCredentials != nil {
		eventPayload.CreatedProject.GitRemoteURL = params.GitCredentials.RemoteURL
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

func provideProvisionedRepository(provisionURL string, params *models.CreateProjectParams) bool {
	if provisionURL != "" && params.GitCredentials == nil {
		return true
	}
	return false
}
