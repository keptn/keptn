package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/handler"
)

const SecretAPIBasePath = "/secret"

type SecretController struct {
	SecretHandler handler.ISecretHandler
}

func NewSecretController(secretHandler handler.ISecretHandler) *SecretController {
	return &SecretController{SecretHandler: secretHandler}
}

func (controller SecretController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST(SecretAPIBasePath, controller.SecretHandler.CreateSecret)
	apiGroup.DELETE(SecretAPIBasePath, controller.SecretHandler.DeleteSecret)
	apiGroup.PUT(SecretAPIBasePath, controller.SecretHandler.UpdateSecret)
}
