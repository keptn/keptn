package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/internal/handler"
)

type SecretController struct {
	SecretHandler handler.ISecretHandler
}

func NewSecretController(secretHandler handler.ISecretHandler) *SecretController {
	return &SecretController{SecretHandler: secretHandler}
}

func (controller SecretController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/secrets", controller.SecretHandler.CreateSecret)
}
