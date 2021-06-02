package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type HealthController struct {
	HealthHandler handler.IHealthHandler
}

func NewHealthController(healthHandler handler.IHealthHandler) Controller {
	return &HealthController{HealthHandler: healthHandler}
}

func (controller HealthController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/health", controller.HealthHandler.Health)
}
