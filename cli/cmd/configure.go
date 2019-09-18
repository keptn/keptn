package cmd

import (
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:          "configure [domain | monitoring]",
	Short:        "configure is the parent command for \"configure domain\" and \"configure monitoring\"",
	Long:         "configure is the parent command for \"configure domain\" and \"configure monitoring\". \"configure\" without subcommand cannot be used.",
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
