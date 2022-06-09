package common_models

import (
	"net/url"
	"strings"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	kerrors "github.com/keptn/keptn/resource-service/errors"
)

// GitCredentials contains git credentials info
type GitCredentials apimodels.GitAuthCredentials

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

func (g GitCredentials) Validate() error {
	if !strings.HasPrefix(g.RemoteURL, "http://") && !strings.HasPrefix(g.RemoteURL, "ssh://") && !strings.HasPrefix(g.RemoteURL, "https://") {
		return kerrors.ErrInvalidRemoteURL
	}
	if g.HttpsAuth != nil {
		if err := g.validateRemoteURLAndToken(); err != nil {
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

func (g GitCredentials) validateRemoteURLAndToken() error {
	_, err := url.Parse(g.RemoteURL)
	if err != nil {
		return kerrors.ErrCredentialsInvalidRemoteURL
	}
	if g.HttpsAuth.Token == "" {
		return kerrors.ErrCredentialsTokenMustNotBeEmpty
	}
	return nil
}
