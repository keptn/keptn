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
		if err == backend.ErrSecretAlreadyExists {
			SetConflictErrorResponse(err, c, "Unable to create secret")
			return
		}
		SetInternalServerErrorResponse(err, c, "Unable to create secret")
		return
	}

	c.JSON(http.StatusCreated, secret)
}

// CreateSecret godoc
// @Summary Update a Secret
// @Description Update an existing Secret
// @Tags Secrets
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param secret body model.Secret true "The updated Secret"
// @Success 200 {object} model.Secret
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /secrets [put]
func (s SecretHandler) UpdateSecret(c *gin.Context) {
	secret := model.Secret{}
	if err := c.ShouldBindJSON(&secret); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	err := s.SecretBackend.UpdateSecret(secret)
	if err != nil {
		if err == backend.ErrSecretNotFound {
			SetNotFoundErrorResponse(err, c, "Unable to update secret")
			return
		}
		SetInternalServerErrorResponse(err, c, "Unable to update secret")
		return
	}
	c.JSON(http.StatusOK, secret)

}

// CreateSecret godoc
// @Summary Delete a Secret
// @Description Delete an existing Secret
// @Tags Secrets
// @Security ApiKeyAuth
// @Param name query string true "The name of the secret"
// @Param scope query string true "The scope of the secret"
// @Success 200
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /secrets [delete]
func (s SecretHandler) DeleteSecret(c *gin.Context) {
	params := &DeleteSecretQueryParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	secret := model.Secret{
		Name:  params.Name,
		Scope: params.Scope,
		Data:  nil,
	}
	err := s.SecretBackend.DeleteSecret(secret)
	if err != nil {
		if err == backend.ErrSecretNotFound {
			SetNotFoundErrorResponse(err, c, "Unable to delete secret")
			return
		}
		SetInternalServerErrorResponse(err, c, "Unable to delete secret")
		return
	}

	c.Status(http.StatusOK)

}

type DeleteSecretQueryParams struct {
	Name  string `form:"name"`
	Scope string `form:"scope"`
}
