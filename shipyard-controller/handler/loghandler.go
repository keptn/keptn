package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type ILogHandler interface {
	CreateLogEntries(context *gin.Context)
	GetLogEntries(context *gin.Context)
	DeleteLogEntries(context *gin.Context)
}

type LogHandler struct {
	logManager ILogManager
}

func NewLogHandler(logManager ILogManager) *LogHandler {
	return &LogHandler{logManager: logManager}
}

// CreateLogEntries persists log entries
// @Summary      Persist a list of log entries
// @Description  Persist a list of log entries
// @Tags         Log
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        integration  body  models.CreateLogsRequest  true  "Logs"
// @Success      200
// @Failure      400  {object}  models.Error  "Invalid payload"
// @Failure      500  {object}  models.Error  "Internal error"
// @Router       /log [post]
func (lh *LogHandler) CreateLogEntries(context *gin.Context) {
	logs := &models.CreateLogsRequest{}
	if err := context.ShouldBindJSON(logs); err != nil {
		SetBadRequestErrorResponse(context, err.Error())
		return
	}

	if err := lh.logManager.CreateLogEntries(*logs); err != nil {
		SetInternalServerErrorResponse(context, err.Error())
		return
	}
	context.JSON(http.StatusOK, models.CreateLogsResponse{})
}

// GetLogEntries Retrieves log entries based on the provided filter
// @Summary      Retrieve logs
// @Description  Retrieve logs
// @Tags         Log
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        integrationId  query     string                  false  "integrationId"
// @Param        fromTime       query     string                  false  "The from time stamp for fetching sequence states"
// @Param        beforeTime     query     string                  false  "The before time stamp for fetching sequence states"
// @Param        pageSize       query     int                     false  "The number of items to return"
// @Param        nextPageKey    query     string                  false  "Pointer to the next set of items"
// @Success      200            {object}  models.GetLogsResponse  "ok"
// @Failure      400            {object}  models.Error            "Invalid payload"
// @Failure      404            {object}  models.Error            "Not found"
// @Router       /log [get]
func (lh *LogHandler) GetLogEntries(context *gin.Context) {
	params := &models.GetLogParams{}
	if err := context.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(context, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	logs, err := lh.logManager.GetLogEntries(*params)
	if err != nil {
		SetNotFoundErrorResponse(context, fmt.Sprintf(UnableRetrieveLogsMsg, err.Error()))
		return
	}
	context.JSON(http.StatusOK, logs)
}

// DeleteLogEntries
// @Summary      Delete logs
// @Deprecated   true
// @Description  INTERNAL Endpoint: Delete logs
// @Tags         Log
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        integrationId  query     string                    false  "integrationId"
// @Param        fromTime       query     string                    false  "The from time stamp for fetching sequence states"
// @Param        beforeTime     query     string                    false  "The before time stamp for fetching sequence states"
// @Success      200            {object}  models.DeleteLogResponse  "ok"
// @Failure      400            {object}  models.Error              "Invalid payload"
// @Failure      500            {object}  models.Error              "Internal error"
// @Router       /log [delete]
func (lh *LogHandler) DeleteLogEntries(context *gin.Context) {
	params := &models.DeleteLogParams{}
	if err := context.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(context, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	if err := lh.logManager.DeleteLogEntries(*params); err != nil {
		SetInternalServerErrorResponse(context, fmt.Sprintf(UnableRetrieveLogsMsg, err.Error()))
		return
	}
	context.JSON(http.StatusOK, models.DeleteLogParams{})
}
