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
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/validator"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
)

var project *string
var chartFilePath *string

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

		if _, err := os.Stat(*chartFilePath); os.IsNotExist(err) {
			return errors.New("Provided Helm chart does not exist")
		}

		chartData, err := ioutil.ReadFile(*chartFilePath)
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

		utils.PrintLog("Starting to onboard service", utils.InfoLevel)

		chartData, err := ioutil.ReadFile(*chartFilePath)
		if err != nil {
			return err
		}
		data := events.ServiceCreateEventData{Project: *project, Service: args[0],
			HelmChart: base64.StdEncoding.EncodeToString(chartData)}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#onboardservice")

		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        events.ServiceCreateEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: data,
		}

		serviceURL := endPoint
		serviceURL.Path = "v1/service"

		utils.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), utils.VerboseLevel)
		if !mocking {
			responseCE, err := utils.Send(serviceURL, event, apiToken)
			if err != nil {
				utils.PrintLog("Onboard service was unsuccessful", utils.QuietLevel)
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				utils.PrintLog("Response CE is nil", utils.QuietLevel)

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

	project = serviceCmd.Flags().StringP("project", "p", "", "The name of the project")
	serviceCmd.MarkFlagRequired("project")

	chartFilePath = serviceCmd.Flags().StringP("chart", "", "",
		"A path to a packed Helm chart. Use `helm package chart_name` to pack your chart")
	serviceCmd.MarkFlagRequired("chart")
}
