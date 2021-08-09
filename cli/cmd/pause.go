package cmd

import "github.com/spf13/cobra"

var pauseCmd = &cobra.Command{
	Use:   "pause [ sequence ]",
	Short: "Pauses the execution of a sequence",
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
