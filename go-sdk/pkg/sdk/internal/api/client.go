package api

import (
	"context"
	"crypto/tls"
	"fmt"
	oauthutils "github.com/keptn/go-utils/pkg/common/oauth2"
	"github.com/keptn/keptn/go-sdk/pkg/sdk/internal/config"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)

// CreateClientGetter returns a HTTPClientGetter implementation based on the values certain properties
// inside the given env configuration
func CreateClientGetter(envConfig config.EnvConfig) HTTPClientGetter {
	if envConfig.OAuthEnabled() {
		logger.Infof("Using Oauth to connect to Keptn wth client ID %s and scopes %v", envConfig.OAuthClientID, envConfig.OAuthScopes)
		return NewOauthClientGetter(envConfig, oauthutils.NewOauthDiscovery(&http.Client{}))
	}
	return New(envConfig)
}

// HTTPClientGetter is responsible for creating an HTTP client
type HTTPClientGetter interface {
	// Get Creates the HTTP Client
	Get() (*http.Client, error)
}

// OAuthClientGetter creates an HTTP client configured for use with SSO/Oauth
type OAuthClientGetter struct {
	*SimpleClientGetter
	envConfig      config.EnvConfig
	oauthDiscovery oauthutils.OauthLocationGetter
}

// NewOauthClientGetter creates a new instance of a OAuthClientGetter
func NewOauthClientGetter(envConfig config.EnvConfig, oauthDiscovery oauthutils.OauthLocationGetter) *OAuthClientGetter {
	return &OAuthClientGetter{
		SimpleClientGetter: &SimpleClientGetter{envConfig: envConfig},
		envConfig:          envConfig,
		oauthDiscovery:     oauthDiscovery,
	}
}

func (g *OAuthClientGetter) Get() (*http.Client, error) {
	c, err := g.SimpleClientGetter.Get()
	if err != nil {
		return nil, err
	}
	if g.envConfig.OAuthClientID == "" || g.envConfig.OAuthClientSecret == "" || len(g.envConfig.OAuthScopes) == 0 {
		return nil, fmt.Errorf("client id or client secret or scopes missing")
	}

	if g.envConfig.OauthTokenURL != "" {
		logger.Infof("Using Token URL for Oauth flow: %s", g.envConfig.OauthTokenURL)
		conf := clientcredentials.Config{
			ClientID:     g.envConfig.OAuthClientID,
			ClientSecret: g.envConfig.OAuthClientSecret,
			Scopes:       g.envConfig.OAuthScopes,
			TokenURL:     g.envConfig.OauthTokenURL,
		}
		return conf.Client(context.WithValue(context.TODO(), oauth2.HTTPClient, c)), nil
	}

	if g.envConfig.OAuthDiscovery != "" {
		logger.Infof("Using Discovery URL for Oauth flow: %s", g.envConfig.OAuthDiscovery)
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()
		discoveryRes, err := g.oauthDiscovery.Discover(ctx, g.envConfig.OAuthDiscovery)
		if err != nil {
			return nil, err
		}

		conf := clientcredentials.Config{
			ClientID:     g.envConfig.OAuthClientID,
			ClientSecret: g.envConfig.OAuthClientSecret,
			Scopes:       g.envConfig.OAuthScopes,
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
		Timeout: g.envConfig.GetAPIProxyHTTPTimeout(),
	}
	return c, nil
}
