package credentialmanager

import (
	"github.com/docker/docker-credential-helpers/osxkeychain"
)

func SetCreds(secret string) error {
	return setCreds(osxkeychain.Osxkeychain{}, secret)
}

func GetCreds() (string, error) {
	return getCreds(osxkeychain.Osxkeychain{})
}
