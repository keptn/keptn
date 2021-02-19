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
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type evaluationStartStruct struct {
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

var evaluationStart evaluationStartStruct

// evaluationStartCmd represents the start-evaluation command
var evaluationStartCmd = &cobra.Command{
	Use: "start-evaluation",
	Short: "Sends an start-evaluation event to Keptn to evaluate a test " +
		"for the specified service in the provided project and stage",
	Long: `Sends a start-evaluation event to Keptn in order to evaluate a test for the specified service in the provided project and stage. 

* This command takes the project (*--project*), stage (*--stage*), and the service (*--service*), which should be evaluated. 
* It is necessary to specify a time frame (*--timeframe*) of the evaluation. If, for example, the 
flag is set to *--timeframe=5m*, the evaluation is conducted for the last 5 minutes. 
* To specify a particular starting point, the flag *--start* flag can be used. In this case, the specified time frame is added to the starting point.
`,
	Deprecated: `Use "keptn trigger evaluation" instead`,
	Example: `keptn send event start-evaluation --project=sockshop --stage=hardening --service=carts --timeframe=5m --start=2019-10-31T11:59:59

keptn send event start-evaluation --project=sockshop --stage=hardening --service=carts --start=2019-10-31T11:59:59 --end=2019-10-31T12:04:59 --labels=test-id=1234,test-name=performance-test
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to send a start-evaluation event to evaluate the service "+
			*evaluationStart.Service+" in project "+*evaluationStart.Project, logging.InfoLevel)

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		startPoint := ""
		if evaluationStart.Start != nil {
			startPoint = *evaluationStart.Start
		}

		endDatePoint := ""
		if evaluationStart.End != nil {
			endDatePoint = *evaluationStart.End
		}

		start, end, err := getStartEndTime(startPoint, endDatePoint, *evaluationStart.Timeframe)
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
				*evaluationStart.Project,
				*evaluationStart.Stage,
				*evaluationStart.Service,
				apimodels.Evaluation{
					Start:  start.Format("2006-01-02T15:04:05"),
					End:    end.Format("2006-01-02T15:04:05"),
					Labels: *evaluationStart.Labels,
				},
			)

			if err != nil {
				logging.PrintLog("Send start-evaluation was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Send start-evaluation was unsuccessful. %s", *err.Message)
			}

			if response == nil {
				logging.PrintLog("No event returned", logging.QuietLevel)
				return nil
			}

			if *evaluationStart.Watch {
				eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
				filter := apiutils.EventFilter{
					KeptnContext: *response.KeptnContext,
					Project:      *evaluationStart.Project,
				}
				watcher := NewDefaultWatcher(eventHandler, filter, time.Duration(*evaluationStart.WatchTime)*time.Second)
				PrintEventWatcher(watcher, *evaluationStart.Output, os.Stdout)
			}

			return nil
		}

		fmt.Println("Skipping send start-evaluation due to mocking flag set to true")
		return nil
	},
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
		"The time frame from which the evaluation data should be gathered (can not be used together with --end)")

	evaluationStart.Start = evaluationStartCmd.Flags().StringP("start", "", "",
		"The starting point from which to start the evaluation in UTC")

	evaluationStart.End = evaluationStartCmd.Flags().StringP("end", "", "",
		"The end point to which the evaluation data should be gathered in UTC (can not be used together with --timeframe)")
	evaluationStart.Labels = evaluationStartCmd.Flags().StringToStringP("labels", "l", nil, "Additional labels to be provided to the lighthouse service")

	evaluationStart.Output = AddOutputFormatFlag(evaluationStartCmd)
	evaluationStart.Watch = AddWatchFlag(evaluationStartCmd)
	evaluationStart.WatchTime = AddWatchTimeFlag(evaluationStartCmd)
}
