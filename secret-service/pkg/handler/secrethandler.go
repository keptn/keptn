package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/model"
)

type ISecretHandler interface {
	CreateSecret(c *gin.Context)
	UpdateSecret(c *gin.Context)
	DeleteSecret(c *gin.Context)
	GetSecrets(c *gin.Context)
}

func NewSecretHandler(backend backend.SecretManager) *SecretHandler {
	return &SecretHandler{
		SecretManager: backend,
	}
}

type SecretHandler struct {
	SecretManager backend.SecretManager
}

// CreateSecret godoc
// @Summary      Create a Secret
// @Description  Create a new Secret
// @Tags         Secrets
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        secret  body      model.Secret  true  "The new secret to be created"
// @Success      201     {object}  model.Secret  "Created"
// @Failure      400     {object}  model.Error   "Invalid Payload"
// @Failure      409     {object}  model.Error   "Conflict"
// @Failure      500     {object}  model.Error   "Internal Server Error"
// @Router       /secret [post]
func (s SecretHandler) CreateSecret(c *gin.Context) {
	secret := model.Secret{}
	if err := c.ShouldBindJSON(&secret); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(ErrInvalidRequestFormatMsg, err.Error()))
		return
	}

	if secret.Scope == "" {
		secret.Scope = model.DefaultSecretScope
	}

	err := s.SecretManager.CreateSecret(secret)
	if err != nil {
		if errors.Is(err, backend.ErrSecretAlreadyExists) {
			SetConflictErrorResponse(c, fmt.Sprintf(ErrCreateSecretMsg, err.Error()))
			return
		}
		if errors.Is(err, backend.ErrTooBigKeySize) || errors.Is(err, backend.ErrScopeNotFound) {
			SetBadRequestErrorResponse(c, fmt.Sprintf(ErrCreateSecretMsg, err.Error()))
			return
		}
		SetInternalServerErrorResponse(c, fmt.Sprintf(ErrCreateSecretMsg, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, secret)
}

// CreateSecret godoc
// @Summary      Update a Secret
// @Description  Update an existing Secret
// @Tags         Secrets
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        secret  body      model.Secret  true  "The updated Secret"
// @Success      200     {object}  model.Secret  "OK"
// @Failure      400     {object}  model.Error   "Invalid payload"
// @Failure      404     {object}  model.Error   "Not Found"
// @Failure      500     {object}  model.Error   "Internal Server Error"
// @Router       /secret [put]
func (s SecretHandler) UpdateSecret(c *gin.Context) {
	secret := model.Secret{}
	if err := c.ShouldBindJSON(&secret); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(ErrInvalidRequestFormatMsg, err.Error()))
		return
	}

	err := s.SecretManager.UpdateSecret(secret)
	if err != nil {
		if errors.Is(err, backend.ErrSecretNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ErrUpdateSecretMsg, err.Error()))
			return
		}
		if errors.Is(err, backend.ErrScopeNotFound) {
			SetBadRequestErrorResponse(c, fmt.Sprintf(ErrUpdateSecretMsg, err.Error()))
			return
		}
		SetInternalServerErrorResponse(c, fmt.Sprintf(ErrUpdateSecretMsg, err.Error()))
		return
	}
	c.JSON(http.StatusOK, secret)

}

// CreateSecret godoc
// @Summary      Delete a Secret
// @Description  Delete an existing Secret
// @Tags         Secrets
// @Security     ApiKeyAuth
// @Param        name   query  string  true  "The name of the secret"
// @Param        scope  query  string  true  "The scope of the secret"
// @Success      200    "OK"
// @Failure      400    {object}  model.Error  "Invalid payload"
// @Failure      404    {object}  model.Error  "Not Found"
// @Failure      500    {object}  model.Error  "Internal Server Error"
// @Router       /secret [delete]
func (s SecretHandler) DeleteSecret(c *gin.Context) {
	params := &DeleteSecretQueryParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(ErrInvalidRequestFormatMsg, err.Error()))
		return
	}

	secret := model.Secret{
		SecretMetadata: model.SecretMetadata{
			Name:  params.Name,
			Scope: params.Scope,
		},
		Data: nil,
	}
	err := s.SecretManager.DeleteSecret(secret)
	if err != nil {
		if errors.Is(err, backend.ErrSecretNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(ErrDeleteSecretMsg, err.Error()))
			return
		}
		if errors.Is(err, backend.ErrScopeNotFound) {
			SetBadRequestErrorResponse(c, fmt.Sprintf(ErrDeleteSecretMsg, err.Error()))
			return
		}
		SetInternalServerErrorResponse(c, fmt.Sprintf(ErrDeleteSecretMsg, err.Error()))
		return
	}

	c.Status(http.StatusOK)

}

// GetSecrets godoc
// @Summary      Get secrets
// @Description  Get secrets
// @Tags         Secrets
// @Security     ApiKeyAuth
// @Success      200  {object}  model.GetSecretsResponse  "OK"
// @Failure      500  {object}  model.Error               "Internal Server Error"
// @Router       /secret [get]
func (s SecretHandler) GetSecrets(c *gin.Context) {
	secrets, err := s.SecretManager.GetSecrets()
	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(ErrGetSecretMsg, err.Error()))
		return
	}

	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, model.GetSecretsResponse{Secrets: secrets})
}

type DeleteSecretQueryParams struct {
	Name  string `form:"name" binding:"required"`
	Scope string `form:"scope" binding:"required"`
}
