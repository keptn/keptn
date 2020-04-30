package cmd

import (
	"github.com/spf13/cobra"
)

type generateCmdParams struct {
	Directory *string
}

// deleteCmd represents the delete command
var generateCmd = &cobra.Command{
	Use: "generate [docs | support-archive]",
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
