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
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

func (g GitCredentials) Validate() error {
	if strings.HasPrefix(g.RemoteURI, "https://") || strings.HasPrefix(g.RemoteURI, "http://") {
		_, err := url.Parse(g.RemoteURI)
		if err != nil {
			return kerrors.ErrCredentialsInvalidRemoteURI
		}
		if g.Token == "" {
			return kerrors.ErrCredentialsTokenMustNotBeEmpty
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
