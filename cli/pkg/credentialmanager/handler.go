package credentialmanager

import (
	"fmt"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/keptn/cli/pkg/common"
	"github.com/keptn/keptn/cli/pkg/config"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

var testEndPoint = url.URL{Scheme: "https", Host: "my-endpoint"}

const testAPIToken = "super-secret"
const testNamespace = "keptn-test-namespace"

const credsLab = "keptn"
const serverURL = "https://keptn.sh"
const installCredsKey = "https://keptn-install.sh"

//go:generate moq -pkg credentialmanager_mock -skip-ensure -out ./fake/credential_manager_mock.go . CredentialManagerInterface
type CredentialManagerInterface interface {
	SetCreds(endPoint url.URL, apiToken string, namespace string) error
	GetCreds(namespace string) (url.URL, string, error)
	SetInstallCreds(creds string) error
	GetInstallCreds() (string, error)
	SetCurrentKubeConfig(kubeConfig KubeConfigFileType)
	GetCurrentKubeConfig() KubeConfigFileType
	SetCurrentKeptnCLIConfig(cliConfig config.CLIConfig)
	GetCurrentKeptnCLIConfig() config.CLIConfig
}

var MockAuthCreds bool
var MockKubeConfigCheck bool

// GlobalCheckForContextChange ...Since credential manager is called at multiple times, we don't want to check for context switch for one command at multiple places,
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

type KubeConfigFileType struct {
	CurrentContext string `yaml:"current-context"`
}

var kubeConfigFile KubeConfigFileType

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
	fileContent, err := fileutils.ReadFile(configLocation)
	if err != nil {
		return url.URL{}, "", err
	}

	var keptnConfig keptnConfigFile

	yaml.Unmarshal(fileContent, &keptnConfig)
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
func initChecks(autoApplyNewContext bool, cm CredentialManagerInterface) {
	cliConfigManager := config.NewCLIConfigManager()
	cliConfig, err := cliConfigManager.LoadCLIConfig()
	if err != nil {
		log.Fatal(err)
	}
	if cliConfig.KubeContextCheck && !GlobalCheckForContextChange {
		getCurrentContextFromKubeConfig()
		updatedCLIConfig, kubeConfig, err := checkForContextChange(cliConfigManager, autoApplyNewContext)
		if err != nil {
			log.Fatal(err)
		}
		if updatedCLIConfig != nil {
			cm.SetCurrentKeptnCLIConfig(*updatedCLIConfig)
		}
		if kubeConfig != nil {
			cm.SetCurrentKubeConfig(*kubeConfig)
		}
		GlobalCheckForContextChange = true
	} else {
		cm.SetCurrentKeptnCLIConfig(cliConfig)
		cm.SetCurrentKubeConfig(kubeConfigFile)
	}
}

func getCurrentContextFromKubeConfig() {
	kubeConfigFile.CurrentContext = ""
	keptnContext = ""
	if MockAuthCreds || MockKubeConfigCheck {
		return
	}

	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = keptnutils.ExpandTilde(os.Getenv("KUBECONFIG"))
	} else {
		kubeconfig = filepath.Join(
			keptnutils.UserHomeDir(), ".kube", "config",
		)
	}

	fileContent, err := fileutils.ReadFile(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not parse KUBECONFIG file: "+err.Error()+"\n")
		fmt.Println("Hint: If you don't have a 'kubeconfig' file, you can disable this check via 'keptn set config KubeContextCheck false'")
	} else if err = yaml.Unmarshal(fileContent, &kubeConfigFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not parse KUBECONFIG file: "+err.Error()+"\n")
	}
}

func checkForContextChange(cliConfigManager *config.CLIConfigManager, autoApplyNewContext bool) (*config.CLIConfig, *KubeConfigFileType, error) {
	if MockAuthCreds || MockKubeConfigCheck {
		// Do nothing
		return nil, nil, nil
	}
	cliConfig, err := cliConfigManager.LoadCLIConfig()
	if err != nil {
		return &cliConfig, &kubeConfigFile, err
	}

	if cliConfig.KubeContextCheck {
		// Setting keptnContext from ~/.keptn/config file
		keptnContext = cliConfig.CurrentContext
		if kubeConfigFile.CurrentContext != "" && keptnContext != kubeConfigFile.CurrentContext {
			fmt.Printf("Kube context has been changed to %s\n", kubeConfigFile.CurrentContext)
			if keptnContext != "" {
				userConfirmation := common.NewUserInput().AskBool("Do you want to switch to the new Kube context with the Keptn running there?", &common.UserInputOptions{AssumeYes: autoApplyNewContext})
				fmt.Println("Info: You can turn off the Kube context check by executing: keptn set config KubeContextCheck false")
				if !userConfirmation {
					return &cliConfig, &kubeConfigFile, nil
				}
			}
			cliConfig.CurrentContext = kubeConfigFile.CurrentContext
			keptnContext = kubeConfigFile.CurrentContext
			err = cliConfigManager.StoreCLIConfig(cliConfig)
			if err != nil {
				return &cliConfig, &kubeConfigFile, err
			}
		}
	}
	return &cliConfig, &kubeConfigFile, nil
}
