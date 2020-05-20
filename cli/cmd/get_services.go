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

// evaluationDoneCmd represents the evaluation-done command
var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Get all keptn services",
	Long:  `Get a list with all keptn services`,
	Example: `keptn get services
NAME           CREATION DATE                 
carts          2020-04-06T14:35:40.210Z
carts-db       2020-04-06T14:52:52.210Z
catalogue      2020-04-12T16:00:40.210Z
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		servicesHandler := apiutils.NewAuthenticatedServiceHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)

		if !mocking {
			projects, err := projectsHandler.GetAllProjects()
			if err != nil {
				fmt.Println(err)
				return errors.New(err.Error())
			}
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)

			fmt.Fprintln(w, "NAME\tPROJECT\tCREATION DATE")

			for _, project := range projects {

				for _, stage := range project.Stages {
					services, err := servicesHandler.GetAllServices(project.ProjectName, stage.StageName)
					if err != nil {
						return errors.New(err.Error())
					}
					for _, service := range services {

						creationDateInt64, err := strconv.ParseInt(service.CreationDate, 10, 64)
						if err != nil {
							panic(err)
						}
						tm := time.Unix(0, creationDateInt64)

						fmt.Fprintln(w, service.ServiceName+"\t"+project.ProjectName+"\t"+tm.Format("2006-01-02T15:04:05Z07:00"))
					}
				}
			}
			w.Flush()
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getServicesCmd)
}
