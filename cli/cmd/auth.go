package cmd

import (
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
	and an api-token. Usage of \"auth\":

keptn auth --endpoint=myendpoint.com --api-token`,
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

		err = utils.Send(req, *apiToken)
		if err != nil {
			utils.Error.Printf("Endpoint invalid or keptn unreachable. Details: %v", err)
			return err
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
	apiToken = authCmd.Flags().StringP("api-token", "a", "", "The api token provided by keptn")
	authCmd.MarkFlagRequired("api-token")
}
