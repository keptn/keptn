package credentialmanager

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/keptn/keptn/cli/pkg/config"
	"github.com/keptn/keptn/cli/pkg/file"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"gopkg.in/yaml.v2"
)

var testEndPoint = url.URL{Scheme: "https", Host: "my-endpoint"}

const testAPIToken = "super-secret"
const testNamespace = "keptn-test-namespace"

const credsLab = "keptn"
const serverURL = "https://keptn.sh"
const installCredsKey = "https://keptn-install.sh"

var MockAuthCreds bool
var MockKubeConfigCheck bool

type credsConfig struct {
	APIToken string `json:"api_token"`
	Endpoint string `json:"endpoint"`
}

type kubeConfigFileType struct {
	CurrentContext string `yaml:"current-context"`
}

var kubeConfigFile kubeConfigFileType

func init() {
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

func setCreds(h credentials.Helper, endPoint url.URL, apiToken string, namespace string) error {
	if MockAuthCreds {
		// Do nothing
		return nil
	}
	customServerURL := serverURL + "/" + kubeConfigFile.CurrentContext + "/" + namespace
	c := &credentials.Credentials{
		ServerURL: customServerURL,
		Username:  endPoint.String(),
		Secret:    apiToken,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper, namespace string) (url.URL, string, error) {

	if MockAuthCreds {
		return url.URL{}, "", nil
	}

	customServerURL := serverURL + "/" + kubeConfigFile.CurrentContext + "/" + namespace
	// Check if creds file is specified in the 'KEPTNCONFIG' environment variable
	if customCredsLocation, ok := os.LookupEnv("KEPTNCONFIG"); ok {
		if customCredsLocation != "" {
			return handleCustomCreds(customCredsLocation)
		}
	}
	endPointStr, apiToken, err := h.Get(customServerURL)
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

// initChecks needs to be run when credentialManager is called or initilized
func initChecks() {
	cliConfigManager := config.NewCLIConfigManager()
	currentContext, err := getCurrentContextFromKubeConfig()
	if err != nil {
		log.Fatal(err)
	}
	checkForContextChange(currentContext, cliConfigManager)
}

func getCurrentContextFromKubeConfig() (string, error) {
	if MockAuthCreds || MockKubeConfigCheck {
		// Do nothing
		kubeConfigFile.CurrentContext = ""
		return "", nil
	}

	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = keptnutils.ExpandTilde(os.Getenv("KUBECONFIG"))
	} else {
		kubeconfig = filepath.Join(
			keptnutils.UserHomeDir(), ".kube", "config",
		)
	}

	fileContent, err := file.ReadFile(kubeconfig)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal([]byte(fileContent), &kubeConfigFile)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return kubeConfigFile.CurrentContext, nil
}

func checkForContextChange(currentContext string, cliConfigManager *config.CLIConfigManager) error {

	if MockAuthCreds || MockKubeConfigCheck {
		// Do nothing
		return nil
	}
	cliConfig, err := cliConfigManager.LoadCLIConfig()
	if err != nil {
		log.Fatal(err)
	}
	if cliConfig.CurrentContext != currentContext {
		fmt.Printf("Kube context has been changed to %s", currentContext)
		fmt.Println()
		fmt.Println("Do you want to continue with this? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.ToLower(strings.TrimSpace(in))
		if !(in == "y" || in == "yes") {
			err := errors.New("stopping installation")
			log.Fatal(err)
		}
		cliConfig.CurrentContext = currentContext
		err = cliConfigManager.StoreCLIConfig(cliConfig)
		if err != nil {
			return err
		}
	}
	return nil
}
