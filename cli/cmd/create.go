package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd implements the create command
var createCmd = &cobra.Command{
	Use:   "create [project,service]",
	Short: `"create" can be used with the subcommand "project" or "service"`,
	Long:  `"create" can be used with the subcommand "project" or "service"`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
