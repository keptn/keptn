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
	apiGroup.PUT("/uniform/registration/:integrationID/ping", controller.UniformIntegrationHandler.KeepAlive)
	apiGroup.POST("/uniform/registration", controller.UniformIntegrationHandler.Register)
	apiGroup.POST("/uniform/registration/:integrationID/subscription", controller.UniformIntegrationHandler.CreateSubscription)
	apiGroup.PUT("/uniform/registration/:integrationID/subscription/:subscriptionID", controller.UniformIntegrationHandler.UpdateSubscription)
	apiGroup.DELETE("/uniform/registration/:integrationID", controller.UniformIntegrationHandler.Unregister)
	apiGroup.DELETE("/uniform/registration/:integrationID/subscription/:subscriptionID", controller.UniformIntegrationHandler.DeleteSubscription)

}
