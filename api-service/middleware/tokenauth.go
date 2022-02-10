package middleware

import (
	"net/http"
	"os"

	openapierrors "github.com/go-openapi/errors"
	"github.com/keptn/keptn/api-service/model"
	log "github.com/sirupsen/logrus"
)

//go:generate moq -pkg middleware_mock --skip-ensure -out ./fake/tokenvalidator_mock.go . TokenValidator
type TokenValidator interface {
	ValidateToken(token string) (*model.Principal, error)
}

type BasicTokenValidator struct{}

func (b *BasicTokenValidator) ValidateToken(token string) (*model.Principal, error) {
	if token == os.Getenv("SECRET_TOKEN") {
		prin := model.Principal(token)
		return &prin, nil
	}
	log.Errorf("Access attempt with incorrect api key auth: %s", token)
	return nil, openapierrors.New(http.StatusUnauthorized, "incorrect api key auth")
}
