package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type ServiceDefaultResourceController struct {
	ServiceDefaultResourceHandler handler.IServiceDefaultResourceHandler
}

func NewServiceDefaultResourceController(serviceDefaultResourceHandler handler.IServiceDefaultResourceHandler) Controller {
	return &ServiceDefaultResourceController{ServiceDefaultResourceHandler: serviceDefaultResourceHandler}
}

func (controller ServiceDefaultResourceController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:projectName/service/:serviceName/resource", controller.ServiceDefaultResourceHandler.CreateServiceDefaultResources)
	apiGroup.GET("/project/:projectName/service/:serviceName/resource", controller.ServiceDefaultResourceHandler.GetServiceDefaultResources)
	apiGroup.PUT("/project/:projectName/service/:serviceName/resource", controller.ServiceDefaultResourceHandler.UpdateServiceDefaultResources)
	apiGroup.GET("/project/:projectName/service/:serviceName/resource/:resourceURI", controller.ServiceDefaultResourceHandler.GetServiceDefaultResource)
	apiGroup.PUT("/project/:projectName/service/:serviceName/resource/:resourceURI", controller.ServiceDefaultResourceHandler.UpdateServiceDefaultResource)
	apiGroup.DELETE("/project/:projectName/service/:serviceName/resource/:resourceURI", controller.ServiceDefaultResourceHandler.DeleteServiceDefaultResource)
}
