package common_models

import (
	git2go "github.com/libgit2/git2go/v34"
	"net/url"
	"strings"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/go-git/go-git/v5/plumbing/transport"
	kerrors "github.com/keptn/keptn/resource-service/errors"
)

const GitInitDefaultBranchName = "master"

// GitCredentials contains git credentials info
type GitCredentials apimodels.GitAuthCredentials

type Git2GoAuth struct {
	CredCallback git2go.CredentialsCallback
	CertCallback git2go.CertificateCheckCallback
	ProxyOptions git2go.ProxyOptions
}

type AuthMethod struct {
	GoGitAuth  transport.AuthMethod
	Git2GoAuth Git2GoAuth
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
	AuthMethod  AuthMethod
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
