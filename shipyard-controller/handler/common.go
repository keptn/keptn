package handler

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/models"
	"net/http"
	"strings"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventsender.go . IEventSender
type IEventSender interface {
	SendEvent(event cloudevents.Event) error
	Send(ctx context.Context, event cloudevents.Event) error
}

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

const invalidRequestFormatMsg = "Invalid request format"

func SetFailedDependencyErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusFailedDependency, models.Error{
		Code:    http.StatusFailedDependency,
		Message: &msg,
	})
}

func SetNotFoundErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusNotFound, models.Error{
		Code:    http.StatusNotFound,
		Message: &msg,
	})
}

func SetInternalServerErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    http.StatusInternalServerError,
		Message: &msg,
	})
}

func SetBadRequestErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusBadRequest, models.Error{
		Code:    http.StatusBadRequest,
		Message: &msg,
	})
}

func SetConflictErrorResponse(err error, c *gin.Context, message ...string) {
	msg := errMsg(err, message)
	c.JSON(http.StatusConflict, models.Error{
		Code:    http.StatusConflict,
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
