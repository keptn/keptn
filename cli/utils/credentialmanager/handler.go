package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/credentials"
)

const credsLab = "keptn"
const serverURL = "https://keptn.sh"

type bot interface {
	SetCreds(endPoint string, apiToken string) error
	GetCreds() (string, string, error)
}

func setCreds(h credentials.Helper, endPoint string, apiToken string) error {
	credentials.SetCredsLabel(credsLab)
	c := &credentials.Credentials{
		ServerURL: serverURL,
		Username:  endPoint,
		Secret:    apiToken,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper) (string, string, error) {
	return h.Get(serverURL)
}
