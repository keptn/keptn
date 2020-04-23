package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/keptn/keptn/cli/pkg/file"

	"github.com/keptn/keptn/cli/pkg/websockethelper"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type createProjectCmdParams struct {
	Shipyard  *string
	GitUser   *string
	GitToken  *string
	RemoteURL *string
}

var createProjectParams *createProjectCmdParams

const gitErrMsg = `Please specify a 'git-user', 'git-token', and 'git-remote-url' as flags for configuring a Git upstream repository`

// crProjectCmd represents the project command
var crProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --shipyard=FILEPATH",
	Short: "Creates a new project",
	Long: `Creates a new project with the provided name and shipyard file. 
The shipyard file describes the used stages. These stages are defined by name, 
deployment-, test-, and remediation strategy.

By executing the *create project* command, Keptn initializes an internal Git repository that is used to maintain all project-related resources. 
To upstream this internal Git repository to a remote repository, the Git user (*--git-user*), an access token (*--git-token*), and the remote URL (*--git-remote-url*) are required.

For more information about Shipyard files, creating projects or upstream repositories visit https://keptn.sh/docs/develop/manage/project/ .
`,
	Example: `keptn create project PROJECTNAME --shipyard=FILEPATH
keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL`,
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

		if !keptn.ValidateKeptnEntityName(args[0]) {
			errorMsg := "Project name contains upper case letter(s) or special character(s).\n"
			errorMsg += "Keptn relies on the following conventions: "
			errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
			errorMsg += "Please update project name and try again."
			return errors.New(errorMsg)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(keptnutils.ExpandTilde(*createProjectParams.Shipyard)); os.IsNotExist(err) {
			return fmt.Errorf("Cannot find shipyard file %s", *createProjectParams.Shipyard)
		}

		// validate shipyard file
		content, err := file.ReadFile(*createProjectParams.Shipyard)
		if err != nil {
			return err
		}

		shipyard, err := parseShipyard(content)
		if err != nil {
			return fmt.Errorf("Invalid shipyard file because parsing failed: %s", err.Error())
		}

		// check stage names
		for _, stage := range shipyard.Stages {
			if !keptn.ValidateKeptnEntityName(stage.Name) {
				errorMsg := "Stage " + stage.Name + " contains upper case letter(s) or special character(s).\n"
				errorMsg += "Keptn relies on the following conventions: "
				errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
				errorMsg += "Please update stage name in your shipyard and try again."
				return errors.New(errorMsg)
			}
		}

		return checkGitCredentials()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to create project", logging.InfoLevel)

		content, _ := file.ReadFile(*createProjectParams.Shipyard)
		shipyard := base64.StdEncoding.EncodeToString([]byte(content))
		project := apimodels.CreateProject{
			Name:     &args[0],
			Shipyard: &shipyard,
		}

		if *createProjectParams.GitUser != "" && *createProjectParams.GitToken != "" && *createProjectParams.RemoteURL != "" {
			project.GitUser = *createProjectParams.GitUser
			project.GitToken = *createProjectParams.GitToken
			project.GitRemoteURL = *createProjectParams.RemoteURL
		}

		projectHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := projectHandler.CreateProject(project)
			if err != nil {
				fmt.Println("Create project was unsuccessful")
				return fmt.Errorf("Create project was unsuccessful. %s", *err.Message)
			}

			// if eventContext is available, open WebSocket communication
			if eventContext != nil && !SuppressWSCommunication {
				return websockethelper.PrintWSContentEventContext(eventContext, endPoint)
			}

			return nil
		}

		fmt.Println("Skipping create project due to mocking flag set to true")
		return nil
	},
}

func checkGitCredentials() error {
	if *createProjectParams.GitUser == "" && *createProjectParams.GitToken == "" && *createProjectParams.RemoteURL == "" {
		return nil
	}

	if *createProjectParams.GitUser != "" && *createProjectParams.GitToken != "" && *createProjectParams.RemoteURL != "" {
		return nil
	}
	return errors.New(gitErrMsg)
}

func parseShipyard(shipyardContent string) (*keptn.Shipyard, error) {
	shipyard := keptn.Shipyard{}
	err := yaml.Unmarshal([]byte(shipyardContent), &shipyard)
	if err != nil {
		return nil, err
	}
	return &shipyard, nil
}

func init() {
	createCmd.AddCommand(crProjectCmd)
	createProjectParams = &createProjectCmdParams{}
	createProjectParams.Shipyard = crProjectCmd.Flags().StringP("shipyard", "s", "", "The shiypard file specifying the environment")
	crProjectCmd.MarkFlagRequired("shipyard")

	createProjectParams.GitUser = crProjectCmd.Flags().StringP("git-user", "u", "", "The git user of the upstream target")
	createProjectParams.GitToken = crProjectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")
	createProjectParams.RemoteURL = crProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")
}
