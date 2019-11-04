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
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type evaluationStartStruct struct {
	Project   *string `json:"project"`
	Stage     *string `json:"stage"`
	Service   *string `json:"service"`
	Timeframe *string `json:"timeframe"`
}

var evaluationStart evaluationStartStruct

// evaluationStartCmd represents the evaluation.start command
var evaluationStartCmd = &cobra.Command{
	Use: "evaluation.start",
	Short: "Sends an evaluation.start event to Keptn in order to evaluate a test" +
		"for the specified service in the provided project",
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

		end, start, err := getStartEndTime(*evaluationStart.Timeframe)
		if end == nil || start == nil || err != nil {
			logging.PrintLog(fmt.Sprintf("Start and end time of evaluation time frame not set: %s", err.Error()), logging.QuietLevel)
			return fmt.Errorf("Start and end time of evaluation time frame not set: %s", err.Error())
		}

		startEvaluationEventData := keptnevents.StartEvaluationEventData{
			Project: *evaluationStart.Project,
			Service: *evaluationStart.Service,
			Stage:   *evaluationStart.Stage,
			Start:   start.Format("2006-01-02T15:04:05-0700"),
			End:     end.Format("2006-01-02T15:04:05-0700"),
		}

		keptnContext := uuid.New().String()
		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuration-change")
		contentType := "application/json"
		sdkEvent := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          keptnContext,
				Type:        keptnevents.StartEvaluationEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: startEvaluationEventData,
		}

		eventByte, err := sdkEvent.MarshalJSON()
		if err != nil {
			return fmt.Errorf("Failed to marshal cloud event. %s", err.Error())
		}

		apiEvent := apimodels.Event{}
		err = json.Unmarshal(eventByte, &apiEvent)
		if err != nil {
			return fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
		}

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			responseEvent, err := eventHandler.SendEvent(apiEvent)
			if err != nil {
				logging.PrintLog("Send evaluation.start was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Send evaluation.start was unsuccessful. %s", *err.Message)
			}

			if responseEvent == nil {
				logging.PrintLog("No event returned", logging.QuietLevel)
				return nil
			}

			fmt.Println("ID of Keptn context: " + *responseEvent.KeptnContext)
			return nil
		}

		fmt.Println("Skipping send evaluation.start due to mocking flag set to true")
		return nil
	},
}

func getStartEndTime(timeframe string) (*time.Time, *time.Time, error) {
	end := time.Now()
	start := time.Now()

	errMsg := "The time frame format is invalid. Use the format [duration]m, e.g.: 5m"

	i := strings.Index(timeframe, "m")
	if i > -1 {
		minutesStr := timeframe[:i]
		minutes, err := strconv.Atoi(minutesStr)
		if err != nil {
			return nil, nil, fmt.Errorf(errMsg)
		}
		minutesOffset := time.Minute * time.Duration(-minutes)
		start = start.Add(minutesOffset)

		return &end, &start, nil
	}

	return nil, nil, fmt.Errorf(errMsg)
}

func init() {
	sendEventCmd.AddCommand(evaluationStartCmd)

	evaluationStart.Project = evaluationStartCmd.Flags().StringP("project", "", "",
		"The project containing the service to be evaluated")
	evaluationStartCmd.MarkFlagRequired("project")

	evaluationStart.Stage = evaluationStartCmd.Flags().StringP("stage", "", "",
		"The stage containing the service to be evaluated")
	evaluationStartCmd.MarkFlagRequired("stage")

	evaluationStart.Service = evaluationStartCmd.Flags().StringP("service", "", "",
		"The service to be evaluated")
	evaluationStartCmd.MarkFlagRequired("service")

	evaluationStart.Timeframe = evaluationStartCmd.Flags().StringP("timeframe", "", "",
		"The time frame from which the evaluation data should be gathered")
	evaluationStartCmd.MarkFlagRequired("service")
}
