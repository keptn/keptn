package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/secretservice"
)

func SetCreds(secret string) error {
	return setCreds(secretservice.Secretservice{}, secret)
}

func GetCreds() (string, error) {
	return getCreds(secretservice.Secretservice{})
}
