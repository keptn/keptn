package cmd

import (
	"errors"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

type deleteSecretCmdParams struct {
	Scope *string
}

var deleteSecretParams *deleteSecretCmdParams

var deleteSecretCommand = &cobra.Command{
	Use:          `secret SECRETNAME --scope=my-scope"`,
	Short:        "Deletes a secret from the given scope",
	Example:      `keptn delete secret SECRETNAME --scope=my-scope"`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SECRETNAME not set")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := NewSecretCmdHandler(credentialmanager.NewCredentialManager(assumeYes))
		if err != nil {
			return nil
		}
		return handler.DeleteSecret(args[0], deleteSecretParams.Scope)
	},
}

func init() {
	deleteCmd.AddCommand(deleteSecretCommand)
	deleteSecretParams = &deleteSecretCmdParams{}
	deleteSecretParams.Scope = deleteSecretCommand.Flags().StringP("scope", "s", defaultSecretScope, "The scope of the secret")
}
