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
	"strings"
	"text/tabwriter"

	"github.com/keptn/go-utils/pkg/api/models"

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

// getServiceCmd represents the get service command
var getServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Get service details",
	Long:    `Get all services or details for a given service within a Keptn project`,
	Example: `keptn get service carts --project=sockshop
NAME           CREATION DATE                 
carts          sockshop        2020-05-28T10:25:58+02:00

keptn get services                                   # List all services in keptn

keptn get services --project=sockshop                # List all services in the sockshop project

keptn get services carts --project=sockshop -o=json  # Get details of the carts service in the sockshop project as json output
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
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
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		servicesHandler := apiutils.NewAuthenticatedServiceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

		if !mocking {
			projects, err := projectsHandler.GetAllProjects()
			if err != nil {
				return err
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 10, 8, 2, '\t', 0)

			if *getService.outputFormat == "" {
				fmt.Fprintln(w, "NAME\tPROJECT\tSTAGE\tCREATION DATE")
			}

			filteredProjects := filterProjects(projects, *getService.project)
			if len(filteredProjects) == 0 {
				if *getService.project != "" {
					fmt.Printf("No project %s found\n", *getService.project)
				} else {
					fmt.Println("No projects found")
				}
				return nil
			}

			for _, project := range filteredProjects {
				//stagesHandler := apiutils.NewStageHandler(endPoint.String())
				for _, stage := range project.Stages {
					services, err := servicesHandler.GetAllServices(project.ProjectName, stage.StageName)
					if err != nil {
						return errors.New(err.Error())
					}
					filteredServices := filterServices(services, getServiceNameFromArgs(args))
					if len(filteredServices) == 0 {
						if len(args) == 1 {
							fmt.Printf("No services %s found in project %s", args[0], *getService.project)
						} else {
							fmt.Printf("No services found in project %s", *getService.project)
						}
						return nil
					}

					for _, service := range filteredServices {
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
							fmt.Fprintln(w, service.ServiceName+"\t"+project.ProjectName+"\t"+stage.StageName+"\t"+parseCreationDate(service.CreationDate))
						}
					}
				}

			}
			w.Flush()

		}
		return nil
	},
}

func getServiceNameFromArgs(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return ""
}

func filterServices(services []*models.Service, serviceName string) []*models.Service {
	filteredServices := make([]*models.Service, 0)
	for _, service := range services {
		if serviceName == "" || serviceName == service.ServiceName {
			filteredServices = append(filteredServices, service)
		}
	}
	return filteredServices
}

func init() {
	getCmd.AddCommand(getServiceCmd)

	getService.project = getServiceCmd.Flags().StringP("project", "", "",
		"keptn project name")
	getService.outputFormat = getServiceCmd.Flags().StringP("output", "o", "",
		"Output format. One of json|yaml")
}
