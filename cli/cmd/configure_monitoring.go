package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/mitchellh/mapstructure"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type configureMonitoringCmdParams struct {
	Type              *string
	Project           *string
	Service           *string
	ServiceIndicators *string
	ServiceObjectives *string
	Remediation       *string
}

var params *configureMonitoringCmdParams

var monitoringCmd = &cobra.Command{
	// Use:          "monitoring <monitoring_provider> --project=<project> --service=<service> --service-indicators=<service_indicators_file_path> --service-objectives=<service_objectives_file_path> --remediation=<remediation_file_path>",
	Use:   "monitoring <monitoring_provider> --project=<project> --service=<service>",
	Short: "Configures a monitoring provider",
	Long: `Configure a monitoring solution for the deployments managed by Keptn.

**Note:** If you are executing *keptn configure monitoring dynatrace*, the service flag is optional since Keptn automatically detects the services of a project. 
See https://keptn.sh/docs/` + keptnReleaseDocsURL + `/monitoring/dynatrace/install/ for more information.
`,
	Example: `keptn configure monitoring dynatrace --project=PROJECTNAME
keptn configure monitoring prometheus --project=PROJECTNAME --service=SERVICENAME`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires a monitoring provider as argument")
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "prometheus" {
			if *params.Project == "" {
				return errors.New("Please specify a project")
			}
			if *params.Service == "" {
				return errors.New("Please specify a service")
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		configureMonitoringEventData := &keptn.ConfigureMonitoringEventData{
			Type:    args[0],
			Project: *params.Project,
			Service: *params.Service,
		}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuremonitoring")

		sdkEvent := cloudevents.NewEvent()
		sdkEvent.SetID(uuid.New().String())
		sdkEvent.SetType(keptn.ConfigureMonitoringEventType)
		sdkEvent.SetSource(source.String())
		sdkEvent.SetDataContentType(cloudevents.ApplicationJSON)
		sdkEvent.SetData(cloudevents.ApplicationJSON, configureMonitoringEventData)

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		eventByte, err := json.Marshal(sdkEvent)
		if err != nil {
			return fmt.Errorf("Failed to marshal cloud event. %s", err.Error())
		}

		apiEvent := apimodels.KeptnContextExtendedCE{}
		err = json.Unmarshal(eventByte, &apiEvent)
		if err != nil {
			return fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
		}

		if !mocking {
			eventContext, err := apiHandler.SendEvent(apiEvent)
			if err != nil {
				logging.PrintLog("Sending configure-monitoring event was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Sending configure-monitoring event was unsuccessful. %s", *err.Message)
			}

			maxFetchEventRetries := 10
			fetchEventRetryTime := 3
			// if eventContext is available, try to fetch the correlating configure-monitoring.finished event
			gotEvent := false
			if eventContext != nil {
				for i := 0; i < maxFetchEventRetries; i = i + i {
					events, errObj := eventHandler.GetEvents(&apiutils.EventFilter{
						KeptnContext: *eventContext.KeptnContext,
						EventType:    keptnv2.GetFinishedEventType("configure-monitoring"),
					})
					if errObj == nil && len(events) > 0 {
						for _, event := range events {
							eventData := &keptnv2.EventData{}
							decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
								Squash: true,
								Result: eventData,
							})
							if err != nil {
								return fmt.Errorf("could not decode event: " + err.Error())
							}

							if err := decoder.Decode(event.Data); err != nil {
								return fmt.Errorf("could not decode event: " + err.Error())
							}

							logging.PrintLog(eventData.Message, logging.InfoLevel)
						}
						gotEvent = true
						break
					}
					<-time.After(time.Duration(fetchEventRetryTime) * time.Second)
				}

				if !gotEvent {
					logging.PrintLog("Could not retrieve configure-monitoring.finished event with KeptnContext "+*eventContext.KeptnContext+". Please check your monitoring solution if it has been configured correctly.", logging.InfoLevel)
				}
			}

			return nil
		}

		fmt.Println("Skipping configure monitoring due to mocking flag set to true")
		return nil
	},
}

func init() {
	configureCmd.AddCommand(monitoringCmd)
	params = &configureMonitoringCmdParams{}
	params.Project = monitoringCmd.Flags().StringP("project", "p", "", "The name of the project")
	params.Service = monitoringCmd.Flags().StringP("service", "s", "", "The name of the service within the project")
}
