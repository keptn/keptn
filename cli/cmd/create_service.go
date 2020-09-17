package cmd

import (
	"errors"
	"fmt"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	"github.com/keptn/keptn/cli/pkg/websockethelper"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type createServiceCmdParams struct {
	Project *string
}

var createServiceParams *createServiceCmdParams

// crProjectCmd represents the project command
var crServiceCmd = &cobra.Command{
	Use:   "service SERVICENAME --project=PROJECTNAME",
	Short: "Creates a new service",
	Long: `Creates a new service with the provided name in the specified project.

**Note:** This command is different from keptn onboard service which requires a Helm chart.
`,
	Example:      `keptn create service carts --project=sockshop`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument SERVICENAME not set")
		}

		if !keptncommon.ValididateUnixDirectoryName(args[0]) {
			return errors.New("Service name contains special character(s)." +
				"The service name has to be a valid Unix directory name. For details see " +
				"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to create service", logging.InfoLevel)

		service := apimodels.CreateService{
			ServiceName: &args[0],
		}

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+`
Possible reasons:
* The Keptn API server is currently not available. Check if your Kubernetes cluster is available.
* Your Keptn CLI points to the wrong API server (verify using 'keptn status')`,
				endPointErr)
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := apiHandler.CreateService(*createServiceParams.Project, service)
			if err != nil {
				fmt.Println("Create service was unsuccessful")
				return fmt.Errorf("Create service was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if eventContext != nil && !SuppressWSCommunication {
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
