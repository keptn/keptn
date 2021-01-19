package credentialmanager

import (
	"bufio"
	"fmt"
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

// GlobalCheckForContextChange ...Since credential manager is called at multiple times, we dont want to check for context switch for one command at multiple places,
// it should be called only for the first time.
var GlobalCheckForContextChange bool

type keptnConfig struct {
	APIToken       string `yaml:"api_token"`
	Endpoint       string `yaml:"endpoint"`
	ContextName    string `yaml:"name"`
	KeptnNamespace string `yaml:"namespace"`
}

type keptnConfigFile struct {
	Contexts []keptnConfig `yaml:"contexts"`
}

type kubeConfigFileType struct {
	CurrentContext string `yaml:"current-context"`
}

var kubeConfigFile kubeConfigFileType

var keptnContext string

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
	customServerURL := serverURL + "/" + keptnContext + "/" + namespace
	c := &credentials.Credentials{
		ServerURL: customServerURL,
		Username:  url.QueryEscape(endPoint.String()),
		Secret:    apiToken,
	}
	return h.Add(c)
}

func getCreds(h credentials.Helper, namespace string) (url.URL, string, error) {

	if MockAuthCreds {
		return url.URL{}, "", nil
	}

	customServerURL := serverURL + "/" + keptnContext + "/" + namespace
	// Check if creds file is specified in the 'KEPTNCONFIG' environment variable
	if customCredsLocation, ok := os.LookupEnv("KEPTNCONFIG"); ok {
		if customCredsLocation != "" {
			endPoint, apiToken, err := handleCustomCreds(customCredsLocation, namespace)
			// If credential is not found in KEPTNCONFIG, use fallback credential manager
			if apiToken != "" || err != nil {
				return endPoint, apiToken, err
			}
		}
	}
	endPointStr, apiToken, err := h.Get(customServerURL)
	if err != nil {
		return url.URL{}, "", err
	}
	outURL, _ := url.QueryUnescape(endPointStr)
	url, err := url.Parse(outURL)
	return *url, apiToken, err
}

func handleCustomCreds(configLocation string, namespace string) (url.URL, string, error) {
	fileContent, err := file.ReadFile(configLocation)
	if err != nil {
		return url.URL{}, "", err
	}

	var keptnConfig keptnConfigFile

	yaml.Unmarshal([]byte(fileContent), &keptnConfig)
	for _, context := range keptnConfig.Contexts {

		// Keeping default namespace to keptn
		if &context.KeptnNamespace == nil || context.KeptnNamespace == "" {
			context.KeptnNamespace = "keptn"
		}

		if context.ContextName == keptnContext && context.KeptnNamespace == namespace {
			parsedURL, err := url.Parse(context.Endpoint)
			if err != nil {
				return url.URL{}, "", err
			}
			return *parsedURL, context.APIToken, nil
		}
	}
	return url.URL{}, "", nil
}

// initChecks needs to be run when credentialManager is called or initialized
func initChecks(autoApplyNewContext bool) {
	if !GlobalCheckForContextChange {
		cliConfigManager := config.NewCLIConfigManager()
		err := getCurrentContextFromKubeConfig()
		if err != nil {
			log.Fatal(err)
		}
		checkForContextChange(cliConfigManager, autoApplyNewContext)
		GlobalCheckForContextChange = true
	}
}

func getCurrentContextFromKubeConfig() error {
	kubeConfigFile.CurrentContext = ""
	keptnContext = ""
	if MockAuthCreds || MockKubeConfigCheck {
		return nil
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
		fmt.Fprintf(os.Stderr, "Warning: could not open KUBECONFIG file: "+err.Error()+"\n")
		return nil
	}

	err = yaml.Unmarshal([]byte(fileContent), &kubeConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not parse KUBECONFIG file: "+err.Error()+"\n")
		return nil
	}
	return nil
}

func checkForContextChange(cliConfigManager *config.CLIConfigManager, autoApplyNewContext bool) error {

	if MockAuthCreds || MockKubeConfigCheck {
		// Do nothing
		return nil
	}
	cliConfig, err := cliConfigManager.LoadCLIConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Setting keptnContext from ~/.keptn/config file
	keptnContext = cliConfig.CurrentContext
	if kubeConfigFile.CurrentContext != "" && keptnContext != kubeConfigFile.CurrentContext {
		fmt.Printf("Kube context has been changed to %s", kubeConfigFile.CurrentContext)
		fmt.Println()
		if !autoApplyNewContext && keptnContext != "" {
			fmt.Println("Do you want to switch to the new Kube context with the Keptn running there? (y/n)")
			reader := bufio.NewReader(os.Stdin)
			in, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			in = strings.ToLower(strings.TrimSpace(in))
			if !(in == "y" || in == "yes") {
				return nil
			}
		}
		cliConfig.CurrentContext = kubeConfigFile.CurrentContext
		keptnContext = kubeConfigFile.CurrentContext
		err = cliConfigManager.StoreCLIConfig(cliConfig)
		if err != nil {
			return err
		}
	}
	return nil
}
