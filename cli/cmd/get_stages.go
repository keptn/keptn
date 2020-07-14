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
	"strings"
	"text/tabwriter"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

type getStageStruct struct {
	project *string
}

var stageParameter getStageStruct

// evaluationDoneCmd represents the evaluation-done command
var getStageCmd = &cobra.Command{
	Use:     "stage",
	Aliases: []string{"stages"},
	Short:   "Get details of a stage",
	Long:    `Get all stages or details of a stage from a given keptn project`,
	Example: `keptn get stages --project=sockshop
NAME           CREATION DATE                 
staging        2020-04-06T14:37:45.210Z
production     2020-04-06T14:37:45.210Z

keptn get stage staging --project sockshop
NAME           CREATION DATE                 
staging        2020-04-06T14:37:45.210Z
	`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		var stageName string
		stageName = strings.Join(args, " ")

		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		stagesHandler := apiutils.NewAuthenticatedStageHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		if !mocking {
			stages, err := stagesHandler.GetAllStages(*stageParameter.project)
			if err != nil {
				fmt.Println(err)
				return errors.New("Could not retrieve stages for project " + *stageParameter.project)
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)
			fmt.Fprintln(w, "NAME\tCREATION DATE")

			if stageName != "" {

				for _, stage := range stages {
					if stage.StageName == stageName {
						fmt.Fprintln(w, stage.StageName+"\tn/a")
					}
				}
			} else {
				for _, stage := range stages {
					fmt.Fprintln(w, stage.StageName+"\tn/a")
				}
			}
			w.Flush()
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getStageCmd)

	stageParameter.project = getStageCmd.Flags().StringP("project", "", "",
		"keptn project name")

	getStageCmd.MarkFlagRequired("project")
}
