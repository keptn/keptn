package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/api-service/handler"
)

type AuthController struct {
	AuthHandler handler.IAuthHandler
}

func NewAuthController(authHandler handler.IAuthHandler) *AuthController {
	return &AuthController{AuthHandler: authHandler}
}

func (controller AuthController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/auth", controller.AuthHandler.VerifyToken)
}
