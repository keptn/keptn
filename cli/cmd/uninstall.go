//go:build !nokubectl
// +build !nokubectl

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Deprecated:   fmt.Sprintf(MsgDeprecatedUseHelm, Version),
	Use:          "uninstall",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("this command is deprecated")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "")
	installCmd.PersistentFlags().MarkHidden("insecure-skip-tls-verify")
}
