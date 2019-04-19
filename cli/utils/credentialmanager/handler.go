package credentialmanager

import (
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
)

var testEndPoint = url.URL{Scheme: "https", Host: "my-endpoint"}

const testAPIToken = "super-secret"

const credsLab = "keptn"
const serverURL = "https://keptn.sh"
const installCredsKey = "https://keptn-install.sh"

type bot interface {
	SetInstallCreds(creds string) error
	GetInstallCreds() (string, error)
	SetCreds(endPoint string, apiToken string) error
	GetCreds() (string, string, error)
}

// MockCreds shows whether the get and set should be mocked by a file
// named "endPoint.txt"
var MockCreds bool

var apiTokenFileURI string
var mockAPItokenFileURI string

var credsFileURI string
var mockCredsFileURI string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	apiTokenFileURI = usr.HomeDir + string(os.PathSeparator) + ".keptn"
	mockAPItokenFileURI = usr.HomeDir + string(os.PathSeparator) + ".keptnmock"

	credsFileURI = usr.HomeDir + string(os.PathSeparator) + ".keptn-creds"
	mockCredsFileURI = usr.HomeDir + string(os.PathSeparator) + ".keptn-credsmock"

	credentials.SetCredsLabel(credsLab)
}

func setInstallCreds(h credentials.Helper, creds string) error {
	if MockCreds {
		// Do nothing
		return nil
	}

	c := &credentials.Credentials{
		ServerURL: installCredsKey,
		Username:  "creds",
		Secret:    creds,
	}
	return h.Add(c)
}

func getInstallCreds(h credentials.Helper) (string, error) {
	if MockCreds {
		// Do nothing
		return "", nil
	}
	_, creds, err := h.Get(installCredsKey)
	if err != nil {
		return "", err
	}
	return creds, err
}

func setCreds(h credentials.Helper, endPoint url.URL, apiToken string) error {
	if MockCreds {
		// Do nothing
		return nil
	}

	c := &credentials.Credentials{
		ServerURL: serverURL,
		Username:  endPoint.String(),
		Secret:    apiToken,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper) (url.URL, string, error) {

	if MockCreds {
		return readCredsFromFile()
	}
	endPointStr, apiToken, err := h.Get(serverURL)
	if err != nil {
		return url.URL{}, "", err
	}
	url, err := url.Parse(endPointStr)
	return *url, apiToken, err
}

// readCredsFromFile reads the credentials from a file named "endPoint.txt".
// This function is used for testing
func readCredsFromFile() (url.URL, string, error) {
	var data []byte
	var err error
	if MockCreds {
		data, err = ioutil.ReadFile(mockAPItokenFileURI)
	} else {
		data, err = ioutil.ReadFile(apiTokenFileURI)
	}
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

// readInstallCredsFromFile reads the credentials from a file named "creds.json".
// This function is used for testing
func readInstallCredsFromFile() (string, error) {
	var data []byte
	var err error
	if MockCreds {
		data, err = ioutil.ReadFile(mockCredsFileURI)
	} else {
		data, err = ioutil.ReadFile(credsFileURI)
	}
	if err != nil {
		return "", err
	}
	dataStr := strings.TrimSpace(strings.Replace(string(data), "\r\n", "\n", -1))
	return dataStr, err
}
