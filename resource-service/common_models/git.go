package common_models

import (
	"net/url"
	"strings"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/go-git/go-git/v5/plumbing/transport"
	kerrors "github.com/keptn/keptn/resource-service/errors"
)

// GitCredentials contains git credentials info
type GitCredentials apimodels.GitAuthCredentials

type GitContext struct {
	Project     string
	Credentials *GitCredentials
	AuthMethod  transport.AuthMethod
}

func (g GitCredentials) Validate() error {
	if !strings.HasPrefix(g.RemoteURL, "http://") && !strings.HasPrefix(g.RemoteURL, "ssh://") && !strings.HasPrefix(g.RemoteURL, "https://") {
		return kerrors.ErrInvalidRemoteURL
	}
	if g.HttpsAuth != nil && !strings.HasPrefix(g.RemoteURL, "ssh://") {
		if err := g.validateRemoteURLAndToken(); err != nil {
			return err
		}
		if err := g.validateProxy(); err != nil {
			return err
		}
	} else if g.SshAuth != nil && strings.HasPrefix(g.RemoteURL, "ssh://") {
		if g.SshAuth.PrivateKey == "" {
			return kerrors.ErrCredentialsPrivateKeyMustNotBeEmpty
		}
	} else {
		return kerrors.ErrInvalidCredentials
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
	return nil
}
