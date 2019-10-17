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
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
	"net/url"
)

type evaluationStartStruct struct {
	Project *string `json:"project"`
	Service *string `json:"service"`
}

var evaluationStart evaluationStartStruct

// evaluationStartCmd represents the evaluation.start command
var evaluationStartCmd = &cobra.Command{
	Use: "evaluation.start",
	Short: "Sends an evaluation.start event to Keptn in order to evaluate a test" +
		"for the specified service in the provided project.",
	Long: `Sends an evaluation.start event to Keptn in order to evaluate a test
for the specified service in the provided project.
	
Example:
	keptn send event evaluation.start --project=sockshop --service=carts`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to send an evaluation.start event to evaluate the service "+
			*evaluationStart.Service+" in project "+*evaluationStart.Project, logging.InfoLevel)

		evaluationStartEvent := keptnevents.EvaluationStartEventData{
			Project:      *evaluationStart.Project,
			Service:      *evaluationStart.Service,
		}

		keptnContext := uuid.New().String()
		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuration-change")
		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          keptnContext,
				Type:        keptnevents.EvaluationStartEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: evaluationStartEvent,
		}

		eventURL := endPoint
		eventURL.Path = "v1/event"

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", eventURL.String()), logging.VerboseLevel)
		if !mocking {
			_, response, err := utils.Send(eventURL, event, apiToken)
			if err != nil {
				logging.PrintLog("Send evaluation.start was unsuccessful", logging.QuietLevel)
				return err
			}

			if response == nil {
				logging.PrintLog("Response is nil", logging.QuietLevel)
				return nil
			}

			return websockethelper.PrintWSContentCEResponse(response, endPoint)
		} 
		
		fmt.Println("Skipping send evaluation.start due to mocking flag set to true")
		return nil
	},
}

func init() {
	sendEventCmd.AddCommand(evaluationStartCmd)

	evaluationStart.Project = evaluationStartCmd.Flags().StringP("project", "", "",
		"The project containing the service which will be evaluated")
	evaluationStartCmd.MarkFlagRequired("project")

	evaluationStart.Service = evaluationStartCmd.Flags().StringP("service", "", "",
		"The service which will be evaluated")
	evaluationStartCmd.MarkFlagRequired("service")
}
