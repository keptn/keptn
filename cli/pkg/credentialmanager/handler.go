package credentialmanager

import (
	"log"
	"net/url"

	"github.com/docker/docker-credential-helpers/credentials"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

var testEndPoint = url.URL{Scheme: "https", Host: "my-endpoint"}

const testAPIToken = "super-secret"

const credsLab = "keptn"
const serverURL = "https://keptn.sh"
const installCredsKey = "https://keptn-install.sh"

var MockAuthCreds bool

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
	endPointStr, apiToken, err := h.Get(serverURL)
	if err != nil {
		return url.URL{}, "", err
	}
	url, err := url.Parse(endPointStr)
	return *url, apiToken, err
}
