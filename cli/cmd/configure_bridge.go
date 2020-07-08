package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"

	keptnutils "github.com/keptn/go-utils/pkg/api/utils"
)

type configureBridgeCmdParams struct {
	User     *string
	Password *string
}

type configureBridgeAPIPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type exposeBridgeAPIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var configureBridgeParams *configureBridgeCmdParams

const basicAuthDocuURL = "https://keptn.sh/docs/0.7.0/reference/keptnsbridge/#enable-authentication"

var bridgeCmd = &cobra.Command{
	Use:          "bridge --user=<user> --password=<password>",
	Short:        "Configures the credentials for the Keptn Bridge",
	Long:         `Configures the credentials for the Keptn Bridge.`,
	Example:      `keptn configure bridge --user=<user> --password=<passsord>`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return verifyConfigureBridgeParams(configureBridgeParams)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		configureBridgeEndpoint := endpoint.Scheme + "://" + endpoint.Host + "/v1/configure/bridge"
		return configureBridge(configureBridgeEndpoint, apiToken, configureBridgeParams)
	},
}

func configureBridge(endpoint string, apiToken string, configureBridgeParams *configureBridgeCmdParams) error {
	exposeBridgeParams := configureBridgeAPIPayload{
		User:     *configureBridgeParams.User,
		Password: *configureBridgeParams.Password,
	}

	payload, err := json.Marshal(exposeBridgeParams)
	if err != nil {
		fmt.Println("Could not complete command: " + err.Error())
		return err
	}
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

	if err != nil {
		return errors.New("Could not configure bridge: " + err.Error())
	}
	fmt.Println("Bridge credentials configured successfully")
	return nil
}

func verifyConfigureBridgeParams(configureBridgeParams *configureBridgeCmdParams) error {
	if configureBridgeParams.User == nil || *configureBridgeParams.User == "" {
		return errors.New("please specify a user name for exposing the bridge using the '--user=<username>' flag")
	}
	if configureBridgeParams.Password == nil || *configureBridgeParams.Password == "" {
		return errors.New("please specify a password for exposing the bridge using the '--password=<password>' flag")
	}
	return nil
}

func init() {
	configureCmd.AddCommand(bridgeCmd)
	configureBridgeParams = &configureBridgeCmdParams{}

	configureBridgeParams.User = bridgeCmd.Flags().StringP("user", "u", "", "The user name to login to the bridge")
	_ = configureCmd.MarkFlagRequired("user")
	configureBridgeParams.Password = bridgeCmd.Flags().StringP("password", "p", "", "The password to login to the bridge")
	_ = configureCmd.MarkFlagRequired("password")
}
