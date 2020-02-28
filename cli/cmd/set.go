package cmd

import "github.com/spf13/cobra"

// setCmd implements the command set
var setCmd = &cobra.Command{
	Use:   "set [config]",
	Short: `"set" can be used with the subcommand "config"`,
	Long:  `"set" can be used with the subcommand "config"`,
}

func init() {
	rootCmd.AddCommand(setCmd)
}
