package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
)

var project *string
var deploymentFilePath *string
var valuesFilePath *string
var serviceFilePath *string
var manifestFilePath *string

type serviceData map[string]interface{}

const manifestEnabled = false

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:          "service",
	Short:        "Onboards a new service.",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if manifestEnabled {
			if *valuesFilePath == "" && *manifestFilePath == "" {
				return errors.New("Specify either a Helm description using the flags values, deployment, and service or a k8s manifest using the flag manifest")
			}
			if *valuesFilePath != "" && *manifestFilePath != "" {
				return errors.New("Error specifying a Helm description as well as a k8s manifest. Only use one option")
			}
		} else {
			if *valuesFilePath == "" {
				cmd.SilenceUsage = false
				return errors.New("Provide a Helm values file")
			}
		}

		if *valuesFilePath != "" {
			valuesData, err := utils.ReadFile(*valuesFilePath)
			if err != nil {
				return err
			}

			if _, err := utils.UnmarshalString(valuesData); err != nil {
				return err
			}

			if *deploymentFilePath != "" {
				if _, err := utils.ReadFile(*deploymentFilePath); err != nil {
					return err
				}
			}

			if *serviceFilePath != "" {
				if _, err := utils.ReadFile(*serviceFilePath); err != nil {
					return err
				}
			}
		} else {
			if *deploymentFilePath != "" {
				utils.PrintLog("The specified deployment file is ignored", utils.InfoLevel)

			}
			if *serviceFilePath != "" {
				utils.PrintLog("The specified service file is ignored", utils.InfoLevel)

			}
			manifestData, err := utils.ReadFile(*manifestFilePath)
			if err != nil {
				return err
			}

			if _, err := utils.UnmarshalString(manifestData); err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		utils.PrintLog("Starting to onboard service", utils.InfoLevel)

		svcData := serviceData{}
		svcData["project"] = *project

		if *valuesFilePath != "" {
			valuesData, err := utils.ReadFile(*valuesFilePath)
			if err != nil {
				return err
			}
			data, err := utils.UnmarshalString(valuesData)
			if err != nil {
				return err
			}
			svcData["values"] = data

			// check if service name is valid
			svcName := ""
			if svcData["values"] != nil && svcData["values"].(map[string]interface{})["service"] != nil && svcData["values"].(map[string]interface{})["service"].(map[string]interface{})["name"] != nil {
				svcName = svcData["values"].(map[string]interface{})["service"].(map[string]interface{})["name"].(string)
			}
			if svcName == "" {
				return errors.New("Service name could not be retrieved. Please verify that a service name is defined in your .yaml file")
			}

			if !utils.ValidateK8sName(svcName) {
				errorMsg := "Service name as defined in the .yaml file includes invalid characters or is not well-formed.\n"
				errorMsg += "keptn relies on Helm charts and thus these conventions have to be followed: "
				errorMsg += "start with a lower case letter, then lower case letters, dash and numbers are allowed.\n"
				errorMsg += "You can find the guidelines here: https://github.com/helm/helm/blob/master/docs/chart_best_practices/conventions.md#chart-names\n"
				errorMsg += "Please update service name and try again."
				return errors.New(errorMsg)
			}

			deployment := make(map[string]string)

			if *deploymentFilePath != "" {
				content, err := utils.ReadFile(*deploymentFilePath)
				if err != nil {
					return err
				}
				deployment["deployment"] = content
			}

			if *serviceFilePath != "" {
				content, err := utils.ReadFile(*serviceFilePath)
				if err != nil {
					return err
				}
				deployment["service"] = content
			}

			if len(deployment) > 0 {
				svcData["templates"] = deployment
			}
		} else {
			manifestData, err := utils.ReadFile(*manifestFilePath)
			if err != nil {
				return err
			}
			data, err := utils.UnmarshalString(manifestData)
			if err != nil {
				return err
			}
			svcData["manifest"] = data
		}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#onboardservice")

		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        "onboard.service",
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: svcData,
		}

		serviceURL := endPoint
		serviceURL.Path = "service"

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
				return websockethelper.PrintWSContent(responseCE)
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

	// Flags for onboarding a service using Helm
	valuesFilePath = serviceCmd.Flags().StringP("values", "", "", "The values file")
	deploymentFilePath = serviceCmd.Flags().StringP("deployment", "d", "", "The deployment file")
	serviceFilePath = serviceCmd.Flags().StringP("service", "s", "", "The service file")

	if manifestEnabled {
		// Flags for onboarding using "pure" k8s manifests
		manifestFilePath = serviceCmd.Flags().StringP("manifest", "m", "", "The manifest file")
		serviceCmd.Long = `Onboards a new service in the provided project. Therefore, this command 
either takes a service description given as yaml or takes a Helm values, deployment and service description.

Examples:
	# Using a k8s service description
	keptn onboard service --project=carts --manifest=serviceDesc.yaml 
	# Using a Helm chart
	keptn onboard service --project=carts --values=values.yaml --deployment=deployment.yaml --service=service.yaml`
	} else {
		serviceCmd.Long = `Onboards a new service in the provided project. Therefore, this command 
takes a Helm values, deployment and service description.

Examples:
	keptn onboard service --project=carts --values=values.yaml --deployment=deployment.yaml --service=service.yaml`
	}
}
