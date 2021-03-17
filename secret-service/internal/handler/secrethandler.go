package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/internal/backend"
	"github.com/keptn/keptn/secret-service/internal/model"
	"net/http"
)

type ISecretHandler interface {
	CreateSecret(c *gin.Context)
	UpdateSecret(c *gin.Context)
	DeleteSecret(c *gin.Context)
}

func NewSecretHandler(store backend.SecretStore) *SecretHandler {
	return &SecretHandler{
		SecretBackend: store,
	}
}

type SecretHandler struct {
	SecretBackend backend.SecretStore
}

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

	c.Status(http.StatusCreated)
}

func (s SecretHandler) UpdateSecret(c *gin.Context) {
	panic("implement me")
}

func (s SecretHandler) DeleteSecret(c *gin.Context) {
	panic("implement me")
}
