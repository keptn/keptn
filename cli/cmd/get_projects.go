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
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type getProjectStruct struct {
	outputFormat *string
}

var getProject getProjectStruct

var getProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"projects"},
	Short:   "Get details of a keptn project",
	Long:    `Get all projects or details of a given keptn project`,
	Example: `keptn get projects
NAME            CREATION DATE
sockshop        2020-05-28T10:25:50+02:00
shirtshop       2020-05-28T10:26:22+02:00
	
keptn get project sockshop
NAME           CREATION DATE                 
sockshop       2020-04-06T14:35:40.210Z

# Returns project details in YAML format
keptn get project sockshop -output=yaml

# Returns project details in JSON format
keptn get project sockshop -output=json
	`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if *getProject.outputFormat != "" {
			if *getProject.outputFormat != "yaml" && *getProject.outputFormat != "json" {
				return errors.New("Invalid output format, only yaml or json allowed")
			}
		}

		if len(args) == 1 {
			if !keptn.ValidateKeptnEntityName(args[0]) {
				errorMsg := "Project name contains upper case letter(s) or special character(s).\n"
				errorMsg += "Keptn relies on the following conventions: "
				errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
				errorMsg += "Please update project name and try again."
				return errors.New(errorMsg)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		var projectName string
		projectName = strings.Join(args, " ")

		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)

		if !mocking {

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)
			if *getProject.outputFormat == "" {
				fmt.Fprintln(w, "NAME\tCREATION DATE")
			}

			var prj models.Project
			prj.ProjectName = projectName

			projects, err := projectsHandler.GetAllProjects()
			if err != nil {
				fmt.Println(err)
				return errors.New(authErrorMsg)
			}

			if projectName != "" {
				for _, project := range projects {

					if project.ProjectName == projectName {

						if strings.ToLower(*getProject.outputFormat) == "yaml" {
							yamlBytes, err := yaml.Marshal(project)
							if err != nil {
								return errors.New(err.Error())
							}
							yamlString := string(yamlBytes)
							fmt.Println(yamlString)

						} else if strings.ToLower(*getProject.outputFormat) == "json" {
							jsonBytes, err := json.MarshalIndent(project, "", "   ")
							if err != nil {
								return errors.New(err.Error())
							}

							jsonString := string(jsonBytes)
							fmt.Println(jsonString)
						} else {
							creationDateInt64, err := strconv.ParseInt(project.CreationDate, 10, 64)
							if err != nil {
								panic(err)
							}
							tm := time.Unix(0, creationDateInt64)
							fmt.Fprintln(w, project.ProjectName+"\t"+tm.Format("2006-01-02T15:04:05Z07:00"))
						}
					}
				}
			} else {
				for _, project := range projects {
					if strings.ToLower(*getProject.outputFormat) == "yaml" {
						yamlBytes, err := yaml.Marshal(project)
						if err != nil {
							return errors.New(err.Error())
						}
						yamlString := string(yamlBytes)
						fmt.Println(yamlString)

					} else if strings.ToLower(*getProject.outputFormat) == "json" {
						jsonBytes, err := json.MarshalIndent(project, "", "   ")
						if err != nil {
							return errors.New(err.Error())
						}

						jsonString := string(jsonBytes)
						fmt.Println(jsonString)
					} else {
						creationDateInt64, err := strconv.ParseInt(project.CreationDate, 10, 64)
						if err != nil {
							panic(err)
						}
						tm := time.Unix(0, creationDateInt64)
						fmt.Fprintln(w, project.ProjectName+"\t"+tm.Format("2006-01-02T15:04:05Z07:00"))
					}
				}
			}
			w.Flush()
		}
		return nil
	},
}

func init() {

	getCmd.AddCommand(getProjectCmd)
	getProject.outputFormat = getProjectCmd.Flags().StringP("output", "o", "",
		"Output format. One of json|yaml")

}
