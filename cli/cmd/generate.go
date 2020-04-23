package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var generateCmd = &cobra.Command{
	Use:   "generate [docs]",
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
