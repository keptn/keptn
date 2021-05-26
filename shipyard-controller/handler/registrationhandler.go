package handler

import "github.com/gin-gonic/gin"

type IRegistrationHandler interface {
	CreateRegistration(context *gin.Context)
	DeleteRegistration(context *gin.Context)
	GetRegistrations(context *gin.Context)
}

type IRegistrationManager interface {
}

type RegistrationHandler struct {
	registrationManager IRegistrationManager
}

// CreateRegistration registers a uniform integration
// @Summary Register a uniform integration
// @Description Register a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param integration body models.Integration true "Integration"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration [post]
func (RegistrationHandler) CreateRegistration(context *gin.Context) {
	panic("implement me")
}

// DeleteRegistration Unregisters a uniform integration
// @Summary Unregister a uniform integration
// @Description Unregister a uniform integration
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration/{id} [delete]
func (RegistrationHandler) DeleteRegistration(context *gin.Context) {
	panic("implement me")
}

// GetRegistrations Retrieves uniform integrations matching the provided filter
// @Summary Retrieve uniform integrations
// @Description Retrieve uniform integrations
// @Tags Uniform
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string false "id"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /uniform/registration [get]
func (RegistrationHandler) GetRegistrations(context *gin.Context) {
	panic("implement me")
}
