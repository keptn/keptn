package client

import (
	auth "github.com/keptn/go-utils/pkg/common/oauth"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleClientGetter_Get(t *testing.T) {
	cfg := config.EnvConfig{}
	getter := New(cfg)
	c, err := getter.Get()
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestSSOClientGetter_Get(t *testing.T) {
	cfg := config.EnvConfig{
		SSOClientID:     "client-id",
		SSOClientSecret: "client-secret",
		SSOScopes:       []string{"scope"},
		SSOTokenURL:     "http://token-url.com/token",
	}
	oauthDiscovery := &auth.OauthDiscoveryMock{}
	getter, err := NewSSOClientGetter(cfg, oauthDiscovery).Get()
	assert.NotNil(t, getter)
	assert.Nil(t, err)
}
