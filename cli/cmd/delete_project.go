package cmd

import (
	"errors"
	"fmt"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type deleteProjectCmdParams struct {
	KeepServices *bool
}

var deleteProjectParams *deleteProjectCmdParams

// delProjectCmd represents the project command
var delProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME",
	Short: "Deletes a project identified by project name",
	Long: `Deletes a project identified by project name. 

**Notes:**
* If a Git upstream is configured for this project, the referenced upstream repository (e.g., on GitHub) will not be deleted. 
* Services that have been deployed to the Kubernetes cluster are not deleted.
* Namespaces that have been created on the Kubernetes cluster are not deleted.
* Helm-releases created for deployments are not deleted. To clean-up deployed Helm releases, pelease see [Clean-up after deleting a project](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/continuous_delivery/deployment_helm/#clean-up-after-deleting-a-project)
`,
	Example:      `keptn delete project sockshop`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
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
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to delete project", logging.InfoLevel)

		project := apimodels.Project{
			ProjectName: args[0],
		}

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		projectsHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			if deleteProjectParams.KeepServices == nil || !*deleteProjectParams.KeepServices {
				apiProject, err := projectsHandler.GetProject(project)

				if err != nil {
					logging.PrintLog("Could not retrieve information about project "+project.ProjectName+": "+*err.Message, logging.InfoLevel)
					return fmt.Errorf("Could not retrieve information about project %s: %s", project.ProjectName, *err.Message)
				} else if apiProject == nil {
					msg := "Project " + project.ProjectName + " not found"
					logging.PrintLog(msg, logging.InfoLevel)
					return fmt.Errorf(msg)
				}

				if len(apiProject.Stages) > 0 {
					fmt.Println("Deleting services of project " + project.ProjectName + "...")
					for _, service := range apiProject.Stages[0].Services {
						logging.PrintLog("Deleting service "+service.ServiceName, logging.InfoLevel)
						deleteResp, err := apiHandler.DeleteService(project.ProjectName, service.ServiceName)
						if err != nil {
							logging.PrintLog("Delete service was unsuccessful", logging.InfoLevel)
							return fmt.Errorf("Delete service was unsuccessful. %s", *err.Message)
						}
						logging.PrintLog("Service deleted successfully", logging.InfoLevel)
						if len(deleteResp.Message) > 0 {
							logging.PrintLog(deleteResp.Message, logging.InfoLevel)
						}
					}
				}
			}

			deleteResp, err := apiHandler.DeleteProject(project)
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
	deleteProjectParams = &deleteProjectCmdParams{}
	deleteProjectParams.KeepServices = delProjectCmd.Flags().BoolP("keep-services", "", false, "Indicate whether the helm releases that are part of the project should be deleted as well, or not")
}
