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
	"text/tabwriter"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

type getStageStruct struct {
	project *string
}

var stageParameter getStageStruct

// getStageCmd represents the get stage command
var getStageCmd = &cobra.Command{
	Use:     "stage",
	Aliases: []string{"stages"},
	Short:   "Get details of a stage",
	Long:    `Get all stages or details of a stage from a given Keptn project`,
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
		_, _, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		stagesHandler := apiutils.NewAuthenticatedStageHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		if !mocking {
			stages, err := stagesHandler.GetAllStages(*stageParameter.project)
			if err != nil {
				return fmt.Errorf("Failed to retrieve stages for project %s: %v", *stageParameter.project, err)
			}

			if len(stages) == 0 {
				fmt.Println("No stages found")
				return nil
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)
			fmt.Fprintln(w, "NAME\tCREATION DATE")

			for _, stage := range stages {
				if len(args) == 1 && stage.StageName == args[0] || len(args) == 0 {
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
