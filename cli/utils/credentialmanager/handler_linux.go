package credentialmanager

import (
	"errors"
	"os"
	"strings"

	"io/ioutil"

	"github.com/docker/docker-credential-helpers/pass"
	"github.com/keptn/keptn/cli/utils"
)

// TODO: Write documentation
// For using pass.Pass{} the following commands need to executed:
// 1. sudo apt-get install gpg pass -y
// 2. gpg --generate-key (Use your name and e-mail); Use "find / | xargs file" for generating random bytes; Copy generate pub key
// 3. pass init [generated pub key]

var passwordStoreDirectory string
var apiTokenURI string

func init() {
	passwordStoreDirectory = os.Getenv("HOME") + "/.password-store"
	apiTokenURI = os.Getenv("HOME") + "/.keptn"
}

// SetCreds stores the credentials consisting of an endpoint and an api token using pass or into a file in case
// pass is unavailable.
func SetCreds(endPoint string, apiToken string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		utils.Warning.Println("Use a file-based storage for the key because the password-store seems to be not set up.")

		return ioutil.WriteFile(apiTokenURI, []byte(endPoint+"\n"+apiToken), 0644)
	}
	return setCreds(pass.Pass{}, endPoint, apiToken)
}

// GetCreds reads the credentials and returns an endpoint, the api token, or potentially an error.
func GetCreds() (string, string, error) {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		utils.Warning.Println("Use a file-based storage for the key because the password-store seems to be not set up.")

		data, err := ioutil.ReadFile(apiTokenURI)
		if err != nil {
			return "", "", err
		}
		creds := strings.Split(string(data), "\n")
		if len(creds) != 2 {
			return "", "", errors.New("Format of file-based key storage is invalid!")
		}
		return creds[0], creds[1], err
	}
	return getCreds(pass.Pass{})
}
