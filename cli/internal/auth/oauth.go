package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

const (
	loginSuccessHTML = `<p><strong>Login successful!</strong></p>
						<script type="text/javascript">	
								setTimeout(function(){
											close();
								},1500)
						</script>`
	redirectURL = "http://localhost:3000/oauth/redirect"
)

// OAuthenticator represents just the interface for a component performing OAuth authentication
type OAuthenticator interface {
	Auth(clientValues OauthClientValues) error
	OauthClient(ctx context.Context) (*http.Client, error)
	TokenStore() OauthStore
}

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
	if err := clientValues.ValidateMandatoryFields(); err != nil {
		return err
	}

	discoveryInfo, err := a.discovery.Discover(context.TODO(), clientValues.OauthDiscoveryURL)
	if err != nil {
		return fmt.Errorf("failed to perform OAuth Discovery using URL %s: %w: ", clientValues.OauthDiscoveryURL, err)
	}

	config := &oauth2.Config{
		ClientID:     clientValues.OauthClientID,
		ClientSecret: clientValues.OauthClientSecret,
		Scopes:       clientValues.OauthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  discoveryInfo.AuthorizationEndpoint,
			TokenURL: discoveryInfo.TokenEndpoint,
		},
		RedirectURL: redirectURL,
	}

	enforceOpenIDScope(config)

	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return fmt.Errorf("failed to generate code verifier: %w", err)
	}
	sum := sha256.Sum256(codeVerifier)
	codeChallenge := strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:]), "=")

	state, err := State(10)
	if err != nil {
		return fmt.Errorf("failed to generate random state query parameter")
	}
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	if err := a.browser.Open(authURL); err != nil {
		return fmt.Errorf("failed to open user Browser: %w", err)
	}

	token, err := a.redirectHandler.Handle(codeVerifier, config, state)
	if err != nil {
		return fmt.Errorf("failed to handle redirect: %w", err)
	}

	oauthInfo := &OauthInfo{
		DiscoveryInfo: discoveryInfo,
		ClientValues:  &clientValues,
		Token:         token,
	}
	if err := a.tokenStore.StoreOauthInfo(oauthInfo); err != nil {
		return fmt.Errorf("failed to sotre oauth information: %w", err)
	}
	return nil
}

// GetOauthClient will eventually return an already ready to use http client which is configured to use
// a OAUth Access Token
func (a *OauthAuthenticator) OauthClient(ctx context.Context) (*http.Client, error) {

	oauthInfo, err := a.tokenStore.GetOauthInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth HTTP client: %w", err)
	}

	config := &oauth2.Config{
		ClientSecret: oauthInfo.ClientValues.OauthClientSecret,
		ClientID:     oauthInfo.ClientValues.OauthClientID,
		Scopes:       oauthInfo.ClientValues.OauthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauthInfo.DiscoveryInfo.AuthorizationEndpoint,
			TokenURL: oauthInfo.DiscoveryInfo.TokenEndpoint,
		},
		RedirectURL: redirectURL,
	}

	enforceOpenIDScope(config)

	nrts := &NotifyRefreshTokenSource{
		config:     config,
		tokenStore: a.tokenStore,
	}
	return oauth2.NewClient(ctx, nrts), nil
}

func (a *OauthAuthenticator) TokenStore() OauthStore {
	return a.tokenStore
}

func enforceOpenIDScope(config *oauth2.Config) {
	openIDScopePresent := false
	for _, s := range config.Scopes {
		if s == "openid" {
			openIDScopePresent = true
			break
		}
	}
	if !openIDScopePresent {
		config.Scopes = append(config.Scopes, "openid")
	}
}

// OauthClientValues are values set by the user when performing OAuth flow
type OauthClientValues struct {
	OauthDiscoveryURL string   `json:"oauth_discovery_url"`
	OauthClientID     string   `json:"oauth_client_id"`
	OauthClientSecret string   `json:"oauth_client_secret"`
	OauthScopes       []string `json:"oauth_scopes"`
}

func (v *OauthClientValues) ValidateMandatoryFields() error {
	if v.OauthClientID == "" || v.OauthDiscoveryURL == "" {
		return fmt.Errorf("client values invalid: client id and discovery URL must be set")
	}
	return nil
}
