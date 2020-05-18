package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

type sendApprovalFinishedStruct struct {
	Project *string `json:"project"`
	Stage   *string `json:"stage"`
	ID      *string `json:"id"`
}

var sendApprovalFinishedOptions sendApprovalFinishedStruct

var approvalFinishedCmd = &cobra.Command{
	Use: "approval.finished",
	Short: "Sends an approval.finished event to Keptn in order to confirm an open approval " +
		"with the specified ID in the provided project and stage",
	Long: `Sends an approval.finished event to Keptn in order to confirm an open approval
with the specified ID in the provided project and stage. 

This command takes the project (*--project*), stage (*--stage*). Besides, it is necessary to specify the ID (*--id*) of the corresponding approval.triggered event.
The open approval.triggered events and their ID can be retrieved using the "keptn get event approval.triggered --project=<project> --stage=<stage>"" command."
`,
	Example:      `keptn send event approval.finished --project=sockshop --stage=hardening --id=1234-5678-9123`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendApprovalFinishedEvent(sendApprovalFinishedOptions)
	},
}

func sendApprovalFinishedEvent(sendApprovalFinishedOptions sendApprovalFinishedStruct) error {
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

	logging.PrintLog("Starting to send approval.finished event", logging.InfoLevel)

	eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	events, errorObj := eventHandler.GetEvents(&apiutils.EventFilter{
		Project:   *sendApprovalFinishedOptions.Project,
		Stage:     *sendApprovalFinishedOptions.Stage,
		EventType: keptnevents.ApprovalTriggeredEventType,
		EventID:   *sendApprovalFinishedOptions.ID,
	})

	if errorObj != nil {
		logging.PrintLog("Cannot retrieve approval.triggered event with ID "+*sendApprovalFinishedOptions.ID+": "+err.Error(), logging.InfoLevel)
		return errors.New(*errorObj.Message)
	}

	if len(events) == 0 {
		logging.PrintLog("No open approval.triggered event with the ID "+*sendApprovalFinishedOptions.ID+" has been found", logging.InfoLevel)
		return nil
	}

	approvalTriggeredEvent := &keptnevents.ApprovalTriggeredEventData{}

	err = mapstructure.Decode(events[0].Data, approvalTriggeredEvent)
	if err != nil {
		logging.PrintLog("Cannot decode approval.triggered event: "+err.Error(), logging.InfoLevel)
		return err
	}

	approvalFinishedEvent := &keptnevents.ApprovalFinishedEventData{
		Project:            approvalTriggeredEvent.Project,
		Service:            approvalTriggeredEvent.Service,
		Stage:              approvalTriggeredEvent.Stage,
		TestStrategy:       approvalTriggeredEvent.TestStrategy,
		DeploymentStrategy: approvalTriggeredEvent.DeploymentStrategy,
		Tag:                approvalTriggeredEvent.Tag,
		Image:              approvalTriggeredEvent.Image,
		Labels:             approvalTriggeredEvent.Labels,
		Approval: keptnevents.ApprovalData{
			TriggeredID: events[0].ID,
			Result:      "pass",
			Status:      "succeeded",
		},
	}

	keptnContext := events[0].Shkeptncontext
	ID := uuid.New().String()
	source, _ := url.Parse("https://github.com/keptn/keptn/cli#approval.finished")
	contentType := "application/json"
	sdkEvent := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          ID,
			Type:        keptnevents.ApprovalFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnContext},
		}.AsV02(),
		Data: approvalFinishedEvent,
	}

	eventByte, err := sdkEvent.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal cloud event. %s", err.Error())
	}

	apiEvent := apimodels.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, &apiEvent)
	if err != nil {
		return fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
	}

	responseEvent, errorObj := eventHandler.SendEvent(apiEvent)
	if errorObj != nil {
		logging.PrintLog("Send approval.triggered was unsuccessful", logging.QuietLevel)
		return fmt.Errorf("Send approval.triggered was unsuccessful. %s", *errorObj.Message)
	}

	if responseEvent == nil {
		logging.PrintLog("No event returned", logging.QuietLevel)
		return nil
	}

	return nil
}

func init() {
	sendEventCmd.AddCommand(approvalFinishedCmd)

	sendApprovalFinishedOptions.Project = approvalFinishedCmd.Flags().StringP("project", "", "",
		"The project containing the service to be approved")
	approvalFinishedCmd.MarkFlagRequired("project")

	sendApprovalFinishedOptions.Stage = approvalFinishedCmd.Flags().StringP("stage", "", "",
		"The stage containing the service to be approved")
	approvalFinishedCmd.MarkFlagRequired("stage")

	sendApprovalFinishedOptions.ID = approvalFinishedCmd.Flags().StringP("id", "", "",
		"The ID of the approval.triggered event to be approved")
	approvalFinishedCmd.MarkFlagRequired("id")

}
