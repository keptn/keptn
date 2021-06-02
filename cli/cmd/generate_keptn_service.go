package cmd

import (
	"fmt"
	fileUtils "github.com/keptn/keptn/cli/pkg/file"
	keptnUtils "github.com/keptn/keptn/cli/pkg/git"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type generateKeptnServiceStruct struct {
	Service *string   `json:"service"`
	Image   *string   `json:"image"`
	Events  *[]string `json:"events"`
}

var generateKeptnService generateKeptnServiceStruct
var serviceTemplateRepoUrl = "https://github.com/keptn-sandbox/keptn-service-template-go"
var serviceTemplateRefUrl = "https://api.github.com/repos/keptn-sandbox/keptn-service-template-go/git/refs/tags"

var generateKeptnServiceCmd = &cobra.Command{
	Use:   "keptn-service --service=SERVICE_NAME --image=IMAGE_NAME --events=EVENT1,EVENT2",
	Short: "Generates a keptn service with your image and cloud events you wish to listen for",
	Long: `Creates a new keptn service with your image and events you wish to listen for.

	This command can be used to automatically clone the template from https://github.com/keptn-sandbox/keptn-service-template-go according to your keptn cluster version, and set the requirements as mentioned in the README.md.

	A new folder with the service name will be generated which will hold all the template requirements for your own keptn-service.
	`,
	Example: `keptn generate keptn-service --service=myService --image=SOME_IMAGE_NAME --events=sh.keptn.events.problem,sh.keptn.events.deployment-finished`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Keptn Service CLI generation work in progress \n")
		return generateServiceTemplate(generateKeptnService)
	},
}

func generateServiceTemplate(generateKeptnService generateKeptnServiceStruct) error {
	var err error
	keptnVersion, err := getInstalledKeptnVersion()
	if err != nil {
		return err
	}
	refTag, err := keptnUtils.GetGitHubRefs(serviceTemplateRefUrl, keptnVersion)

	fmt.Printf("Cloning the template\n")
	err = keptnUtils.CloneGitHubUrl(*generateKeptnService.Service, serviceTemplateRepoUrl, refTag)
	if err != nil {
		return err
	}
	fmt.Printf("Replacing Image name to: %s\n", *generateKeptnService.Image)
	err = replaceImageFiles(*generateKeptnService.Image, *generateKeptnService.Service)

	if err != nil {
		return err
	}

	fmt.Printf("Creating your service named: %s\n", *generateKeptnService.Service)
	err = replaceServiceName(*generateKeptnService.Service)
	if err != nil {
		return err
	}
	err = removeGeneratedGitFolder(*generateKeptnService.Service)
	if err != nil {
		return err
	}
	return nil
}

func removeGeneratedGitFolder(serviceName string) error {
	err := os.RemoveAll(serviceName + "/.git")
	if err != nil {
		return err
	}
	return nil
}

func replaceImageFiles(imageName string, serviceName string) error {
	var err error
	filePatterns := []string{"*"}
	err = filepath.Walk(serviceName, fileUtils.RecursiveRefactor("keptnsandbox/keptn-service-template-go", imageName, filePatterns))
	if err != nil {
		return err
	}
	return nil
}

func replaceServiceName(serviceName string) error {
	filePatterns := []string{"*"}
	err := filepath.Walk(serviceName, fileUtils.RecursiveRefactor("keptn-service-template-go", serviceName, filePatterns))
	if err != nil {
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
