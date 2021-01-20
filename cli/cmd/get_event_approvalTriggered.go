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
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
	"os"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type approvalTriggeredStruct struct {
	Project *string `json:"project"`
	Stage   *string `json:"stage"`
	Service *string `json:"service"`
}

var approvalTriggered approvalTriggeredStruct

// approvalTriggeredCmd represents the approval.triggered command
var approvalTriggeredCmd = &cobra.Command{
	Use:          "approval.triggered",
	Short:        "Returns the latest Keptn sh.keptn.event.approval.triggered event from a specific project/stage/service",
	Long:         `Returns the latest Keptn sh.keptn.event.approval.triggered event from a specific project/stage/service.`,
	Example:      `keptn get event approval.triggered --project=sockshop --stage=staging --service=carts`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getApprovalTriggeredEvents(approvalTriggered)
	},
}

func getApprovalTriggeredEvents(approvalTriggered approvalTriggeredStruct) error {
	var endPoint url.URL
	var apiToken string
	var err error
	if !mocking {
		endPoint, apiToken, err = credentialmanager.NewCredentialManager(false).GetCreds(namespace)
	} else {
		endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
		endPoint = *endPointPtr
		apiToken = ""
	}
	if err != nil {
		return errors.New(authErrorMsg)
	}

	logging.PrintLog("Starting to get approval.triggered events", logging.InfoLevel)

	if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
		return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
			endPointErr)
	}

	scHandler := apiutils.NewAuthenticatedShipyardControllerHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
	eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	if approvalTriggered.Service == nil || *approvalTriggered.Service == "" {
		return getAllApprovalEventsInStage(approvalTriggered, scHandler, eventHandler)
	}
	return getAllApprovalEventsInService(approvalTriggered, scHandler, eventHandler)
}

func getAllApprovalEventsInService(approvalTriggered approvalTriggeredStruct, scHandler *apiutils.ShipyardControllerHandler, eventHandler *apiutils.EventHandler) error {
	allEvents, err := scHandler.GetOpenTriggeredEvents(apiutils.EventFilter{
		Stage:     *approvalTriggered.Stage,
		Project:   *approvalTriggered.Project,
		Service:   *approvalTriggered.Service,
		EventType: keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName),
	})
	if err != nil {
		return err
	}

	printApprovalEvents(allEvents)
	return nil
}

func getAllApprovalEventsInStage(approvalTriggered approvalTriggeredStruct, scHandler *apiutils.ShipyardControllerHandler, eventHandler *apiutils.EventHandler) error {
	allEvents, err := scHandler.GetOpenTriggeredEvents(apiutils.EventFilter{
		Project:   *approvalTriggered.Project,
		Stage:     *approvalTriggered.Stage,
		EventType: keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName),
	})
	if err != nil {
		return err
	}

	printApprovalEvents(allEvents)
	return nil
}

func printApprovalEvents(allEvents []*apimodels.KeptnContextExtendedCE) {
	if len(allEvents) == 0 {
		logging.PrintLog("No approval.triggered events have been found", logging.InfoLevel)
	} else if len(allEvents) == 1 {
		eventsJSON, _ := json.MarshalIndent(allEvents[0], "", "	")
		fmt.Println(string(eventsJSON))
	} else {
		eventsJSON, _ := json.MarshalIndent(allEvents, "", "	")
		fmt.Println(string(eventsJSON))
	}
}

func retrieveApprovalEventsFromService(svc *apimodels.Service, project, stage string, eventHandler *apiutils.EventHandler, allEvents []*apimodels.KeptnContextExtendedCE) ([]*apimodels.KeptnContextExtendedCE, error) {
	for _, approval := range svc.OpenApprovals {
		events, err := eventHandler.GetEvents(&apiutils.EventFilter{
			Project: project,
			Stage:   stage,
			Service: svc.ServiceName,
			EventID: approval.EventID,
		})

		if err != nil {
			logging.PrintLog("Get approval.triggered event was unsuccessful", logging.InfoLevel)
			return nil, fmt.Errorf("%s", *err.Message)
		}

		if events != nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func init() {
	getEventCmd.AddCommand(approvalTriggeredCmd)

	approvalTriggered.Project = approvalTriggeredCmd.Flags().StringP("project", "p", "",
		"The name of a project from which to retrieve an approval.triggered event")
	approvalTriggeredCmd.MarkFlagRequired("project")
	approvalTriggered.Stage = approvalTriggeredCmd.Flags().StringP("stage", "s", "",
		"The name of a stage within a project from which to retrieve an approval.triggered event")
	approvalTriggeredCmd.MarkFlagRequired("stage")
	approvalTriggered.Service = approvalTriggeredCmd.Flags().StringP("service", "", "",
		"The name of a service within a project from which to retrieve an approval.triggered event")
}
