package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

type IProjectHandler interface {
	CreateProject(context *gin.Context)
	UpdateProject(context *gin.Context)
	DeleteProject(context *gin.Context)
}

type ProjectHandler struct {
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
			Message: stringp(err.Error()),
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

	if secretCreated, err := pm.createProject(createProjectParams); err != nil {
		if secretCreated {
			if err2 := pm.secretStore.DeleteSecret(getUpstreamRepoCredsSecretName(*createProjectParams.Name)); err2 != nil {
				pm.logger.Error(fmt.Sprintf("could not delete git credentials for project %s: %s", *createProjectParams.Name, err.Error()))
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
// @Param   project     body    operations.CreateProjectParams     true        "Project"
// @Success 200 {object} operations.CreateProjectResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [put]
func (service *ProjectHandler) UpdateProject(c *gin.Context) {
	// validate the input
	createProjectParams := &operations.CreateProjectParams{}
	if err := c.ShouldBindJSON(createProjectParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateUpdateProjectParams(createProjectParams); err != nil {
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

	if err := pm.updateProject(createProjectParams); err != nil {
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

	pm, err := newProjectManager()
	if err != nil {

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Message: stringp("Could not process request: " + err.Error()),
		})
		return
	}

	response, err := pm.deleteProject(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

func NewProjectHandler() IProjectHandler {
	return &ProjectHandler{}
}
