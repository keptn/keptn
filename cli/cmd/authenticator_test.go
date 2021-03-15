package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAuthenticate(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			return
		}),
	)

	t.Run("TestAuthenticate", func(t *testing.T) {
		tsURL, _ := url.Parse(ts.URL)
		expectedURL, _ := url.Parse(ts.URL + "/api")

		credentialManagerMock := MockedCredentialGetSetter{
			SetCredsFunc: func(endPoint url.URL, apiToken string, namespace string) error {
				assert.Equal(t, *expectedURL, endPoint)
				return nil
			},
			GetCredsFunc: func(namespace string) (url.URL, string, error) {
				return *tsURL, "TOKEN", nil
			},
		}

		instance := NewAuthenticator("keptn", &credentialManagerMock)
		err := instance.Auth(AuthenticatorOptions{})
		assert.Nil(t, err)
	})

	t.Run("TestAuthenticate_WithExplicitEndpointAndToken", func(t *testing.T) {
		tsURL, _ := url.Parse(ts.URL)
		expectedURL, _ := url.Parse(ts.URL + "/api")

		credentialManagerMock := MockedCredentialGetSetter{
			SetCredsFunc: func(endPoint url.URL, apiToken string, namespace string) error {
				assert.Equal(t, *expectedURL, endPoint)
				return nil
			},
			GetCredsFunc: func(namespace string) (url.URL, string, error) {
				return *tsURL, "", nil
			},
		}
		instance := NewAuthenticator("keptn", &credentialManagerMock)
		err := instance.Auth(AuthenticatorOptions{Endpoint: ts.URL})
		assert.Nil(t, err)
	})
	t.Run("TestAuthenticate_GettingCredentialsFails", func(t *testing.T) {
		tsURL, _ := url.Parse(ts.URL)

		credentialManagerMock := MockedCredentialGetSetter{
			SetCredsFunc: func(endPoint url.URL, apiToken string, namespace string) error {
				return nil
			},
			GetCredsFunc: func(namespace string) (url.URL, string, error) {
				return *tsURL, "TOKEN", fmt.Errorf("whoops")
			},
		}
		instance := NewAuthenticator("keptn", &credentialManagerMock)
		err := instance.Auth(AuthenticatorOptions{})
		assert.NotNil(t, err)
	})

	t.Run("TestAuthenticate_SettingCredentialsFails", func(t *testing.T) {
		tsURL, _ := url.Parse(ts.URL)

		credentialManagerMock := MockedCredentialGetSetter{
			SetCredsFunc: func(endPoint url.URL, apiToken string, namespace string) error {
				return fmt.Errorf("whoops")
			},
			GetCredsFunc: func(namespace string) (url.URL, string, error) {
				return *tsURL, "TOKEN", nil
			},
		}
		instance := NewAuthenticator("keptn", &credentialManagerMock)
		err := instance.Auth(AuthenticatorOptions{})
		assert.NotNil(t, err)
	})

}
