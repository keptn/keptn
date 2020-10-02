package credentialmanager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/docker/docker-credential-helpers/credentials"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

var testEndPoint = url.URL{Scheme: "https", Host: "my-endpoint"}

const testAPIToken = "super-secret"

const credsLab = "keptn"
const serverURL = "https://keptn.sh"
const installCredsKey = "https://keptn-install.sh"

var MockAuthCreds bool

type credsConfig struct {
	APIToken string `json:"api_token"`
	Endpoint string `json:"endpoint"`
}

func init() {

	_, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		log.Fatal(err)
	}
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

	if customCredsLocation, ok := os.LookupEnv("KEPTNCONFIG"); ok {
		return handleCustomCreds(customCredsLocation)
	}

	endPointStr, apiToken, err := h.Get(serverURL)
	if err != nil {
		return url.URL{}, "", err
	}
	url, err := url.Parse(endPointStr)
	return *url, apiToken, err
}

func handleCustomCreds(configLocation string) (url.URL, string, error) {
	fd, err := os.Open(configLocation)
	if err != nil {
		return url.URL{}, "", err
	}

	defer fd.Close()

	byteValue, _ := ioutil.ReadAll(fd)

	var credsConfig credsConfig

	json.Unmarshal(byteValue, &credsConfig)

	parsedURL, err := url.Parse(credsConfig.Endpoint)
	if err != nil {
		return url.URL{}, "", err
	}

	return *parsedURL, credsConfig.APIToken, nil
}
