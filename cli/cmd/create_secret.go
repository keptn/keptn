package cmd

import (
	"errors"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

type createSecretCmdParams struct {
	Data  []string
	Scope *string
}

var createSecretParams *createSecretCmdParams

var createSecretCommand = &cobra.Command{
	Use:          `secret SECRET_NAME --from-literal="key1=value1"" --from-literal="key2=value2" --scope=my-scope`,
	Short:        "Creates a new secret",
	Example:      `keptn create secret SECRET_NAME --from-literal="key1=value1"" --from-literal="key2=value2" --scope=my-scope`,
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
		return handler.CreateSecret(args[0], createSecretParams.Data, createSecretParams.Scope)
	},
}

func init() {
	createCmd.AddCommand(createSecretCommand)
	createSecretParams = &createSecretCmdParams{}
	createSecretCommand.Flags().StringArrayVar(&createSecretParams.Data, "from-literal", createSecretParams.Data, "Specify a key and literal value to insert in secret (i.e. my-key=some-value)")
	createSecretParams.Scope = createSecretCommand.Flags().StringP("scope", "s", defaultSecretScope, "The scope of the secret")
}
