package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

const loginSuccessHTML = `<p><strong>Login successful!</strong></p>`

// Authenticator is responsible for authenticate the user using SSO/Oauth2
type Authenticator interface {
	// Auth is triggering the authentication
	Auth(discovery OauthLocationGetter, tokenStore TokenStore, redirectURL string) error
}

// OauthAuthenticator is an implementation of Authenticator which implements the Oauth2 Authorization Code Flow
type OauthAuthenticator struct {
	discovery  OauthLocationGetter
	tokenStore TokenStore
	browser    URLOpener
}

// NewOauthAuthenticator is creating a new OauthAuthenticator
func NewOauthAuthenticator(discovery OauthLocationGetter, tokenStore TokenStore, browser URLOpener) *OauthAuthenticator {
	return &OauthAuthenticator{
		discovery:  discovery,
		tokenStore: tokenStore,
		browser:    browser,
	}
}

// Auth tries to start the Oauth2 Authorization Code Flow
func (a *OauthAuthenticator) Auth() error {
	discoveryInfo, err := a.tokenStore.GetTokenDiscovery()
	if err != nil {
		return err
	}
	if discoveryInfo == nil {
		if discoveryInfo, err = a.discovery.Discover(context.TODO()); err != nil {
			return err
		}
		if err := a.tokenStore.StoreTokenDiscovery(discoveryInfo); err != nil {
			return err
		}
	}

	config := &oauth2.Config{
		ClientID: "dt0s03.cloudautomation-keptn-local",
		Scopes:   []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  discoveryInfo.AuthorizationEndpoint,
			TokenURL: discoveryInfo.TokenEndpoint,
		},
		RedirectURL: "http://localhost:3000/oauth/redirect",
	}

	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return err
	}
	sum := sha256.Sum256(codeVerifier)
	codeChallenge := strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:]), "=")

	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	err = a.browser.Open(authURL)
	if err != nil {
		return err
	}

	redirectHandler := ClosingRedirectHandler{
		codeVerifier: codeVerifier,
		oauthConfig:  config,
	}

	token, err := redirectHandler.Handle()
	if err != nil {
		return err
	}

	err = a.tokenStore.StoreToken(token)
	if err != nil {
		return err
	}

	return nil
}

// GetOauthClient will eventually return an already ready to use http client which is configuered to use
// the correct access token
func (a *OauthAuthenticator) GetOauthClient(ctx context.Context) (*http.Client, error) {
	result, err := a.discovery.Discover(ctx)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID: "dt0s03.cloudautomation-keptn-local",
		Scopes:   []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  result.AuthorizationEndpoint,
			TokenURL: result.TokenEndpoint,
		},
		RedirectURL: "http://localhost:3000/oauth/redirect",
	}
	nrts := &NotifyRefreshTokenSource{
		config:     config,
		tokenStore: a.tokenStore,
	}
	return oauth2.NewClient(ctx, nrts), nil
}

//// GetOauthConfig uses a given OauthLocationGetter to determine the parameters
//// for the Oauth flow
//func GetOauthConfig(discovery OauthLocationGetter) (*oauth2.Config, *OauthDiscoveryResult, error) {
//	result, err := discovery.Discover(context.TODO())
//	if err != nil {
//		return nil, nil, err
//	}
//	oauthConfig := &oauth2.Config{
//		ClientID: "dt0s03.cloudautomation-keptn-local",
//		Scopes:   []string{"openid"},
//		Endpoint: oauth2.Endpoint{
//			AuthURL:  result.AuthorizationEndpoint,
//			TokenURL: result.TokenEndpoint,
//		},
//		RedirectURL: "http://localhost:3000/oauth/redirect",
//	}
//
//	return oauthConfig, result, nil
//}
