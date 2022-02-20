package cmd

import (
	"github.com/keptn/keptn/cli/internal/auth"
	"net/url"
	"strings"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Args:  cobra.NoArgs,
	Short: "Checks the status of the CLI",
	Long: `Checks the status of the CLI. This includes a test whether the CLI is authenticated against the Keptn API. 
`,
	Example:      `keptn status`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		credentialManager := credentialmanager.NewCredentialManager(assumeYes)
		return NewAuthenticator(namespace, credentialManager, auth.NewLocalFileOauthStore()).Auth(AuthenticatorOptions{})
	},
}

func getBridgeURLFromAPIURL(endpointURL url.URL) string {
	return strings.Replace(endpointURL.String(), endpointURL.Path, "/bridge", 1)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
