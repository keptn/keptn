package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
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

type exposeBridgeAPIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var configureBridgeParams *configureBridgeCmdParams

var bridgeCmd = &cobra.Command{
	Use:          "bridge --user=<user> --password=<password>",
	Short:        "Configures the credentials for the Keptn Bridge",
	Long:         `Configures the credentials for the Keptn Bridge.`,
	Example:      `keptn configure bridge --user=<user> --password=<password>`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return verifyConfigureBridgeParams(configureBridgeParams)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		endpoint.Path = path.Join(endpoint.Path, "/v1/config/bridge")

		return configureBridge(endpoint.String(), apiToken, configureBridgeParams)
	},
}

func configureBridge(endpoint string, apiToken string, configureBridgeParams *configureBridgeCmdParams) error {
	if configureBridgeParams.Read != nil && *configureBridgeParams.Read {
		creds, err := retrieveBridgeCredentials(endpoint, apiToken)
		if err != nil {
			fmt.Println("Could not retrieve bridge credentials: " + err.Error())
			return err
		}
		fmt.Println("user: " + creds.User)
		fmt.Println("password: " + creds.Password)
		return nil
	}
	return configureBridgeCredentials(endpoint, apiToken, configureBridgeParams)
}

func retrieveBridgeCredentials(endpoint string, apiToken string) (*configureBridgeAPIPayload, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("x-token", apiToken)
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not send request: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("Received not successful response: " + string(body))
	}

	res := &configureBridgeAPIPayload{}
	json.Unmarshal(body, res)

	if err != nil {
		return nil, errors.New("Could not decode bridge credentials: " + err.Error())
	}
	return res, nil
}

func configureBridgeCredentials(endpoint string, apiToken string, configureBridgeParams *configureBridgeCmdParams) error {
	bridgeCredentials := configureBridgeAPIPayload{
		User:     *configureBridgeParams.User,
		Password: *configureBridgeParams.Password,
	}

	payload, err := json.Marshal(bridgeCredentials)
	if err != nil {
		fmt.Println("Could not marshal response payload: " + err.Error())
		return err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader([]byte(payload)))
	req.Header.Add("x-token", apiToken)
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not send request: " + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Received not successful response: " + string(body))
	}

	if err != nil {
		return errors.New("Could not configure bridge: " + err.Error())
	}
	fmt.Println("Bridge credentials configured successfully")
	return nil
}

func verifyConfigureBridgeParams(configureBridgeParams *configureBridgeCmdParams) error {
	if configureBridgeParams.Read != nil && *configureBridgeParams.Read {
		return nil
	}
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
	configureBridgeParams.Password = bridgeCmd.Flags().StringP("password", "p", "", "The password to login to the bridge")

	configureBridgeParams.Read = bridgeCmd.Flags().BoolP("output", "o", false, "Print the current credentials")
}
