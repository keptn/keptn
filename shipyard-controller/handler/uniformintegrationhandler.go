package handler

import (
	"github.com/gin-gonic/gin"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"net/http"
)

type IUniformIntegrationHandler interface {
	Register(context *gin.Context)
	Unregister(context *gin.Context)
	GetRegistrations(context *gin.Context)
}

type UniformIntegrationHandler struct {
	integrationManager IUniformIntegrationManager
}

func NewUniformIntegrationHandler(im IUniformIntegrationManager) *UniformIntegrationHandler {
	return &UniformIntegrationHandler{integrationManager: im}
}

// CreateRegistration registers a uniform integration
// @Summary Register a uniform integration
// @Description Register a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integration body models.Integration true "Integration"
// @Success 200
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration [post]
func (rh *UniformIntegrationHandler) Register(c *gin.Context) {
	integration := &models.Integration{}

	if err := c.ShouldBindJSON(integration); err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}

	integrationID := keptnmodels.IntegrationID{
		Name:      integration.Name,
		Namespace: integration.MetaData.KubernetesMetaData.Namespace,
		Project:   integration.Subscription.Filter.Project,
		Stage:     integration.Subscription.Filter.Stage,
		Service:   integration.Subscription.Filter.Service,
	}

	hash, err := integrationID.Hash()
	if err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}

	integration.ID = hash

	if err := rh.integrationManager.Register(*integration); err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}
	c.JSON(http.StatusOK, &models.RegisterResponse{
		ID: integration.ID,
	})
}

// DeleteRegistration Unregisters a uniform integration
// @Summary Unregister a uniform integration
// @Description Unregister a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration/{id} [delete]
func (rh *UniformIntegrationHandler) Unregister(c *gin.Context) {
	integrationID := c.Param("id")

	if err := rh.integrationManager.Unregister(integrationID); err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}
	c.JSON(http.StatusOK, &models.UnregisterResponse{})
}

// GetRegistrations Retrieves uniform integrations matching the provided filter
// @Summary Retrieve uniform integrations
// @Description Retrieve uniform integrations
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string false "id"
// @Param name query string false "name"
// @Param project query string false "project"
// @Param stage query string false "stage"
// @Param service query string false "service"
// @Success 200 {object} []models.Integration "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration [get]
func (rh *UniformIntegrationHandler) GetRegistrations(c *gin.Context) {
	params := &models.GetUniformIntegrationParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}
	uniformIntegrations, err := rh.integrationManager.GetRegistrations(*params)
	if err != nil {
		SetInternalServerErrorResponse(err, c, "Unable to query uniform integrations repository")
		return
	}

	c.JSON(http.StatusOK, uniformIntegrations)
}
