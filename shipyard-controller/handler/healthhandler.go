package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/db"
	"net/http"
)

type IHealthHandler interface {
	Health(context *gin.Context)
}

type HealthHandler struct {
	db.MongoDBConnection
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
