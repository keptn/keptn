package client

import (
	"context"
	"crypto/tls"
	"fmt"
	auth "github.com/keptn/go-utils/pkg/common/oauth"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)

// CreateClientGetter returns a HTTPClientGetter implementation based on the values certain properties
// inside the given env configuration
func CreateClientGetter(envConfig config.EnvConfig) HTTPClientGetter {
	if envConfig.SSOEnabled() {
		logger.Infof("Using Oauth to connect to Keptn wth client ID %s and scopes %v", envConfig.SSOClientID, envConfig.SSOScopes)
		return NewSSOClientGetter(envConfig, auth.NewOauthDiscovery(&http.Client{}))
	}
	return New(envConfig)
}

// HTTPClientGetter is responsible for creating an HTTP client
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
	if g.envConfig.SSOClientID == "" || g.envConfig.SSOClientSecret == "" || len(g.envConfig.SSOScopes) == 0 {
		return nil, fmt.Errorf("client id or client secret or scopes missing")
	}

	if g.envConfig.SSOTokenURL != "" {
		logger.Infof("Using Token URL for Oauth flow: %s", g.envConfig.SSOTokenURL)
		conf := clientcredentials.Config{
			ClientID:     g.envConfig.SSOClientID,
			ClientSecret: g.envConfig.SSOClientSecret,
			Scopes:       g.envConfig.SSOScopes,
			TokenURL:     g.envConfig.SSOTokenURL,
		}
		return conf.Client(context.WithValue(context.TODO(), oauth2.HTTPClient, c)), nil
	}

	if g.envConfig.SSODiscovery != "" {
		logger.Infof("Using Discovery URL for Oauth flow: %s", g.envConfig.SSODiscovery)
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()
		discoveryRes, err := g.oauthDiscovery.Discover(ctx, g.envConfig.SSODiscovery)
		if err != nil {
			return nil, err
		}

		conf := clientcredentials.Config{
			ClientID:     g.envConfig.SSOClientID,
			ClientSecret: g.envConfig.SSOClientSecret,
			Scopes:       g.envConfig.SSOScopes,
			TokenURL:     discoveryRes.TokenEndpoint,
		}
		return conf.Client(context.WithValue(context.TODO(), oauth2.HTTPClient, c)), nil
	}
	return nil, fmt.Errorf("no discovery or token url is provided")
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
