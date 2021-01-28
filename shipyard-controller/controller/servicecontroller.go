package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type ServiceController struct {
	ServiceHandler handler.IServiceHandler
}

func NewServiceController(serviceHandler handler.IServiceHandler) Controller {
	return &ServiceController{ServiceHandler: serviceHandler}
}

func (controller ServiceController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:project/service", controller.ServiceHandler.CreateService)
	apiGroup.PUT("/project/:project/service/:service", controller.ServiceHandler.DeleteService)
}
