package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

const (
	loginSuccessHTML = `<p><strong>Login successful!</strong></p>`
	redirectURL      = "http://localhost:3000/oauth/redirect"
	openIDScope      = "openid"
)

// OauthAuthenticator is an implementation of Authenticator which implements the Oauth2 Authorization Code Flow
type OauthAuthenticator struct {
	discovery       OauthLocationGetter
	tokenStore      OauthStore
	browser         URLOpener
	redirectHandler TokenGetter
}

// NewOauthAuthenticator is creating a new OauthAuthenticator
func NewOauthAuthenticator(discovery OauthLocationGetter, tokenStore OauthStore, browser URLOpener, redirectHandler TokenGetter) *OauthAuthenticator {
	return &OauthAuthenticator{
		discovery:       discovery,
		tokenStore:      tokenStore,
		browser:         browser,
		redirectHandler: redirectHandler,
	}
}

// Auth tries to start the Oauth2 Authorization Code Flow
func (a *OauthAuthenticator) Auth(clientValues OauthClientValues) error {
	discoveryInfo, err := a.discovery.Discover(context.TODO(), clientValues.OauthDiscoveryURL)
	if err != nil {
		return err
	}

	config := &oauth2.Config{
		ClientID:     clientValues.OauthClientID,
		ClientSecret: clientValues.OauthClientSecret,
		Scopes:       []string{openIDScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  discoveryInfo.AuthorizationEndpoint,
			TokenURL: discoveryInfo.TokenEndpoint,
		},
		RedirectURL: redirectURL,
	}

	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return err
	}
	sum := sha256.Sum256(codeVerifier)
	codeChallenge := strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:]), "=")

	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	if err := a.browser.Open(authURL); err != nil {
		return err
	}

	token, err := a.redirectHandler.Handle(codeVerifier, config)
	if err != nil {
		return err
	}

	oauthInfo := &OauthInfo{
		DiscoveryInfo: discoveryInfo,
		ClientValues:  &clientValues,
		Token:         token,
	}
	return a.tokenStore.StoreOauthInfo(oauthInfo)
}

// GetOauthClient will eventually return an already ready to use http client which is configured to use
// a OAUth Access Token
func (a *OauthAuthenticator) GetOauthClient(ctx context.Context) (*http.Client, error) {
	oauthInfo, err := a.tokenStore.GetOauthInfo()
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientSecret: oauthInfo.ClientValues.OauthClientSecret,
		ClientID:     oauthInfo.ClientValues.OauthClientID,
		Scopes:       []string{openIDScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauthInfo.DiscoveryInfo.AuthorizationEndpoint,
			TokenURL: oauthInfo.DiscoveryInfo.TokenEndpoint,
		},
		RedirectURL: redirectURL,
	}
	nrts := &NotifyRefreshTokenSource{
		config:     config,
		tokenStore: a.tokenStore,
	}
	return oauth2.NewClient(ctx, nrts), nil
}

// OauthClientValues are values set by the user when performing SSO
type OauthClientValues struct {
	OauthDiscoveryURL string `json:"oauth_discovery_url"`
	OauthClientID     string `json:"oauth_client_id"`
	OauthClientSecret string `json:"oauth_client_secret"`
}
