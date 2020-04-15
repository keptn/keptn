package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd implements the create command
var createCmd = &cobra.Command{
	Use:   "create [project | service]",
}

func init() {
	rootCmd.AddCommand(createCmd)
}
