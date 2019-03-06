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
	"io/ioutil"
	"os"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/knative/pkg/cloudevents"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var project *string
var deploymentFilePath *string
var valuesFilePath *string
var serviceFilePath *string
var manifestFilePath *string

type serviceData map[string]interface{}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Onboards a new service.",
	Long: `Onboards a new service. Therefore, this command takes a 
	service description given as yaml and onboards 
	this service in the provided project. 
	Optionally, this command takes a Helm deployment and service description.
	Usage of \"onboard service\":

keptn onboard service --project=carts --values=values.yaml --deployment=deployment.yaml --service=service.yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if (*valuesFilePath != "" || *manifestFilePath != "") && (*valuesFilePath != "" && *manifestFilePath != "") {
			return errors.New("Error specifying a Helm description as well as a k8s manifest. Only use one option")
		}

		if *valuesFilePath != "" {
			valuesData, err := readFile(*valuesFilePath)
			if err != nil {
				return err
			}

			if _, err := unmarshalValues(valuesData); err != nil {
				return err
			}

			if *deploymentFilePath != "" {
				if _, err := readFile(*deploymentFilePath); err != nil {
					return err
				}
			}

			if *serviceFilePath != "" {
				if _, err := readFile(*serviceFilePath); err != nil {
					return err
				}
			}
		} else {
			if *deploymentFilePath != "" {
				fmt.Printf("The specified deployment file is ignored")
			}
			if *serviceFilePath != "" {
				fmt.Println("The specified service file is ignored")
			}
			manifestData, err := readFile(*manifestFilePath)
			if err != nil {
				return err
			}

			if _, err := unmarshalValues(manifestData); err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to onboard service")

		svcData := serviceData{}
		svcData["project"] = *project

		if *valuesFilePath != "" {
			valuesData, err := readFile(*valuesFilePath)
			if err != nil {
				return err
			}
			data, err := unmarshalValues(valuesData)
			if err != nil {
				return err
			}
			svcData["values"] = data

			deployment := make(map[string]string)

			if *deploymentFilePath != "" {
				content, err := readFile(*deploymentFilePath)
				if err != nil {
					return err
				}
				deployment["deployment"] = content
			}

			if *serviceFilePath != "" {
				content, err := readFile(*serviceFilePath)
				if err != nil {
					return err
				}
				deployment["service"] = content
			}

			if len(deployment) > 0 {
				svcData["templates"] = deployment
			}
		} else {
			manifestData, err := readFile(*manifestFilePath)
			if err != nil {
				return err
			}
			data, err := unmarshalValues(manifestData)
			if err != nil {
				return err
			}
			svcData["manifest"] = data
		}

		builder := cloudevents.Builder{
			Source:    "https://github.com/keptn/keptn/cli#onboardservice",
			EventType: "onboard.service",
			Encoding:  cloudevents.StructuredV01,
		}
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil || endPoint == "" {
			fmt.Printf("onboard service called without beeing authenticated.")
			return errors.New("This command requires to be authenticated. See \"keptn auth\" for details")
		}

		req, err := builder.Build(endPoint+"service", svcData)
		if err != nil {
			return err
		}

		err = utils.Send(req, apiToken)
		if err != nil {
			fmt.Printf("onboard service command was unsuccessful. Details: %v", err)
			return err
		}
		fmt.Printf("Successfully onboarded service in project %v\n", svcData["project"])
		return nil
	},
}

func readFile(fileName string) (string, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", fmt.Errorf("Cannot find file %s", fileName)
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshalValues(data string) (interface{}, error) {
	var body interface{}
	err := yaml.Unmarshal([]byte(data), &body)
	if err != nil {
		return nil, err
	}
	return convert(body), nil
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func init() {
	onboardCmd.AddCommand(serviceCmd)

	project = serviceCmd.Flags().StringP("project", "p", "", "The name of the project")
	serviceCmd.MarkFlagRequired("project")

	// Flags for onboarding a service using Helm
	valuesFilePath = serviceCmd.Flags().StringP("values", "v", "", "The values file")
	deploymentFilePath = serviceCmd.Flags().StringP("deployment", "d", "", "The deployment file")
	serviceFilePath = serviceCmd.Flags().StringP("service", "s", "", "The service file")

	// Flags for onboarding using "pure" k8s manifests
	manifestFilePath = serviceCmd.Flags().StringP("manifest", "m", "", "The manifest file")
}
