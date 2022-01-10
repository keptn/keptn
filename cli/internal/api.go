package internal

import (
	"context"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	auth2 "github.com/keptn/keptn/cli/internal/auth"
	"net/http"
)

func GetApiSet(baseURL string, authToken string, authHeader string, scheme string) (*apiutils.ApiSet, error) {
	var client *http.Client
	tokenStore := auth2.NewLocalFileTokenStore()
	if tokenInitialized, _ := tokenStore.Location(); tokenInitialized {
		oauthConfig, err := auth2.GetOauthConfig(auth2.StaticOauthDiscovery{})
		if err != nil {
			return nil, err
		}
		browser := auth2.NewBrowser()
		oauth := auth2.NewOauthAuthenticator(oauthConfig, tokenStore, browser)

		client, err = oauth.GetOauthClient(context.Background())
		if err != nil {
			return nil, err
		}
	}
	apiset, err := apiutils.NewApiSet(baseURL, authToken, authHeader, client, scheme)
	if err != nil {
		return nil, err
	}
	return apiset, nil
}
