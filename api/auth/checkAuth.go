// This file is safe to edit. Once it exists it will not be overwritten

package auth

import (
	"os"

	errors "github.com/go-openapi/errors"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations"
)

// CheckToken checkes whether the token is correct
func CheckToken(api *operations.API, token string) (*models.Principal, error) {

	if token == os.Getenv("keptn-api-token") {
		prin := models.Principal(token)
		return &prin, nil
	}
	api.Logger("Access attempt with incorrect api key auth: %s", token)
	return nil, errors.New(401, "incorrect api key auth")
}
