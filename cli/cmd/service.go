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
				return errors.New("Provide a Helm values file\n")
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
				fmt.Printf("The specified deployment file is ignored")
			}
			if *serviceFilePath != "" {
				fmt.Println("The specified service file is ignored")
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

		fmt.Println("Starting to onboard service")

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

		fmt.Println("Connecting to server ", endPoint.String())
		_, err = utils.Send(serviceURL, event, apiToken, utils.AddXKeptnSignatureHeader)

		if err != nil {
			fmt.Println("Onboard service was unsuccessful")
			return err
		}

		fmt.Printf("Successfully onboarded service in project %v\n", svcData["project"])
		return nil
	},
}

func init() {
	onboardCmd.AddCommand(serviceCmd)

	project = serviceCmd.Flags().StringP("project", "p", "", "The name of the project")
	serviceCmd.MarkFlagRequired("project")

	// Flags for onboarding a service using Helm
	valuesFilePath = serviceCmd.Flags().StringP("values", "v", "", "The values file")
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
