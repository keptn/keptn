package cmd

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

var endPoint *string
var apiToken *string

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth --endpoint=https://api.keptn.MY.DOMAIN.COM --api-token=SECRET_TOKEN",
	Short: "Authenticates the Keptn CLI against a Keptn installation",
	Long: `Authenticates the Keptn CLI against a Keptn installation using an endpoint
and an API token. The endpoint and API token are exposed during the Keptn installation.
If the authentication is successful, the endpoint and the API token are stored in a password store. 

Example:
	keptn auth --endpoint=https://api.keptn.my.domain.com --api-token=abcd-0123-wxyz-7890`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.PrintLog("Starting to authenticate", logging.InfoLevel)

		url, err := url.Parse(*endPoint)
		if err != nil {
			return err
		}

		authHandler := apiutils.NewAuthenticatedAuthHandler(url.String(), *apiToken, "x-token", nil, "https")

		if !mocking {
			_, err := authHandler.Authenticate()
			if err != nil {
				errMsg := fmt.Sprintf("Authentication was unsuccessful. %s", *err.Message)
				logging.PrintLog(errMsg, logging.QuietLevel)
				return fmt.Errorf("Authentication was unsuccessful. %s", *err.Message)
			}

			logging.PrintLog("Successfully authenticated", logging.InfoLevel)
			return credentialmanager.SetCreds(*url, *apiToken)
		}

		fmt.Println("skipping auth due to mocking flag set to true")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	endPoint = authCmd.Flags().StringP("endpoint", "e", "", "The endpoint exposed by keptn")
	authCmd.MarkFlagRequired("endpoint")
	apiToken = authCmd.Flags().StringP("api-token", "a", "", "The API token provided by keptn")
	authCmd.MarkFlagRequired("api-token")
}

func authUsingKube() error {

	apiToken, err := getAPITokenUsingKube()

	const errorMsg = `Could not retrieve keptn API token: %s
To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/.`

	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	var keptnEndpoint string
	for retries := 0; retries < 15; time.Sleep(5 * time.Second) {

		out, err := getEndpointUsingKube()
		if err != nil || strings.TrimSpace(string(out)) == "" {
			logging.PrintLog("API endpoint not yet available... trying again in 5s", logging.InfoLevel)
		} else {
			keptnEndpoint = "https://api.keptn." + strings.TrimSpace(string(out))
			break
		}
		retries++
	}
	if keptnEndpoint == "" {
		return errors.New("Cannot obtain endpoint of api")
	}
	return authenticate(keptnEndpoint, apiToken)
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
