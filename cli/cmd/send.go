package cmd

import "github.com/spf13/cobra"

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [event]",
	Short: `send in combination with the subcommand "event" allows to send a keptn event`,
	Long:  `send in combination with the subcommand "event" allows to send a keptn event. Send without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
