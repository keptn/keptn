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
	"encoding/json"
	"errors"
	"fmt"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"

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
	Short: "Sends a Keptn event",
	Long: `Allows to send an arbitrary Keptn event that is defined in the passed file.

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

		apiEvent := apimodels.Event{}
		err = json.Unmarshal([]byte(eventString), &apiEvent)
		if err != nil {
			logging.PrintLog("Could not unmarshal event", logging.QuietLevel)
			return fmt.Errorf("Could not unmarshal event. %s", err.Error())
		}

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := eventHandler.SendEvent(apiEvent)
			if err != nil {
				logging.PrintLog("Send event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Send event was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available and stream-websocket flag is true, open WebSocket communication
			if eventContext != nil && openWebSocketConnection {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint)
			}
			return nil
		}

		fmt.Println("Skipping send-new artifact due to mocking flag set to true")
		return nil
	},
}

func init() {
	sendCmd.AddCommand(sendEventCmd)
	eventFilePath = sendEventCmd.Flags().StringP("file", "f", "", "The file containing the event as Cloud Event in JSON.")
	sendCmd.PersistentFlags().BoolVarP(&openWebSocketConnection, "stream-websocket", "s", false, "Stream websocket communication of keptn messages")
}
