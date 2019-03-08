package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/knative/pkg/cloudevents"
	"github.com/spf13/cobra"
)

var endPoint *string
var apiToken *string

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticates the keptn CLI against a keptn installation.",
	Long: `Authenticates the keptn CLI against a keptn installation using an endpoint
and an API token. Both, the endpoint and API token are exposed during the keptn installation.
If the authentication is successful, the endpoint and the API token are stored in a password store. 

Example:
	keptn auth --endpoint=myendpoint.com --api-token=xyz`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to authenticate")
		builder := cloudevents.Builder{
			Source:    "https://github.com/keptn/keptn/cli#auth",
			EventType: "auth",
			Encoding:  cloudevents.StructuredV01,
		}

		if !strings.HasSuffix(*endPoint, "/") {
			*endPoint += "/"
		}

		var data interface{}
		req, err := builder.Build(*endPoint+"auth", data)
		if err != nil {
			return err
		}

		resp, err := utils.Send(req, *apiToken)
		if err != nil {
			fmt.Println("Authentication was unsuccessful")
			return err
		}
		if resp.StatusCode == 404 {
			return errors.New("Endpoint not found")
		}
		if resp.StatusCode == 401 {
			return errors.New("Unauthorized request")
		}
		if resp.StatusCode != 200 {
			fmt.Println("Authentication was unsuccessful")
			return errors.New(resp.Status)
		}

		// Store endpoint and api token as credentials
		fmt.Println("Successfully authenticated")
		return credentialmanager.SetCreds(*endPoint, *apiToken)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	endPoint = authCmd.Flags().StringP("endpoint", "e", "", "The endpoint exposed by keptn")
	authCmd.MarkFlagRequired("endpoint")
	apiToken = authCmd.Flags().StringP("api-token", "a", "", "The API token provided by keptn")
	authCmd.MarkFlagRequired("api-token")
}
