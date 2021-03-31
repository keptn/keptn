package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"net/http"
	"strings"
)

func SetBadRequestErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusBadRequest, model.Error{
		Code:    http.StatusBadRequest,
		Message: &msg,
	})
}

func SetInternalServerErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusInternalServerError, model.Error{
		Code:    http.StatusInternalServerError,
		Message: &msg,
	})
}

func SetConflictErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusConflict, model.Error{
		Code:    http.StatusConflict,
		Message: &msg,
	})
}

func SetNotFoundErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusNotFound, model.Error{
		Code:    http.StatusNotFound,
		Message: &msg,
	})
}

func errMsg(err error, message []string) string {
	var sb strings.Builder

	if len(message) > 0 {
		if err != nil {
			sb.WriteString(fmt.Sprintf("%s: %s", message[0], err.Error()))
		} else {
			sb.WriteString(message[0])
		}
	} else {
		if err != nil {
			sb.WriteString(err.Error())
		}
	}
	return sb.String()
}
