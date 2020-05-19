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
	"text/tabwriter"
	"time"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
)

/*type evaluationDoneStruct struct {
	KeptnContext *string `json:"keptnContext"`
}

var evaluationDone evaluationDoneStruct
*/
// evaluationDoneCmd represents the evaluation-done command
var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get all keptn projects",
	Long:  `Get a list with all keptn projects`,
	Example: `keptn get projects
NAME           CREATION DATE                 
sockshop       2020-04-06T14:35:40.210Z
musicshop      2020-04-06T14:37:45.210Z	
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		//fmt.Println(endPoint.String())

		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)

		if !mocking {

			projects, err := projectsHandler.GetAllProjects()
			if err != nil {
				fmt.Println(err)
				return errors.New(authErrorMsg)
			}
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)
			fmt.Fprintln(w, "NAME\tCREATION DATE")
			for _, project := range projects {
				creationDateInt64, err := strconv.ParseInt(project.CreationDate, 10, 64)
				if err != nil {
					panic(err)
				}
				tm := time.Unix(0, creationDateInt64)
				fmt.Fprintln(w, project.ProjectName+"\t"+tm.Format("2006-01-02T15:04:05Z07:00"))
			}
			w.Flush()
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getProjectsCmd)
}
