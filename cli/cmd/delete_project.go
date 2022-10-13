package cmd

import (
	"errors"
	"fmt"

	"github.com/keptn/keptn/cli/internal"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

// delProjectCmd represents the project command
var delProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME",
	Short: "Deletes a project identified by project name",
	Long: `Deletes a project identified by project name. 

**Notes:**
* If a Git upstream is configured for this project, the referenced upstream repository (e.g., on GitHub) will not be deleted. 
* Services that have been deployed to the Kubernetes cluster are not deleted.
* Namespaces that have been created on the Kubernetes cluster are not deleted.
* Helm-releases created for deployments are not deleted. To clean-up deployed Helm releases, please see [Clean-up after deleting a project](https://keptn.sh/docs/` + getReleaseDocsURL() + `/continuous_delivery/deployment_helm/#clean-up-after-deleting-a-project)
`,
	Example:      `keptn delete project sockshop`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument PROJECTNAME not set")
		} else if len(args) >= 2 {
			cmd.SilenceUsage = false
			return errors.New("too many arguments set")
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to delete project", logging.InfoLevel)

		project := apimodels.Project{
			ProjectName: args[0],
		}

		api, err := internal.APIProvider(endPoint.String(), apiToken)
		if err != nil {
			return internal.OnAPIError(err)
		}

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {

			deleteResp, err := api.APIV1().DeleteProject(project)
			if err != nil {
				logging.PrintLog("Delete project was unsuccessful", logging.InfoLevel)
				return fmt.Errorf("Delete project was unsuccessful. %s", *err.Message)
			}
			if deleteResp != nil {
				logging.PrintLog("Project deleted successfully", logging.InfoLevel)
				logging.PrintLog(deleteResp.Message, logging.InfoLevel)
			}

			return nil
		}

		fmt.Println("Skipping delete project due to mocking flag set to true")
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(delProjectCmd)
}
