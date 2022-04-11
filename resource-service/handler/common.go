package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

const pathParamProjectName = "projectName"
const pathParamStageName = "stageName"
const pathParamServiceName = "serviceName"
const pathParamResourceURI = "resourceURI"

func OnAPIError(c *gin.Context, err error) {
	logger.Infof("Could not complete request %s %s: %v", c.Request.Method, c.Request.RequestURI, err)
	if errors.Is(err, errors2.ErrProjectAlreadyExists) {
		SetConflictErrorResponse(c, "Project already exists")
	} else if errors.Is(err, errors2.ErrStageAlreadyExists) || errors.Is(err, errors2.ErrBranchExists) {
		SetConflictErrorResponse(c, "Stage already exists")
	} else if errors.Is(err, errors2.ErrServiceAlreadyExists) {
		SetConflictErrorResponse(c, "Service already exists")
	} else if errors.Is(err, errors2.ErrInvalidGitToken) {
		SetFailedDependencyErrorResponse(c, "Invalid git token")
	} else if errors.Is(err, errors2.ErrCredentialsTokenMustNotBeEmpty) {
		SetBadRequestErrorResponse(c, "Upstream repository not found")
	} else if errors.Is(err, errors2.ErrCredentialsNotFound) {
		SetNotFoundErrorResponse(c, "Could not find credentials for upstream repository")
	} else if errors.Is(err, errors2.ErrMalformedCredentials) {
		SetFailedDependencyErrorResponse(c, "Could not decode credentials for upstream repository")
	} else if errors.Is(err, errors2.ErrCredentialsInvalidRemoteURI) {
		SetBadRequestErrorResponse(c, "Upstream repository not found")
	} else if errors.Is(err, errors2.ErrRepositoryNotFound) {
		SetNotFoundErrorResponse(c, "Upstream repository not found")
	} else if errors.Is(err, errors2.ErrProjectNotFound) {
		SetNotFoundErrorResponse(c, "Project not found")
	} else if errors.Is(err, errors2.ErrStageNotFound) || errors.Is(err, errors2.ErrReferenceNotFound) {
		SetNotFoundErrorResponse(c, "Stage not found")
	} else if errors.Is(err, errors2.ErrServiceNotFound) {
		SetNotFoundErrorResponse(c, "Service not found")
	} else if errors.Is(err, errors2.ErrResourceNotFound) {
		SetNotFoundErrorResponse(c, "Resource not found")
	} else {
		logger.Errorf("Encountered unknown error: %v", err)
		SetInternalServerErrorResponse(c, "Internal server error")
	}
}

func SetFailedDependencyErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusFailedDependency, models.Error{
		Code:    http.StatusFailedDependency,
		Message: msg,
	})
}

func SetNotFoundErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.Error{
		Code:    http.StatusNotFound,
		Message: msg,
	})
}

func SetInternalServerErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}

func SetBadRequestErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.Error{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

func SetConflictErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, models.Error{
		Code:    http.StatusConflict,
		Message: msg,
	})
}
