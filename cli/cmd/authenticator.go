package cmd

import (
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/internal"
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
	APIToken string
	OAuth    bool
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

func (a *Authenticator) GetAPISet() (*api.APISet, error) {
	endpoint, apiToken, err := a.CredentialManager.GetCreds(a.Namespace)
	if err != nil {
		return nil, err
	}
	return internal.APIProvider(endpoint.String(), apiToken)
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
		apiToken = authenticatorOptions.APIToken
	}

	logging.PrintLog("Starting to authenticate", logging.InfoLevel)

	if endpoint.Path == "" || endpoint.Path == "/" {
		endpoint.Path = "/api"
	}

	api, err := internal.APIProvider(endpoint.String(), apiToken)
	if err != nil {
		return err
	}

	if !LookupHostname(endpoint.Hostname(), net.LookupHost, time.Sleep) {
		return fmt.Errorf("Authentication was unsuccessful - could not resolve hostname.")
	}

	// Try to call the /auth endpoint of Keptn
	// NOTE: We expect that some component is setting the api token using the
	// "x-token" HTTP header. This could've either be done via the CLI
	// by using the `--api-token` flag or some other component (e.g. load balancer, ...)
	authenticated := false
	errMsg := ""
	// try to authenticate (and retry it)
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

	if !authenticated {
		if authenticatorOptions.OAuth && apiToken == "" {
			fmt.Println("WARNING: You are using the OAuth integration feature without a Keptn API Token. Please verify that your configuration allows it")
		} else {
			return fmt.Errorf(errMsg)
		}
	} else {
		logging.PrintLog("Successfully authenticated against the Keptn cluster "+endpoint.String(), logging.InfoLevel)
	}

	return a.CredentialManager.SetCreds(endpoint, apiToken, namespace)
}

func (a *Authenticator) parseURL(rawURL string) (url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return url.URL{}, err
	}
	return *parsedURL, nil
}
