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

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/internal"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

type sendApprovalFinishedStruct struct {
	Project *string            `json:"project"`
	Stage   *string            `json:"stage"`
	Service *string            `json:"service"`
	ID      *string            `json:"id"`
	Labels  *map[string]string `json:"labels"`
}

var sendApprovalFinishedOptions sendApprovalFinishedStruct

var approvalFinishedCmd = &cobra.Command{
	Use:  "approval.finished",
	Args: cobra.NoArgs,
	Short: "Sends an approval.finished event to Keptn in order to confirm an open approval " +
		"with the specified ID in the provided project and stage",
	Long: `Sends an approval.finished event to Keptn in order to confirm an open approval with the specified ID in the provided project and stage. 

* This command takes the project (*--project*) and stage (*--stage*). 
* It is optional to specify the ID (*--id*) of the corresponding approval.triggered event. If the ID is not provided, the command asks the user which open approval should be accepted or declined.
* The open approval.triggered events and their IDs can be retrieved using the "keptn get event approval.triggered --project=<project> --stage=<stage>" command.
`,
	Example: `keptn send event approval.finished --project=sockshop --stage=hardening --id=1234-5678-9123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := deSendEventApprovalFinishedPreRunCheck(); err != nil {
			return err
		}
		return sendApprovalFinishedEvent(sendApprovalFinishedOptions)
	},
	SilenceUsage: true,
}

func deSendEventApprovalFinishedPreRunCheck() error {
	if *sendApprovalFinishedOptions.ID == "" && *sendApprovalFinishedOptions.Service == "" {
		logging.PrintLog("Either ID or service must be provided", logging.InfoLevel)
		return errors.New("either ID or service must be provided")
	} else if *sendApprovalFinishedOptions.ID != "" && *sendApprovalFinishedOptions.Service != "" {
		logging.PrintLog("Either ID or service must be provided", logging.InfoLevel)
		return errors.New("either ID or service must be provided")
	}
	return nil
}

func sendApprovalFinishedEvent(sendApprovalFinishedOptions sendApprovalFinishedStruct) error {
	var endPoint url.URL
	var apiToken string
	var err error
	if !mocking {
		endPoint, apiToken, err = credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
	} else {
		endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
		endPoint = *endPointPtr
		apiToken = ""
	}
	if err != nil {
		return errors.New(authErrorMsg)
	}

	logging.PrintLog("Starting to send approval.finished event", logging.InfoLevel)

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)
	api, err := internal.APIProvider(endPoint.String(), apiToken)
	if err != nil {
		return err
	}

	var keptnContext string
	var triggeredID string
	var approvalFinishedEvent *keptnv2.ApprovalFinishedEventData

	if *sendApprovalFinishedOptions.ID != "" {
		keptnContext, triggeredID, approvalFinishedEvent, err = getApprovalFinishedForID(api.EventsV1(), sendApprovalFinishedOptions)
	} else if *sendApprovalFinishedOptions.Service != "" {
		keptnContext, triggeredID, approvalFinishedEvent, err = getApprovalFinishedForService(api.EventsV1(),
			api.ShipyardControlV1(), sendApprovalFinishedOptions)
	}
	if err != nil {
		return err
	}

	if approvalFinishedEvent == nil {
		return nil
	}

	responseEvent, err := sendEvent(keptnContext, triggeredID, keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), approvalFinishedEvent, api.APIV1())
	if err != nil {
		return err
	}

	if responseEvent == nil {
		logging.PrintLog("No event returned", logging.QuietLevel)
		return nil
	}

	return nil
}

func sendEvent(keptnContext, triggeredID, eventType string, approvalFinishedEvent interface{}, apiHandler apiutils.APIV1Interface) (*apimodels.EventContext, error) {
	ID := uuid.New().String()
	source, _ := url.Parse("https://github.com/keptn/keptn/cli#" + eventType)

	sdkEvent := cloudevents.NewEvent()
	sdkEvent.SetID(ID)
	sdkEvent.SetType(eventType)
	sdkEvent.SetSource(source.String())
	sdkEvent.SetDataContentType(cloudevents.ApplicationJSON)
	sdkEvent.SetExtension("shkeptncontext", keptnContext)
	sdkEvent.SetExtension("triggeredid", triggeredID)
	sdkEvent.SetData(cloudevents.ApplicationJSON, approvalFinishedEvent)

	eventByte, err := json.Marshal(sdkEvent)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal cloud event. %s", err.Error())
	}

	apiEvent := apimodels.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, &apiEvent)
	if err != nil {
		return nil, fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
	}

	responseEvent, errorObj := apiHandler.SendEvent(apiEvent)
	if errorObj != nil {
		logging.PrintLog("Send "+eventType+" was unsuccessful", logging.QuietLevel)
		return nil, fmt.Errorf("Send %s was unsuccessful. %s", eventType, *errorObj.Message)
	}
	return responseEvent, nil
}

func getApprovalFinishedForService(eventHandler apiutils.EventsV1Interface, scHandler apiutils.ShipyardControlV1Interface,
	approvalFinishedOptions sendApprovalFinishedStruct) (string, string, *keptnv2.ApprovalFinishedEventData, error) {

	allEvents, err := scHandler.GetOpenTriggeredEvents(apiutils.EventFilter{
		Project:   *approvalFinishedOptions.Project,
		Stage:     *approvalFinishedOptions.Stage,
		Service:   *approvalFinishedOptions.Service,
		EventType: keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName),
	})
	if err != nil {
		logging.PrintLog("Open approval.triggered event for service "+*approvalFinishedOptions.Service+" could not be retrieved: "+err.Error(), logging.InfoLevel)
		return "", "", nil, err
	}
	if len(allEvents) == 0 {
		logging.PrintLog("No open approval.triggered event for service "+*approvalFinishedOptions.Service+" has been found", logging.InfoLevel)
		return "", "", nil, nil
	}

	// print all available options
	printApprovalOptions(allEvents, eventHandler, approvalFinishedOptions)

	// select option
	nrOfOptions := len(allEvents)
	selectedOption, err := selectApprovalOption(nrOfOptions)
	if err != nil {
		return "", "", nil, err
	}

	index := selectedOption - 1
	eventToBeApproved := allEvents[index]

	// approve or decline?
	approve := approveOrDecline()

	approvalTriggeredEvent := &keptnv2.ApprovalTriggeredEventData{}

	err = keptnv2.Decode(eventToBeApproved.Data, approvalTriggeredEvent)
	if err != nil {
		logging.PrintLog("Cannot decode approval.triggered event: "+err.Error(), logging.InfoLevel)
		return "", "", nil, err
	}

	var approvalResult keptnv2.ResultType
	if approve {
		approvalResult = keptnv2.ResultPass
	} else {
		approvalResult = keptnv2.ResultFailed
	}

	approvalFinishedEvent := &keptnv2.ApprovalFinishedEventData{
		EventData: keptnv2.EventData{
			Project: approvalTriggeredEvent.Project,
			Stage:   approvalTriggeredEvent.Stage,
			Service: approvalTriggeredEvent.Service,
			Labels:  approvalTriggeredEvent.Labels,
			Status:  keptnv2.StatusSucceeded,
			Result:  approvalResult,
			Message: "",
		},
	}

	return eventToBeApproved.Shkeptncontext, eventToBeApproved.ID, approvalFinishedEvent, nil
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

func printApprovalOptions(approvals []*apimodels.KeptnContextExtendedCE, eventHandler apiutils.EventsV1Interface, approvalFinishedOptions sendApprovalFinishedStruct) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 2, '\t', 0)

	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t%s\t", "OPTION", "IMAGE", "EVALUATION")

	for index, approval := range approvals {
		score := getScoreForApprovalTriggeredEvent(eventHandler, approvalFinishedOptions, approval)
		image := getApprovalImageEvent(approval)
		appendOptionToWriter(w, index, image, score)
	}
	fmt.Fprintf(w, "\n")
}

func appendOptionToWriter(w *tabwriter.Writer, index int, commitID, score string) {
	fmt.Fprintf(w, "\n (%d)\t%s\t%s\t", index+1, commitID, score)
}

func getScoreForApprovalTriggeredEvent(eventHandler apiutils.EventsV1Interface, approvalFinishedOptions sendApprovalFinishedStruct, approval *apimodels.KeptnContextExtendedCE) string {
	score := "n/a"
	evaluationDoneEvents, errorObj := eventHandler.GetEvents(&apiutils.EventFilter{
		Project:      *approvalFinishedOptions.Project,
		Stage:        *approvalFinishedOptions.Stage,
		Service:      *approvalFinishedOptions.Service,
		EventType:    keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName),
		KeptnContext: approval.Shkeptncontext,
	})
	if errorObj != nil {
		return score
	}
	if len(evaluationDoneEvents) == 0 {
		return score
	}
	evaluationDoneData := &keptnv2.EvaluationFinishedEventData{}

	err := mapstructure.Decode(evaluationDoneEvents[0].Data, evaluationDoneData)
	if err != nil {
		return score
	}

	score = fmt.Sprintf("%f", evaluationDoneData.Evaluation.Score)
	return score
}

func getApprovalImageEvent(approval *apimodels.KeptnContextExtendedCE) string {
	unknownImage := "n/a"

	// the approval.triggered event should also include the configurationChange property (see https://github.com/keptn/keptn/issues/3199)
	// therefore, we can cast its data property to a DeploymentTriggeredEventData struct and use the property from this struct
	deploymentTriggeredData := &keptnv2.DeploymentTriggeredEventData{}

	err := keptnv2.Decode(approval.Data, deploymentTriggeredData)

	if err != nil {
		return unknownImage
	}

	if deploymentTriggeredData.ConfigurationChange.Values != nil {
		if image, ok := deploymentTriggeredData.ConfigurationChange.Values["image"].(string); ok {
			return image
		}
	}
	return unknownImage
}

func getApprovalFinishedForID(eventHandler apiutils.EventsV1Interface, sendApprovalFinishedOptions sendApprovalFinishedStruct) (string,
	string, *keptnv2.ApprovalFinishedEventData, error) {
	events, errorObj := eventHandler.GetEvents(&apiutils.EventFilter{
		Project:   *sendApprovalFinishedOptions.Project,
		Stage:     *sendApprovalFinishedOptions.Stage,
		EventType: keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName),
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

	approvalTriggeredEvent := &keptnv2.ApprovalTriggeredEventData{}

	if err := keptnv2.Decode(events[0].Data, approvalTriggeredEvent); err != nil {
		logging.PrintLog("Cannot decode approval.triggered event: "+err.Error(), logging.InfoLevel)
		return "", "", nil, err
	}

	approvalFinishedEvent := &keptnv2.ApprovalFinishedEventData{
		EventData: keptnv2.EventData{
			Project: approvalTriggeredEvent.Project,
			Stage:   approvalTriggeredEvent.Stage,
			Service: approvalTriggeredEvent.Service,
			Labels:  approvalTriggeredEvent.Labels,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
			Message: "",
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
	sendApprovalFinishedOptions.Labels = approvalFinishedCmd.Flags().StringToStringP("labels", "l", nil, "Additional labels to be provided for the service that is to be approved")

}
