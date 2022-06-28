package api

import (
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/config"
	"net/http"
	"net/url"
	"strings"
)

// Initializer implements both methods of creating a new keptn API with internal or remote execution plane
type Initializer struct {
	Remote   func(baseURL string, options ...func(*keptnapi.APISet)) (*keptnapi.APISet, error)
	Internal func(client *http.Client, apiMappings ...api.InClusterAPIMappings) (*api.InternalAPISet, error)
}

func CreateKeptnAPI(httpClient *http.Client, env config.EnvConfig) (keptnapi.KeptnInterface, error) {
	return createAPI(httpClient, env, Initializer{keptnapi.New, api.NewInternal})
}

func createAPI(httpClient *http.Client, env config.EnvConfig, apiInit Initializer) (keptnapi.KeptnInterface, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		scheme := "http"
		parsed, err := url.ParseRequestURI(env.KeptnAPIEndpoint)
		if err != nil {
			return nil, err
		}

		// accepts either "" or http
		if env.KeptnAPIEndpoint != "" && (parsed.Scheme == "" || !strings.HasPrefix(parsed.Scheme, "http")) {
			return nil, fmt.Errorf("invalid scheme for keptn endpoint, %s is not http or https", env.KeptnAPIEndpoint)
		}

		if strings.HasPrefix(parsed.Scheme, "http") {
			// if no value is assigned to the endpoint than we keep the default scheme
			scheme = parsed.Scheme
		}
		return apiInit.Remote(env.KeptnAPIEndpoint, keptnapi.WithScheme(scheme), keptnapi.WithHTTPClient(httpClient), keptnapi.WithAuthToken(env.KeptnAPIToken))
	}

	return apiInit.Internal(httpClient)
}
