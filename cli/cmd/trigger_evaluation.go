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
	"os"
	"strconv"
	"strings"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type triggerEvaluationStruct struct {
	Project   *string            `json:"project"`
	Stage     *string            `json:"stage"`
	Service   *string            `json:"service"`
	Timeframe *string            `json:"timeframe"`
	Start     *string            `json:"start"`
	End       *string            `json:"end"`
	Labels    *map[string]string `json:"labels"`
	Watch     *bool
	WatchTime *int
	Output    *string
}

var triggerEvaluation triggerEvaluationStruct

var triggerEvaluationCmd = &cobra.Command{
	Use:   "evaluation",
	Short: "Triggers the evaluation of a test for a service in a project and stage",
	Long: `Triggers the evaluation of a test for a service in a project and stage 

* This command takes the project (--project), stage (--stage), and the service (--service), which should be evaluated. 
* It is necessary to specify a time frame (--timeframe) of the evaluation. If, for example, the 
flag is set to --timeframe=5m, the evaluation is conducted for the last 5 minutes. 
* To specify a particular starting point, the --start flag can be used. In this case, the specified time frame is added to the starting point.
`,
	Example: `keptn trigger evaluation --project=sockshop --stage=hardening --service=carts --timeframe=5m --start=2019-10-31T11:59:59
keptn trigger evaluation --project=sockshop --stage=hardening --service=carts --start=2019-10-31T11:59:59 --end=2019-10-31T12:04:59 --labels=test-id=1234,test-name=performance-test
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return doTriggerEvaluation(triggerEvaluation)
	},
}

func doTriggerEvaluation(triggerEvaluationData triggerEvaluationStruct) error {
	endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
	if err != nil {
		return errors.New(authErrorMsg)
	}

	logging.PrintLog("Starting to trigger evaluation of the service "+
		*triggerEvaluationData.Service+" in project "+*triggerEvaluationData.Project, logging.InfoLevel)

	if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
		return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
			endPointErr)
	}

	startPoint := ""
	if triggerEvaluationData.Start != nil {
		startPoint = *triggerEvaluationData.Start
	}

	endDatePoint := ""
	if triggerEvaluationData.End != nil {
		endDatePoint = *triggerEvaluationData.End
	}

	start, end, err := getStartEndTime(startPoint, endDatePoint, *triggerEvaluationData.Timeframe)
	if start == nil || end == nil || err != nil {
		logging.PrintLog(fmt.Sprintf("Start and end time of evaluation time frame not set: %s", err.Error()), logging.QuietLevel)
		return fmt.Errorf("Start and end time of evaluation time frame not set: %s", err.Error())
	}

	if err != nil {
		return fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
	}

	apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	if !mocking {
		response, err := apiHandler.TriggerEvaluation(
			*triggerEvaluationData.Project,
			*triggerEvaluationData.Stage,
			*triggerEvaluationData.Service,
			apimodels.Evaluation{
				Start:  start.Format("2006-01-02T15:04:05"),
				End:    end.Format("2006-01-02T15:04:05"),
				Labels: *triggerEvaluationData.Labels,
			},
		)

		if err != nil {
			logging.PrintLog("trigger evaluation was unsuccessful", logging.QuietLevel)
			return fmt.Errorf("trigger evaluation was unsuccessful. %s", *err.Message)
		}

		if response == nil {
			logging.PrintLog("No event returned", logging.QuietLevel)
			return nil
		}

		if *triggerEvaluationData.Watch {
			eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
			filter := apiutils.EventFilter{
				KeptnContext: *response.KeptnContext,
				Project:      *triggerEvaluationData.Project,
			}
			watcher := NewDefaultWatcher(eventHandler, filter, time.Duration(*triggerEvaluationData.WatchTime)*time.Second)
			PrintEventWatcher(rootCmd.Context(), watcher, *triggerEvaluationData.Output, os.Stdout)
		}

		return nil
	}

	fmt.Println("Skipping trigger evaluation due to mocking flag set to true")
	return nil
}

func getStartEndTime(startDatePoint string, endDatePoint string, timeframe string) (*time.Time, *time.Time, error) {
	// set default values for start and end time
	dateLayout := "2006-01-02T15:04:05"
	var err error

	minutes := 5 // default timeframe

	// input validation
	if startDatePoint != "" && endDatePoint == "" {
		// if a start date is set, but no end date is set, we require the timeframe to be set
		if timeframe == "" {
			errMsg := "Please provide a timeframe, e.g., --timeframe=5m, or an end date using --end=..."

			return nil, nil, fmt.Errorf(errMsg)
		}
	}
	if endDatePoint != "" && timeframe != "" {
		// can not use end date and timeframe at the same time
		errMsg := "You can not use --end together with --timeframe"

		return nil, nil, fmt.Errorf(errMsg)
	}
	if endDatePoint != "" && startDatePoint == "" {
		errMsg := "start date is required when using an end date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	// parse timeframe
	if timeframe != "" {
		errMsg := "The time frame format is invalid. Use the format [duration]m, e.g.: 5m"

		i := strings.Index(timeframe, "m")

		if i > -1 {
			minutesStr := timeframe[:i]
			minutes, err = strconv.Atoi(minutesStr)
			if err != nil {
				return nil, nil, fmt.Errorf(errMsg)
			}
		} else {
			return nil, nil, fmt.Errorf(errMsg)
		}
	}

	// initialize default values for end and start time
	end := time.Now().UTC()
	start := time.Now().UTC().Add(-time.Duration(minutes) * time.Minute)

	// Parse start date
	if startDatePoint != "" {
		start, err = time.Parse(dateLayout, startDatePoint)

		if err != nil {
			return nil, nil, err
		}
	}

	// Parse end date
	if endDatePoint != "" {
		end, err = time.Parse(dateLayout, endDatePoint)

		if err != nil {
			return nil, nil, err
		}
	}

	// last but not least: if a start date and a timeframe is provided, we set the end date to start date + timeframe
	if startDatePoint != "" && endDatePoint == "" && timeframe != "" {
		minutesOffset := time.Minute * time.Duration(minutes)
		end = start.Add(minutesOffset)
	}

	// ensure end date is greater than start date
	diff := end.Sub(start).Minutes()

	if diff < 1 {
		errMsg := "end date must be at least 1 minute after start date"

		return nil, nil, fmt.Errorf(errMsg)
	}

	return &start, &end, nil
}

func init() {
	triggerCmd.AddCommand(triggerEvaluationCmd)

	triggerEvaluation.Project = triggerEvaluationCmd.Flags().StringP("project", "", "",
		"The project containing the service to be evaluated")
	triggerEvaluationCmd.MarkFlagRequired("project")

	triggerEvaluation.Stage = triggerEvaluationCmd.Flags().StringP("stage", "", "",
		"The stage containing the service to be evaluated")
	triggerEvaluationCmd.MarkFlagRequired("stage")

	triggerEvaluation.Service = triggerEvaluationCmd.Flags().StringP("service", "", "",
		"The service to be evaluated")
	triggerEvaluationCmd.MarkFlagRequired("service")

	triggerEvaluation.Timeframe = triggerEvaluationCmd.Flags().StringP("timeframe", "", "",
		"The time frame from which the evaluation data should be gathered (can not be used together with --end)")

	triggerEvaluation.Start = triggerEvaluationCmd.Flags().StringP("start", "", "",
		"The starting point from which to start the evaluation in UTC")

	triggerEvaluation.End = triggerEvaluationCmd.Flags().StringP("end", "", "",
		"The end point to which the evaluation data should be gathered in UTC (can not be used together with --timeframe)")
	triggerEvaluation.Labels = triggerEvaluationCmd.Flags().StringToStringP("labels", "l", nil, "Additional labels to be provided to the lighthouse service")

	triggerEvaluation.Output = AddOutputFormatFlag(triggerEvaluationCmd)
	triggerEvaluation.Watch = AddWatchFlag(triggerEvaluationCmd)
	triggerEvaluation.WatchTime = AddWatchTimeFlag(triggerEvaluationCmd)
}
