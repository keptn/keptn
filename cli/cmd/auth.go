package cmd

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

var endPoint *string
var apiToken *string

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth --endpoint=https://api.keptn.MY.DOMAIN.COM --api-token=SECRET_TOKEN",
	Short: "Authenticates the keptn CLI against a keptn installation.",
	Long: `Authenticates the keptn CLI against a keptn installation using an endpoint
and an API token. Both, the endpoint and API token are exposed during the keptn installation.
If the authentication is successful, the endpoint and the API token are stored in a password store. 

Example:
	keptn auth --endpoint=https://api.keptn.my.domain.com --api-token=xyz0123`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.PrintLog("Starting to authenticate", logging.InfoLevel)

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#auth")
		contentType := "application/json"
		var data interface{}
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        "auth",
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: data,
		}

		u, err := url.Parse(*endPoint)
		if err != nil {
			return err
		}

		authURL := *u
		authURL.Path = "v1/auth"

		if !mocking {
			_, err = utils.Send(authURL, event, *apiToken)
			if err != nil {
				logging.PrintLog("Authentication was unsuccessful", logging.QuietLevel)
				return err
			}
			logging.PrintLog("Successfully authenticated", logging.InfoLevel)
			return credentialmanager.SetCreds(*u, *apiToken)
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
			keptnEndpoint = "https://" + strings.TrimSpace(string(out))
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
		"virtualservice",
		"api",
		"-n",
		"keptn",
		"-ojsonpath={.spec.hosts[0]}"}
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
