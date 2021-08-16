package cmd

import "github.com/spf13/cobra"

var abortCmd = &cobra.Command{
	Use:   "abort [ sequence ]",
	Short: "Aborts the execution of a sequence",
}

func init() {
	rootCmd.AddCommand(abortCmd)
}
