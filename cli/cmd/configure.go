package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type configData struct {
	Org   *string `json:"org"`
	User  *string `json:"user"`
	Token *string `json:"token"`
}

var config configData

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configures the GitHub organization, the GitHub user, and the GitHub personal access token belonging to that user in the keptn installation.",
	Long: `Configures the GitHub organization, the GitHub user, and the GitHub personal access token belonging to that user in the keptn installation.

Example:
	keptn configure --org=MyOrg --user=keptnUser --token=XYZ`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		fmt.Println("Starting to configure the GitHub organization, the GitHub user, and the GitHub personal access token")

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configure")

		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        "configure",
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: config,
		}

		configURL := endPoint
		configURL.Path = "config"

		fmt.Println("Connecting to server ", endPoint.String())
		_, err = utils.Send(configURL, event, apiToken)
		if err != nil {
			fmt.Println("Configure was unsuccessful")
			return err
		}

		fmt.Println("Successfully configured the GitHub organization, the GitHub user, and the GitHub personal access token")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	config.Org = configureCmd.Flags().StringP("org", "o", "", "The GitHub organization")
	configureCmd.MarkFlagRequired("org")
	config.User = configureCmd.Flags().StringP("user", "u", "", "The GitHub user")
	configureCmd.MarkFlagRequired("user")
	config.Token = configureCmd.Flags().StringP("token", "t", "", "The GitHub personal access token")
	configureCmd.MarkFlagRequired("token")
}
