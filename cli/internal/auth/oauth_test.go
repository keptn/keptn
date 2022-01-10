package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
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

	config, _ := GetOauthConfig(discovery)
	tokenStore := &TokenStoreMock{}
	browser := &BrowserMock{
		openFn: func(string) error { return nil },
	}
	authenticator := NewOauthAuthenticator(config, tokenStore, browser)
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
	discovery := &OauthDiscoveryMock{
		discoverFn: func() (*OauthDiscoveryResult, error) {
			return &OauthDiscoveryResult{}, nil
		},
	}

	config, _ := GetOauthConfig(discovery)

	type fields struct {
		config     *oauth2.Config
		tokenStore TokenStore
		browser    URLOpener
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{"open browser fails", fields{
			config: config,
			browser: &BrowserMock{
				openFn: func(string) error { return errors.New("NOPE") },
			},
		}, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &OauthAuthenticator{
				config:     tt.fields.config,
				tokenStore: tt.fields.tokenStore,
				browser:    tt.fields.browser,
			}
			tt.wantErr(t, a.Auth(), fmt.Sprintf("Auth()"))
		})
	}
}

func TestOauthAuthenticator_GetOauthClient(t *testing.T) {
	type fields struct {
		config     *oauth2.Config
		tokenStore TokenStore
		browser    URLOpener
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantClient assert.ValueAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		{"", fields{config: &oauth2.Config{}}, args{context.TODO()}, assert.NotNil, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &OauthAuthenticator{
				config:     tt.fields.config,
				tokenStore: tt.fields.tokenStore,
				browser:    tt.fields.browser,
			}
			got, err := a.GetOauthClient(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("GetOauthClient(%v)", tt.args.ctx)) {
				return
			}
			if !tt.wantClient(t, got, fmt.Sprintf("GetOauthClient(%v)", tt.args.ctx)) {
				return
			}
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
