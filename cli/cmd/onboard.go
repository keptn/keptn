package cmd

import (
	"github.com/spf13/cobra"
)

// onboardCmd represents the onboard command
var onboardCmd = &cobra.Command{
	Use:   "onboard [service]",
	Short: "onboard allows to create a new service",
	Long:  `onboard currently allows to create a new service with "onboard service". onboard without sub-command cannot be used.`,
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
