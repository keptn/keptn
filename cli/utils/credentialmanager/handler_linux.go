package credentialmanager

import (
	"fmt"
	"net/url"
	"os"

	"io/ioutil"

	"github.com/docker/docker-credential-helpers/pass"
)

// TODO: Write documentation
// For using pass.Pass{} the following commands need to executed:
// 1. sudo apt-get install gpg pass -y
// 2. gpg --generate-key (Use your name and e-mail); Use "find / | xargs file" for generating random bytes; Copy generate pub key
// 3. pass init [generated pub key]

var passwordStoreDirectory string

func init() {
	passwordStoreDirectory = os.Getenv("HOME") + "/.password-store"
}

// SetCreds stores the credentials consisting of an endpoint and an api token using pass or into a file in case
// pass is unavailable.
func SetCreds(endPoint url.URL, apiToken string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		fmt.Println("Using a file-based storage for the key because the password-store seems to be not set up.")

		return ioutil.WriteFile(apiTokenFileURI, []byte(endPoint.String()+"\n"+apiToken), 0644)
	}
	return setCreds(pass.Pass{}, endPoint, apiToken)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func GetCreds() (url.URL, string, error) {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		return readCredsFromFile()
	}
	return getCreds(pass.Pass{})
}

// SetInstallCreds sets the install credentials
func SetInstallCreds(creds string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		fmt.Println("Using a file-based storage for the key because the password-store seems to be not set up.")

		return ioutil.WriteFile(credsFileURI, []byte(creds), 0644)
	}
	return setInstallCreds(pass.Pass{}, creds)
}

// GetInstallCreds gets the install credentials
func GetInstallCreds() (string, error) {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		return readInstallCredsFromFile()
	}
	return getInstallCreds(pass.Pass{})
}
