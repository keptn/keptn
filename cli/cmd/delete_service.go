package cmd

import (
	"errors"
	"fmt"

	"github.com/keptn/keptn/cli/pkg/websockethelper"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type deleteServiceCmdParams struct {
	Project *string
}

var deleteServiceParams *deleteServiceCmdParams

// delServiceCmd allows to delete a service
var delServiceCmd = &cobra.Command{
	Use:   "service SERVICENAME --project=PROJECTNAME",
	Short: "Deletes a service from a project",
	Long: `Deletes a service from a project by deleting the configuration in the GIT repository.
Furthermore, if Keptn is used for continuous delivery (i.e. services have been onboarded), this command will also uninstall the associated Helm releases.
`,
	Example:      `keptn delete service carts --project=sockshop`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to delete service", logging.InfoLevel)

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := apiHandler.DeleteService(*deleteServiceParams.Project, args[0])
			if err != nil {
				return fmt.Errorf("Delete service was unsuccessful: %s", *err.Message)
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
	deleteServiceParams = &deleteServiceCmdParams{}
	deleteServiceParams.Project = delServiceCmd.Flags().StringP("project", "p", "", "The project in which to create the service")
	delServiceCmd.MarkFlagRequired("project")
}
