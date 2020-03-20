package credentialmanager

import (
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
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

// MockAuthCreds shows whether the get and set for the auth-creds should be mocked
var MockAuthCreds bool

var apiTokenFileURI string

var credsFileURI string

func init() {
	dir, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		log.Fatal(err)
	}

	apiTokenFileURI = dir + ".keptn"

	credsFileURI = dir + ".keptn-creds"

	credentials.SetCredsLabel(credsLab)
}

func setInstallCreds(h credentials.Helper, creds string) error {
	c := &credentials.Credentials{
		ServerURL: installCredsKey,
		Username:  "creds",
		Secret:    creds,
	}
	return h.Add(c)
}

func getInstallCreds(h credentials.Helper) (string, error) {
	_, creds, err := h.Get(installCredsKey)
	if err != nil {
		return "", err
	}
	return creds, err
}

func setCreds(h credentials.Helper, endPoint url.URL, apiToken string) error {
	if MockAuthCreds {
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

	if MockAuthCreds {
		return url.URL{}, "", nil
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
	data, err = ioutil.ReadFile(apiTokenFileURI)
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
	data, err = ioutil.ReadFile(credsFileURI)
	if err != nil {
		return "", err
	}
	dataStr := strings.TrimSpace(strings.Replace(string(data), "\r\n", "\n", -1))
	return dataStr, err
}
