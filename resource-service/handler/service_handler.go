package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
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
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := sh.ServiceManager.CreateService(*params)
	if err != nil {
		if errors.Is(err, common.ErrServiceAlreadyExists) {
			SetConflictErrorResponse(c, "Service already exists")
		} else if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else if errors.Is(err, common.ErrStageNotFound) {
			SetNotFoundErrorResponse(c, "Stage not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
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
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project does not exist")
		} else if errors.Is(err, common.ErrStageNotFound) {
			SetNotFoundErrorResponse(c, "Stage does not exist")
		} else if errors.Is(err, common.ErrServiceNotFound) {
			SetNotFoundErrorResponse(c, "Stage does not exist")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}
	c.String(http.StatusNoContent, "")
}
