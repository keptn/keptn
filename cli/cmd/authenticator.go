package cmd

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/keptn/keptn/cli/internal"
	"github.com/keptn/keptn/cli/internal/auth"
)

type MockedCredentialGetSetter struct {
	SetCredsFunc func(endPoint url.URL, apiToken string, namespace string) error
	GetCredsFunc func(namespace string) (url.URL, string, error)
}

func (m *MockedCredentialGetSetter) SetCreds(endPoint url.URL, apiToken string, namespace string) error {
	return m.SetCredsFunc(endPoint, apiToken, namespace)
}
func (m *MockedCredentialGetSetter) GetCreds(namespace string) (url.URL, string, error) {
	return m.GetCredsFunc(namespace)
}

type CredentialGetSetter interface {
	SetCreds(endPoint url.URL, apiToken string, namespace string) error
	GetCreds(namespace string) (url.URL, string, error)
}

type Authenticator struct {
	Namespace         string
	CredentialManager CredentialGetSetter
	OauthStore        auth.OauthStore
}

type AuthenticatorOptions struct {
	Endpoint string
	APIToken string
}

func NewAuthenticator(namespace string, credentialManager CredentialGetSetter, oauthStore auth.OauthStore) *Authenticator {
	return &Authenticator{
		Namespace:         namespace,
		CredentialManager: credentialManager,
		OauthStore:        oauthStore,
	}
}

func (a *Authenticator) GetCredentials() (url.URL, string, error) {
	return a.CredentialManager.GetCreds(a.Namespace)
}

func (a *Authenticator) Auth(authenticatorOptions AuthenticatorOptions) error {
	var endpoint url.URL
	var apiToken string
	var err error
	if authenticatorOptions.Endpoint == "" {
		endpoint, apiToken, err = a.CredentialManager.GetCreds(a.Namespace)
		if err != nil {
			return internal.OnAPIError(err)
		}
	} else {
		endpoint, err = a.parseURL(authenticatorOptions.Endpoint)
		if err != nil {
			return internal.OnAPIError(err)
		}
		apiToken = authenticatorOptions.APIToken
	}

	fmt.Println("Starting to authenticate")

	if endpoint.Path == "" || endpoint.Path == "/" {
		endpoint.Path = "/api"
	}

	api, err := internal.APIProvider(endpoint.String(), apiToken)
	if err != nil {
		return err
	}

	if !LookupHostname(endpoint.Hostname(), net.LookupHost, time.Sleep) {
		return fmt.Errorf("Authentication was unsuccessful - could not resolve hostname")
	}

	// Skip usual auth call if we use OAuth
	if a.OauthStore.Created() {
		fmt.Printf("Successfully authenticated against the Keptn cluster %s\n", endpoint.String())
		fmt.Printf("Bridge URL: %s\n", getBridgeURLFromAPIURL(endpoint))
		return a.CredentialManager.SetCreds(endpoint, apiToken, namespace)
	}

	// Try to call Keptn Auth endpoint
	authenticated := false
	errMsg := ""
	for retries := 0; retries < 3; time.Sleep(5 * time.Second) {
		_, err := api.AuthV1().Authenticate()
		if err != nil {
			errMsg = fmt.Sprintf("Authentication was unsuccessful. %s", *err.Message)
			retries++
		} else {
			authenticated = true
			break
		}
	}

	// Authentication failed
	if !authenticated {
		return fmt.Errorf(errMsg)
	}

	// Authentication succeeded
	fmt.Println("Successfully authenticated against the Keptn cluster " + endpoint.String())
	fmt.Println("Bridge URL: " + getBridgeURLFromAPIURL(endpoint))
	return a.CredentialManager.SetCreds(endpoint, apiToken, namespace)
}

func (a *Authenticator) parseURL(rawURL string) (url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}
	return *parsedURL, nil
}
