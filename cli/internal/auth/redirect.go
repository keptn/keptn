package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"net/url"
)

// TokenGetter handles the retrieval of oauth access tokens
type TokenGetter interface {
	Handle(codeVerifier []byte, oauthConfig *oauth2.Config, state string) (*oauth2.Token, error)
}

// ClosingRedirectHandler is an implementation of TokenGetter
// It opens a local http server with the hard-coded path "/oauth/redirect"
// which serves as a callback to transfer the access tokens retrieved during the Oauth flow
type ClosingRedirectHandler struct{}

// Handle opens a server at port 3000 and performs the exchange of the authorization code into an access token when the
// user was redirect to the local server
// It returns the obtained oauth2 token or an error
// TODO: close handler after a timeout
// TODO: get rid of hard-coded path and port
func (r *ClosingRedirectHandler) Handle(codeVerifier []byte, oauthConfig *oauth2.Config, oauthState string) (*oauth2.Token, error) {
	server := &http.Server{}
	var tokenExchangeErr error
	var acquiredToken *oauth2.Token

	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, req *http.Request) {
		defer func() { go server.Close() }()
		queryParts, _ := url.ParseQuery(req.URL.RawQuery)
		state := queryParts["state"][0]
		if state != oauthState {
			tokenExchangeErr = fmt.Errorf("invalid oauth state")
			return
		}
		code := queryParts["code"][0]
		tok, err := oauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", string(codeVerifier)))
		if err != nil {
			tokenExchangeErr = err
			return
		}
		acquiredToken = tok
		fmt.Fprint(w, loginSuccessHTML)
	})
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		return nil, err
	}
	server.Serve(l)
	return acquiredToken, tokenExchangeErr
}
