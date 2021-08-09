package cmd

import (
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the status of the CLI",
	Long: `Checks the status of the CLI. This includes a test whether the CLI is authenticated against the Keptn API. 
`,
	Example:      `keptn status`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		credentialManager := credentialmanager.NewCredentialManager(assumeYes)
		authenticator := NewAuthenticator(namespace, credentialManager)

		err := authenticator.Auth(AuthenticatorOptions{})
		if err != nil {
			return err
		}

		endpoint, _, err := credentialManager.GetCreds(namespace)
		if err != nil {
			return err
		}
		fmt.Println("Bridge URL: " + getBridgeURLFromAPIURL(endpoint))
		return nil
	},
}

func getBridgeURLFromAPIURL(endpointURL url.URL) string {
	return strings.Replace(endpointURL.String(), endpointURL.Path, "/bridge", 1)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
