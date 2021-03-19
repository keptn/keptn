package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"net/http"
)

type ISecretHandler interface {
	CreateSecret(c *gin.Context)
	UpdateSecret(c *gin.Context)
	DeleteSecret(c *gin.Context)
}

func NewSecretHandler(backend backend.SecretBackend) *SecretHandler {
	return &SecretHandler{
		SecretBackend: backend,
	}
}

type SecretHandler struct {
	SecretBackend backend.SecretBackend
}

// CreateSecret godoc
// @Summary Create a Secret
// @Description Create a new Secret
// @Tags Secrets
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param secret body model.Secret true "The new secret to be created"
// @Success 200 {object} model.Secret
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /secrets [post]
func (s SecretHandler) CreateSecret(c *gin.Context) {

	secret := model.Secret{}
	if err := c.ShouldBindJSON(&secret); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	err := s.SecretBackend.CreateSecret(secret)
	if err != nil {
		SetInternalServerErrorResponse(err, c, "Unable to create secret")
		return
	}

	c.JSON(http.StatusCreated, secret)
}

func (s SecretHandler) UpdateSecret(c *gin.Context) {
	panic("implement me")
}

func (s SecretHandler) DeleteSecret(c *gin.Context) {
	panic("implement me")
}
