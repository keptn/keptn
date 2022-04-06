package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
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
	GitPemCertificate *string
	GitProxyInsecure  *bool
}

var createProjectParams *createProjectCmdParams

const gitErrMsg = `Please specify a 'git-user' and 'git-remote-url' as flags for configuring a Git upstream repository together with 'git-token' or 'git-private-key' depending on auth method. Please be aware that authentication with public/private key is supported only when using resource-service.`
const gitMissingUpstream = `WARNING: Creating a project without Git upstream repository is not recommended and will not be supported in the future anymore.
You can configure a Git upstream repository using: 

keptn update project PROJECTNAME --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN

or (if your repository provider allows to use user and password)

keptn update project PROJECTNAME --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN

or (only for resource-service)

keptn update project PROJECTNAME --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-private-key=PRIVATE_KEY_PATH --git-private-key-pass=PRIVATE_KEY_PASSPHRASE

or (only for resource-service)

keptn update project PROJECTNAME --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN --git-proxy-url=PROXY_IP --git-proxy-scheme=SCHEME --git-proxy-user=PROXY_USER --git-proxy-password=PROXY_PASS --git-proxy-insecure

Please be aware that authentication with public/private key and via proxy is supported only when using resource-service.
`

// crProjectCmd represents the project command
var crProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --shipyard=FILEPATH",
	Short: "Creates a new project",
	Long: `Creates a new project with the provided name and Shipyard. 
The shipyard file describes the used stages. These stages are defined by name, as well as their task sequences.

By executing the *create project* command, Keptn initializes an internal Git repository that is used to maintain all project-related resources. 
To upstream this internal Git repository to a remote repository, the remote URL (*--git-remote-url*) is required
together with private key (*--git-private-key*) or access token (*--git-token*). The Git user (*--git-user*) can be specified if the repository allows it. 
For using proxy please specify proxy IP address together with port (*--git-proxy-url*) and
used scheme (*--git-proxy-scheme=*) to connect to proxy. Please be aware that authentication with public/private key and via proxy is 
supported only when using resource-service.

For more information about Shipyard, creating projects, or upstream repositories, please go to [Manage Keptn](https://keptn.sh/docs/` + getReleaseDocsURL() + `/manage/)
`,
	Example: `keptn create project PROJECTNAME --shipyard=FILEPATH
keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN

or (only for resource-service)

keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-private-key=PRIVATE_KEY_PATH --git-private-key-pass=PRIVATE_KEY_PASSPHRASE

or (only for resource-service)

keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN --git-proxy-url=PROXY_IP --git-proxy-scheme=SCHEME --git-proxy-user=PROXY_USER --git-proxy-password=PROXY_PASS --git-proxy-insecure
`,
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

		if *createProjectParams.RemoteURL != "" {
			if *createProjectParams.GitToken == "" && *createProjectParams.GitPrivateKey == "" {
				return errors.New("Access token or private key must be set")
			}

			project.GitUser = *createProjectParams.GitUser
			project.GitToken = *createProjectParams.GitToken
			project.GitRemoteURL = *createProjectParams.RemoteURL

			if *createProjectParams.GitProxyURL != "" && strings.HasPrefix(*createProjectParams.RemoteURL, "ssh://") {
				return errors.New("Proxy cannot be set with SSH")
			}

			if *createProjectParams.GitProxyURL != "" && *createProjectParams.GitProxyScheme == "" {
				return errors.New("Proxy cannot be set without scheme")
			}

			project.GitProxyURL = *createProjectParams.GitProxyURL
			project.GitProxyScheme = *createProjectParams.GitProxyScheme
			project.GitProxyUser = *createProjectParams.GitProxyUser
			project.GitProxyPassword = *createProjectParams.GitProxyPassword
			project.GitProxyInsecure = *createProjectParams.GitProxyInsecure

			if strings.HasPrefix(*createProjectParams.RemoteURL, "ssh://") {
				content, err := ioutil.ReadFile(*createProjectParams.GitPrivateKey)
				if err != nil {
					return fmt.Errorf("unable to read privateKey file: %s\n", err.Error())
				}

				project.GitPrivateKey = string(base64.StdEncoding.EncodeToString(content))
				project.GitPrivateKeyPass = *createProjectParams.GitPrivateKeyPass
			}

			if *createProjectParams.GitPemCertificate != "" {
				content, err := ioutil.ReadFile(*createProjectParams.GitPemCertificate)
				if err != nil {
					return fmt.Errorf("unable to read PEM Certificate file: %s\n", err.Error())
				}

				project.GitPemCertificate = string(base64.StdEncoding.EncodeToString(content))
			}
		}

		api, err := internal.APIProvider(endPoint.String(), apiToken)
		if err != nil {
			return internal.OnAPIError(err)
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
	if *createProjectParams.GitToken == "" && *createProjectParams.RemoteURL == "" {
		fmt.Println(gitMissingUpstream)
		return nil
	}

	if *createProjectParams.GitToken != "" && *createProjectParams.GitPrivateKey != "" {
		return errors.New("Access token or private key cannot be set together")
	}

	if *createProjectParams.GitUser != "" && *createProjectParams.RemoteURL != "" {
		return nil
	}

	if *createProjectParams.GitToken != "" && *createProjectParams.RemoteURL == "" {
		return errors.New(gitErrMsg)
	}

	return nil
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
	createProjectParams.RemoteURL = crProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")

	createProjectParams.GitToken = crProjectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")

	createProjectParams.GitPrivateKey = crProjectCmd.Flags().StringP("git-private-key", "k", "", "The SSH git private key of the git user")
	createProjectParams.GitPrivateKeyPass = crProjectCmd.Flags().StringP("git-private-key-pass", "l", "", "The passphrase of git private key")

	createProjectParams.GitProxyURL = crProjectCmd.Flags().StringP("git-proxy-url", "p", "", "The git proxy URL and port")
	createProjectParams.GitProxyScheme = crProjectCmd.Flags().StringP("git-proxy-scheme", "j", "", "The git proxy scheme")
	createProjectParams.GitProxyUser = crProjectCmd.Flags().StringP("git-proxy-user", "w", "", "The git proxy user")
	createProjectParams.GitProxyPassword = crProjectCmd.Flags().StringP("git-proxy-password", "e", "", "The git proxy password")
	createProjectParams.GitProxyInsecure = crProjectCmd.Flags().BoolP("git-proxy-insecure", "x", false, "The git proxy insecure TLS connection")

	createProjectParams.GitPemCertificate = crProjectCmd.Flags().StringP("git-pem-certificate", "g", "", "The git PEM Certificate file")

}
