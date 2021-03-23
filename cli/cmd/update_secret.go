package cmd

import (
	"errors"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

type updateSecretCmdParams struct {
	Data  []string
	Scope *string
}

var updateSecretParams *updateSecretCmdParams

var updateSecretCommand = &cobra.Command{
	Use:          `secret SECRETNAME --from-literal="key1=value1"" --from-literal="key2=value2 --scope=my-scope"`,
	Short:        "Updates an existing secret",
	Example:      `keptn update secret SECRETNAME --from-literal="key1=value1"" --from-literal="key2=value2 --scope=my-scope"`,
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
		return handler.UpdateSecret(args[0], createSecretParams.Data, updateSecretParams.Scope)
	},
}

func init() {
	updateCmd.AddCommand(updateSecretCommand)
	updateSecretParams = &updateSecretCmdParams{}
	updateSecretCommand.Flags().StringArrayVar(&updateSecretParams.Data, "from-literal", updateSecretParams.Data, "Specify a key and literal value to insert in secret (i.e. mykey=somevalue)")
	updateSecretParams.Scope = updateSecretCommand.Flags().StringP("scope", "s", defaultSecretScope, "The scope of the secret")
}
