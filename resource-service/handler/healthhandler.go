package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IHealthHandler interface {
	Health(context *gin.Context)
}

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
