package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/wincred"
)

func SetCreds(secret string) error {
	return setCreds(wincred.Wincred{}, secret)
}

func GetCreds() (string, error) {
	return getCreds(wincred.Wincred{})
}
