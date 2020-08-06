package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd implements the create command
var updateCmd = &cobra.Command{
	Use:   "update [project]",
	Short: "Updates an existing Keptn project",
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
