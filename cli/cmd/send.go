package cmd

import "github.com/spf13/cobra"

// sendCmd implements the send command
var sendCmd = &cobra.Command{
	Use:   "send [event]",
	Short: `"send" can be used with the subcommand "event"`,
	Long:  `"send" can be used with the subcommand "event"`,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
