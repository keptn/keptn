package cmd

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

type upgradeProjectCmdParams struct {
	Shipyard    *bool
	FromVersion *string
	ToVersion   *string
}

var upgradeProjectParams *upgradeProjectCmdParams

var supportedFromVersions = []string{"0.1", "0.1.0"}
var supportedToVersions = []string{"0.2", "0.2.0"}

const defaultFromVersion = "0.1"
const defaultToVersion = "0.2"

// upgradeProjectCmd represents the project command
var upgradeProjectCmd = &cobra.Command{
	Use:   "upgrade project PROJECTNAME --shipyard --fromVersion=CURRENT_SHIPYARD_VERSION --toVersion=TARGET_SHIPYARD_VERSION",
	Short: "Upgrades an existing Keptn project",
	Long: `Upgrades an existing Keptn project with the provided name. 

This command will upgrade the shipyard of the project to the specified version

By executing the update project command, Keptn will fetch the current shipyard.yaml file of the project and convert it to the version specified in the 'toVersion'' flag.

For more information about upgrading projects, go to [Manage Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/upgrade)
`,
	Example:      `keptn upgrade project PROJECTNAME --shipyard --fromVersion=0.1.0 --toVersion=0.2.0`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		_, _, err := credentialmanager.NewCredentialManager().GetCreds(namespace)
		if err != nil {
			return errors.New(authErrorMsg)
		}

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("required argument PROJECTNAME not set")
		}

		if upgradeProjectParams.FromVersion == nil || *upgradeProjectParams.FromVersion == "" {
			*upgradeProjectParams.ToVersion = defaultFromVersion
		} else if err := checkFromVersion(upgradeProjectParams.FromVersion); err != nil {
			return err
		}

		if upgradeProjectParams.ToVersion == nil || *upgradeProjectParams.ToVersion == "" {
			*upgradeProjectParams.ToVersion = defaultToVersion
		} else if err := checkToVersion(upgradeProjectParams.ToVersion); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var endPoint url.URL
		var apiToken string
		var err error
		if !mocking {
			endPoint, apiToken, err = credentialmanager.NewCredentialManager().GetCreds(namespace)
		} else {
			endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
			endPoint = *endPointPtr
			apiToken = ""
		}
		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		resourceHandler := apiutils.NewAuthenticatedResourceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		shipyardResource, err := resourceHandler.GetProjectResource(*newArtifact.Project, "shipyard.yaml")
		if err != nil {
			return fmt.Errorf("Error while retrieving shipyard.yaml for project %s: %s:", *newArtifact.Project, err.Error())
		}

		shipyard := &keptn.Shipyard{}

		if err := yaml.Unmarshal([]byte(shipyardResource.ResourceContent), shipyard); err != nil {
			return fmt.Errorf("Error while decoding shipyard.yaml for project %s: %s", *newArtifact.Project, err.Error())
		}

		return nil
	},
}

func checkFromVersion(fromVersion *string) error {
	for _, value := range supportedFromVersions {
		if value == *fromVersion {
			return nil
		}
	}
	return fmt.Errorf("invalid fromVersion %s. Please enter one of the following: %v", fromVersion, supportedFromVersions)
}

func checkToVersion(toVersion *string) error {
	for _, value := range supportedToVersions {
		if value == *toVersion {
			return nil
		}
	}
	return fmt.Errorf("invalid toVersion %s. Please enter one of the following: %v", toVersion, supportedToVersions)
}

func init() {
	rootCmd.AddCommand(upgradeProjectCmd)

	upgradeProjectParams = &upgradeProjectCmdParams{}

	upgradeProjectParams.Shipyard = upgradeProjectCmd.Flags().BoolP("shipyard", "", false, "Upgrade the shipyard file of the project")
	upgradeProjectParams.FromVersion = upgradeProjectCmd.Flags().StringP("fromVersion", "", "", "The current version of the shipyard")
	upgradeProjectParams.ToVersion = upgradeProjectCmd.Flags().StringP("toVersion", "", "", "The new target version of the shipyard")
	upgradeProjectCmd.MarkFlagRequired("shipyard")
}
