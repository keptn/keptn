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

type TokenGetter interface {
	Handle() (*oauth2.Token, error)
}

type ClosingRedirectHandler struct {
	server       *http.Server
	codeVerifier []byte
	oauthConfig  *oauth2.Config
}

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
