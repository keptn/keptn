package credentialmanager

import (
	"net/url"

	"github.com/docker/docker-credential-helpers/wincred"
)

// SetCreds stores the credentials consisting of an endpoint and an api token in the keychain.
func SetCreds(endPoint url.URL, apiToken string) error {
	return setCreds(wincred.Wincred{}, endPoint, apiToken)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func GetCreds() (url.URL, string, error) {
	return getCreds(wincred.Wincred{})
}
