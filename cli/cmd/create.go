package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "create currently allows to create a project",
	Long:  `create currently allows to create a project with "create project". create without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
