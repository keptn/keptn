package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
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

		logging.PrintLog("Starting to configure the GitHub organization, the GitHub user, and the GitHub personal access token", logging.InfoLevel)

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
		configURL.Path = "v1/configure"

		if !mocking {
			logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)
			responseCE, err := utils.Send(configURL, event, apiToken)
			if err != nil {
				logging.PrintLog("Configure was unsuccessful", logging.QuietLevel)
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				logging.PrintLog("Response CE is nil", logging.QuietLevel)
				return nil
			}
			if responseCE.Data != nil {
				return websockethelper.PrintWSContentCEResponse(responseCE, endPoint)
			}
		} else {
			fmt.Println("Skipping configure due to mocking flag set to true")
		}
		// fmt.Println("Successfully configured the GitHub organization, the GitHub user, and the GitHub personal access token")
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

func configure(githubOrg, githubUserName, githubPersonalAccessToken string) error {

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"configure",
		fmt.Sprintf("--org=%s", githubOrg),
		fmt.Sprintf("--user=%s", githubUserName),
		fmt.Sprintf("--token=%s", githubPersonalAccessToken),
	}
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func checkIfConfiguredUsingKube() (bool, error) {

	ops := options{"get",
		"secrets",
		"-n",
		"keptn"}
	ops.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", ops)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "github-credentials"), nil
}
