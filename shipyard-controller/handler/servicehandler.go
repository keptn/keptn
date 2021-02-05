package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

type IServiceHandler interface {
	CreateService(context *gin.Context)
	DeleteService(context *gin.Context)
}

type ServiceHandler struct {
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Param   service     body    operations.CreateServiceParams     true        "Project"
// @Success 200 {object} operations.CreateServiceResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project/service [post]
func (service *ServiceHandler) CreateService(c *gin.Context) {
	projectName := c.Param("project")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}
	// validate the input
	createServiceParams := &operations.CreateServiceParams{}
	if err := c.ShouldBindJSON(createServiceParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateCreateServiceParams(createServiceParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Could not validate payload: " + err.Error()),
		})
		return
	}

	sm, err := newServiceManager()
	if err != nil {

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Message: stringp("Could not process request: " + err.Error()),
		})
		return
	}

	if err := sm.createService(projectName, createServiceParams); err != nil {
		if err == errServiceAlreadyExists {
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

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Param   service     path    string     true        "Service"
// @Success 200 {object} operations.DeleteServiceResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project/service/:service [delete]
func (service *ServiceHandler) DeleteService(c *gin.Context) {
	projectName := c.Param("project")
	serviceName := c.Param("service")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a service name"),
		})
	}

	sm, err := newServiceManager()
	if err != nil {

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Message: stringp("Could not process request: " + err.Error()),
		})
		return
	}

	if err := sm.deleteService(projectName, serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, &operations.DeleteServiceResponse{})
}

func NewServiceHandler() IServiceHandler {
	return &ServiceHandler{}
}
