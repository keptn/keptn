package auth

func NewOauthDiscovery(discoveryURL string) *OauthDiscovery {
	return &OauthDiscovery{
		discoveryURL,
	}
}

type OauthDiscovery struct {
	discoveryURL string
}

type OauthLocationGetter interface {
	Discover() (*OauthDiscoveryResult, error)
}

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

func (d OauthDiscovery) Discover() (*OauthDiscoveryResult, error) {
	panic("not yet implemented")
	return nil, nil
}

type StaticOauthDiscovery struct {
}

func (d StaticOauthDiscovery) Discover() (*OauthDiscoveryResult, error) {
	return &OauthDiscoveryResult{
		Issuer:                 "https://sso-dev.dynatracelabs.com:443",
		AuthorizationEndpoint:  "https://sso-dev.dynatracelabs.com:443/oauth2/authorize",
		TokenEndpoint:          "https://sso-dev.dynatracelabs.com:443/sso/oauth2/token",
		UserinfoEndpoint:       "https://sso-dev.dynatracelabs.com:443/sso/oauth2/userinfo",
		EndSessionEndpoint:     "https://sso-dev.dynatracelabs.com:443/oauth2/end_session",
		ResponseTypesSupported: []string{"code"},
		GrantTypesSupported:    []string{"authorization_code", "refresh_token", "password", "client_credentials", "urn:ietf:params:oauth:grant-type:token-exchange"},
		JwksURI:                "https://sso-dev.dynatracelabs.com:443/.well-known/jwks.json",
		IntrospectionEndpoint:  "https://sso-dev.dynatracelabs.com:443/sso/oauth2/tokeninfo",
	}, nil
}

type OauthDiscoveryMock struct {
	discoverFn func() (*OauthDiscoveryResult, error)
}

func (o *OauthDiscoveryMock) Discover() (*OauthDiscoveryResult, error) {
	if o != nil && o.discoverFn != nil {
		return o.discoverFn()
	}
	return &OauthDiscoveryResult{}, nil

}
