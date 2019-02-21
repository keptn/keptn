package credentialmanager

import "github.com/docker/docker-credential-helpers/credentials"

const credsLab = "keptn"
const serverURL = "https://test.keptn.sh"
const username = "token"

type bot interface {
	SetCreds(secret string) error
	GetCreds() (string, error)
}

func setCreds(h credentials.Helper, secret string) error {
	credentials.SetCredsLabel(credsLab)
	c := &credentials.Credentials{
		ServerURL: serverURL,
		Username:  username,
		Secret:    secret,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper) (string, error) {
	_, secret, err := h.Get(serverURL)
	return secret, err
}
