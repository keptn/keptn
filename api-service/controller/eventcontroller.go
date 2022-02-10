package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/api-service/handler"
)

type EventController struct {
	EventHandler handler.IEventHandler
}

func NewEventController(eventHandler handler.IEventHandler) *EventController {
	return &EventController{EventHandler: eventHandler}
}

func (controller EventController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/event", controller.EventHandler.ForwardEvent)
	apiGroup.GET("/event", controller.EventHandler.GetEvent)
}
