package api

import (
	"github.com/keptn/go-utils/pkg/api/utils"
	"net/http"
)

// Initializer implements both methods of creating a new keptn API with internal or remote control plane
type Initializer struct {
	Remote   func(baseURL string, options ...func(*api.APISet)) (*api.APISet, error)
	Internal func(client *http.Client, apiMappings ...InClusterAPIMappings) (*InternalAPISet, error)
}
