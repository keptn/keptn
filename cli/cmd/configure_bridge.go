package cmd

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strconv"

	keptnutils "github.com/keptn/go-utils/pkg/api/utils"
)

type configureBridgeCmdParams struct {
	Action *string
}

type exposeBridgeAPIPayload struct {
	Expose bool `json:"expose"`
}

type exposeBridgeAPIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var configureBridgeParams *configureBridgeCmdParams

const actionExpose = "expose"
const actionLockdown = "lockdown"

var bridgeCmd = &cobra.Command{
	Use:   "bridge --action=[expose|lockdown]",
	Short: "Exposes or locks down the bridge",
	Long: `Exposes or locks down the Keptn's bridge.

When exposing Keptn's Bridge it will be available publicly. 
Make sure to protect Keptn's Bridge using Basic authentication.
`,
	Example: `keptn configure bridge --action=expose
	keptn configure bridge --action=lockdown`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return verifyConfigureBridgeParams(configureBridgeParams)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		configureBridgeEndpoint := endpoint.Scheme + "://" + endpoint.Host + "/v1/configure/bridge/expose"
		return configureBridge(configureBridgeEndpoint, apiToken, configureBridgeParams)
	},
}

func configureBridge(endpoint string, apiToken string, configureBridgeParams *configureBridgeCmdParams) error {
	doExpose := *configureBridgeParams.Action == "expose"
	payload := strconv.FormatBool(doExpose)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext:     keptnutils.ResolveXipIoWithContext,
		},
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader([]byte(payload)))
	req.Header.Add("x-token", apiToken)
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not complete command: " + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Could not complete command: " + string(body))
	}

	if doExpose {
		if err != nil {
			return errors.New("Could not " + *configureBridgeParams.Action + " bridge: " + err.Error())
		}
		fmt.Printf("Bridge exposed successfully. You can reach it here: https://%s\n", strings.Trim(strings.TrimSpace(string(body)), "\""))
		// Todo: migrate docs for exposing keptn's bridge into keptn.github.io
		fmt.Printf("Make sure to add basic authentication as described here: https://github.com/keptn/keptn/blob/master/bridge/README.md#setting-up-basic-authentication")
	} else {
		if err != nil {
			return errors.New("Could not " + *configureBridgeParams.Action + " bridge: " + err.Error())
		}
		fmt.Println("Bridge locked down successfully. Disabled public access.")
	}
	return nil
}

func verifyConfigureBridgeParams(configureBridgeParams *configureBridgeCmdParams) error {
	if *configureBridgeParams.Action != actionExpose && *configureBridgeParams.Action != actionLockdown {
		return errors.New("Invalid value " + *configureBridgeParams.Action + " 'action'. Must provide either '--action=expose' or '--action=lockdown'")
	}
	return nil
}

func init() {
	configureCmd.AddCommand(bridgeCmd)
	configureBridgeParams = &configureBridgeCmdParams{}

	configureBridgeParams.Action = bridgeCmd.Flags().StringP("action", "a", "", "The action to perform [expose,lockdown]")
	_ = configureCmd.MarkFlagRequired("action")
}
