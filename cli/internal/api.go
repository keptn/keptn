package internal

import (
	"context"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	auth2 "github.com/keptn/keptn/cli/internal/auth"
	"net/http"
)

var PublicDiscovery = auth2.NewOauthDiscovery(&http.Client{})

// APIProvider is used to get a handle to the Keptn API clients
var APIProvider = getAPISet

// getAPISet retrieves the ApiSet containing all Keptn client APIs
func getAPISet(baseURL string, authToken string, authHeader string, scheme string) (*apiutils.APISet, error) {
	var client *http.Client
	var err error
	tokenStore := auth2.NewLocalFileTokenStore()
	if tokenInitialized, _ := tokenStore.Location(); tokenInitialized {
		oauth := auth2.NewOauthAuthenticator(PublicDiscovery, tokenStore, auth2.NewBrowser())
		client, err = oauth.GetOauthClient(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return apiutils.NewAPISet(baseURL, authToken, authHeader, client, scheme)
}
