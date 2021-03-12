package cmd

import (
	"fmt"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"net"
	"net/url"
	"time"
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
}

type AuthenticatorOptions struct {
	Endpoint string
	ApiToken string
}

func NewAuthenticator(namespace string, credentialManager CredentialGetSetter) *Authenticator {
	return &Authenticator{
		Namespace:         namespace,
		CredentialManager: credentialManager,
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
			return err
		}
	} else {
		endpoint, err = a.parseURL(authenticatorOptions.Endpoint)
		if err != nil {
			return err
		}
		apiToken = authenticatorOptions.ApiToken
	}

	logging.PrintLog("Starting to authenticate", logging.InfoLevel)

	if endpoint.Path == "" || endpoint.Path == "/" {
		endpoint.Path = "/api"
	}

	authHandler := apiutils.NewAuthenticatedAuthHandler(endpoint.String(), apiToken, "x-token", nil, endpoint.Scheme)

	if !LookupHostname(endpoint.Hostname(), net.LookupHost, time.Sleep) {
		return fmt.Errorf("Authentication was unsuccessful - could not resolve hostname.")
	}

	if endPointErr := CheckEndpointStatus(endpoint.String()); endPointErr != nil {
		return fmt.Errorf("Authentication was unsuccessful: %s"+endPointErrorReasons,
			endPointErr)
	}

	authenticated := false
	// try to authenticate (and retry it)
	for retries := 0; retries < 3; time.Sleep(5 * time.Second) {
		_, err := authHandler.Authenticate()
		if err != nil {
			errMsg := fmt.Sprintf("Authentication was unsuccessful. %s", *err.Message)
			logging.PrintLog(errMsg, logging.QuietLevel)
			logging.PrintLog("Retrying...", logging.InfoLevel)
			retries++
		} else {
			authenticated = true
			break
		}
	}

	if !authenticated {
		return fmt.Errorf("Authentication was unsuccessful - could not authenticate against the server.")
	}

	logging.PrintLog("Successfully authenticated against the Keptn cluster "+endpoint.String(), logging.InfoLevel)
	return a.CredentialManager.SetCreds(endpoint, apiToken, namespace)
}

func (a *Authenticator) parseURL(rawURL string) (url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}
	return *parsedURL, nil
}
