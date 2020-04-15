package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var generateCmd = &cobra.Command{
	Use:   "generate [docs]",
	Short: `generate is the parent command of "generate docs"`,
	Long:  `generate is the parent command of "generate docs". "generate" without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
