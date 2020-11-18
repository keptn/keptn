package cmd

import (
	"github.com/spf13/cobra"
)

type upgradeProjectCmdParams struct {
	Shipyard    *bool
	FromVersion *string
	ToVersion   *string
}

var upgradeProjectParams *upgradeProjectCmdParams

// upgradeProjectCmd represents the project command
var upgradeProjectCmd = &cobra.Command{
	Use:   "upgrade project PROJECTNAME --shipyard --fromVersion=CURRENT_SHIPYARD_VERSION --toVersion=TARGET_SHIPYARD_VERSION",
	Short: "Upgrades an existing Keptn project",
	Long: `Upgrades an existing Keptn project with the provided name. 

This command will upgrade the shipyard of the project to the specified version

By executing the update project command, Keptn will fetch the current shipyard.yaml file of the project and convert it to the version specified in the 'toVersion'' flag.

For more information about upgrading projects, go to [Manage Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/upgrade)
`,
	Example:      `keptn upgrade project PROJECTNAME --shipyard --fromVersion=0.1.0 --toVersion=0.2.0`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeProjectCmd)

	upgradeProjectParams = &upgradeProjectCmdParams{}

	upgradeProjectParams.Shipyard = upgradeProjectCmd.Flags().BoolP("shipyard", "", false, "Upgrade the shipyard file of the project")
	upgradeProjectParams.FromVersion = upgradeProjectCmd.Flags().StringP("fromVersion", "", "", "The current version of the shipyard")
	upgradeProjectParams.ToVersion = upgradeProjectCmd.Flags().StringP("toVersion", "", "", "The new target version of the shipyard")
	upgradeProjectCmd.MarkFlagRequired("shipyard")
	upgradeProjectCmd.MarkFlagRequired("fromVersion")
	upgradeProjectCmd.MarkFlagRequired("toVersion")
}
