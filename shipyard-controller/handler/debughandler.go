package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type IDebugHandler interface {
	GetAllProjects(context *gin.Context)
	GetSequenceByID(context *gin.Context)
	GetAllSequencesForProject(context *gin.Context)
	GetAllEvents(context *gin.Context)
	GetEventByID(context *gin.Context)
}

type DebugHandler struct {
	DebugManager IDebugManager
}

func NewDebugHandler(debugManager IDebugManager) *DebugHandler {
	return &DebugHandler{
		DebugManager: debugManager,
	}
}

// GetAllProjects godoc
// @Summary      Get all keptn projects
// @Description  Get all keptn projects
// @Tags         Project
// @Success      200                  {object}  []apimodels.ExpandedProject     "ok"
// @Failure      400                  {object}  models.Error                    "Bad Request"
// @Failure      500                  {object}  models.Error                    "Internal error"
// @Router       /debug/project [get]
func (dh *DebugHandler) GetAllProjects(c *gin.Context) {
	params := &models.GetProjectParams{}

	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	projects, err := dh.DebugManager.GetAllProjects()

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].ProjectName < projects[j].ProjectName
	})

	var payload = &apimodels.ExpandedProjects{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Projects:    []*apimodels.ExpandedProject{},
	}

	paginationInfo := common.Paginate(len(projects), params.PageSize, params.NextPageKey)
	totalCount := len(projects)
	if paginationInfo.NextPageKey < int64(totalCount) {
		payload.Projects = append(payload.Projects, projects[paginationInfo.NextPageKey:paginationInfo.EndIndex]...)
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	c.JSON(http.StatusOK, payload)
}

// GetAllSequencesForProject godoc
// @Summary      Get all sequences for specific project
// @Description  Get all the sequences which are present in a project
// @Tags         Sequence
// @Param        project              path      string                    true "The name of the project"
// @Success      200                  {object}  []models.SequenceState    "ok"
// @Failure      400                  {object}  models.Error              "Bad Request"
// @Failure      404                  {object}  models.Error              "not found"
// @Failure      500                  {object}  models.Error              "Internal error"
// @Router       /debug/project/{project} [get]
func (dh *DebugHandler) GetAllSequencesForProject(c *gin.Context) {
	projectName := c.Param("project")
	payload, err := dh.DebugManager.GetAllSequencesForProject(projectName)

	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ProjectNotFoundMsg, projectName))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, payload)
}

// GetSequenceByID godoc
// @Summary      Get a sequence with the shkeptncontext
// @Description  Get a specific sequence of a project which is identified by the shkeptncontext
// @Tags         Sequence
// @Param        project              path      string                    true  "The name of the project"
// @Param        shkeptncontext       path      string                    true  "The shkeptncontext"
// @Success      200                  {object}  models.SequenceState      "ok"
// @Failure      400                  {object}  models.Error              "Bad Request"
// @Failure      404                  {object}  models.Error              "not found"
// @Failure      500                  {object}  models.Error              "Internal error"
// @Router       /debug/project/{project}/shkeptncontext/{shkeptncontext} [get]
func (dh *DebugHandler) GetSequenceByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")
	sequence, err := dh.DebugManager.GetSequenceByID(projectName, shkeptncontext)

	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ProjectNotFoundMsg, projectName))
			return
		}

		if errors.Is(err, ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(UnableFindSequenceMsg, shkeptncontext))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, sequence)
}

// GetAllEvents godoc
// @Summary      Get all the Events
// @Description  Gets all the events of a project with the given shkeptncontext
// @Tags         Sequence
// @Param        project              path      string                             true  "The name of the project"
// @Param        shkeptncontext       path      string                             true  "The shkeptncontext"
// @Success      200                  {object}  []models.KeptnContextExtendedCE    "ok"
// @Failure      400                  {object}  models.Error                       "Bad Request"
// @Failure      404                  {object}  models.Error                       "not found"
// @Failure      500                  {object}  models.Error                       "Internal error"
// @Router       /debug/project/{project}/shkeptncontext/{shkeptncontext}/event [get]
func (dh *DebugHandler) GetAllEvents(c *gin.Context) {
	params := &models.GetProjectParams{}
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")

	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	events, err := dh.DebugManager.GetAllEvents(projectName, shkeptncontext)

	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ProjectNotFoundMsg, projectName))
			return
		}

		if errors.Is(err, ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(UnableFindSequenceMsg, shkeptncontext))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	var payload = &apimodels.Events{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Events:      []*apimodels.KeptnContextExtendedCE{},
	}

	paginationInfo := common.Paginate(len(events), params.PageSize, params.NextPageKey)
	totalCount := len(events)
	if paginationInfo.NextPageKey < int64(totalCount) {
		payload.Events = append(payload.Events, events[paginationInfo.NextPageKey:paginationInfo.EndIndex]...)
	}

	payload.TotalCount = float64(totalCount)
	c.JSON(http.StatusOK, payload)
}

// GetEventByID godoc
// @Summary      Get a single Event
// @Description  Gets a single event of a project with the given shkeptncontext and eventId
// @Tags         Sequence
// @Param        project              path      string                             true  "The name of the project"
// @Param        shkeptncontext       path      string                             true  "The shkeptncontext"
// @Param        eventId              path      string                             true  "The Id of the event"
// @Success      200                  {object}  models.KeptnContextExtendedCE      "ok"
// @Failure      400                  {object}  models.Error                       "Bad Request"
// @Failure      404                  {object}  models.Error                       "not found"
// @Failure      500                  {object}  models.Error                       "Internal error"
// @Router       /debug/project/{project}/shkeptncontext/{shkeptncontext}/event/{eventId} [get]
func (dh *DebugHandler) GetEventByID(c *gin.Context) {
	params := &models.GetProjectParams{}

	shkeptncontext := c.Param("shkeptncontext")
	eventId := c.Param("eventId")
	projectName := c.Param("project")

	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	event, err := dh.DebugManager.GetEventByID(projectName, shkeptncontext, eventId)

	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ProjectNotFoundMsg, projectName))
			return
		}

		if errors.Is(err, ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(UnableFindSequenceMsg, shkeptncontext))
			return
		}

		if errors.Is(err, ErrNoMatchingEvent) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(EventNotFoundMsg, eventId))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, event)
}
