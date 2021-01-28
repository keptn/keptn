package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type EventController struct {
	EventHandler handler.IEventHandler
}

func NewEventController(eventHandler handler.IEventHandler) Controller {
	return &EventController{EventHandler: eventHandler}
}

func (controller EventController) Inject(engine *gin.Engine) {
	engine.GET("/event/triggered/:eventType", controller.EventHandler.GetTriggeredEvents)
	engine.POST("/event", controller.EventHandler.HandleEvent)
}
