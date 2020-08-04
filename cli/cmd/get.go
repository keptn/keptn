package cmd

import "github.com/spf13/cobra"

// getCmd represents the send command
var getCmd = &cobra.Command{
	Use:   "get [event | project | projects | stage | stages | service | services]",
	Short: "Displays an event or Keptn entities such as project, stage, or service",
	Long:  `Displays an event or Keptn entities such as project, stage, or service.`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
