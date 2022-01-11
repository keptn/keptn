package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"net/url"
	"os"
)

// TokenGetter handles the retrieval of oauth access tokens
type TokenGetter interface {
	Handle() (*oauth2.Token, error)
}

// CLosingRedirectHandler is an implementation of TokenHandler
// It opens a local http server with the hard-coded path "/oauth/redirect"
// which serves as a callback to transfer the access tokens retrieved during the Oauth flow
type ClosingRedirectHandler struct {
	server       *http.Server
	codeVerifier []byte
	oauthConfig  *oauth2.Config
}

// Handle openc a server at port 3000 and performs the exchange of the authorization code into an access token when the
// user was redirect to the local server
// It returns the obtained oauth2 token or an error
func (r *ClosingRedirectHandler) Handle() (*oauth2.Token, error) {
	server := &http.Server{}
	var tokenExchangeErr error
	var acquiredToken *oauth2.Token

	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, req *http.Request) {
		defer func() { go server.Close() }()
		queryParts, _ := url.ParseQuery(req.URL.RawQuery)
		code := queryParts["code"][0]

		tok, err := r.oauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", string(r.codeVerifier)))
		if err != nil {
			fmt.Println(err.Error())
			tokenExchangeErr = err
			return
		}
		acquiredToken = tok
		fmt.Fprint(w, loginSuccessHTML)
	})
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Printf("can't listen to port %s: %s\n", ":3000", err)
		os.Exit(1)
	}
	server.Serve(l)
	return acquiredToken, tokenExchangeErr
}
