package client

import (
	"context"
	"crypto/tls"
	auth "github.com/keptn/go-utils/pkg/common/oauth"
	"github.com/keptn/keptn/distributor/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)

// CreateClientGetter returns a HTTPClientGetter implementation based on the values certain properties
// inside the given env configuration
func CreateClientGetter(envConfig config.EnvConfig) HTTPClientGetter {
	if envConfig.SSOEnabled() {
		return NewSSOClientGetter(envConfig, auth.NewOauthDiscovery(&http.Client{}))
	}
	return New(envConfig)
}

// HTTPClientGetter is responsible for creating a HTTP client
type HTTPClientGetter interface {
	// Get Creates the HTTP Client
	Get() (*http.Client, error)
}

// SSOClientGetter creates an HTTP client configured for use with SSO/Oauth
type SSOClientGetter struct {
	*SimpleClientGetter
	envConfig      config.EnvConfig
	oauthDiscovery auth.OauthLocationGetter
}

// NewSSOClientGetter creates a new instance of a SSOClientGetter
func NewSSOClientGetter(envConfig config.EnvConfig, oauthDiscovery auth.OauthLocationGetter) *SSOClientGetter {
	return &SSOClientGetter{
		SimpleClientGetter: &SimpleClientGetter{envConfig: envConfig},
		envConfig:          envConfig,
		oauthDiscovery:     oauthDiscovery,
	}
}

func (g *SSOClientGetter) Get() (*http.Client, error) {
	c, err := g.SimpleClientGetter.Get()
	if err != nil {
		return nil, err
	}

	conf := clientcredentials.Config{
		ClientID:     g.envConfig.SSOClientID,
		ClientSecret: g.envConfig.SSOClientSecret,
		Scopes:       g.envConfig.SSOScopes,
		TokenURL:     g.envConfig.SSOTokenURL,
	}
	return conf.Client(context.WithValue(context.TODO(), oauth2.HTTPClient, c)), nil
}

// SimpleClientGetter creates a basic HTTP client
type SimpleClientGetter struct {
	envConfig config.EnvConfig
}

// New Creates a new instance of a SimpleClientGetter
func New(envConfig config.EnvConfig) *SimpleClientGetter {
	return &SimpleClientGetter{envConfig: envConfig}
}

func (g *SimpleClientGetter) Get() (*http.Client, error) {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !g.envConfig.VerifySSL}, //nolint:gosec
		},
		Timeout: 5 * time.Second,
	}
	return c, nil
}
