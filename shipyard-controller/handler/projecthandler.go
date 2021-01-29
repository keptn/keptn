package handler

import (
	"github.com/gin-gonic/gin"
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
	ProjectManager IProjectManager
}

func NewProjectHandler(projectmanager IProjectManager) *ProjectHandler {
	return &ProjectHandler{ProjectManager: projectmanager}
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

	allProjects, err := service.ProjectManager.GetProjects()
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

	project, err := service.ProjectManager.GetProjectByName(params.ProjectName)
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

	// validate input
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
			Message: stringp(err.Error()),
		})
		return
	}

	if secretCreated, err := service.ProjectManager.CreateProject(createProjectParams); err != nil {
		if secretCreated {
			if err2 := service.ProjectManager.DeleteSecret(getUpstreamRepoCredsSecretName(*createProjectParams.Name)); err2 != nil {
				//TODO//pm.logger.Error(fmt.Sprintf("could not delete git credentials for project %s: %s", *createProjectParams.Name, err.Error()))
			}
		}
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
	// validate the input
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

	if err := service.ProjectManager.UpdateProject(params); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
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
	projectName := c.Param("project")

	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}

	response, err := service.ProjectManager.DeleteProject(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
