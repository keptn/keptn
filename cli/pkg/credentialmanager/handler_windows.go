package credentialmanager

import (
	"net/url"

	"github.com/docker/docker-credential-helpers/wincred"
)

type CredentialManager struct {
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager(autoApplyNewContext bool) (cm *CredentialManager) {
	cm := &CredentialManager{}
	initChecks(autoApplyNewContext, cm)
	return
}

// SetCreds stores the credentials consisting of an endpoint and an api token in the keychain.
func (cm *CredentialManager) SetCreds(endPoint url.URL, apiToken string, namespace string) error {
	return setCreds(wincred.Wincred{}, endPoint, apiToken, namespace)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func (cm *CredentialManager) GetCreds(namespace string) (url.URL, string, error) {
	return getCreds(wincred.Wincred{}, namespace)
}

// SetInstallCreds sets the install credentials
func (cm *CredentialManager) SetInstallCreds(creds string) error {
	return setInstallCreds(wincred.Wincred{}, creds)
}

// GetInstallCreds gets the install credentials
func (cm *CredentialManager) GetInstallCreds() (string, error) {
	return getInstallCreds(wincred.Wincred{})
}

func (cm *CredentialManager) SetCurrentKubeConfig(kubeConfig KubeConfigFileType) {
	cm.kubeConfig = kubeConfig
}

func (cm *CredentialManager) GetCurrentKubeConfig() KubeConfigFileType {
	return cm.kubeConfig
}

func (cm *CredentialManager) SetCurrentKeptnCLIConfig(cliConfig config.CLIConfig) {
	cm.cliConfig = cliConfig
}

func (cm *CredentialManager) GetCurrentKeptnCLIConfig() config.CLIConfig {
	return cm.cliConfig
}
