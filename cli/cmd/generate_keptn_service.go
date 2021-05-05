package cmd

import (
	"fmt"
	fileUtils "github.com/keptn/keptn/cli/pkg/file"
	keptnUtils "github.com/keptn/keptn/cli/pkg/git"
	"github.com/spf13/cobra"
	"path/filepath"
)

type generateKeptnServiceStruct struct {
	Service *string   `json:"service"`
	Image   *string   `json:"image"`
	Events  *[]string `json:"events"`
}

var generateKeptnService generateKeptnServiceStruct
var serviceTemplateRepoUrl = "https://github.com/keptn-sandbox/keptn-service-template-go"

var generateKeptnServiceCmd = &cobra.Command{
	Use:          "keptn-service",
	Short:        "Generates keptn service",
	Long:         `Generates keptn service with version check`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Keptn Service CLI generation work in progress \n")
		return generateServiceTemplate(generateKeptnService)
	},
}

func generateServiceTemplate(generateKeptnService generateKeptnServiceStruct) error {
	var err error
	fmt.Printf("Cloning the template\n")
	err = keptnUtils.CloneGitHubUrl(*generateKeptnService.Service, serviceTemplateRepoUrl)
	if err != nil{
		return err
	}
	fmt.Printf("Replacing Image name to: %s\n", *generateKeptnService.Image)
	err = replaceImageFiles(*generateKeptnService.Image, *generateKeptnService.Service)

	if err != nil{
		return err
	}

	fmt.Printf("Creating your service named: %s\n", *generateKeptnService.Service)
	err = replaceServiceName(*generateKeptnService.Service)
	if err != nil{
		return err
	}


	return nil
}


func replaceImageFiles(imageName string, serviceName string) error {
	var err error
	filePatterns := []string{"*"}
	err = filepath.Walk(serviceName, fileUtils.RecursiveRefactor("keptnsandbox/keptn-service-template-go",imageName,filePatterns))
	if err != nil{
		return err
	}
	return nil
}

func replaceServiceName(serviceName string) error {
	filePatterns := []string{"*"}
	err := filepath.Walk(serviceName, fileUtils.RecursiveRefactor("keptn-service-template-go",serviceName,filePatterns))
	if err != nil{
		return err
	}
	return nil
}

func init() {
	generateCmd.AddCommand(generateKeptnServiceCmd)
	generateKeptnService.Service = generateKeptnServiceCmd.Flags().StringP("service", "s", "",
		"Name of the service to be generated ")
	generateKeptnServiceCmd.MarkFlagRequired("service")
	generateKeptnService.Events = generateKeptnServiceCmd.Flags().StringArrayP("events", "e", nil,
		"Comma separated list of cloud-events to listen for")
	generateKeptnServiceCmd.MarkFlagRequired("events")
	generateKeptnService.Image = generateKeptnServiceCmd.Flags().StringP("image", "i", "",
		"The name of the docker image name (organisation/image)")
	generateKeptnServiceCmd.MarkFlagRequired("image")

}
