package cmd

import "github.com/spf13/cobra"

// sendCmd implements the send command
var sendCmd = &cobra.Command{
	Use:   "send [event]",
	Short: "Sends an event to Keptn",
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
