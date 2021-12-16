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

type Authenticator interface {
	Authorize(discovery OauthLocationGetter, tokenStore TokenStore, redirectURL string) error
}

type OauthAuthenticator struct {
	discovery  OauthLocationGetter
	tokenStore TokenStore
	browser    URLOpener
}

func NewOauthAuthenticator(discovery OauthLocationGetter, tokenStore TokenStore, browser URLOpener) *OauthAuthenticator {
	return &OauthAuthenticator{
		discovery:  discovery,
		tokenStore: tokenStore,
		browser:    browser,
	}
}

func (a *OauthAuthenticator) Auth() error {
	r, err := a.discovery.Discover()
	if err != nil {
		return err
	}
	oauthConfig := &oauth2.Config{
		ClientID:     "dt0s03.cloudautomation-keptn-local",
		ClientSecret: "",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  r.AuthorizationEndpoint,
			TokenURL: r.TokenEndpoint,
		},
		RedirectURL: "http://localhost:3000/oauth/redirect",
	}

	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return err
	}
	sum := sha256.Sum256(codeVerifier)
	codeChallenge := strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:]), "=")

	authURL := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	err = a.browser.Open(authURL)
	if err != nil {
		return err
	}

	redirectHandler := ClosingRedirectHandler{
		codeVerifier: codeVerifier,
		oauthConfig:  oauthConfig,
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

func (a *OauthAuthenticator) GetOauthClient(ctx context.Context) (*http.Client, error) {
	r, err := a.discovery.Discover()
	if err != nil {
		return nil, err
	}
	oauthConfig := &oauth2.Config{
		ClientID: "dt0s03.cloudautomation-keptn-local",
		Scopes:   []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  r.AuthorizationEndpoint,
			TokenURL: r.TokenEndpoint,
		},
		RedirectURL: "http://localhost:3000/oauth/redirect",
	}

	nrts := &NotifyRefreshTokenSource{
		config:     oauthConfig,
		tokenStore: a.tokenStore,
	}
	return oauth2.NewClient(ctx, nrts), nil
}
