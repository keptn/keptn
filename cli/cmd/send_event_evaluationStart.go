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

// NOTE: THIS COMMAND WILL BE REMOVED, THUS THE WHOLE FILE WILL BE REMOVED IN A FUTURE RELEASE

import (
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
		triggerEvaluation := triggerEvaluationStruct{
			Project:   evaluationStart.Project,
			Stage:     evaluationStart.Stage,
			Service:   evaluationStart.Service,
			Timeframe: evaluationStart.Timeframe,
			Start:     evaluationStart.Start,
			End:       evaluationStart.End,
			Labels:    evaluationStart.Labels,
			Watch:     evaluationStart.Watch,
			WatchTime: evaluationStart.WatchTime,
			Output:    evaluationStart.Output,
		}
		return doTriggerEvaluation(triggerEvaluation)
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
