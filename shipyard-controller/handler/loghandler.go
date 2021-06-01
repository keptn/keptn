package handler

import (
	"github.com/gin-gonic/gin"
)

type ILogHandler interface {
	CreateLogEntries(context *gin.Context)
	GetLogEntries(context *gin.Context)
}

type LogHandler struct {
	logManager ILogManager
}

func NewLogHandler(logManager ILogManager) *LogHandler {
	return &LogHandler{logManager: logManager}
}

// CreateLogEntries persists log entries
// @Summary Persist a list of log entries
// @Description Persist a list of log entries
// @Tags Log
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integration body models.CreateLogsRequest true "Logs"
// @Success 200
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /log [post]
func (LogHandler) CreateLogEntries(context *gin.Context) {
	panic("implement me")
}

// GetLogEntries Retrieves log entries based on the provided filter
// @Summary Retrieve logs
// @Description Retrieve logs
// @Tags Log
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integrationId query string false "integrationId"
// @Param	fromTime			query	string	false	"The from time stamp for fetching sequence states"
// @Param 	beforeTime			query	string	false	"The before time stamp for fetching sequence states"
// @Param	pageSize			query	int		false	"The number of items to return"
// @Param   nextPageKey     	query   string  false	"Pointer to the next set of items"
// @Success 200 {object} models.GetLogsResponse "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /log [get]
func (LogHandler) GetLogEntries(context *gin.Context) {
	panic("implement me")
}
