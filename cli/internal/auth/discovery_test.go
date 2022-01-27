package auth

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDiscovery(t *testing.T) {
	response := `{"issuer":"https://issuer.com:443","authorization_endpoint":"https://endpoint.com:443/oauth2/authorize","token_endpoint":"https://tokenendpoint.com:443/sso/oauth2/token","userinfo_endpoint":"https://userinfo.com:443/sso/oauth2/userinfo","end_session_endpoint":"https://endsession.com:443/oauth2/end_session","response_types_supported":["code"],"grant_types_supported":["authorization_code","refresh_token","password","client_credentials","urn:ietf:params:oauth:grant-type:token-exchange"],"jwks_uri":"https://jwksuri.com:443/.well-known/jwks.json","introspection_endpoint":"https://introspection.com:443/sso/oauth2/tokeninfo","subject_types_supported":["public"],"id_token_signing_alg_values_supported":["ECDSA256"]}`
	client := &HTTPClientMock{}
	discovery := NewOauthDiscovery(client)
	tt := []struct {
		Body       string
		StatusCode int
		expResult  *OauthDiscoveryResult
		expErr     error
	}{
		{
			Body:       response,
			StatusCode: 200,
			expResult: &OauthDiscoveryResult{
				Issuer:                 "https://issuer.com:443",
				AuthorizationEndpoint:  "https://endpoint.com:443/oauth2/authorize",
				TokenEndpoint:          "https://tokenendpoint.com:443/sso/oauth2/token",
				UserinfoEndpoint:       "https://userinfo.com:443/sso/oauth2/userinfo",
				EndSessionEndpoint:     "https://endsession.com:443/oauth2/end_session",
				ResponseTypesSupported: []string{"code"},
				GrantTypesSupported:    []string{"authorization_code", "refresh_token", "password", "client_credentials", "urn:ietf:params:oauth:grant-type:token-exchange"},
				JwksURI:                "https://jwksuri.com:443/.well-known/jwks.json",
				IntrospectionEndpoint:  "https://introspection.com:443/sso/oauth2/tokeninfo",
			},
			expErr: nil,
		}, {
			Body:       "",
			StatusCode: http.StatusNotFound,
			expResult:  nil,
			expErr:     fmt.Errorf(http.StatusText(http.StatusNotFound)),
		}, {
			Body:       "",
			StatusCode: http.StatusBadRequest,
			expResult:  nil,
			expErr:     fmt.Errorf(http.StatusText(http.StatusBadRequest)),
		},
	}

	for _, test := range tt {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}
		r, err := discovery.Discover(context.Background(), "http://well-known-discovery-url.com")

		assert.Equal(t, test.expErr, err)
		assert.Equal(t, test.expResult, r)
	}

	result, err := discovery.Discover(context.TODO(), "http://well-known-discovery-url.com")
	if err != nil {
		return
	}
	fmt.Println(result.TokenEndpoint)
}
