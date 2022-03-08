package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/model"
)

var ErrCreateSecret = "Unable to create secret: %s"
var ErrInvalidRequestFormat = "Invalid request format: %s"
var ErrUpdateSecret = "Unable to update secret: %s"
var ErrGetSecret = "Unable to get secret: %s"
var ErrDeleteSecret = "Unable to delete secret: %s"
var ErrGetScopes = "Unable to get scopes: %s"

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
