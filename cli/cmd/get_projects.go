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

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	"github.com/keptn/keptn/cli/pkg/logging"

	"github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
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
	Short:   "Get details of a Keptn project",
	Long:    `Get all projects or details of a given Keptn project`,
	Example: `keptn get projects
NAME            CREATION DATE
sockshop        2020-05-28T10:25:50+02:00
shirtshop       2020-05-28T10:26:22+02:00
	
keptn get project sockshop
NAME           CREATION DATE                 
sockshop       2020-04-06T14:35:40.210Z

keptn get project sockshop -output=yaml  # Returns project details in YAML format

keptn get project sockshop -output=json  # Returns project details in JSON format
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if *getProject.outputFormat != "" {
			if *getProject.outputFormat != "yaml" && *getProject.outputFormat != "json" {
				return errors.New("Invalid output format, only yaml or json allowed")
			}
		}

		if len(args) == 1 {
			if !keptncommon.ValidateKeptnEntityName(args[0]) {
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
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

		if !mocking {

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 2, '\t', 0)
			if *getProject.outputFormat == "" {
				fmt.Fprintln(w, "NAME\tCREATION DATE\tSHIPYARD VERSION")
			}

			projects, err := projectsHandler.GetAllProjects()
			if err != nil {
				return err
			}

			filteredProjects := filterProjects(projects, getProjectNameFromArgs(args))
			if len(filteredProjects) == 0 {
				if len(args) == 1 {
					fmt.Printf("No project %s found\n", args[0])
				} else {
					fmt.Println("No projects found")
				}
				return nil
			}

			for _, project := range filteredProjects {
				if strings.ToLower(*getProject.outputFormat) == "yaml" {
					yamlBytes, err := yaml.Marshal(project)
					if err != nil {
						return errors.New(err.Error())
					}
					fmt.Println(string(yamlBytes))
				} else if strings.ToLower(*getProject.outputFormat) == "json" {
					jsonBytes, err := json.MarshalIndent(project, "", "   ")
					if err != nil {
						return errors.New(err.Error())
					}
					fmt.Println(string(jsonBytes))
				} else {
					fmt.Fprintln(w, project.ProjectName+"\t"+parseCreationDate(project.CreationDate)+"\t"+project.ShipyardVersion)
				}
			}
			w.Flush()
		}
		return nil
	},
}

func getProjectNameFromArgs(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return ""
}

func filterProjects(projects []*models.Project, projectName string) []*models.Project {
	filteredProjects := make([]*models.Project, 0)
	for _, project := range projects {
		if projectName == "" || projectName == project.ProjectName {
			filteredProjects = append(filteredProjects, project)
		}
	}
	return filteredProjects
}

func parseCreationDate(creationDate string) string {

	const na = "n/a"
	if creationDate == "" {
		return na
	}
	creationDateInt64, err := strconv.ParseInt(creationDate, 10, 64)
	if err != nil {
		logging.PrintLog("Failed to parse Creation Date", logging.InfoLevel)
		return na
	}
	tm := time.Unix(0, creationDateInt64)
	return tm.Format("2006-01-02T15:04:05Z07:00")
}

func init() {

	getCmd.AddCommand(getProjectCmd)
	getProject.outputFormat = getProjectCmd.Flags().StringP("output", "o", "",
		"Output format. One of json|yaml")

}
