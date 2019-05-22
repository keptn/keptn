// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
)

var eventFilePath *string

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends a keptn event.",
	Long: `Allows to send arbitrary keptn events, which are defined in the passed file.

Example:
	keptn send --event=new_artifact.json`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		eventString, err := utils.ReadFile(*eventFilePath)
		if err != nil {
			return err
		}
		var body interface{}
		return json.Unmarshal([]byte(eventString), &body)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		eventString, err := utils.ReadFile(*eventFilePath)
		if err != nil {
			return err
		}

		utils.PrintLog("Starting to send an event", utils.InfoLevel)

		eventURL := endPoint
		eventURL.Path = "event"
		utils.PrintLog(fmt.Sprintf("Connecting to server %s", eventURL.String()), utils.VerboseLevel)
		if !mocking {
			bodyBytes := []byte(eventString)
			req, err := http.NewRequest("POST", eventURL.String(), bytes.NewBuffer(bodyBytes))

			mac := hmac.New(sha1.New, []byte(apiToken))
			mac.Write(bodyBytes)
			signatureVal := mac.Sum(nil)
			sha1Hash := "sha1=" + fmt.Sprintf("%x", signatureVal)

			// Add signature header
			req.Header.Set("X-Keptn-Signature", sha1Hash)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				utils.PrintLog("Send event was unsuccessful", utils.QuietLevel)
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				utils.PrintLog("Send event was unsuccessful", utils.QuietLevel)
				return errors.New(resp.Status)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			// check for responseCE to include token
			if body == nil || len(body) == 0 {
				utils.PrintLog("Response is empty", utils.InfoLevel)
				return nil
			}
			return websockethelper.PrintWSContentByteResponse(body)
		} else {
			fmt.Println("Skipping send-new artifact due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	eventFilePath = serviceCmd.Flags().StringP("event", "e", "", "The file containing the event as Cloud Event in JSON.")
}
