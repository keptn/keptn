package cmd

import "github.com/spf13/cobra"

// setCmd implements the command set
var setCmd = &cobra.Command{
	Use:   "set [config]",
	Short: `Sets flags of the CLI configuration`,
	Long:  `Sets flags of the CLI configuration.`,
}

func init() {
	rootCmd.AddCommand(setCmd)
}
