package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventsender.go . IEventSender
type IEventSender interface {
	SendEvent(event cloudevents.Event) error
	Send(ctx context.Context, event cloudevents.Event) error
}

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

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

func SetUnprocessableEntityResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusUnprocessableEntity, models.Error{
		Code:    http.StatusUnprocessableEntity,
		Message: &msg,
	})
}

func DecodeInputData(body io.ReadCloser, params any) error {
	jsonData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	d := json.NewDecoder(strings.NewReader(string(jsonData)))
	d.DisallowUnknownFields()

	if err := d.Decode(&params); err != nil {
		return err
	}
	return nil
}

func mapError(c *gin.Context, err error) {
	if errors.Is(err, common.ErrProjectAlreadyExists) {
		SetConflictErrorResponse(c, err.Error())
		return
	}
	if err.Error() == common.AlreadyInitializedRepositoryMsg {
		SetConflictErrorResponse(c, err.Error())
		return
	}
	if errors.Is(err, common.ErrConfigStoreUpstreamNotFound) {
		SetNotFoundErrorResponse(c, err.Error())
		return
	}
	if errors.Is(err, common.ErrConfigStoreInvalidToken) {
		SetFailedDependencyErrorResponse(c, err.Error())
		return
	}
	if errors.Is(err, common.ErrProjectNotFound) {
		SetNotFoundErrorResponse(c, err.Error())
		return
	}
	if errors.Is(err, common.ErrInvalidStageChange) {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}
	SetInternalServerErrorResponse(c, err.Error())
	return
}
