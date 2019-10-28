package cmd

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/utils/websockethelper"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type createServiceCmdParams struct {
	Project  *string
}

var createServiceParams *createServiceCmdParams

// crProjectCmd represents the project command
var crServiceCmd = &cobra.Command{
	Use:   "service SERVICENAME --project=PROJECTNAME",
	Short: "Creates a new service",
	Long: `Creates a new service with the provided name and in the specified project. 

Example:
	keptn create service carts --project=sockshop`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SERVICENAME not set")
		}

		if !utils.ValidateK8sName(args[0]) {
			errorMsg := "Service name contains upper case letter(s) or special character(s).\n"
			return errors.New(errorMsg)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to create service", logging.InfoLevel)

		service := apimodels.Service{
			ServiceName:     &args[0],
		}

		serviceHandler := apiutils.NewAuthenticatedServiceHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := serviceHandler.CreateService(*createServiceParams.Project, service)
			if err != nil {
				fmt.Println("Create service was unsuccessful")
				return fmt.Errorf("Create service was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if eventContext != nil {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint)
			}

			return nil
		}

		fmt.Println("Skipping create service due to mocking flag set to true")
		return nil
	},
}

func init() {
	createCmd.AddCommand(crServiceCmd)
	createServiceParams = &createServiceCmdParams{}
	createServiceParams.Project = crServiceCmd.Flags().StringP("project", "p", "", "The project in which to create the service")
	crServiceCmd.MarkFlagRequired("project")
}
