package cmd

import (
	"errors"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

type getSecretsStruct struct {
	outputFormat *string
}

var getSecrets getSecretsStruct

var getSecretsCommand = &cobra.Command{
	Use:     `secret`,
	Aliases: []string{"secrets"},
	Short:   "Gets the list of secrets managed by the Keptn secret-service",
	Example: `keptn get secrets
NAME
my-secret-1
my-secret-2

keptn get secrets -output=yaml  # Returns secret list in YAML format

keptn get secrets -output=json  # Returns secret list in JSON format
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if *getSecrets.outputFormat != "" {
			if *getProject.outputFormat != "yaml" && *getProject.outputFormat != "json" {
				return errors.New("Invalid output format, only yaml or json allowed")
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := NewSecretCmdHandler(credentialmanager.NewCredentialManager(assumeYes))
		if err != nil {
			return nil
		}
		return handler.GetSecrets(*getSecrets.outputFormat)
	},
}

func init() {
	getCmd.AddCommand(getSecretsCommand)
	getSecrets.outputFormat = getSecretsCommand.Flags().StringP("output", "o", "",
		"Output format. One of json|yaml")
}
