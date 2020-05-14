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

type approvalTriggeredStruct struct {
	Project *string `json:"project"`
	Stage   *string `json:"stage"`
	Service *string `json:"service"`
}

var approvalTriggered approvalTriggeredStruct

// evaluationDoneCmd represents the approval.triggered command
var approvalTriggeredCmd = &cobra.Command{
	Use:          "approval.triggered",
	Short:        "Returns the latest Keptn sh.keptn.events.approval.triggered event from a specific project/stage/service",
	Long:         `Returns the latest Keptn sh.keptn.events.approval.triggered event from a specific project/stage/service.`,
	Example:      `keptn get event approval.triggered --project=sockshop --stage=staging --service=carts`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to get approval.triggered event", logging.InfoLevel)

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			evaluationDoneEvt, err := eventHandler.GetEvent(*evaluationDone.KeptnContext, keptnevents.EvaluationDoneEventType)
			if err != nil {
				logging.PrintLog("Get approval.triggered event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("%s", *err.Message)
			}

			if evaluationDoneEvt == nil {
				logging.PrintLog("No event returned", logging.QuietLevel)
				return nil
			}

			event, _ := json.Marshal(evaluationDoneEvt)
			fmt.Println(string(event))

		} else {
			fmt.Println("Skipping execution due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	getEventCmd.AddCommand(approvalTriggeredCmd)

	approvalTriggered.Project = approvalTriggeredCmd.Flags().StringP("project", "p", "",
		"The name of a project from which to retrieve an approval.triggered event")
	approvalTriggeredCmd.MarkFlagRequired("project")
	approvalTriggered.Stage = approvalTriggeredCmd.Flags().StringP("stage", "s", "",
		"The name of a stage within a project from which to retrieve an approval.triggered event")
	approvalTriggeredCmd.MarkFlagRequired("stage")
	approvalTriggered.Service = approvalTriggeredCmd.Flags().StringP("service", "s", "",
		"The name of a service within a project from which to retrieve an approval.triggered event")
}
