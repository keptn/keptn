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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
)

const timeout = 60

var eventFilePath *string
var openWebSocketConnection bool

// sendEventCmd represents the send command
var sendEventCmd = &cobra.Command{
	Use:   "event --file=FILEPATH --stream-websocket",
	Short: "Sends a keptn event.",
	Long: `Allows to send an arbitrary keptn event that is defined in the passed file.

Example:
	keptn send event --file=./new_artifact_event.json --stream-websocket`,
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

		var body interface{}
		json.Unmarshal([]byte(eventString), &body)

		eventURL := endPoint
		eventURL.Path = "v1/event"

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", eventURL.String()), logging.VerboseLevel)
		if !mocking {
			data, err := json.Marshal(body)
			req, err := http.NewRequest("POST", eventURL.String(), bytes.NewReader(data))

			// Add signature header
			req.Header.Set("x-token", apiToken)
			req.Header.Set("Content-Type", "application/cloudevents+json")

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			client := &http.Client{Timeout: timeout * time.Second, Transport: tr}
			resp, err := client.Do(req)
			if err != nil {
				logging.PrintLog("Send event was unsuccessful", logging.QuietLevel)
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				logging.PrintLog("Event could not be processed", logging.QuietLevel)
				return errors.New(resp.Status)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			// check for responseCE to include token
			if body == nil || len(body) == 0 {
				logging.PrintLog("Response is empty", logging.InfoLevel)
				return nil
			}

			logging.PrintLog("Event sent to keptn", logging.InfoLevel)
			// open a web socket connection if the open-web-socket flag is set
			if openWebSocketConnection {
				return websockethelper.PrintWSContentByteResponse(body, endPoint)
			}

		} else {
			fmt.Println("Skipping send-new artifact due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	sendCmd.AddCommand(sendEventCmd)
	eventFilePath = sendEventCmd.Flags().StringP("file", "f", "", "The file containing the event as Cloud Event in JSON.")
	sendCmd.PersistentFlags().BoolVarP(&openWebSocketConnection, "stream-websocket", "s", false, "Stream websocket communication of keptn messages")
}
