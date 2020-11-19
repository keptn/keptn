package credentialmanager

import (
	"net/url"

	"github.com/docker/docker-credential-helpers/osxkeychain"
)

type CredentialManager struct {
}

func NewCredentialManager(autoApplyNewContext bool) (cm *CredentialManager) {
	initChecks(autoApplyNewContext)
	return &CredentialManager{}
}

// SetCreds stores the credentials consisting of an endpoint and an api token in the keychain.
func (cm *CredentialManager) SetCreds(endPoint url.URL, apiToken string, namespace string) error {
	return setCreds(osxkeychain.Osxkeychain{}, endPoint, apiToken, namespace)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func (cm *CredentialManager) GetCreds(namespace string) (url.URL, string, error) {
	return getCreds(osxkeychain.Osxkeychain{}, namespace)
}

// SetInstallCreds sets the install credentials
func (cm *CredentialManager) SetInstallCreds(creds string) error {
	return setInstallCreds(osxkeychain.Osxkeychain{}, creds)
}

// GetInstallCreds gets the install credentials
func (cm *CredentialManager) GetInstallCreds() (string, error) {
	return getInstallCreds(osxkeychain.Osxkeychain{})
}
