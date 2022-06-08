package common_models

import (
	"net/url"
	"strings"

	kerrors "github.com/keptn/keptn/resource-service/errors"
)

// GitCredentials contains git credentials info
type GitCredentials struct {
	User      string        `json:"user,omitempty"`
	RemoteURL string        `json:"remoteURL,omitempty"`
	HttpsAuth *HttpsGitAuth `json:"https,omitempty"`
	SshAuth   *SshGitAuth   `json:"ssh,omitempty"`
}

// HttpsGitAuth stores HTTPS git credentials
type HttpsGitAuth struct {
	Token       string `json:"token"`
	Certificate string `json:"certificate,omitempty"`
	// omitempty property is missing due to fallback of this
	// parameter to "undefined" when marshalling/unmarshalling data
	// when "false" value is present
	InsecureSkipTLS bool          `json:"insecureSkipTLS"`
	Proxy           *ProxyGitAuth `json:"proxy,omitempty"`
}

// SshGitAuth stores SSH git credentials
type SshGitAuth struct {
	PrivateKey     string `json:"privateKey"`
	PrivateKeyPass string `json:"privateKeyPass,omitempty"`
}

// ProxyGitAuth stores proxy git credentials
type ProxyGitAuth struct {
	URL      string `json:"url"`
	Scheme   string `json:"scheme"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

func (g GitCredentials) Validate() error {
	if g.HttpsAuth != nil {
		if err := g.validateRemoteURIAndToken(); err != nil {
			return err
		}
		if err := g.validateProxy(); err != nil {
			return err
		}
	} else if g.SshAuth != nil {
		if g.SshAuth.PrivateKey == "" {
			return kerrors.ErrCredentialsPrivateKeyMustNotBeEmpty
		}
	} else {
		return kerrors.ErrCredentialsInvalidRemoteURL
	}
	return nil
}

func (g GitCredentials) validateProxy() error {
	if g.HttpsAuth.Proxy != nil {
		if g.HttpsAuth.Proxy.Scheme != "http" && g.HttpsAuth.Proxy.Scheme != "https" {
			return kerrors.ErrProxyInvalidScheme
		}
		if !strings.Contains(g.HttpsAuth.Proxy.URL, ":") {
			return kerrors.ErrProxyInvalidURL
		}
	}
	return nil
}

func (g GitCredentials) validateRemoteURIAndToken() error {
	_, err := url.Parse(g.RemoteURL)
	if err != nil {
		return kerrors.ErrCredentialsInvalidRemoteURL
	}
	if g.HttpsAuth.Token == "" {
		return kerrors.ErrCredentialsTokenMustNotBeEmpty
	}
	return nil
}
