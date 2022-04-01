package internal

import (
	"context"
	"fmt"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal/auth"
	"net/http"
	"strings"
)

const ErrWithStatusCode = "error with status code %d"
const ErrNotAuthenticated = "You are not authenticated. Use the keptn auth command to authenticate based on the API token or an OAuth client"
const ErrForbidden = "You do not have enough permissions"
const ErrInternalServerError = "Keptn API seems to be down"

var PublicDiscovery = auth.NewOauthDiscovery(&http.Client{})

// APIProvider is used to get a handle to the Keptn API clients
var APIProvider = getAPISet

// getAPISet will create and return an API Set that already
// contains the correct HTTP Client to be used for contacting the API endpoints
// Depending on whether OAuth is in use or not it will create an OAUth enabled client or not.
// Further, authToken and httpClient is optional. If the user does not provide an auth token,
// it will not be part of the requests made via the entities of the API set. If an HTTP client is given
// it is used without further modifications (apart from the modifications being done in go-utils)
// It is also possible to provide nil as httpClient parameter, in which case a fresh HTTP client will be created
// but does NOT support OAuth
func getAPISet(baseURL string, keptnXToken string, httpClient ...*http.Client) (*apiutils.APISet, error) {
	return getAPISetWithOauthGetter(baseURL, keptnXToken, auth.NewOauthAuthenticator(PublicDiscovery, auth.NewLocalFileOauthStore(), auth.NewBrowser(), &auth.ClosingRedirectHandler{}), httpClient...)
}

func getAPISetWithOauthGetter(baseURL string, keptnXToken string, oauthAuthenticator auth.OAuthenticator, httpClient ...*http.Client) (*apiutils.APISet, error) {
	// if an HTTP client was explicitly given,
	// we just create and return a new APISet with that client
	if len(httpClient) > 0 {
		if keptnXToken == "" {
			return apiutils.New(baseURL, apiutils.WithHTTPClient(httpClient[0]))
		}
		return apiutils.New(baseURL, apiutils.WithAuthToken(keptnXToken), apiutils.WithHTTPClient(httpClient[0]))
	}
	// else, depending on whether OAuth is in use or not,
	// create and return a APISet with an OAuth enabled HTTP client or not
	var client *http.Client

	// check whether OAuth is in use
	if storeCreated := oauthAuthenticator.TokenStore().Created(); storeCreated {
		oauthInfo, err := oauthAuthenticator.TokenStore().GetOauthInfo()
		if err != nil {
			return nil, err
		}
		// get the ready to use HTTP client
		client, err := oauthAuthenticator.OauthClient(context.Background())
		if err != nil {
			return nil, err
		}
		// check if the HTTP client is still usable, i.e.
		// make a call to the token endpoint
		// If it fails, we assume that it can be fixed by
		// starting the authorization code flow again
		// TODO: Check for a way to determine whether the error is related to an invalid refresh token
		_, err = client.Head(oauthInfo.DiscoveryInfo.TokenEndpoint)
		if err != nil {
			// start the authorization code flow
			err = oauthAuthenticator.Auth(*oauthInfo.ClientValues)
			if err != nil {
				return nil, err
			}
		}
		if keptnXToken == "" {
			return apiutils.New(baseURL, apiutils.WithHTTPClient(client))
		}
	}
	return apiutils.New(baseURL, apiutils.WithAuthToken(keptnXToken), apiutils.WithHTTPClient(client))
}

func OnAPIError(err error) error {
	switch 0 {
	case compareError(err, ErrWithStatusCode, 401):
		return fmt.Errorf(ErrNotAuthenticated)
	case compareError(err, ErrWithStatusCode, 403):
		return fmt.Errorf(ErrForbidden)
	case compareError(err, ErrWithStatusCode, 500):
		return fmt.Errorf(ErrInternalServerError)
	default:
		return err
	}
}

func compareError(err error, msg string, code int) int {
	return strings.Compare(err.Error(), fmt.Sprintf(msg, code))
}
