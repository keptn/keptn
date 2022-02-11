package client

import (
	"fmt"
	auth "github.com/keptn/go-utils/pkg/common/oauth"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSimpleClientGetter_Get(t *testing.T) {
	cfg := config.EnvConfig{}
	getter := New(cfg)
	c, err := getter.Get()
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestOAuthClientGetter_Get(t *testing.T) {
	t.Run("Get - No Discovery, nor Token URL given", func(t *testing.T) {
		cfg := config.EnvConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
		}
		oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
	t.Run("Get - With Discovery URL", func(t *testing.T) {
		cfg := config.EnvConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
			OAuthDiscovery:    "http://some-url.com",
		}
		oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.NotNil(t, c)
		assert.Nil(t, err)
	})
	t.Run("Get - With Token URL", func(t *testing.T) {
		cfg := config.EnvConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
			OauthTokenURL:     "http://some-url.com",
		}
		oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.NotNil(t, c)
		assert.Nil(t, err)
	})
	t.Run("Get - missing scopes", func(t *testing.T) {
		cfg := config.EnvConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthDiscovery:    "http://some-url.com",
		}
		oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
	t.Run("Get - missing client id", func(t *testing.T) {
		cfg := config.EnvConfig{
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
			OAuthDiscovery:    "http://some-url.com",
		}
		oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
	t.Run("Get - missing client secret", func(t *testing.T) {
		cfg := config.EnvConfig{
			OAuthClientID:  "client-id",
			OAuthScopes:    []string{"scope"},
			OAuthDiscovery: "http://some-url.com",
		}
		oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
}

func TestOAuthClientGetter_Get_TokenEndpointIsCalled(t *testing.T) {
	tokenURLCalled := false
	tokenURLSrv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		tokenURLCalled = true
		rw.WriteHeader(http.StatusOK)
	}))
	defer tokenURLSrv.Close()

	cfg := config.EnvConfig{
		OAuthClientID:     "client-id",
		OAuthClientSecret: "client-secret",
		OAuthScopes:       []string{"scope"},
		OAuthDiscovery:    "http://some-wellknown-url.com",
	}
	oauthDiscovery := &auth.StaticOauthDiscovery{DiscoveryValues: &auth.OauthDiscoveryResult{
		TokenEndpoint: tokenURLSrv.URL,
	}}

	c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
	assert.NotNil(t, c)
	assert.Nil(t, err)

	// next line will obviously fail,
	// but we only want to check whether the token endpoint
	// is called
	c.Get("localhost")
	assert.Eventually(t, func() bool {
		fmt.Println(tokenURLCalled)
		return tokenURLCalled == true
	}, time.Second, 100*time.Millisecond)

}
