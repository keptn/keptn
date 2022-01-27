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
	var handleRedirectErr error
	var acquiredToken *oauth2.Token

	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, req *http.Request) {
		defer func() { go server.Close() }()
		queryParts, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			handleRedirectErr = err
			return
		}
		state := queryParts["state"]
		if len(state) == 0 {
			handleRedirectErr = fmt.Errorf("no oauth state param found")
			return
		}
		if state[0] != oauthState {
			handleRedirectErr = fmt.Errorf("invalid oauth state")
			return
		}
		code := queryParts["code"]
		if len(code) == 0 {
			handleRedirectErr = fmt.Errorf("no code param fround")
			return
		}
		tok, err := oauthConfig.Exchange(context.Background(), code[0], oauth2.SetAuthURLParam("code_verifier", string(codeVerifier)))
		if err != nil {
			handleRedirectErr = err
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
	return acquiredToken, handleRedirectErr
}
