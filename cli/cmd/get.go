package cmd

import "github.com/spf13/cobra"

// getCmd represents the send command
var getCmd = &cobra.Command{
	Use:   "get [event | project | projects | stage | stages | service | services]",
	Short: `Display an event or one or many keptn entities such as project, stage, or service.`,
	Long:  `Display an event or one or many Keptn entities such as project, stage, or service.`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
