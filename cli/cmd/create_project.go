package cmd

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/internal/projectcreator"
	"time"

	"github.com/keptn/keptn/cli/internal"

	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"

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
	InsecureSkipTLS   *bool
}

var createProjectParams *createProjectCmdParams

// crProjectCmd represents the project command
var crProjectCmd = &cobra.Command{
	Use:   "project PROJECTNAME --shipyard=FILEPATH",
	Short: "Creates a new project",
	Long: `Creates a new project with the provided name and shipyard file.
The shipyard file describes the used stages. These stages are defined by name, as well as their task sequences.
If the Keptn installation is setup with the "automatic GIT repository provisioning" feature disabled, this command requires the user to pass GIT related information using
the "--git-remote-url", as well as the "--git-token" or "--git-private-key" flag.
If the GIT repository allows it, you can specify the a GIT user using the "--git-user" flag. For using a proxy please specify the proxy IP address together with port using the ""--git-proxy-url" flag and
the scheme using the "--git-proxy-scheme" flag. Please be aware that authentication with public/private key and via proxy is supported only when using resource-service.

For more information about Shipyard, creating projects, or upstream repositories, please go to [Manage Keptn](https://keptn.sh/docs/` + getReleaseDocsURL() + `/manage/)
`,
	Example: `keptn create project PROJECTNAME --shipyard=FILEPATH
keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN

or (only for resource-service)

keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-private-key=PRIVATE_KEY_PATH --git-private-key-pass=PRIVATE_KEY_PASSPHRASE

or (only for resource-service)

keptn create project PROJECTNAME --shipyard=FILEPATH --git-user=GIT_USER --git-remote-url=GIT_REMOTE_URL --git-token=GIT_TOKEN --git-proxy-url=PROXY_IP --git-proxy-scheme=SCHEME --git-proxy-user=PROXY_USER --git-proxy-password=PROXY_PASS --insecure-skip-tls
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
		logging.PrintLog("Starting to create project", logging.InfoLevel)

		// getting Keptn endpoint and token
		endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		// loading shipyard file from specified location
		shipyard, err := retrieveShipyard(*createProjectParams.Shipyard)
		if err != nil {
			return fmt.Errorf("Failed to read and parse shipyard file - %s", err.Error())
		}

		// getting API set
		api, err := internal.APIProvider(endPoint.String(), apiToken)
		if err != nil {
			return internal.OnAPIError(err)
		}

		// creating project using project creator
		projectCreator := projectcreator.New(api.APIV1(), func(string) ([]byte, error) { return shipyard, nil }, fileOpener{})
		projectInfo := projectcreator.ProjectInfo{
			Name:              args[0],
			Shipyard:          *createProjectParams.Shipyard,
			GitUser:           *createProjectParams.GitUser,
			GitToken:          *createProjectParams.GitToken,
			RemoteURL:         *createProjectParams.RemoteURL,
			GitPrivateKey:     *createProjectParams.GitPrivateKey,
			GitPrivateKeyPass: *createProjectParams.GitPrivateKeyPass,
			GitProxyURL:       *createProjectParams.GitProxyURL,
			GitProxyScheme:    *createProjectParams.GitProxyScheme,
			GitProxyUser:      *createProjectParams.GitProxyUser,
			GitProxyPassword:  *createProjectParams.GitProxyPassword,
			GitPemCertificate: *createProjectParams.GitPemCertificate,
			InsecureSkipTLS:   *createProjectParams.InsecureSkipTLS,
		}
		err = projectCreator.CreateProject(projectInfo)

		if err != nil {
			return err
		}
		return nil
	},
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
	createProjectParams.GitUser = crProjectCmd.Flags().StringP("git-user", "u", "", "The git user of the upstream target")
	createProjectParams.RemoteURL = crProjectCmd.Flags().StringP("git-remote-url", "r", "", "The remote url of the upstream target")
	createProjectParams.GitToken = crProjectCmd.Flags().StringP("git-token", "t", "", "The git token of the git user")
	createProjectParams.GitPrivateKey = crProjectCmd.Flags().StringP("git-private-key", "k", "", "The SSH git private key of the git user")
	createProjectParams.GitPrivateKeyPass = crProjectCmd.Flags().StringP("git-private-key-pass", "l", "", "The passphrase of git private key")
	createProjectParams.GitProxyURL = crProjectCmd.Flags().StringP("git-proxy-url", "p", "", "The git proxy URL and port")
	createProjectParams.GitProxyScheme = crProjectCmd.Flags().StringP("git-proxy-scheme", "j", "", "The git proxy scheme")
	createProjectParams.GitProxyUser = crProjectCmd.Flags().StringP("git-proxy-user", "w", "", "The git proxy user")
	createProjectParams.GitProxyPassword = crProjectCmd.Flags().StringP("git-proxy-password", "e", "", "The git proxy password")
	createProjectParams.InsecureSkipTLS = crProjectCmd.Flags().BoolP("insecure-skip-tls", "x", false, "Disable TLS verification to allow connections to servers using self signed certificates")
	createProjectParams.GitPemCertificate = crProjectCmd.Flags().StringP("git-pem-certificate", "g", "", "The git PEM Certificate file")
	crProjectCmd.MarkFlagRequired("shipyard")
}
