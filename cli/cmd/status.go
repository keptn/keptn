package cmd

import (
	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"github.com/spf13/cobra"
)

// statusCmdCmd represents the auth command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the status of the CLI",
	Long: `Checks the status of the CLI. This includes a test whether the CLI is authenticated against the Keptn API. 
`,
	Example:      `keptn status`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		authenticator := Authenticator{
			Namespace:        namespace,
			CedentialManager: credentialmanager.NewCredentialManager(assumeYes),
		}
		return authenticator.Auth(AuthenticatorOptions{})
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
