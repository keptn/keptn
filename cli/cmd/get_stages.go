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

type getStagesStruct struct {
	project *string
}

var projectName getStagesStruct

// evaluationDoneCmd represents the evaluation-done command
var getStagesCmd = &cobra.Command{
	Use:   "stages",
	Short: "Get all stages of a given project",
	Long:  `Get all stages of a given project`,
	Example: `keptn get stages --project sockshop
NAME           CREATION DATE                 
dev            2020-04-06T14:35:40.210Z
staging        2020-04-06T14:37:45.210Z
production     2020-04-06T14:37:45.210Z
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		stagesHandler := apiutils.NewAuthenticatedStageHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)

		if !mocking {
			stages, err := stagesHandler.GetAllStages(*projectName.project)
			if err != nil {
				fmt.Println(err)
				return errors.New("Could not retrieve stages for project " + *projectName.project)
			}
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)
			fmt.Fprintln(w, "NAME\tCREATION DATE")
			for _, stage := range stages {
				/*creationDateInt64, err := strconv.ParseInt(stage.CreationDate, 10, 64)
				if err != nil {
					panic(err)
				}
				tm := time.Unix(0, creationDateInt64)*/

				fmt.Fprintln(w, stage.StageName+"\t n/a") // + "\t" + tm.Format("2006-01-02T15:04:05Z07:00"))

			}

			w.Flush()
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getStagesCmd)

	projectName.project = getStagesCmd.Flags().StringP("project", "", "",
		"keptn project name")

	getStagesCmd.MarkFlagRequired("project")
}
