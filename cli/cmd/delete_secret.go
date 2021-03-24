package cmd

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type deleteSecretCmdParams struct {
	Scope *string
}

var deleteSecretParams *deleteSecretCmdParams

var deleteSecretCommand = &cobra.Command{
	Use:          `secret SECRET_NAME --scope=my-scope"`,
	Short:        "Deletes a secret from the given scope",
	Example:      `keptn delete secret SECRET_NAME --scope=my-scope"`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SECRET_NAME not set")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := NewSecretCmdHandler(credentialmanager.NewCredentialManager(assumeYes))
		if err != nil {
			return nil
		}
		if err := handler.DeleteSecret(args[0], deleteSecretParams.Scope); err != nil {
			return err
		}
		logging.PrintLog(fmt.Sprintf("Secret %s has been deleted", args[0]), logging.InfoLevel)
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteSecretCommand)
	deleteSecretParams = &deleteSecretCmdParams{}
	deleteSecretParams.Scope = deleteSecretCommand.Flags().StringP("scope", "s", defaultSecretScope, "The scope of the secret")
}
