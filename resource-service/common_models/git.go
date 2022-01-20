package common_models

import (
	kerrors "github.com/keptn/keptn/resource-service/errors"
	"net/url"
)

// GitCredentials contains git credentials info
type GitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

func (g GitCredentials) Validate() error {
	_, err := url.Parse(g.RemoteURI)
	if err != nil {
		return kerrors.ErrCredentialsInvalidRemoteURI
	}
	if g.Token == "" {
		return kerrors.ErrCredentialsTokenMustNotBeEmpty
	}
	return nil
}
