package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/keptn/keptn/cli/internal"

	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/spf13/cobra"
)

type createProjectCmdParams struct {
	Shipyard          *string
	GitUser           *string
	GitToken          *string
	RemoteURL         *string
	GitPrivateKey     *string
	GitPrivateKeyPass *string
	GitProxyURL       *string
	GitProxyScheme    *string
	GitProxyUser      *string
	GitProxyPassword  *string
	GitProxySecure    *bool
	GitPublicCert     *string
}

var createProjectParams *createProjectCmdParams

const gitErrMsg = `Please specify a 'git-user', 'git-token', and 'git-remote-url' as flags for configuring a Git upstream repository`
const gitMissingUpstream = `WARNING: Creating a project without Git upstream repository is not recommended and will not be supported in the future anymore.
You can configure a Git upstream repository using: 

keptn update project PROJECTNAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL
`

// crProjectCmd represents the project command
var crProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --shipyard=FILEPATH",
	Short: "Creates a new project",
	Long: `Creates a new project with the provided name and Shipyard. 
The shipyard file describes the used stages. These stages are defined by name, as well as their task sequences.

By executing the *create project* command, Keptn initializes an internal Git repository that is used to maintain all project-related resources. 
To upstream this internal Git repository to a remote repository, the Git user (*--git-user*), an access token (*--git-token*), and the remote URL (*--git-remote-url*) are required.

For more information about Shipyard, creating projects, or upstream repositories, please go to [Manage Keptn](https://keptn.sh/docs/` + getReleaseDocsURL() + `/manage/)
`,
	Example: `keptn create project PROJECTNAME --shipyard=FILEPATH
keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument PROJECTNAME not set")
		} else if len(args) >= 2 {
			cmd.SilenceUsage = false
			return errors.New("too many arguments set")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		shipyard, err := retrieveShipyard(*createProjectParams.Shipyard)
		if err != nil {
			return fmt.Errorf("Failed to read and parse shipyard file - %s", err.Error())
		}

		if err := checkGitCredentials(); err != nil {
			return err
		}

		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}
		logging.PrintLog("Starting to create project", logging.InfoLevel)

		encodedShipyardContent := base64.StdEncoding.EncodeToString(shipyard)
		project := apimodels.CreateProject{
			Name:     &args[0],
			Shipyard: &encodedShipyardContent,
		}

		if *createProjectParams.GitUser != "" && *createProjectParams.RemoteURL != "" {
			project.GitUser = *createProjectParams.GitUser
			project.GitRemoteURL = *createProjectParams.RemoteURL
			project.GitToken = *createProjectParams.GitToken

			content, err := ioutil.ReadFile(*createProjectParams.GitPrivateKey)
			if err != nil {
				fmt.Errorf("unable to read privateKey file: %s\n", err.Error())
			}
			project.GitPrivateKey = string(content)
			project.GitPrivateKeyPass = *createProjectParams.GitPrivateKeyPass

			project.GitProxyURL = *createProjectParams.GitProxyURL
			project.GitProxyScheme = *createProjectParams.GitProxyScheme
			project.GitProxyUser = *createProjectParams.GitProxyUser
			project.GitProxyPassword = *createProjectParams.GitProxyPassword
			project.GitProxySecure = *createProjectParams.GitProxySecure

			content, err = ioutil.ReadFile(*createProjectParams.GitPublicCert)
			if err != nil {
				fmt.Errorf("unable to read TLS public certificate: %s\n", err.Error())
			}
			project.GitPublicCert = string(content)
		}

		api, err := internal.APIProvider(endPoint.String(), apiToken)
		if err != nil {
			return err
		}

		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		if !mocking {
			_, err := api.APIV1().CreateProject(project)
			if err != nil {
				return fmt.Errorf("Create project was unsuccessful.\n%s", *err.Message)
			}

			logging.PrintLog("Project created successfully", logging.InfoLevel)

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

	if *createProjectParams.GitUser != "" && *createProjectParams.RemoteURL != "" {
		return nil
	}
	return errors.New(gitErrMsg)
}

func retrieveShipyard(location string) ([]byte, error) {
	var content []byte
	var err error
	if httputils.IsValidURL(location) {
		content, err = httputils.NewDownloader(httputils.WithTimeout(5 * time.Second)).DownloadFromURL(location)
	} else {
		content, err = fileutils.ReadFile(location)
	}
	if err != nil {
		return nil, err
	}
	return content, nil
}

func init() {
	createCmd.AddCommand(crProjectCmd)
	createProjectParams = &createProjectCmdParams{}
	createProjectParams.Shipyard = crProjectCmd.Flags().StringP("shipyard", "s", "", "The path or URL to the shipyard file specifying the environment")
	crProjectCmd.MarkFlagRequired("shipyard")

	createProjectParams.GitUser = crProjectCmd.Flags().StringP("git-user", "u", "", "The git user of the upstream target")
	crProjectCmd.MarkFlagRequired("git-user")
	createProjectParams.RemoteURL = crProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")
	crProjectCmd.MarkFlagRequired("git-remote-url")
	createProjectParams.GitToken = crProjectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")

	createProjectParams.GitPrivateKey = crProjectCmd.Flags().StringP("git-private-key", "k", "", "The SSH git private key of the git user")
	createProjectParams.GitPrivateKeyPass = crProjectCmd.Flags().StringP("git-private-key-pass", "l", "", "The passphrase of git private key")

	createProjectParams.GitProxyURL = crProjectCmd.Flags().StringP("git-proxy-url", "p", "", "The git proxy URL and port")
	createProjectParams.GitProxyScheme = crProjectCmd.Flags().StringP("git-proxy-scheme", "j", "", "The git proxy scheme")
	createProjectParams.GitProxyUser = crProjectCmd.Flags().StringP("git-proxy-user", "w", "", "The git proxy user")
	createProjectParams.GitProxyPassword = crProjectCmd.Flags().StringP("git-proxy-password", "e", "", "The git proxy password")
	createProjectParams.GitProxySecure = crProjectCmd.Flags().BoolP("git-proxy-insecure", "e", true, "The git proxy secure TLS connection")

	createProjectParams.GitPublicCert = crProjectCmd.Flags().StringP("git-public-cert", "c", "", "The TLS git public certificate")
}
