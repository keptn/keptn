package cmd

import "github.com/spf13/cobra"

// getCmd represents the send command
var getCmd = &cobra.Command{
	Use:   "get [event]",
	Short: `get in combination with the subcommand "event" allows to retrieve a Keptn event`,
	Long:  `get in combination with the subcommand "event" allows to retrieve a Keptn event. Get without subcommand cannot be used.`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
