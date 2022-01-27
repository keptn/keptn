package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd implements the create command
var createCmd = &cobra.Command{
	Use:   "create [project | service | secret]",
	Short: `Creates a new project, service or secret`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
