package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type authCmdParams struct {
	endPoint     *string
	apiToken     *string
	exportConfig *bool
}

var authParams *authCmdParams
var exportEndPoint url.URL
var exportAPIToken string

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth --endpoint=https://api.keptn.MY.DOMAIN.COM --api-token=SECRET_TOKEN",
	Short: "Authenticates the Keptn CLI against a Keptn installation",
	Long: `Authenticates the Keptn CLI against a Keptn installation using an endpoint and an API token. 
The endpoint and API token are automatically configured during the Keptn installation.
If the authentication is successful, the endpoint and the API token are stored in a password store of the underlying operating system.
More precisely, the Keptn CLI stores the endpoint and API token using *pass* in case of Linux, using *Keychain* in case of macOS, or *Wincred* in case of Windows.

**Note:** If you receive a warning *Using a file-based storage for the key because the password-store seems to be not set up.* this is because a password store could not be found in your environment. In this case, the credentials are stored in *~/.keptn/.keptn* in your home directory.
	`,
	Example:      `keptn auth --endpoint=https://api.keptn.MY.DOMAIN.COM --api-token=abcd-0123-wxyz-7890`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return verifyAuthParams(authParams)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error
		// User wants to print current auth credentials
		if *authParams.exportConfig {
			exportEndPoint, exportAPIToken, err = credentialmanager.NewCredentialManager().GetCreds()
			if err != nil {
				return err
			}
			fmt.Println("Endpoint: ", exportEndPoint.String())
			fmt.Println("API Token: ", exportAPIToken)
			return nil
		}

		logging.PrintLog("Starting to authenticate", logging.InfoLevel)

		url, err := url.Parse(*authParams.endPoint)
		if err != nil {
			logging.PrintLog("Error parsing Keptn API URL", logging.InfoLevel)
			return err
		}

		if url.Path == "" || url.Path == "/" {
			url.Path = "/api"
		}

		authHandler := apiutils.NewAuthenticatedAuthHandler(url.String(), *authParams.apiToken, "x-token", nil, url.Scheme)

		if !mocking {
			authenticated := false

			if !lookupHostname(url.Hostname()) {
				return fmt.Errorf("Authentication was unsuccessful - could not resolve hostname.")
			}

			if endPointErr := checkEndPointStatus(*authParams.endPoint); endPointErr != nil {
				return fmt.Errorf("Authentication was unsucessful: %s"+endPointErrorReasons,
					endPointErr)
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

			logging.PrintLog("Successfully authenticated against the Keptn cluster "+*authParams.endPoint, logging.InfoLevel)
			return credentialmanager.NewCredentialManager().SetCreds(*url, *authParams.apiToken)
		}

		fmt.Println("skipping auth due to mocking flag set to true")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authParams = &authCmdParams{}

	authParams.endPoint = authCmd.Flags().StringP("endpoint", "e", "", "The endpoint exposed by the Keptn installation (e.g., api.keptn.127.0.0.1.xip.io)")
	authParams.apiToken = authCmd.Flags().StringP("api-token", "a", "", "The API token to communicate with the Keptn installation")
	authParams.exportConfig = authCmd.Flags().BoolP("export", "c", false, "To export the current cluster config i.e API token and Endpoint")
}

func verifyAuthParams(authParams *authCmdParams) error {

	if *authParams.exportConfig {
		return nil
	}

	if !mocking {
		if (authParams.endPoint == nil || *authParams.endPoint == "") && (authParams.apiToken == nil || *authParams.apiToken == "") {
			return errors.New("required flag(s) \"api-token\", \"endpoint\" not set")
		} else if authParams.endPoint == nil || *authParams.endPoint == "" {
			return errors.New("required flag \"endpoint\" not set")
		} else if authParams.apiToken == nil || *authParams.apiToken == "" {
			return errors.New("required flag \"api-token\" not set")
		}
	}
	return nil
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

// try to authenticate towards the given endpoint with the provided apiToken
func authenticate(endPoint string, apiToken string) error {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"auth",
		fmt.Sprintf("--endpoint=%s", endPoint),
		fmt.Sprintf("--api-token=%s", apiToken),
	}
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
