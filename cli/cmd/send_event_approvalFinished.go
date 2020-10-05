package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

type sendApprovalFinishedStruct struct {
	Project *string `json:"project"`
	Stage   *string `json:"stage"`
	Service *string `json:"service"`
	ID      *string `json:"id"`
}

var sendApprovalFinishedOptions sendApprovalFinishedStruct

var approvalFinishedCmd = &cobra.Command{
	Use: "approval.finished",
	Short: "Sends an approval.finished event to Keptn in order to confirm an open approval " +
		"with the specified ID in the provided project and stage",
	Long: `Sends an approval.finished event to Keptn in order to confirm an open approval with the specified ID in the provided project and stage. 

* This command takes the project (*--project*) and stage (*--stage*). 
* It is optional to specify the ID (*--id*) of the corresponding approval.triggered event. If the ID is not provided, the command asks the user which open approval should be accepted or declined.
* The open approval.triggered events and their ID can be retrieved using the "keptn get event approval.triggered --project=<project> --stage=<stage>" command.
`,
	Example: `keptn send event approval.finished --project=sockshop --stage=hardening --id=1234-5678-9123`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if *sendApprovalFinishedOptions.ID == "" && *sendApprovalFinishedOptions.Service == "" {
			logging.PrintLog("Either ID or service must be provided", logging.InfoLevel)
			return errors.New("either ID or service must be provided")
		} else if *sendApprovalFinishedOptions.ID != "" && *sendApprovalFinishedOptions.Service != "" {
			logging.PrintLog("Either ID or service must be provided", logging.InfoLevel)
			return errors.New("either ID or service must be provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendApprovalFinishedEvent(sendApprovalFinishedOptions)
	},
	SilenceUsage: true,
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

	if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
		return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
			endPointErr)
	}

	apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme,
		ClientCertPath, ClientKeyPath, RootCertPath)
	eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme,
		ClientCertPath, ClientKeyPath, RootCertPath)

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	var keptnContext string
	var triggeredID string
	var approvalFinishedEvent *keptnevents.ApprovalFinishedEventData

	if *sendApprovalFinishedOptions.ID != "" {
		keptnContext, triggeredID, approvalFinishedEvent, err = getApprovalFinishedForID(eventHandler, sendApprovalFinishedOptions)
	} else if *sendApprovalFinishedOptions.Service != "" {
		serviceHandler := apiutils.NewAuthenticatedServiceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme,
			ClientCertPath, ClientKeyPath, RootCertPath)
		keptnContext, triggeredID, approvalFinishedEvent, err = getApprovalFinishedForService(eventHandler,
			serviceHandler, sendApprovalFinishedOptions)
	}
	if err != nil {
		return err
	}

	if approvalFinishedEvent == nil {
		return nil
	}

	ID := uuid.New().String()
	source, _ := url.Parse("https://github.com/keptn/keptn/cli#approval.finished")

	sdkEvent := cloudevents.NewEvent()
	sdkEvent.SetID(ID)
	sdkEvent.SetType(keptnevents.ApprovalFinishedEventType)
	sdkEvent.SetSource(source.String())
	sdkEvent.SetDataContentType(cloudevents.ApplicationJSON)
	sdkEvent.SetExtension("shkeptncontext", keptnContext)
	sdkEvent.SetExtension("triggeredid", triggeredID)
	sdkEvent.SetData(cloudevents.ApplicationJSON, approvalFinishedEvent)

	eventByte, err := json.Marshal(sdkEvent)
	if err != nil {
		return fmt.Errorf("Failed to marshal cloud event. %s", err.Error())
	}

	apiEvent := apimodels.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, &apiEvent)
	if err != nil {
		return fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
	}

	responseEvent, errorObj := apiHandler.SendEvent(apiEvent)
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

