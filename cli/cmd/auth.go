package cmd

import (
	"fmt"
	"strings"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/knative/pkg/cloudevents"
	"github.com/spf13/cobra"
)

type authData struct {
}

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
		}

		data := authData{}
		if !strings.HasSuffix(*endPoint, "/") {
			*endPoint += "/"
		}
		authRestEndPoint := *endPoint + "auth"

		err := utils.Send(authRestEndPoint, *apiToken, builder, data)
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
