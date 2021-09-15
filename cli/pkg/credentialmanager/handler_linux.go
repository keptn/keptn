package credentialmanager

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"io/ioutil"

	"github.com/docker/docker-credential-helpers/pass"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

// For using pass.Pass{} the following commands need to executed:
// 1. sudo apt-get install gpg pass -y
// 2. gpg --generate-key (Use your name and e-mail); Use "find / | xargs file" for generating random bytes; Copy generate pub key
// 3. pass init [generated pub key]

var passwordStoreDirectory string

func init() {
	passwordStoreDirectory = os.Getenv("HOME") + "/.password-store"
}

type CredentialManager struct {
	// MockAuthCreds shows whether the get and set for the auth-creds should be mocked
	apiTokenFile string
	credsFile    string
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager(autoApplyNewContext bool) (cm *CredentialManager) {

	dir, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		log.Fatal(err)
	}
	cm := &CredentialManager{apiTokenFile: dir + ".keptn", credsFile: dir + ".keptn-creds"}
	initChecks(autoApplyNewContext, cm)
	return cm
}

// SetCreds stores the credentials consisting of an endpoint and an api token using pass or into a file in case
// pass is unavailable.
func (cm *CredentialManager) SetCreds(endPoint url.URL, apiToken string, namespace string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		fmt.Println("Using a file-based storage for the key because the password-store seems to be not set up.")
		apiTokenFile := cm.getLinuxApiTokenFile(namespace)
		return ioutil.WriteFile(apiTokenFile, []byte(endPoint.String()+"\n"+apiToken), 0644)
	}
	return setCreds(pass.Pass{}, endPoint, apiToken, namespace)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func (cm *CredentialManager) GetCreds(namespace string) (url.URL, string, error) {
	// mock credentials if necessary
	if MockAuthCreds {
		return url.URL{}, "", nil
	}

	// Check if creds file is specified in the 'KEPTNCONFIG' environment variable
	if customCredsLocation, ok := os.LookupEnv("KEPTNCONFIG"); ok {
		if customCredsLocation != "" {
			endPoint, apiToken, err := handleCustomCreds(customCredsLocation, namespace)
			// If credential is not found in KEPTNCONFIG, use fallback credential manager
			if apiToken != "" || err != nil {
				return endPoint, apiToken, err
			}
		}
	}

	// try to read credentials from password-store
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		// password-store not found, read credentials from apiTokenFile
		apiTokenFile := cm.getLinuxApiTokenFile(namespace)
		data, err := ioutil.ReadFile(apiTokenFile)
		if err != nil {
			return url.URL{}, "", err
		}
		dataStr := strings.TrimSpace(strings.Replace(string(data), "\r\n", "\n", -1))
		creds := strings.Split(dataStr, "\n")
		if len(creds) != 2 {
			return url.URL{}, "", errors.New("Format of file-based key storage is invalid")
		}
		url, err := url.Parse(creds[0])
		return *url, creds[1], err
	}
	return getCreds(pass.Pass{}, namespace)
}

// SetInstallCreds sets the install credentials
func (cm *CredentialManager) SetInstallCreds(creds string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		fmt.Println("Using a file-based storage for the key because the password-store seems to be not set up.")

		return ioutil.WriteFile(cm.credsFile, []byte(creds), 0644)
	}
	return setInstallCreds(pass.Pass{}, creds)
}

// GetInstallCreds gets the install credentials
func (cm *CredentialManager) GetInstallCreds() (string, error) {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		data, err := ioutil.ReadFile(cm.credsFile)
		if err != nil {
			return "", err
		}
		dataStr := strings.TrimSpace(strings.Replace(string(data), "\r\n", "\n", -1))
		return dataStr, nil
	}
	return getInstallCreds(pass.Pass{})
}

func (cm *CredentialManager) getLinuxApiTokenFile(namespace string) string {
	sanitizedCurrentContext := strings.ReplaceAll(keptnContext, "/", "-")
	return cm.apiTokenFile + "__" + sanitizedCurrentContext + "__" + namespace
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
