package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
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
	ProjectManager *ProjectManager
	EventSender    keptn.EventSender
}

func NewProjectHandler(projectManager *ProjectManager, eventSender keptn.EventSender) *ProjectHandler {
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
func (service *ProjectHandler) GetAllProjects(c *gin.Context) {

	params := &operations.GetProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}

	allProjects, err := service.ProjectManager.Get()
	if err != nil {
		sendInternalServerErrorResponse(err, c)
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
// @Param	projectName		path	string	true	"The name of the project"
// @Success 200 {object} models.ExpandedProject	"ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal Error)
// @Router /project/{projectName} [get]
func (service *ProjectHandler) GetProjectByName(c *gin.Context) {
	params := &operations.GetProjectProjectNameParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}

	project, err := service.ProjectManager.GetByName(params.ProjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
	if project == nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp("Project not found: " + params.ProjectName),
		})
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
// @Success 200 {object} operations.CreateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [post]
func (service *ProjectHandler) CreateProject(c *gin.Context) {
	keptnContext := uuid.New().String()

	createProjectParams := &operations.CreateProjectParams{}
	if err := c.ShouldBindJSON(createProjectParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateCreateProjectParams(createProjectParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp(err.Error()),
		})
		return
	}

	if err := service.sendProjectCreateStartedEvent(keptnContext, createProjectParams); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
	}

	err, rollback := service.ProjectManager.Create(createProjectParams)
	if err != nil {
		if err := service.sendProjectCreateFailFinishedEvent(keptnContext, createProjectParams); err != nil {
			//LOG MESSAGE ONLY
		}
		rollback()
		if err == ErrProjectAlreadyExists {
			c.JSON(http.StatusConflict, models.Error{
				Code:    http.StatusConflict,
				Message: stringp(err.Error()),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, models.Error{
				Code:    http.StatusInternalServerError,
				Message: stringp(err.Error()),
			})
			return
		}
	}
	if err := service.sendProjectCreateSuccessFinishedEvent(keptnContext, createProjectParams); err != nil {
		//LOG MESSAGE ONLY
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
func (service *ProjectHandler) UpdateProject(c *gin.Context) {
	//validate the input
	params := &operations.UpdateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateUpdateProjectParams(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Could not validate payload: " + err.Error()),
		})
		return
	}

	err, rollback := service.ProjectManager.Update(params)
	if err != nil {
		rollback()
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
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
//// @Router /project/:project [delete]
func (service *ProjectHandler) DeleteProject(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")

	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}

	err, response := service.ProjectManager.Delete(projectName)
	if err != nil {
		if err := service.sendProjectDeleteFailFinishedEvent(keptnContext, projectName); err != nil {
			//LOG MESSAGE ONLY
		}

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}

	if err := service.sendProjectDeleteSuccessFinishedEvent(keptnContext, projectName); err != nil {
		//LOG MESSAGE ONLY
	}

	c.JSON(http.StatusOK, response)
}

func (service *ProjectHandler) sendProjectCreateStartedEvent(keptnContext string, params *operations.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
		},
	}
	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), eventPayload)
	return service.EventSender.SendEvent(ce)

}

func (service *ProjectHandler) sendProjectCreateSuccessFinishedEvent(keptnContext string, params *operations.CreateProjectParams) error {
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
	return service.EventSender.SendEvent(ce)
}

func (service *ProjectHandler) sendProjectCreateFailFinishedEvent(keptnContext string, params *operations.CreateProjectParams) error {
	eventPayload := keptnv2.ProjectCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), eventPayload)
	return service.EventSender.SendEvent(ce)
}

func (service *ProjectHandler) sendProjectDeleteSuccessFinishedEvent(keptnContext, projectName string) error {
	eventPayload := keptnv2.ProjectDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName), eventPayload)
	return service.EventSender.SendEvent(ce)
}

func (service *ProjectHandler) sendProjectDeleteFailFinishedEvent(keptnContext, projectName string) error {
	eventPayload := keptnv2.ProjectDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	ce := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName), eventPayload)
	return service.EventSender.SendEvent(ce)
}
