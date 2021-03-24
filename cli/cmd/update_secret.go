package cmd

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type updateSecretCmdParams struct {
	Data  []string
	Scope *string
}

var updateSecretParams *updateSecretCmdParams

var updateSecretCommand = &cobra.Command{
	Use:          `secret SECRET_NAME --from-literal="key1=value1"" --from-literal="key2=value2 --scope=my-scope"`,
	Short:        "Updates an existing secret",
	Example:      `keptn update secret SECRET_NAME --from-literal="key1=value1"" --from-literal="key2=value2 --scope=my-scope"`,
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
		if err := handler.UpdateSecret(args[0], createSecretParams.Data, updateSecretParams.Scope); err != nil {
			return err
		}
		logging.PrintLog(fmt.Sprintf("Secret %s updated successfully", args[0]), logging.InfoLevel)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateSecretCommand)
	updateSecretParams = &updateSecretCmdParams{}
	updateSecretCommand.Flags().StringArrayVar(&updateSecretParams.Data, "from-literal", updateSecretParams.Data, "Specify a key and literal value to insert in secret (i.e. my-key=some-value)")
	updateSecretParams.Scope = updateSecretCommand.Flags().StringP("scope", "s", defaultSecretScope, "The scope of the secret")
}
