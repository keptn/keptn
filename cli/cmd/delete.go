package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [project]",
	Short: "delete is the parent command of \"delete project\"",
	Long:  `delete is the parent command of \"delete project\". \"delete\" without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
