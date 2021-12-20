package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type ServiceController struct {
	ServiceHandler handler.IServiceHandler
}

func NewServiceController(serviceHandler handler.IServiceHandler) Controller {
	return &ServiceController{ServiceHandler: serviceHandler}
}

func (controller ServiceController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:projectName/stage/:stageName/service", controller.ServiceHandler.CreateService)
	apiGroup.PUT("/project/:projectName/stage/:stageName/service/:serviceName", controller.ServiceHandler.UpdateService)
	apiGroup.DELETE("/project/:projectName/stage/:stageName/service/:serviceName", controller.ServiceHandler.DeleteService)
}
