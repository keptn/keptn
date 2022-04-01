package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/keptn/keptn/cli/internal"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type updateProjectCmdParams struct {
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

var updateProjectParams *updateProjectCmdParams

// upProjectCmd represents the project command
var upProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL",
	Short: "Updates an existing Keptn project",
	Long: `Updates an existing Keptn project with the provided name. 

Updating a shipyard file is not possible.

By executing the update project command, Keptn will add the provided upstream repository to the existing internal Git repository that is used to maintain all project-related resources. 
To upstream this internal Git repository to a remote repository and the remote URL (*--git-remote-url*) are required
together with private key (*--git-private-key*) or access token (*--git-token*). The Git user (--git-user) can be added if required. For using proxy please specify proxy IP address together with port (*--git-proxy-url*) and
used scheme (*--git-proxy-scheme=*) to connect to proxy. Please be aware that authentication with public/private key and via proxy is 
supported only when using resource-service.

For more information about updating projects or upstream repositories, please go to [Manage Keptn](https://keptn.sh/docs/` + getReleaseDocsURL() + `/manage/)
`,
	Example: `keptn update project PROJECTNAME --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL

or (if the git repo requires user)

keptn update project PROJECTNAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL

or (only for resource-service)

keptn update project PROJECTNAME --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-private-key=PRIVATE_KEY_PATH --git-private-key-pass=PRIVATE_KEY_PASSPHRASE

or (only for resource-service)

keptn update project PROJECTNAME --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN --git-proxy-url=PROXY_IP --git-proxy-scheme=SCHEME --git-proxy-user=PROXY_USER --git-proxy-password=PROXY_PASS --git-proxy-insecure`,
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

		project := apimodels.CreateProject{
			Name: &args[0],
		}

		if *updateProjectParams.GitUser != "" && *updateProjectParams.RemoteURL != "" {
			if *updateProjectParams.GitToken == "" && *updateProjectParams.GitPrivateKey == "" {
				return errors.New("Access token or private key must be set")
			}

			if *updateProjectParams.GitToken != "" && *updateProjectParams.GitPrivateKey != "" {
				return errors.New("Access token or private key cannot be set together")
			}

			project.GitUser = *updateProjectParams.GitUser
			project.GitToken = *updateProjectParams.GitToken
			project.GitRemoteURL = *updateProjectParams.RemoteURL

			if *updateProjectParams.GitProxyURL != "" && strings.HasPrefix(*updateProjectParams.RemoteURL, "ssh://") {
				return errors.New("Proxy cannot be set with SSH")
			}

			if *updateProjectParams.GitProxyURL != "" && *updateProjectParams.GitProxyScheme == "" {
				return errors.New("Proxy cannot be set without scheme")
			}

			project.GitProxyURL = *updateProjectParams.GitProxyURL
			project.GitProxyScheme = *updateProjectParams.GitProxyScheme
			project.GitProxyUser = *updateProjectParams.GitProxyUser
			project.GitProxyPassword = *updateProjectParams.GitProxyPassword
			project.GitProxyInsecure = *updateProjectParams.GitProxyInsecure

			if strings.HasPrefix(*updateProjectParams.RemoteURL, "ssh://") {
				content, err := ioutil.ReadFile(*updateProjectParams.GitPrivateKey)
				if err != nil {
					return fmt.Errorf("unable to read privateKey file: %s\n", err.Error())
				}

				project.GitPrivateKey = string(base64.StdEncoding.EncodeToString(content))
				project.GitPrivateKeyPass = *updateProjectParams.GitPrivateKeyPass
			}

			if *updateProjectParams.GitPemCertificate != "" {
				content, err := ioutil.ReadFile(*updateProjectParams.GitPemCertificate)
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
			_, err := api.APIV1().UpdateProject(project)
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
	upProjectCmd.MarkFlagRequired("git-remote-url")

	updateProjectParams.RemoteURL = upProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")

	updateProjectParams.GitPrivateKey = upProjectCmd.Flags().StringP("git-private-key", "k", "", "The SSH git private key of the git user")
	updateProjectParams.GitPrivateKeyPass = upProjectCmd.Flags().StringP("git-private-key-pass", "l", "", "The passphrase of git private key")

	updateProjectParams.GitProxyURL = upProjectCmd.Flags().StringP("git-proxy-url", "p", "", "The git proxy URL and port")
	updateProjectParams.GitProxyScheme = upProjectCmd.Flags().StringP("git-proxy-scheme", "j", "", "The git proxy scheme")
	updateProjectParams.GitProxyUser = upProjectCmd.Flags().StringP("git-proxy-user", "w", "", "The git proxy user")
	updateProjectParams.GitProxyPassword = upProjectCmd.Flags().StringP("git-proxy-password", "e", "", "The git proxy password")
	updateProjectParams.GitProxyInsecure = upProjectCmd.Flags().BoolP("git-proxy-insecure", "x", false, "The git proxy insecure TLS connection")

	updateProjectParams.GitPemCertificate = upProjectCmd.Flags().StringP("git-pem-certificate", "g", "", "The git PEM Certificate file")
}
