package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/model"
)

var ErrCreateSecretMsg = "Unable to create secret: %s"
var ErrInvalidRequestFormatMsg = "Invalid request format: %s"
var ErrUpdateSecretMsg = "Unable to update secret: %s"
var ErrGetSecretMsg = "Unable to get secret: %s"
var ErrDeleteSecretMsg = "Unable to delete secret: %s"
var ErrGetScopesMsg = "Unable to get scopes: %s"

func SetBadRequestErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, model.Error{
		Code:    http.StatusBadRequest,
		Message: &msg,
	})
}

func SetInternalServerErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, model.Error{
		Code:    http.StatusInternalServerError,
		Message: &msg,
	})
}

func SetConflictErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, model.Error{
		Code:    http.StatusConflict,
		Message: &msg,
	})
}

func SetNotFoundErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, model.Error{
		Code:    http.StatusNotFound,
		Message: &msg,
	})
}
