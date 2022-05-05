package api

import (
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"net/http"
)

// Initializer implements both methods of creating a new keptn API with internal or remote execution plane
type Initializer struct {
	Remote   func(baseURL string, options ...func(*keptnapi.APISet)) (*keptnapi.APISet, error)
	Internal func(client *http.Client, apiMappings ...InClusterAPIMappings) (*InternalAPISet, error)
}
