package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/validator"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
)

type onboardServiceCmdParams struct {
	Project            *string
	ChartFilePath      *string
	DeploymentStrategy *string
}

var onboardServiceParams *onboardServiceCmdParams

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service SERVICENAME --project=PROJECTNAME --chart=FILEPATH",
	Short: "Onboards a new service.",
	Long: `Onboards a new service in the provided project. Therefore, this command 
takes a folder to a Helm chart or an already packed Helm chart as .tgz.
	
Example:
	keptn onboard service carts --project=sockshop --chart=./carts-chart.tgz`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires SERVICENAME")
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		// validate deployment strategy flag
		if *onboardServiceParams.DeploymentStrategy != "" && (*onboardServiceParams.DeploymentStrategy != "direct" && *onboardServiceParams.DeploymentStrategy != "blue_green_service") {
			return errors.New("The provided deployment strategy is not supported. Select: [direct|blue_green_service]")
		}

		// validate chart flag
		*onboardServiceParams.ChartFilePath = keptnutils.ExpandTilde(*onboardServiceParams.ChartFilePath)

		if _, err := os.Stat(*onboardServiceParams.ChartFilePath); os.IsNotExist(err) {
			return errors.New("Provided Helm chart does not exist")
		}

		ch, err := keptnutils.LoadChartFromPath(*onboardServiceParams.ChartFilePath)
		if err != nil {
			return err
		}

		res, err := validator.ValidateHelmChart(ch)
		if err != nil {
			return err
		}

		if !res {
			return errors.New("The provided Helm chart is invalid. Please checkout the requirements")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to onboard service", logging.InfoLevel)

		ch, err := keptnutils.LoadChartFromPath(*onboardServiceParams.ChartFilePath)
		if err != nil {
			return err
		}

		chartData, err := keptnutils.PackageChart(ch)
		if err != nil {
			return err
		}

		data := events.ServiceCreateEventData{
			Project:   *onboardServiceParams.Project,
			Service:   args[0],
			HelmChart: base64.StdEncoding.EncodeToString(chartData),
		}

		if *onboardServiceParams.DeploymentStrategy != "" {
			deplStrategies := make(map[string]events.DeploymentStrategy)

			if *onboardServiceParams.DeploymentStrategy == "direct" {
				deplStrategies["*"] = events.Direct
			} else if *onboardServiceParams.DeploymentStrategy == "blue_green_service" {
				deplStrategies["*"] = events.Duplicate
			} else {
				return fmt.Errorf("The provided deployment strategy %s is not supported. Select: [direct|blue_green_service]", *onboardServiceParams.DeploymentStrategy)
			}

			data.DeploymentStrategies = deplStrategies
		}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#onboardservice")

		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        events.InternalServiceCreateEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: data,
		}

		serviceURL := endPoint
		serviceURL.Path = "v1/service"

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)
		if !mocking {
			responseCE, err := utils.Send(serviceURL, event, apiToken)
			if err != nil {
				logging.PrintLog("Onboard service was unsuccessful", logging.QuietLevel)
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				logging.PrintLog("Response CE is nil", logging.QuietLevel)

				return nil
			}
			if responseCE.Data != nil {
				return websockethelper.PrintWSContentCEResponse(responseCE, endPoint)
			}
		} else {
			fmt.Println("Skipping onboard service due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	onboardCmd.AddCommand(serviceCmd)
	onboardServiceParams = &onboardServiceCmdParams{}
	onboardServiceParams.Project = serviceCmd.Flags().StringP("project", "p", "", "The name of the project")
	serviceCmd.MarkFlagRequired("project")

	onboardServiceParams.ChartFilePath = serviceCmd.Flags().StringP("chart", "", "", "A path to a Helm chart folder or an already archived Helm chart")
	serviceCmd.MarkFlagRequired("chart")

	onboardServiceParams.DeploymentStrategy = serviceCmd.Flags().StringP("deployment-strategy", "", "", "Allows to define a deployment strategy that overrides the shipyard definition for this service")
}
