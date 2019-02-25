package credentialmanager

import (
	"os"

	"github.com/docker/docker-credential-helpers/pass"
	"github.com/keptn/keptn/cli/utils"
	"io/ioutil"
)

// TODO: Write documentation
// For using pass.Pass{} the following commands need to executed:
// 1. sudo apt-get install gpg pass -y
// 2. gpg --generate-key (Use your name and e-mail); Use "find / | xargs file" for generating random bytes; Copy generate pub key
// 3. pass init [generated pub key]

const passwordStoreDirectory = "~/.password-store"
const apiTokenFileName = "api-token.txt"

func SetCreds(secret string) error {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		utils.Warning.Println("Use a file-based storage for the key because the password-store seems to be not set up.")

		return ioutil.WriteFile(apiTokenFileName,, []byte(secret), 0644)
	}
	return setCreds(pass.Pass{}, secret)
}

func GetCreds() (string, error) {
	if _, err := os.Stat(passwordStoreDirectory); os.IsNotExist(err) {
		utils.Warning.Println("Use a file-based storage for the key because the password-store seems to be not set up.")

		dat, err := ioutil.ReadFile(apiTokenFileName)
		return string(dat), err
	}
	return getCreds(pass.Pass{})
}
