package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/pass"
)

func SetCreds(secret string) error {
	return setCreds(pass.Pass{}, secret)
}

func GetCreds() (string, error) {
	return getCreds(pass.Pass{})
}
