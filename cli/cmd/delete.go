package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [project]",
	Short: "Deletes a project",
	Long:  `Deletes a project`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
