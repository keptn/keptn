package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type UniformIntegrationController struct {
	UniformIntegrationHandler handler.IUniformIntegrationHandler
}

func NewUniformIntegrationController(uniformIntegrationHandler handler.IUniformIntegrationHandler) Controller {
	return &UniformIntegrationController{UniformIntegrationHandler: uniformIntegrationHandler}
}

func (controller UniformIntegrationController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/uniform/registration", controller.UniformIntegrationHandler.GetRegistrations)
	apiGroup.PUT("/uniform/registration/:id", controller.UniformIntegrationHandler.KeepAlive)
	apiGroup.POST("/uniform/registration", controller.UniformIntegrationHandler.Register)
	apiGroup.DELETE("/uniform/registration/:id", controller.UniformIntegrationHandler.Unregister)

}
