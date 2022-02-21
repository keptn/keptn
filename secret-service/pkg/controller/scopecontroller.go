package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/handler"
)

const ScopeAPIBasePath = "/scope"

type ScopeController struct {
	ScopeHandler handler.IScopeHandler
}

func NewScopeController(scopeHandler handler.IScopeHandler) *ScopeController {
	return &ScopeController{ScopeHandler: scopeHandler}
}

func (controller ScopeController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET(ScopeAPIBasePath, controller.ScopeHandler.GetScopes)
}

