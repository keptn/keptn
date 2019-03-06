package cmd

import (
	"errors"

	"github.com/keptn/keptn/cli/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "create currently allows to create a project",
	Long:  `create currently allows to create a project with \"create project\". create without subcommand cannot be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Info.Println("create called")
		return errors.New("create can only be called in combination with \"project\"")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
