package cmd

import (
	"github.com/spf13/cobra"
)

// onboardCmd represents the onboard command
var onboardCmd = &cobra.Command{
	Use:   "onboard [service]",
	Short: "Creates a new service and uploads its Helm chart to the branches in the Git repository",
	Long:  `Creates a new service and uploads its Helm chart to the branches in the Git repository.`,
	Deprecated: "please use \"create\" instead.",
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
