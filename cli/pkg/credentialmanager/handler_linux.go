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

func NewCredentialManager() (cm *CredentialManager) {

	dir, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		log.Fatal(err)
	}
	return &CredentialManager{apiTokenFile: dir + ".keptn", credsFile: dir + ".keptn-creds"}
}

// SetCreds stores the credentials consisting of an endpoint and an api token using pass or into a file in case
// pass is unavailable.
func (cm *CredentialManager) SetCreds(endPoint url.URL, apiToken string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		fmt.Println("Using a file-based storage for the key because the password-store seems to be not set up.")

		return ioutil.WriteFile(cm.apiTokenFile, []byte(endPoint.String()+"\n"+apiToken), 0644)
	}
	return setCreds(pass.Pass{}, endPoint, apiToken)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func (cm *CredentialManager) GetCreds() (url.URL, string, error) {
	// mock credentials if encessary
	if MockAuthCreds {
		return url.URL{}, "", nil
	}

	// try to read credentials from password-store
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		// password-store not found, read credentials from apiTokenFile
		data, err := ioutil.ReadFile(cm.apiTokenFile)
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
	return getCreds(pass.Pass{})
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
