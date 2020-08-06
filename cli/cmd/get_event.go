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
	"github.com/ghodss/yaml"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
	"log"
)

type getEventStruct struct {
	KeptnContext 	*string
	Project			*string
	Stage			*string
	Service			*string
	Output			*string
}

var getEvent getEventStruct

// getEventCmd represents the get command
var getEventCmd = &cobra.Command{
	Use:     "event [eventType]",
	Aliases: []string{"events"},
	Short:   `Returns the latest Keptn event specified by event type`,
	Long:    `Returns the latest Keptn event specified by event type. The event type is defined here: https://github.com/keptn/spec/blob/0.1.4/cloudevents.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("please provide an event type as an argument")
		}

		eventType := args[0]

		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()

		if err != nil {
			log.Fatal(err)
		}

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

		if !mocking {
			events, err := eventHandler.GetEvents(&apiutils.EventFilter{
				KeptnContext: *getEvent.KeptnContext,
				Service: *getEvent.Service,
				Stage: *getEvent.Stage,
				Project: *getEvent.Project,
				EventType:    eventType,
			})

			if err != nil {
				log.Fatal(err)
			}

			if len(events) == 0 {
				logging.PrintLog("No event returned", logging.QuietLevel)
				return nil
			}

			for _, event := range events {
				if *getEvent.Output == "yaml" {
					event, _ := yaml.Marshal(event)
					fmt.Println(string(event))
				} else {
					event, _ := json.MarshalIndent(event, "", "    ")
					fmt.Println(string(event))
				}
			}
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getEventCmd)

	getEvent.KeptnContext = getEventCmd.Flags().StringP("keptn-context", "", "",
		"The ID of a Keptn context from which to retrieve the event")

	getEvent.Project = getEventCmd.Flags().StringP("project", "", "",
		"The Keptn project name from which to retrieve the event")

	getEvent.Stage = getEventCmd.Flags().StringP("stage", "", "",
		"The name of a stage within a project from which to retrieve the event")

	getEvent.Service = getEventCmd.Flags().StringP("service", "", "",
		"The name of a service within a project from which to retrieve the event")

	getEvent.Output = getEventCmd.Flags().StringP("output", "o", "",
		" Output format. One of: json|yaml")
}
