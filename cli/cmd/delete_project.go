package cmd

import (
	"errors"
	"fmt"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

// crprojectCmd represents the project command
var delProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME",
	Short: "Deletes a project.",
	Long: `Deletes a new project with the provided name. 

Example:
	keptn delete project sockshop`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires PROJECTNAME")
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
		logging.PrintLog("Starting to delete project", logging.InfoLevel)

		project := apimodels.Project{
			Name: args[0],
		}

		if *createProjectParams.GitUser != "" && *createProjectParams.GitToken != "" && *createProjectParams.RemoteURL != "" {
			project.GitUser = *createProjectParams.GitUser
			project.GitToken = *createProjectParams.GitToken
			project.GitRemoteURL = *createProjectParams.RemoteURL
		}

		projectHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			response, err := projectHandler.DeleteProject(project)
			if err != nil {
				fmt.Println("Delete project was unsuccessful")
				return err
			}

			// check for response, which is of type apimodels.Error
			if response == nil {
				return nil
			}

			if response.Code > 299 {
				fmt.Sprintf("Delete project was unsuccessful. %s", *response.Message)
				return fmt.Errorf("Delete project was unsuccessfu. %s", *response.Message)
			}

			return fmt.Errorf("Received unexpected return code: %d", response.Code)
		} else {
			fmt.Println("Skipping delete project due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	deleteCmd.AddCommand(delProjectCmd)
}
