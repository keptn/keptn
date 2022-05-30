package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/internal"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type sequenceStruct struct {
	Project  *string `json:"project"`
	Service  *string `json:"service"`
	Stage    *string `json:"stage"`
}

var sequence = sequenceStruct{}

// Shipyard can have multiple sequences with an arbitrary name (my-sequence-1, my-sequence-2) in a stage
// The sequence name is used to identify the sequence to be triggered

/*
apiVersion: "spec.keptn.sh/0.2.2"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "my-sequence-1"
          tasks:
            - name: "deployment"
            - name: "test"
            - name: "evaluation"
            - name: "release"

        - name: "my-sequence-2"
          tasks:
            - name: "deployment"
            - name: "test"
            - name: "evaluation"
            - name: "release"
*/

var triggerSequenceCmd = &cobra.Command{
	Use:     "sequence",
	Aliases: []string{"sequence"},
	Short:   "Triggers the execution of any sequence in a project",
	Long: `Triggers the execution of any sequence in a project with an arbitrary name.
The name of the sequence has to be provided as an argument to the command. The sequence name is used to identify the sequence to be triggered.
`,
	Example:      `keptn trigger sequence <sequence-name> --project=<project> --service=<service> --stage=<stage>`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return doTriggerSequence(sequence, args[0])
	},
}

func doTriggerSequence(sequenceInputData sequenceStruct, sequenceName string) error {
	var endPoint url.URL
	var apiToken string
	var err error
	if !mocking {
		endPoint, apiToken, err = credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
	} else {
		endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
		endPoint = *endPointPtr
		apiToken = os.Getenv("MOCK_API_TOKEN")
	}
	if err != nil {
		return errors.New(authErrorMsg)
	}

	logging.PrintLog("Triggering sequence "+sequenceName+" in project "+*sequenceInputData.Project+" in stage "+*sequenceInputData.Stage+" service "+*sequenceInputData.Service, logging.InfoLevel)

	api, err := internal.APIProvider(endPoint.String(), apiToken)
	if err != nil {
		return err
	}

	projectServices, err := api.ServicesV1().GetAllServices(*sequenceInputData.Project, *sequenceInputData.Stage)

	if err != nil {
		return fmt.Errorf("error while retrieving information for service %s: %s", *sequenceInputData.Service, err.Error())
	}
	if !ServiceInSlice(*sequenceInputData.Service, projectServices) {
		return fmt.Errorf("could not start sequence because service %s has not been found in project %s", *sequenceInputData.Service, *sequenceInputData.Project)
	}

	triggeredEvent := keptnv2.DeploymentTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: *sequenceInputData.Project,
			Stage:   *sequenceInputData.Stage,
			Service: *sequenceInputData.Service,
		},
	}

	sdkEvent := cloudevents.NewEvent()
	sdkEvent.SetID(uuid.New().String())
	sdkEvent.SetType(keptnv2.GetTriggeredEventType(*sequenceInputData.Stage + "." + sequenceName))
	sdkEvent.SetSource("https://github.com/keptn/keptn/cli#configuration-change")
	sdkEvent.SetDataContentType(cloudevents.ApplicationJSON)
	sdkEvent.SetData(cloudevents.ApplicationJSON, triggeredEvent)

	eventByte, err := sdkEvent.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal cloud event. %s", err.Error())
	}

	apiEvent := apimodels.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, &apiEvent)
	if err != nil {
		return fmt.Errorf("failed to map cloud event to API event model. %v", err)
	}

	apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	_, err2 := apiHandler.SendEvent(apiEvent)
	if err2 != nil {
		return fmt.Errorf("trigger sequence was unsuccessful. %s", *err2.Message)
	}

	return nil
}

func doTriggerSequencePreRunCheck(sequenceInputData sequenceStruct) error {
	if sequenceInputData.Project == nil {
		return fmt.Errorf("Project has to be provided")
	}

	if sequenceInputData.Service == nil {
		return fmt.Errorf("Service has to be provided")
	}

	if sequenceInputData.Stage == nil {
		return fmt.Errorf("Stage has to be provided")
	}

	return nil
}

func init() {
	triggerCmd.AddCommand(triggerSequenceCmd)
	sequence.Project = triggerSequenceCmd.Flags().StringP("project", "", "",
		"The project containing the service for which the new artifact will be triggered")
	triggerSequenceCmd.MarkFlagRequired("project")

	sequence.Service = triggerSequenceCmd.Flags().StringP("service", "", "",
		"The service for which the new artifact will be triggered")
	triggerSequenceCmd.MarkFlagRequired("service")

	sequence.Stage = triggerSequenceCmd.Flags().StringP("stage", "", "",
		"The stage in which the new artifact will be triggered")
	triggerSequenceCmd.MarkFlagRequired("stage")

}
