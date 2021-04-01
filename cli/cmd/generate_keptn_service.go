package cmd

import (
	"fmt"
	keptnutils "github.com/keptn/keptn/cli/pkg/git"
	"github.com/spf13/cobra"
)

type generateKeptnServiceStruct struct {
	Service *string            `json:"service"`
	Image   *string            `json:"image"`
	Events  *map[string]string `json:"events"`
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
	err = keptnutils.CloneGitHubUrl(*generateKeptnService.Service, serviceTemplateRepoUrl)
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
	generateKeptnService.Events = generateKeptnServiceCmd.Flags().StringToStringP("events", "e", nil,
		"Comma separated list of cloud-events to listen for")
	generateKeptnServiceCmd.MarkFlagRequired("events")
	generateKeptnService.Image = generateKeptnServiceCmd.Flags().StringP("image", "i", "",
		"The name of the docker image name (organisation/image)")
	generateKeptnServiceCmd.MarkFlagRequired("image")

}
