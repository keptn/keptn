package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"net/http"
)

type IScopeHandler interface {
	GetScopes(c *gin.Context)
}

type ScopeHandler struct {
	ScopeBackend backend.ScopeBackend
}

func NewScopeHandler(backend backend.ScopeBackend) *ScopeHandler {
	return &ScopeHandler{
		ScopeBackend: backend,
	}
}


// GetScopes godoc
// @Summary Get scopes
// @Description Get scopes
// @Tags Scopes
// @Security ApiKeyAuth
// @Success 200 {object} model.GetScopesResponse
// @Failure 500 {object} model.Error
// @Router /scope [get]
func (s ScopeHandler) GetScopes(c *gin.Context) {
	scopes, err := s.ScopeBackend.GetScopes()
	if err != nil {
		SetInternalServerErrorResponse(err, c, "Unable to get scopes")
		return
	}

	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, model.GetScopesResponse{Scopes: scopes})
}
