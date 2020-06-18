package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/websockethelper"
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

var allowedMonitoringTypes = []string{
	"prometheus",
	"dynatrace",
}

var monitoringCmd = &cobra.Command{
	// Use:          "monitoring <monitoring_provider> --project=<project> --service=<service> --service-indicators=<service_indicators_file_path> --service-objectives=<service_objectives_file_path> --remediation=<remediation_file_path>",
	Use:   "monitoring <monitoring_provider> --project=<project> --service=<service>",
	Short: "Configures monitoring provider",
	Long: `Configure a monitoring solution for a Keptn cluster. This command sets up Dynatrace or Prometheus monitoring if it is not installed. Before executing the command, the dynatrace-service or prometheus-service has to be deployed.

**Note:** If you are executing *keptn configure monitoring dynatrace*, the service flag is optional since Keptn automatically detects the services of a project.

See https://keptn.sh/docs/develop/reference/monitoring/ for more information.
`,
	Example: `keptn configure monitoring dynatrace --project=PROJECTNAME
keptn configure monitoring prometheus --project=PROJECTNAME --service=SERVICENAME`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires a monitoring provider as argument")
		}

		for _, monitoringType := range allowedMonitoringTypes {
			if monitoringType == args[0] {
				return nil
			}
		}

		errorMsg := "Invalid monitoring type. Must be one of: "
		for _, monitoringType := range allowedMonitoringTypes {
			errorMsg = errorMsg + "\n - " + monitoringType
		}

		return errors.New(errorMsg)
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
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		configureMonitoringEventData := &keptn.ConfigureMonitoringEventData{
			Type:    args[0],
			Project: *params.Project,
			Service: *params.Service,
		}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuremonitoring")
		contentType := "application/json"
		sdkEvent := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        keptn.ConfigureMonitoringEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: configureMonitoringEventData,
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		eventByte, err := sdkEvent.MarshalJSON()
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

			// if eventContext is available, open WebSocket communication
			if eventContext != nil && !SuppressWSCommunication {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint, *scheme == "https")
			}

			return nil
		}

		fmt.Println("Skipping send-new artifact due to mocking flag set to true")
		return nil
	},
}

func init() {
	configureCmd.AddCommand(monitoringCmd)
	params = &configureMonitoringCmdParams{}
	params.Project = monitoringCmd.Flags().StringP("project", "p", "", "The name of the project")
	params.Service = monitoringCmd.Flags().StringP("service", "s", "", "The name of the service within the project")
}
