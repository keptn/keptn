package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type configureBridgeCmdParams struct {
	User     *string
	Password *string
	Read     *bool
}

type configureBridgeAPIPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
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
		fmt.Println("You can get your login credentials e.g. by using the following kubectl commands:")
		fmt.Println("Username - kubectl get secret -n keptn bridge-credentials -o jsonpath=\"{.data.BASIC_AUTH_USERNAME}\" | base64 --decode")
		fmt.Println("Password - kubectl get secret -n keptn bridge-credentials -o jsonpath=\"{.data.BASIC_AUTH_PASSWORD}\" | base64 --decode")
		fmt.Println()
		fmt.Println("For editing the login credentials please use: 'kubectl edit secrets -n keptn bridge-credentials'")
		fmt.Println("In order to apply the new credentials you need to restart the Keptn bridge:")
		fmt.Println("kubectl -n keptn rollout restart deployment bridge")
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
