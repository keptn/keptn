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
func getAPISet(baseURL string, authToken string) (*apiutils.APISet, error) {
	var client *http.Client
	var err error
	tokenStore := auth2.NewLocalFileOauthStore()
	if storeCreated := tokenStore.Created(); storeCreated {
		oauth := auth2.NewOauthAuthenticator(PublicDiscovery, tokenStore, auth2.NewBrowser(), &auth2.ClosingRedirectHandler{})
		client, err = oauth.GetOauthClient(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return apiutils.New(baseURL, apiutils.WithAuthToken(authToken), apiutils.WithHTTPClient(client))
}
