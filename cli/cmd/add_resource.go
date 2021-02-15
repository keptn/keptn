package cmd

import (
	"errors"
	"fmt"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"io/ioutil"
	"os"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
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

- *--project* - is mandatory. The resource will be added to the root folder in the master branch. 
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
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		*addResourceCmdParams.Resource = keptnutils.ExpandTilde(*addResourceCmdParams.Resource)
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

		resourceHandler := apiutils.NewAuthenticatedResourceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Resource %s could not be uploaded: %s"+endPointErrorReasons,
				*addResourceCmdParams.Resource, endPointErr)
		}

		// Handle different cases of adding resource to a projects default branch, stage branch, and/or service sub-directory
		if isStringFlagSet(addResourceCmdParams.Service) && isBoolFlagSet(addResourceCmdParams.AllStages) {
			// add to all stages
			logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" to all stages of project "+*addResourceCmdParams.Project, logging.InfoLevel)
		} else if areStringFlagsSet(addResourceCmdParams.Service, addResourceCmdParams.Stage) {
			// add to service and stage
			logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" to service "+*addResourceCmdParams.Service+" in stage "+*addResourceCmdParams.Stage+" in project "+*addResourceCmdParams.Project, logging.InfoLevel)
		} else if !isStringFlagSet(addResourceCmdParams.Service) && isStringFlagSet(addResourceCmdParams.Stage) {
			// service is empty, add to stage
			logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" to stage "+*addResourceCmdParams.Stage+" in project "+*addResourceCmdParams.Project, logging.InfoLevel)
		} else if !isStringFlagSet(addResourceCmdParams.Service) && !isStringFlagSet(addResourceCmdParams.Stage) {
			// service and stage are empty, add to default branch
			logging.PrintLog("Adding resource "+*addResourceCmdParams.Resource+" to project "+*addResourceCmdParams.Project, logging.InfoLevel)
		} else {
			return errors.New("Flag 'stage' is missing")
		}

		if !mocking {
			if addResourceCmdParams.AllStages != nil && *addResourceCmdParams.AllStages {
				// Upload to all stages
				// get stages
				stagesHandler := apiutils.NewAuthenticatedStageHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

				stages, err := stagesHandler.GetAllStages(*addResourceCmdParams.Project)
				if err != nil {
					return fmt.Errorf("Failed to retrieve stages for project %s: %v", *addResourceCmdParams.Project, err)
				}

				if len(stages) == 0 {
					return fmt.Errorf("No stages found")
				}

				for _, stage := range stages {
					_, errorObj := resourceHandler.CreateResources(*addResourceCmdParams.Project, stage.StageName, *addResourceCmdParams.Service, resources)
					if errorObj != nil {
						return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded: " + *errorObj.Message)
					}
				}
			} else {
				// upload to specific project/service/stage
				_, errorObj := resourceHandler.CreateResources(*addResourceCmdParams.Project, *addResourceCmdParams.Stage, *addResourceCmdParams.Service, resources)
				if errorObj != nil {
					return errors.New("Resource " + *addResourceCmdParams.Resource + " could not be uploaded: " + *errorObj.Message)
				}
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

		// When setting --all-stages, project and service needs to be set
		if isBoolFlagSet(addResourceCmdParams.AllStages) &&
			!areStringFlagsSet(addResourceCmdParams.Service, addResourceCmdParams.Project) {
			return errors.New("--service and --project need to be supplied when using --all-stages")
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
