package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/wincred"
)

func SetCreds(endPoint string, secret string) error {
	return setCreds(wincred.Wincred{}, endPoint, secret)
}

func GetCreds() (string, string, error) {
	return getCreds(wincred.Wincred{})
}
