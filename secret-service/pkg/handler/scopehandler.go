package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/model"
)

type IScopeHandler interface {
	GetScopes(c *gin.Context)
}

type ScopeHandler struct {
	ScopeBackend backend.ScopeManager
}

func NewScopeHandler(backend backend.ScopeManager) *ScopeHandler {
	return &ScopeHandler{
		ScopeBackend: backend,
	}
}

// GetScopes godoc
// @Summary      Get scopes
// @Description  Get scopes
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}secrets:read</span>
// @Tags         Scopes
// @Security     ApiKeyAuth
// @Success      200  {object}  model.GetScopesResponse  "OK"
// @Failure      500  {object}  model.Error              "Internal Server Error"
// @Router       /scope [get]
func (s ScopeHandler) GetScopes(c *gin.Context) {
	scopes, err := s.ScopeBackend.GetScopes()
	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(ErrGetScopesMsg, err.Error()))
		return
	}

	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, model.GetScopesResponse{Scopes: scopes})
}
