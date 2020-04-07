package cmd

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
	"net/http"
)

type configureBridgeCmdParams struct {
	Action *string
}

var configureBridgeParams *configureBridgeCmdParams

const actionExpose = "expose"
const actionLockdown = "lockdown"

var bridgeCmd = &cobra.Command{
	Use:   "--action=[expose|lockdown]",
	Short: "Exposes or locks down the bridge",
	Long: `Exposes or locks down the Keptn's bridge.

Example:
	keptn configure bridge --action=expose
	keptn configure bridge --action=lockdown`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if configureBridgeParams.Action == nil {
			return errors.New("Missing flag 'action'. Must provide either '--action=expose' or '--action=lockdown'")
		}
		if *configureBridgeParams.Action != actionExpose && *configureBridgeParams.Action != actionLockdown {
			return errors.New("Invalid value " + *configureBridgeParams.Action + " 'action'. Must provide either '--action=expose' or '--action=lockdown'")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		// TODO: send the request to the endpoint implemented in https://github.com/keptn/keptn/issues/1153
		client := &http.Client{}
		req, err := http.NewRequest("POST", endPoint.RequestURI()+"", nil)
		req.Header.Add("x-token", apiToken)

		resp, err := client.Do(req)
		fmt.Printf("%v", resp)

		return nil
	},
}

func init() {
	configureCmd.AddCommand(bridgeCmd)
	configureBridgeParams = &configureBridgeCmdParams{}

	configureBridgeParams.Action = bridgeCmd.Flags().StringP("action", "a", "", "The action to perform [expose,lockdown]")
}
