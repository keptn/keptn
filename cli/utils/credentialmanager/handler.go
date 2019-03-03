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

var apiTokenFileURI string
var mockAPItokenFileURI string

func init() {
	apiTokenFileURI = os.Getenv("HOME") + "/.keptn"
	mockAPItokenFileURI = os.Getenv("HOME") + "/.keptnmock"
}

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
		return readCredsFromFile()
	}
	return h.Get(serverURL)
}

// readCredsFromFile reads the credentials from a file named "endPoint.txt".
// This function is used for testing
func readCredsFromFile() (string, string, error) {
	var data []byte
	var err error
	if MockCreds {
		data, err = ioutil.ReadFile(mockAPItokenFileURI)
	} else {
		data, err = ioutil.ReadFile(apiTokenFileURI)
	}
	if err != nil {
		return "", "", err
	}
	creds := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(creds) != 2 {
		return "", "", errors.New("Format of file-based key storage is invalid")
	}
	if !strings.HasSuffix(creds[0], "/") {
		creds[0] += "/"
	}
	return creds[0], creds[1], err
}
