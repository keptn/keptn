package common_models

import (
	"net/url"

	kerrors "github.com/keptn/keptn/resource-service/errors"
)

// GitCredentials contains git credentials info
type GitCredentials struct {
	User       string `json:"user,omitempty"`
	Token      string `json:"token,omitempty"`
	RemoteURI  string `json:"remoteURI,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
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
