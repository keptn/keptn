package credentialmanager

import "github.com/docker/docker-credential-helpers/osxkeychain"

func SetCreds(endPoint string, secret string) error {
	return setCreds(osxkeychain.Osxkeychain{}, endPoint, secret)
}

func GetCreds() (string, string, error) {
	return getCreds(osxkeychain.Osxkeychain{})
}
