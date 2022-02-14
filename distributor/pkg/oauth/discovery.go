package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OauthLocationGetter is used to get the location parameters
// used in an oauth flow
type OauthLocationGetter interface {
	// Discover is responsible for determining the parameters used for an oauth flow
	// and returns them as a OauthDiscoveryResult
	Discover(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error)
}

// NewOauthDiscovery creates a new OauthDiscovery
func NewOauthDiscovery(client HTTPClient) *OauthDiscovery {
	return &OauthDiscovery{
		c: client,
	}
}

// OauthDiscoveryResult is the result of a OauthLocation discovery call
// and contains all the parameters usable for a following oauth flow
type OauthDiscoveryResult struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	UserinfoEndpoint       string   `json:"userinfo_endpoint"`
	EndSessionEndpoint     string   `json:"end_session_endpoint"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	GrantTypesSupported    []string `json:"grant_types_supported"`
	JwksURI                string   `json:"jwks_uri"`
	IntrospectionEndpoint  string   `json:"introspection_endpoint"`
}

// OauthDiscovery is an implementation of OauthLocationGetter which calls
// a known URL to get the parameters
type OauthDiscovery struct {
	c HTTPClient
}

// Discover starts the OAuth discovery and eventually returns the results
func (d OauthDiscovery) Discover(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := d.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(http.StatusText(resp.StatusCode))
	}

	var result OauthDiscoveryResult
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}

// StaticOauthDiscovery has a static/hard-coded set of oauth parameters
// and does not actually do a discovery call to get the parameters but just returns
// the hard-coded values as OauthDiscoveryResult
type StaticOauthDiscovery struct {
	DiscoveryValues *OauthDiscoveryResult
}

// Discover tries to determine the parameters used for an oauth flow
// and returns them as a OauthDiscoveryResult
func (d StaticOauthDiscovery) Discover(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
	return d.DiscoveryValues, nil
}

// OauthDiscoveryMock is an implementation of OauthLocationGetter usable
// as a mock implementation in tests
type OauthDiscoveryMock struct {
	DiscoverFn func(context.Context, string) (*OauthDiscoveryResult, error)
}

// Discover calls the mocked function of the OauthDiscoveryMock
func (o *OauthDiscoveryMock) Discover(ctx context.Context, discoveryURL string) (*OauthDiscoveryResult, error) {
	if o != nil && o.DiscoverFn != nil {
		return o.DiscoverFn(ctx, discoveryURL)
	}
	return &OauthDiscoveryResult{}, nil
}

// HTTPClient is an interface that models *http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
