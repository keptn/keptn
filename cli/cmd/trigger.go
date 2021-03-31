package cmd

import "github.com/spf13/cobra"

var triggerCmd = &cobra.Command{
	Use:   "trigger [delivery | evaluation] ",
	Short: "Triggers the execution of an action in keptn",
	Long:  "Triggers the execution of an action in keptn",
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}
