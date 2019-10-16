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
	"errors"
	"fmt"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type evaluationDoneStruct struct {
	KeptnContext *string `json:"keptnContext"`
}

var evaluationDone evaluationDoneStruct

// evaluationDoneCmd represents the evaluation-done command
var evaluationDoneCmd = &cobra.Command{
	Use:   "evaluation-done",
	Short: "",
	Long: `
	
Example:
	keptn get event evaluation-done --keptn-context 1234-5678-90`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to get evaluation-done event", logging.InfoLevel)

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			evaluationDoneEvt, err := eventHandler.GetEvent("", keptnevents.EvaluationDoneEventType)
			if err != nil {
				logging.PrintLog("Get evaluation-done event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Get evaluation-done event was unsuccessful. %s", *err.Message)
			}

			if evaluationDoneEvt == nil {
				logging.PrintLog("Response is nil", logging.QuietLevel)
				return nil
			}

			fmt.Println("print Event")

		} else {
			fmt.Println("Skipping send evaluation-start due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	sendEventCmd.AddCommand(evaluationDoneCmd)

	evaluationDone.KeptnContext = evaluationDoneCmd.Flags().StringP("keptn-context", "", "",
		"The ID of a Keptn context of an evaluation step")
	evaluationDoneCmd.MarkFlagRequired("keptn-context")
}
