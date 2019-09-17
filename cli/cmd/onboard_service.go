package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
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
	Project       *string
	ChartFilePath *string
	Direct        bool
}

var onboardServiceParams *onboardServiceCmdParams

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service service_name",
	Short: "Onboards a new service.",
	Long: `Onboards a new service in the provided project. Therefore, this command 
takes a Helm chart as packaged .tgz.
	
Examples:
	keptn onboard service service_name --project=carts --chart=chart.tgz`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires service_name")
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		*onboardServiceParams.ChartFilePath = keptnutils.ExpandTilde(*onboardServiceParams.ChartFilePath)

		if _, err := os.Stat(*onboardServiceParams.ChartFilePath); os.IsNotExist(err) {
			return errors.New("Provided Helm chart does not exist")
		}

		chartData, err := ioutil.ReadFile(*onboardServiceParams.ChartFilePath)
		if err != nil {
			return err
		}

		res, err := validator.ValidateHelmChart(chartData)
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

		chartData, err := ioutil.ReadFile(*onboardServiceParams.ChartFilePath)
		if err != nil {
			return err
		}
		data := events.ServiceCreateEventData{Project: *onboardServiceParams.Project, Service: args[0],
			HelmChart: base64.StdEncoding.EncodeToString(chartData)}
		if onboardServiceParams.Direct {
			deplStrategies := make(map[string]events.DeploymentStrategy)
			deplStrategies["*"] = events.Direct
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

	onboardServiceParams.ChartFilePath = serviceCmd.Flags().StringP("chart", "", "",
		"A path to a packed Helm chart. Use `helm package chart_name` to pack your chart")
	serviceCmd.MarkFlagRequired("chart")

	serviceCmd.PersistentFlags().BoolVarP(&onboardServiceParams.Direct, "direct", "", false, "allows to set the deployment strategy to direct for the onboarded service")
}