func getApprovalFinishedForService(eventHandler *apiutils.EventHandler, serviceHandler *apiutils.ServiceHandler,
	approvalFinishedOptions sendApprovalFinishedStruct) (string, string, *keptnevents.ApprovalFinishedEventData, error) {
	svc, err := serviceHandler.GetService(*approvalFinishedOptions.Project, *approvalFinishedOptions.Stage, *approvalFinishedOptions.Service)
	if err != nil {
		logging.PrintLog("Open approval.triggered event for service "+*approvalFinishedOptions.Service+" could not be retrieved: "+err.Error(), logging.InfoLevel)
		return "", "", nil, err
	}
	if svc == nil {
		logging.PrintLog("Service "+*approvalFinishedOptions.Service+" could not be found", logging.InfoLevel)
		return "", "", nil, nil
	}

	if len(svc.OpenApprovals) == 0 {
		logging.PrintLog("No open approval.triggered event for service "+*approvalFinishedOptions.Service+" has been found", logging.InfoLevel)
		return "", "", nil, nil
	}

	// print all available options
	printApprovalOptions(svc.OpenApprovals, eventHandler, approvalFinishedOptions)

	// select option
	nrOfOptions := len(svc.OpenApprovals)
	selectedOption, err := selectApprovalOption(nrOfOptions)
	if err != nil {
		return "", "", nil, err
	}

	index := selectedOption - 1
	eventToBeApproved := svc.OpenApprovals[index]

	// approve or decline?
	approve := approveOrDecline()

	events, errorObj := eventHandler.GetEvents(&apiutils.EventFilter{
		Project:   *approvalFinishedOptions.Project,
		Stage:     *approvalFinishedOptions.Stage,
		EventType: keptnevents.ApprovalTriggeredEventType,
		EventID:   eventToBeApproved.EventID,
	})

	if errorObj != nil {
		logging.PrintLog("Cannot retrieve approval.triggered event with ID "+*approvalFinishedOptions.ID+": "+*errorObj.Message, logging.InfoLevel)
		return "", "", nil, errors.New(*errorObj.Message)
	}

	if len(events) == 0 {
		logging.PrintLog("No open approval.triggered event with the ID "+*approvalFinishedOptions.ID+" has been found", logging.InfoLevel)
		return "", "", nil, nil
	}

	approvalTriggeredEvent := &keptnevents.ApprovalTriggeredEventData{}

	err = mapstructure.Decode(events[0].Data, approvalTriggeredEvent)
	if err != nil {
		logging.PrintLog("Cannot decode approval.triggered event: "+err.Error(), logging.InfoLevel)
		return "", "", nil, err
	}

	var approvalResult string
	if approve {
		approvalResult = "pass"
	} else {
		approvalResult = "failed"
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
			Result: approvalResult,
			Status: "succeeded",
		},
	}

	return eventToBeApproved.KeptnContext, eventToBeApproved.EventID, approvalFinishedEvent, nil
}

func approveOrDecline() bool {
	var approve bool
	keepAsking := true
	for keepAsking {
		logging.PrintLog("Do you want to (a)pprove or (d)ecline: ", logging.InfoLevel)
		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			logging.PrintLog("Invalid option. Please enter either 'a' to approve, or 'd' to decline", logging.InfoLevel)
		}
		in = strings.TrimSpace(in)
		if in != "a" && in != "d" {
			logging.PrintLog("Invalid option. Please enter either 'a' to approve, or 'd' to decline", logging.InfoLevel)
		} else {
			keepAsking = false
		}
		if in == "a" {
			approve = true
		} else if in == "d" {
			approve = false
		}
	}

	return approve
}

func selectApprovalOption(nrOfOptions int) (int, error) {
	var selectedOption int

	keepAsking := true
	for keepAsking {
		logging.PrintLog("Select the option to approve or decline: ", logging.InfoLevel)
		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			logging.PrintLog(fmt.Sprintf("Invalid option. Please enter a value between 1 and %d", nrOfOptions), logging.InfoLevel)
		}
		in = strings.TrimSpace(in)
		selectedOption, err = strconv.Atoi(in)

		if err != nil || selectedOption < 1 || selectedOption > nrOfOptions {
			logging.PrintLog(fmt.Sprintf("Invalid option. Please enter a value between 1 and %d", nrOfOptions), logging.InfoLevel)
		} else {
			keepAsking = false
		}
	}
	return selectedOption, nil
}

