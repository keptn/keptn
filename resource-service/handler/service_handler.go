package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"net/http"
)

type IServiceHandler interface {
	CreateService(context *gin.Context)
	DeleteService(context *gin.Context)
}

type ServiceHandler struct {
	ServiceManager IServiceManager
}

func NewServiceHandler(serviceManager IServiceManager) *ServiceHandler {
	return &ServiceHandler{
		ServiceManager: serviceManager,
	}
}

func (sh *ServiceHandler) CreateService(c *gin.Context) {
	params := &models.CreateServiceParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:   models.Stage{StageName: c.Param(pathParamStageName)},
	}

	createService := &models.CreateServicePayload{}
	if err := c.ShouldBindJSON(createService); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.CreateServicePayload = *createService

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := sh.ServiceManager.CreateService(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
}

func (sh *ServiceHandler) DeleteService(c *gin.Context) {
	params := &models.DeleteServiceParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:   models.Stage{StageName: c.Param(pathParamStageName)},
		Service: models.Service{ServiceName: c.Param(pathParamServiceName)},
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	if err := sh.ServiceManager.DeleteService(*params); err != nil {
		OnAPIError(c, err)
		return
	}
	c.String(http.StatusNoContent, "")
}
