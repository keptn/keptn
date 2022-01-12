package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOauthAuthenticator_Auth_StoresTokenAndDiscoveryInfoLocally(t *testing.T) {
	server, serverCLoser := setupMockOAuthServer()
	defer serverCLoser()

	discovery := &OauthDiscoveryMock{
		discoverFn: func(ctx context.Context) (*OauthDiscoveryResult, error) {
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
		token, _ := tokenStore.GetToken()
		tokenDiscovery, _ := tokenStore.GetTokenDiscovery()
		return token != nil && token.AccessToken == "mocked-token" && tokenDiscovery != nil
	}, 5*time.Second, 1*time.Second, "")
}

func TestOauthAuthenticator_Auth(t *testing.T) {
	discovery := &OauthDiscoveryMock{
		discoverFn: func(ctx context.Context) (*OauthDiscoveryResult, error) {
			return &OauthDiscoveryResult{}, nil
		},
	}

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
		{"open browser fails", fields{
			discovery: discovery,
			browser: &BrowserMock{
				openFn: func(string) error { return errors.New("NOPE") },
			},
			tokenStore: &TokenStoreMock{
				getTokenDiscoveryFn: func() (*OauthDiscoveryResult, error) {
					return &OauthDiscoveryResult{}, nil
				},
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
			tt.wantErr(t, a.Auth(), "Auth()")
		})
	}
}

func TestOauthAuthenticator_GetOauthClient(t *testing.T) {
	discovery := &OauthDiscoveryMock{
		discoverFn: func(ctx context.Context) (*OauthDiscoveryResult, error) {
			return &OauthDiscoveryResult{}, nil
		},
	}

	type fields struct {
		discovery  OauthLocationGetter
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
		{"", fields{discovery: discovery}, args{context.TODO()}, assert.NotNil, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &OauthAuthenticator{
				discovery:  tt.fields.discovery,
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
