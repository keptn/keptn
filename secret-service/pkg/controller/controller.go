package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	Inject(apiGroup *gin.RouterGroup)
}
