package cmd

import (
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:          "configure [monitoring | bridge]",
	Short:        "Configures one of the specified parts of Keptn",
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
