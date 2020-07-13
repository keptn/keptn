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

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type getServiceStruct struct {
	project      *string
	outputFormat *string
}

var getService getServiceStruct

// evaluationDoneCmd represents the evaluation-done command
var getServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Get service details",
	Long:    `Get all services or details for a given service within a keptn project`,
	Example: `keptn get service carts --project=sockshop
NAME           CREATION DATE                 
carts          sockshop        2020-05-28T10:25:58+02:00

List all services in keptn
# keptn get services

# List all services in the sockshop project
keptn get services --project=sockshop

# Get details of the carts service in the sockshop project as json output
keptn get services carts --project=sockshop -o=json
	`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if *getService.outputFormat != "" {
			if *getService.outputFormat != "yaml" && *getService.outputFormat != "json" {
				return errors.New("Invalid output format, only yaml or json allowed")
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		var serviceName string
		if len(args) > 0 {
			serviceName = strings.Join(args, " ")
		} else {
			serviceName = ""
		}

		var projectName string
		projectName = *getService.project

		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		servicesHandler := apiutils.NewAuthenticatedServiceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

		if !mocking {

			projects, err := projectsHandler.GetAllProjects()
			if err != nil {
				fmt.Println(err)
				return errors.New(err.Error())
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 0, '\t', 0)

			if *getService.outputFormat == "" {
				fmt.Fprintln(w, "NAME\tPROJECT\tCREATION DATE")
			}

			for _, project := range projects {
				// make the project filter optional => print all services with matching names
				if projectName == "" || project.ProjectName == projectName {
					//stagesHandler := apiutils.NewStageHandler(endPoint.String())
					for _, stage := range project.Stages {
						services, err := servicesHandler.GetAllServices(project.ProjectName, stage.StageName)
						if err != nil {
							return errors.New(err.Error())
						}
						for _, service := range services {

							var returnService bool

							// return service if a serviceName is set and its matching
							if serviceName != "" && service.ServiceName == serviceName {
								returnService = true
							} else if serviceName == "" && project.ProjectName == projectName {
								// return all services of a project if no service name is defined
								returnService = true
							} else if serviceName == "" && projectName == "" {
								// print all services if no project or service is defined
								returnService = true
							} else {
								returnService = false
							}

							if returnService {

								if strings.ToLower(*getService.outputFormat) == "yaml" {
									yamlBytes, err := yaml.Marshal(service)
									if err != nil {
										return errors.New(err.Error())
									}
									yamlString := string(yamlBytes)
									fmt.Println(yamlString)
								} else if strings.ToLower(*getService.outputFormat) == "json" {
									jsonBytes, err := json.MarshalIndent(service, "", "   ")
									if err != nil {
										return errors.New(err.Error())
									}

									jsonString := string(jsonBytes)
									fmt.Println(jsonString)
								} else {

									creationDateInt64, err := strconv.ParseInt(service.CreationDate, 10, 64)
									if err != nil {
										panic(err)
									}
									tm := time.Unix(0, creationDateInt64)

									fmt.Fprintln(w, service.ServiceName+"\t"+project.ProjectName+"\t"+tm.Format("2006-01-02T15:04:05Z07:00"))
								}
							}
						}
					}
				}

			}
			w.Flush()

		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getServiceCmd)

	getService.project = getServiceCmd.Flags().StringP("project", "", "",
		"keptn project name")
	getService.outputFormat = getServiceCmd.Flags().StringP("output", "o", "",
		"Output format. One of json|yaml")

	// getServiceCmd.MarkFlagRequired("project")
}
