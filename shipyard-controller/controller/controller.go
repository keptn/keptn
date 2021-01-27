package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	Inject(engine *gin.Engine)
}
