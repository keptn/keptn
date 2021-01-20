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
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type evaluationDoneStruct struct {
	KeptnContext *string `json:"keptnContext"`
}

var evaluationDone evaluationDoneStruct

// getEvaluationFinishedCmd represents the evaluation.finished command
var getEvaluationFinishedCmd = &cobra.Command{
	Use:          "evaluation.finished",
	Short:        "Returns the latest Keptn sh.keptn.event.evaluation.finished event from a specific Keptn context",
	Long:         `Returns the latest Keptn sh.keptn.event.evaluation.finished event from a specific Keptn context.`,
	Example:      `keptn get event evaluation.finished --keptn-context=1234-5678-90ab-cdef`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(`NOTE: The "keptn get event evaluation.finished" command is DEPRECATED and will be removed in a future release`)
		fmt.Println(`Use "keptn get event evaluation.finished" instead`)
		fmt.Println()

		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to get evaluation.finished event", logging.InfoLevel)

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			evaluationDoneEvts, err := eventHandler.GetEvents(&apiutils.EventFilter{
				KeptnContext: *evaluationDone.KeptnContext,
				EventType:    keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
			})
			if err != nil {
				logging.PrintLog("Get evaluation.finished event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("%s", *err.Message)
			}

			if len(evaluationDoneEvts) == 0 {
				logging.PrintLog("No event returned", logging.QuietLevel)
				return nil
			} else if len(evaluationDoneEvts) == 1 {
				eventsJSON, _ := json.MarshalIndent(evaluationDoneEvts[0], "", "	")
				fmt.Println(string(eventsJSON))
			} else {
				eventsJSON, _ := json.MarshalIndent(evaluationDoneEvts, "", "	")
				fmt.Println(string(eventsJSON))
			}
		} else {
			fmt.Println("Skipping send evaluation-start due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	getEventCmd.AddCommand(getEvaluationFinishedCmd)

	evaluationDone.KeptnContext = getEvaluationFinishedCmd.Flags().StringP("keptn-context", "", "",
		"The ID of a Keptn context from which to retrieve an evaluation.finished event")
	getEvaluationFinishedCmd.MarkFlagRequired("keptn-context")
}