func printApprovalOptions(approvals []*apimodels.Approval, eventHandler *apiutils.EventHandler, approvalFinishedOptions sendApprovalFinishedStruct) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "OPTION", "VERSION", "EVALUATION")

	for index, approval := range approvals {
		score := getScoreForApprovalTriggeredEvent(eventHandler, approvalFinishedOptions, approval)

		appendOptionToWriter(w, index, approval, score)
	}
	fmt.Fprintf(w, "\n")
}

func appendOptionToWriter(w *tabwriter.Writer, index int, approval *apimodels.Approval, score string) {
	fmt.Fprintf(w, "\n (%d)\t%s\t%s\t", index+1, approval.Tag, score)
}

func getScoreForApprovalTriggeredEvent(eventHandler *apiutils.EventHandler, approvalFinishedOptions sendApprovalFinishedStruct, approval *apimodels.Approval) string {
	score := "n/a"
	evaluationDoneEvents, errorObj := eventHandler.GetEvents(&apiutils.EventFilter{
		Project:      *approvalFinishedOptions.Project,
		Stage:        *approvalFinishedOptions.Stage,
		Service:      *approvalFinishedOptions.Service,
		EventType:    keptnevents.EvaluationDoneEventType,
		KeptnContext: approval.KeptnContext,
	})
	if errorObj != nil {
		return score
	}
	if len(evaluationDoneEvents) == 0 {
		return score
	}
	evaluationDoneData := &keptnevents.EvaluationDoneEventData{}

	err := mapstructure.Decode(evaluationDoneEvents[0].Data, evaluationDoneData)
	if err != nil {
		return score
	}

	if evaluationDoneData.EvaluationDetails != nil {
		score = fmt.Sprintf("%f", evaluationDoneData.EvaluationDetails.Score)
	}
	return score
}

func getApprovalFinishedForID(eventHandler *apiutils.EventHandler, sendApprovalFinishedOptions sendApprovalFinishedStruct) (string,
	string, *keptnevents.ApprovalFinishedEventData, error) {
	events, errorObj := eventHandler.GetEvents(&apiutils.EventFilter{
		Project:   *sendApprovalFinishedOptions.Project,
		Stage:     *sendApprovalFinishedOptions.Stage,
		EventType: keptnevents.ApprovalTriggeredEventType,
		EventID:   *sendApprovalFinishedOptions.ID,
	})

	if errorObj != nil {
		logging.PrintLog("Cannot retrieve approval.triggered event with ID "+*sendApprovalFinishedOptions.ID+": "+*errorObj.Message, logging.InfoLevel)
		return "", "", nil, errors.New(*errorObj.Message)
	}

	if len(events) == 0 {
		logging.PrintLog("No open approval.triggered event with the ID "+*sendApprovalFinishedOptions.ID+" has been found", logging.InfoLevel)
		return "", "", nil, nil
	}

	approvalTriggeredEvent := &keptnevents.ApprovalTriggeredEventData{}

	err := mapstructure.Decode(events[0].Data, approvalTriggeredEvent)
	if err != nil {
		logging.PrintLog("Cannot decode approval.triggered event: "+err.Error(), logging.InfoLevel)
		return "", "", nil, err
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
			Result: "pass",
			Status: "succeeded",
		},
	}
	return events[0].Shkeptncontext, events[0].ID, approvalFinishedEvent, nil
}

func init() {
	sendEventCmd.AddCommand(approvalFinishedCmd)

	sendApprovalFinishedOptions.Project = approvalFinishedCmd.Flags().StringP("project", "", "",
		"The project containing the service to be approved")
	approvalFinishedCmd.MarkFlagRequired("project")

	sendApprovalFinishedOptions.Stage = approvalFinishedCmd.Flags().StringP("stage", "", "",
		"The stage containing the service to be approved")
	approvalFinishedCmd.MarkFlagRequired("stage")

	sendApprovalFinishedOptions.Service = approvalFinishedCmd.Flags().StringP("service", "", "",
		"The service to be approved")

	sendApprovalFinishedOptions.ID = approvalFinishedCmd.Flags().StringP("id", "", "",
		"The ID of the approval.triggered event to be approved")
	// approvalFinishedCmd.MarkFlagRequired("id")

}
