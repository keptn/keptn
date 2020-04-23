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

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type evaluationDoneStruct struct {
	KeptnContext *string `json:"keptnContext"`
}

var evaluationDone evaluationDoneStruct

// evaluationDoneCmd represents the evaluation-done command
var evaluationDoneCmd = &cobra.Command{
	Use:   "evaluation-done",
	Short: "Returns the latest Keptn sh.keptn.events.evaluation-done event from a specific Keptn context",
	Long: `Returns the latest Keptn sh.keptn.events.evaluation-done event from a specific Keptn context.`,
	Example: `keptn get event evaluation-done --keptn-context=1234-5678-90ab-cdef`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to get evaluation-done event", logging.InfoLevel)

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			evaluationDoneEvt, err := eventHandler.GetEvent(*evaluationDone.KeptnContext, keptnevents.EvaluationDoneEventType)
			if err != nil {
				logging.PrintLog("Get evaluation-done event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("%s", *err.Message)
			}

			if evaluationDoneEvt == nil {
				logging.PrintLog("No event returned", logging.QuietLevel)
				return nil
			}

			event, _ := json.Marshal(evaluationDoneEvt)
			fmt.Println(string(event))

		} else {
			fmt.Println("Skipping send evaluation-start due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	getEventCmd.AddCommand(evaluationDoneCmd)

	evaluationDone.KeptnContext = evaluationDoneCmd.Flags().StringP("keptn-context", "", "",
		"The ID of a Keptn context from which to retrieve an evaluation-done event")
	evaluationDoneCmd.MarkFlagRequired("keptn-context")
}
