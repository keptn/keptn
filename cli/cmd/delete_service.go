package cmd

import (
	"errors"
	"fmt"

	"github.com/keptn/keptn/cli/internal"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type deleteServiceCmdParams struct {
	Project *string
	Service *string
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
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SERVICENAME not set")
		} else if len(args) >= 2 {
			cmd.SilenceUsage = false
			return errors.New("too many arguments set")
		}
		service := args[0]

		logging.PrintLog("Starting to delete service", logging.InfoLevel)

		api, err := internal.APIProvider(endPoint.String(), apiToken)
		if err != nil {
			return internal.OnAPIError(err)
		}

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			deleteResp, err := api.APIV1().DeleteService(*deleteServiceParams.Project, service)
			if err != nil {
				logging.PrintLog("Delete project was unsuccessful", logging.InfoLevel)
				return fmt.Errorf("Delete project was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if deleteResp != nil {
				logging.PrintLog("Service deleted successfully", logging.InfoLevel)
				logging.PrintLog(deleteResp.Message, logging.InfoLevel)
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
	deleteServiceParams.Project = delServiceCmd.Flags().StringP("project", "p", "", "The project from which to delete the service")
	delServiceCmd.MarkFlagRequired("project")
}
