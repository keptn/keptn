package common_models

import (
	"net/url"
	"strings"

	kerrors "github.com/keptn/keptn/resource-service/errors"
)

// GitCredentials contains git credentials info
type GitCredentials struct {
	User              string `json:"user,omitempty"`
	Token             string `json:"token,omitempty"`
	RemoteURI         string `json:"remoteURI,omitempty"`
	GitPrivateKey     string `json:"privateKey,omitempty"`
	GitPrivateKeyPass string `json:"privateKeyPass,omitempty"`
	GitProxyURL       string `json:"gitProxyUrl,omitempty"`
	GitProxyScheme    string `json:"gitProxyScheme,omitempty"`
	GitProxyUser      string `json:"gitProxyUser,omitempty"`
	GitProxyPassword  string `json:"gitProxyPassword,omitempty"`
	GitPemCertificate string `json:"gitPemCertificate,omitempty"`
	// omitempty property is missing due to fallback of this
	// parameter to "undefined" when marshalling/unmarshalling data
	// when "false" value is present
	GitProxyInsecure bool `json:"gitProxyInsecure"`
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

func (g GitCredentials) Validate() error {
	if strings.HasPrefix(g.RemoteURI, "https://") || strings.HasPrefix(g.RemoteURI, "http://") {
		if err := g.validateRemoteURIAndToken(); err != nil {
			return err
		}
		if err := g.validateProxy(); err != nil {
			return err
		}
	} else if strings.HasPrefix(g.RemoteURI, "ssh://") {
		if g.GitPrivateKey == "" {
			return kerrors.ErrCredentialsPrivateKeyMustNotBeEmpty
		}
	} else {
		return kerrors.ErrCredentialsInvalidRemoteURI
	}
	return nil
}

func (g GitCredentials) validateProxy() error {
	if g.GitProxyURL != "" {
		if g.GitProxyScheme != "http" && g.GitProxyScheme != "https" {
			return kerrors.ErrProxyInvalidScheme
		}
		if !strings.Contains(g.GitProxyURL, ":") {
			return kerrors.ErrProxyInvalidURL
		}
	}
	return nil
}

func (g GitCredentials) validateRemoteURIAndToken() error {
	_, err := url.Parse(g.RemoteURI)
	if err != nil {
		return kerrors.ErrCredentialsInvalidRemoteURI
	}
	if g.Token == "" {
		return kerrors.ErrCredentialsTokenMustNotBeEmpty
	}
	return nil
}
