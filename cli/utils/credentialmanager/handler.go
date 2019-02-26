package credentialmanager

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
)

const credsLab = "keptn"
const serverURL = "https://keptn.sh"

type bot interface {
	SetCreds(endPoint string, apiToken string) error
	GetCreds() (string, string, error)
}

// MockCreds shows whether the get and set should be mocked by a file
// named "endPoint.txt"
var MockCreds bool

func setCreds(h credentials.Helper, endPoint string, apiToken string) error {
	if MockCreds {
		// Do nothing
		return nil
	}

	credentials.SetCredsLabel(credsLab)
	c := &credentials.Credentials{
		ServerURL: serverURL,
		Username:  endPoint,
		Secret:    apiToken,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper) (string, string, error) {

	if MockCreds {
		return ReadCredsFromFile()
	}
	return h.Get(serverURL)
}

// ReadCredsFromFile reads the credentials from a file named "endPoint.txt".
// This function is used for testing
func ReadCredsFromFile() (string, string, error) {
	const endPointFile = "endPoint.txt"
	if _, err := os.Stat(endPointFile); os.IsNotExist(err) {
		return "", "", err
	}

	data, err := ioutil.ReadFile(endPointFile)
	if err != nil {
		return "", "", err
	}
	creds := strings.Split(string(data), "\n")
	if len(creds) != 2 {
		return "", "", errors.New("Format of endPoint.txt file is invalid")
	}
	if !strings.HasSuffix(creds[0], "/") {
		creds[0] += "/"
	}
	return creds[0], creds[1], nil
}
