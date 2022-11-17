package middleware

import (
	openapierrors "github.com/go-openapi/errors"
	"github.com/keptn/keptn/api/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

//go:generate moq -pkg middleware_mock --skip-ensure -out ./fake/tokenvalidator_mock.go . TokenValidator
type TokenValidator interface {
	ValidateToken(token string) (*models.Principal, error)
}

type BasicTokenValidator struct{}

func (b *BasicTokenValidator) ValidateToken(token string) (*models.Principal, error) {
	if token == os.Getenv("SECRET_TOKEN") {
		prin := models.Principal(token)
		return &prin, nil
	}
	log.Warn("Access attempt with incorrect API token")
	return nil, openapierrors.New(http.StatusUnauthorized, "incorrect api key auth")
}
