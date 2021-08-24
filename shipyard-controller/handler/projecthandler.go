package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

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
}

func NewProjectHandler(projectManager IProjectManager, eventSender common.EventSender) *ProjectHandler {
	return &ProjectHandler{
		ProjectManager: projectManager,
		EventSender:    eventSender,
	}
}

// GetTriggeredEvents godoc
// @Summary Get all projects
// @Description Get the list of all projects
// @Tags Projects
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Param   disableUpstreamSync	query		boolean		false	"Disable sync of upstream repo before reading content"
// @Success 200 {object} models.ExpandedProjects	"ok"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [get]
func (ph *ProjectHandler) GetAllProjects(c *gin.Context) {

	params := &operations.GetProjectParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	allProjects, err := ph.ProjectManager.Get()
	if err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}

	sort.Slice(allProjects, func(i, j int) bool {
		return allProjects[i].ProjectName < allProjects[j].ProjectName
	})

	var payload = &models.ExpandedProjects{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Projects:    []*models.ExpandedProject{},
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
	c.JSON(http.StatusOK, payload)
}

// GetTriggeredEvents godoc
// @Summary Get a project by name
// @Description Get a project by its name
// @Tags Projects
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	project		path	string	true	"The name of the project"
// @Success 200 {object} models.ExpandedProject	"ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal Error)
// @Router /project/{project} [get]
func (ph *ProjectHandler) GetProjectByName(c *gin.Context) {
	projectName := c.Param("project")

	project, err := ph.ProjectManager.GetByName(projectName)
	if err != nil {
		if project == nil && err == errProjectNotFound {
			SetNotFoundErrorResponse(nil, c, "Project not found: "+projectName)
			return
		}

		SetInternalServerErrorResponse(err, c)
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
// @Param   project     body    operations.CreateProjectParams     true        "Project"
// @Success 201 {object} operations.CreateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	keptnContext := uuid.New().String()

	createProjectParams := &operations.CreateProjectParams{}
	if err := c.ShouldBindJSON(createProjectParams); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}
	if err := createProjectParams.Validate(); err != nil {
		SetBadRequestErrorResponse(err, c, "Could not validate payload")
		return
	}

	common.LockProject(*createProjectParams.Name)
	defer common.UnlockProject(*createProjectParams.Name)

	if err := ph.sendProjectCreateStartedEvent(keptnContext, createProjectParams); err != nil {
		log.Errorf("could not send project.create.started event: %s", err.Error())
	}

	err, rollback := ph.ProjectManager.Create(createProjectParams)
	if err != nil {
		if err := ph.sendProjectCreateFailFinishedEvent(keptnContext, createProjectParams); err != nil {
			log.Errorf("could not send project.create.finished event: %s", err.Error())
		}
		rollback()
		if err == ErrProjectAlreadyExists {
			SetConflictErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}
	if err := ph.sendProjectCreateSuccessFinishedEvent(keptnContext, createProjectParams); err != nil {
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
// @Param   project     body    operations.UpdateProjectParams     true        "Project"
// @Success 200 {object} operations.UpdateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [put]
func (ph *ProjectHandler) UpdateProject(c *gin.Context) {
	//validate the input
	params := &operations.UpdateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}
	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(err, c, "Could not validate payload")
		return
	}

	common.LockProject(*params.Name)
	defer common.UnlockProject(*params.Name)

	err, rollback := ph.ProjectManager.Update(params)
	if err != nil {
		rollback()
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.Status(http.StatusCreated)
}

//// DeleteProject godoc
//// @Summary Delete a project
//// @Description Delete a project
//// @Tags Projects
//// @Security ApiKeyAuth
//// @Accept  json
//// @Produce  json
//// @Param   project     path    string     true        "Project name"
//// @Success 200 {object} operations.DeleteProjectResponse	"ok"
//// @Failure 400 {object} models.Error "Invalid payload"
//// @Failure 500 {object} models.Error "Internal error"
//// @Router /project/{project} [delete]
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")

	common.LockProject(projectName)
	defer common.UnlockProject(projectName)
	responseMessage, err := ph.ProjectManager.Delete(projectName)
	if err != nil {
		log.Errorf("failed to delete project %s: %s", projectName, err.Error())
		if err := ph.sendProjectDeleteFailFinishedEvent(keptnContext, projectName); err != nil {
			log.Errorf("failed to send finished event: %s", err.Error())
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	if err := ph.sendProjectDeleteSuccessFinishedEvent(keptnContext, projectName); err != nil {
		log.Errorf("failed to send finished event: %s", err.Error())
	}

	c.JSON(http.StatusOK, operations.DeleteProjectResponse{
		Message: responseMessage,
	})
}

func (ph *ProjectHandler) sendProjectCreateStartedEvent(keptnContext string, params *operations.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
		},
	}
	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), eventPayload)
	return ph.EventSender.SendEvent(ce)

}

func (ph *ProjectHandler) sendProjectCreateSuccessFinishedEvent(keptnContext string, params *operations.CreateProjectParams) error {
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

func (ph *ProjectHandler) sendProjectCreateFailFinishedEvent(keptnContext string, params *operations.CreateProjectParams) error {
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
