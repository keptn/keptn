package cmd

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
)

var endPoint *string
var apiToken *string

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth --endpoint=https://api.keptn.MY.DOMAIN.COM --api-token=SECRET_TOKEN",
	Short: "Authenticates the Keptn CLI against a Keptn installation",
	Long: `Authenticates the Keptn CLI against a Keptn installation using an endpoint
and an API token. The endpoint and API token are automatically configured during the Keptn installation.
If the authentication is successful, the endpoint and the API token are stored in a password store of the underlying operating system.
More precisely, the keptn CLI stores the endpoint and API token using *pass* in case of Linux, using *Keychain* in case of macOS, or *Wincred* in case of Windows.

**Note**: If you receive a warning *Using a file-based storage for the key because the password-store seems to be not set up.* this is because a password store could not be found in your environment. In this case, the credentials are stored in *~/.keptn/.keptn* in your home directory.
	`,
	Example:      `keptn auth --endpoint=https://api.keptn.MY.DOMAIN.COM --api-token=abcd-0123-wxyz-7890`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.PrintLog("Starting to authenticate", logging.InfoLevel)

		url, err := url.Parse(*endPoint)
		if err != nil {
			logging.PrintLog("Error parsing Keptn API URL", logging.InfoLevel)
			return err
		}

		authHandler := apiutils.NewAuthenticatedAuthHandler(url.String(), *apiToken, "x-token", nil, *scheme)

		if !mocking {
			authenticated := false

			if !lookupHostname(url.Hostname()) {
				return fmt.Errorf("Authentication was unsuccessful - could not resolve hostname.")
			}

			// try to authenticate (and retry it)
			for retries := 0; retries < 3; time.Sleep(5 * time.Second) {
				_, err := authHandler.Authenticate()
				if err != nil {
					errMsg := fmt.Sprintf("Authentication was unsuccessful. %s", *err.Message)
					logging.PrintLog(errMsg, logging.QuietLevel)
					logging.PrintLog("Retrying...", logging.InfoLevel)
					retries++
				} else {
					authenticated = true
					break
				}
			}

			if !authenticated {
				return fmt.Errorf("Authentication was unsuccessful - could not authenticate against the server.")
			}

			logging.PrintLog("Successfully authenticated", logging.InfoLevel)
			return credentialmanager.NewCredentialManager().SetCreds(*url, *apiToken)
		}

		fmt.Println("skipping auth due to mocking flag set to true")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	endPoint = authCmd.Flags().StringP("endpoint", "e", "", "The endpoint exposed by the Keptn installation (e.g., api.keptn.127.0.0.1.xip.io)")
	authCmd.MarkFlagRequired("endpoint")
	apiToken = authCmd.Flags().StringP("api-token", "a", "", "The API token to communicate with the Keptn installation")
	authCmd.MarkFlagRequired("api-token")
}

func lookupHostname(hostname string) bool {
	if strings.HasSuffix(hostname, "xip.io") {
		logging.PrintLog("Skipping lookup of xip.io domain", logging.InfoLevel)
		return true
	} else {
		// first, try to resolve the domain (and retry it)
		for retries := 0; retries < 3; time.Sleep(5 * time.Second) {
			_, err := net.LookupHost(hostname)
			if err != nil {
				logging.PrintLog("Failed to resolve hostname "+hostname, logging.InfoLevel)
				logging.PrintLog("Retrying...", logging.InfoLevel)
				retries++
			} else {
				return true
			}
		}
	}

	return false
}

// authenticate using secrets obtained via kubectl
func authUsingKube(apiSvcType apiServiceType) error {

	// get api token
	apiToken, err := getAPITokenUsingKube()

	const errorMsg = `Could not retrieve keptn API token: %s
To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/.`

	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	// try to obtain endpoint using kubectl (retry a couple of times in case it is not yet available)
	var keptnEndpoint string
	for retries := 0; retries < 15; time.Sleep(5 * time.Second) {

		var out string
		var err error
		switch apiSvcType {
		case LoadBalancer:
			out, err = getLoadBalancerEndpoint()
			break
		case NodePort:
			out, err = getNodePortEndpoint()
			break
		case ClusterIP:
			logging.PrintLog("Cannot reach API service of type 'ClusterIP'", logging.InfoLevel)
			return errors.New("cannot reach API service of type 'ClusterIP'")
			break
		default:
			break
		}
		if err != nil || strings.TrimSpace(out) == "" {
			logging.PrintLog("API endpoint not yet available... trying again in 5s", logging.InfoLevel)
		} else {
			logging.PrintLog("Received Keptn Domain: "+out, logging.InfoLevel)
			keptnEndpoint = "http://" + strings.TrimSpace(out)
			break
		}
		retries++
	}
	// fail if we did not receive a correct endpoint
	if keptnEndpoint == "" {
		return errors.New("Cannot obtain endpoint of api")
	}

	return authenticate(keptnEndpoint, apiToken)
}

func getNodePortEndpoint() (string, error) {
	ops := options{"get",
		"nodes",
		"-ojsonpath={ $.items[0].status.addresses[?(@.type==\"ExternalIP\")].address }"}
	ops.appendIfNotEmpty(kubectlOptions)
	return keptnutils.ExecuteCommand("kubectl", ops)
}

func getLoadBalancerEndpoint() (string, error) {
	ops := options{"get",
		"svc",
		"api-gateway-nginx",
		"-n",
		"keptn",
		"-ojsonpath={.status.loadBalancer.ingress[0].ip}"}
	ops.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", ops)
	if err != nil {
		return "", err
	}

	if out != "" {
		return out, nil
	}

	ops2 := options{"get",
		"svc",
		"api-gateway-nginx",
		"-n",
		"keptn",
		"-ojsonpath={.status.loadBalancer.ingress[0].hostname}"}
	ops2.appendIfNotEmpty(kubectlOptions)
	return keptnutils.ExecuteCommand("kubectl", ops2)
}

func getEndpointUsingKube() (string, error) {
	ops := options{"get",
		"cm",
		"keptn-domain",
		"-n",
		"keptn",
		"-ojsonpath={.data.app_domain}"}
	ops.appendIfNotEmpty(kubectlOptions)
	return keptnutils.ExecuteCommand("kubectl", ops)
}

func getAPITokenUsingKube() (string, error) {
	ops := options{"get",
		"secret",
		"keptn-api-token",
		"-n",
		"keptn",
		"-ojsonpath={.data.keptn-api-token}"}
	ops.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", ops)
	if err != nil {
		return out, err
	}
	apiToken, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		return "", err
	}
	return string(apiToken), nil
}

// try to authenticate towards the given endpoint with the provided apiToken
func authenticate(endPoint string, apiToken string) error {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"auth",
		"--scheme=http",
		fmt.Sprintf("--endpoint=%s", endPoint),
		fmt.Sprintf("--api-token=%s", apiToken),
	}
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
