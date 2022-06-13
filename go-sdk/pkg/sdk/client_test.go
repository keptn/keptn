package sdk

import (
	"fmt"
	oauthutils "github.com/keptn/go-utils/pkg/common/oauth2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSimpleClientGetter_Get(t *testing.T) {
	cfg := envConfig{}
	getter := New(cfg)
	c, err := getter.Get()
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestOAuthClientGetter_Get(t *testing.T) {
	t.Run("Get - No Discovery, nor Token URL given", func(t *testing.T) {
		cfg := envConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
		}
		oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
	t.Run("Get - With Discovery URL", func(t *testing.T) {
		cfg := envConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
			OAuthDiscovery:    "http://some-url.com",
		}
		oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.NotNil(t, c)
		assert.Nil(t, err)
	})
	t.Run("Get - With Token URL", func(t *testing.T) {
		cfg := envConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
			OauthTokenURL:     "http://some-url.com",
		}
		oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.NotNil(t, c)
		assert.Nil(t, err)
	})
	t.Run("Get - missing scopes", func(t *testing.T) {
		cfg := envConfig{
			OAuthClientID:     "client-id",
			OAuthClientSecret: "client-secret",
			OAuthDiscovery:    "http://some-url.com",
		}
		oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
	t.Run("Get - missing client id", func(t *testing.T) {
		cfg := envConfig{
			OAuthClientSecret: "client-secret",
			OAuthScopes:       []string{"scope"},
			OAuthDiscovery:    "http://some-url.com",
		}
		oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{}}

		c, err := NewOauthClientGetter(cfg, oauthDiscovery).Get()
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
	t.Run("Get - missing client secret", func(t *testing.T) {
		cfg := envConfig{
			OAuthClientID:  "client-id",
			OAuthScopes:    []string{"scope"},
			OAuthDiscovery: "http://some-url.com",
		}
		oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{}}

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

	cfg := envConfig{
		OAuthClientID:     "client-id",
		OAuthClientSecret: "client-secret",
		OAuthScopes:       []string{"scope"},
		OAuthDiscovery:    "http://some-wellknown-url.com",
	}
	oauthDiscovery := &oauthutils.StaticOauthDiscovery{DiscoveryValues: &oauthutils.OauthDiscoveryResult{
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
