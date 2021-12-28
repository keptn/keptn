package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	"net/http"
)

const pathParamProjectName = "projectName"
const pathParamStageName = "stageName"
const pathParamServiceName = "serviceName"
const pathParamResourceURI = "resourceURI"

func OnAPIError(c *gin.Context, err error) {
	if errors.Is(err, common.ErrProjectAlreadyExists) {
		SetConflictErrorResponse(c, "Project already exists")
	} else if errors.Is(err, common.ErrStageAlreadyExists) {
		SetConflictErrorResponse(c, "Stage already exists")
	} else if errors.Is(err, common.ErrServiceAlreadyExists) {
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
	} else if errors.Is(err, common.ErrServiceNotFound) {
		SetNotFoundErrorResponse(c, "Service not found")
	} else if errors.Is(err, common.ErrResourceNotFound) {
		SetNotFoundErrorResponse(c, "Resource not found")
	} else {
		SetInternalServerErrorResponse(c, "Internal server error")
	}
}

func SetFailedDependencyErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusFailedDependency, models.Error{
		Code:    http.StatusFailedDependency,
		Message: &msg,
	})
}

func SetNotFoundErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.Error{
		Code:    http.StatusNotFound,
		Message: &msg,
	})
}

func SetInternalServerErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    http.StatusInternalServerError,
		Message: &msg,
	})
}

func SetBadRequestErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.Error{
		Code:    http.StatusBadRequest,
		Message: &msg,
	})
}

func SetConflictErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, models.Error{
		Code:    http.StatusConflict,
		Message: &msg,
	})
}
