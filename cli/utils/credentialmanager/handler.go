package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/credentials"
)

const credsLab = "keptn"
const serverURL = "https://keptn.sh"
const username = "token"

type bot interface {
	SetCreds(secret string) error
	GetCreds() (string, error)
}

func setCreds(h credentials.Helper, endPoint string, secret string) error {
	credentials.SetCredsLabel(credsLab)
	c := &credentials.Credentials{
		ServerURL: serverURL,
		Username:  endPoint,
		Secret:    secret,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper) (string, string, error) {
	return h.Get(serverURL)
}
