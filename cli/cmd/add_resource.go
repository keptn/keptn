package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/keptn/cli/internal"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type addResourceCommandParameters struct {
	Project     *string
	Stage       *string
	AllStages   *bool
	Service     *string
	Resource    *string
	ResourceURI *string
}

var addResourceCmdParams *addResourceCommandParameters

var addResourceCmd = &cobra.Command{
	Use:   "add-resource --project=PROJECT --stage=STAGE --service=SERVICE --resource=FILEPATH --resourceUri=FILEPATH",
	Short: "Adds a local resource to a service within your project in the specified stage",
	Long: `Adds a local resource to a service within your project in the specified stage. The resource is then stored within the Git repository.

This command allows adding, for example, *test files* to a service, which will then be used by a test service (e.g., jmeter-service) during the continuous delivery.

To specify a unique resource identifier (URI) for this resource, the optional flag *--resourceUri* can be set to a file path. 
By default, the URI is set to the file path specified at the *--resource* flag. 
From a technical perspective, the file provided via the *--resource* flag is stored with the path and name specified within *--resourceUri* flag.

**The target location of the resource:**

- *--project* - is mandatory. The resource will be added to the root folder in the master branch. Do not use this command alone but add --all-stages if you are using resource service in branch mode.
- *--stage* - is optional (when the *--service* flag is not used). The resource will be added to the root folder in the stage branch.
- *--service* - is optional. The resource will be added to the service folder in the stage branch.
`,
	Example: `keptn add-resource --project=musicshop --stage=hardening --service=catalogue --resource=slo.yaml
keptn add-resource --project=musicshop --stage=hardening --service=catalogue --resource=slo-quality-gates.yaml --resourceUri=slo.yaml
keptn add-resource --project=sockshop --stage=dev --service=carts --resource=./jmeter.jmx --resourceUri=jmeter/functional.jmx
keptn add-resource --project=rockshop --stage=production --service=shop --resource=./basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx
keptn add-resource --project=keptn --service=keptn-control-plane --all-stages --resource=0.7.3_keptn-installer.tgz --resourceUri=helm/keptn-control-plane.tgz`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		*addResourceCmdParams.Resource = fileutils.ExpandTilde(*addResourceCmdParams.Resource)
		if !fileExists(*addResourceCmdParams.Resource) {
			return errors.New("File " + *addResourceCmdParams.Resource + " not found on local file system")
		}

		resourceContent, err := ioutil.ReadFile(*addResourceCmdParams.Resource)
		if err != nil {
			return errors.New("File " + *addResourceCmdParams.Resource + " could not be read")
		}

		if *addResourceCmdParams.ResourceURI == "" {
			addResourceCmdParams.ResourceURI = addResourceCmdParams.Resource
		}

		resourceContentStr := string(resourceContent)
		resources := []*apimodels.Resource{
			{
				ResourceContent: resourceContentStr,
				ResourceURI:     addResourceCmdParams.ResourceURI,
			},
		}

		api, err := internal.APIProvider(endPoint.String(), apiToken)
		if err != nil {
			return internal.OnAPIError(err)
		}

		if !mocking {
			if !isStringFlagSet(addResourceCmdParams.Service) && !isBoolFlagSet(addResourceCmdParams.AllStages) && !isStringFlagSet(addResourceCmdParams.Stage) {
				// project resource
				logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" as a project resource to project "+*addResourceCmdParams.Project, logging.InfoLevel)
				_, errorObj := api.ResourcesV1().CreateResources(*addResourceCmdParams.Project, "", "", resources)
				if errorObj != nil {
					return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded as a project resource: " + *errorObj.Message)
				}
			} else if !isStringFlagSet(addResourceCmdParams.Service) && isBoolFlagSet(addResourceCmdParams.AllStages) {
				// stage resource to all stages
				logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" as a stage resource to all stages for project "+*addResourceCmdParams.Project, logging.InfoLevel)
				api, err := internal.APIProvider(endPoint.String(), apiToken)
				if err != nil {
					return err
				}

				stages, err := api.StagesV1().GetAllStages(*addResourceCmdParams.Project)
				if err != nil {
					return fmt.Errorf("Failed to retrieve stages for project %s: %v", *addResourceCmdParams.Project, err)
				}

				if len(stages) == 0 {
					return fmt.Errorf("No stages found")
				}

				for _, stage := range stages {
					_, errorObj := api.ResourcesV1().CreateResources(*addResourceCmdParams.Project, stage.StageName, "", resources)
					if errorObj != nil {
						return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded as stage resource: " + *errorObj.Message)
					}
				}
			} else if !isStringFlagSet(addResourceCmdParams.Service) && isStringFlagSet(addResourceCmdParams.Stage) {
				// stage resource to a defined stage
				logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" as a stage resource to stage "+*addResourceCmdParams.Stage+" for project "+*addResourceCmdParams.Project, logging.InfoLevel)
				_, errorObj := api.ResourcesV1().CreateResources(*addResourceCmdParams.Project, *addResourceCmdParams.Stage, "", resources)
				if errorObj != nil {
					return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded as a stage resource: " + *errorObj.Message)
				}
			} else if isStringFlagSet(addResourceCmdParams.Service) && isBoolFlagSet(addResourceCmdParams.AllStages) && !isStringFlagSet(addResourceCmdParams.Stage) {
				// service resource to all stages for a defined service
				logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" as a service resource to service "+*addResourceCmdParams.Service+" for all stages for project "+*addResourceCmdParams.Project, logging.InfoLevel)
				api, err := internal.APIProvider(endPoint.String(), apiToken)
				if err != nil {
					return err
				}

				stages, err := api.StagesV1().GetAllStages(*addResourceCmdParams.Project)
				if err != nil {
					return fmt.Errorf("Failed to retrieve stages for project %s: %v", *addResourceCmdParams.Project, err)
				}

				if len(stages) == 0 {
					return fmt.Errorf("No stages found")
				}

				for _, stage := range stages {
					_, errorObj := api.ResourcesV1().CreateResources(*addResourceCmdParams.Project, stage.StageName, *addResourceCmdParams.Service, resources)
					if errorObj != nil {
						return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded as service resource: " + *errorObj.Message)
					}
				}
			} else if isStringFlagSet(addResourceCmdParams.Service) && !isBoolFlagSet(addResourceCmdParams.AllStages) && isStringFlagSet(addResourceCmdParams.Stage) {
				// service resource to defined stage for a defined service
				logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" as a service resource to service "+*addResourceCmdParams.Service+" for stage "+*addResourceCmdParams.Stage+" for project "+*addResourceCmdParams.Project, logging.InfoLevel)
				_, errorObj := api.ResourcesV1().CreateResources(*addResourceCmdParams.Project, *addResourceCmdParams.Stage, *addResourceCmdParams.Service, resources)
				if errorObj != nil {
					return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded as a service resource: " + *errorObj.Message)
				}
			} else {
				return errors.New("Invalid combination of input parameters")
			}

			logging.PrintLog("Resource has been uploaded.", logging.InfoLevel)
			return nil
		}

		fmt.Println("Skipping add resource due to mocking flag set to true")
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Prevent setting --stage and --all-stages at the same time
		if isStringFlagSet(addResourceCmdParams.Stage) && isBoolFlagSet(addResourceCmdParams.AllStages) {
			return errors.New("Cannot use --stage and --all-stages at the same time")
		}

		if !isStringFlagSet(addResourceCmdParams.Stage) && !isBoolFlagSet(addResourceCmdParams.AllStages) && isStringFlagSet(addResourceCmdParams.Service) {
			return errors.New("Flag 'stage' is missing")
		}

		return nil
	},
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	rootCmd.AddCommand(addResourceCmd)
	addResourceCmdParams = &addResourceCommandParameters{}

	addResourceCmdParams.Project = addResourceCmd.Flags().StringP("project", "p", "", "The name of the project")
	addResourceCmd.MarkFlagRequired("project")

	addResourceCmdParams.Stage = addResourceCmd.Flags().StringP("stage", "s", "", "The name of the stage (cannot be used together with (--all-stages)")

	addResourceCmdParams.AllStages = addResourceCmd.Flags().Bool("all-stages", false, "Add resource to all stages (can not be used together with (--stage)")

	addResourceCmdParams.Service = addResourceCmd.Flags().StringP("service", "", "", "The name of the service within the project")

	addResourceCmdParams.Resource = addResourceCmd.Flags().StringP("resource", "r", "", "Path pointing to the resource on your local file system")
	addResourceCmd.MarkFlagRequired("resource")

	addResourceCmdParams.ResourceURI = addResourceCmd.Flags().StringP("resourceUri", "", "", "Optional: Location where the resource should be stored within the config repo. If empty, The name of the resource will be the same as on your local file system")

}
