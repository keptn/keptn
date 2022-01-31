package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func Test_Auth_DependenciesFail(t *testing.T) {

	type fields struct {
		discovery       OauthLocationGetter
		tokenStore      OauthStore
		browser         URLOpener
		redirectHandler TokenGetter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{"Auth() - discovery fails",
			fields{
				discovery: &OauthDiscoveryMock{
					discoverFn: func(context.Context, string) (*OauthDiscoveryResult, error) {
						return nil, fmt.Errorf("disocvery failed ")
					},
				},
			}, assert.Error,
		},
		{"Auth() - open browser fails",
			fields{
				discovery: &OauthDiscoveryMock{
					discoverFn: func(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				browser: &BrowserMock{
					openFn: func(string) error { return errors.New("browser open failed") },
				},
				tokenStore: &TokenStoreMock{
					getTokenDiscoveryFn: func() (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
			}, assert.Error,
		},
		{"Auth() - callback handler fails",
			fields{
				discovery: &OauthDiscoveryMock{
					discoverFn: func(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				browser: &BrowserMock{
					openFn: func(string) error { return nil },
				},
				tokenStore: &TokenStoreMock{
					getTokenDiscoveryFn: func() (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				redirectHandler: RedirectHandlerMock{
					handleFn: func(bytes []byte, config *oauth2.Config, state string) (*oauth2.Token, error) {
						return nil, errors.New("callback handler failed")
					},
				},
			}, assert.Error,
		},
		{"Auth() - success",
			fields{
				discovery: &OauthDiscoveryMock{
					discoverFn: func(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				browser: &BrowserMock{
					openFn: func(string) error { return nil },
				},
				tokenStore: &TokenStoreMock{
					getTokenDiscoveryFn: func() (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				redirectHandler: RedirectHandlerMock{
					handleFn: func(bytes []byte, config *oauth2.Config, state string) (*oauth2.Token, error) {
						return &oauth2.Token{}, nil
					},
				},
			}, assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &OauthAuthenticator{
				discovery:       tt.fields.discovery,
				tokenStore:      tt.fields.tokenStore,
				browser:         tt.fields.browser,
				redirectHandler: tt.fields.redirectHandler,
			}
			tt.wantErr(t, a.Auth(OauthClientValues{"http://well-known-discovery-url.com", "clientID", "", []string{}}), "Auth()")
		})
	}
}

func Test_Auth_Scopes(t *testing.T) {
	discovery := &OauthDiscoveryMock{}
	tokenStore := &TokenStoreMock{}
	browser := &BrowserMock{}
	redirectHandler := &RedirectHandlerMock{}
	authenticator := &OauthAuthenticator{
		discovery:       discovery,
		tokenStore:      tokenStore,
		browser:         browser,
		redirectHandler: redirectHandler,
	}
	t.Run("Auth - default openid scope is set", func(t *testing.T) {
		redirectHandler.handleFn = func(b []byte, c *oauth2.Config, s string) (*oauth2.Token, error) {
			assert.Equal(t, 1, len(c.Scopes))
			assert.Equal(t, "openid", c.Scopes[0])
			return nil, nil
		}
		authenticator.Auth(OauthClientValues{"http://well-known-discovery-url.com", "clientID", "", []string{}})
	})
	t.Run("Auth - scopes always contain default openid scope", func(t *testing.T) {
		redirectHandler.handleFn = func(b []byte, c *oauth2.Config, s string) (*oauth2.Token, error) {
			assert.Equal(t, 2, len(c.Scopes))
			assert.Contains(t, c.Scopes, "openid")
			assert.Contains(t, c.Scopes, "somescope")
			return nil, nil
		}
		authenticator.Auth(OauthClientValues{"http://well-known-discovery-url.com", "clientID", "", []string{"somescope"}})
	})
}

func Test_Auth_MissingOauthInfo(t *testing.T) {
	discovery := &OauthDiscoveryMock{}
	tokenStore := &TokenStoreMock{}
	browser := &BrowserMock{}
	redirectHandler := &RedirectHandlerMock{}
	authenticator := &OauthAuthenticator{
		discovery:       discovery,
		tokenStore:      tokenStore,
		browser:         browser,
		redirectHandler: redirectHandler,
	}
	t.Run("Auth - client id set and client secret is optional", func(t *testing.T) {
		{
			redirectHandler.handleFn = func(b []byte, c *oauth2.Config, s string) (*oauth2.Token, error) {
				assert.Equal(t, "clientID", c.ClientID)
				assert.Equal(t, "", c.ClientSecret)
				return nil, nil
			}
			authenticator.Auth(OauthClientValues{"http://well-known-discovery-url.com", "clientID", "", []string{}})
		}
	})
	t.Run("Auth - client id and secret given", func(t *testing.T) {
		{
			redirectHandler.handleFn = func(b []byte, c *oauth2.Config, s string) (*oauth2.Token, error) {
				assert.Equal(t, "clientID", c.ClientID)
				assert.Equal(t, "clientSecret", c.ClientSecret)
				return nil, nil
			}
			authenticator.Auth(OauthClientValues{"http://well-known-discovery-url.com", "clientID", "clientSecret", []string{}})
		}
	})
	t.Run("Auth - client id missing", func(t *testing.T) {
		{
			err := authenticator.Auth(OauthClientValues{"http://well-known-discovery-url.com", "", "", []string{}})
			assert.NotNil(t, err)
		}
	})
}

func Test_GetOauthClient(t *testing.T) {
	type fields struct {
		discovery  OauthLocationGetter
		tokenStore OauthStore
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
		{"GetOauthClient - no persisted oauth info",
			fields{
				discovery: &OauthDiscoveryMock{
					discoverFn: func(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				tokenStore: &TokenStoreMock{
					getOauthInfoFn: func() (*OauthInfo, error) {
						return nil, fmt.Errorf("not found")
					},
				},
			},
			args{
				context.TODO(),
			},
			assert.Nil,
			assert.Error},
		{"GetOauthClient - success",
			fields{
				discovery: &OauthDiscoveryMock{
					discoverFn: func(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
						return &OauthDiscoveryResult{}, nil
					},
				},
				tokenStore: &TokenStoreMock{
					getOauthInfoFn: func() (*OauthInfo, error) {
						return &OauthInfo{
							DiscoveryInfo: &OauthDiscoveryResult{},
							ClientValues:  &OauthClientValues{},
							Token:         &oauth2.Token{},
						}, nil
					},
				},
			},
			args{
				context.TODO(),
			},
			assert.NotNil,
			assert.NoError},
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
