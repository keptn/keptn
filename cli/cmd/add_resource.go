package cmd

import (
	"github.com/spf13/cobra"
)

type addResourceCommandParameters struct {
	Project     *string
	Stage       *string
	Service     *string
	Resource    *string
	ResourceURI *string
}

var addResourceCmdParams *addResourceCommandParameters

var addResourceCmd = &cobra.Command{
	Use:   "add-resource",
	Short: "Adds a resource to a service within your project in the specified stage",
	Long: `Adds a resource to a service within your project in the specified stage
	
	Example: 
	./keptn add-resource --project=sockshop --stage=dev --service=carts --resource=./jmeter.jmx --resourceUri=jmeter/functional.jmx
	`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addResourceCmd)
	addResourceCmdParams = &addResourceCommandParameters{}
	addResourceCmdParams.Project = addResourceCmd.Flags().StringP("project", "p", "", "The name of the project")
	addResourceCmdParams.Stage = addResourceCmd.Flags().StringP("stage", "s", "", "The name of the stage")
	addResourceCmdParams.Service = addResourceCmd.Flags().StringP("service", "svc", "", "The name of the service within the project")
	addResourceCmdParams.Resource = addResourceCmd.Flags().StringP("resource", "r", "", "Path pointing to the resource on your local file system")
	addResourceCmdParams.ResourceURI = addResourceCmd.Flags().StringP("resourceUri", "ru", "", "Optional: Location where the resource should be stored within the config repo. If empty, The name of the resource will be the same as on your local file system")

}
