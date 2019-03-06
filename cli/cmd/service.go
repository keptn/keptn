package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/knative/pkg/cloudevents"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var project *string
var deploymentFilePath *string
var valuesFilePath *string
var serviceFilePath *string

type serviceData map[string]interface{}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Onboards a new service.",
	Long: `Onboards a new service. Therefore, this command takes a 
	service description given as yaml and onboards 
	this service in the provided project. 
	Optionally, this command takes a Helm deployment and service description.
	Usage of \"onboard service\":

keptn onboard service --project=carts --values=values.yaml --deployment=deployment.yaml --service=service.yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		valuesData, err := readFile(*valuesFilePath)
		if err != nil {
			return err
		}

		if _, err := unmarshalValues(valuesData); err != nil {
			return err
		}

		if *deploymentFilePath != "" {
			if _, err := readFile(*deploymentFilePath); err != nil {
				return err
			}
		}

		if *serviceFilePath != "" {
			if _, err := readFile(*serviceFilePath); err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to onboard service")

		svcData := serviceData{}
		svcData["project"] = *project

		valuesData, err := readFile(*valuesFilePath)
		if err != nil {
			return err
		}
		data, err := unmarshalValues(valuesData)
		if err != nil {
			return err
		}
		svcData["values"] = data

		deployment := make(map[string]string)

		if *deploymentFilePath != "" {
			content, err := readFile(*deploymentFilePath)
			if err != nil {
				return err
			}
			deployment["deployment"] = content
		}

		if *serviceFilePath != "" {
			content, err := readFile(*serviceFilePath)
			if err != nil {
				return err
			}
			deployment["service"] = content
		}

		if len(deployment) > 0 {
			svcData["templates"] = deployment
		}

		builder := cloudevents.Builder{
			Source:    "https://github.com/keptn/keptn/cli#onboardservice",
			EventType: "onboard.service",
		}
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil || endPoint == "" {
			fmt.Printf("onboard service called without beeing authenticated.")
			return errors.New("This command requires to be authenticated. See \"keptn auth\" for details")
		}
		projectEndPoint := endPoint + "service"
		err = utils.Send(projectEndPoint, apiToken, builder, svcData)
		if err != nil {
			fmt.Printf("onboard service command was unsuccessful. Details: %v", err)
			return err
		}
		fmt.Printf("Successfully onboarded service in project %v\n", svcData["Project"])
		return nil
	},
}

func readFile(fileName string) (string, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", fmt.Errorf("Cannot find file %s", fileName)
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshalValues(data string) (interface{}, error) {
	var body interface{}
	err := yaml.Unmarshal([]byte(data), &body)
	if err != nil {
		return nil, err
	}
	return convert(body), nil
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func init() {
	onboardCmd.AddCommand(serviceCmd)

	project = serviceCmd.Flags().StringP("project", "p", "", "The name of the project")
	serviceCmd.MarkFlagRequired("project")

	valuesFilePath = serviceCmd.Flags().StringP("values", "v", "", "The values file")
	serviceCmd.MarkFlagRequired("values")

	deploymentFilePath = serviceCmd.Flags().StringP("deployment", "d", "", "The deployment file")

	serviceFilePath = serviceCmd.Flags().StringP("service", "s", "", "The service file")
}
