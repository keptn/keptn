package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type LogController struct {
	LogHandler handler.ILogHandler
}

func NewLogController(logHandler handler.ILogHandler) *LogController {
	return &LogController{LogHandler: logHandler}
}

func (controller LogController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/log", controller.LogHandler.GetLogEntries)
	apiGroup.POST("/log", controller.LogHandler.CreateLogEntries)
	apiGroup.DELETE("/log", controller.LogHandler.DeleteLogEntries)
}
