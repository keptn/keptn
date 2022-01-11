package internal

import (
	"context"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	auth2 "github.com/keptn/keptn/cli/internal/auth"
	"net/http"
	"time"
)

var PublicDiscovery = auth2.NewOauthDiscovery(&http.Client{}, "https://sso-dev.dynatracelabs.com/.well-known/openid-configuration", 10*time.Second)

// APIProvider is used to get a handle to the kepn API clients
var APIProvider = getAPISet

// getAPISet retrieves the ApiSet containing all keptn client APIs
func getAPISet(baseURL string, authToken string, authHeader string, scheme string) (*apiutils.ApiSet, error) {
	var client *http.Client
	tokenStore := auth2.NewLocalFileTokenStore()
	if tokenInitialized, _ := tokenStore.Location(); tokenInitialized {
		oauthConfig, err := auth2.GetOauthConfig(PublicDiscovery)
		if err != nil {
			return nil, err
		}
		oauth := auth2.NewOauthAuthenticator(oauthConfig, tokenStore, auth2.NewBrowser())
		client, err = oauth.GetOauthClient(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return apiutils.NewApiSet(baseURL, authToken, authHeader, client, scheme)
}
