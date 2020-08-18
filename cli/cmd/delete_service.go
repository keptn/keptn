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

// delServiceCmd allows to delete a service
var delServiceCmd = &cobra.Command{
	Use:   "delete SERVICENAME --project=PROJECTNAME",
	Short: "Deletes a service from a project",
	Long: `Deletes a service from a project by deleting the configuration in the GIT repository.
Furthermore, if Keptn is used for continuous delivery (i.e. services have been onboarded), this command will also uninstall the associated Helm releases.
`,
	Example:      `keptn delete project sockshop`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to delete project", logging.InfoLevel)

		project := apimodels.Project{
			ProjectName: args[0],
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			// TODO: API needs to be implemented first
			eventContext, err := apiHandler.DeleteProject(project)
			if err != nil {
				fmt.Println("Delete project was unsuccessful")
				return fmt.Errorf("Delete project was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if eventContext != nil && !SuppressWSCommunication {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint)
			}

			return nil
		}

		fmt.Println("Skipping delete service due to mocking flag set to true")
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(delServiceCmd)
}
