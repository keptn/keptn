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
	"github.com/keptn/go-utils/pkg/commonutils"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

var eventFilePath *string

// sendEventCmd represents the send command
var sendEventCmd = &cobra.Command{
	Use:   "event --file=FILEPATH",
	Short: "Sends an event to Keptn",
	Long: `Sends an arbitrary Keptn event that is defined in the provided JSON file.
An event has to follow the CloudEvents specification (https://cloudevents.io/) in version 0.2 and has to be written in JSON.
In addition, the payload of the CloudEvent needs to follow the Keptn spec (https://github.com/keptn/spec/blob/0.1.4/cloudevents.md).
`,
	Example:      `keptn send event --file=./new_artifact_event.json`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := doSendEventPreRunChecks(); err != nil {
			return err
		}
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}
		eventString, err := commonutils.ReadFile(*eventFilePath)
		if err != nil {
			return err
		}

		apiEvent := apimodels.KeptnContextExtendedCE{}
		err = json.Unmarshal(eventString, &apiEvent)
		if err != nil {
			return fmt.Errorf("Failed to map event to API event model. %s", err.Error())
		}

		if endPointErr := CheckEndpointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			_, err := apiHandler.SendEvent(apiEvent)
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

func doSendEventPreRunChecks() error {
	eventString, err := commonutils.ReadFile(*eventFilePath)
	if err != nil {
		return err
	}
	var body interface{}
	return json.Unmarshal(eventString, &body)
}

func init() {
	sendCmd.AddCommand(sendEventCmd)
	eventFilePath = sendEventCmd.Flags().StringP("file", "f", "", "The file containing the event as Cloud Event in JSON.")
	sendEventCmd.MarkFlagRequired("file")
}
