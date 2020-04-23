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

	"github.com/keptn/keptn/cli/pkg/file"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

const timeout = 60

var eventFilePath *string

// sendEventCmd represents the send command
var sendEventCmd = &cobra.Command{
	Use:   "event --file=FILEPATH --stream-websocket",
	Short: "Sends a Keptn event",
	Long: `Allows to send an arbitrary Keptn event that is defined in the provided JSON file.
An event has to follow the Cloud Events specification (https://cloudevents.io/) in version 0.2 and has to be written in JSON.
In addition, the payload of the Cloud Event needs to follow the Keptn spec (https://github.com/keptn/spec/blob/0.1.3/cloudevents.md).

For convenience, this command offers the *--stream-websocket* flag to open a web socket communication to Keptn. Consequently, messages from the receiving Keptn service, which processes the event, are sent to the CLI via WebSocket.
	`,
	Example: `keptn send event --file=./new_artifact_event.json --stream-websocket`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		eventString, err := file.ReadFile(*eventFilePath)
		if err != nil {
			return err
		}
		var body interface{}
		return json.Unmarshal([]byte(eventString), &body)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		eventString, err := file.ReadFile(*eventFilePath)
		if err != nil {
			return err
		}

		apiEvent := apimodels.KeptnContextExtendedCE{}
		err = json.Unmarshal([]byte(eventString), &apiEvent)
		if err != nil {
			return fmt.Errorf("Failed to map event to API event model. %s", err.Error())
		}

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			_, err := eventHandler.SendEvent(apiEvent)
			if err != nil {
				logging.PrintLog("Send event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Send event was unsuccessful. %s", *err.Message)
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
	sendEventCmd.MarkFlagRequired("file")
}
