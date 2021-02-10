package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/validator"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

type onboardServiceCmdParams struct {
	Project       *string
	ChartFilePath *string
}

var onboardServiceParams *onboardServiceCmdParams

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service SERVICENAME --project=PROJECTNAME --chart=FILEPATH",
	Short: "Onboards a new service and its Helm chart to a project",
	Long: `Onboards a new service and its Helm chart to the provided project. 
Therefore, this command takes a folder to a Helm chart or an already packed Helm chart as .tgz.
`,
	Example: `keptn onboard service SERVICENAME --project=PROJECTNAME --chart=FILEPATH

keptn onboard service SERVICENAME --project=PROJECTNAME --chart=HELM_CHART.tgz
`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SERVICENAME not set")
		}
		if !keptncommon.ValidateKeptnEntityName(args[0]) {
			errorMsg := "Service name contains upper case letter(s) or special character(s).\n"
			return errors.New(errorMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := doOnboardServicePreRunChecks(args); err != nil {
			return err
		}
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to onboard service", logging.InfoLevel)

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		chart, err := keptnutils.LoadChartFromPath(*onboardServiceParams.ChartFilePath)
		if err != nil {
			return err
		}

		chartData, err := keptnutils.PackageChart(chart) // ToDo: use chartPackager.package
		if err != nil {
			return err
		}

		//helmChart := base64.StdEncoding.EncodeToString(chartData)
		serviceName := args[0]

		service := apimodels.CreateService{
			ServiceName: &serviceName,
			// HelmChart:   helmChart,
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			_, err := apiHandler.CreateService(*onboardServiceParams.Project, service)
			if err != nil {
				logging.PrintLog("Onboard service: Create service was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Onboard service: Create service was unsuccessful. %s", *err.Message)
			}

			// initialize handlers
			resourceHandler := apiutils.NewAuthenticatedResourceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
			stagesHandler := apiutils.NewAuthenticatedStageHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
			chartStorer := keptnutils.NewChartStorer(resourceHandler)

			// get all stages
			stages, err2 := stagesHandler.GetAllStages(*onboardServiceParams.Project)

			if err2 != nil {
				return fmt.Errorf("Failed to retrieve stages for project %s: %v", *stageParameter.project, err)
			}

			if len(stages) == 0 {
				fmt.Println("No stages found")
				return nil
			}

			for _, stage := range stages {
				storeOpts := keptnutils.StoreChartOptions{
					Project:   *onboardServiceParams.Project,
					Service:   serviceName,
					Stage:     stage.StageName,
					ChartName: serviceName,
					HelmChart: chartData,
				}

				if _, err := chartStorer.Store(storeOpts); err != nil {
					logging.PrintLog("Error when storing the Helm Chart: " + err.Error(), logging.QuietLevel)
					return fmt.Errorf("Onboard service: Storing charts was unsuccessful. %v", err)
				}
			}

			logging.PrintLog("Service onboarded successfully", logging.InfoLevel)

			return nil
		}

		logging.PrintLog("Skipping onboard service due to mocking flag set to true", logging.InfoLevel)
		return nil
	},
}

func doOnboardServicePreRunChecks(args []string) error {
	// validate chart flag
	*onboardServiceParams.ChartFilePath = keptnutils.ExpandTilde(*onboardServiceParams.ChartFilePath)

	if _, err := os.Stat(*onboardServiceParams.ChartFilePath); os.IsNotExist(err) {
		return errors.New("Provided Helm chart does not exist")
	}

	ch, err := keptnutils.LoadChartFromPath(*onboardServiceParams.ChartFilePath)
	if err != nil {
		return err
	}

	res, err := validator.ValidateHelmChart(ch, args[0])
	if err != nil {
		return err
	}

	if !res {
		return errors.New("The provided Helm chart is invalid. Please checkout the requirements")
	}

	return nil
}

func init() {
	onboardCmd.AddCommand(serviceCmd)
	onboardServiceParams = &onboardServiceCmdParams{}
	onboardServiceParams.Project = serviceCmd.Flags().StringP("project", "p", "", "The name of the project")
	serviceCmd.MarkFlagRequired("project")

	onboardServiceParams.ChartFilePath = serviceCmd.Flags().StringP("chart", "", "", "A path to a Helm chart folder or an already archived Helm chart")
	serviceCmd.MarkFlagRequired("chart")
}
