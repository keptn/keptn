package cmd

import (
	"errors"
	"fmt"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type updateProjectCmdParams struct {
	GitUser   *string
	GitToken  *string
	RemoteURL *string
}

var updateProjectParams *updateProjectCmdParams

// upProjectCmd represents the project command
var upProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL",
	Short: "Updates an existing Keptn project",
	Long: `Updates an existing Keptn project with the provided name. 

Updating a shipyard file is not possible.

By executing the update project command, Keptn will add the provided upstream repository to the existing internal Git repository that is used to maintain all project-related resources. 
To upstream this internal Git repository to a remote repository, the Git user (--git-user), an access token (--git-token), and the remote URL (--git-remote-url) are required.

For more information about updating projects or upstream repositories, please go to [Manage Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/manage/)
`,
	Example:      `keptn update project PROJECTNAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument PROJECTNAME not set")
		}

		if !keptncommon.ValidateKeptnEntityName(args[0]) {
			errorMsg := "Project name contains upper case letter(s) or special character(s).\n"
			errorMsg += "Keptn relies on the following conventions: "
			errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
			errorMsg += "Please update project name and try again."
			return errors.New(errorMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to update project", logging.InfoLevel)

		if endPointErr := CheckEndpointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		project := apimodels.CreateProject{
			Name: &args[0],
		}

		if *updateProjectParams.GitUser != "" && *updateProjectParams.GitToken != "" && *updateProjectParams.RemoteURL != "" {
			project.GitUser = *updateProjectParams.GitUser
			project.GitToken = *updateProjectParams.GitToken
			project.GitRemoteURL = *updateProjectParams.RemoteURL
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			_, err := apiHandler.UpdateProject(project)
			if err != nil {
				return fmt.Errorf("Update project was unsuccessful. %s", *err.Message)
			}

			logging.PrintLog("Project updated successfully", logging.InfoLevel)
			return nil
		}

		fmt.Println("Skipping update project due to mocking flag set to true")
		return nil
	},
}

func init() {
	updateCmd.AddCommand(upProjectCmd)
	updateProjectParams = &updateProjectCmdParams{}

	updateProjectParams.GitUser = upProjectCmd.Flags().StringP("git-user", "u", "", "The git user of the upstream target")
	updateProjectParams.GitToken = upProjectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")
	updateProjectParams.RemoteURL = upProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")
	upProjectCmd.MarkFlagRequired("git-user")
	upProjectCmd.MarkFlagRequired("git-token")
	upProjectCmd.MarkFlagRequired("git-remote-url")
}
