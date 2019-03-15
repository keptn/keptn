package cmd

import (
	"errors"
	"fmt"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/knative/pkg/cloudevents"
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

		builder := cloudevents.Builder{
			Source:    "https://github.com/keptn/keptn/cli#configure",
			EventType: "configure",
			Encoding:  cloudevents.StructuredV01,
		}
		configURL := endPoint
		configURL.Path = "config"

		req, err := builder.Build(configURL.String(), config)
		if err != nil {
			return err
		}

		resp, err := utils.Send(req, apiToken)
		if err != nil {
			fmt.Println("Configure was unsuccessful")
			return err
		}
		if resp.StatusCode != 200 {
			fmt.Println("Configure was unsuccessful")
			return errors.New(resp.Status)
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
