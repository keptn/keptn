package cmd

import (
	"errors"

	"github.com/keptn/keptn/cli/utils"
	"github.com/spf13/cobra"
)

// onboardCmd represents the onboard command
var onboardCmd = &cobra.Command{
	Use:   "onboard [service]",
	Short: "onboard allows to onbard a new service",
	Long:  `onbaord currently allows to onboard a new service with \"onboard service\". onboard without subcommand cannot be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Info.Println("onboard called")
		return errors.New("onboard can only be called in combination with \"service\"")
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
