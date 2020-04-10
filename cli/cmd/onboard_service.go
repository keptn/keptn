package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/keptn/keptn/cli/utils/websockethelper"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/validator"
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
	Short: "Onboards a new service to a project",
	Long: `Onboards a new service to the provided project. Therefore, this command 
takes a folder to a Helm chart or an already packed Helm chart as .tgz.
	
Example:
	keptn onboard service carts --project=sockshop --chart=./carts-chart.tgz`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SERVICENAME not set")
		}
		if !keptnutils.ValidateKeptnEntityName(args[0]) {
			errorMsg := "Service name contains upper case letter(s) or special character(s).\n"
			return errors.New(errorMsg)
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

		chart, err := keptnutils.LoadChartFromPath(*onboardServiceParams.ChartFilePath)
		if err != nil {
			return err
		}

		chartData, err := keptnutils.PackageChart(chart)
		if err != nil {
			return err
		}

		helmChart := base64.StdEncoding.EncodeToString(chartData)
		service := apimodels.Service{
			ServiceName: &args[0],
			HelmChart:   helmChart,
		}

		if *onboardServiceParams.DeploymentStrategy != "" {
			deplStrategies := make(map[string]string)

			if *onboardServiceParams.DeploymentStrategy == "direct" {
				deplStrategies["*"] = events.Direct.String()
			} else if *onboardServiceParams.DeploymentStrategy == "blue_green_service" {
				deplStrategies["*"] = events.Duplicate.String()
			} else {
				return fmt.Errorf("The provided deployment strategy %s is not supported. Select: [direct|blue_green_service]", *onboardServiceParams.DeploymentStrategy)
			}

			service.DeploymentStrategies = deplStrategies
		}

		serviceHandler := apiutils.NewAuthenticatedServiceHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := serviceHandler.CreateService(*onboardServiceParams.Project, service)
			if err != nil {
				logging.PrintLog("Onboard service was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Onboard service was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if eventContext != nil && !SuppressWSCommunication {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint, *scheme == "https")
			}

			return nil
		}

		fmt.Println("Skipping onboard service due to mocking flag set to true")
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
