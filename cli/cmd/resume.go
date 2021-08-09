package cmd

import "github.com/spf13/cobra"

var resumeCmd = &cobra.Command{
	Use:   "resume [ sequence ]",
	Short: "Resumes the execution of a sequence",
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
