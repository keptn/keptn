package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd implements the create command
var createCmd = &cobra.Command{
	Use:   "create [project | service]",
	Short: `Creates a new project or service`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
