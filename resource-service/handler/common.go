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

	if check, resourceType := alreadyExists(err); check {
		SetConflictErrorResponse(c, resourceType+" already exists")
	} else if errors.Is(err, errors2.ErrInvalidGitToken) {
		SetFailedDependencyErrorResponse(c, "Invalid git token")
	} else if errors.Is(err, errors2.ErrCredentialsNotFound) {
		SetNotFoundErrorResponse(c, "Could not find credentials for upstream repository")
	} else if errors.Is(err, errors2.ErrMalformedCredentials) {
		SetFailedDependencyErrorResponse(c, "Could not decode credentials for upstream repository")
	} else if errors.Is(err, errors2.ErrCredentialsInvalidRemoteURI) || errors.Is(err, errors2.ErrCredentialsTokenMustNotBeEmpty) {
		SetBadRequestErrorResponse(c, "Upstream repository not found")
	} else if errors.Is(err, errors2.ErrRepositoryNotFound) {
		SetNotFoundErrorResponse(c, "Upstream repository not found")
	} else if check, resourceType := resourceNotFound(err); check {
		SetNotFoundErrorResponse(c, resourceType+" not found")
	} else {
		logger.Errorf("Encountered unknown error: %v", err)
		SetInternalServerErrorResponse(c, "Internal server error")
	}
}

func alreadyExists(err error) (bool, string) {
	if errors.Is(err, errors2.ErrProjectAlreadyExists) {
		return true, "Project"
	} else if errors.Is(err, errors2.ErrStageAlreadyExists) || errors.Is(err, errors2.ErrBranchExists) {
		return true, "Stage"
	} else if errors.Is(err, errors2.ErrServiceAlreadyExists) {
		return true, "Service"
	}
	return false, ""
}

func resourceNotFound(err error) (bool, string) {
	if errors.Is(err, errors2.ErrProjectNotFound) {
		return true, "Project"
	} else if errors.Is(err, errors2.ErrStageNotFound) || errors.Is(err, errors2.ErrReferenceNotFound) {
		return true, "Stage"
	} else if errors.Is(err, errors2.ErrServiceNotFound) {
		return true, "Service"
	} else if errors.Is(err, errors2.ErrResourceNotFound) {
		return true, "Resource"
	}
	return false, ""
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
