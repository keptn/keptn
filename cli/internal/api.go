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

// getAPISet will create and return an API Set that already
// contains the correct HTTP Client to be used for contacting the API endpoints
// Depending on wheather SSO is in use or not it will create an OAUth enabled client or not.
// Further, authToken and httpClient is optional. If ther user does not provide an auth token,
// it will not be part of the requests made via the entities of the API set. If a HTTP client is given
// it is used without further modifications (apart from the modifications being done in go-utils)
// It is also possible to provide nil as httpClient parameter, in which case a fresh HTTP client will be created
// but does NOT support OAuth
func getAPISet(baseURL string, authToken string, httpClient ...*http.Client) (*apiutils.APISet, error) {
	// if a HTTP client was explicitly given,
	// we just create and return a new APISet with that client
	if len(httpClient) > 0 {
		if authToken == "" {
			return apiutils.New(baseURL, apiutils.WithHTTPClient(httpClient[0]))
		}
		return apiutils.New(baseURL, apiutils.WithAuthToken(authToken), apiutils.WithHTTPClient(httpClient[0]))
	}
	// else, depending on whether SSO is in use or not,
	// create and return a APISet with an Oauth enabled HTTP client or not
	var client *http.Client
	var err error
	tokenStore := auth2.NewLocalFileOauthStore()
	if storeCreated := tokenStore.Created(); storeCreated {
		oauth := auth2.NewOauthAuthenticator(PublicDiscovery, tokenStore, auth2.NewBrowser(), &auth2.ClosingRedirectHandler{})
		client, err = oauth.GetOauthClient(context.Background())
		if err != nil {
			return nil, err
		}
		if authToken == "" {
			return apiutils.New(baseURL, apiutils.WithHTTPClient(client))
		}
	}
	return apiutils.New(baseURL, apiutils.WithAuthToken(authToken), apiutils.WithHTTPClient(client))
}
