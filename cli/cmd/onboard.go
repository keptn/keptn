package cmd

import (
	"github.com/spf13/cobra"
)

// onboardCmd represents the onboard command
var onboardCmd = &cobra.Command{
	Use:   "onboard [service]",
	Short: "onboard allows to onbard a new service",
	Long:  `onbaord currently allows to onboard a new service with "onboard service". onboard without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
