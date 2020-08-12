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
	"net/url"
	"os"
)

type GetEventStruct struct {
	KeptnContext 	*string
	Project			*string
	Stage			*string
	Service			*string
	PageSize		*string
	Output			*string
	NumOfPages		*int
}

var getEventParams GetEventStruct

// getEventCmd represents the get command
var getEventCmd = &cobra.Command{
	Use:     "event [eventType]",
	Aliases: []string{"events"},
	Short:   `Returns the latest Keptn event specified by event type`,
	Long:    `Returns the latest Keptn event specified by event type. The event type is defined here: https://github.com/keptn/spec/blob/0.1.4/cloudevents.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getEvent(getEventParams, args)
	},
}

func getEvent(eventStruct GetEventStruct, args []string) error {
	if len(args) == 0 {
		return errors.New("please provide an event type as an argument")
	}

	pageSize := setParameterValue(*eventStruct.PageSize, "1")

	eventType := args[0]

	var endPoint url.URL
	var apiToken string
	var err error
	if !mocking {
		endPoint, apiToken, err = credentialmanager.NewCredentialManager().GetCreds()
	} else {
		endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
		endPoint = *endPointPtr
		apiToken = ""
	}

	if err != nil {
		return errors.New(authErrorMsg)
	}

	eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

	events, modErr := eventHandler.GetEvents(&apiutils.EventFilter{
		KeptnContext: *eventStruct.KeptnContext,
		Service: *eventStruct.Service,
		Stage: *eventStruct.Stage,
		Project: *eventStruct.Project,
		EventType:    eventType,
		PageSize: pageSize,
		NumberOfPages: *eventStruct.NumOfPages,
	})

	if modErr != nil {
		log.Fatal(modErr)
	}

	if len(events) == 0 {
		logging.PrintLog("No event returned", logging.QuietLevel)
		return nil
	} else if len(events) == 1 {
		printEvents(events[0], *eventStruct.Output)
	} else {
		printEvents(events, *eventStruct.Output)
	}

	return nil
}

func printEvents(events interface{}, outputType string) {
	if outputType == "yaml" {
		eventsYAML, _ := yaml.Marshal(events)
		fmt.Println(string(eventsYAML))
	} else {
		eventsJSON, _ := json.MarshalIndent(events, "", "    ")
		fmt.Println(string(eventsJSON))
	}
}

func init() {
	getCmd.AddCommand(getEventCmd)

	getEventParams.KeptnContext = getEventCmd.Flags().StringP("keptn-context", "", "",
		"The ID of a Keptn context from which to retrieve the event")

	getEventParams.Project = getEventCmd.Flags().StringP("project", "", "",
		"The Keptn project name from which to retrieve the event")
	getEventCmd.MarkFlagRequired("project")

	getEventParams.Stage = getEventCmd.Flags().StringP("stage", "", "",
		"The name of a stage within a project from which to retrieve the event")

	getEventParams.Service = getEventCmd.Flags().StringP("service", "", "",
		"The name of a service within a project from which to retrieve the event")

	getEventParams.Output = getEventCmd.Flags().StringP("output", "o", "",
		" Output format. One of: json|yaml")

	getEventParams.PageSize = getEventCmd.Flags().StringP("page-size", "", "",
		"Max number of return events per page (Default 1)")

	getEventParams.NumOfPages = getEventCmd.Flags().IntP("num-of-pages", "", 1,
		"Number of pages that should be returned (Default 1).")
}

func setParameterValue(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}