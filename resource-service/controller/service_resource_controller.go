package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type ServiceResourceController struct {
	ServiceResourceHandler handler.IServiceResourceHandler
}

func NewServiceResourceController(serviceResourceHandler handler.IServiceResourceHandler) Controller {
	return &ServiceResourceController{ServiceResourceHandler: serviceResourceHandler}
}

func (controller ServiceResourceController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:projectName/stage/:stageName/service/:serviceName/resource", controller.ServiceResourceHandler.CreateServiceResources)
	apiGroup.GET("/project/:projectName/stage/:stageName/service/:serviceName/resource", controller.ServiceResourceHandler.GetServiceResources)
	apiGroup.PUT("/project/:projectName/stage/:stageName/service/:serviceName/resource", controller.ServiceResourceHandler.UpdateServiceResources)
	apiGroup.GET("/project/:projectName/stage/:stageName/service/:serviceName/resource/:resourceURI", controller.ServiceResourceHandler.GetServiceResource)
	apiGroup.PUT("/project/:projectName/stage/:stageName/service/:serviceName/resource/:resourceURI", controller.ServiceResourceHandler.UpdateServiceResource)
	apiGroup.DELETE("/project/:projectName/stage/:stageName/service/:serviceName/resource/:resourceURI", controller.ServiceResourceHandler.DeleteServiceResource)
}
