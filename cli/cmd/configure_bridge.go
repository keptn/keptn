package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type configureBridgeCmdParams struct {
	User     *string
	Password *string
	Read     *bool
}

var configureBridgeParams *configureBridgeCmdParams

var bridgeCmd = &cobra.Command{
	Use:          "bridge --user=<user> --password=<password>",
	Short:        "Configures the credentials for the Keptn Bridge",
	Long:         `Configures the credentials for the Keptn Bridge.`,
	Example:      `keptn configure bridge --user=<user> --password=<password>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Warning: From version 0.9.0 of Keptn this command is not supported anymore!")
		fmt.Println()
		if isBoolFlagSet(configureBridgeParams.Read) {
			fmt.Println("You can get your login credentials e.g. by using the following kubectl commands:")
			fmt.Println("Username - kubectl get secret -n keptn bridge-credentials -o jsonpath=\"{.data.BASIC_AUTH_USERNAME}\" | base64 --decode")
			fmt.Println("Password - kubectl get secret -n keptn bridge-credentials -o jsonpath=\"{.data.BASIC_AUTH_PASSWORD}\" | base64 --decode")
		} else {
			fmt.Println(getReplaceSecretCommand(*configureBridgeParams))
		}

		fmt.Println()
		fmt.Println("The URL to your Keptn Bridge can be retrieved using 'keptn status'")
		return nil
	},
}

func init() {
	configureCmd.AddCommand(bridgeCmd)
	configureBridgeParams = &configureBridgeCmdParams{}
	configureBridgeParams.Read = bridgeCmd.Flags().BoolP("output", "o", false, "Print the current credentials")
	configureBridgeParams.User = bridgeCmd.Flags().StringP("user", "u", "", "The user name to login to the bridge")
	configureBridgeParams.Password = bridgeCmd.Flags().StringP("password", "p", "", "The password to login to the bridge")
}

func getReplaceSecretCommand(cmdParams configureBridgeCmdParams) string {
	user := "${BRIDGE_USER}"
	password := "${BRIDGE_PASSWORD}"
	if isStringFlagSet(cmdParams.User) {
		user = *cmdParams.User
	}
	if isStringFlagSet(cmdParams.Password) {
		password = *cmdParams.Password
	}

	builder := strings.Builder{}

	builder.WriteString("For editing the login credentials please use the following command:\n\n")
	builder.WriteString(fmt.Sprintf("kubectl create secret -n ${KEPTN_NAMESPACE} generic bridge-credentials --from-literal=\"BASIC_AUTH_USERNAME=%s\" --from-literal=\"BASIC_AUTH_PASSWORD=%s\" -oyaml --dry-run=client | kubectl replace -f -\n", user, password))
	builder.WriteString("kubectl -n keptn rollout restart deployment bridge")

	return builder.String()
}
