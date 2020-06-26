package cmd

import (
	"errors"
	"fmt"

	"github.com/keptn/keptn/cli/pkg/websockethelper"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

// crprojectCmd represents the project command
var delProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME",
	Short: "Deletes a project identified by project name",
	Long: `Deletes a project identified by project name. 

**Known Limitations**:

* If a Git upstream is configured for this project, the referenced upstream repository (e.g., on GitHub) will not be deleted. 
* Services that have been deployed to the Kubernetes cluster are not deleted (same goes for the namespaces).
* Helm-releases created for deployments are not deleted - see https://keptn.sh/docs/develop/reference/helm/#clean-up-after-deleting-a-project
`,
	Example:      `keptn delete project sockshop`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument PROJECTNAME not set")
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to delete project", logging.InfoLevel)

		project := apimodels.Project{
			ProjectName: args[0],
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := apiHandler.DeleteProject(project)
			if err != nil {
				fmt.Println("Delete project was unsuccessful")
				return fmt.Errorf("Delete project was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if eventContext != nil && !SuppressWSCommunication {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint, *scheme == "https")
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
