package cmd

import (
	"fmt"

	"github.com/keptn/keptn/cli/utils/credentialmanager"

	"github.com/spf13/cobra"
)

// statusCmdCmd represents the auth command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the status of the CLI",
	Long: `Checks the status of the CLI. This includes a test whether the CLI is authenticated against the Keptn API. 

Example:
	keptn status`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil || endPoint.String() == "" || apiToken == "" {
			fmt.Println("CLI is not authenticated against any Keptn cluster.  For authenticating your CLI use \"keptn auth\"")
			return nil
		}

		err = authenticate(endPoint.String(), apiToken)
		if err != nil {
			fmt.Printf("CLI cannot be authenticated against Keptn cluster %s\n", endPoint.String())
			return err
		}
		fmt.Printf("CLI is authenticated against the Keptn cluster %s\n", endPoint.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
