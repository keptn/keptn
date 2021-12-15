package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const loginSuccessHTML = `<p><strong>Login successful!</strong></p>`

type Authenticator interface {
	Authorize(discovery OauthLocationGetter, tokenStore TokenStore, redirectURL string) error
}

type OauthAuthenticator struct {
	discovery   OauthLocationGetter
	tokenStore  TokenStore
	browser     Browser
	redirectURL string
}

func NewOauthAuthenticator(discovery OauthLocationGetter, tokenStore TokenStore, browser Browser) *OauthAuthenticator {
	return &OauthAuthenticator{
		discovery:   discovery,
		tokenStore:  tokenStore,
		browser:     browser,
		redirectURL: "http://localhost:3000/oauth/redirect",
	}
}

func (a *OauthAuthenticator) Authorize() error {
	oauthConfig, err := a.discover()
	if err != nil {
		return err
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
	server := &http.Server{Addr: a.redirectURL}
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		defer cleanup(server)
		queryParts, _ := url.ParseQuery(r.URL.RawQuery)
		code := queryParts["code"][0]

		tok, err := oauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", string(codeVerifier)))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = a.tokenStore.StoreToken(tok)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Fprint(w, loginSuccessHTML)
	})

	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Printf("can't listen to port %s: %s\n", ":3000", err)
		os.Exit(1)
	}
	server.Serve(l)
	return nil
}

func cleanup(server io.Closer) {
	go server.Close()
}

func (a *OauthAuthenticator) GetOauthClient(ctx context.Context, tokenStore TokenStore) (*http.Client, error) {
	oauthConfig, err := a.discover()
	if err != nil {
		return nil, err
	}
	nrts := &NotifyRefreshTokenSource{
		config:     oauthConfig,
		tokenStore: tokenStore,
	}
	return oauth2.NewClient(ctx, nrts), nil
}

func (a *OauthAuthenticator) discover() (*oauth2.Config, error) {
	r, err := a.discovery.Discover()
	if err != nil {
		return nil, err
	}
	return &oauth2.Config{
		ClientID:     "dt0s03.cloudautomation-keptn-local",
		ClientSecret: "",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  r.AuthorizationEndpoint,
			TokenURL: r.TokenEndpoint,
		},
		RedirectURL: a.redirectURL,
	}, nil
}
