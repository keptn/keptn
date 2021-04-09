package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/handler"
)

var k string

type SecretController struct {
	SecretHandler handler.ISecretHandler
}

func NewSecretController(secretHandler handler.ISecretHandler) *SecretController {
	return &SecretController{SecretHandler: secretHandler}
}

func (controller SecretController) Inject(apiGroup *gin.RouterGroup) {

	apiGroup.POST("/secrets", controller.SecretHandler.CreateSecret)
	apiGroup.DELETE("/secrets", controller.SecretHandler.DeleteSecret)
	apiGroup.PUT("/secrets", controller.SecretHandler.UpdateSecret)
}
