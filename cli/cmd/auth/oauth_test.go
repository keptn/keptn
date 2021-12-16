package auth

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOauthAuthenticator_Auth_StoresTokenInTokenStore(t *testing.T) {
	server, serverCLoser := setupMockOAuthServer()
	defer serverCLoser()

	discovery := &OauthDiscoveryMock{
		discoverFn: func() (*OauthDiscoveryResult, error) {
			return &OauthDiscoveryResult{
				AuthorizationEndpoint: server.URL + "/auth",
				TokenEndpoint:         server.URL + "/token",
			}, nil
		},
	}
	tokenStore := &TokenStoreMock{}
	browser := &BrowserMock{
		openFn: func(string) error { return nil },
	}
	authenticator := NewOauthAuthenticator(discovery, tokenStore, browser)
	go func() {
		err := authenticator.Auth()
		assert.Nil(t, err)
	}()
	assert.Eventuallyf(t, func() bool {
		_, err := http.Get("http://localhost:3000/oauth/redirect?code=code") //nolint:bodyclose
		return err == nil
	}, 5*time.Second, 1*time.Second, "")

	assert.Eventuallyf(t, func() bool {
		return tokenStore.storedToken != nil && (tokenStore.storedToken.AccessToken == "mocked-token")
	}, 5*time.Second, 1*time.Second, "")
}

func TestOauthAuthenticator_Auth1(t *testing.T) {
	type fields struct {
		discovery  OauthLocationGetter
		tokenStore TokenStore
		browser    URLOpener
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{"discovery fails", fields{
			discovery: &OauthDiscoveryMock{
				discoverFn: func() (*OauthDiscoveryResult, error) { return nil, errors.New("NOPE") }},
		}, assert.Error,
		},
		{"open browser fails", fields{
			discovery: &OauthDiscoveryMock{},
			browser: &BrowserMock{
				openFn: func(string) error { return errors.New("NOPE") },
			},
		}, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &OauthAuthenticator{
				discovery:  tt.fields.discovery,
				tokenStore: tt.fields.tokenStore,
				browser:    tt.fields.browser,
			}
			tt.wantErr(t, a.Auth(), fmt.Sprintf("Auth()"))
		})
	}
}

func setupMockOAuthServer() (*httptest.Server, func()) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=mocked-token&scope=user&token_type=bearer"))
	})

	server := httptest.NewServer(mux)

	return server, func() {
		server.Close()
	}
}
