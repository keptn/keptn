package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/keptn/keptn/cli/pkg/file"
	"github.com/keptn/keptn/cli/pkg/websockethelper"

	"github.com/asaskevich/govalidator"
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
const gitMissingUpstream = `WARNING: Creating a project without Git upstream repository is not recommended. 
You can configure a Git upstream repository using: 

keptn update project PROJECTNAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL
`

// crProjectCmd represents the project command
var crProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --shipyard=FILEPATH",
	Short: "Creates a new project",
	Long: `Creates a new project with the provided name and Shipyard. 
The shipyard file describes the used stages. These stages are defined by name, 
deployment-, test-, and remediation strategy.

By executing the *create project* command, Keptn initializes an internal Git repository that is used to maintain all project-related resources. 
To upstream this internal Git repository to a remote repository, the Git user (*--git-user*), an access token (*--git-token*), and the remote URL (*--git-remote-url*) are required.

For more information about Shipyard, creating projects, or upstream repositories, please go to [Manage Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/manage/)
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

		if !keptncommon.ValidateKeptnEntityName(args[0]) {
			errorMsg := "Project name contains upper case letter(s) or special character(s).\n"
			errorMsg += "Keptn relies on the following conventions: "
			errorMsg += "start with a lower case letter, then lower case letters, numbers, and hyphens are allowed.\n"
			errorMsg += "Please update project name and try again."
			return errors.New(errorMsg)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		shipyard := keptn.Shipyard{}
		err := getAndParseYaml(*createProjectParams.Shipyard, &shipyard)
		if err != nil {
			return fmt.Errorf("Failed to read and parse shipyard file - %s", err.Error())
		}

		// check stage names
		for _, stage := range shipyard.Stages {
			if !keptncommon.ValidateKeptnEntityName(stage.Name) {
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

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			eventContext, err := apiHandler.CreateProject(project)
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
		fmt.Println(gitMissingUpstream)
		return nil
	}

	if *createProjectParams.GitUser != "" && *createProjectParams.GitToken != "" && *createProjectParams.RemoteURL != "" {
		return nil
	}
	return errors.New(gitErrMsg)
}

func getAndParseYaml(arg string, out interface{}) error {
	var content string
	var err error
	if govalidator.IsURL(arg) {
		content, err = getYamlFromURL(arg)
	} else {
		content, err = getYamlFromFile(arg)
	}
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(content), out)
	if err != nil {
		return err
	}
	return nil
}

func getYamlFromURL(arg string) (string, error) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Get(arg)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getYamlFromFile(arg string) (string, error) {
	if _, err := os.Stat(keptnutils.ExpandTilde(arg)); os.IsNotExist(err) {
		return "", fmt.Errorf("Cannot find file %s", arg)
	}
	return file.ReadFile(arg)
}

func init() {
	createCmd.AddCommand(crProjectCmd)
	createProjectParams = &createProjectCmdParams{}
	createProjectParams.Shipyard = crProjectCmd.Flags().StringP("shipyard", "s", "", "The path or URL to the shipyard file specifying the environment")
	crProjectCmd.MarkFlagRequired("shipyard")

	createProjectParams.GitUser = crProjectCmd.Flags().StringP("git-user", "u", "", "The git user of the upstream target")
	createProjectParams.GitToken = crProjectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")
	createProjectParams.RemoteURL = crProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")
}
