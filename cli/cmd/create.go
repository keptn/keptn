package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "create is the parent command of \"create project\"",
	Long:  `create is the parent command of \"create project\". \"create\" without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
